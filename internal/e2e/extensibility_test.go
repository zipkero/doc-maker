// Package e2e는 새 llms.txt 타입 사이트의 추가성(extensibility) e2e 통합 테스트를 담는다.
//
// 검증 목적(SPEC §5.8 / task-012):
//   - 수집·번역 본체 코드(internal/* 및 cmd/*)를 전혀 수정하지 않고,
//     새 사이트의 사이트 폴더(config.json + glossary.json)만 추가하면
//     그 사이트의 수집→번역 파이프라인 전 과정이 동작함을 보장한다.
//
// 테스트 구성:
//   - 실제 네트워크 대신 모킹 fetcher를 주입해 결정적(deterministic)으로 동작한다.
//   - TempDir 안에 사이트 폴더(config.json·glossary.json·llms.txt 응답)를 직접 구성하므로
//     기존 sites/ 디렉터리의 영구 변경이 필요 없다.
//   - "본체 코드 무수정"의 증거: 이 파일은 내부 API(collector.Collect,
//     translator.SelectTargets, translator.SaveTranslation, config.Load,
//     source.NewLLMsSourceWithFetcher)를 호출만 할 뿐, 수정하거나 확장하지 않는다.
package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"doc-maker/internal/collector"
	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
	"doc-maker/internal/source"
	"doc-maker/internal/translator"
)

// TestNewSiteExtensibility_E2E는 새 llms.txt 타입 사이트를
// 사이트 폴더(config.json·glossary.json)만 추가해(코드 변경 없이) 수집→번역할 수 있음을 검증한다.
//
// 시나리오:
//  1. TempDir에 "docsite" 사이트 폴더(docsite/config.json, docsite/glossary.json) 생성
//     — 프로덕션 sites/ 디렉터리를 건드리지 않는다.
//  2. 모킹 fetcher로 llms.txt 목록과 두 페이지 원문을 제공한다.
//  3. config.Load → source.NewLLMsSourceWithFetcher → collector.Collect 실행
//     — 원문이 sites/docsite/raw/ 아래 원본 경로 구조로 저장됨을 확인.
//  4. translator.SelectTargets로 번역 대상 선별
//     — 두 페이지 모두 미번역 상태이므로 전부 대상에 포함됨을 확인.
//  5. translator.SaveTranslation으로 번역문 저장
//     — 출력 경로(sites/docsite/output/)·구조·매니페스트 TranslatedHash 기록을 확인.
//  6. translator.SelectTargets 재실행
//     — 무변경이므로 대상 0건임을 확인(증분 재실행 성립).
func TestNewSiteExtensibility_E2E(t *testing.T) {
	tmp := t.TempDir()

	// ─── 1. 사이트 폴더 구성 ─────────────────────────────────────────────────
	// sites/docsite/ 폴더 하나가 곧 한 사이트다(D4).
	sitesRoot := filepath.Join(tmp, "sites")
	siteDir := filepath.Join(sitesRoot, "docsite")
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}

	// ─── 2. 새 사이트 config.json 생성 ──────────────────────────────────────
	// base_url, source_type, include/exclude 패턴만 담는다(D10 — output/glossary/manifest 경로는 규약).
	cfgData, err := json.Marshal(map[string]interface{}{
		"base_url":         "https://docs.example.com",
		"source_type":      "llms.txt",
		"include_patterns": []string{"docs/**"},
		"exclude_patterns": []string{},
	})
	if err != nil {
		t.Fatalf("설정 JSON 생성 실패: %v", err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "config.json"), cfgData, 0o644); err != nil {
		t.Fatalf("설정 파일 쓰기 실패: %v", err)
	}

	// ─── 3. 새 사이트 glossary.json 생성 ─────────────────────────────────────
	glossaryData, err := json.Marshal(map[string]string{
		"gateway": "게이트웨이",
		"plugin":  "플러그인",
		"route":   "라우트",
	})
	if err != nil {
		t.Fatalf("용어집 JSON 생성 실패: %v", err)
	}
	if err := os.WriteFile(filepath.Join(siteDir, "glossary.json"), glossaryData, 0o644); err != nil {
		t.Fatalf("용어집 파일 쓰기 실패: %v", err)
	}

	// ─── 4. 설정 파일 로드 (config.Load — 폴더명이 사이트 식별자) ──────────────
	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("config.Load 실패: %v", err)
	}
	if site.ID != "docsite" {
		t.Errorf("사이트 식별자: 기대 %q, 실제 %q", "docsite", site.ID)
	}
	if site.SourceType != config.SourceLLMsTxt {
		t.Errorf("SourceType: 기대 %q, 실제 %q", config.SourceLLMsTxt, site.SourceType)
	}

	// ─── 5. 모킹 fetcher 구성 ────────────────────────────────────────────────
	llmsTxtContent := strings.Join([]string{
		"# Example Site Documentation",
		"",
		"[Getting Started](https://docs.example.com/docs/getting-started)",
		"[Plugin Guide](https://docs.example.com/docs/plugin-guide)",
	}, "\n")

	page1Content := []byte("# Getting Started\n\nThis guide explains the gateway setup.")
	page2Content := []byte("# Plugin Guide\n\nLearn how to install and configure a plugin.")

	mockContents := map[string][]byte{
		"https://docs.example.com/llms.txt":             []byte(llmsTxtContent),
		"https://docs.example.com/docs/getting-started": page1Content,
		"https://docs.example.com/docs/plugin-guide":    page2Content,
	}

	mockFetch := func(u string) ([]byte, error) {
		if data, ok := mockContents[u]; ok {
			return data, nil
		}
		return nil, fmt.Errorf("모킹 fetcher: 미등록 URL — %s", u)
	}

	// ─── 6. source 생성 ──────────────────────────────────────────────────────
	src := source.NewLLMsSourceWithFetcher(site, mockFetch)

	// ─── 7. 수집 실행 (collector.Collect) ────────────────────────────────────
	collectResult, err := collector.Collect(site, src, mockFetch)
	if err != nil {
		t.Fatalf("collector.Collect 실패: %v", err)
	}

	// 수집 결과 검증: 두 페이지 모두 신규 저장
	if collectResult.Updated != 2 {
		t.Errorf("수집 Updated: 기대 2, 실제 %d", collectResult.Updated)
	}
	if collectResult.Skipped != 0 {
		t.Errorf("수집 Skipped: 기대 0, 실제 %d", collectResult.Skipped)
	}
	if collectResult.Failed != 0 {
		t.Errorf("수집 Failed: 기대 0, 실제 %d", collectResult.Failed)
	}

	// 원문 저장 경로가 sites/<siteID>/raw/<URL경로구조>를 반영하는지 확인
	page1RawPath := filepath.Join(siteDir, "raw", "docs", "getting-started.md")
	page2RawPath := filepath.Join(siteDir, "raw", "docs", "plugin-guide.md")

	for _, p := range []string{page1RawPath, page2RawPath} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("원문 파일이 없음: %s", p)
		}
	}

	// 원문 내용 확인
	saved1, err := os.ReadFile(page1RawPath)
	if err != nil {
		t.Fatalf("page1 원문 읽기 실패: %v", err)
	}
	if string(saved1) != string(page1Content) {
		t.Errorf("page1 원문 내용 불일치:\n 기대: %q\n 실제: %q", page1Content, saved1)
	}

	// 매니페스트 기록 확인(사이트 폴더 안 manifest.json)
	mfAfterCollect, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("수집 후 매니페스트 로드 실패: %v", err)
	}
	entry1, ok := mfAfterCollect.Get("https://docs.example.com/docs/getting-started")
	if !ok {
		t.Error("매니페스트에 page1 항목 없음")
	}
	if entry1.SourceHash == "" {
		t.Error("매니페스트 page1 SourceHash가 비어 있음")
	}
	if entry1.SourcePath == "" {
		t.Error("매니페스트 page1 SourcePath가 비어 있음")
	}
	// SourcePath는 사이트 폴더 기준 상대경로여야 한다
	if entry1.SourcePath != "raw/docs/getting-started.md" {
		t.Errorf("매니페스트 page1 SourcePath: 기대 %q, 실제 %q", "raw/docs/getting-started.md", entry1.SourcePath)
	}

	// ─── 8. 용어집 로드 확인 (LoadGlossary — 규약 경로) ─────────────────────
	glossary, err := translator.LoadGlossary(site.GlossaryPath())
	if err != nil {
		t.Fatalf("LoadGlossary 실패: %v", err)
	}
	if glossary["gateway"] != "게이트웨이" {
		t.Errorf("용어집 'gateway': 기대 %q, 실제 %q", "게이트웨이", glossary["gateway"])
	}
	if glossary["plugin"] != "플러그인" {
		t.Errorf("용어집 'plugin': 기대 %q, 실제 %q", "플러그인", glossary["plugin"])
	}

	// ─── 9. 번역 대상 선별 (translator.SelectTargets) ─────────────────────────
	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("SelectTargets 실패: %v", err)
	}
	if len(targets) != 2 {
		t.Errorf("번역 대상 수: 기대 2, 실제 %d", len(targets))
	}

	// ─── 10. 번역문 저장 (translator.SaveTranslation) ──────────────────────
	for _, tgt := range targets {
		translated := []byte("번역된 내용: " + tgt.PageURL + " (게이트웨이, 플러그인)")
		translatedPath := translator.TranslatedPath(site, tgt.PageURL)

		if err := translator.SaveTranslation(
			site, tgt.PageURL, translatedPath, translated, tgt.SourceHash,
		); err != nil {
			t.Fatalf("SaveTranslation 실패 (%s): %v", tgt.PageURL, err)
		}

		// 번역문 파일이 생성됐는지 확인
		if _, err := os.Stat(translatedPath); os.IsNotExist(err) {
			t.Errorf("번역문 파일이 없음: %s", translatedPath)
		}
	}

	// 번역문 경로가 sites/<siteID>/output/<URL경로구조>를 따르는지 확인
	page1TransPath := filepath.Join(siteDir, "output", "docs", "getting-started.md")
	page2TransPath := filepath.Join(siteDir, "output", "docs", "plugin-guide.md")

	for _, p := range []string{page1TransPath, page2TransPath} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("번역문 파일이 없음(경로 구조 검증): %s", p)
		}
	}

	// 매니페스트 TranslatedHash 기록 확인
	mfAfterTranslate, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("번역 후 매니페스트 로드 실패: %v", err)
	}
	entryAfter, ok := mfAfterTranslate.Get("https://docs.example.com/docs/getting-started")
	if !ok {
		t.Error("번역 후 매니페스트에 page1 항목 없음")
	}
	if entryAfter.TranslatedHash == "" {
		t.Error("번역 후 매니페스트 TranslatedHash가 비어 있음")
	}
	if entryAfter.TranslatedHash != entryAfter.SourceHash {
		t.Errorf("TranslatedHash가 SourceHash와 다름: translated=%q source=%q",
			entryAfter.TranslatedHash, entryAfter.SourceHash)
	}

	// ─── 11. 번역 증분 재실행 확인 ────────────────────────────────────────────
	targets2, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("2회차 SelectTargets 실패: %v", err)
	}
	if len(targets2) != 0 {
		t.Errorf("무변경 2회차 번역 대상: 기대 0, 실제 %d (재번역 불필요 페이지 포함됨)", len(targets2))
	}
}

// TestNewSiteExtensibility_CodeChangeNotRequired는 새 사이트 추가에
// 본체 코드 변경이 필요 없음을 구조적으로 설명하는 주석 전용 테스트다.
func TestNewSiteExtensibility_CodeChangeNotRequired(t *testing.T) {
	t.Log("새 llms.txt 타입 사이트 추가에 본체 코드 변경이 불필요함은 " +
		"TestNewSiteExtensibility_E2E의 구성에서 확인된다: " +
		"사이트 폴더(config.json·glossary.json)만 TempDir에 추가해 수집→번역 전 과정을 성공시킨다.")
}

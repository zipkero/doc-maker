// translate CLI 테스트: 대상출력→번역문 생성→완료기록→대상출력(빠짐) 흐름을 검증한다(task-013).
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
)

// writeTempSite는 임시 사이트 폴더를 구성하고 Site를 반환한다.
//
// 구성:
//   - sites/<siteID>/config.json (ollama 형식 최소 설정)
//   - sites/<siteID>/raw/<relPath> = content (원문 파일)
//   - sites/<siteID>/manifest.json (Entry: SourceHash 기록, TranslatedHash 없음)
func writeTempSite(t *testing.T, siteID string, pages map[string][]byte) (*config.Site, string) {
	t.Helper()

	sitesRoot := t.TempDir()
	siteDir := filepath.Join(sitesRoot, siteID)

	// config.json 최소 설정
	cfgDir := siteDir
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}
	cfgJSON := `{"base_url":"https://example.com","source_type":"llms.txt"}`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(cfgJSON), 0o644); err != nil {
		t.Fatalf("config.json 쓰기 실패: %v", err)
	}

	// Site 로드
	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("설정 로드 실패: %v", err)
	}

	// 원문 파일과 매니페스트 초기 구성
	mf := manifest.New()
	for pageURL, content := range pages {
		// pageURL에서 파일 경로 추출: "https://example.com/path/to/page" → "path/to/page.md"
		relPath := strings.TrimPrefix(pageURL, "https://example.com/") + ".md"
		absPath := filepath.Join(site.RawDir(), filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			t.Fatalf("raw 디렉터리 생성 실패: %v", err)
		}
		if err := os.WriteFile(absPath, content, 0o644); err != nil {
			t.Fatalf("원문 파일 쓰기 실패: %v", err)
		}

		// 매니페스트에 Entry 등록(사이트 폴더 기준 상대경로)
		rel, err := filepath.Rel(siteDir, absPath)
		if err != nil {
			t.Fatalf("상대경로 변환 실패: %v", err)
		}
		mf.Set(pageURL, manifest.Entry{
			SourceHash: manifest.HashContent(content),
			SourcePath: filepath.ToSlash(rel),
		})
	}
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	return site, sitesRoot
}

// captureOutput은 os.Pipe를 통해 runPlan/runCommit의 출력을 문자열로 캡처한다.
func captureOutput(t *testing.T, fn func(w *os.File) error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe 생성 실패: %v", err)
	}
	defer r.Close()

	runErr := fn(w)
	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("출력 읽기 실패: %v", err)
	}
	return buf.String(), runErr
}

// TestPlanShowsUntranslated는 plan이 미번역 대상을 나열하는지 확인한다.
func TestPlanShowsUntranslated(t *testing.T) {
	pages := map[string][]byte{
		"https://example.com/docs/intro": []byte("# Introduction"),
		"https://example.com/docs/guide": []byte("# Guide"),
	}
	_, sitesRoot := writeTempSite(t, "testsite", pages)

	out, err := captureOutput(t, func(w *os.File) error {
		return runPlan("testsite", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("runPlan 오류: %v", err)
	}

	// 대상 2건이 나열되어야 한다.
	if !strings.Contains(out, "2 페이지") {
		t.Errorf("plan 출력에 '2 페이지'가 없음:\n%s", out)
	}
	if !strings.Contains(out, "https://example.com/docs/intro") {
		t.Errorf("plan 출력에 intro URL이 없음:\n%s", out)
	}
	if !strings.Contains(out, "https://example.com/docs/guide") {
		t.Errorf("plan 출력에 guide URL이 없음:\n%s", out)
	}
}

// TestPlanExcludesTranslated는 이미 번역된 페이지가 plan 출력에서 제외되는지 확인한다.
func TestPlanExcludesTranslated(t *testing.T) {
	pages := map[string][]byte{
		"https://example.com/page": []byte("content"),
	}
	site, sitesRoot := writeTempSite(t, "site2", pages)

	// 미리 번역 완료 상태로 만든다: TranslatedHash == SourceHash
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("매니페스트 로드 실패: %v", err)
	}
	entry, _ := mf.Get("https://example.com/page")
	entry.TranslatedHash = entry.SourceHash
	mf.Set("https://example.com/page", entry)
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	out, err := captureOutput(t, func(w *os.File) error {
		return runPlan("site2", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("runPlan 오류: %v", err)
	}

	if strings.Contains(out, "https://example.com/page") {
		t.Errorf("이미 번역된 페이지가 plan에 포함되어서는 안 됨:\n%s", out)
	}
	if !strings.Contains(out, "번역할 페이지가 없습니다") {
		t.Errorf("plan 출력에 '번역할 페이지가 없습니다'가 없음:\n%s", out)
	}
}

// TestCommitFlow는 대상출력→번역문 생성(output/ 파일 작성)→완료기록→대상출력(빠짐) 흐름을
// 검증하는 핵심 통합 테스트다(task-013 검증 조건).
func TestCommitFlow(t *testing.T) {
	pages := map[string][]byte{
		"https://example.com/docs/intro": []byte("# Introduction"),
	}
	site, sitesRoot := writeTempSite(t, "flowsite", pages)

	// 1단계: plan — 1건의 대상이 나열되어야 한다.
	out1, err := captureOutput(t, func(w *os.File) error {
		return runPlan("flowsite", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("1단계 runPlan 오류: %v", err)
	}
	if !strings.Contains(out1, "1 페이지") {
		t.Errorf("1단계 plan: '1 페이지' 없음:\n%s", out1)
	}

	// 2단계: 번역문을 output/에 직접 작성한다(Claude가 하는 역할을 시뮬레이션).
	outPath := site.OutputDir()
	translatedFile := filepath.Join(outPath, "docs", "intro.md")
	if err := os.MkdirAll(filepath.Dir(translatedFile), 0o755); err != nil {
		t.Fatalf("output 디렉터리 생성 실패: %v", err)
	}
	if err := os.WriteFile(translatedFile, []byte("# 소개\n\n한국어 번역 내용"), 0o644); err != nil {
		t.Fatalf("번역문 파일 쓰기 실패: %v", err)
	}

	// 3단계: commit — 1건이 기록되어야 한다.
	out2, err := captureOutput(t, func(w *os.File) error {
		return runCommit("flowsite", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("3단계 runCommit 오류: %v", err)
	}
	if !strings.Contains(out2, "1건 기록") {
		t.Errorf("commit 출력에 '1건 기록' 없음:\n%s", out2)
	}

	// 매니페스트에 TranslatedHash가 기록되었는지 직접 확인한다.
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("매니페스트 로드 실패: %v", err)
	}
	entry, ok := mf.Get("https://example.com/docs/intro")
	if !ok {
		t.Fatal("매니페스트에 intro 항목이 없음")
	}
	if entry.TranslatedHash == "" {
		t.Error("TranslatedHash가 기록되지 않음")
	}
	if entry.TranslatedHash != entry.SourceHash {
		t.Errorf("TranslatedHash(%q) != SourceHash(%q)", entry.TranslatedHash, entry.SourceHash)
	}

	// 4단계: plan — 이제 대상이 없어야 한다(증분이 닫힘).
	out3, err := captureOutput(t, func(w *os.File) error {
		return runPlan("flowsite", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("4단계 runPlan 오류: %v", err)
	}
	if !strings.Contains(out3, "번역할 페이지가 없습니다") {
		t.Errorf("commit 후 plan에 여전히 대상이 있음:\n%s", out3)
	}
}

// TestCommitSkipsWithoutOutputFile은 output/ 파일이 없는 페이지는 건너뛰는지 확인한다.
func TestCommitSkipsWithoutOutputFile(t *testing.T) {
	pages := map[string][]byte{
		"https://example.com/page": []byte("content"),
	}
	site, sitesRoot := writeTempSite(t, "skipsite", pages)

	// output/ 파일 없이 commit 실행
	out, err := captureOutput(t, func(w *os.File) error {
		return runCommit("skipsite", sitesRoot, w)
	})
	if err != nil {
		t.Fatalf("runCommit 오류: %v", err)
	}

	if !strings.Contains(out, "0건 기록") {
		t.Errorf("출력 파일 없을 때 기록 건수가 0이어야 함:\n%s", out)
	}

	// 매니페스트의 TranslatedHash가 여전히 비어 있어야 한다.
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("매니페스트 로드 실패: %v", err)
	}
	entry, ok := mf.Get("https://example.com/page")
	if !ok {
		t.Fatal("매니페스트에 page 항목이 없음")
	}
	if entry.TranslatedHash != "" {
		t.Errorf("번역문 없는 페이지의 TranslatedHash가 기록됨: %q", entry.TranslatedHash)
	}
}

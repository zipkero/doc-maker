package translator_test

// TestIncrementalTranslation_RoundTrip은 SelectTargets+SaveTranslation 조합이
// 재실행 시나리오에서 올바르게 동작하는지 검증하는 통합 테스트다(task-011 검증 조건).
//
// 시나리오:
//  1. 1회차: 4개 페이지 모두 미번역 → SelectTargets가 4개 반환 → 모두 SaveTranslation
//  2. 2회차(무변경): 매니페스트 그대로 → SelectTargets가 0개 반환 (새로 번역되는 페이지 없음)
//  3. 3회차(일부 원문 변경): pageB·pageD 원문 변경 → SelectTargets가 2개만 반환
//
// 번역 텍스트는 임의 문자열로 대체한다. 번역 품질이 아니라 증분 선별·스킵이 검증 대상이다.

import (
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/manifest"
	"doc-maker/internal/translator"
)

func TestIncrementalTranslation_RoundTrip(t *testing.T) {
	site := testSite(t, "incremental-site")
	siteRawDir := site.RawDir()

	// 원문 초기 콘텐츠 준비
	pages := []struct {
		relPath string
		pageURL string
		content []byte
	}{
		{"docs/page-a.md", "https://example.com/docs/page-a", []byte("page A original")},
		{"docs/page-b.md", "https://example.com/docs/page-b", []byte("page B original")},
		{"section/page-c.md", "https://example.com/section/page-c", []byte("page C original")},
		{"section/page-d.md", "https://example.com/section/page-d", []byte("page D original")},
	}

	// 원문 파일 생성 및 초기 매니페스트 구성
	mf := manifest.New()
	for _, p := range pages {
		absPath := writeRawFile(t, siteRawDir, p.relPath, p.content)
		hash := manifest.HashContent(p.content)
		mf.Set(p.pageURL, manifest.Entry{
			SourceHash: hash,
			SourcePath: toSiteRel(t, site.SiteDir, absPath),
		})
	}
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("초기 매니페스트 저장 실패: %v", err)
	}

	// ── 1회차: 전체 미번역 ──────────────────────────────────────────────────────
	targets1, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("1회차 SelectTargets 오류: %v", err)
	}
	if len(targets1) != 4 {
		t.Errorf("1회차 대상 수: 기대 4, 실제 %d", len(targets1))
	}

	// 선별된 모든 페이지를 번역 완료로 기록 (번역 텍스트는 임의 문자열)
	for _, tgt := range targets1 {
		translatedContent := []byte("번역: " + tgt.PageURL)
		translatedPath := translator.TranslatedPath(site, tgt.PageURL)
		if err := translator.SaveTranslation(site, tgt.PageURL, translatedPath, translatedContent, tgt.SourceHash); err != nil {
			t.Fatalf("1회차 SaveTranslation 오류 (%s): %v", tgt.PageURL, err)
		}
	}

	// ── 2회차(무변경): SelectTargets 대상 0건 기대 ─────────────────────────────
	targets2, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("2회차 SelectTargets 오류: %v", err)
	}
	if len(targets2) != 0 {
		t.Errorf("무변경 2회차 대상 수: 기대 0, 실제 %d (재번역되어서는 안 되는 페이지 포함됨)", len(targets2))
		for _, tgt := range targets2 {
			t.Logf("  불필요하게 포함된 대상: %s", tgt.PageURL)
		}
	}

	// ── 3회차: pageB·pageD 원문 변경 후 재실행 ──────────────────────────────────
	// pageB 원문 변경
	pageBPath := filepath.Join(siteRawDir, "docs", "page-b.md")
	newContentB := []byte("page B UPDATED content")
	if err := os.WriteFile(pageBPath, newContentB, 0o644); err != nil {
		t.Fatalf("pageB 원문 변경 실패: %v", err)
	}
	// pageD 원문 변경
	pageDPath := filepath.Join(siteRawDir, "section", "page-d.md")
	newContentD := []byte("page D UPDATED content")
	if err := os.WriteFile(pageDPath, newContentD, 0o644); err != nil {
		t.Fatalf("pageD 원문 변경 실패: %v", err)
	}
	// 수집 단계가 매니페스트 SourceHash를 갱신하는 역할을 담당하므로,
	// 여기서는 변경된 원문의 SourceHash를 매니페스트에 반영한다
	// (collector가 수집 후 SourceHash를 갱신하는 것과 동일한 상태를 만든다).
	mf3, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("3회차 사전 매니페스트 로드 실패: %v", err)
	}
	for _, p := range []struct {
		pageURL    string
		newContent []byte
		newPath    string
	}{
		{"https://example.com/docs/page-b", newContentB, pageBPath},
		{"https://example.com/section/page-d", newContentD, pageDPath},
	} {
		entry, ok := mf3.Get(p.pageURL)
		if !ok {
			t.Fatalf("3회차 사전 준비: 매니페스트에 %s 없음", p.pageURL)
		}
		entry.SourceHash = manifest.HashContent(p.newContent)
		entry.SourcePath = toSiteRel(t, site.SiteDir, p.newPath)
		mf3.Set(p.pageURL, entry)
	}
	if err := manifest.Save(site.ManifestDir(), mf3); err != nil {
		t.Fatalf("3회차 사전 매니페스트 저장 실패: %v", err)
	}

	targets3, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("3회차 SelectTargets 오류: %v", err)
	}
	if len(targets3) != 2 {
		t.Errorf("원문 변경 2건 후 대상 수: 기대 2, 실제 %d", len(targets3))
		for _, tgt := range targets3 {
			t.Logf("  3회차 포함된 대상: %s", tgt.PageURL)
		}
	}

	// pageA·pageC가 3회차 대상에 포함되지 않는지 확인한다
	byURL3 := make(map[string]translator.TranslationTarget)
	for _, tgt := range targets3 {
		byURL3[tgt.PageURL] = tgt
	}
	for _, unchanged := range []string{
		"https://example.com/docs/page-a",
		"https://example.com/section/page-c",
	} {
		if _, ok := byURL3[unchanged]; ok {
			t.Errorf("3회차: 원문 무변경 페이지가 대상에 포함됨 — %s", unchanged)
		}
	}

	// pageB·pageD가 3회차 대상에 포함되는지 확인한다
	for _, changed := range []string{
		"https://example.com/docs/page-b",
		"https://example.com/section/page-d",
	} {
		if _, ok := byURL3[changed]; !ok {
			t.Errorf("3회차: 원문 변경 페이지가 대상에서 빠짐 — %s", changed)
		}
	}

	// 3회차 변경분도 번역 완료 처리 후 다시 SelectTargets하면 0건이어야 한다
	for _, tgt := range targets3 {
		translatedContent := []byte("재번역: " + tgt.PageURL)
		translatedPath := translator.TranslatedPath(site, tgt.PageURL)
		if err := translator.SaveTranslation(site, tgt.PageURL, translatedPath, translatedContent, tgt.SourceHash); err != nil {
			t.Fatalf("3회차 SaveTranslation 오류 (%s): %v", tgt.PageURL, err)
		}
	}
	targets3b, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("3회차 완료 후 SelectTargets 오류: %v", err)
	}
	if len(targets3b) != 0 {
		t.Errorf("3회차 번역 완료 후 대상 수: 기대 0, 실제 %d", len(targets3b))
	}
}

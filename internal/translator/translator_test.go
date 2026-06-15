package translator_test

import (
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
	"doc-maker/internal/translator"
)

// 테스트 헬퍼: siteRawDir/<relPath>에 content를 담은 파일을 생성하고 절대 경로를 반환한다.
func writeRawFile(t *testing.T, siteRawDir, relPath string, content []byte) string {
	t.Helper()
	abs := filepath.Join(siteRawDir, relPath)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatalf("디렉터리 생성 실패: %v", err)
	}
	if err := os.WriteFile(abs, content, 0o644); err != nil {
		t.Fatalf("파일 쓰기 실패: %v", err)
	}
	return abs
}

// testSite는 임시 siteDir을 가리키는 테스트용 Site를 반환한다.
func testSite(t *testing.T, id string) *config.Site {
	t.Helper()
	siteDir := filepath.Join(t.TempDir(), "sites", id)
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}
	return &config.Site{
		ID:      id,
		SiteDir: siteDir,
	}
}

// TestSelectTargets_MixedStates는 미번역·무변경·변경 세 상태 페이지가 섞인 입력에서
// 선별 결과가 기대와 일치하는지 확인한다(task-009 핵심 검증 조건).
//
// 시나리오:
//   - pageA: 미번역 (TranslatedHash == "")               → 대상 포함
//   - pageB: 이미 번역된 무변경 (TranslatedHash == SourceHash)  → 제외
//   - pageC: 이미 번역됐으나 원문 변경 (TranslatedHash != 현재 해시) → 대상 포함
func TestSelectTargets_MixedStates(t *testing.T) {
	site := testSite(t, "testsite")
	siteRawDir := site.RawDir()

	// 원문 콘텐츠 준비
	contentA := []byte("page A original content")
	contentB := []byte("page B original content")
	contentC_old := []byte("page C old content")
	contentC_new := []byte("page C NEW content")

	hashA := manifest.HashContent(contentA)
	hashB := manifest.HashContent(contentB)
	hashC_old := manifest.HashContent(contentC_old)
	hashC_new := manifest.HashContent(contentC_new)

	// 원문 파일 생성
	pathA := writeRawFile(t, siteRawDir, "docs/page-a.md", contentA)
	pathB := writeRawFile(t, siteRawDir, "docs/page-b.md", contentB)
	// pageC는 원문이 변경된 상태로 디스크에 기록한다.
	pathC := writeRawFile(t, siteRawDir, "docs/page-c.md", contentC_new)

	// 매니페스트의 SourcePath는 사이트 폴더 기준 상대경로로 기록
	relA := toSiteRel(t, site.SiteDir, pathA)
	relB := toSiteRel(t, site.SiteDir, pathB)
	relC := toSiteRel(t, site.SiteDir, pathC)

	// 매니페스트 구성
	mf := manifest.New()

	// pageA: 수집됐으나 번역 미수행 (TranslatedHash == "")
	mf.Set("https://example.com/docs/page-a", manifest.Entry{
		SourceHash: hashA,
		SourcePath: relA,
	})

	// pageB: 수집·번역 모두 완료, 원문 무변경
	mf.Set("https://example.com/docs/page-b", manifest.Entry{
		SourceHash:     hashB,
		SourcePath:     relB,
		TranslatedHash: hashB,
	})

	// pageC: 번역된 적 있으나 원문이 새로 변경됨
	mf.Set("https://example.com/docs/page-c", manifest.Entry{
		SourceHash:     hashC_new,
		SourcePath:     relC,
		TranslatedHash: hashC_old,
	})

	// 매니페스트 저장
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	// 선별 실행
	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("SelectTargets 오류: %v", err)
	}

	// pageB는 제외되어야 하고, pageA·pageC는 포함되어야 한다.
	if len(targets) != 2 {
		t.Errorf("대상 페이지 수: 기대 2, 실제 %d", len(targets))
		for _, tgt := range targets {
			t.Logf("  포함된 대상: %s", tgt.PageURL)
		}
	}

	// URL별로 인덱싱해 개별 확인한다.
	byURL := make(map[string]translator.TranslationTarget)
	for _, tgt := range targets {
		byURL[tgt.PageURL] = tgt
	}

	// pageA: 미번역 → 포함
	tgtA, ok := byURL["https://example.com/docs/page-a"]
	if !ok {
		t.Error("pageA(미번역)가 선별 대상에 없음")
	} else {
		if tgtA.SourceHash != hashA {
			t.Errorf("pageA SourceHash: 기대 %q, 실제 %q", hashA, tgtA.SourceHash)
		}
		if tgtA.LocalPath != pathA {
			t.Errorf("pageA LocalPath: 기대 %q, 실제 %q", pathA, tgtA.LocalPath)
		}
	}

	// pageB: 무변경 번역본 → 제외
	if _, ok := byURL["https://example.com/docs/page-b"]; ok {
		t.Error("pageB(무변경 번역)가 선별 대상에 포함되어서는 안 됨")
	}

	// pageC: 원문 변경 → 포함
	tgtC, ok := byURL["https://example.com/docs/page-c"]
	if !ok {
		t.Error("pageC(원문 변경)가 선별 대상에 없음")
	} else {
		if tgtC.SourceHash != hashC_new {
			t.Errorf("pageC SourceHash: 기대 %q, 실제 %q", hashC_new, tgtC.SourceHash)
		}
	}
}

// TestSelectTargets_AllUntranslated는 모든 페이지가 미번역인 경우 전체가 대상에 포함되는지 확인한다.
func TestSelectTargets_AllUntranslated(t *testing.T) {
	site := testSite(t, "site1")
	siteRawDir := site.RawDir()

	contents := [][]byte{
		[]byte("first page"),
		[]byte("second page"),
	}
	paths := []string{
		writeRawFile(t, siteRawDir, "p1.md", contents[0]),
		writeRawFile(t, siteRawDir, "p2.md", contents[1]),
	}
	urls := []string{
		"https://site1.com/p1",
		"https://site1.com/p2",
	}

	mf := manifest.New()
	for i, u := range urls {
		mf.Set(u, manifest.Entry{
			SourceHash: manifest.HashContent(contents[i]),
			SourcePath: toSiteRel(t, site.SiteDir, paths[i]),
		})
	}
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("SelectTargets 오류: %v", err)
	}

	if len(targets) != 2 {
		t.Errorf("전체 미번역 시 대상 수: 기대 2, 실제 %d", len(targets))
	}
}

// TestSelectTargets_AllTranslatedUnchanged는 모든 페이지가 이미 번역된 무변경인 경우
// 대상 목록이 비어야 함을 확인한다.
func TestSelectTargets_AllTranslatedUnchanged(t *testing.T) {
	site := testSite(t, "site2")
	siteRawDir := site.RawDir()

	content := []byte("already translated content")
	hash := manifest.HashContent(content)
	path := writeRawFile(t, siteRawDir, "page.md", content)

	mf := manifest.New()
	mf.Set("https://site2.com/page", manifest.Entry{
		SourceHash:     hash,
		SourcePath:     toSiteRel(t, site.SiteDir, path),
		TranslatedHash: hash,
	})
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("SelectTargets 오류: %v", err)
	}

	if len(targets) != 0 {
		t.Errorf("모두 번역된 무변경 시 대상 수: 기대 0, 실제 %d", len(targets))
	}
}

// TestSelectTargets_NoRawDir는 원문 디렉터리가 없으면 빈 목록을 반환하는지 확인한다(수집 전 상태).
func TestSelectTargets_NoRawDir(t *testing.T) {
	site := testSite(t, "missing-site")
	// raw/ 디렉터리를 생성하지 않은 상태

	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("원문 디렉터리 없을 때 오류 반환(기대: nil): %v", err)
	}
	if len(targets) != 0 {
		t.Errorf("원문 디렉터리 없을 때 대상 수: 기대 0, 실제 %d", len(targets))
	}
}

// TestSelectTargets_TargetFields는 선별된 항목의 필드(PageURL, LocalPath, SourceHash)가
// 모두 채워져 있는지 확인한다(task-010이 사용하는 계약 확인).
func TestSelectTargets_TargetFields(t *testing.T) {
	site := testSite(t, "fieldcheck")
	siteRawDir := site.RawDir()

	content := []byte("some markdown content")
	hash := manifest.HashContent(content)
	path := writeRawFile(t, siteRawDir, "doc/guide.md", content)
	pageURL := "https://example.com/doc/guide"

	mf := manifest.New()
	mf.Set(pageURL, manifest.Entry{
		SourceHash: hash,
		SourcePath: toSiteRel(t, site.SiteDir, path),
		// TranslatedHash 없음 → 미번역
	})
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("매니페스트 저장 실패: %v", err)
	}

	targets, err := translator.SelectTargets(site)
	if err != nil {
		t.Fatalf("SelectTargets 오류: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("선별 대상 수: 기대 1, 실제 %d", len(targets))
	}

	tgt := targets[0]
	if tgt.PageURL == "" {
		t.Error("PageURL이 비어 있음")
	}
	if tgt.LocalPath == "" {
		t.Error("LocalPath가 비어 있음")
	}
	if tgt.SourceHash == "" {
		t.Error("SourceHash가 비어 있음")
	}
	if tgt.PageURL != pageURL {
		t.Errorf("PageURL: 기대 %q, 실제 %q", pageURL, tgt.PageURL)
	}
	if tgt.SourceHash != hash {
		t.Errorf("SourceHash: 기대 %q, 실제 %q", hash, tgt.SourceHash)
	}
}

// toSiteRel은 absPath를 siteDir 기준 슬래시 구분 상대경로로 변환한다.
func toSiteRel(t *testing.T, siteDir, absPath string) string {
	t.Helper()
	rel, err := filepath.Rel(filepath.Clean(siteDir), filepath.Clean(absPath))
	if err != nil {
		t.Fatalf("상대경로 변환 실패: %v", err)
	}
	return filepath.ToSlash(rel)
}

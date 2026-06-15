package translator_test

import (
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
	"doc-maker/internal/translator"
)

// testSiteForStore는 임시 siteDir을 가리키는 테스트용 Site를 반환한다.
func testSiteForStore(t *testing.T, id string) (*config.Site, string) {
	t.Helper()
	siteDir := filepath.Join(t.TempDir(), "sites", id)
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}
	return &config.Site{
		ID:      id,
		SiteDir: siteDir,
	}, siteDir
}

// TestTranslatedPath_Structure는 TranslatedPath가 site.OutputDir()/<URL경로구조>를
// 올바르게 생성하는지 확인한다(ANALYSIS D5 — collector의 경로 규칙과 대응).
func TestTranslatedPath_Structure(t *testing.T) {
	cases := []struct {
		name       string
		siteDir    string
		siteID     string
		pageURL    string
		wantSuffix string // filepath.ToSlash 기준 경로 접미사
	}{
		{
			name:       "확장자 없는 URL에 .md 붙임",
			siteDir:    "sites/ollama",
			siteID:     "ollama",
			pageURL:    "https://ollama.com/docs/api/overview",
			wantSuffix: "sites/ollama/output/docs/api/overview.md",
		},
		{
			name:       "이미 .md 확장자 있는 URL",
			siteDir:    "sites/ollama",
			siteID:     "ollama",
			pageURL:    "https://ollama.com/docs/guide.md",
			wantSuffix: "sites/ollama/output/docs/guide.md",
		},
		{
			name:       "비-.md 확장자는 .md로 교체",
			siteDir:    "sites/ollama",
			siteID:     "ollama",
			pageURL:    "https://ollama.com/openapi.yaml",
			wantSuffix: "sites/ollama/output/openapi.md",
		},
		{
			name:       "중첩 경로 구조 보존",
			siteDir:    "/tmp/sites/mysite",
			siteID:     "mysite",
			pageURL:    "https://example.com/section/sub/page",
			wantSuffix: "section/sub/page.md",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			site := &config.Site{ID: tc.siteID, SiteDir: tc.siteDir}
			got := translator.TranslatedPath(site, tc.pageURL)
			// OS 경로 구분자를 슬래시로 정규화해 비교
			gotSlash := filepath.ToSlash(got)

			// 전체 경로가 wantSuffix로 끝나는지 확인한다
			if len(gotSlash) < len(tc.wantSuffix) {
				t.Errorf("경로가 너무 짧음: %q (기대 접미사: %q)", gotSlash, tc.wantSuffix)
				return
			}
			suffix := gotSlash[len(gotSlash)-len(tc.wantSuffix):]
			if suffix != tc.wantSuffix {
				t.Errorf("TranslatedPath(site{SiteDir=%q}, %q)\n 기대 접미사: %q\n 실제: %q",
					tc.siteDir, tc.pageURL, tc.wantSuffix, gotSlash)
			}
		})
	}
}

// TestTranslatedPath_MatchesCollectorConvention은 TranslatedPath가
// collector의 raw 경로와 동일한 하위 경로를 생성하는지 확인한다.
// raw/<원본경로>와 output/<원본경로>가 동일한 상대 경로를 공유해야 한다.
func TestTranslatedPath_MatchesCollectorConvention(t *testing.T) {
	siteDir := "sites/ollama"
	pageURL := "https://ollama.com/docs/api/generate"

	site := &config.Site{ID: "ollama", SiteDir: siteDir}
	rawPath := filepath.Join(siteDir, "raw", "docs", "api", "generate.md")
	translatedPath := translator.TranslatedPath(site, pageURL)

	// 두 경로에서 raw/ 또는 output/ 이후의 상대 경로가 같아야 한다.
	rawBase := filepath.Join(siteDir, "raw")
	outBase := filepath.Join(siteDir, "output")
	rawRel := filepath.ToSlash(rawPath[len(rawBase)+1:])
	transRel := filepath.ToSlash(translatedPath[len(outBase)+1:])

	if rawRel != transRel {
		t.Errorf("collector 경로와 불일치:\n  raw 상대경로: %q\n  translated 상대경로: %q", rawRel, transRel)
	}
}

// TestSaveTranslation_Basic은 번역문 저장과 매니페스트 TranslatedHash 기록이
// 올바르게 동작하는지 확인한다(task-010 핵심 검증 조건).
func TestSaveTranslation_Basic(t *testing.T) {
	site, siteDir := testSiteForStore(t, "testsite")
	pageURL := "https://example.com/docs/intro"
	translatedPath := filepath.Join(siteDir, "output", "docs", "intro.md")
	content := []byte("# 소개\n\n한국어 번역 내용입니다.")
	sourceHash := manifest.HashContent([]byte("# Introduction\n\nOriginal content."))

	// 사전 조건: 매니페스트에 수집 항목 등록(TranslatedHash 없음)
	mf := manifest.New()
	mf.Set(pageURL, manifest.Entry{
		SourceHash: sourceHash,
		SourcePath: "raw/docs/intro.md",
	})
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("사전 매니페스트 저장 실패: %v", err)
	}

	// SaveTranslation 실행
	err := translator.SaveTranslation(site, pageURL, translatedPath, content, sourceHash)
	if err != nil {
		t.Fatalf("SaveTranslation 오류: %v", err)
	}

	// 1. 파일이 translatedPath에 생성되었는지 확인
	if _, err := os.Stat(translatedPath); os.IsNotExist(err) {
		t.Errorf("번역문 파일이 생성되지 않음: %s", translatedPath)
	}

	// 2. 파일 내용이 정확한지 확인
	saved, err := os.ReadFile(translatedPath)
	if err != nil {
		t.Fatalf("저장된 파일 읽기 실패: %v", err)
	}
	if string(saved) != string(content) {
		t.Errorf("저장된 내용 불일치:\n 기대: %q\n 실제: %q", content, saved)
	}

	// 3. 매니페스트에 TranslatedHash가 sourceHash로 기록되었는지 확인
	mfAfter, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("저장 후 매니페스트 로드 실패: %v", err)
	}
	entry, ok := mfAfter.Get(pageURL)
	if !ok {
		t.Fatalf("매니페스트에 pageURL 항목 없음: %s", pageURL)
	}
	if entry.TranslatedHash != sourceHash {
		t.Errorf("TranslatedHash: 기대 %q, 실제 %q", sourceHash, entry.TranslatedHash)
	}

	// 4. NeedsTranslation이 false(재번역 불필요)로 판정되는지 확인
	if mfAfter.NeedsTranslation(pageURL, sourceHash) {
		t.Error("TranslatedHash 기록 후 NeedsTranslation이 여전히 true — 증분 판정 실패")
	}
}

// TestSaveTranslation_PreservesDirectoryStructure는 중간 디렉터리가 없을 때
// 자동 생성되어 경로 구조가 보존되는지 확인한다(ANALYSIS D5).
func TestSaveTranslation_PreservesDirectoryStructure(t *testing.T) {
	site, siteDir := testSiteForStore(t, "site")
	pageURL := "https://example.com/section/subsection/page"
	translatedPath := filepath.Join(siteDir, "output", "section", "subsection", "page.md")
	content := []byte("번역 내용")
	sourceHash := manifest.HashContent([]byte("original"))

	// 매니페스트 사전 준비 없이 SaveTranslation 호출 — 매니페스트가 없어도 동작해야 함
	err := translator.SaveTranslation(site, pageURL, translatedPath, content, sourceHash)
	if err != nil {
		t.Fatalf("SaveTranslation 오류: %v", err)
	}

	// 경로 내 중간 디렉터리가 모두 생성되었는지 확인
	if _, err := os.Stat(translatedPath); os.IsNotExist(err) {
		t.Errorf("경로 구조 보존 실패 — 파일이 없음: %s", translatedPath)
	}
}

// TestSaveTranslation_UpdatesExistingEntry는 기존 매니페스트 항목의
// TranslatedHash만 갱신되고 SourceHash·SourcePath는 보존되는지 확인한다.
func TestSaveTranslation_UpdatesExistingEntry(t *testing.T) {
	site, siteDir := testSiteForStore(t, "site")
	pageURL := "https://example.com/page"
	translatedPath := filepath.Join(siteDir, "output", "page.md")

	originalSourceHash := manifest.HashContent([]byte("original source"))

	// 매니페스트 초기 상태: SourceHash, SourcePath 기록됨
	mf := manifest.New()
	mf.Set(pageURL, manifest.Entry{
		SourceHash: originalSourceHash,
		SourcePath: "raw/page.md",
	})
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		t.Fatalf("초기 매니페스트 저장 실패: %v", err)
	}

	err := translator.SaveTranslation(site, pageURL, translatedPath, []byte("번역"), originalSourceHash)
	if err != nil {
		t.Fatalf("SaveTranslation 오류: %v", err)
	}

	mfAfter, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("저장 후 매니페스트 로드 실패: %v", err)
	}
	entry, ok := mfAfter.Get(pageURL)
	if !ok {
		t.Fatal("항목이 없음")
	}

	// SourceHash가 보존되었는지 확인
	if entry.SourceHash != originalSourceHash {
		t.Errorf("SourceHash가 변경됨: 기대 %q, 실제 %q", originalSourceHash, entry.SourceHash)
	}
	// SourcePath가 보존되었는지 확인
	if entry.SourcePath != "raw/page.md" {
		t.Errorf("SourcePath가 변경됨: 기대 %q, 실제 %q", "raw/page.md", entry.SourcePath)
	}
	// TranslatedHash가 갱신되었는지 확인
	if entry.TranslatedHash != originalSourceHash {
		t.Errorf("TranslatedHash: 기대 %q, 실제 %q", originalSourceHash, entry.TranslatedHash)
	}
}

package collector_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/collector"
	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
	"doc-maker/internal/source"
)

// mockSource는 테스트용 고정 페이지 목록을 반환하는 source.Source 구현체다.
type mockSource struct {
	pages []source.Page
	err   error
}

func (m *mockSource) Pages() ([]source.Page, error) {
	return m.pages, m.err
}

// mockFetch는 URL별 응답을 맵으로 지정할 수 있는 테스트용 fetcher다.
// contents[url]이 존재하면 해당 바이트를 반환하고, 없으면 오류를 반환한다.
type mockFetch struct {
	contents map[string][]byte
	err      map[string]error // URL별 오류 (nil이면 정상 반환)
}

func (m *mockFetch) fetch(u string) ([]byte, error) {
	if m.err != nil {
		if err, ok := m.err[u]; ok {
			return nil, err
		}
	}
	if data, ok := m.contents[u]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("mockFetch: URL 없음: %s", u)
}

// testSite는 테스트용 Site를 생성한다. siteDir을 사이트 폴더로 설정한다.
func testSite(id, siteDir string) *config.Site {
	return &config.Site{
		ID:         id,
		SiteDir:    siteDir,
		BaseURL:    "https://example.com",
		SourceType: config.SourceLLMsTxt,
	}
}

// TestCollect_FirstRun은 1회차 실행에서 대상 원문 모두 저장되고 경로 구조가 원본을 반영하는지 확인한다.
func TestCollect_FirstRun(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")

	site := testSite("example", siteDir)

	pages := []source.Page{
		{URL: "https://example.com/docs/intro", FetchPath: "https://example.com/docs/intro"},
		{URL: "https://example.com/docs/api/overview", FetchPath: "https://example.com/docs/api/overview"},
	}

	mf := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/intro":        []byte("# Intro"),
			"https://example.com/docs/api/overview": []byte("# Overview"),
		},
	}

	res, err := collector.Collect(site, &mockSource{pages: pages}, mf.fetch)
	if err != nil {
		t.Fatalf("Collect 오류: %v", err)
	}

	// 갱신 건수 확인
	if res.Updated != 2 {
		t.Errorf("Updated: 기대 2, 실제 %d", res.Updated)
	}
	if res.Skipped != 0 {
		t.Errorf("Skipped: 기대 0, 실제 %d", res.Skipped)
	}
	if res.Failed != 0 {
		t.Errorf("Failed: 기대 0, 실제 %d", res.Failed)
	}

	// 경로 구조가 원본 URL을 반영하는지 확인(사이트 폴더 안 raw/ 아래)
	introPath := filepath.Join(siteDir, "raw", "docs", "intro.md")
	overviewPath := filepath.Join(siteDir, "raw", "docs", "api", "overview.md")

	if _, err := os.Stat(introPath); os.IsNotExist(err) {
		t.Errorf("intro.md 파일 없음: %s", introPath)
	}
	if _, err := os.Stat(overviewPath); os.IsNotExist(err) {
		t.Errorf("overview.md 파일 없음: %s", overviewPath)
	}

	// 파일 내용 확인
	data, _ := os.ReadFile(introPath)
	if string(data) != "# Intro" {
		t.Errorf("intro.md 내용: 기대 %q, 실제 %q", "# Intro", string(data))
	}

	// 매니페스트가 저장되었는지 확인
	mfLoaded, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("매니페스트 로드 실패: %v", err)
	}
	e, ok := mfLoaded.Get("https://example.com/docs/intro")
	if !ok {
		t.Error("매니페스트에 intro 항목 없음")
	}
	if e.SourcePath == "" {
		t.Error("매니페스트 SourcePath가 비어 있음")
	}
	// SourcePath는 사이트 폴더 기준 상대경로여야 한다
	if e.SourcePath != "raw/docs/intro.md" {
		t.Errorf("SourcePath: 기대 %q, 실제 %q", "raw/docs/intro.md", e.SourcePath)
	}
}

// TestCollect_SecondRunNoChange는 2회차 실행(무변경)에서 갱신 0건이 보고되는지 확인한다.
func TestCollect_SecondRunNoChange(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")

	site := testSite("example", siteDir)

	pages := []source.Page{
		{URL: "https://example.com/docs/intro", FetchPath: "https://example.com/docs/intro"},
	}

	mf := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/intro": []byte("# Intro"),
		},
	}

	// 1회차
	res1, err := collector.Collect(site, &mockSource{pages: pages}, mf.fetch)
	if err != nil {
		t.Fatalf("1회차 Collect 오류: %v", err)
	}
	if res1.Updated != 1 {
		t.Errorf("1회차 Updated: 기대 1, 실제 %d", res1.Updated)
	}

	// 2회차(콘텐츠 동일)
	res2, err := collector.Collect(site, &mockSource{pages: pages}, mf.fetch)
	if err != nil {
		t.Fatalf("2회차 Collect 오류: %v", err)
	}

	if res2.Updated != 0 {
		t.Errorf("2회차 Updated: 기대 0(무변경), 실제 %d", res2.Updated)
	}
	if res2.Skipped != 1 {
		t.Errorf("2회차 Skipped: 기대 1, 실제 %d", res2.Skipped)
	}
}

// TestCollect_PartialChange는 일부 원문 변경 시 변경분만 갱신 건수로 보고되는지 확인한다.
func TestCollect_PartialChange(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")

	site := testSite("example", siteDir)

	pages := []source.Page{
		{URL: "https://example.com/docs/intro", FetchPath: "https://example.com/docs/intro"},
		{URL: "https://example.com/docs/guide", FetchPath: "https://example.com/docs/guide"},
	}

	// 1회차 fetch
	mf1 := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/intro": []byte("# Intro v1"),
			"https://example.com/docs/guide": []byte("# Guide v1"),
		},
	}

	_, err := collector.Collect(site, &mockSource{pages: pages}, mf1.fetch)
	if err != nil {
		t.Fatalf("1회차 Collect 오류: %v", err)
	}

	// 2회차: intro만 변경, guide는 동일
	mf2 := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/intro": []byte("# Intro v2"), // 변경됨
			"https://example.com/docs/guide": []byte("# Guide v1"), // 동일
		},
	}

	res2, err := collector.Collect(site, &mockSource{pages: pages}, mf2.fetch)
	if err != nil {
		t.Fatalf("2회차 Collect 오류: %v", err)
	}

	if res2.Updated != 1 {
		t.Errorf("부분 변경 후 Updated: 기대 1, 실제 %d", res2.Updated)
	}
	if res2.Skipped != 1 {
		t.Errorf("부분 변경 후 Skipped: 기대 1, 실제 %d", res2.Skipped)
	}

	// 변경된 intro.md 내용이 갱신되었는지 확인
	introPath := filepath.Join(siteDir, "raw", "docs", "intro.md")
	data, err := os.ReadFile(introPath)
	if err != nil {
		t.Fatalf("intro.md 읽기 실패: %v", err)
	}
	if string(data) != "# Intro v2" {
		t.Errorf("intro.md 내용: 기대 %q, 실제 %q", "# Intro v2", string(data))
	}

	// guide.md는 변경되지 않았으므로 v1 내용 유지
	guidePath := filepath.Join(siteDir, "raw", "docs", "guide.md")
	gdata, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatalf("guide.md 읽기 실패: %v", err)
	}
	if string(gdata) != "# Guide v1" {
		t.Errorf("guide.md 내용: 기대 %q, 실제 %q", "# Guide v1", string(gdata))
	}
}

// TestCollect_FetchFailurePreservesExisting은 취득 실패 시 기존 원문을 덮어쓰지 않는지 확인한다.
func TestCollect_FetchFailurePreservesExisting(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")

	site := testSite("example", siteDir)

	pages := []source.Page{
		{URL: "https://example.com/docs/intro", FetchPath: "https://example.com/docs/intro"},
	}

	// 1회차: 정상 취득
	mf1 := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/intro": []byte("# Intro original"),
		},
	}

	_, err := collector.Collect(site, &mockSource{pages: pages}, mf1.fetch)
	if err != nil {
		t.Fatalf("1회차 Collect 오류: %v", err)
	}

	// 2회차: 취득 실패
	mf2 := &mockFetch{
		err: map[string]error{
			"https://example.com/docs/intro": fmt.Errorf("네트워크 오류"),
		},
	}

	res2, err := collector.Collect(site, &mockSource{pages: pages}, mf2.fetch)
	if err != nil {
		t.Fatalf("2회차 Collect 오류(실패 페이지 있어도 전체 오류가 아님): %v", err)
	}

	// 취득 실패는 Failed로 카운트
	if res2.Failed != 1 {
		t.Errorf("Failed: 기대 1, 실제 %d", res2.Failed)
	}
	if res2.Updated != 0 {
		t.Errorf("Updated: 기대 0, 실제 %d", res2.Updated)
	}

	// 기존 원문이 보존되어 있는지 확인
	introPath := filepath.Join(siteDir, "raw", "docs", "intro.md")
	data, err := os.ReadFile(introPath)
	if err != nil {
		t.Fatalf("intro.md 읽기 실패: %v", err)
	}
	if string(data) != "# Intro original" {
		t.Errorf("취득 실패 후 기존 원문이 덮어쓰여짐: %q", string(data))
	}
}

// TestCollect_PathStructure는 저장 경로가 URL 경로 구조를 정확히 반영하는지 확인한다(ANALYSIS D5).
func TestCollect_PathStructure(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "mysite")

	site := testSite("mysite", siteDir)

	pages := []source.Page{
		// 깊이 3인 경로
		{URL: "https://example.com/a/b/c", FetchPath: "https://example.com/a/b/c"},
		// 확장자가 있는 경로(.md)
		{URL: "https://example.com/docs/page.md", FetchPath: "https://example.com/docs/page.md"},
		// 루트 직속 경로
		{URL: "https://example.com/top", FetchPath: "https://example.com/top"},
	}

	mf := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/a/b/c":        []byte("deep"),
			"https://example.com/docs/page.md": []byte("page with ext"),
			"https://example.com/top":           []byte("top level"),
		},
	}

	_, err := collector.Collect(site, &mockSource{pages: pages}, mf.fetch)
	if err != nil {
		t.Fatalf("Collect 오류: %v", err)
	}

	cases := []struct {
		desc     string
		expected string
	}{
		{"깊이 3 경로", filepath.Join(siteDir, "raw", "a", "b", "c.md")},
		{"확장자 있는 경로", filepath.Join(siteDir, "raw", "docs", "page.md")},
		{"루트 직속 경로", filepath.Join(siteDir, "raw", "top.md")},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if _, err := os.Stat(tc.expected); os.IsNotExist(err) {
				t.Errorf("파일 없음: %s", tc.expected)
			}
		})
	}
}

// TestCollect_NilFetch는 fetch가 nil일 때 오류를 반환하는지 확인한다.
func TestCollect_NilFetch(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")
	site := testSite("example", siteDir)
	pages := []source.Page{{URL: "https://example.com/page", FetchPath: "https://example.com/page"}}

	_, err := collector.Collect(site, &mockSource{pages: pages}, nil)
	if err == nil {
		t.Error("fetch가 nil일 때 오류를 반환해야 함")
	}
}

// TestCollect_SourceError는 source.Pages()가 오류를 반환할 때 Collect도 오류를 반환하는지 확인한다.
func TestCollect_SourceError(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")
	site := testSite("example", siteDir)
	mf := &mockFetch{contents: map[string][]byte{}}
	srcErr := &mockSource{err: fmt.Errorf("목록 확보 실패")}

	_, err := collector.Collect(site, srcErr, mf.fetch)
	if err == nil {
		t.Error("source.Pages() 오류 시 Collect도 오류를 반환해야 함")
	}
}

// TestCollect_ManifestUpdated는 수집 후 매니페스트에 SourceHash와 SourcePath가 기록되는지 확인한다.
func TestCollect_ManifestUpdated(t *testing.T) {
	tmp := t.TempDir()
	siteDir := filepath.Join(tmp, "sites", "example")

	site := testSite("example", siteDir)

	content := []byte("# Content")
	pages := []source.Page{
		{URL: "https://example.com/docs/page", FetchPath: "https://example.com/docs/page"},
	}

	mf := &mockFetch{
		contents: map[string][]byte{
			"https://example.com/docs/page": content,
		},
	}

	_, err := collector.Collect(site, &mockSource{pages: pages}, mf.fetch)
	if err != nil {
		t.Fatalf("Collect 오류: %v", err)
	}

	loaded, err := manifest.Load(site.ManifestDir())
	if err != nil {
		t.Fatalf("매니페스트 로드 실패: %v", err)
	}

	entry, ok := loaded.Get("https://example.com/docs/page")
	if !ok {
		t.Fatal("매니페스트에 항목 없음")
	}

	expectedHash := manifest.HashContent(content)
	if entry.SourceHash != expectedHash {
		t.Errorf("SourceHash: 기대 %q, 실제 %q", expectedHash, entry.SourceHash)
	}
	if entry.SourcePath == "" {
		t.Error("SourcePath가 비어 있음")
	}
	// SourcePath는 사이트 폴더 기준 상대경로여야 한다
	if entry.SourcePath != "raw/docs/page.md" {
		t.Errorf("SourcePath: 기대 %q, 실제 %q", "raw/docs/page.md", entry.SourcePath)
	}
}

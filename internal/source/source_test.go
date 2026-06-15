package source_test

import (
	"errors"
	"fmt"
	"testing"

	"doc-maker/internal/config"
	"doc-maker/internal/source"
)

// makeSite는 주어진 source 타입으로 최소한의 Site 구조체를 만든다.
func makeSite(st config.SourceType) *config.Site {
	return &config.Site{
		ID:         "test",
		SiteDir:    "./sites/test",
		BaseURL:    "https://docs.example.com",
		SourceType: st,
	}
}

// makeLLMsSite는 포함·제외 패턴을 지정한 llms.txt 타입 Site를 만든다.
func makeLLMsSite(include, exclude []string) *config.Site {
	return &config.Site{
		ID:              "test",
		SiteDir:         "./sites/test",
		BaseURL:         "https://docs.example.com",
		SourceType:      config.SourceLLMsTxt,
		IncludePatterns: include,
		ExcludePatterns: exclude,
	}
}

// TestNew_LLMsTxt는 llms.txt 타입 선택 시 Source 구현체가 반환되는지 확인한다.
func TestNew_LLMsTxt(t *testing.T) {
	site := makeSite(config.SourceLLMsTxt)

	src, err := source.New(site)
	if err != nil {
		t.Fatalf("llms.txt 타입 선택 시 오류가 발생해서는 안 됨: %v", err)
	}
	if src == nil {
		t.Fatal("llms.txt 타입 선택 시 nil이 아닌 구현체가 반환되어야 함")
	}
}

// TestNew_Sitemap_Unimplemented는 sitemap 타입 선택 시 미구현 오류로 중단되는지 확인한다(SPEC §5.7).
func TestNew_Sitemap_Unimplemented(t *testing.T) {
	site := makeSite(config.SourceSitemap)

	src, err := source.New(site)
	if err == nil {
		t.Fatal("sitemap 타입 선택 시 오류가 반환되어야 함")
	}
	if src != nil {
		t.Fatal("sitemap 타입 선택 시 구현체가 nil이어야 함")
	}
	if !errors.Is(err, source.ErrNotImplemented) {
		t.Errorf("sitemap 오류는 ErrNotImplemented를 감싸야 함: 실제 오류 = %v", err)
	}
}

// TestNew_Crawl_Unimplemented는 crawl 타입 선택 시 미구현 오류로 중단되는지 확인한다(SPEC §5.7).
func TestNew_Crawl_Unimplemented(t *testing.T) {
	site := makeSite(config.SourceCrawl)

	src, err := source.New(site)
	if err == nil {
		t.Fatal("crawl 타입 선택 시 오류가 반환되어야 함")
	}
	if src != nil {
		t.Fatal("crawl 타입 선택 시 구현체가 nil이어야 함")
	}
	if !errors.Is(err, source.ErrNotImplemented) {
		t.Errorf("crawl 오류는 ErrNotImplemented를 감싸야 함: 실제 오류 = %v", err)
	}
}

// sampleLLMsTxt는 단위 테스트에 사용할 표본 llms.txt 내용이다.
// 마크다운 링크 형태와 순수 URL 형태를 모두 포함한다.
const sampleLLMsTxt = `# Ollama Documentation

## Getting Started
- [Introduction](https://docs.example.com/docs/introduction)
- [Quickstart](https://docs.example.com/docs/quickstart)

## API Reference
- [REST API](https://docs.example.com/api/rest)
- [Python Library](https://docs.example.com/api/python)

## Tutorials
https://docs.example.com/tutorials/getting-started
https://docs.example.com/tutorials/advanced

## Models
- [Model Library](https://docs.example.com/models/library)
`

// mockFetcher는 llms.txt URL 요청에 고정 내용을 반환하는 테스트용 Fetcher를 만든다.
func mockFetcher(content string) source.Fetcher {
	return func(url string) ([]byte, error) {
		return []byte(content), nil
	}
}

// mockFetcherError는 항상 오류를 반환하는 테스트용 Fetcher를 만든다.
func mockFetcherError(msg string) source.Fetcher {
	return func(url string) ([]byte, error) {
		return nil, fmt.Errorf("%s", msg)
	}
}

// TestPages_AllPages는 패턴 없이 호출하면 llms.txt의 모든 페이지가 반환되는지 확인한다.
func TestPages_AllPages(t *testing.T) {
	site := makeLLMsSite(nil, nil)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(sampleLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}

	// sampleLLMsTxt에는 7개 URL이 있다
	const wantCount = 7
	if len(pages) != wantCount {
		t.Errorf("페이지 수: got %d, want %d", len(pages), wantCount)
		for i, p := range pages {
			t.Logf("  [%d] URL=%s FetchPath=%s", i, p.URL, p.FetchPath)
		}
	}
}

// TestPages_IncludePattern은 포함 패턴이 적용되어 해당 경로만 남는지 확인한다.
func TestPages_IncludePattern(t *testing.T) {
	// docs/* 경로만 포함 (선행 '/' 없음 — filterPages가 url.Path에서 '/'를 제거해 패턴과 맞춤)
	site := makeLLMsSite([]string{"docs/*"}, nil)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(sampleLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}

	// /docs/introduction, /docs/quickstart 두 개만 남아야 한다
	const wantCount = 2
	if len(pages) != wantCount {
		t.Errorf("포함 패턴 적용 후 페이지 수: got %d, want %d", len(pages), wantCount)
		for _, p := range pages {
			t.Logf("  URL=%s", p.URL)
		}
	}
	for _, p := range pages {
		if p.URL != "https://docs.example.com/docs/introduction" &&
			p.URL != "https://docs.example.com/docs/quickstart" {
			t.Errorf("예상치 않은 페이지 포함됨: %s", p.URL)
		}
	}
}

// TestPages_ExcludePattern은 제외 패턴이 적용되어 해당 경로가 제거되는지 확인한다.
func TestPages_ExcludePattern(t *testing.T) {
	// api/* 경로를 제외 (선행 '/' 없음 — filterPages가 url.Path에서 '/'를 제거해 패턴과 맞춤)
	site := makeLLMsSite(nil, []string{"api/*"})
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(sampleLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}

	// 전체 7개 중 /api/rest, /api/python 2개가 제외되어 5개가 남아야 한다
	const wantCount = 5
	if len(pages) != wantCount {
		t.Errorf("제외 패턴 적용 후 페이지 수: got %d, want %d", len(pages), wantCount)
		for _, p := range pages {
			t.Logf("  URL=%s", p.URL)
		}
	}
	for _, p := range pages {
		if p.URL == "https://docs.example.com/api/rest" ||
			p.URL == "https://docs.example.com/api/python" {
			t.Errorf("제외 패턴 대상이 결과에 포함됨: %s", p.URL)
		}
	}
}

// TestPages_IncludeAndExclude는 포함 패턴과 제외 패턴이 함께 적용되는지 확인한다.
func TestPages_IncludeAndExclude(t *testing.T) {
	// docs/* 또는 api/* 포함, 단 api/python 제외 (선행 '/' 없음)
	site := makeLLMsSite(
		[]string{"docs/*", "api/*"},
		[]string{"api/python"},
	)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(sampleLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}

	// /docs/introduction, /docs/quickstart, /api/rest 세 개가 남아야 한다
	const wantCount = 3
	if len(pages) != wantCount {
		t.Errorf("포함+제외 패턴 적용 후 페이지 수: got %d, want %d", len(pages), wantCount)
		for _, p := range pages {
			t.Logf("  URL=%s", p.URL)
		}
	}
}

// ollamaLLMsTxt는 ollama.json의 실제 패턴 검증용 표본 llms.txt다.
// docs 하위에 여러 깊이의 URL과 blog 하위, 그 외 경로를 포함한다.
const ollamaLLMsTxt = `# Ollama Docs
- [Getting Started](https://ollama.com/docs/foo)
- [Deep Page](https://ollama.com/docs/foo/bar)
- [Very Deep](https://ollama.com/docs/foo/bar/baz)
- [Blog Post](https://ollama.com/docs/api/blog/post)
- [Blog Sub](https://ollama.com/docs/api/blog/post/sub)
- [Other](https://ollama.com/other/x)
`

// TestPages_OllamaPatterns는 configs/ollama.json의 실제 패턴을 그대로 사용해
// ** 깊이 매칭과 exclude가 올바르게 동작하는지 확인한다.
//
// include: ["docs/**"]  — docs 하위 임의 깊이 포함
// exclude: ["docs/api/blog/**"]  — blog 하위는 제외
func TestPages_OllamaPatterns(t *testing.T) {
	// ollama.json의 실제 패턴 그대로 사용
	site := &config.Site{
		ID:              "ollama",
		SiteDir:         "./sites/ollama",
		BaseURL:         "https://ollama.com/docs",
		SourceType:      config.SourceLLMsTxt,
		IncludePatterns: []string{"docs/**"},
		ExcludePatterns: []string{"docs/api/blog/**", "docs/changelog/**"},
	}
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(ollamaLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}

	// 기대 결과:
	//   포함: docs/foo, docs/foo/bar, docs/foo/bar/baz  (docs/** 매칭)
	//   제외: docs/api/blog/post, docs/api/blog/post/sub  (docs/api/blog/** 매칭)
	//   제외: other/x  (docs/** 미매칭)
	// → 최종 3개
	const wantCount = 3
	if len(pages) != wantCount {
		t.Errorf("ollama 패턴 적용 후 페이지 수: got %d, want %d", len(pages), wantCount)
		for _, p := range pages {
			t.Logf("  URL=%s", p.URL)
		}
	}

	// docs 하위 여러 깊이가 모두 포함되는지 확인
	wantURLs := map[string]bool{
		"https://ollama.com/docs/foo":         true,
		"https://ollama.com/docs/foo/bar":     true,
		"https://ollama.com/docs/foo/bar/baz": true,
	}
	for _, p := range pages {
		if !wantURLs[p.URL] {
			t.Errorf("예상치 않은 페이지 포함됨: %s", p.URL)
		}
	}

	// blog 하위와 other는 결과에 없어야 한다
	for _, p := range pages {
		if p.URL == "https://ollama.com/docs/api/blog/post" ||
			p.URL == "https://ollama.com/docs/api/blog/post/sub" ||
			p.URL == "https://ollama.com/other/x" {
			t.Errorf("제외되어야 할 페이지가 결과에 포함됨: %s", p.URL)
		}
	}
}

// TestGlobMatch_DoubleStarDepths는 ** 패턴이 여러 깊이를 모두 매칭하는지 단위 검증한다.
func TestGlobMatch_DoubleStarDepths(t *testing.T) {
	cases := []struct {
		pattern string
		target  string
		want    bool
	}{
		// docs/** — 1단계 이상 모두 포함
		{"docs/**", "docs/foo", true},
		{"docs/**", "docs/foo/bar", true},
		{"docs/**", "docs/foo/bar/baz", true},
		// docs/api/blog/** — blog 하위만 매칭
		{"docs/api/blog/**", "docs/api/blog/x", true},
		{"docs/api/blog/**", "docs/api/blog/x/y", true},
		{"docs/api/blog/**", "docs/foo", false},
		// 단일 * — 한 세그먼트 안에서만 매칭
		{"docs/*", "docs/foo", true},
		{"docs/*", "docs/foo/bar", false}, // * 는 '/' 미포함
		// 비매칭
		{"docs/**", "other/x", false},
	}
	for _, tc := range cases {
		// source 패키지의 globMatch는 비공개이므로 filterPages를 통해 간접 검증한다.
		// 직접 검증: makeLLMsSite + mockFetcher로 단일 URL 필터링
		urlStr := "https://host/" + tc.target
		llmsTxt := "- [P](" + urlStr + ")\n"
		site := &config.Site{
			ID:              "t",
			SiteDir:         "./sites/t",
			BaseURL:         "https://host",
			SourceType:      config.SourceLLMsTxt,
			IncludePatterns: []string{tc.pattern},
		}
		src := source.NewLLMsSourceWithFetcher(site, mockFetcher(llmsTxt))
		pages, err := src.Pages()
		if err != nil {
			t.Fatalf("pattern=%q target=%q: Pages() 오류 %v", tc.pattern, tc.target, err)
		}
		got := len(pages) == 1
		if got != tc.want {
			t.Errorf("pattern=%q target=%q: got matched=%v, want %v", tc.pattern, tc.target, got, tc.want)
		}
	}
}

// TestPages_FetcherError는 llms.txt 취득 실패 시 오류가 전달되는지 확인한다.
func TestPages_FetcherError(t *testing.T) {
	site := makeLLMsSite(nil, nil)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcherError("연결 실패"))

	_, err := src.Pages()
	if err == nil {
		t.Fatal("fetcher 오류 시 Pages()가 오류를 반환해야 함")
	}
}

// TestPages_EmptyLLMsTxt는 빈 llms.txt에서 빈 목록이 반환되는지 확인한다.
func TestPages_EmptyLLMsTxt(t *testing.T) {
	site := makeLLMsSite(nil, nil)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(""))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("빈 llms.txt에서 오류 발생: %v", err)
	}
	if len(pages) != 0 {
		t.Errorf("빈 llms.txt에서 페이지가 반환되어서는 안 됨: got %d", len(pages))
	}
}

// TestPages_URLAndFetchPathEqual은 Page의 URL과 FetchPath가 같은지 확인한다.
// task-006에서 FetchPath를 별도로 변환하기 전까지 URL = FetchPath 여야 한다.
func TestPages_URLAndFetchPathEqual(t *testing.T) {
	site := makeLLMsSite(nil, nil)
	src := source.NewLLMsSourceWithFetcher(site, mockFetcher(sampleLLMsTxt))

	pages, err := src.Pages()
	if err != nil {
		t.Fatalf("Pages() 오류 발생: %v", err)
	}
	for _, p := range pages {
		if p.URL != p.FetchPath {
			t.Errorf("URL과 FetchPath가 달라야 함: URL=%s FetchPath=%s", p.URL, p.FetchPath)
		}
	}
}

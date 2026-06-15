// Package source는 페이지 목록 확보 방식을 추상화하는 인터페이스와
// source 타입 값으로 구현체를 선택하는 팩토리를 제공한다(ANALYSIS §3, D7).
package source

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"doc-maker/internal/config"
)

// Page는 source가 산출하는 대상 페이지 하나를 나타낸다.
// URL은 페이지를 식별하는 주소이고, FetchPath는 실제 원문(마크다운)을 내려받을 경로다.
type Page struct {
	// 페이지를 식별하는 URL
	URL string
	// 원문 취득 경로(HTTP URL 또는 로컬 경로)
	FetchPath string
}

// Source는 페이지 목록 확보 방식의 단일 계약이다(ANALYSIS §3, SPEC §5.7).
// 베이스 URL·포함/제외 패턴을 설정으로 받아
// 대상 페이지 목록과 각 페이지 원문 취득 경로를 산출한다.
type Source interface {
	// Pages는 대상 페이지 목록을 반환한다.
	// 설정의 포함/제외 패턴이 적용된 결과를 돌려준다.
	Pages() ([]Page, error)
}

// ErrNotImplemented는 인터페이스는 정의되어 있지만 아직 구현되지 않은
// source 타입을 선택했을 때 반환되는 오류 종류다(SPEC §5.7).
var ErrNotImplemented = errors.New("미구현 source 타입")

// Fetcher는 지정된 URL의 내용을 바이트 슬라이스로 반환하는 함수 타입이다.
// 기본값은 defaultFetch(net/http 사용)이며, 테스트 시 모킹으로 교체할 수 있다.
// task-006의 rate limit 취득기와 합칠 때도 이 타입으로 교체한다.
type Fetcher func(url string) ([]byte, error)

// defaultFetch는 net/http를 이용해 URL의 내용을 반환하는 기본 Fetcher다.
func defaultFetch(rawURL string) ([]byte, error) {
	resp, err := http.Get(rawURL) //nolint:noctx // 기본 fetcher는 context 없이 동작
	if err != nil {
		return nil, fmt.Errorf("HTTP GET 실패 (%s): %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 응답 오류 (%s): 상태 코드 %d", rawURL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 본문 읽기 실패 (%s): %w", rawURL, err)
	}
	return data, nil
}

// New는 site 설정의 SourceType 값으로 알맞은 Source 구현체를 선택해 반환한다.
// sitemap·crawl 타입은 현재 미구현으로, 선택 즉시 ErrNotImplemented를 감싼 오류를 반환한다.
// llms.txt 타입은 LLMsSource 구현체를 반환한다.
func New(site *config.Site) (Source, error) {
	switch site.SourceType {
	case config.SourceLLMsTxt:
		return &LLMsSource{site: site, fetch: defaultFetch}, nil
	case config.SourceSitemap:
		return nil, fmt.Errorf("source 타입 %q: %w", site.SourceType, ErrNotImplemented)
	case config.SourceCrawl:
		return nil, fmt.Errorf("source 타입 %q: %w", site.SourceType, ErrNotImplemented)
	default:
		return nil, fmt.Errorf("알 수 없는 source 타입: %q", site.SourceType)
	}
}

// NewLLMsSourceWithFetcher는 커스텀 Fetcher를 주입한 LLMsSource를 반환한다.
// 단위 테스트에서 네트워크 없이 검증할 때 사용한다.
func NewLLMsSourceWithFetcher(site *config.Site, fetch Fetcher) Source {
	return &LLMsSource{site: site, fetch: fetch}
}

// LLMsSource는 llms.txt 방식으로 페이지 목록을 확보하는 Source 구현체다.
// fetch 필드를 통해 HTTP 취득 의존성을 주입할 수 있어 단위 테스트에서 네트워크 없이 검증 가능하다.
type LLMsSource struct {
	site  *config.Site
	fetch Fetcher
}

// mdLinkRe는 마크다운 링크 형태 `[title](url)` 또는 `- [title](url)` 패턴을 인식한다.
// llms.txt 관례상 마크다운 링크 목록 형태로 페이지가 나열된다.
var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// Pages는 llms.txt에서 페이지 목록을 취득하고 포함/제외 패턴을 적용해 반환한다.
//
// 가정한 llms.txt 포맷:
//   - 마크다운 링크 `[title](url)` 또는 `- [title](url)` 형태의 라인
//   - 또는 http:// / https:// 로 시작하는 순수 URL 라인
//   - 빈 줄과 `#` 으로 시작하는 주석 라인은 무시
//
// 패턴 매칭은 path.Match 기반 glob를 URL 경로 부분에 적용한다.
// IncludePatterns가 비어 있으면 전체 포함, ExcludePatterns가 비어 있으면 제외 없음.
func (s *LLMsSource) Pages() ([]Page, error) {
	// llms.txt URL 구성: baseURL + "/llms.txt"
	llmsURL, err := buildLLMsURL(s.site.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("llms.txt URL 구성 실패: %w", err)
	}

	// llms.txt 내용 취득
	data, err := s.fetch(llmsURL)
	if err != nil {
		return nil, fmt.Errorf("llms.txt 취득 실패 (%s): %w", llmsURL, err)
	}

	// 페이지 목록 파싱
	pages, err := parseLLMsTxt(data, s.site.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("llms.txt 파싱 실패: %w", err)
	}

	// 포함/제외 패턴 필터 적용
	filtered := filterPages(pages, s.site.IncludePatterns, s.site.ExcludePatterns)

	return filtered, nil
}

// buildLLMsURL은 베이스 URL에 "/llms.txt" 경로를 붙여 반환한다.
func buildLLMsURL(baseURL string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("베이스 URL 파싱 실패: %w", err)
	}
	// 경로 끝 슬래시를 제거하고 /llms.txt를 붙인다
	u.Path = strings.TrimRight(u.Path, "/") + "/llms.txt"
	return u.String(), nil
}

// parseLLMsTxt는 llms.txt 내용을 파싱해 Page 슬라이스로 반환한다.
// baseURL은 상대 경로를 절대 URL로 변환할 때 사용한다.
func parseLLMsTxt(data []byte, baseURL string) ([]Page, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("베이스 URL 파싱 실패: %w", err)
	}

	var pages []Page
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 빈 줄과 주석(#으로 시작) 무시
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 마크다운 링크 형태 `[title](url)` 파싱 시도
		if m := mdLinkRe.FindStringSubmatch(line); m != nil {
			rawURL := strings.TrimSpace(m[2])
			page, ok := makePageFromURL(rawURL, base)
			if ok {
				pages = append(pages, page)
			}
			continue
		}

		// 순수 URL 라인 (http:// 또는 https://)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			page, ok := makePageFromURL(line, base)
			if ok {
				pages = append(pages, page)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("llms.txt 스캔 오류: %w", err)
	}

	return pages, nil
}

// makePageFromURL은 rawURL로부터 Page를 생성한다.
// URL이 절대 주소이면 그대로, 상대 주소이면 base를 기준으로 절대화한다.
// 파싱 실패 시 (zero, false)를 반환한다.
func makePageFromURL(rawURL string, base *url.URL) (Page, bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Page{}, false
	}
	abs := base.ResolveReference(u)
	pageURL := abs.String()
	return Page{
		URL:       pageURL,
		FetchPath: pageURL,
	}, true
}

// filterPages는 pages에 IncludePatterns·ExcludePatterns를 적용해 결과를 반환한다.
//
// 패턴 매칭 규칙 (globMatch 기반):
//   - 패턴은 URL의 경로(path) 부분에만 적용한다.
//   - url.Path의 선행 '/'를 제거해 "docs/foo" 형태로 패턴과 맞춘다.
//     (패턴은 "docs/**"처럼 선행 '/' 없이 작성한다.)
//   - IncludePatterns가 비어 있으면 전체 포함으로 간주한다.
//   - 페이지가 IncludePatterns 중 하나 이상에 매칭되어야 후보가 된다.
//   - 후보 중 ExcludePatterns 중 하나라도 매칭되면 제외된다.
func filterPages(pages []Page, includePatterns, excludePatterns []string) []Page {
	var result []Page
	for _, p := range pages {
		u, err := url.Parse(p.URL)
		if err != nil {
			// URL 파싱 실패한 항목은 건너뜀
			continue
		}
		// url.Path는 "/docs/foo" 형태이므로 선행 '/'를 제거해 패턴과 통일한다.
		// 패턴(configs/*.json)은 "docs/**"처럼 선행 '/' 없이 작성되어 있다.
		urlPath := strings.TrimPrefix(u.Path, "/")

		// 포함 패턴 검사
		if len(includePatterns) > 0 {
			included := matchesAny(urlPath, includePatterns)
			if !included {
				continue
			}
		}

		// 제외 패턴 검사
		if len(excludePatterns) > 0 {
			excluded := matchesAny(urlPath, excludePatterns)
			if excluded {
				continue
			}
		}

		result = append(result, p)
	}
	return result
}

// matchesAny는 target이 patterns 중 하나라도 globMatch로 매칭되면 true를 반환한다.
func matchesAny(target string, patterns []string) bool {
	for _, pat := range patterns {
		if globMatch(pat, target) {
			return true
		}
	}
	return false
}

// globMatch는 ** と * を含む glob パターンと target を照合する。
//
// 매칭 의미:
//   - '*'  : 단일 경로 세그먼트 안에서만 임의 문자 매칭 ('/' 미포함)
//   - '**' : 임의 깊이의 경로 전체 매칭 ('/' 포함, 0개 이상 세그먼트)
//
// 구현 전략: 패턴을 세그먼트로 분리한 뒤 재귀적으로 매칭.
// "**" 세그먼트는 0개 이상의 target 세그먼트를 소비하는 와일드카드로 처리한다.
func globMatch(pattern, target string) bool {
	patSegs := strings.Split(pattern, "/")
	tgtSegs := strings.Split(target, "/")
	return globMatchSegs(patSegs, tgtSegs)
}

// globMatchSegs는 패턴 세그먼트 슬라이스와 대상 세그먼트 슬라이스를 재귀 매칭한다.
func globMatchSegs(patSegs, tgtSegs []string) bool {
	// 패턴과 대상을 모두 소비하면 매칭 성공
	if len(patSegs) == 0 && len(tgtSegs) == 0 {
		return true
	}
	// 패턴이 남아 있지만 대상이 없는 경우: "**"만 남아 있으면 0개 매칭으로 성공
	if len(patSegs) > 0 && len(tgtSegs) == 0 {
		// 남은 패턴이 모두 "**"인 경우에만 성공
		for _, seg := range patSegs {
			if seg != "**" {
				return false
			}
		}
		return true
	}
	// 패턴이 없는데 대상이 남아 있으면 실패
	if len(patSegs) == 0 {
		return false
	}

	pat := patSegs[0]
	if pat == "**" {
		// '**' 세그먼트: 0개 이상의 대상 세그먼트를 소비하며 시도
		// 0개 소비(건너뜀) 또는 1개씩 소비하며 나머지와 재귀 매칭
		for i := 0; i <= len(tgtSegs); i++ {
			if globMatchSegs(patSegs[1:], tgtSegs[i:]) {
				return true
			}
		}
		return false
	}

	// 단일 '*' 또는 리터럴 세그먼트: path.Match로 단일 세그먼트만 매칭
	// path.Match의 '*'는 '/'를 포함하지 않으므로 세그먼트 단위에서는 올바르게 동작한다.
	matched, err := path.Match(pat, tgtSegs[0])
	if err != nil || !matched {
		return false
	}
	return globMatchSegs(patSegs[1:], tgtSegs[1:])
}

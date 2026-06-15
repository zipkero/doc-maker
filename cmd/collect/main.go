// collect는 사이트 식별자를 인자로 받아 원문 수집을 실행하는 CLI 진입점이다.
//
// 사용법:
//
//	collect [플래그] <사이트-식별자>
//
// 사이트 식별자는 sites 디렉터리 아래의 폴더명이다.
// 예: sites/ollama/ → 식별자 "ollama"
//
// 기본값:
//
//	-sites ./sites     사이트 루트 디렉터리
//
// 종료 코드:
//
//	0 정상 완료
//	1 인자 오류 / 설정 로드 실패 / source 초기화 실패 / 수집 오류
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"doc-maker/internal/collector"
	"doc-maker/internal/config"
	"doc-maker/internal/fetcher"
	"doc-maker/internal/source"
)

func main() {
	// 플래그 정의
	sitesRoot := flag.String("sites", "./sites", "사이트 루트 디렉터리")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: collect [플래그] <사이트-식별자>\n\n")
		fmt.Fprintf(os.Stderr, "사이트 식별자 예시: ollama  (→ sites/ollama/ 폴더)\n\n")
		fmt.Fprintf(os.Stderr, "플래그:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "오류: 사이트 식별자가 필요합니다.")
		flag.Usage()
		os.Exit(1)
	}

	siteID := flag.Arg(0)

	if err := run(siteID, *sitesRoot, fetcher.NewDefault().Fetch); err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}

// run은 수집 흐름 전체를 실행한다. 의존성을 인자로 주입받으므로 테스트에서 모킹 가능하다.
//
// 흐름:
//  1. sitesRoot/<siteID>/config.json 설정 로드
//  2. source.New로 Source 생성 (미구현 타입이면 오류 노출)
//  3. collector.Collect 호출
//  4. 갱신/스킵/실패 건수 출력
func run(
	siteID string,
	sitesRoot string,
	fetch source.Fetcher,
) error {
	// 사이트 폴더 경로: sitesRoot/<siteID>
	siteDir := filepath.Join(sitesRoot, siteID)

	site, err := config.Load(siteDir)
	if err != nil {
		return fmt.Errorf("설정 로드 실패 (%s): %w", siteDir, err)
	}

	// Source 생성: 미구현 타입(sitemap, crawl)은 ErrNotImplemented로 표면화.
	var src source.Source
	switch site.SourceType {
	case config.SourceLLMsTxt:
		src = source.NewLLMsSourceWithFetcher(site, fetch)
	default:
		// sitemap, crawl 등 미구현 타입: source.New를 통해 ErrNotImplemented를 얻는다
		_, err = source.New(site)
		if errors.Is(err, source.ErrNotImplemented) {
			return fmt.Errorf("source 타입 %q은 아직 구현되지 않았습니다: %w", site.SourceType, source.ErrNotImplemented)
		}
		return fmt.Errorf("source 초기화 실패: %w", err)
	}

	fmt.Printf("수집 시작: 사이트=%s, 폴더=%s\n", siteID, siteDir)

	// 수집 실행
	result, err := collector.Collect(site, src, fetch)
	if err != nil {
		return fmt.Errorf("수집 실행 실패: %w", err)
	}

	// 결과 출력
	fmt.Printf("수집 완료: 갱신=%d, 스킵=%d, 실패=%d\n", result.Updated, result.Skipped, result.Failed)
	if result.Failed > 0 {
		fmt.Fprintf(os.Stderr, "경고: %d건 취득 실패 (기존 원문 보존됨)\n", result.Failed)
	}

	return nil
}

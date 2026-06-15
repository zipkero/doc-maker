// translate는 번역 대상 출력과 완료 기록을 제공하는 CLI 진입점이다.
//
// 사용법:
//
//	translate [플래그] <서브커맨드> <사이트-식별자>
//
// 서브커맨드:
//
//	plan    미번역·변경분 번역 대상 목록을 출력한다.
//	commit  번역문이 output/에 존재하는 페이지의 매니페스트 translated_hash를 기록한다.
//
// 기본값:
//
//	-sites ./sites     사이트 루트 디렉터리
//
// 종료 코드:
//
//	0 정상 완료
//	1 인자 오류 / 설정 로드 실패 / 실행 오류
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"doc-maker/internal/config"
	"doc-maker/internal/translator"
)

func main() {
	// 최상위 플래그
	sitesRoot := flag.String("sites", "./sites", "사이트 루트 디렉터리")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: translate [플래그] <서브커맨드> <사이트-식별자>\n\n")
		fmt.Fprintf(os.Stderr, "서브커맨드:\n")
		fmt.Fprintf(os.Stderr, "  plan    미번역·변경분 대상 목록 출력\n")
		fmt.Fprintf(os.Stderr, "  commit  번역 완료 해시 기록(output/에 파일이 존재하는 페이지 한정)\n\n")
		fmt.Fprintf(os.Stderr, "플래그:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "오류: 서브커맨드와 사이트 식별자가 필요합니다.")
		flag.Usage()
		os.Exit(1)
	}

	subCmd := args[0]
	siteID := args[1]

	var err error
	switch subCmd {
	case "plan":
		err = runPlan(siteID, *sitesRoot, os.Stdout)
	case "commit":
		err = runCommit(siteID, *sitesRoot, os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "오류: 알 수 없는 서브커맨드 %q\n", subCmd)
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}

// runPlan은 사이트의 미번역·변경분 번역 대상 목록을 w에 출력한다.
//
// 출력 형식(각 항목):
//
//	[N] PageURL
//	    로컬 경로: <LocalPath>
//	    원문 해시: <SourceHash 앞 12자>...
func runPlan(siteID, sitesRoot string, w *os.File) error {
	site, err := loadSite(siteID, sitesRoot)
	if err != nil {
		return err
	}

	targets, err := translator.SelectTargets(site)
	if err != nil {
		return fmt.Errorf("번역 대상 선별 실패: %w", err)
	}

	if len(targets) == 0 {
		fmt.Fprintln(w, "번역할 페이지가 없습니다 (모두 최신 번역 상태).")
		return nil
	}

	fmt.Fprintf(w, "번역 대상: %d 페이지 (사이트=%s)\n\n", len(targets), siteID)
	for i, tgt := range targets {
		// 원문 해시는 앞 12자만 표시해 가독성을 높인다.
		hashPrefix := tgt.SourceHash
		if len(hashPrefix) > 12 {
			hashPrefix = hashPrefix[:12] + "..."
		}
		fmt.Fprintf(w, "[%d] %s\n", i+1, tgt.PageURL)
		fmt.Fprintf(w, "    로컬 경로: %s\n", tgt.LocalPath)
		fmt.Fprintf(w, "    원문 해시: %s\n", hashPrefix)
	}

	return nil
}

// runCommit은 SelectTargets 대상 중 TranslatedPath에 파일이 실제로 존재하는
// 페이지를 찾아, 매니페스트의 TranslatedHash를 현재 원문 해시(SourceHash)로 기록한다.
//
// 번역문 파일(output/)은 건드리지 않는다. 파일이 없는 페이지는 건너뛴다.
func runCommit(siteID, sitesRoot string, w *os.File) error {
	site, err := loadSite(siteID, sitesRoot)
	if err != nil {
		return err
	}

	targets, err := translator.SelectTargets(site)
	if err != nil {
		return fmt.Errorf("번역 대상 선별 실패: %w", err)
	}

	if len(targets) == 0 {
		fmt.Fprintln(w, "기록할 번역 완료 페이지가 없습니다 (모두 최신 번역 상태).")
		return nil
	}

	committed := 0
	skipped := 0
	for _, tgt := range targets {
		// 번역문 출력 경로 계산
		outPath := translator.TranslatedPath(site, tgt.PageURL)

		// 번역문 파일이 실제로 존재하는지 확인한다.
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			// 파일 없음: 아직 번역되지 않은 페이지 → 건너뜀
			skipped++
			continue
		}

		// 번역문 파일이 존재: 매니페스트의 TranslatedHash를 SourceHash로 기록한다.
		if err := translator.CommitTranslation(site, tgt.PageURL, tgt.SourceHash); err != nil {
			return fmt.Errorf("해시 기록 실패 (%s): %w", tgt.PageURL, err)
		}
		fmt.Fprintf(w, "기록: %s\n", tgt.PageURL)
		committed++
	}

	fmt.Fprintf(w, "\n완료 기록: %d건 기록, %d건 스킵(번역문 없음)\n", committed, skipped)
	return nil
}

// loadSite는 sitesRoot/<siteID>/config.json을 읽어 Site를 반환한다.
func loadSite(siteID, sitesRoot string) (*config.Site, error) {
	siteDir := filepath.Join(sitesRoot, siteID)
	site, err := config.Load(siteDir)
	if err != nil {
		return nil, fmt.Errorf("설정 로드 실패 (%s): %w", siteDir, err)
	}
	return site, nil
}

// Package collector는 수집 파이프라인의 증분 판정과 사이트별 원문 저장 단계를 담당한다.
// source로부터 페이지 목록을 받아, 각 페이지의 원문을 취득하고 해시 비교 후 신규·변경분만
// 저장하며 매니페스트를 갱신한다(SPEC §5.2, §5.3 / ANALYSIS §2, §5 D4, D5, D8).
//
// 저장 레이아웃 규약(D4, D5):
//
//	sites/<siteID>/
//	├─ config.json
//	├─ glossary.json
//	├─ raw/<원본경로>.md     ← 이 패키지가 저장하는 위치
//	├─ output/<원본경로>.md
//	└─ manifest.json         ← 이 패키지가 갱신하는 파일
package collector

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
	"doc-maker/internal/source"
)

// Result는 수집 실행 결과 통계를 담는다(SPEC §5.3).
type Result struct {
	// Updated는 이번 실행에서 신규 추가되거나 변경된 원문 건수다.
	Updated int
	// Skipped는 해시가 동일해 재취득을 건너뛴 건수다.
	Skipped int
	// Failed는 취득 최종 실패로 갱신하지 못한 건수다.
	Failed int
}

// Collect는 site의 페이지 목록을 src에서 가져와 증분 판정 후 원문을 저장한다.
//
// 매개변수:
//   - site: 사이트 설정(ID, SiteDir 등). site.RawDir()이 원문 저장 루트로 쓰인다.
//   - src: 페이지 목록 공급자(source.Source)
//   - fetch: 원문 취득 함수(source.Fetcher 타입과 호환). nil이면 오류를 반환한다.
//
// 원문 저장 경로 규칙 (D5):
//
//	site.RawDir()/<URL 경로 구조>
//	예) site.SiteDir="sites/ollama", URL path="/api/chat"
//	    → sites/ollama/raw/api/chat.md
//
// 매니페스트 SourcePath 규칙:
//
//	사이트 폴더 기준 상대경로(예: "raw/api/chat.md")로 기록해
//	레이아웃 이동 시에도 상대 참조가 유지되도록 한다.
//
// 취득 최종 실패(fetch가 nil/err 반환) 시 해당 페이지는 갱신 실패로 분류하고
// 기존 보관 원문을 덮어쓰지 않는다(SPEC §3 원문 보존 / task-006 계약).
func Collect(
	site *config.Site,
	src source.Source,
	fetch source.Fetcher,
) (Result, error) {
	if fetch == nil {
		return Result{}, fmt.Errorf("fetch 함수가 nil입니다")
	}

	// 페이지 목록 확보
	pages, err := src.Pages()
	if err != nil {
		return Result{}, fmt.Errorf("페이지 목록 확보 실패: %w", err)
	}

	// 기존 매니페스트 로드(파일 없으면 빈 Manifest 반환)
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		return Result{}, fmt.Errorf("매니페스트 로드 실패: %w", err)
	}

	// 원문 저장 루트: sites/<siteID>/raw/
	rawDir := site.RawDir()

	var res Result

	for _, page := range pages {
		// 원문 저장 로컬 경로 도출
		localPath := pageURLToLocalPath(rawDir, page.URL)

		// 원문 취득
		data, fetchErr := fetch(page.FetchPath)
		if fetchErr != nil || data == nil {
			// 취득 최종 실패: 기존 원문을 보존하고 건너뜀
			res.Failed++
			continue
		}

		// 콘텐츠 해시 계산
		hash := manifest.HashContent(data)

		// 매니페스트와 비교해 변경 여부 판정(ANALYSIS D8)
		if !mf.IsChanged(page.URL, hash) {
			// 해시 동일 → 재저장하지 않음
			res.Skipped++
			continue
		}

		// 신규 또는 변경: 로컬 디렉터리 생성 후 원문 저장
		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			return Result{}, fmt.Errorf("디렉터리 생성 실패 (%s): %w", filepath.Dir(localPath), err)
		}
		if err := os.WriteFile(localPath, data, 0o644); err != nil {
			return Result{}, fmt.Errorf("원문 파일 쓰기 실패 (%s): %w", localPath, err)
		}

		// 매니페스트 SourcePath: 사이트 폴더 기준 상대경로(슬래시 통일)
		relPath := toSiteDirRelPath(site.SiteDir, localPath)

		// 매니페스트 항목 갱신(SourceHash, SourcePath 기록)
		mf.Set(page.URL, manifest.Entry{
			SourceHash: hash,
			SourcePath: relPath,
		})
		res.Updated++
	}

	// 매니페스트 저장
	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		return Result{}, fmt.Errorf("매니페스트 저장 실패: %w", err)
	}

	return res, nil
}

// toSiteDirRelPath는 절대/OS경로인 localPath를 siteDir 기준 슬래시 구분 상대경로로 변환한다.
// 예) siteDir="sites/ollama", localPath="sites\ollama\raw\api\chat.md"
//
//	→ "raw/api/chat.md"
func toSiteDirRelPath(siteDir, localPath string) string {
	// 두 경로를 모두 filepath.Clean 처리해 OS 구분자로 통일
	cleanSiteDir := filepath.Clean(siteDir)
	cleanLocal := filepath.Clean(localPath)

	// siteDir을 prefix로 제거
	rel, err := filepath.Rel(cleanSiteDir, cleanLocal)
	if err != nil {
		// 상대경로 변환 실패 시 절대경로 그대로 반환 (폴백)
		return filepath.ToSlash(cleanLocal)
	}
	// Windows 백슬래시를 슬래시로 통일
	return filepath.ToSlash(rel)
}

// pageURLToLocalPath는 URL에서 경로 부분을 추출해 rawDir 아래의 로컬 파일 경로를 반환한다.
//
// 변환 규칙:
//  1. url.Parse로 URL 경로(u.Path) 추출
//  2. 선행 '/' 제거
//  3. 경로 구분자를 OS 구분자로 변환(filepath.FromSlash) — Windows 안전
//  4. 확장자가 없거나 비어 있으면 ".md"를 붙임
//  5. filepath.Clean으로 중복 구분자·상대 참조 제거
//  6. rawDir과 결합
//
// 예) rawDir="sites/ollama/raw", URL="https://ollama.com/api/chat"
//
//	→ sites/ollama/raw/api/chat.md
func pageURLToLocalPath(rawDir, pageURL string) string {
	u, err := url.Parse(pageURL)
	if err != nil {
		return filepath.Join(rawDir, "unknown.md")
	}

	// 선행 '/'를 제거해 상대 경로로 만든다.
	rawPath := strings.TrimPrefix(u.Path, "/")

	// '/'를 OS 경로 구분자로 변환(Windows: '\')
	localRelPath := filepath.FromSlash(rawPath)

	// 확장자가 없으면 ".md"를 붙인다.
	if filepath.Ext(localRelPath) == "" {
		localRelPath += ".md"
	}

	return filepath.Clean(filepath.Join(rawDir, localRelPath))
}

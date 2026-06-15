// Package translator는 번역 경계의 입력 로드와 증분 선별 단계를 담당한다.
// 수집 경계가 남긴 로컬 원문·매니페스트를 읽어, 미번역이거나 원문 해시가 직전 번역분과
// 달라진 페이지만 번역 대상으로 골라 반환한다(SPEC §5.5 / ANALYSIS §2, §5 D8).
//
// 사이트 폴더 레이아웃 규약(D4, D5):
//
//	sites/<siteID>/
//	├─ raw/<원본경로>.md     ← 원문(SelectTargets가 읽음)
//	├─ output/<원본경로>.md  ← 번역문(SaveTranslation이 저장)
//	├─ glossary.json         ← 용어집(LoadGlossary가 읽음)
//	└─ manifest.json         ← 증분 기록(SelectTargets·SaveTranslation이 읽고 씀)
//
// 용어집 로드와 실제 번역 텍스트 생성·저장·TranslatedHash 갱신은 task-010의 소관이다.
package translator

import (
	"fmt"
	"os"
	"path/filepath"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
)

// TranslationTarget은 번역 대상으로 선별된 페이지 하나의 정보다.
// task-010이 번역 수행·저장·TranslatedHash 갱신에 필요한 값을 모두 담는다.
type TranslationTarget struct {
	// PageURL은 페이지 식별자이자 매니페스트 키다.
	PageURL string

	// LocalPath는 로컬 원문 파일 경로다(sites/<siteID>/raw/... 아래).
	LocalPath string

	// SourceHash는 현재 원문의 SHA-256 해시다(hex).
	// task-010이 번역 완료 후 매니페스트의 TranslatedHash에 기록한다.
	SourceHash string
}

// SelectTargets는 site의 raw 디렉터리에서 원문 목록을 읽고,
// 매니페스트의 NeedsTranslation 판정으로 미번역·변경분만 선별해 반환한다.
//
// 매개변수:
//   - site: 사이트 설정. site.RawDir()에서 원문을 읽고, site.ManifestDir()에서 매니페스트를 읽는다.
//
// 반환:
//   - 번역이 필요한 TranslationTarget 슬라이스(순서는 파일시스템 Walk 순서).
//   - 오류(매니페스트 로드 실패, 원문 읽기 실패 등).
//
// 선별 기준(SPEC §5.5 / ANALYSIS D8):
//   - 매니페스트 TranslatedHash == "" (미번역)   → 대상 포함
//   - 매니페스트 TranslatedHash != SourceHash    → 대상 포함 (원문 변경)
//   - 매니페스트 TranslatedHash == SourceHash    → 제외 (이미 번역된 무변경)
func SelectTargets(site *config.Site) ([]TranslationTarget, error) {
	// 매니페스트 로드(파일 없으면 빈 Manifest 반환)
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		return nil, fmt.Errorf("매니페스트 로드 실패 (siteID=%s): %w", site.ID, err)
	}

	// 원문 디렉터리: sites/<siteID>/raw/
	rawDir := site.RawDir()

	// 디렉터리가 없으면 수집이 아직 실행되지 않은 것으로 보고 빈 목록 반환
	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		return nil, nil
	}

	var targets []TranslationTarget

	// rawDir 아래의 모든 파일을 순회해 매니페스트 항목과 대조한다.
	// 매니페스트는 페이지 URL을 키로 삼으므로, SourcePath로 역매핑한다.
	//
	// SourcePath는 사이트 폴더 기준 슬래시 구분 상대경로(예: "raw/api/chat.md").
	// 역매핑 시 site.SiteDir + "/" + sourcePath 로 절대화한다.
	type entryWithURL struct {
		url   string
		entry manifest.Entry
	}
	pathIndex := make(map[string]entryWithURL)
	for pageURL, entry := range mf.Entries {
		// 사이트 폴더 기준 상대경로 → 실제 파일 경로
		absPath := filepath.Clean(filepath.Join(site.SiteDir, filepath.FromSlash(entry.SourcePath)))
		pathIndex[absPath] = entryWithURL{url: pageURL, entry: entry}
	}

	// rawDir 아래의 파일을 Walk하며 매니페스트 항목과 대조한다.
	err = filepath.Walk(rawDir, func(path string, info os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}
		if info.IsDir() {
			return nil
		}

		// 매니페스트 역인덱스에서 해당 파일의 URL과 Entry를 조회한다.
		key := filepath.Clean(path)
		ev, found := pathIndex[key]
		if !found {
			// 매니페스트에 없는 파일: 수집 결과 아님 → 건너뜀
			return nil
		}

		// 원문 파일을 읽어 현재 해시를 계산한다.
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("원문 파일 읽기 실패 (%s): %w", path, err)
		}
		currentHash := manifest.HashContent(data)

		// NeedsTranslation: TranslatedHash != currentHash (미번역·변경 모두 포함)
		if mf.NeedsTranslation(ev.url, currentHash) {
			targets = append(targets, TranslationTarget{
				PageURL:    ev.url,
				LocalPath:  path,
				SourceHash: currentHash,
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("원문 디렉터리 탐색 실패 (%s): %w", rawDir, err)
	}

	return targets, nil
}

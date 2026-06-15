// 번역문 저장: 번역 완료된 텍스트를 출력 경로에 원본 구조 보존하여 저장하고,
// 매니페스트의 TranslatedHash를 기록한다(SPEC §5.4, §5.6 / ANALYSIS §2, §5 D1, D2, D5, D8).
package translator

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"doc-maker/internal/config"
	"doc-maker/internal/manifest"
)

// TranslatedPath는 번역문 출력 경로를 계산한다.
//
// 경로 규칙(ANALYSIS D5 — collector의 pageURLToLocalPath와 동일 변환 로직):
//
//	site.OutputDir()/<URL 경로 구조>
//	예) site.SiteDir="sites/ollama", pageURL="https://ollama.com/api/chat"
//	    → sites/ollama/output/api/chat.md
//
// 변환 규칙:
//  1. url.Parse로 URL 경로(u.Path) 추출
//  2. 선행 '/' 제거
//  3. '/'를 OS 경로 구분자로 변환(filepath.FromSlash)
//  4. 확장자를 ".md"로 정규화한다(번역문은 항상 마크다운이므로):
//     확장자가 없으면 ".md"를 붙이고, ".md"가 아닌 확장자(예: ".yaml")는 ".md"로 교체한다.
//     따라서 비-.md 원문(raw/openapi.yaml)의 출력은 output/openapi.md가 되어
//     raw와 확장자만 달라진다(하위 경로 구조는 동일).
//  5. filepath.Clean 적용
//  6. site.OutputDir() 아래에 결합
func TranslatedPath(site *config.Site, pageURL string) string {
	outputDir := site.OutputDir()

	u, err := url.Parse(pageURL)
	if err != nil {
		return filepath.Join(outputDir, "unknown.md")
	}

	rawPath := strings.TrimPrefix(u.Path, "/")
	localRelPath := filepath.FromSlash(rawPath)

	if ext := filepath.Ext(localRelPath); ext == "" {
		localRelPath += ".md"
	} else if ext != ".md" {
		localRelPath = strings.TrimSuffix(localRelPath, ext) + ".md"
	}

	return filepath.Clean(filepath.Join(outputDir, localRelPath))
}

// CommitTranslation은 번역문 파일을 건드리지 않고, 매니페스트의 TranslatedHash만
// sourceHash로 갱신한다.
//
// SaveTranslation이 "파일 저장 + 해시 기록"을 함께 수행하는 것과 달리,
// CommitTranslation은 Claude가 output/에 이미 번역문을 작성한 뒤 CLI가
// 매니페스트만 갱신할 때 사용한다(task-013 commit 서브커맨드 용도).
//
// 매개변수:
//   - site: 사이트 설정(매니페스트 경로 도출에 사용).
//   - pageURL: 페이지 식별자이자 매니페스트 키.
//   - sourceHash: 번역한 원문의 SHA-256 해시(TranslatedHash에 기록할 값).
func CommitTranslation(site *config.Site, pageURL, sourceHash string) error {
	// 매니페스트 로드
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		return fmt.Errorf("매니페스트 로드 실패 (siteID=%s): %w", site.ID, err)
	}

	// 기존 Entry를 가져와 TranslatedHash만 교체한다. 기존 Entry가 없으면 새로 생성한다.
	entry, _ := mf.Get(pageURL)
	entry.TranslatedHash = sourceHash
	mf.Set(pageURL, entry)

	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		return fmt.Errorf("매니페스트 저장 실패 (siteID=%s): %w", site.ID, err)
	}

	return nil
}

// SaveTranslation은 번역된 텍스트를 translatedPath에 저장하고,
// site.ManifestDir()의 매니페스트에서 pageURL 항목의 TranslatedHash를 sourceHash로 갱신한다.
//
// 매개변수:
//   - site: 사이트 설정(매니페스트 경로 도출에 사용).
//   - pageURL: 페이지 식별자이자 매니페스트 키.
//   - translatedPath: 번역문을 기록할 로컬 파일 경로(TranslatedPath로 계산한 경로).
//   - content: 번역된 텍스트(바이트 슬라이스).
//   - sourceHash: 번역한 원문의 SHA-256 해시(TranslationTarget.SourceHash). TranslatedHash로 기록된다.
//
// 동작:
//  1. translatedPath의 디렉터리를 생성한다.
//  2. content를 translatedPath에 저장한다.
//  3. site.ManifestDir()의 매니페스트를 로드한다.
//  4. pageURL 항목의 TranslatedHash를 sourceHash로 갱신하고 저장한다.
func SaveTranslation(site *config.Site, pageURL, translatedPath string, content []byte, sourceHash string) error {
	// 1. 번역문 디렉터리 생성
	if err := os.MkdirAll(filepath.Dir(translatedPath), 0o755); err != nil {
		return fmt.Errorf("번역문 디렉터리 생성 실패 (%s): %w", filepath.Dir(translatedPath), err)
	}

	// 2. 번역문 저장
	if err := os.WriteFile(translatedPath, content, 0o644); err != nil {
		return fmt.Errorf("번역문 파일 쓰기 실패 (%s): %w", translatedPath, err)
	}

	// 3. 매니페스트 로드
	mf, err := manifest.Load(site.ManifestDir())
	if err != nil {
		return fmt.Errorf("매니페스트 로드 실패 (siteID=%s): %w", site.ID, err)
	}

	// 4. TranslatedHash 갱신
	// 기존 Entry를 가져와 TranslatedHash만 교체한다. 기존 Entry가 없으면 새로 생성한다.
	entry, _ := mf.Get(pageURL)
	entry.TranslatedHash = sourceHash
	mf.Set(pageURL, entry)

	if err := manifest.Save(site.ManifestDir(), mf); err != nil {
		return fmt.Errorf("매니페스트 저장 실패 (siteID=%s): %w", site.ID, err)
	}

	return nil
}

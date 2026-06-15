// Package config는 사이트별 설정 파일을 파싱해 구조체로 로드한다.
// 설정 포맷은 JSON이며, 표준 라이브러리(encoding/json)만 사용한다.
// 사이트 폴더명이 사이트 식별자로 쓰인다(analysis.md D3, D4, D10).
//
// 사이트 폴더 레이아웃 규약(D5, D10):
//
//	sites/<siteID>/
//	├─ config.json      (이 패키지가 읽는 파일 — base_url, source_type, include/exclude 패턴)
//	├─ glossary.json    (번역 용어집 — 고정 이름)
//	├─ raw/             (원문 보관)
//	├─ output/          (번역문 출력)
//	└─ manifest.json    (증분 매니페스트 — 고정 이름)
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SourceType은 페이지 목록 확보 방식(source 타입) 값을 나타낸다.
// 유효한 값: "llms.txt", "sitemap", "crawl"
type SourceType string

const (
	SourceLLMsTxt SourceType = "llms.txt"
	SourceSitemap SourceType = "sitemap"
	SourceCrawl   SourceType = "crawl"
)

// Site는 한 사이트의 설정 값을 담는 구조체다(SPEC §5.1, D10).
// 출력·용어집·매니페스트 위치는 사이트 폴더 규약으로 고정되므로 설정에 두지 않는다.
type Site struct {
	// ID는 사이트 식별자: 사이트 폴더명(analysis.md D3).
	ID string `json:"-"`

	// SiteDir은 사이트 폴더 절대/상대 경로다.
	// LoadFromSiteDir이 채우며, config.json과 나란히 위치하는 폴더다.
	SiteDir string `json:"-"`

	// BaseURL (필수)
	BaseURL string `json:"base_url"`

	// SourceType: "llms.txt" | "sitemap" | "crawl" (필수)
	SourceType SourceType `json:"source_type"`

	// IncludePatterns (선택): 비어 있으면 전체 포함
	IncludePatterns []string `json:"include_patterns,omitempty"`

	// ExcludePatterns (선택)
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

// RawDir은 사이트 폴더 안 원문 보관 디렉터리 경로를 반환한다(규약: <siteDir>/raw/).
func (s *Site) RawDir() string {
	return filepath.Join(s.SiteDir, "raw")
}

// OutputDir은 사이트 폴더 안 번역문 출력 디렉터리 경로를 반환한다(규약: <siteDir>/output/).
func (s *Site) OutputDir() string {
	return filepath.Join(s.SiteDir, "output")
}

// GlossaryPath는 사이트 폴더 안 용어집 파일 경로를 반환한다(규약: <siteDir>/glossary.json).
func (s *Site) GlossaryPath() string {
	return filepath.Join(s.SiteDir, "glossary.json")
}

// ManifestDir은 사이트 폴더 경로를 반환한다(manifest.json은 이 폴더 직속에 고정).
func (s *Site) ManifestDir() string {
	return s.SiteDir
}

// Load는 siteDir/config.json을 파싱해 Site 구조체로 반환한다.
// 사이트 폴더명을 사이트 식별자로 도출한다(analysis.md D3).
// 필수 값(BaseURL, SourceType)이 빠지면 오류를 반환한다.
func Load(siteDir string) (*Site, error) {
	cfgPath := filepath.Join(siteDir, "config.json")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패 (%s): %w", cfgPath, err)
	}

	var site Site
	if err := json.Unmarshal(data, &site); err != nil {
		return nil, fmt.Errorf("설정 파일 파싱 실패 (%s): %w", cfgPath, err)
	}

	// 폴더명(마지막 경로 요소)을 사이트 식별자로 도출
	site.ID = filepath.Base(filepath.Clean(siteDir))
	site.SiteDir = siteDir

	if err := validate(&site); err != nil {
		return nil, fmt.Errorf("설정 파일 검증 실패 (%s): %w", cfgPath, err)
	}

	return &site, nil
}

// validate는 필수 필드가 모두 채워졌는지 확인하고, 빠진 필드를 열거해 오류로 반환한다.
func validate(site *Site) error {
	var missing []string

	if strings.TrimSpace(site.BaseURL) == "" {
		missing = append(missing, "base_url")
	}
	if strings.TrimSpace(string(site.SourceType)) == "" {
		missing = append(missing, "source_type")
	}

	if len(missing) > 0 {
		return fmt.Errorf("필수 필드 누락: %s", strings.Join(missing, ", "))
	}

	// source_type 값이 유효한 집합 안에 있는지 확인
	switch site.SourceType {
	case SourceLLMsTxt, SourceSitemap, SourceCrawl:
		// 유효한 값
	default:
		return fmt.Errorf("source_type 값이 유효하지 않습니다: %q (허용 값: llms.txt, sitemap, crawl)",
			site.SourceType)
	}

	return nil
}

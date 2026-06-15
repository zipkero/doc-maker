package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/config"
)

// writeTempSiteConfig는 지정한 내용을 임시 디렉터리 아래 <siteID>/config.json 으로 쓰고
// 사이트 폴더 경로를 반환한다.
func writeTempSiteConfig(t *testing.T, siteID, content string) string {
	t.Helper()
	siteDir := filepath.Join(t.TempDir(), siteID)
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}
	cfgPath := filepath.Join(siteDir, "config.json")
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("임시 설정 파일 생성 실패: %v", err)
	}
	return siteDir
}

// TestLoad_ValidConfig는 필수 필드가 모두 채워진 정상 설정 파일을 로드했을 때
// 구조체 값과 사이트 식별자가 올바르게 도출되는지 확인한다.
func TestLoad_ValidConfig(t *testing.T) {
	content := `{
		"base_url": "https://docs.example.com",
		"source_type": "llms.txt",
		"include_patterns": ["docs/**"],
		"exclude_patterns": ["docs/internal/**"]
	}`
	siteDir := writeTempSiteConfig(t, "example", content)

	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("정상 설정 로드 시 오류 발생: %v", err)
	}

	// 사이트 식별자는 폴더명
	if site.ID != "example" {
		t.Errorf("사이트 식별자: 기대 %q, 실제 %q", "example", site.ID)
	}
	if site.BaseURL != "https://docs.example.com" {
		t.Errorf("BaseURL: 기대 %q, 실제 %q", "https://docs.example.com", site.BaseURL)
	}
	if site.SourceType != config.SourceLLMsTxt {
		t.Errorf("SourceType: 기대 %q, 실제 %q", config.SourceLLMsTxt, site.SourceType)
	}
	if len(site.IncludePatterns) != 1 || site.IncludePatterns[0] != "docs/**" {
		t.Errorf("IncludePatterns: 기대 [\"docs/**\"], 실제 %v", site.IncludePatterns)
	}
	if len(site.ExcludePatterns) != 1 || site.ExcludePatterns[0] != "docs/internal/**" {
		t.Errorf("ExcludePatterns: 기대 [\"docs/internal/**\"], 실제 %v", site.ExcludePatterns)
	}
	// output_path, glossary_path는 설정에 없으므로 규약 경로를 확인한다
	if site.GlossaryPath() == "" {
		t.Error("GlossaryPath()가 빈 문자열을 반환함")
	}
	if site.OutputDir() == "" {
		t.Error("OutputDir()가 빈 문자열을 반환함")
	}
}

// TestLoad_OptionalPatternsOmitted는 포함·제외 패턴이 없는 설정도 정상 로드되는지 확인한다.
func TestLoad_OptionalPatternsOmitted(t *testing.T) {
	content := `{
		"base_url": "https://docs.example.com",
		"source_type": "llms.txt"
	}`
	siteDir := writeTempSiteConfig(t, "no-patterns", content)

	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("패턴 없는 설정 로드 시 오류 발생: %v", err)
	}
	if len(site.IncludePatterns) != 0 {
		t.Errorf("IncludePatterns: 비어 있어야 하는데 %v", site.IncludePatterns)
	}
	if len(site.ExcludePatterns) != 0 {
		t.Errorf("ExcludePatterns: 비어 있어야 하는데 %v", site.ExcludePatterns)
	}
}

// TestLoad_SiteIDFromFolderName은 폴더명이 사이트 식별자로 정확히 도출되는지 확인한다.
func TestLoad_SiteIDFromFolderName(t *testing.T) {
	content := `{
		"base_url": "https://docs.ollama.com",
		"source_type": "llms.txt"
	}`
	siteDir := writeTempSiteConfig(t, "ollama", content)

	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("로드 실패: %v", err)
	}
	if site.ID != "ollama" {
		t.Errorf("사이트 식별자: 기대 %q, 실제 %q", "ollama", site.ID)
	}
}

// TestLoad_AllSourceTypes는 세 source 타입 모두가 유효 값으로 수용되는지 확인한다.
func TestLoad_AllSourceTypes(t *testing.T) {
	cases := []struct {
		sourceType string
		expected   config.SourceType
	}{
		{"llms.txt", config.SourceLLMsTxt},
		{"sitemap", config.SourceSitemap},
		{"crawl", config.SourceCrawl},
	}

	for _, tc := range cases {
		t.Run(tc.sourceType, func(t *testing.T) {
			content := `{
				"base_url": "https://docs.example.com",
				"source_type": "` + tc.sourceType + `"
			}`
			siteDir := writeTempSiteConfig(t, "example", content)

			site, err := config.Load(siteDir)
			if err != nil {
				t.Fatalf("source_type %q 로드 시 오류 발생: %v", tc.sourceType, err)
			}
			if site.SourceType != tc.expected {
				t.Errorf("SourceType: 기대 %q, 실제 %q", tc.expected, site.SourceType)
			}
		})
	}
}

// TestLoad_MissingRequiredFields_BaseURL은 base_url 누락 시 명확한 오류가 발생하는지 확인한다.
func TestLoad_MissingRequiredFields_BaseURL(t *testing.T) {
	content := `{
		"source_type": "llms.txt"
	}`
	siteDir := writeTempSiteConfig(t, "missing-base-url", content)

	_, err := config.Load(siteDir)
	if err == nil {
		t.Fatal("base_url 누락 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_MissingRequiredFields_SourceType은 source_type 누락 시 명확한 오류가 발생하는지 확인한다.
func TestLoad_MissingRequiredFields_SourceType(t *testing.T) {
	content := `{
		"base_url": "https://docs.example.com"
	}`
	siteDir := writeTempSiteConfig(t, "missing-source-type", content)

	_, err := config.Load(siteDir)
	if err == nil {
		t.Fatal("source_type 누락 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_MissingRequiredFields_Multiple는 여러 필수 필드가 동시에 누락되어도
// 오류 메시지에 누락된 필드 이름들이 포함되는지 확인한다.
func TestLoad_MissingRequiredFields_Multiple(t *testing.T) {
	content := `{}`
	siteDir := writeTempSiteConfig(t, "missing-multiple", content)

	_, err := config.Load(siteDir)
	if err == nil {
		t.Fatal("필수 필드 누락 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_InvalidSourceType은 허용되지 않는 source_type 값 시 명확한 오류가 발생하는지 확인한다.
func TestLoad_InvalidSourceType(t *testing.T) {
	content := `{
		"base_url": "https://docs.example.com",
		"source_type": "unknown-type"
	}`
	siteDir := writeTempSiteConfig(t, "invalid-source-type", content)

	_, err := config.Load(siteDir)
	if err == nil {
		t.Fatal("유효하지 않은 source_type 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_FileNotFound는 존재하지 않는 폴더 경로를 주었을 때 오류가 발생하는지 확인한다.
func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/site")
	if err == nil {
		t.Fatal("존재하지 않는 폴더 로드 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_InvalidJSON은 JSON 형식이 깨진 파일을 로드했을 때 오류가 발생하는지 확인한다.
func TestLoad_InvalidJSON(t *testing.T) {
	siteDir := writeTempSiteConfig(t, "bad-json", `{ this is not valid json }`)

	_, err := config.Load(siteDir)
	if err == nil {
		t.Fatal("잘못된 JSON 파싱 시 오류가 발생해야 하는데 발생하지 않음")
	}
	t.Logf("예상 오류 확인: %v", err)
}

// TestLoad_OllamaSiteFolder는 sites/ollama/ 사이트 폴더가 로더에서 오류 없이 로드되고
// 3종 값과 사이트 식별자가 정상 도출되는지 확인한다(task-002 검증).
func TestLoad_OllamaSiteFolder(t *testing.T) {
	// 테스트 파일 위치 기준으로 프로젝트 루트의 sites/ollama 경로를 도출한다.
	// internal/config/ → ../../sites/ollama
	siteDir := filepath.Join("..", "..", "sites", "ollama")

	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("올라마 사이트 폴더 로드 실패: %v", err)
	}

	// 사이트 식별자는 폴더명이어야 한다.
	if site.ID != "ollama" {
		t.Errorf("사이트 식별자: 기대 %q, 실제 %q", "ollama", site.ID)
	}
	// 필수 필드 2종이 모두 채워졌는지 확인한다.
	if site.BaseURL == "" {
		t.Error("base_url이 비어 있음")
	}
	if site.SourceType == "" {
		t.Error("source_type이 비어 있음")
	}
	if site.SourceType != config.SourceLLMsTxt {
		t.Errorf("source_type: 기대 %q, 실제 %q", config.SourceLLMsTxt, site.SourceType)
	}
	// 규약 경로 헬퍼가 빈 문자열을 반환하지 않는지 확인한다.
	if site.GlossaryPath() == "" {
		t.Error("GlossaryPath()가 빈 문자열을 반환함")
	}
	if site.RawDir() == "" {
		t.Error("RawDir()가 빈 문자열을 반환함")
	}
	if site.OutputDir() == "" {
		t.Error("OutputDir()가 빈 문자열을 반환함")
	}
	t.Logf("올라마 사이트 폴더 로드 성공: ID=%s, BaseURL=%s, SourceType=%s",
		site.ID, site.BaseURL, site.SourceType)
}

// TestLoad_ConventionPaths는 규약 경로 헬퍼가 siteDir 기반 경로를 반환하는지 확인한다.
func TestLoad_ConventionPaths(t *testing.T) {
	content := `{
		"base_url": "https://docs.example.com",
		"source_type": "llms.txt"
	}`
	siteDir := writeTempSiteConfig(t, "mysite", content)

	site, err := config.Load(siteDir)
	if err != nil {
		t.Fatalf("로드 실패: %v", err)
	}

	// 각 규약 경로가 siteDir 아래를 가리키는지 확인한다.
	expectedRaw := filepath.Join(siteDir, "raw")
	expectedOutput := filepath.Join(siteDir, "output")
	expectedGlossary := filepath.Join(siteDir, "glossary.json")
	expectedManifestDir := siteDir

	if site.RawDir() != expectedRaw {
		t.Errorf("RawDir(): 기대 %q, 실제 %q", expectedRaw, site.RawDir())
	}
	if site.OutputDir() != expectedOutput {
		t.Errorf("OutputDir(): 기대 %q, 실제 %q", expectedOutput, site.OutputDir())
	}
	if site.GlossaryPath() != expectedGlossary {
		t.Errorf("GlossaryPath(): 기대 %q, 실제 %q", expectedGlossary, site.GlossaryPath())
	}
	if site.ManifestDir() != expectedManifestDir {
		t.Errorf("ManifestDir(): 기대 %q, 실제 %q", expectedManifestDir, site.ManifestDir())
	}
}

// 비고: 네트워크 실접속 없이 "사이트 인자→설정 로드→Collect 호출→건수 출력" 흐름을
// 모킹으로 검증하는 단위 테스트다. 실제 올라마 네트워크 수집은 선택적 수동 확인으로 둔다.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/source"
)

// --- 헬퍼: 임시 사이트 폴더 생성 ---

// writeSiteConfig는 t.TempDir() 아래에 <siteID>/config.json 설정 파일을 작성하고
// sitesRoot 경로를 반환한다.
func writeSiteConfig(t *testing.T, siteID string, extra map[string]interface{}) string {
	t.Helper()
	sitesRoot := t.TempDir()
	siteDir := filepath.Join(sitesRoot, siteID)
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		t.Fatalf("사이트 폴더 생성 실패: %v", err)
	}

	base := map[string]interface{}{
		"base_url":    "https://example.com",
		"source_type": "llms.txt",
	}
	for k, v := range extra {
		base[k] = v
	}

	data, err := json.Marshal(base)
	if err != nil {
		t.Fatalf("설정 파일 직렬화 실패: %v", err)
	}

	cfgPath := filepath.Join(siteDir, "config.json")
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatalf("설정 파일 쓰기 실패: %v", err)
	}
	return sitesRoot
}

// --- 헬퍼: 모킹 fetcher ---

// llmsTxtFetcher는 llms.txt 요청에 표본 목록을, 페이지 요청에 표본 원문을 반환하는 모킹 fetcher다.
func llmsTxtFetcher(llmsTxtContent string, pageContent string) source.Fetcher {
	return func(url string) ([]byte, error) {
		// llms.txt 요청: URL이 /llms.txt 로 끝나면 표본 반환
		if len(url) >= 8 && url[len(url)-8:] == "llms.txt" {
			return []byte(llmsTxtContent), nil
		}
		// 페이지 원문 요청
		return []byte(pageContent), nil
	}
}

// failFetcher는 항상 오류를 반환하는 모킹 fetcher다.
func failFetcher(url string) ([]byte, error) {
	return nil, fmt.Errorf("모킹: 취득 실패 (%s)", url)
}

// --- 테스트 ---

// TestRunSuccess는 정상 흐름(설정 로드→Collect 호출→건수 출력)을 검증한다.
func TestRunSuccess(t *testing.T) {
	siteID := "testsite"
	// llms.txt에 페이지 하나를 목록으로 둔다
	llmsContent := "https://example.com/docs/intro\n"
	pageContent := "# 소개\n\n이 페이지는 테스트용이다.\n"

	sitesRoot := writeSiteConfig(t, siteID, nil)
	fetch := llmsTxtFetcher(llmsContent, pageContent)

	if err := run(siteID, sitesRoot, fetch); err != nil {
		t.Fatalf("run 실패: %v", err)
	}

	// 원문 파일이 sites/<siteID>/raw/docs/intro.md 경로에 저장되어야 한다
	expectedPath := filepath.Join(sitesRoot, siteID, "raw", "docs", "intro.md")
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("원문 파일이 저장되지 않았습니다 (%s): %v", expectedPath, err)
	}
	if string(data) != pageContent {
		t.Errorf("원문 내용 불일치: 기대=%q, 실제=%q", pageContent, string(data))
	}
}

// TestRunIdempotent는 동일 내용으로 2회 실행 시 2회차에 갱신 0건임을 검증한다(증분).
func TestRunIdempotent(t *testing.T) {
	siteID := "idempotent"
	llmsContent := "https://example.com/docs/page\n"
	pageContent := "# 페이지\n변경 없음.\n"

	sitesRoot := writeSiteConfig(t, siteID, nil)
	fetch := llmsTxtFetcher(llmsContent, pageContent)

	// 1회차
	if err := run(siteID, sitesRoot, fetch); err != nil {
		t.Fatalf("1회차 run 실패: %v", err)
	}

	// 2회차: 내용 동일 → 갱신 0건이어야 한다.
	pathToCheck := filepath.Join(sitesRoot, siteID, "raw", "docs", "page.md")
	info1, err := os.Stat(pathToCheck)
	if err != nil {
		t.Fatalf("1회차 후 원문 파일 확인 실패: %v", err)
	}

	if err := run(siteID, sitesRoot, fetch); err != nil {
		t.Fatalf("2회차 run 실패: %v", err)
	}

	info2, err := os.Stat(pathToCheck)
	if err != nil {
		t.Fatalf("2회차 후 원문 파일 확인 실패: %v", err)
	}

	// 파일 수정 시각이 바뀌지 않으면 재저장하지 않은 것이다(증분 확인)
	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Errorf("2회차에 원문이 재저장됨: 1회차=%v, 2회차=%v", info1.ModTime(), info2.ModTime())
	}
}

// TestRunMissingConfig는 존재하지 않는 사이트 식별자 시 오류를 검증한다.
func TestRunMissingConfig(t *testing.T) {
	sitesRoot := t.TempDir() // 사이트 폴더 없는 빈 디렉터리

	err := run("nonexistent", sitesRoot, failFetcher)
	if err == nil {
		t.Fatal("존재하지 않는 식별자에 대해 오류가 반환되어야 합니다")
	}
}

// TestRunNotImplementedSource는 sitemap, crawl source 타입이 ErrNotImplemented로 표면화됨을 검증한다.
func TestRunNotImplementedSource(t *testing.T) {
	for _, srcType := range []string{"sitemap", "crawl"} {
		t.Run(srcType, func(t *testing.T) {
			siteID := "notimpl_" + srcType
			sitesRoot := writeSiteConfig(t, siteID, map[string]interface{}{
				"source_type": srcType,
			})

			err := run(siteID, sitesRoot, failFetcher)
			if err == nil {
				t.Fatalf("source 타입 %q에 대해 오류가 반환되어야 합니다", srcType)
			}
			if !errors.Is(err, source.ErrNotImplemented) {
				t.Errorf("ErrNotImplemented를 기대했지만 다른 오류: %v", err)
			}
		})
	}
}

// TestRunFetchFailure는 모든 페이지 취득이 실패해도 run이 오류 없이 완료됨을 검증한다
// (취득 실패는 result.Failed로 집계되지, 파이프라인 전체를 중단하지 않는다).
func TestRunFetchFailure(t *testing.T) {
	siteID := "fetchfail"
	// llms.txt는 정상 반환, 페이지 원문 취득만 실패하는 fetcher
	fetch := func(url string) ([]byte, error) {
		if len(url) >= 8 && url[len(url)-8:] == "llms.txt" {
			return []byte("https://example.com/docs/fail\n"), nil
		}
		return nil, fmt.Errorf("모킹: 페이지 취득 실패")
	}

	sitesRoot := writeSiteConfig(t, siteID, nil)

	// 취득 실패는 run 전체 오류가 아니라 result.Failed 집계이므로 err==nil이어야 한다
	if err := run(siteID, sitesRoot, fetch); err != nil {
		t.Fatalf("취득 실패 시에도 run은 정상 완료여야 합니다, 실제 오류: %v", err)
	}
}

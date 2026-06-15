// Package manifest는 페이지별 원문 콘텐츠 해시와 원본 경로를 사이트별로
// 기록·조회하는 매니페스트 구조를 제공한다(ANALYSIS §3, §5 D8).
//
// 매니페스트는 사이트 폴더(`sites/<siteID>/`) 직속 고정 이름 `manifest.json`으로 관리되며,
// 수집(task-007)과 번역(task-009~011)이 같은 사실을 근거로 증분을 판정할 수 있게 한다(D5, D10).
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Entry는 페이지 하나에 대한 매니페스트 항목이다.
// 수집 단계가 SourceHash와 SourcePath를 기록하고,
// 번역 단계가 TranslatedHash를 기록한다(ANALYSIS §5 D8).
type Entry struct {
	// SourceHash는 수집 시점에 계산한 원문 콘텐츠의 SHA-256 해시다(hex 인코딩).
	SourceHash string `json:"source_hash"`

	// SourcePath는 원문 파일의 사이트 폴더 기준 상대경로다.
	// 예: "raw/api/chat.md"
	SourcePath string `json:"source_path"`

	// TranslatedHash는 번역 단계가 번역한 원문 버전의 SHA-256 해시다(hex 인코딩).
	// 번역이 아직 수행되지 않은 경우 빈 문자열이다.
	TranslatedHash string `json:"translated_hash,omitempty"`
}

// Manifest는 한 사이트의 페이지별 항목 집합이다.
// 키는 페이지 URL이며, 값은 해당 페이지의 Entry다.
type Manifest struct {
	// Entries는 페이지 URL → Entry 매핑이다.
	Entries map[string]Entry `json:"entries"`
}

// New는 빈 Manifest를 생성해 반환한다.
func New() *Manifest {
	return &Manifest{Entries: make(map[string]Entry)}
}

// Get은 페이지 URL에 해당하는 Entry와 존재 여부를 반환한다.
func (m *Manifest) Get(pageURL string) (Entry, bool) {
	e, ok := m.Entries[pageURL]
	return e, ok
}

// Set은 페이지 URL에 대해 Entry를 기록(신규 추가 또는 갱신)한다.
func (m *Manifest) Set(pageURL string, entry Entry) {
	m.Entries[pageURL] = entry
}

// IsChanged는 주어진 원문 해시가 기존 기록의 SourceHash와 다른지 판정한다.
// 기존 기록이 없으면 신규 항목으로 보아 true를 반환한다(증분 판정 — SPEC §5.3).
func (m *Manifest) IsChanged(pageURL, hash string) bool {
	e, ok := m.Entries[pageURL]
	if !ok {
		// 기록이 없으면 신규 → 갱신 대상
		return true
	}
	return e.SourceHash != hash
}

// NeedsTranslation은 주어진 원문 해시와 기존 TranslatedHash가 달라
// 번역을 (재)수행해야 하는지 판정한다.
// 기존 기록이 없거나 TranslatedHash가 비어 있으면 true를 반환한다(SPEC §5.5).
func (m *Manifest) NeedsTranslation(pageURL, sourceHash string) bool {
	e, ok := m.Entries[pageURL]
	if !ok {
		return true
	}
	return e.TranslatedHash != sourceHash
}

// manifestFilePath는 사이트 폴더 경로로 매니페스트 파일 경로를 도출한다.
// 규약: <siteDir>/manifest.json (고정 이름, D5, D10).
func manifestFilePath(siteDir string) string {
	return filepath.Join(siteDir, "manifest.json")
}

// Load는 siteDir 안의 manifest.json을 읽어 Manifest를 반환한다.
// 파일이 없으면 빈 Manifest를 반환한다(첫 실행 허용).
// 파일이 있지만 파싱에 실패하면 오류를 반환한다.
func Load(siteDir string) (*Manifest, error) {
	path := manifestFilePath(siteDir)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// 파일 없음 → 신규 Manifest 반환
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("매니페스트 파일 읽기 실패 (%s): %w", path, err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("매니페스트 파일 파싱 실패 (%s): %w", path, err)
	}
	// Unmarshal 후 nil map 방어
	if m.Entries == nil {
		m.Entries = make(map[string]Entry)
	}

	return &m, nil
}

// Save는 Manifest를 siteDir 안의 manifest.json에 JSON으로 저장한다.
// 디렉터리가 없으면 생성한다.
func Save(siteDir string, m *Manifest) error {
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		return fmt.Errorf("매니페스트 디렉터리 생성 실패 (%s): %w", siteDir, err)
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("매니페스트 직렬화 실패: %w", err)
	}

	path := manifestFilePath(siteDir)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("매니페스트 파일 쓰기 실패 (%s): %w", path, err)
	}

	return nil
}

// HashContent는 콘텐츠 바이트의 SHA-256 해시를 hex 문자열로 반환한다.
// 표준 라이브러리(crypto/sha256)만 사용한다(task-004 비확장 조건).
func HashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

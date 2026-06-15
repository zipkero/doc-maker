package translator_test

import (
	"os"
	"path/filepath"
	"testing"

	"doc-maker/internal/translator"
)

// TestLoadGlossary_Normal은 유효한 JSON 용어집 파일을 로드해 매핑이 올바른지 확인한다.
func TestLoadGlossary_Normal(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "glossary.json")

	content := `{"model": "모델", "embedding": "임베딩", "inference": "추론"}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("용어집 파일 쓰기 실패: %v", err)
	}

	g, err := translator.LoadGlossary(path)
	if err != nil {
		t.Fatalf("LoadGlossary 오류: %v", err)
	}

	cases := map[string]string{
		"model":     "모델",
		"embedding": "임베딩",
		"inference": "추론",
	}
	for en, ko := range cases {
		if got := g[en]; got != ko {
			t.Errorf("g[%q]: 기대 %q, 실제 %q", en, ko, got)
		}
	}
}

// TestLoadGlossary_FileNotExist는 파일이 없을 때 빈 Glossary가 반환되는지 확인한다.
func TestLoadGlossary_FileNotExist(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "nonexistent.json")

	g, err := translator.LoadGlossary(path)
	if err != nil {
		t.Fatalf("파일 없을 때 오류 반환(기대: nil): %v", err)
	}
	if len(g) != 0 {
		t.Errorf("파일 없을 때 빈 Glossary 기대, 실제 길이 %d", len(g))
	}
}

// TestLoadGlossary_EmptyFile은 빈 파일일 때 빈 Glossary가 반환되는지 확인한다.
func TestLoadGlossary_EmptyFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "empty.json")

	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("빈 파일 쓰기 실패: %v", err)
	}

	g, err := translator.LoadGlossary(path)
	if err != nil {
		t.Fatalf("빈 파일일 때 오류 반환(기대: nil): %v", err)
	}
	if len(g) != 0 {
		t.Errorf("빈 파일일 때 빈 Glossary 기대, 실제 길이 %d", len(g))
	}
}

// TestLoadGlossary_InvalidJSON은 잘못된 JSON 파일에서 오류가 반환되는지 확인한다.
func TestLoadGlossary_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "bad.json")

	if err := os.WriteFile(path, []byte("{not valid json}"), 0o644); err != nil {
		t.Fatalf("파일 쓰기 실패: %v", err)
	}

	_, err := translator.LoadGlossary(path)
	if err == nil {
		t.Error("잘못된 JSON에서 오류가 반환되어야 함")
	}
}

// TestLoadGlossary_EmptyObject는 빈 JSON 오브젝트 "{}"일 때 빈 Glossary를 반환하는지 확인한다.
func TestLoadGlossary_EmptyObject(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "empty_obj.json")

	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatalf("파일 쓰기 실패: %v", err)
	}

	g, err := translator.LoadGlossary(path)
	if err != nil {
		t.Fatalf("빈 오브젝트 파싱 오류: %v", err)
	}
	if len(g) != 0 {
		t.Errorf("빈 오브젝트일 때 빈 Glossary 기대, 실제 길이 %d", len(g))
	}
}

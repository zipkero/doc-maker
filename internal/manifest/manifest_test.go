package manifest_test

import (
	"path/filepath"
	"testing"

	"doc-maker/internal/manifest"
)

// TestHashContent는 동일 콘텐츠에 대해 일관된 SHA-256 해시가 반환되는지 확인한다.
func TestHashContent(t *testing.T) {
	content := []byte("hello, world")

	h1 := manifest.HashContent(content)
	h2 := manifest.HashContent(content)

	if h1 == "" {
		t.Fatal("HashContent가 빈 문자열을 반환함")
	}
	if h1 != h2 {
		t.Errorf("같은 입력에 대해 다른 해시 반환: %q vs %q", h1, h2)
	}
}

// TestHashContent_DifferentContent는 다른 콘텐츠에 대해 서로 다른 해시가 반환되는지 확인한다.
func TestHashContent_DifferentContent(t *testing.T) {
	h1 := manifest.HashContent([]byte("content A"))
	h2 := manifest.HashContent([]byte("content B"))

	if h1 == h2 {
		t.Error("서로 다른 콘텐츠에 대해 같은 해시가 반환되어서는 안 됨")
	}
}

// TestHashContent_EmptyBytes는 빈 바이트 슬라이스에 대해서도 유효한 해시가 반환되는지 확인한다.
func TestHashContent_EmptyBytes(t *testing.T) {
	h := manifest.HashContent([]byte{})
	if h == "" {
		t.Fatal("빈 입력에 대한 HashContent가 빈 문자열을 반환함")
	}
}

// TestNew는 빈 Manifest가 올바르게 초기화되는지 확인한다.
func TestNew(t *testing.T) {
	m := manifest.New()
	if m == nil {
		t.Fatal("New()가 nil을 반환함")
	}
	_, ok := m.Get("https://example.com/page")
	if ok {
		t.Error("빈 Manifest에서 항목이 조회되어서는 안 됨")
	}
}

// TestSet_Get는 Entry를 Set한 뒤 Get으로 동일 값이 조회되는지 확인한다.
func TestSet_Get(t *testing.T) {
	m := manifest.New()

	pageURL := "https://docs.example.com/guide"
	entry := manifest.Entry{
		SourceHash: "abc123",
		SourcePath: "raw/example/guide.md",
	}
	m.Set(pageURL, entry)

	got, ok := m.Get(pageURL)
	if !ok {
		t.Fatal("Set한 항목이 Get으로 조회되지 않음")
	}
	if got.SourceHash != entry.SourceHash {
		t.Errorf("SourceHash: 기대 %q, 실제 %q", entry.SourceHash, got.SourceHash)
	}
	if got.SourcePath != entry.SourcePath {
		t.Errorf("SourcePath: 기대 %q, 실제 %q", entry.SourcePath, got.SourcePath)
	}
}

// TestSet_TranslatedHash는 TranslatedHash가 Entry에 함께 기록·조회되는지 확인한다(D8 번역 증분).
func TestSet_TranslatedHash(t *testing.T) {
	m := manifest.New()

	pageURL := "https://docs.example.com/guide"
	entry := manifest.Entry{
		SourceHash:     "def456",
		SourcePath:     "raw/example/guide.md",
		TranslatedHash: "def456",
	}
	m.Set(pageURL, entry)

	got, ok := m.Get(pageURL)
	if !ok {
		t.Fatal("Set한 항목이 Get으로 조회되지 않음")
	}
	if got.TranslatedHash != entry.TranslatedHash {
		t.Errorf("TranslatedHash: 기대 %q, 실제 %q", entry.TranslatedHash, got.TranslatedHash)
	}
}

// TestIsChanged_NewPage는 기록이 없는 페이지는 IsChanged가 true를 반환하는지 확인한다.
func TestIsChanged_NewPage(t *testing.T) {
	m := manifest.New()
	if !m.IsChanged("https://docs.example.com/new-page", "somehash") {
		t.Error("기록 없는 페이지에 대해 IsChanged는 true여야 함")
	}
}

// TestIsChanged_SameHash는 기존 기록과 같은 해시면 IsChanged가 false를 반환하는지 확인한다.
func TestIsChanged_SameHash(t *testing.T) {
	m := manifest.New()
	pageURL := "https://docs.example.com/page"
	hash := "aabbcc"
	m.Set(pageURL, manifest.Entry{SourceHash: hash, SourcePath: "raw/page.md"})

	if m.IsChanged(pageURL, hash) {
		t.Error("같은 해시일 때 IsChanged는 false여야 함")
	}
}

// TestIsChanged_DifferentHash는 기존 기록과 다른 해시면 IsChanged가 true를 반환하는지 확인한다.
func TestIsChanged_DifferentHash(t *testing.T) {
	m := manifest.New()
	pageURL := "https://docs.example.com/page"
	m.Set(pageURL, manifest.Entry{SourceHash: "oldhash", SourcePath: "raw/page.md"})

	if !m.IsChanged(pageURL, "newhash") {
		t.Error("다른 해시일 때 IsChanged는 true여야 함")
	}
}

// TestNeedsTranslation_NewPage는 기록이 없는 페이지는 NeedsTranslation이 true를 반환하는지 확인한다.
func TestNeedsTranslation_NewPage(t *testing.T) {
	m := manifest.New()
	if !m.NeedsTranslation("https://docs.example.com/new-page", "somehash") {
		t.Error("기록 없는 페이지에 대해 NeedsTranslation은 true여야 함")
	}
}

// TestNeedsTranslation_NoTranslatedHash는 TranslatedHash가 비어 있으면 true를 반환하는지 확인한다.
func TestNeedsTranslation_NoTranslatedHash(t *testing.T) {
	m := manifest.New()
	pageURL := "https://docs.example.com/page"
	m.Set(pageURL, manifest.Entry{SourceHash: "abc", SourcePath: "raw/page.md"})

	if !m.NeedsTranslation(pageURL, "abc") {
		t.Error("TranslatedHash 없는 페이지에 대해 NeedsTranslation은 true여야 함")
	}
}

// TestNeedsTranslation_TranslatedMatchesCurrent는 TranslatedHash가 현재 원문 해시와 같으면
// NeedsTranslation이 false를 반환하는지 확인한다(이미 번역된 무변경 페이지).
func TestNeedsTranslation_TranslatedMatchesCurrent(t *testing.T) {
	m := manifest.New()
	pageURL := "https://docs.example.com/page"
	hash := "aabbcc"
	m.Set(pageURL, manifest.Entry{
		SourceHash:     hash,
		SourcePath:     "raw/page.md",
		TranslatedHash: hash,
	})

	if m.NeedsTranslation(pageURL, hash) {
		t.Error("TranslatedHash가 현재 해시와 같을 때 NeedsTranslation은 false여야 함")
	}
}

// TestNeedsTranslation_SourceChanged는 원문 해시가 바뀌면 NeedsTranslation이 true를 반환하는지 확인한다.
func TestNeedsTranslation_SourceChanged(t *testing.T) {
	m := manifest.New()
	pageURL := "https://docs.example.com/page"
	// 이전에 "oldhash"를 번역했음
	m.Set(pageURL, manifest.Entry{
		SourceHash:     "newhash",
		SourcePath:     "raw/page.md",
		TranslatedHash: "oldhash",
	})

	// 현재 원문 해시는 "newhash"로 바뀜
	if !m.NeedsTranslation(pageURL, "newhash") {
		t.Error("원문 해시가 바뀌었을 때 NeedsTranslation은 true여야 함")
	}
}

// TestSaveLoad_RoundTrip는 Manifest를 Save한 뒤 Load하면 동일 내용이 보존되는지 확인한다(SPEC §5.3).
func TestSaveLoad_RoundTrip(t *testing.T) {
	siteDir := t.TempDir()

	original := manifest.New()
	original.Set("https://docs.example.com/page1", manifest.Entry{
		SourceHash: "hash1",
		SourcePath: "raw/example/page1.md",
	})
	original.Set("https://docs.example.com/page2", manifest.Entry{
		SourceHash:     "hash2",
		SourcePath:     "raw/example/page2.md",
		TranslatedHash: "hash2",
	})

	if err := manifest.Save(siteDir, original); err != nil {
		t.Fatalf("Save 실패: %v", err)
	}

	loaded, err := manifest.Load(siteDir)
	if err != nil {
		t.Fatalf("Load 실패: %v", err)
	}

	// page1 확인
	e1, ok := loaded.Get("https://docs.example.com/page1")
	if !ok {
		t.Fatal("Load 후 page1 항목이 없음")
	}
	if e1.SourceHash != "hash1" {
		t.Errorf("page1 SourceHash: 기대 %q, 실제 %q", "hash1", e1.SourceHash)
	}
	if e1.SourcePath != "raw/example/page1.md" {
		t.Errorf("page1 SourcePath: 기대 %q, 실제 %q", "raw/example/page1.md", e1.SourcePath)
	}
	if e1.TranslatedHash != "" {
		t.Errorf("page1 TranslatedHash: 비어 있어야 하는데 %q", e1.TranslatedHash)
	}

	// page2 확인 (TranslatedHash 포함)
	e2, ok := loaded.Get("https://docs.example.com/page2")
	if !ok {
		t.Fatal("Load 후 page2 항목이 없음")
	}
	if e2.SourceHash != "hash2" {
		t.Errorf("page2 SourceHash: 기대 %q, 실제 %q", "hash2", e2.SourceHash)
	}
	if e2.TranslatedHash != "hash2" {
		t.Errorf("page2 TranslatedHash: 기대 %q, 실제 %q", "hash2", e2.TranslatedHash)
	}
}

// TestLoad_NoFile은 매니페스트 파일이 없으면 빈 Manifest가 반환되는지 확인한다(첫 실행).
func TestLoad_NoFile(t *testing.T) {
	siteDir := t.TempDir()
	m, err := manifest.Load(siteDir)
	if err != nil {
		t.Fatalf("파일 없을 때 Load는 오류 없이 빈 Manifest를 반환해야 함: %v", err)
	}
	if m == nil {
		t.Fatal("Load가 nil을 반환함")
	}
}

// TestSave_CreatesDirectory는 디렉터리가 없어도 Save가 디렉터리를 만드는지 확인한다.
func TestSave_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	// 존재하지 않는 하위 디렉터리를 지정한다.
	siteDir := filepath.Join(base, "sites", "mysite")

	m := manifest.New()
	m.Set("https://docs.example.com/page", manifest.Entry{
		SourceHash: "abc",
		SourcePath: "raw/page.md",
	})

	if err := manifest.Save(siteDir, m); err != nil {
		t.Fatalf("디렉터리 없을 때 Save 실패: %v", err)
	}

	// 저장 후 Load로 확인
	loaded, err := manifest.Load(siteDir)
	if err != nil {
		t.Fatalf("저장된 파일 Load 실패: %v", err)
	}
	_, ok := loaded.Get("https://docs.example.com/page")
	if !ok {
		t.Error("저장된 항목이 Load 후 없음")
	}
}

// TestSiteIsolation은 서로 다른 사이트 디렉터리의 매니페스트가 분리되어 관리되는지 확인한다(ANALYSIS D5).
func TestSiteIsolation(t *testing.T) {
	base := t.TempDir()
	siteDirA := filepath.Join(base, "sites", "siteA")
	siteDirB := filepath.Join(base, "sites", "siteB")

	// siteA에 항목 저장
	mA := manifest.New()
	mA.Set("https://siteA.com/page", manifest.Entry{SourceHash: "hashA", SourcePath: "raw/siteA/page.md"})
	if err := manifest.Save(siteDirA, mA); err != nil {
		t.Fatalf("siteA Save 실패: %v", err)
	}

	// siteB에 항목 저장
	mB := manifest.New()
	mB.Set("https://siteB.com/page", manifest.Entry{SourceHash: "hashB", SourcePath: "raw/siteB/page.md"})
	if err := manifest.Save(siteDirB, mB); err != nil {
		t.Fatalf("siteB Save 실패: %v", err)
	}

	// siteA 로드 후 siteB 항목이 없는지 확인
	loadedA, err := manifest.Load(siteDirA)
	if err != nil {
		t.Fatalf("siteA Load 실패: %v", err)
	}
	_, okB := loadedA.Get("https://siteB.com/page")
	if okB {
		t.Error("siteA Manifest에 siteB 항목이 존재해서는 안 됨")
	}

	// siteB 로드 후 siteA 항목이 없는지 확인
	loadedB, err := manifest.Load(siteDirB)
	if err != nil {
		t.Fatalf("siteB Load 실패: %v", err)
	}
	_, okA := loadedB.Get("https://siteA.com/page")
	if okA {
		t.Error("siteB Manifest에 siteA 항목이 존재해서는 안 됨")
	}
}

// TestSaveLoad_IsChanged는 Save → Load 후 IsChanged 판정이 정상 동작하는지 확인한다.
func TestSaveLoad_IsChanged(t *testing.T) {
	siteDir := t.TempDir()
	pageURL := "https://docs.example.com/page"
	hash := manifest.HashContent([]byte("page content"))

	m := manifest.New()
	m.Set(pageURL, manifest.Entry{SourceHash: hash, SourcePath: "raw/example/page.md"})

	if err := manifest.Save(siteDir, m); err != nil {
		t.Fatalf("Save 실패: %v", err)
	}

	loaded, err := manifest.Load(siteDir)
	if err != nil {
		t.Fatalf("Load 실패: %v", err)
	}

	// 같은 해시 → 변경 없음
	if loaded.IsChanged(pageURL, hash) {
		t.Error("Load 후 같은 해시에 대해 IsChanged는 false여야 함")
	}

	// 다른 해시 → 변경됨
	newHash := manifest.HashContent([]byte("updated content"))
	if !loaded.IsChanged(pageURL, newHash) {
		t.Error("Load 후 다른 해시에 대해 IsChanged는 true여야 함")
	}
}

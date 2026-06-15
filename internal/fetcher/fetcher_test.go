package fetcher_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"doc-maker/internal/fetcher"
)

// mockHTTPDoer는 테스트용 HTTPDoer다.
// responses 슬라이스에 순서대로 응답을 넣으면 Do() 호출마다 순차 반환한다.
// 호출 횟수를 CallCount로 추적한다.
type mockHTTPDoer struct {
	responses []mockResponse
	idx       int
	CallCount int
}

type mockResponse struct {
	statusCode int
	body       string
	err        error
}

func (m *mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	m.CallCount++
	if m.idx >= len(m.responses) {
		// 응답 목록 소진 시 마지막 응답을 반복한다
		m.idx = len(m.responses) - 1
	}
	r := m.responses[m.idx]
	m.idx++
	if r.err != nil {
		return nil, r.err
	}
	return &http.Response{
		StatusCode: r.statusCode,
		Body:       io.NopCloser(strings.NewReader(r.body)),
	}, nil
}

// testConfig는 테스트에서 sleep 없이 실행되도록 Delay=0으로 설정한다.
func testConfig(maxRetries int) fetcher.Config {
	return fetcher.Config{
		MaxRetries:    maxRetries,
		Delay:         0,             // sleep 없이 빠르게 실행
		BackoffFactor: 0,
	}
}

// TestFetch_Success는 200 OK 응답에서 원문이 정상 반환되는지 확인한다.
func TestFetch_Success(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 200, body: "# Hello World"},
		},
	}
	f := fetcher.New(testConfig(3), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err != nil {
		t.Fatalf("정상 응답에서 오류가 발생해서는 안 됨: %v", err)
	}
	if string(data) != "# Hello World" {
		t.Errorf("원문 내용: got %q, want %q", string(data), "# Hello World")
	}
	if mock.CallCount != 1 {
		t.Errorf("정상 응답 시 요청 횟수: got %d, want 1", mock.CallCount)
	}
}

// TestFetch_TransientErrorRetries는 5xx 응답에서 재시도가 일어나는지 확인한다.
// 처음 두 번은 500, 세 번째에 200을 반환 → 최종 성공을 기대한다.
func TestFetch_TransientErrorRetries(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 500, body: ""},
			{statusCode: 500, body: ""},
			{statusCode: 200, body: "# Recovered"},
		},
	}
	f := fetcher.New(testConfig(3), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err != nil {
		t.Fatalf("재시도 후 성공해야 함: %v", err)
	}
	if string(data) != "# Recovered" {
		t.Errorf("원문 내용: got %q, want %q", string(data), "# Recovered")
	}
	// 500 두 번 + 200 한 번 = 3회 요청
	if mock.CallCount != 3 {
		t.Errorf("요청 횟수: got %d, want 3", mock.CallCount)
	}
}

// TestFetch_NetworkErrorRetries는 네트워크 오류에서 재시도가 일어나는지 확인한다.
// 처음 두 번은 네트워크 오류, 세 번째에 200을 반환 → 최종 성공을 기대한다.
func TestFetch_NetworkErrorRetries(t *testing.T) {
	netErr := fmt.Errorf("연결 거부")
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{err: netErr},
			{err: netErr},
			{statusCode: 200, body: "# Recovered"},
		},
	}
	f := fetcher.New(testConfig(3), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err != nil {
		t.Fatalf("네트워크 오류 후 재시도 성공해야 함: %v", err)
	}
	if string(data) != "# Recovered" {
		t.Errorf("원문 내용: got %q, want %q", string(data), "# Recovered")
	}
	if mock.CallCount != 3 {
		t.Errorf("요청 횟수: got %d, want 3", mock.CallCount)
	}
}

// TestFetch_FinalFailure는 모든 재시도를 소진했을 때 nil 콘텐츠와 오류를 반환하는지 확인한다.
// nil 콘텐츠를 받으면 호출자는 기존 원문을 보존해야 한다(SPEC §3 원문 보존 계약).
func TestFetch_FinalFailure(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 503, body: ""},
			{statusCode: 503, body: ""},
			{statusCode: 503, body: ""},
			{statusCode: 503, body: ""},
		},
	}
	// MaxRetries=3 → 최초 1회 + 재시도 3회 = 총 4회 시도
	f := fetcher.New(testConfig(3), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err == nil {
		t.Fatal("최종 실패 시 오류를 반환해야 함")
	}
	// 콘텐츠 미반환 계약: nil이어야 함 — 호출자가 기존 원문을 덮어쓰지 않도록
	if data != nil {
		t.Errorf("최종 실패 시 콘텐츠는 nil이어야 함 (원문 보존 계약): got %q", string(data))
	}
	if mock.CallCount != 4 {
		t.Errorf("요청 횟수: got %d, want 4 (최초 1 + 재시도 3)", mock.CallCount)
	}
}

// TestFetch_PermanentFailure_404는 4xx 응답에서 재시도 없이 즉시 실패하는지 확인한다.
func TestFetch_PermanentFailure_404(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 404, body: ""},
		},
	}
	f := fetcher.New(testConfig(3), mock)

	data, err := f.Fetch("https://example.com/docs/missing")
	if err == nil {
		t.Fatal("4xx 응답 시 오류를 반환해야 함")
	}
	if data != nil {
		t.Errorf("4xx 실패 시 콘텐츠는 nil이어야 함: got %q", string(data))
	}
	// 4xx는 영구 실패 — 재시도 없이 1회만 요청해야 한다
	if mock.CallCount != 1 {
		t.Errorf("4xx 시 요청 횟수: got %d, want 1 (재시도 없음)", mock.CallCount)
	}
	// ErrPermanent로 감싸져 있어야 한다
	if !errors.Is(err, fetcher.ErrPermanent) {
		t.Errorf("4xx 오류는 ErrPermanent를 감싸야 함: 실제 오류 = %v", err)
	}
}

// TestFetch_PermanentFailure_401은 401 Unauthorized도 영구 실패로 처리되는지 확인한다.
func TestFetch_PermanentFailure_401(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 401, body: ""},
		},
	}
	f := fetcher.New(testConfig(3), mock)

	_, err := f.Fetch("https://example.com/docs/private")
	if !errors.Is(err, fetcher.ErrPermanent) {
		t.Errorf("401 오류는 ErrPermanent를 감싸야 함: 실제 오류 = %v", err)
	}
	if mock.CallCount != 1 {
		t.Errorf("401 시 요청 횟수: got %d, want 1 (재시도 없음)", mock.CallCount)
	}
}

// TestFetch_NetworkErrorFinalFailure는 네트워크 오류가 MaxRetries를 초과하면
// nil 콘텐츠와 오류를 반환하는지 확인한다.
func TestFetch_NetworkErrorFinalFailure(t *testing.T) {
	netErr := fmt.Errorf("타임아웃")
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{err: netErr},
			{err: netErr},
		},
	}
	// MaxRetries=1 → 최초 1회 + 재시도 1회 = 총 2회
	f := fetcher.New(testConfig(1), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err == nil {
		t.Fatal("최종 실패 시 오류를 반환해야 함")
	}
	if data != nil {
		t.Errorf("최종 실패 시 콘텐츠는 nil이어야 함: got %q", string(data))
	}
	if mock.CallCount != 2 {
		t.Errorf("요청 횟수: got %d, want 2 (최초 1 + 재시도 1)", mock.CallCount)
	}
}

// TestFetch_ZeroRetries는 MaxRetries=0일 때 재시도 없이 1회만 요청하는지 확인한다.
func TestFetch_ZeroRetries(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 503, body: ""},
		},
	}
	f := fetcher.New(testConfig(0), mock)

	data, err := f.Fetch("https://example.com/docs/page")
	if err == nil {
		t.Fatal("실패 시 오류를 반환해야 함")
	}
	if data != nil {
		t.Errorf("실패 시 콘텐츠는 nil이어야 함")
	}
	if mock.CallCount != 1 {
		t.Errorf("MaxRetries=0 시 요청 횟수: got %d, want 1", mock.CallCount)
	}
}

// TestFetch_FetcherCompatibility는 Fetcher.Fetch 메서드 시그니처가
// source.Fetcher 함수 타입(func(string)([]byte,error))과 호환되는지 확인한다.
// source.Fetcher = func(url string)([]byte, error) 와 동일 시그니처이므로
// f.Fetch를 source.Fetcher 변수에 대입 가능해야 한다.
func TestFetch_FetcherCompatibility(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 200, body: "content"},
		},
	}
	f := fetcher.New(testConfig(0), mock)

	// source.Fetcher 타입과의 호환성: func(string)([]byte,error) 에 대입
	var fn func(string) ([]byte, error) = f.Fetch
	data, err := fn("https://example.com/page")
	if err != nil {
		t.Fatalf("호환성 확인 중 오류: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("데이터: got %q, want %q", string(data), "content")
	}
}

// TestFetch_DelayBetweenRetries는 Delay>0 설정 시 재시도 간 대기가 발생하는지
// 시간 측정으로 확인한다. 단위가 짧아 실제 sleep 여부만 대략 검증한다.
func TestFetch_DelayBetweenRetries(t *testing.T) {
	mock := &mockHTTPDoer{
		responses: []mockResponse{
			{statusCode: 500, body: ""},
			{statusCode: 200, body: "ok"},
		},
	}
	cfg := fetcher.Config{
		MaxRetries:    1,
		Delay:         50 * time.Millisecond, // 짧은 delay로 테스트 속도 확보
		BackoffFactor: 1.0,
	}
	f := fetcher.New(cfg, mock)

	start := time.Now()
	_, err := f.Fetch("https://example.com/docs/page")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("재시도 후 성공해야 함: %v", err)
	}
	// 재시도 1회 × delay 50ms = 최소 50ms 이상 소요되어야 한다
	if elapsed < 40*time.Millisecond {
		t.Errorf("Delay가 적용되지 않음: elapsed=%v, want >= 50ms", elapsed)
	}
}

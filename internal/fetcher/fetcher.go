// Package fetcher는 rate limit(요청 간격·재시도)을 적용한 HTTP 원문 취득기를 제공한다.
// source.Fetcher 함수 타입을 내부 HTTP 호출로 주입받아 단위 테스트에서 네트워크 없이 검증 가능하다.
package fetcher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ErrPermanent는 4xx 응답처럼 재시도해도 회복될 가능성이 없는 영구 실패를 나타낸다.
// 이 오류로 감싸진 경우 Fetch는 재시도 없이 즉시 실패를 반환한다.
var ErrPermanent = errors.New("영구 실패 (재시도 불가)")

// 재시도 대상 판정 기준:
//   - 네트워크 오류(err != nil, 응답 없음) → 재시도
//   - 5xx 서버 오류 → 재시도 (서버 일시 장애)
//   - 4xx 클라이언트 오류 → 즉시 실패 (요청 자체가 잘못됨)
//   - 3xx 리다이렉트 → http.Client가 자동 처리하므로 도달 시 즉시 실패로 간주
//   - 2xx(200 제외) → 즉시 실패

// Config는 rate limit 취득기의 동작 파라미터를 담는다.
// 테스트에서 Delay=0, BackoffFactor=0 으로 설정하면 sleep 없이 빠르게 실행된다.
type Config struct {
	// MaxRetries는 최초 요청 이후 재시도 최대 횟수다(0이면 재시도 없음).
	MaxRetries int
	// Delay는 요청 간 최소 대기 시간이다. 0이면 대기하지 않는다.
	Delay time.Duration
	// BackoffFactor는 재시도 간격에 곱해지는 배율이다(0이면 고정 간격).
	// 예: Delay=1s, BackoffFactor=2 → 1s → 2s → 4s
	BackoffFactor float64
}

// DefaultConfig는 외부 사이트 대상 합리적 기본값이다.
// 요청 간 1초, 최대 3회 재시도, 지수 backoff 배율 2.
var DefaultConfig = Config{
	MaxRetries:    3,
	Delay:         1 * time.Second,
	BackoffFactor: 2.0,
}

// HTTPDoer는 단일 HTTP 요청을 실행하는 인터페이스다.
// *http.Client가 이 인터페이스를 만족하며, 테스트에서 모킹으로 교체 가능하다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Fetcher는 rate limit·재시도를 적용해 URL의 원문을 취득하는 취득기다.
type Fetcher struct {
	cfg    Config
	client HTTPDoer
}

// New는 지정한 Config와 HTTP 클라이언트로 Fetcher를 만든다.
// client에 nil을 전달하면 기본 http.Client를 사용한다.
func New(cfg Config, client HTTPDoer) *Fetcher {
	if client == nil {
		client = &http.Client{}
	}
	return &Fetcher{cfg: cfg, client: client}
}

// NewDefault는 DefaultConfig와 기본 http.Client로 Fetcher를 만든다.
func NewDefault() *Fetcher {
	return New(DefaultConfig, nil)
}

// Fetch는 url의 원문을 취득해 반환한다.
//
// 동작:
//   - 정상(200 OK): 원문 바이트를 반환한다.
//   - 일시 오류(네트워크 오류 / 5xx): cfg.MaxRetries 횟수까지 재시도한다.
//     재시도 간격은 cfg.Delay * cfg.BackoffFactor^시도횟수 로 증가한다.
//   - 최종 실패: (nil, error)를 반환한다. 호출자는 nil 콘텐츠를 받으면
//     기존 원문을 그대로 유지해야 한다(SPEC §3 원문 보존 계약).
//   - 영구 실패(4xx): 재시도 없이 즉시 (nil, ErrPermanent 감싼 오류)를 반환한다.
//
// 요청 간격: 각 요청(최초 포함) 이전에 cfg.Delay만큼 대기한다.
// Delay=0이면 대기하지 않으므로 테스트에서 sleep 없이 실행할 수 있다.
func (f *Fetcher) Fetch(url string) ([]byte, error) {
	var lastErr error
	delay := f.cfg.Delay

	for attempt := 0; attempt <= f.cfg.MaxRetries; attempt++ {
		// 요청 간격 적용: 첫 번째 요청 포함, Delay>0 일 때만 대기
		if delay > 0 && attempt > 0 {
			time.Sleep(delay)
			// backoff: 다음 재시도 간격 계산
			if f.cfg.BackoffFactor > 0 {
				delay = time.Duration(float64(delay) * f.cfg.BackoffFactor)
			}
		}

		data, err := f.doRequest(url)
		if err == nil {
			return data, nil
		}

		// 영구 실패(4xx)는 즉시 반환 — 재시도해도 회복 불가
		if errors.Is(err, ErrPermanent) {
			return nil, err
		}

		// 일시 오류: 재시도 대상
		lastErr = err
	}

	// 모든 재시도 소진 후 최종 실패 — 콘텐츠를 반환하지 않아 호출자가 기존 원문을 보존할 수 있다
	return nil, fmt.Errorf("최종 실패 (재시도 %d회 소진): %w", f.cfg.MaxRetries, lastErr)
}

// doRequest는 단일 HTTP GET 요청을 실행하고 결과를 반환한다.
// 오류 분류:
//   - 네트워크 오류: 일반 error 반환 → 재시도 대상
//   - 5xx: 일반 error 반환 → 재시도 대상
//   - 4xx: ErrPermanent 감싼 error 반환 → 즉시 실패
//   - 200 외 2xx/3xx: ErrPermanent 감싼 error 반환 → 즉시 실패
func (f *Fetcher) doRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// 잘못된 URL 등 요청 생성 오류는 영구 실패로 처리한다
		return nil, fmt.Errorf("%w: 요청 생성 실패 (%s): %v", ErrPermanent, url, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		// 네트워크 오류(DNS 실패, 연결 거부, 타임아웃 등) → 재시도 대상
		return nil, fmt.Errorf("네트워크 오류 (%s): %w", url, err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		// 정상 응답
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			// 본문 읽기 실패는 일시 오류로 처리해 재시도한다
			return nil, fmt.Errorf("응답 본문 읽기 실패 (%s): %w", url, err)
		}
		return data, nil

	case resp.StatusCode >= 500:
		// 5xx 서버 오류 → 재시도 대상 (서버 일시 장애)
		return nil, fmt.Errorf("서버 오류 (%s): 상태 코드 %d", url, resp.StatusCode)

	case resp.StatusCode >= 400:
		// 4xx 클라이언트 오류 → 영구 실패 (요청 자체가 잘못됨: 404 Not Found, 401 Unauthorized 등)
		return nil, fmt.Errorf("%w: 클라이언트 오류 (%s): 상태 코드 %d", ErrPermanent, url, resp.StatusCode)

	default:
		// 3xx 리다이렉트(http.Client 자동 처리 한계 초과)·1xx 등 → 영구 실패
		return nil, fmt.Errorf("%w: 처리 불가 응답 (%s): 상태 코드 %d", ErrPermanent, url, resp.StatusCode)
	}
}

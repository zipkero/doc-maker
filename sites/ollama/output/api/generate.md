# 응답 생성

`POST /api/generate`

주어진 프롬프트에 대한 응답을 생성합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `prompt` | string |  | 모델이 응답을 생성할 입력 텍스트 |
| `suffix` | string |  | fill-in-the-middle 모델용. 사용자 프롬프트 뒤이자 모델 응답 앞에 오는 텍스트 |
| `images` | array |  | 이미지 입력을 지원하는 모델용 base64 인코딩 이미지 목록 |
| `format` | string \| object |  | 구조화된 출력 형식. 문자열 `"json"` 또는 JSON 스키마 객체 |
| `system` | string |  | 응답 생성에 사용할 시스템 프롬프트 |
| `stream` | boolean |  | `true`면 부분 응답을 스트리밍. 기본값 `true` |
| `think` | boolean \| string |  | `true`면 content와 별도로 사고 출력을 반환. 불리언 또는 `high`/`medium`/`low`(지원 모델 한정) |
| `raw` | boolean |  | `true`면 프롬프트 템플릿을 적용하지 않은 원시(raw) 응답을 반환 |
| `keep_alive` | string \| number |  | 모델 메모리 유지 시간(예: `5m`, 또는 즉시 해제는 `0`) |
| `options` | object |  | 생성을 제어하는 런타임 옵션(아래 참조) |
| `logprobs` | boolean |  | 출력 토큰의 로그 확률 반환 여부 |
| `top_logprobs` | integer |  | 각 토큰 위치에서 반환할 상위 토큰 수(`logprobs` 활성화 시) |

### 생성 옵션 (`options`)

| 필드 | 타입 | 설명 |
|---|---|---|
| `seed` | integer | 재현 가능한 출력을 위한 난수 시드 |
| `temperature` | number | 생성의 무작위성 제어(높을수록 무작위) |
| `top_k` | integer | 다음 토큰 후보를 가장 가능성 높은 K개로 제한 |
| `top_p` | number | 뉴클리어스 샘플링의 누적 확률 임계값 |
| `min_p` | number | 토큰 선택의 최소 확률 임계값 |
| `stop` | string \| array | 생성을 중단시킬 정지 시퀀스 |
| `num_ctx` | integer | 컨텍스트 길이(토큰 수) |
| `num_predict` | integer | 생성할 최대 토큰 수 |

## 응답

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 모델 이름 |
| `created_at` | string | 응답 생성 시각(ISO 8601) |
| `response` | string | 모델이 생성한 텍스트 응답 |
| `thinking` | string | 모델이 생성한 사고 출력 |
| `done` | boolean | 생성 완료 여부 |
| `done_reason` | string | 생성이 종료된 이유 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 프롬프트의 입력 토큰 수 |
| `prompt_eval_duration` | integer | 프롬프트 평가 시간(나노초) |
| `eval_count` | integer | 응답으로 생성된 출력 토큰 수 |
| `eval_duration` | integer | 토큰 생성 시간(나노초) |
| `logprobs` | array | 생성된 토큰의 로그 확률 정보(`logprobs` 활성화 시. 아래 참조) |

스트리밍 모드(`stream: true`, 기본값)에서는 위 형태의 부분 이벤트가 여러 번 전송되며, 각 이벤트의
`response`에는 해당 청크의 텍스트가, 마지막 이벤트의 `done`에는 `true`가 담깁니다.

### 로그 확률 객체 (`logprobs` 항목)

| 필드 | 타입 | 설명 |
|---|---|---|
| `token` | string | 토큰의 텍스트 표현 |
| `logprob` | number | 이 토큰의 로그 확률 |
| `bytes` | array | 토큰의 원시 바이트 표현 |
| `top_logprobs` | array | 이 위치에서 가장 가능성 높은 토큰들과 각각의 로그 확률 |

`top_logprobs`의 각 항목도 `token`·`logprob`·`bytes` 필드를 가집니다.

### 응답 예시

```json
{
  "model": "gemma4",
  "created_at": "2025-10-17T23:14:07.414671Z",
  "response": "Hello! How can I help you today?",
  "done": true,
  "done_reason": "stop",
  "total_duration": 174560334,
  "load_duration": 101397084,
  "prompt_eval_count": 11,
  "prompt_eval_duration": 13074791,
  "eval_count": 18,
  "eval_duration": 52479709
}
```

## 사용 예시

기본 요청:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "Why is the sky blue?"
}'
```

비스트리밍 요청:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "Why is the sky blue?",
  "stream": false
}'
```

옵션 지정:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "Why is the sky blue?",
  "options": {
    "temperature": 0.8,
    "top_p": 0.9,
    "seed": 42
  }
}'
```

구조화된 출력:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "What are the populations of the United States and Canada?",
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "countries": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "country": {"type": "string"},
            "population": {"type": "integer"}
          },
          "required": ["country", "population"]
        }
      }
    },
    "required": ["countries"]
  }
}'
```

이미지 입력:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "What is in this picture?",
  "images": ["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"]
}'
```

모델 로드:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4"
}'
```

모델 언로드:

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "keep_alive": 0
}'
```

> 원문: https://docs.ollama.com/api/generate

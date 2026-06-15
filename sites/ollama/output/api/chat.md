# 채팅 메시지 생성

`POST /api/chat`

사용자와 어시스턴트 간의 대화에서 다음 채팅 메시지를 생성합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `messages` | array | ✅ | 채팅 기록. 각 메시지 객체는 `role`과 `content`를 가집니다 |
| `tools` | array |  | 모델이 채팅 중 호출할 수 있는 함수 도구 목록 |
| `format` | string \| object |  | 응답 형식. `json` 또는 JSON 스키마를 지정할 수 있습니다 |
| `options` | object |  | 생성을 제어하는 런타임 옵션(아래 참조) |
| `stream` | boolean |  | 스트리밍 여부. 기본값 `true` |
| `think` | boolean \| string |  | 사고 과정 출력 여부. 불리언 또는 `high`/`medium`/`low`(지원 모델 한정) |
| `keep_alive` | string \| number |  | 모델 메모리 유지 시간(예: `5m`, 또는 즉시 해제는 `0`) |
| `logprobs` | boolean |  | 출력 토큰의 로그 확률 반환 여부 |
| `top_logprobs` | integer |  | 각 토큰 위치에서 반환할 상위 토큰 수(`logprobs` 활성화 시) |

### 메시지 객체 (`messages` 항목)

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `role` | string | ✅ | 작성 주체. `system` / `user` / `assistant` / `tool` |
| `content` | string | ✅ | 메시지 본문 |
| `images` | array |  | 멀티모달 모델용 인라인 이미지(base64) 목록 |
| `tool_calls` | array |  | 모델이 생성한 도구 호출 요청 |

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
| `model` | string | 이 메시지를 생성한 모델 이름 |
| `created_at` | string | 응답 생성 시각(ISO 8601) |
| `message` | object | 생성된 메시지. `role`, `content`, 선택적 `thinking`·`tool_calls`·`images` 포함 |
| `done` | boolean | 응답 완료 여부 |
| `done_reason` | string | 응답이 종료된 이유 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 프롬프트의 토큰 수 |
| `prompt_eval_duration` | integer | 프롬프트 평가 시간(나노초) |
| `eval_count` | integer | 응답으로 생성된 토큰 수 |
| `eval_duration` | integer | 토큰 생성 시간(나노초) |

스트리밍 모드(`stream: true`, 기본값)에서는 위 형태의 부분 이벤트가 여러 번 전송되며,
마지막 이벤트의 `done`이 `true`가 됩니다.

### 응답 예시

```json
{
  "model": "gemma4",
  "created_at": "2025-10-17T23:14:07.414671Z",
  "message": {
    "role": "assistant",
    "content": "Hello! How can I help you today?"
  },
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
curl http://localhost:11434/api/chat -d '{
  "model": "gemma4",
  "messages": [
    { "role": "user", "content": "why is the sky blue?" }
  ]
}'
```

비스트리밍 요청:

```bash
curl http://localhost:11434/api/chat -d '{
  "model": "gemma4",
  "messages": [
    { "role": "user", "content": "why is the sky blue?" }
  ],
  "stream": false
}'
```

도구 호출(tool calling):

```bash
curl http://localhost:11434/api/chat -d '{
  "model": "qwen3",
  "messages": [
    { "role": "user", "content": "What is the weather today in Paris?" }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": { "type": "string", "description": "The location to get the weather for, e.g. San Francisco, CA" },
            "format": { "type": "string", "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'", "enum": ["celsius", "fahrenheit"] }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}'
```

> 원문: https://docs.ollama.com/api/chat

# 모델 생성

`POST /api/create`

기존 모델을 기반으로 새 모델을 만듭니다. 시스템 프롬프트·템플릿·파라미터를 지정하거나 양자화를
적용할 수 있습니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 생성할 모델 이름 |
| `from` | string |  | 기반이 될 기존 모델 |
| `template` | string |  | 모델에 사용할 프롬프트 템플릿 |
| `license` | string \| array |  | 모델의 라이선스 문자열 또는 라이선스 목록 |
| `system` | string |  | 모델에 내장할 시스템 프롬프트 |
| `parameters` | object |  | 모델용 키-값 파라미터 |
| `messages` | array |  | 모델에 사용할 메시지 기록(아래 메시지 객체 참조) |
| `quantize` | string |  | 적용할 양자화 수준(예: `q4_K_M`, `q8_0`) |
| `stream` | boolean |  | 상태 업데이트를 스트리밍할지 여부. 기본값 `true` |

### 메시지 객체 (`messages` 항목)

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `role` | string | ✅ | 작성 주체. `system` / `user` / `assistant` / `tool` |
| `content` | string | ✅ | 메시지 본문 |
| `images` | array |  | 멀티모달 모델용 인라인 이미지(base64) 목록 |
| `tool_calls` | array |  | 모델이 생성한 도구 호출 요청(아래 참조) |

### 도구 호출 객체 (`tool_calls` 항목)

각 도구 호출은 `function` 객체를 가지며, 그 필드는 다음과 같습니다.

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `name` | string | ✅ | 호출할 함수 이름 |
| `description` | string |  | 함수가 하는 일 |
| `arguments` | object |  | 함수에 전달할 인자를 담은 JSON 객체 |

## 응답

생성 진행 상태가 업데이트로 전송됩니다. 비스트리밍 응답은 `application/json`으로 최종 상태만,
스트리밍 응답은 `application/x-ndjson`으로 여러 상태 이벤트를 반환합니다.

### 상태 응답 / 상태 이벤트

| 필드 | 타입 | 설명 |
|---|---|---|
| `status` | string | 현재 상태 메시지 |
| `digest` | string | 해당 상태와 관련된 콘텐츠 다이제스트(있는 경우) |
| `total` | integer | 작업에 예상되는 총 바이트 수 |
| `completed` | integer | 지금까지 전송된 바이트 수 |

비스트리밍 응답 예시:

```json
{
  "status": "success"
}
```

## 사용 예시

기본 요청:

```bash
curl http://localhost:11434/api/create -d '{
  "from": "gemma4",
  "model": "alpaca",
  "system": "You are Alpaca, a helpful AI assistant. You only answer with Emojis."
}'
```

기존 모델로부터 생성:

```bash
curl http://localhost:11434/api/create -d '{
  "model": "ollama",
  "from": "gemma4",
  "system": "You are Ollama the llama."
}'
```

양자화 적용:

```bash
curl http://localhost:11434/api/create -d '{
  "model": "llama3.1:8b-instruct-Q4_K_M",
  "from": "llama3.1:8b-instruct-fp16",
  "quantize": "q4_K_M"
}'
```

> 원문: https://docs.ollama.com/api/create

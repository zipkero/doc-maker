# 오류

## 상태 코드

엔드포인트는 요청의 성공 여부에 따라 HTTP 상태 줄에 적절한 HTTP 상태 코드를 반환합니다(예:
`HTTP/1.1 200 OK` 또는 `HTTP/1.1 400 Bad Request`). 자주 쓰이는 상태 코드는 다음과 같습니다.

| 상태 코드 | 의미 |
|---|---|
| `200` | 성공 |
| `400` | Bad Request (파라미터 누락, 잘못된 JSON 등) |
| `404` | Not Found (모델이 존재하지 않는 경우 등) |
| `429` | Too Many Requests (요청 한도 초과 등) |
| `500` | Internal Server Error |
| `502` | Bad Gateway (클라우드 모델에 연결할 수 없는 경우 등) |

## 오류 메시지

오류는 `application/json` 형식으로 반환되며, 오류 메시지는 `error` 속성에 담깁니다.

```json
{
  "error": "the model failed to generate a response"
}
```

## 스트리밍 중 발생하는 오류

스트리밍 도중 오류가 발생하면, 오류는 `application/x-ndjson` 형식의 객체로 `error` 속성에 담겨
반환됩니다. 이미 응답이 시작된 상태이므로 응답의 상태 코드는 변경되지 않습니다.

```json
{"model":"gemma4","created_at":"2025-10-26T17:21:21.196249Z","response":" Yes","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:21:21.207235Z","response":".","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:21:21.219166Z","response":"I","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:21:21.231094Z","response":"can","done":false}
{"error":"an error was encountered while running the model"}
```

> 원문: https://docs.ollama.com/api/errors

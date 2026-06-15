# 모델 push

`POST /api/push`

지정한 모델을 레지스트리에 업로드(게시)합니다. 기본적으로 진행 상태가 스트리밍으로
전달됩니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 게시할 모델 이름 |
| `insecure` | boolean |  | 안전하지 않은 연결로 게시 허용 |
| `stream` | boolean |  | 진행 상태 스트리밍 여부. 기본값 `true` |

## 응답

스트리밍(`stream: true`, 기본값) 모드에서는 진행 상태 이벤트가 NDJSON으로 여러 번
전송됩니다. 각 이벤트는 다음 필드를 가집니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `status` | string | 사람이 읽을 수 있는 상태 메시지 |
| `digest` | string | 해당 상태와 관련된 콘텐츠 다이제스트(있는 경우) |
| `total` | integer | 작업에 필요한 전체 바이트 수 |
| `completed` | integer | 현재까지 전송된 바이트 수 |

비스트리밍(`stream: false`) 모드에서는 완료 후 상태 메시지 하나만 반환됩니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `status` | string | 현재 상태 메시지 |

### 응답 예시

```json
{
  "status": "success"
}
```

## 사용 예시

기본 요청:

```bash
curl http://localhost:11434/api/push -d '{
  "model": "my-username/my-model"
}'
```

비스트리밍 요청:

```bash
curl http://localhost:11434/api/push -d '{
  "model": "my-username/my-model",
  "stream": false
}'
```

> 원문: https://docs.ollama.com/api/push

# 모델 삭제

`DELETE /api/delete`

모델을 삭제합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 삭제할 모델 이름 |

## 응답

삭제에 성공하면 `200` 상태 코드를 반환합니다.

## 사용 예시

모델 삭제:

```bash
curl -X DELETE http://localhost:11434/api/delete -d '{
  "model": "gemma4"
}'
```

> 원문: https://docs.ollama.com/api/delete

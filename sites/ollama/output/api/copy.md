# 모델 복사

`POST /api/copy`

기존 모델을 새 이름으로 복사합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `source` | string | ✅ | 복사할 원본 모델 이름 |
| `destination` | string | ✅ | 새로 만들 모델 이름 |

## 응답

복사에 성공하면 `200` 상태 코드를 반환합니다.

## 사용 예시

모델을 새 이름으로 복사:

```bash
curl http://localhost:11434/api/copy -d '{
  "source": "gemma4",
  "destination": "gemma4-backup"
}'
```

> 원문: https://docs.ollama.com/api/copy

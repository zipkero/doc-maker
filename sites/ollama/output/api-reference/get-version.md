# 버전 조회

`GET /api/version`

실행 중인 Ollama의 버전을 반환합니다.

## 응답

| 필드 | 타입 | 설명 |
|---|---|---|
| `version` | string | Ollama 버전 |

### 응답 예시

```json
{
  "version": "0.12.6"
}
```

## 사용 예시

```bash
curl http://localhost:11434/api/version
```

> 원문: https://docs.ollama.com/api-reference/get-version

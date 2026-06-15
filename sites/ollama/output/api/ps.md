# 실행 중인 모델 목록 조회

`GET /api/ps`

현재 메모리에 로드되어 실행 중인 모델 목록을 가져옵니다.

## 응답

| 필드 | 타입 | 설명 |
|---|---|---|
| `models` | array | 현재 실행 중인 모델 목록. 각 항목은 아래 모델 객체 형태입니다 |

### 모델 객체 (`models` 항목)

| 필드 | 타입 | 설명 |
|---|---|---|
| `name` | string | 실행 중인 모델 이름 |
| `model` | string | 실행 중인 모델 이름 |
| `size` | integer | 모델 크기(바이트) |
| `digest` | string | 모델의 SHA256 다이제스트 |
| `details` | object | 형식·계열 등 모델 세부 정보 |
| `expires_at` | string | 모델이 메모리에서 해제될 시각 |
| `size_vram` | integer | VRAM 사용량(바이트) |
| `context_length` | integer | 실행 중인 모델의 컨텍스트 길이 |

### 응답 예시

```json
{
  "models": [
    {
      "name": "gemma4",
      "model": "gemma4",
      "size": 6591830464,
      "digest": "c6eb396dbd5992bbe3f5cdb947e8bbc0ee413d7c17e2beaae69f5d569cf982eb",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "gemma4",
        "families": ["gemma4"],
        "parameter_size": "8.0B",
        "quantization_level": "Q4_K_M"
      },
      "expires_at": "2025-10-17T16:47:07.93355-07:00",
      "size_vram": 5333539264,
      "context_length": 4096
    }
  ]
}
```

## 사용 예시

실행 중인 모델 목록 조회:

```bash
curl http://localhost:11434/api/ps
```

> 원문: https://docs.ollama.com/api/ps

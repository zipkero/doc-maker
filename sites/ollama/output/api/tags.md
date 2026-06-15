# 모델 목록 조회

`GET /api/tags`

로컬에서 사용할 수 있는 모델 목록과 각 모델의 세부 정보를 가져옵니다.

## 응답

| 필드 | 타입 | 설명 |
|---|---|---|
| `models` | array | 로컬에서 사용 가능한 모델 목록. 각 항목은 아래 모델 요약 객체 형태입니다 |

### 모델 요약 객체 (`models` 항목)

| 필드 | 타입 | 설명 |
|---|---|---|
| `name` | string | 모델 이름 |
| `model` | string | 모델 이름 |
| `remote_model` | string | 원격 모델인 경우 업스트림 모델 이름 |
| `remote_host` | string | 원격 모델인 경우 업스트림 Ollama 호스트의 URL |
| `modified_at` | string | 마지막 수정 시각(ISO 8601 형식) |
| `size` | integer | 디스크상 모델 전체 크기(바이트) |
| `digest` | string | 모델 내용의 SHA256 다이제스트 식별자 |
| `details` | object | 모델의 형식·계열 등 추가 정보(아래 참조) |

### 모델 세부 정보 (`details`)

| 필드 | 타입 | 설명 |
|---|---|---|
| `format` | string | 모델 파일 형식(예: `gguf`) |
| `family` | string | 기본 모델 계열(예: `llama`) |
| `families` | array | 모델이 속한 모든 계열(해당하는 경우) |
| `parameter_size` | string | 대략적인 파라미터 수 라벨(예: `7B`, `13B`) |
| `quantization_level` | string | 사용된 양자화 수준(예: `Q4_0`) |

### 응답 예시

```json
{
  "models": [
    {
      "name": "gemma4",
      "model": "gemma4",
      "modified_at": "2025-10-03T23:34:03.409490317-07:00",
      "size": 9608350245,
      "digest": "c6eb396dbd5992bbe3f5cdb947e8bbc0ee413d7c17e2beaae69f5d569cf982eb",
      "details": {
        "format": "gguf",
        "family": "gemma4",
        "families": ["gemma4"],
        "parameter_size": "8.0B",
        "quantization_level": "Q4_K_M"
      }
    }
  ]
}
```

## 사용 예시

모델 목록 조회:

```bash
curl http://localhost:11434/api/tags
```

> 원문: https://docs.ollama.com/api/tags

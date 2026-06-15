# 모델 상세 정보 조회

`POST /api/show`

지정한 모델의 파라미터, 라이선스, 템플릿, 지원 기능, 메타데이터 등 상세 정보를 반환합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 조회할 모델 이름 |
| `verbose` | boolean |  | `true`이면 응답에 용량이 큰 상세 필드까지 포함합니다 |

## 응답

| 필드 | 타입 | 설명 |
|---|---|---|
| `parameters` | string | 텍스트로 직렬화된 모델 파라미터 설정 |
| `license` | string | 모델의 라이선스 |
| `modified_at` | string | 마지막 수정 시각(ISO 8601) |
| `details` | object | 모델 개요 정보(형식, 계열, 파라미터 크기, 양자화 수준 등) |
| `template` | string | 프롬프트를 구성할 때 모델이 사용하는 템플릿 |
| `capabilities` | array | 지원하는 기능 목록(예: `completion`, `vision`) |
| `model_info` | object | 추가 모델 메타데이터 |

### 응답 예시

```json
{
  "parameters": "temperature 0.7\nnum_ctx 2048",
  "license": "Gemma Terms of Use \n\nLast modified: February 21, 2024...",
  "capabilities": [
    "completion",
    "vision"
  ],
  "modified_at": "2025-08-14T15:49:43.634137516-07:00",
  "details": {
    "parent_model": "",
    "format": "gguf",
    "family": "gemma4",
    "families": [
      "gemma4"
    ],
    "parameter_size": "8.0B",
    "quantization_level": "Q4_K_M"
  },
  "model_info": {
    "gemma4.attention.head_count": 8,
    "gemma4.attention.head_count_kv": 2,
    "gemma4.block_count": 42,
    "gemma4.context_length": 131072,
    "gemma4.embedding_length": 2560,
    "general.architecture": "gemma4",
    "general.file_type": 15,
    "general.quantization_version": 2,
    "tokenizer.ggml.model": "llama",
    "tokenizer.ggml.bos_token_id": 2,
    "tokenizer.ggml.eos_token_id": 1
  }
}
```

> `model_info`는 모델 아키텍처에 따라 위 예시보다 훨씬 많은 키를 포함할 수 있습니다(어텐션 헤드 수, RoPE 설정, 토크나이저 메타데이터 등).

## 사용 예시

기본 요청:

```bash
curl http://localhost:11434/api/show -d '{
  "model": "gemma4"
}'
```

상세 요청(verbose):

```bash
curl http://localhost:11434/api/show -d '{
  "model": "gemma4",
  "verbose": true
}'
```

> 원문: https://docs.ollama.com/api-reference/show-model-details

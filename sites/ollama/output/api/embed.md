# 임베딩 생성

`POST /api/embed`

입력 텍스트를 표현하는 벡터 임베딩을 생성합니다.

## 요청 본문

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `input` | string \| array | ✅ | 임베딩을 생성할 텍스트 또는 텍스트 배열 |
| `truncate` | boolean |  | `true`면 컨텍스트 윈도를 초과하는 입력을 잘라냄. `false`면 오류 반환. 기본값 `true` |
| `dimensions` | integer |  | 생성할 임베딩의 차원 수 |
| `keep_alive` | string |  | 모델 메모리 유지 시간 |
| `options` | object |  | 생성을 제어하는 런타임 옵션(아래 참조) |

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
| `model` | string | 임베딩을 생성한 모델 이름 |
| `embeddings` | array | 벡터 임베딩 배열(각 항목이 숫자 배열) |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 임베딩 생성을 위해 처리한 입력 토큰 수 |

### 응답 예시

```json
{
  "model": "embeddinggemma",
  "embeddings": [
    [
      0.010071029,
      -0.0017594862,
      0.05007221,
      0.04692972,
      0.054916814,
      0.008599704,
      0.105441414,
      -0.025878139,
      0.12958129,
      0.031952348
    ]
  ],
  "total_duration": 14143917,
  "load_duration": 1019500,
  "prompt_eval_count": 8
}
```

## 사용 예시

기본 요청:

```bash
curl http://localhost:11434/api/embed -d '{
  "model": "embeddinggemma",
  "input": "Why is the sky blue?"
}'
```

여러 입력:

```bash
curl http://localhost:11434/api/embed -d '{
  "model": "embeddinggemma",
  "input": [
    "Why is the sky blue?",
    "Why is the grass green?"
  ]
}'
```

입력 잘라내기(truncation):

```bash
curl http://localhost:11434/api/embed -d '{
  "model": "embeddinggemma",
  "input": "Generate embeddings for this text",
  "truncate": true
}'
```

차원 지정:

```bash
curl http://localhost:11434/api/embed -d '{
  "model": "embeddinggemma",
  "input": "Generate embeddings for this text",
  "dimensions": 128
}'
```

> 원문: https://docs.ollama.com/api/embed

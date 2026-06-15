# 임베딩

텍스트를 숫자 벡터로 변환해 의미 기반 검색, 검색 증강(RAG), 유사도 비교 등에 활용합니다.

임베딩은 텍스트를 벡터로 바꿔 벡터 데이터베이스에 저장하거나, 코사인 유사도로 검색하거나, RAG 파이프라인에 사용할 수 있게 해 줍니다.
벡터 길이는 모델에 따라 다르며 보통 384~1024차원입니다.

## 추천 모델

* [embeddinggemma](https://ollama.com/library/embeddinggemma)
* [qwen3-embedding](https://ollama.com/library/qwen3-embedding)
* [all-minilm](https://ollama.com/library/all-minilm)

## 임베딩 생성

엔드포인트는 `POST /api/embed`이며, `model`과 `input`(문자열)을 전달합니다. 다음은 환경별 호출 방법입니다.

CLI에서 바로 생성:

```shell
ollama run embeddinggemma "Hello world"
```

텍스트를 파이프로 넘겨 생성할 수도 있습니다. 출력은 JSON 배열입니다.

```shell
echo "Hello world" | ollama run embeddinggemma
```

cURL:

```shell
curl -X POST http://localhost:11434/api/embed \
  -H "Content-Type: application/json" \
  -d '{
    "model": "embeddinggemma",
    "input": "The quick brown fox jumps over the lazy dog."
  }'
```

Python:

```python
import ollama

single = ollama.embed(
  model='embeddinggemma',
  input='The quick brown fox jumps over the lazy dog.'
)
print(len(single['embeddings'][0]))  # vector length
```

JavaScript:

```javascript
import ollama from 'ollama'

const single = await ollama.embed({
  model: 'embeddinggemma',
  input: 'The quick brown fox jumps over the lazy dog.',
})
console.log(single.embeddings[0].length) // vector length
```

> `/api/embed` 엔드포인트는 L2 정규화된(단위 길이) 벡터를 반환합니다.

## 임베딩 일괄 생성

`input`에 문자열 배열을 전달하면 여러 텍스트의 임베딩을 한 번에 생성합니다.

cURL:

```shell
curl -X POST http://localhost:11434/api/embed \
  -H "Content-Type: application/json" \
  -d '{
    "model": "embeddinggemma",
    "input": [
      "First sentence",
      "Second sentence",
      "Third sentence"
    ]
  }'
```

Python:

```python
import ollama

batch = ollama.embed(
  model='embeddinggemma',
  input=[
    'The quick brown fox jumps over the lazy dog.',
    'The five boxing wizards jump quickly.',
    'Jackdaws love my big sphinx of quartz.',
  ]
)
print(len(batch['embeddings']))  # number of vectors
```

JavaScript:

```javascript
import ollama from 'ollama'

const batch = await ollama.embed({
  model: 'embeddinggemma',
  input: [
    'The quick brown fox jumps over the lazy dog.',
    'The five boxing wizards jump quickly.',
    'Jackdaws love my big sphinx of quartz.',
  ],
})
console.log(batch.embeddings.length) // number of vectors
```

## 활용 팁

* 대부분의 의미 기반 검색에는 코사인 유사도를 사용하세요.
* 색인할 때와 질의할 때 같은 임베딩 모델을 사용하세요.

> 원문: https://docs.ollama.com/capabilities/embeddings

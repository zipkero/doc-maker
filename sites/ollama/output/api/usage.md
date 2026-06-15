# 사용량 지표

Ollama API 응답에는 성능과 모델 사용량을 측정할 수 있는 지표가 함께 담겨 있습니다.
주요 필드는 다음과 같습니다.

| 필드 | 설명 |
|---|---|
| `total_duration` | 응답을 생성하는 데 걸린 총 시간 |
| `load_duration` | 모델을 로드하는 데 걸린 시간 |
| `prompt_eval_count` | 처리된 입력 토큰 수 |
| `prompt_eval_duration` | 프롬프트를 평가하는 데 걸린 시간 |
| `eval_count` | 생성된 출력 토큰 수 |
| `eval_duration` | 출력 토큰을 생성하는 데 걸린 시간 |

모든 시간 값의 단위는 나노초입니다.

## 응답 예시

사용량 지표를 반환하는 엔드포인트의 응답 본문에는 위 필드가 포함됩니다. 예를 들어
`/api/generate`를 비스트리밍으로 호출하면 다음과 같은 응답을 받을 수 있습니다.

```json
{
  "model": "gemma4",
  "created_at": "2025-10-17T23:14:07.414671Z",
  "response": "Hello! How can I help you today?",
  "done": true,
  "done_reason": "stop",
  "total_duration": 174560334,
  "load_duration": 101397084,
  "prompt_eval_count": 11,
  "prompt_eval_duration": 13074791,
  "eval_count": 18,
  "eval_duration": 52479709
}
```

**스트리밍 응답**을 반환하는 엔드포인트에서는 사용량 필드가 `done`이 `true`인
마지막 청크에 포함됩니다.

> 원문: https://docs.ollama.com/api/usage

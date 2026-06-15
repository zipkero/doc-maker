# 사고 과정 (Thinking)

사고 기능을 지원하는 모델은 `thinking` 필드를 별도로 내보내, 추론 과정(reasoning trace)을 최종 답변과
분리해 줍니다.

이 기능은 모델의 추론 단계를 검토하거나, UI에서 모델이 "생각하는" 모습을 보여주거나, 최종 응답만
필요할 때 추론 과정을 완전히 숨기는 용도로 활용할 수 있습니다.

## 지원 모델

- [Qwen 3](https://ollama.com/library/qwen3)
- [GPT-OSS](https://ollama.com/library/gpt-oss) — `think` 레벨(`low`, `medium`, `high`)을 사용합니다.
  추론 과정을 완전히 비활성화할 수는 없습니다.
- [DeepSeek-v3.1](https://ollama.com/library/deepseek-v3.1)
- [DeepSeek R1](https://ollama.com/library/deepseek-r1)
- 최신 추가 모델은 [thinking models](https://ollama.com/search?c=thinking)에서 확인할 수 있습니다.

## API 호출에서 사고 기능 활성화

chat 또는 generate 요청에 `think` 필드를 설정합니다. 대부분의 모델은 불리언(`true`/`false`)을 받습니다.

GPT-OSS는 대신 `low`, `medium`, `high` 중 하나를 받아 추론 과정의 길이를 조절합니다.

추론 과정은 `message.thinking`(chat 엔드포인트) 또는 `thinking`(generate 엔드포인트) 필드에 담기며,
최종 답변은 `message.content` / `response`에 담깁니다.

cURL:

```shell
curl http://localhost:11434/api/chat -d '{
  "model": "qwen3",
  "messages": [{
    "role": "user",
    "content": "How many letter r are in strawberry?"
  }],
  "think": true,
  "stream": false
}'
```

Python:

```python
from ollama import chat

response = chat(
  model='qwen3',
  messages=[{'role': 'user', 'content': 'How many letter r are in strawberry?'}],
  think=True,
  stream=False,
)

print('Thinking:\n', response.message.thinking)
print('Answer:\n', response.message.content)
```

JavaScript:

```javascript
import ollama from 'ollama'

const response = await ollama.chat({
  model: 'deepseek-r1',
  messages: [{ role: 'user', content: 'How many letter r are in strawberry?' }],
  think: true,
  stream: false,
})

console.log('Thinking:\n', response.message.thinking)
console.log('Answer:\n', response.message.content)
```

> 참고: GPT-OSS는 `think`를 반드시 `"low"`, `"medium"`, `"high"` 중 하나로 설정해야 합니다. 이 모델에서는
> `true`/`false`를 전달해도 무시됩니다.

## 추론 과정 스트리밍

사고 과정 스트림은 답변 토큰보다 먼저 추론 토큰을 내보냅니다. 첫 `thinking` 청크를 감지하면 "사고 중"
영역을 렌더링하고, `message.content`가 도착하면 최종 응답으로 전환하면 됩니다.

Python:

```python
from ollama import chat

stream = chat(
  model='qwen3',
  messages=[{'role': 'user', 'content': 'What is 17 × 23?'}],
  think=True,
  stream=True,
)

in_thinking = False

for chunk in stream:
  if chunk.message.thinking and not in_thinking:
    in_thinking = True
    print('Thinking:\n', end='')

  if chunk.message.thinking:
    print(chunk.message.thinking, end='')
  elif chunk.message.content:
    if in_thinking:
      print('\n\nAnswer:\n', end='')
      in_thinking = False
    print(chunk.message.content, end='')

```

JavaScript:

```javascript
import ollama from 'ollama'

async function main() {
  const stream = await ollama.chat({
    model: 'qwen3',
    messages: [{ role: 'user', content: 'What is 17 × 23?' }],
    think: true,
    stream: true,
  })

  let inThinking = false

  for await (const chunk of stream) {
    if (chunk.message.thinking && !inThinking) {
      inThinking = true
      process.stdout.write('Thinking:\n')
    }

    if (chunk.message.thinking) {
      process.stdout.write(chunk.message.thinking)
    } else if (chunk.message.content) {
      if (inThinking) {
        process.stdout.write('\n\nAnswer:\n')
        inThinking = false
      }
      process.stdout.write(chunk.message.content)
    }
  }
}

main()
```

## CLI 빠른 참조

- 단일 실행에 사고 기능 활성화: `ollama run deepseek-r1 --think "Where should I visit in Lisbon?"`
- 사고 기능 비활성화: `ollama run deepseek-r1 --think=false "Summarize this article"`
- 사고 모델을 쓰되 추론 과정만 숨기기: `ollama run deepseek-r1 --hidethinking "Is 9.9 bigger or 9.11?"`
- 대화형 세션에서는 `/set think` 또는 `/set nothink`로 토글합니다.
- GPT-OSS는 레벨만 받습니다: `ollama run gpt-oss --think=low "Draft a headline"`
  (`low`는 필요에 따라 `medium`이나 `high`로 바꿉니다).

> 참고: 지원 모델에 대해서는 CLI와 API 모두에서 사고 기능이 기본으로 활성화되어 있습니다.

> 원문: https://docs.ollama.com/capabilities/thinking

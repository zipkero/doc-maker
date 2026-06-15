# 스트리밍

스트리밍을 사용하면 모델이 텍스트를 생성하는 즉시 화면에 렌더링할 수 있습니다.

REST API에서는 스트리밍이 기본으로 켜져 있지만, SDK에서는 기본으로 꺼져 있습니다.
SDK에서 스트리밍을 켜려면 `stream` 파라미터를 `True`로 설정하세요.

## 핵심 개념

1. **채팅**: 어시스턴트 메시지를 부분적으로 스트리밍합니다. 각 청크에 `content`가 들어 있어 메시지를 도착하는 대로 렌더링할 수 있습니다.
2. **사고 과정(thinking)**: 사고 기능을 지원하는 모델은 각 청크에서 일반 `content`와 함께 `thinking` 필드를 내보냅니다.
   스트리밍 청크에서 이 필드를 감지하면 최종 답변이 나오기 전의 추론 과정을 보여 주거나 숨길 수 있습니다.
3. **도구 호출(tool calling)**: 각 청크에서 스트리밍되는 `tool_calls`를 확인해 요청된 도구를 실행하고, 그 결과를 다시 대화에 덧붙입니다.

## 스트리밍 청크 처리

> 대화 기록을 유지하려면 부분 필드들을 누적해야 합니다. 특히 도구 호출에서는 모델의 사고 과정, 모델이 만든 도구 호출,
> 그리고 실행된 도구 결과를 다음 요청에서 모델에 다시 전달해야 하므로 누적이 중요합니다.

Python:

```python
from ollama import chat

stream = chat(
  model='qwen3',
  messages=[{'role': 'user', 'content': 'What is 17 × 23?'}],
  stream=True,
)

in_thinking = False
content = ''
thinking = ''
for chunk in stream:
  if chunk.message.thinking:
    if not in_thinking:
      in_thinking = True
      print('Thinking:\n', end='', flush=True)
    print(chunk.message.thinking, end='', flush=True)
    # accumulate the partial thinking 
    thinking += chunk.message.thinking
  elif chunk.message.content:
    if in_thinking:
      in_thinking = False
      print('\n\nAnswer:\n', end='', flush=True)
    print(chunk.message.content, end='', flush=True)
    # accumulate the partial content
    content += chunk.message.content

  # append the accumulated fields to the messages for the next request
  new_messages = [{ role: 'assistant', thinking: thinking, content: content }]
```

JavaScript:

```javascript
import ollama from 'ollama'

async function main() {
  const stream = await ollama.chat({
    model: 'qwen3',
    messages: [{ role: 'user', content: 'What is 17 × 23?' }],
    stream: true,
  })

  let inThinking = false
  let content = ''
  let thinking = ''

  for await (const chunk of stream) {
    if (chunk.message.thinking) {
      if (!inThinking) {
        inThinking = true
        process.stdout.write('Thinking:\n')
      }
      process.stdout.write(chunk.message.thinking)
      // accumulate the partial thinking
      thinking += chunk.message.thinking
    } else if (chunk.message.content) {
      if (inThinking) {
        inThinking = false
        process.stdout.write('\n\nAnswer:\n')
      }
      process.stdout.write(chunk.message.content)
      // accumulate the partial content
      content += chunk.message.content
    }
  }

  // append the accumulated fields to the messages for the next request
  new_messages = [{ role: 'assistant', thinking: thinking, content: content }]
}

main().catch(console.error)
```

> 원문: https://docs.ollama.com/capabilities/streaming

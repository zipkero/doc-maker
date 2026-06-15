# Anthropic 호환

Ollama는 [Anthropic Messages API](https://docs.anthropic.com/en/api/messages)와의 호환성을
제공하여, Claude Code 같은 도구를 비롯한 기존 애플리케이션을 Ollama에 손쉽게 연결할 수 있게
합니다.

## 사용법

### 환경 변수

Claude Code처럼 Anthropic API를 기대하는 도구에서 Ollama를 사용하려면 다음 환경 변수를
설정합니다.

```shell
export ANTHROPIC_AUTH_TOKEN=ollama  # required but ignored
export ANTHROPIC_BASE_URL=http://localhost:11434
```

### `/v1/messages` 기본 예시

Python:

```python
import anthropic

client = anthropic.Anthropic(
    base_url='http://localhost:11434',
    api_key='ollama',  # required but ignored
)

message = client.messages.create(
    model='qwen3-coder',
    max_tokens=1024,
    messages=[
        {'role': 'user', 'content': 'Hello, how are you?'}
    ]
)
print(message.content[0].text)
```

JavaScript:

```javascript
import Anthropic from "@anthropic-ai/sdk";

const anthropic = new Anthropic({
  baseURL: "http://localhost:11434",
  apiKey: "ollama", // required but ignored
});

const message = await anthropic.messages.create({
  model: "qwen3-coder",
  max_tokens: 1024,
  messages: [{ role: "user", content: "Hello, how are you?" }],
});

console.log(message.content[0].text);
```

cURL:

```shell
curl -X POST http://localhost:11434/v1/messages \
-H "Content-Type: application/json" \
-H "x-api-key: ollama" \
-H "anthropic-version: 2023-06-01" \
-d '{
  "model": "qwen3-coder",
  "max_tokens": 1024,
  "messages": [{ "role": "user", "content": "Hello, how are you?" }]
}'
```

### 스트리밍 예시

Python:

```python
import anthropic

client = anthropic.Anthropic(
    base_url='http://localhost:11434',
    api_key='ollama',
)

with client.messages.stream(
    model='qwen3-coder',
    max_tokens=1024,
    messages=[{'role': 'user', 'content': 'Count from 1 to 10'}]
) as stream:
    for text in stream.text_stream:
        print(text, end='', flush=True)
```

JavaScript:

```javascript
import Anthropic from "@anthropic-ai/sdk";

const anthropic = new Anthropic({
  baseURL: "http://localhost:11434",
  apiKey: "ollama",
});

const stream = await anthropic.messages.stream({
  model: "qwen3-coder",
  max_tokens: 1024,
  messages: [{ role: "user", content: "Count from 1 to 10" }],
});

for await (const event of stream) {
  if (
    event.type === "content_block_delta" &&
    event.delta.type === "text_delta"
  ) {
    process.stdout.write(event.delta.text);
  }
}
```

cURL:

```shell
curl -X POST http://localhost:11434/v1/messages \
-H "Content-Type: application/json" \
-d '{
  "model": "qwen3-coder",
  "max_tokens": 1024,
  "stream": true,
  "messages": [{ "role": "user", "content": "Count from 1 to 10" }]
}'
```

### 도구 호출 예시

Python:

```python
import anthropic

client = anthropic.Anthropic(
    base_url='http://localhost:11434',
    api_key='ollama',
)

message = client.messages.create(
    model='qwen3-coder',
    max_tokens=1024,
    tools=[
        {
            'name': 'get_weather',
            'description': 'Get the current weather in a location',
            'input_schema': {
                'type': 'object',
                'properties': {
                    'location': {
                        'type': 'string',
                        'description': 'The city and state, e.g. San Francisco, CA'
                    }
                },
                'required': ['location']
            }
        }
    ],
    messages=[{'role': 'user', 'content': "What's the weather in San Francisco?"}]
)

for block in message.content:
    if block.type == 'tool_use':
        print(f'Tool: {block.name}')
        print(f'Input: {block.input}')
```

JavaScript:

```javascript
import Anthropic from "@anthropic-ai/sdk";

const anthropic = new Anthropic({
  baseURL: "http://localhost:11434",
  apiKey: "ollama",
});

const message = await anthropic.messages.create({
  model: "qwen3-coder",
  max_tokens: 1024,
  tools: [
    {
      name: "get_weather",
      description: "Get the current weather in a location",
      input_schema: {
        type: "object",
        properties: {
          location: {
            type: "string",
            description: "The city and state, e.g. San Francisco, CA",
          },
        },
        required: ["location"],
      },
    },
  ],
  messages: [{ role: "user", content: "What's the weather in San Francisco?" }],
});

for (const block of message.content) {
  if (block.type === "tool_use") {
    console.log("Tool:", block.name);
    console.log("Input:", block.input);
  }
}
```

cURL:

```shell
curl -X POST http://localhost:11434/v1/messages \
-H "Content-Type: application/json" \
-d '{
  "model": "qwen3-coder",
  "max_tokens": 1024,
  "tools": [
    {
      "name": "get_weather",
      "description": "Get the current weather in a location",
      "input_schema": {
        "type": "object",
        "properties": {
          "location": {
            "type": "string",
            "description": "The city and state"
          }
        },
        "required": ["location"]
      }
    }
  ],
  "messages": [{ "role": "user", "content": "What is the weather in San Francisco?" }]
}'
```

## Claude Code와 함께 사용하기

[Claude Code](https://code.claude.com/docs/en/overview)는 백엔드로 Ollama를 사용하도록 설정할
수 있습니다.

### 권장 모델

코딩 용도라면 `glm-4.7`, `minimax-m2.1`, `qwen3-coder` 같은 모델을 권장합니다.

사용 전에 모델을 다운로드합니다.

```shell
ollama pull qwen3-coder
```

> 참고: Qwen 3 coder는 300억(30B) 파라미터 모델로, 원활하게 실행하려면 최소 24GB의 VRAM이
> 필요합니다. 컨텍스트 길이가 길수록 더 많은 메모리가 필요합니다.

```shell
ollama pull glm-4.7:cloud
```

### 빠른 설정

```shell
ollama launch claude
```

이 명령은 모델을 선택하도록 안내하고, Claude Code를 자동으로 설정한 뒤 실행합니다. 실행하지 않고
설정만 하려면 다음을 사용합니다.

```shell
ollama launch claude --config
```

### 수동 설정

환경 변수를 설정한 뒤 Claude Code를 실행합니다.

```shell
ANTHROPIC_AUTH_TOKEN=ollama ANTHROPIC_BASE_URL=http://localhost:11434 claude --model qwen3-coder
```

또는 셸 프로필에 환경 변수를 설정합니다.

```shell
export ANTHROPIC_AUTH_TOKEN=ollama
export ANTHROPIC_BASE_URL=http://localhost:11434
```

그런 다음 원하는 Ollama 모델로 Claude Code를 실행합니다.

```shell
claude --model qwen3-coder
```

## 엔드포인트

### `/v1/messages`

#### 지원 기능

| 기능 | 지원 |
|---|---|
| 메시지(Messages) | ✅ |
| 스트리밍 | ✅ |
| 시스템 프롬프트 | ✅ |
| 멀티턴 대화 | ✅ |
| 비전(이미지) | ✅ |
| 도구(함수 호출) | ✅ |
| 도구 결과(tool results) | ✅ |
| 사고/확장 사고(thinking) | ✅ |

#### 지원 요청 필드

| 필드 | 지원 | 비고 |
|---|---|---|
| `model` | ✅ | |
| `max_tokens` | ✅ | |
| `messages` | ✅ | 텍스트 `content`, 이미지 `content`(base64), 콘텐츠 블록 배열, `tool_use`/`tool_result`/`thinking` 블록 |
| `system` | ✅ | 문자열 또는 배열 |
| `stream` | ✅ | |
| `temperature` | ✅ | |
| `top_p` | ✅ | |
| `top_k` | ✅ | |
| `stop_sequences` | ✅ | |
| `tools` | ✅ | |
| `thinking` | ✅ | |
| `tool_choice` | ❌ | |
| `metadata` | ❌ | |

#### 지원 응답 필드

| 필드 | 지원 | 비고 |
|---|---|---|
| `id` | ✅ | |
| `type` | ✅ | |
| `role` | ✅ | |
| `model` | ✅ | |
| `content` | ✅ | text, tool_use, thinking 블록 |
| `stop_reason` | ✅ | end_turn, max_tokens, tool_use |
| `usage` | ✅ | input_tokens, output_tokens |

#### 스트리밍 이벤트

| 이벤트 | 지원 | 비고 |
|---|---|---|
| `message_start` | ✅ | |
| `content_block_start` | ✅ | |
| `content_block_delta` | ✅ | text_delta, input_json_delta, thinking_delta |
| `content_block_stop` | ✅ | |
| `message_delta` | ✅ | |
| `message_stop` | ✅ | |
| `ping` | ✅ | |
| `error` | ✅ | |

## 모델

Ollama는 로컬 모델과 클라우드 모델을 모두 지원합니다.

### 로컬 모델

사용 전에 로컬 모델을 다운로드합니다.

```shell
ollama pull qwen3-coder
```

권장 로컬 모델:

- `qwen3-coder` - 코딩 작업에 탁월
- `gpt-oss:20b` - 강력한 범용 모델

### 클라우드 모델

클라우드 모델은 별도로 다운로드하지 않아도 즉시 사용할 수 있습니다.

- `glm-4.7:cloud` - 고성능 클라우드 모델
- `minimax-m2.1:cloud` - 빠른 클라우드 모델

### 기본 모델 이름

`claude-3-5-sonnet`처럼 Anthropic 기본 모델 이름에 의존하는 도구를 사용한다면, `ollama cp`로
기존 모델 이름을 복사해 둡니다.

```shell
ollama cp qwen3-coder claude-3-5-sonnet
```

이후 `model` 필드에 이 새 모델 이름을 지정할 수 있습니다.

```shell
curl http://localhost:11434/v1/messages \
    -H "Content-Type: application/json" \
    -d '{
        "model": "claude-3-5-sonnet",
        "max_tokens": 1024,
        "messages": [
            {
                "role": "user",
                "content": "Hello!"
            }
        ]
    }'
```

## Anthropic API와의 차이점

### 동작 차이

- API 키는 받아들이지만 검증하지 않습니다.
- `anthropic-version` 헤더는 받아들이지만 사용하지 않습니다.
- 토큰 수는 기반 모델의 토크나이저를 바탕으로 한 근삿값입니다.

### 미지원 기능

다음 Anthropic API 기능은 현재 지원되지 않습니다.

| 기능 | 설명 |
|---|---|
| `/v1/messages/count_tokens` | 토큰 수 계산 엔드포인트 |
| `tool_choice` | 특정 도구 사용 강제 또는 도구 비활성화 |
| `metadata` | 요청 메타데이터(user_id) |
| 프롬프트 캐싱 | 접두부 캐싱을 위한 `cache_control` 블록 |
| Batches API | 비동기 배치 처리용 `/v1/messages/batches` |
| 인용(Citations) | `citations` 콘텐츠 블록 |
| PDF 지원 | PDF 파일을 담는 `document` 콘텐츠 블록 |
| 서버 전송 오류 | 스트리밍 중 `error` 이벤트(오류는 HTTP 상태로 반환됨) |

### 부분 지원

| 기능 | 상태 |
|---|---|
| 이미지 콘텐츠 | Base64 이미지는 지원, URL 이미지는 미지원 |
| 확장 사고(Extended thinking) | 기본 지원. `budget_tokens`는 받아들이지만 강제하지 않음 |

> 원문: https://docs.ollama.com/api/anthropic-compatibility

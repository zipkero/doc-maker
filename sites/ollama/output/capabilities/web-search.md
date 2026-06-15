# 웹 검색 (Web search)

Ollama의 웹 검색 API는 모델에 최신 정보를 보강해, 환각(hallucination)을 줄이고 정확도를 높이는 데
사용할 수 있습니다.

웹 검색은 REST API로 제공되며, Python·JavaScript 라이브러리에는 더 깊은 도구 통합이 포함되어 있습니다.
이를 통해 OpenAI의 gpt-oss 모델 같은 모델이 장시간에 걸친 리서치 작업을 수행할 수도 있습니다.

## 인증

Ollama 웹 검색 API에 접근하려면 [API key](https://ollama.com/settings/keys)를 발급받아야 합니다.
무료 Ollama 계정이 필요합니다.

## 웹 검색 API

`POST https://ollama.com/api/web_search`

단일 쿼리에 대해 웹 검색을 수행하고 관련 결과를 반환합니다.

### 요청 파라미터

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `query` | string | ✅ | 검색 쿼리 문자열 |
| `max_results` | integer |  | 반환할 최대 결과 수(기본값 5, 최대 10) |

### 응답

응답은 다음을 담은 객체입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `results` | array | 검색 결과 객체의 배열 |

각 검색 결과 객체(`results` 항목)는 다음 필드를 가집니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `title` | string | 웹 페이지의 제목 |
| `url` | string | 웹 페이지의 URL |
| `content` | string | 웹 페이지에서 추출한 관련 내용 일부 |

### 예시

> 참고: `OLLAMA_API_KEY` 환경 변수를 설정하거나, Authorization 헤더에 직접 전달해야 합니다.

cURL 요청:

```bash
curl https://ollama.com/api/web_search \
  --header "Authorization: Bearer $OLLAMA_API_KEY" \
	-d '{
	  "query":"what is ollama?"
	}'
```

응답:

```json
{
  "results": [
    {
      "title": "Ollama",
      "url": "https://ollama.com/",
      "content": "Cloud models are now available..."
    },
    {
      "title": "What is Ollama? Introduction to the AI model management tool",
      "url": "https://www.hostinger.com/tutorials/what-is-ollama",
      "content": "Ariffud M. 6min Read..."
    },
    {
      "title": "Ollama Explained: Transforming AI Accessibility and Language ...",
      "url": "https://www.geeksforgeeks.org/artificial-intelligence/ollama-explained-transforming-ai-accessibility-and-language-processing/",
      "content": "Data Science Data Science Projects Data Analysis..."
    }
  ]
}
```

Python 라이브러리:

```python
import ollama
response = ollama.web_search("What is Ollama?")
print(response)
```

출력 예시:

```python

results = [
    {
        "title": "Ollama",
        "url": "https://ollama.com/",
        "content": "Cloud models are now available in Ollama..."
    },
    {
        "title": "What is Ollama? Features, Pricing, and Use Cases - Walturn",
        "url": "https://www.walturn.com/insights/what-is-ollama-features-pricing-and-use-cases",
        "content": "Our services..."
    },
    {
        "title": "Complete Ollama Guide: Installation, Usage & Code Examples",
        "url": "https://collabnix.com/complete-ollama-guide-installation-usage-code-examples",
        "content": "Join our Discord Server..."
    }
]

```

더 많은 예시는 Ollama [Python 예제](https://github.com/ollama/ollama-python/blob/main/examples/web-search.py)를 참고하세요.

JavaScript 라이브러리:

```tsx
import { Ollama } from "ollama";

const client = new Ollama();
const results = await client.webSearch("what is ollama?");
console.log(JSON.stringify(results, null, 2));
```

출력 예시:

```json
{
  "results": [
    {
      "title": "Ollama",
      "url": "https://ollama.com/",
      "content": "Cloud models are now available..."
    },
    {
      "title": "What is Ollama? Introduction to the AI model management tool",
      "url": "https://www.hostinger.com/tutorials/what-is-ollama",
      "content": "Ollama is an open-source tool..."
    },
    {
      "title": "Ollama Explained: Transforming AI Accessibility and Language Processing",
      "url": "https://www.geeksforgeeks.org/artificial-intelligence/ollama-explained-transforming-ai-accessibility-and-language-processing/",
      "content": "Ollama is a groundbreaking..."
    }
  ]
}
```

더 많은 예시는 Ollama [JavaScript 예제](https://github.com/ollama/ollama-js/blob/main/examples/websearch/websearch-tools.ts)를 참고하세요.

## 웹 페이지 가져오기 API (Web fetch)

`POST https://ollama.com/api/web_fetch`

URL로 단일 웹 페이지를 가져와 그 내용을 반환합니다.

### 요청 파라미터

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `url` | string | ✅ | 가져올 URL |

### 응답

응답은 다음을 담은 객체입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `title` | string | 웹 페이지의 제목 |
| `content` | string | 웹 페이지의 본문 내용 |
| `links` | array | 페이지에서 발견된 링크의 배열 |

### 예시

cURL 요청:

```python
curl --request POST \
  --url https://ollama.com/api/web_fetch \
  --header "Authorization: Bearer $OLLAMA_API_KEY" \
  --header 'Content-Type: application/json' \
  --data '{
      "url": "ollama.com"
  }'
```

응답:

```json
{
  "title": "Ollama",
  "content": "[Cloud models](https://ollama.com/blog/cloud-models) are now available in Ollama...",
  "links": [
    "http://ollama.com/",
    "http://ollama.com/models",
    "https://github.com/ollama/ollama"
  ]

```

Python SDK:

```python
from ollama import web_fetch

result = web_fetch('https://ollama.com')
print(result)
```

결과:

```python
WebFetchResponse(
    title='Ollama',
    content='[Cloud models](https://ollama.com/blog/cloud-models) are now available in Ollama\n\n**Chat & build
with open models**\n\n[Download](https://ollama.com/download) [Explore
models](https://ollama.com/models)\n\nAvailable for macOS, Windows, and Linux',
    links=['https://ollama.com/', 'https://ollama.com/models', 'https://github.com/ollama/ollama']
)
```

JavaScript SDK:

```tsx
import { Ollama } from "ollama";

const client = new Ollama();
const fetchResult = await client.webFetch("https://ollama.com");
console.log(JSON.stringify(fetchResult, null, 2));
```

결과:

```json
{
  "title": "Ollama",
  "content": "[Cloud models](https://ollama.com/blog/cloud-models) are now available in Ollama...",
  "links": [
    "https://ollama.com/",
    "https://ollama.com/models",
    "https://github.com/ollama/ollama"
  ]
}
```

## 검색 에이전트 만들기

Ollama의 웹 검색 API를 도구로 사용해 간단한 검색 에이전트를 만들 수 있습니다. 다음 예시는 Alibaba의
Qwen 3 모델(4B 파라미터)을 사용합니다.

```bash
ollama pull qwen3:4b
```

```python
from ollama import chat, web_fetch, web_search

available_tools = {'web_search': web_search, 'web_fetch': web_fetch}

messages = [{'role': 'user', 'content': "what is ollama's new engine"}]

while True:
  response = chat(
    model='qwen3:4b',
    messages=messages,
    tools=[web_search, web_fetch],
    think=True
    )
  if response.message.thinking:
    print('Thinking: ', response.message.thinking)
  if response.message.content:
    print('Content: ', response.message.content)
  messages.append(response.message)
  if response.message.tool_calls:
    print('Tool calls: ', response.message.tool_calls)
    for tool_call in response.message.tool_calls:
      function_to_call = available_tools.get(tool_call.function.name)
      if function_to_call:
        args = tool_call.function.arguments
        result = function_to_call(**args)
        print('Result: ', str(result)[:200]+'...')
        # Result is truncated for limited context lengths
        messages.append({'role': 'tool', 'content': str(result)[:2000 * 4], 'tool_name': tool_call.function.name})
      else:
        messages.append({'role': 'tool', 'content': f'Tool {tool_call.function.name} not found', 'tool_name': tool_call.function.name})
  else:
    break
```

결과:

```
Thinking:  Okay, the user is asking about Ollama's new engine. I need to figure out what they're referring to. Ollama is a company that develops large language models, so maybe they've released a new model or an updated version of their existing engine....

Tool calls:  [ToolCall(function=Function(name='web_search', arguments={'max_results': 3, 'query': 'Ollama new engine'}))]
Result:  results=[WebSearchResult(content='# New model scheduling\n\n## September 23, 2025\n\nOllama now includes a significantly improved model scheduling system. Ahead of running a model, Ollama’s new engine

Thinking:  Okay, the user asked about Ollama's new engine. Let me look at the search results.

First result is from September 23, 2025, talking about new model scheduling. It mentions improved memory management, reduced crashes, better GPU utilization, and multi-GPU performance. Examples show speed improvements and accurate memory reporting. Supported models include gemma3, llama4, qwen3, etc...

Content:  Ollama has introduced two key updates to its engine, both released in 2025:

1. **Enhanced Model Scheduling (September 23, 2025)**
   - **Precision Memory Management**: Exact memory allocation reduces out-of-memory crashes and optimizes GPU utilization.
   - **Performance Gains**: Examples show significant speed improvements (e.g., 85.54 tokens/s vs 52.02 tokens/s) and full GPU layer utilization.
   - **Multi-GPU Support**: Improved efficiency across multiple GPUs, with accurate memory reporting via tools like `nvidia-smi`.
   - **Supported Models**: Includes `gemma3`, `llama4`, `qwen3`, `mistral-small3.2`, and more.

2. **Multimodal Engine (May 15, 2025)**
   - **Vision Support**: First-class support for vision models, including `llama4:scout` (109B parameters), `gemma3`, `qwen2.5vl`, and `mistral-small3.1`.
   - **Multimodal Tasks**: Examples include identifying animals in multiple images, answering location-based questions from videos, and document scanning.

These updates highlight Ollama's focus on efficiency, performance, and expanded capabilities for both text and vision tasks.
```

### 컨텍스트 길이와 에이전트

웹 검색 결과는 수천 개의 토큰을 반환할 수 있습니다. 모델의 컨텍스트 길이를 최소 약 32000 토큰 이상으로
늘리는 것을 권장합니다. 검색 에이전트는 전체 컨텍스트 길이에서 가장 잘 동작합니다.
[Ollama의 클라우드 모델](https://docs.ollama.com/cloud)은 전체 컨텍스트 길이로 실행됩니다.

## MCP 서버

[Python MCP 서버](https://github.com/ollama/ollama-python/blob/main/examples/web-search-mcp.py)를 통해
어떤 MCP 클라이언트에서든 웹 검색을 활성화할 수 있습니다.

### Cline

Ollama의 웹 검색은 MCP 서버 설정을 사용해 Cline에 쉽게 통합할 수 있습니다.

`Manage MCP Servers` > `Configure MCP Servers`로 이동한 뒤 다음 설정을 추가합니다.

```json
{
  "mcpServers": {
    "web_search_and_fetch": {
      "type": "stdio",
      "command": "uv",
      "args": ["run", "path/to/web-search-mcp.py"],
      "env": { "OLLAMA_API_KEY": "your_api_key_here" }
    }
  }
}
```

### Codex

Ollama는 OpenAI의 Codex 도구와도 잘 동작합니다. `~/.codex/config.toml`에 다음 설정을 추가합니다.

```python
[mcp_servers.web_search]
command = "uv"
args = ["run", "path/to/web-search-mcp.py"]
env = { "OLLAMA_API_KEY" = "your_api_key_here" }
```

### Goose

Ollama는 MCP 기능을 통해 Goose와 통합할 수 있습니다.

### 기타 통합

Ollama는 API 직접 연동, Python·JavaScript 라이브러리, OpenAI 호환 API, MCP 서버 통합 등을 통해 대부분의
도구에 통합할 수 있습니다.

> 원문: https://docs.ollama.com/capabilities/web-search

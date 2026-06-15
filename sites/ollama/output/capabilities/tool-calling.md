# 도구 호출 (Tool calling)

Ollama는 도구 호출(함수 호출이라고도 함)을 지원합니다. 이를 통해 모델이 도구를 호출하고 그 결과를
자신의 응답에 반영할 수 있습니다.

## 단일 도구 호출

도구 하나를 호출하고, 그 결과를 후속 요청에 포함합니다. "단발성(single-shot)" 도구 호출이라고도 합니다.

cURL:

```shell
curl -s http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "qwen3",
  "messages": [{"role": "user", "content": "What is the temperature in New York?"}],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_temperature",
        "description": "Get the current temperature for a city",
        "parameters": {
          "type": "object",
          "required": ["city"],
          "properties": {
            "city": {"type": "string", "description": "The name of the city"}
          }
        }
      }
    }
  ]
}'
```

단일 도구 결과를 담아 최종 응답을 생성:

```shell
curl -s http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "qwen3",
  "messages": [
    {"role": "user", "content": "What is the temperature in New York?"},
    {
      "role": "assistant",
      "tool_calls": [
        {
          "type": "function",
          "function": {
            "index": 0,
            "name": "get_temperature",
            "arguments": {"city": "New York"}
          }
        }
      ]
    },
    {"role": "tool", "tool_name": "get_temperature", "content": "22°C"}
  ],
  "stream": false
}'
```

Python: 먼저 Ollama Python SDK를 설치합니다.

```bash
# with pip
pip install ollama -U

# with uv
uv add ollama    
```

```python
from ollama import chat

def get_temperature(city: str) -> str:
  """Get the current temperature for a city
  
  Args:
    city: The name of the city

  Returns:
    The current temperature for the city
  """
  temperatures = {
    "New York": "22°C",
    "London": "15°C",
    "Tokyo": "18°C",
  }
  return temperatures.get(city, "Unknown")

messages = [{"role": "user", "content": "What is the temperature in New York?"}]

# pass functions directly as tools in the tools list or as a JSON schema
response = chat(model="qwen3", messages=messages, tools=[get_temperature], think=True)

messages.append(response.message)
if response.message.tool_calls:
  # only recommended for models which only return a single tool call
  call = response.message.tool_calls[0]
  result = get_temperature(**call.function.arguments)
  # add the tool result to the messages
  messages.append({"role": "tool", "tool_name": call.function.name, "content": str(result)})

  final_response = chat(model="qwen3", messages=messages, tools=[get_temperature], think=True)
  print(final_response.message.content)
```

JavaScript: 먼저 Ollama JavaScript 라이브러리를 설치합니다.

```bash
# with npm
npm i ollama

# with bun
bun i ollama
```

```typescript
import ollama from 'ollama'

function getTemperature(city: string): string {
  const temperatures: Record<string, string> = {
    'New York': '22°C',
    'London': '15°C',
    'Tokyo': '18°C',
  }
  return temperatures[city] ?? 'Unknown'
}

const tools = [
  {
    type: 'function',
    function: {
      name: 'get_temperature',
      description: 'Get the current temperature for a city',
      parameters: {
        type: 'object',
        required: ['city'],
        properties: {
          city: { type: 'string', description: 'The name of the city' },
        },
      },
    },
  },
]

const messages = [{ role: 'user', content: "What is the temperature in New York?" }]

const response = await ollama.chat({
  model: 'qwen3',
  messages,
  tools,
  think: true,
})

messages.push(response.message)
if (response.message.tool_calls?.length) {
  // only recommended for models which only return a single tool call
  const call = response.message.tool_calls[0]
  const args = call.function.arguments as { city: string }
  const result = getTemperature(args.city)
  // add the tool result to the messages
  messages.push({ role: 'tool', tool_name: call.function.name, content: result })

  // generate the final response
  const finalResponse = await ollama.chat({ model: 'qwen3', messages, tools, think: true })
  console.log(finalResponse.message.content)
}
```

## 병렬 도구 호출

여러 도구를 병렬로 호출하도록 요청한 뒤, 모든 도구 응답을 모아 모델에 다시 전달합니다.

cURL:

```shell
curl -s http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "qwen3",
  "messages": [{"role": "user", "content": "What are the current weather conditions and temperature in New York and London?"}],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_temperature",
        "description": "Get the current temperature for a city",
        "parameters": {
          "type": "object",
          "required": ["city"],
          "properties": {
            "city": {"type": "string", "description": "The name of the city"}
          }
        }
      }
    },
    {
      "type": "function",
      "function": {
        "name": "get_conditions",
        "description": "Get the current weather conditions for a city",
        "parameters": {
          "type": "object",
          "required": ["city"],
          "properties": {
            "city": {"type": "string", "description": "The name of the city"}
          }
        }
      }
    }
  ]
}'
```

여러 도구 결과를 담아 최종 응답을 생성:

```shell
curl -s http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "qwen3",
  "messages": [
    {"role": "user", "content": "What are the current weather conditions and temperature in New York and London?"},
    {
      "role": "assistant",
      "tool_calls": [
        {
          "type": "function",
          "function": {
            "index": 0,
            "name": "get_temperature",
            "arguments": {"city": "New York"}
          }
        },
        {
          "type": "function",
          "function": {
            "index": 1,
            "name": "get_conditions",
            "arguments": {"city": "New York"}
          }
        },
        {
          "type": "function",
          "function": {
            "index": 2,
            "name": "get_temperature",
            "arguments": {"city": "London"}
          }
        },
        {
          "type": "function",
          "function": {
            "index": 3,
            "name": "get_conditions",
            "arguments": {"city": "London"}
          }
        }
      ]
    },
    {"role": "tool", "tool_name": "get_temperature", "content": "22°C"},
    {"role": "tool", "tool_name": "get_conditions", "content": "Partly cloudy"},
    {"role": "tool", "tool_name": "get_temperature", "content": "15°C"},
    {"role": "tool", "tool_name": "get_conditions", "content": "Rainy"}
  ],
  "stream": false
}'
```

Python:

```python
from ollama import chat

def get_temperature(city: str) -> str:
  """Get the current temperature for a city
  
  Args:
    city: The name of the city

  Returns:
    The current temperature for the city
  """
  temperatures = {
    "New York": "22°C",
    "London": "15°C",
    "Tokyo": "18°C"
  }
  return temperatures.get(city, "Unknown")

def get_conditions(city: str) -> str:
  """Get the current weather conditions for a city
  
  Args:
    city: The name of the city

  Returns:
    The current weather conditions for the city
  """
  conditions = {
    "New York": "Partly cloudy",
    "London": "Rainy",
    "Tokyo": "Sunny"
  }
  return conditions.get(city, "Unknown")


messages = [{'role': 'user', 'content': 'What are the current weather conditions and temperature in New York and London?'}]

# The python client automatically parses functions as a tool schema so we can pass them directly
# Schemas can be passed directly in the tools list as well 
response = chat(model='qwen3', messages=messages, tools=[get_temperature, get_conditions], think=True)

# add the assistant message to the messages
messages.append(response.message)
if response.message.tool_calls:
  # process each tool call 
  for call in response.message.tool_calls:
    # execute the appropriate tool
    if call.function.name == 'get_temperature':
      result = get_temperature(**call.function.arguments)
    elif call.function.name == 'get_conditions':
      result = get_conditions(**call.function.arguments)
    else:
      result = 'Unknown tool'
    # add the tool result to the messages
    messages.append({'role': 'tool',  'tool_name': call.function.name, 'content': str(result)})

  # generate the final response
  final_response = chat(model='qwen3', messages=messages, tools=[get_temperature, get_conditions], think=True)
  print(final_response.message.content)
```

JavaScript:

```typescript
import ollama from 'ollama'

function getTemperature(city: string): string {
  const temperatures: { [key: string]: string } = {
    "New York": "22°C",
    "London": "15°C",
    "Tokyo": "18°C"
  }
  return temperatures[city] || "Unknown"
}

function getConditions(city: string): string {
  const conditions: { [key: string]: string } = {
    "New York": "Partly cloudy",
    "London": "Rainy",
    "Tokyo": "Sunny"
  }
  return conditions[city] || "Unknown"
}

const tools = [
  {
    type: 'function',
    function: {
      name: 'get_temperature',
      description: 'Get the current temperature for a city',
      parameters: {
        type: 'object',
        required: ['city'],
        properties: {
          city: { type: 'string', description: 'The name of the city' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'get_conditions',
      description: 'Get the current weather conditions for a city',
      parameters: {
        type: 'object',
        required: ['city'],
        properties: {
          city: { type: 'string', description: 'The name of the city' },
        },
      },
    },
  }
]

const messages = [{ role: 'user', content: 'What are the current weather conditions and temperature in New York and London?' }]

const response = await ollama.chat({
  model: 'qwen3',
  messages,
  tools,
  think: true
})

// add the assistant message to the messages
messages.push(response.message)
if (response.message.tool_calls) {
  // process each tool call 
  for (const call of response.message.tool_calls) {
    // execute the appropriate tool
    let result: string
    if (call.function.name === 'get_temperature') {
      const args = call.function.arguments as { city: string }
      result = getTemperature(args.city)
    } else if (call.function.name === 'get_conditions') {
      const args = call.function.arguments as { city: string }
      result = getConditions(args.city)
    } else {
      result = 'Unknown tool'
    }
    // add the tool result to the messages
    messages.push({ role: 'tool', tool_name: call.function.name, content: result })
  }

  // generate the final response
  const finalResponse = await ollama.chat({ model: 'qwen3', messages, tools, think: true })
  console.log(finalResponse.message.content)
}
```

## 멀티턴 도구 호출 (에이전트 루프)

에이전트 루프를 사용하면 모델이 언제 도구를 호출할지 스스로 결정하고, 그 결과를 응답에 반영할 수
있습니다. 모델에게 "지금 루프 안에 있으며 여러 번 도구를 호출할 수 있다"고 알려 주면 도움이 되기도
합니다.

Python:

```python
from ollama import chat, ChatResponse


def add(a: int, b: int) -> int:
  """Add two numbers"""
  """
  Args:
    a: The first number
    b: The second number

  Returns:
    The sum of the two numbers
  """
  return a + b


def multiply(a: int, b: int) -> int:
  """Multiply two numbers"""
  """
  Args:
    a: The first number
    b: The second number

  Returns:
    The product of the two numbers
  """
  return a * b


available_functions = {
  'add': add,
  'multiply': multiply,
}

messages = [{'role': 'user', 'content': 'What is (11434+12341)*412?'}]
while True:
    response: ChatResponse = chat(
        model='qwen3',
        messages=messages,
        tools=[add, multiply],
        think=True,
    )
    messages.append(response.message)
    print("Thinking: ", response.message.thinking)
    print("Content: ", response.message.content)
    if response.message.tool_calls:
        for tc in response.message.tool_calls:
            if tc.function.name in available_functions:
                print(f"Calling {tc.function.name} with arguments {tc.function.arguments}")
                result = available_functions[tc.function.name](**tc.function.arguments)
                print(f"Result: {result}")
                # add the tool result to the messages
                messages.append({'role': 'tool', 'tool_name': tc.function.name, 'content': str(result)})
    else:
        # end the loop when there are no more tool calls
        break
  # continue the loop with the updated messages
```

JavaScript:

```typescript
import ollama from 'ollama'

type ToolName = 'add' | 'multiply'

function add(a: number, b: number): number {
  return a + b
}

function multiply(a: number, b: number): number {
  return a * b
}

const availableFunctions: Record<ToolName, (a: number, b: number) => number> = {
  add,
  multiply,
}

const tools = [
  {
    type: 'function',
    function: {
      name: 'add',
      description: 'Add two numbers',
      parameters: {
        type: 'object',
        required: ['a', 'b'],
        properties: {
          a: { type: 'integer', description: 'The first number' },
          b: { type: 'integer', description: 'The second number' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'multiply',
      description: 'Multiply two numbers',
      parameters: {
        type: 'object',
        required: ['a', 'b'],
        properties: {
          a: { type: 'integer', description: 'The first number' },
          b: { type: 'integer', description: 'The second number' },
        },
      },
    },
  },
]

async function agentLoop() {
  const messages = [{ role: 'user', content: 'What is (11434+12341)*412?' }]

  while (true) {
    const response = await ollama.chat({
      model: 'qwen3',
      messages,
      tools,
      think: true,
    })

    messages.push(response.message)
    console.log('Thinking:', response.message.thinking)
    console.log('Content:', response.message.content)

    const toolCalls = response.message.tool_calls ?? []
    if (toolCalls.length) {
      for (const call of toolCalls) {
        const fn = availableFunctions[call.function.name as ToolName]
        if (!fn) {
          continue
        }

        const args = call.function.arguments as { a: number; b: number }
        console.log(`Calling ${call.function.name} with arguments`, args)
        const result = fn(args.a, args.b)
        console.log(`Result: ${result}`)
        messages.push({ role: 'tool', tool_name: call.function.name, content: String(result) })
      }
    } else {
      break
    }
  }
}

agentLoop().catch(console.error)
```

## 스트리밍과 함께 쓰는 도구 호출

스트리밍할 때는 `thinking`, `content`, `tool_calls`의 모든 청크를 모은 뒤, 이 필드들을 도구 실행 결과와
함께 후속 요청에 전달합니다.

Python:

```python
from ollama import chat 


def get_temperature(city: str) -> str:
  """Get the current temperature for a city
  
  Args:
    city: The name of the city

  Returns:
    The current temperature for the city
  """
  temperatures = {
    'New York': '22°C',
    'London': '15°C',
  }
  return temperatures.get(city, 'Unknown')


messages = [{'role': 'user', 'content': "What is the temperature in New York?"}]

while True:
  stream = chat(
    model='qwen3',
    messages=messages,
    tools=[get_temperature],
    stream=True,
    think=True,
  )

  thinking = ''
  content = ''
  tool_calls = []

  done_thinking = False
  # accumulate the partial fields
  for chunk in stream:
    if chunk.message.thinking:
      thinking += chunk.message.thinking
      print(chunk.message.thinking, end='', flush=True)
    if chunk.message.content:
      if not done_thinking:
        done_thinking = True
        print('\n')
      content += chunk.message.content
      print(chunk.message.content, end='', flush=True)
    if chunk.message.tool_calls:
      tool_calls.extend(chunk.message.tool_calls)
      print(chunk.message.tool_calls)

  # append accumulated fields to the messages
  if thinking or content or tool_calls:
    messages.append({'role': 'assistant', 'thinking': thinking, 'content': content, 'tool_calls': tool_calls})

  if not tool_calls:
    break

  for call in tool_calls:
    if call.function.name == 'get_temperature':
      result = get_temperature(**call.function.arguments)
    else:
      result = 'Unknown tool'
    messages.append({'role': 'tool', 'tool_name': call.function.name, 'content': result})
```

JavaScript:

```typescript
import ollama from 'ollama'

function getTemperature(city: string): string {
  const temperatures: Record<string, string> = {
    'New York': '22°C',
    'London': '15°C',
  }
  return temperatures[city] ?? 'Unknown'
}

const getTemperatureTool = {
  type: 'function',
  function: {
    name: 'get_temperature',
    description: 'Get the current temperature for a city',
    parameters: {
      type: 'object',
      required: ['city'],
      properties: {
        city: { type: 'string', description: 'The name of the city' },
      },
    },
  },
}

async function agentLoop() {
  const messages = [{ role: 'user', content: "What is the temperature in New York?" }]

  while (true) {
    const stream = await ollama.chat({
      model: 'qwen3',
      messages,
      tools: [getTemperatureTool],
      stream: true,
      think: true,
    })

    let thinking = ''
    let content = ''
    const toolCalls: any[] = []
    let doneThinking = false

    for await (const chunk of stream) {
      if (chunk.message.thinking) {
        thinking += chunk.message.thinking
        process.stdout.write(chunk.message.thinking)
      }
      if (chunk.message.content) {
        if (!doneThinking) {
          doneThinking = true
          process.stdout.write('\n')
        }
        content += chunk.message.content
        process.stdout.write(chunk.message.content)
      }
      if (chunk.message.tool_calls?.length) {
        toolCalls.push(...chunk.message.tool_calls)
        console.log(chunk.message.tool_calls)
      }
    }

    if (thinking || content || toolCalls.length) {
      messages.push({ role: 'assistant', thinking, content, tool_calls: toolCalls } as any)
    }

    if (!toolCalls.length) {
      break
    }

    for (const call of toolCalls) {
      if (call.function.name === 'get_temperature') {
        const args = call.function.arguments as { city: string }
        const result = getTemperature(args.city)
        messages.push({ role: 'tool', tool_name: call.function.name, content: result } )
      } else {
        messages.push({ role: 'tool', tool_name: call.function.name, content: 'Unknown tool' } )
      }
    }
  }
}

agentLoop().catch(console.error)
```

이 루프는 어시스턴트 응답을 스트리밍하면서 부분 필드를 누적하고, 이를 한꺼번에 다시 전달한 뒤 도구
실행 결과를 덧붙여 모델이 답변을 완성할 수 있게 합니다.

## Ollama Python SDK에서 함수를 도구로 사용하기

Python SDK는 함수를 자동으로 도구 스키마로 변환해 주므로, 함수를 그대로 전달할 수 있습니다. 필요하다면
스키마를 직접 전달할 수도 있습니다.

```python
from ollama import chat

def get_temperature(city: str) -> str:
  """Get the current temperature for a city
  
  Args:
    city: The name of the city

  Returns:
    The current temperature for the city
  """
  temperatures = {
    'New York': '22°C',
    'London': '15°C',
  }
  return temperatures.get(city, 'Unknown')

available_functions = {
  'get_temperature': get_temperature,
}
# directly pass the function as part of the tools list
response = chat(model='qwen3', messages=messages, tools=available_functions.values(), think=True)
```

> 원문: https://docs.ollama.com/capabilities/tool-calling

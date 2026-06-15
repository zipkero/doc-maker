# 클라우드

## 클라우드 모델

Ollama 클라우드 모델은 고성능 GPU 없이도 실행할 수 있는 새로운 형태의 모델입니다.
클라우드 모델은 자동으로 Ollama의 클라우드 서비스로 오프로딩되며, 로컬 모델과 동일한 기능을 제공합니다.
덕분에 개인 컴퓨터에는 올릴 수 없는 더 큰 모델을 돌리면서도 평소 쓰던 로컬 도구를 그대로 사용할 수 있습니다.

### 지원 모델

지원 모델 목록은 Ollama [모델 라이브러리](https://ollama.com/search?c=cloud)를 참고하세요.

### 클라우드 모델 실행하기

클라우드 모델을 쓰려면 [ollama.com](https://ollama.com) 계정이 필요합니다.
로그인하거나 새 계정을 만들려면 다음을 실행합니다.

```
ollama signin
```

CLI에서 클라우드 모델을 실행하려면 터미널에서 다음을 입력합니다.

```
ollama run gpt-oss:120b-cloud
```

Python에서 사용하려면 먼저 클라우드 모델을 pull 해 둡니다.

```
ollama pull gpt-oss:120b-cloud
```

[Ollama Python 라이브러리](https://github.com/ollama/ollama-python)를 설치합니다.

```
pip install ollama
```

간단한 스크립트로 실행합니다.

```python
from ollama import Client

client = Client()

messages = [
  {
    'role': 'user',
    'content': 'Why is the sky blue?',
  },
]

for part in client.chat('gpt-oss:120b-cloud', messages=messages, stream=True):
  print(part['message']['content'], end='', flush=True)
```

JavaScript에서 사용할 때도 먼저 클라우드 모델을 pull 합니다.

```
ollama pull gpt-oss:120b-cloud
```

[Ollama JavaScript 라이브러리](https://github.com/ollama/ollama-js)를 설치합니다.

```
npm i ollama
```

라이브러리로 클라우드 모델을 실행합니다.

```typescript
import { Ollama } from "ollama";

const ollama = new Ollama();

const response = await ollama.chat({
  model: "gpt-oss:120b-cloud",
  messages: [{ role: "user", content: "Explain quantum computing" }],
  stream: true,
});

for await (const part of response) {
  process.stdout.write(part.message.content);
}
```

cURL로 사용할 때도 모델을 먼저 pull 한 뒤, Ollama API를 통해 호출합니다.

```
ollama pull gpt-oss:120b-cloud
```

```
curl http://localhost:11434/api/chat -d '{
  "model": "gpt-oss:120b-cloud",
  "messages": [{
    "role": "user",
    "content": "Why is the sky blue?"
  }],
  "stream": false
}'
```

## 클라우드 API 직접 접근

클라우드 모델은 ollama.com의 API로 직접 접근할 수도 있습니다.
이 방식에서는 ollama.com이 원격 Ollama 호스트 역할을 합니다.

### 인증

ollama.com API에 직접 접근하려면 먼저 [API 키](https://ollama.com/settings/keys)를 발급합니다.
그다음 발급받은 키를 `OLLAMA_API_KEY` 환경 변수에 설정합니다.

```
export OLLAMA_API_KEY=your_api_key
```

### 모델 목록 조회

Ollama API로 직접 사용할 수 있는 모델은 다음과 같이 조회합니다.

```
curl https://ollama.com/api/tags
```

### 응답 생성

Python에서 사용하려면 [Ollama Python 라이브러리](https://github.com/ollama/ollama-python)를 설치합니다.

```
pip install ollama
```

호스트와 인증 헤더를 지정해 요청합니다.

```python
import os
from ollama import Client

client = Client(
    host="https://ollama.com",
    headers={'Authorization': 'Bearer ' + os.environ.get('OLLAMA_API_KEY')}
)

messages = [
  {
    'role': 'user',
    'content': 'Why is the sky blue?',
  },
]

for part in client.chat('gpt-oss:120b', messages=messages, stream=True):
  print(part['message']['content'], end='', flush=True)
```

JavaScript에서는 [Ollama JavaScript 라이브러리](https://github.com/ollama/ollama-js)를 설치합니다.

```
npm i ollama
```

호스트와 인증 헤더를 지정해 요청합니다.

```typescript
import { Ollama } from "ollama";

const ollama = new Ollama({
  host: "https://ollama.com",
  headers: {
    Authorization: "Bearer " + process.env.OLLAMA_API_KEY,
  },
});

const response = await ollama.chat({
  model: "gpt-oss:120b",
  messages: [{ role: "user", content: "Explain quantum computing" }],
  stream: true,
});

for await (const part of response) {
  process.stdout.write(part.message.content);
}
```

cURL로는 채팅 API에 다음과 같이 요청합니다.

```
curl https://ollama.com/api/chat \
  -H "Authorization: Bearer $OLLAMA_API_KEY" \
  -d '{
    "model": "gpt-oss:120b",
    "messages": [{
      "role": "user",
      "content": "Why is the sky blue?"
    }],
    "stream": false
  }'
```

## 로컬 전용 모드

[Ollama 클라우드 기능을 비활성화](./faq#how-do-i-disable-ollama-cloud)하면 Ollama를 로컬 전용 모드로 실행할 수 있습니다.

## 모델 지원 종료(deprecation)

Ollama는 더 좋은 오픈소스 모델이 나오면 오래된 클라우드 모델의 지원을 종료하고 폐기하기도 합니다.
Ollama 클라우드 모델에 의존하는 도구나 애플리케이션은 계속 동작하도록 업데이트가 필요할 수 있습니다.
영향을 받는 사용자에게는 모델 지원 종료·폐기에 앞서 미리 안내가 전달되며, 안내는 이메일과 Ollama 웹사이트를 통해 이루어집니다.

> 참고: 클라우드 모델 폐기는 로컬 모델에는 영향을 주지 않습니다.

### 예정된 지원 종료

| 폐기일 | 모델 | 권장 대체 모델 |
| --- | --- | --- |
| 2026년 6월 16일 | `kimi-k2-thinking` | `kimi-k2.6` |
| 2026년 6월 16일 | `kimi-k2:1t` | `kimi-k2.6` |
| 2026년 6월 16일 | `minimax-m2` | `minimax-m3` |
| 2026년 6월 16일 | `glm-4.6` | `glm-5.1` |
| 2026년 6월 16일 | `qwen3-next:80b` | `qwen3.5` |
| 2026년 6월 16일 | `qwen3-vl:235b` | `qwen3.5` |
| 2026년 6월 16일 | `qwen3-vl:235b-instruct` | `qwen3.5` |
| 2026년 6월 16일 | `cogito-2.1:671b` | `deepseek-v4-flash` |

> 원문: https://docs.ollama.com/cloud

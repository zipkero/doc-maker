# 구조화된 출력 (Structured Outputs)

> 참고: Ollama Cloud는 현재 구조화된 출력을 지원하지 않습니다.

구조화된 출력은 모델 응답에 JSON 스키마를 강제하는 기능입니다. 이를 통해 구조화된 데이터를 안정적으로
추출하거나, 이미지를 일정한 형식으로 설명하거나, 모든 응답을 일관된 형태로 유지할 수 있습니다.

## 구조화된 JSON 생성

`format` 필드에 `json`을 지정하면 응답이 JSON 형식으로 반환됩니다.

cURL:

```shell
curl -X POST http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "gpt-oss",
  "messages": [{"role": "user", "content": "Tell me about Canada in one line"}],
  "stream": false,
  "format": "json"
}'
```

Python:

```python
from ollama import chat

response = chat(
  model='gpt-oss',
  messages=[{'role': 'user', 'content': 'Tell me about Canada.'}],
  format='json'
)
print(response.message.content)
```

JavaScript:

```javascript
import ollama from 'ollama'

const response = await ollama.chat({
  model: 'gpt-oss',
  messages: [{ role: 'user', content: 'Tell me about Canada.' }],
  format: 'json'
})
console.log(response.message.content)
```

## 스키마를 지정한 구조화된 JSON 생성

`format` 필드에 JSON 스키마를 직접 전달하면, 응답이 그 스키마에 맞춰 생성됩니다.

> 참고: 동일한 JSON 스키마를 프롬프트에도 문자열로 함께 전달하면 모델 응답을 더 잘 유도할 수 있습니다.

cURL:

```shell
curl -X POST http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "gpt-oss",
  "messages": [{"role": "user", "content": "Tell me about Canada."}],
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "name": {"type": "string"},
      "capital": {"type": "string"},
      "languages": {
        "type": "array",
        "items": {"type": "string"}
      }
    },
    "required": ["name", "capital", "languages"]
  }
}'
```

Python: Pydantic 모델을 정의하고 `model_json_schema()`를 `format`에 전달한 뒤, 응답을 검증합니다.

```python
from ollama import chat
from pydantic import BaseModel

class Country(BaseModel):
  name: str
  capital: str
  languages: list[str]

response = chat(
  model='gpt-oss',
  messages=[{'role': 'user', 'content': 'Tell me about Canada.'}],
  format=Country.model_json_schema(),
)

country = Country.model_validate_json(response.message.content)
print(country)
```

JavaScript: Zod 스키마를 `z.toJSONSchema()`로 직렬화하고, 구조화된 응답을 파싱합니다.

```javascript
import ollama from 'ollama'
import * as z from 'zod'

const Country = z.object({
  name: z.string(),
  capital: z.string(),
  languages: z.array(z.string()),
})

const response = await ollama.chat({
  model: 'gpt-oss',
  messages: [{ role: 'user', content: 'Tell me about Canada.' }],
  format: z.toJSONSchema(Country),
})

const country = Country.parse(JSON.parse(response.message.content))
console.log(country)
```

## 예시: 구조화된 데이터 추출

반환받고 싶은 객체를 정의하면 모델이 각 필드를 채워 줍니다.

```python
from ollama import chat
from pydantic import BaseModel

class Pet(BaseModel):
  name: str
  animal: str
  age: int
  color: str | None
  favorite_toy: str | None

class PetList(BaseModel):
  pets: list[Pet]

response = chat(
  model='gpt-oss',
  messages=[{'role': 'user', 'content': 'I have two cats named Luna and Loki...'}],
  format=PetList.model_json_schema(),
)

pets = PetList.model_validate_json(response.message.content)
print(pets)
```

## 예시: 비전 모델과 구조화된 출력

비전 모델도 동일한 `format` 파라미터를 받으므로, 이미지를 결정론적인 형태로 설명하게 만들 수 있습니다.

```python
from ollama import chat
from pydantic import BaseModel
from typing import Literal, Optional

class Object(BaseModel):
  name: str
  confidence: float
  attributes: str

class ImageDescription(BaseModel):
  summary: str
  objects: list[Object]
  scene: str
  colors: list[str]
  time_of_day: Literal['Morning', 'Afternoon', 'Evening', 'Night']
  setting: Literal['Indoor', 'Outdoor', 'Unknown']
  text_content: Optional[str] = None

response = chat(
  model='gemma4',
  messages=[{
    'role': 'user',
    'content': 'Describe this photo and list the objects you detect.',
    'images': ['path/to/image.jpg'],
  }],
  format=ImageDescription.model_json_schema(),
  options={'temperature': 0},
)

image_description = ImageDescription.model_validate_json(response.message.content)
print(image_description)
```

## 안정적인 구조화된 출력을 위한 팁

- 스키마는 Pydantic(Python)이나 Zod(JavaScript)로 정의해 검증에 재사용할 수 있도록 합니다.
- 더 결정론적인 완성을 원한다면 온도를 낮춥니다(예: `0`으로 설정).
- 구조화된 출력은 OpenAI 호환 API에서 `response_format`을 통해서도 동작합니다.

> 원문: https://docs.ollama.com/capabilities/structured-outputs

# 비전 (Vision)

비전 모델은 텍스트와 함께 이미지를 입력받아, 모델이 본 내용을 설명하거나 분류하고 그에 대한 질문에
답할 수 있습니다.

## 빠른 시작

```shell
ollama run gemma4 ./image.png whats in this image?
```

## Ollama API에서 사용하기

`images` 배열을 전달합니다. SDK는 파일 경로, URL, 원시 바이트를 받지만, REST API는 base64로 인코딩된
이미지 데이터를 받습니다.

cURL:

```shell
# 1. Download a sample image
curl -L -o test.jpg "https://upload.wikimedia.org/wikipedia/commons/3/3a/Cat03.jpg"

# 2. Encode the image
IMG=$(base64 < test.jpg | tr -d '\n')

# 3. Send it to Ollama
curl -X POST http://localhost:11434/api/chat \
-H "Content-Type: application/json" \
-d '{
    "model": "gemma4",
    "messages": [{
    "role": "user",
    "content": "What is in this image?",
    "images": ["'"$IMG"'"]
    }],
    "stream": false
}'
```

Python:

```python
from ollama import chat
# from pathlib import Path

# Pass in the path to the image
path = input('Please enter the path to the image: ')

# You can also pass in base64 encoded image data
# img = base64.b64encode(Path(path).read_bytes()).decode()
# or the raw bytes
# img = Path(path).read_bytes()

response = chat(
  model='gemma4',
  messages=[
    {
      'role': 'user',
      'content': 'What is in this image? Be concise.',
      'images': [path],
    }
  ],
)

print(response.message.content)
```

JavaScript:

```javascript
import ollama from 'ollama'

const imagePath = '/absolute/path/to/image.jpg'
const response = await ollama.chat({
  model: 'gemma4',
  messages: [
    { role: 'user', content: 'What is in this image?', images: [imagePath] }
  ],
  stream: false,
})

console.log(response.message.content)
```

> 원문: https://docs.ollama.com/capabilities/vision

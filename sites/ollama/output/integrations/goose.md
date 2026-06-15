# Goose

## Goose Desktop

[Goose](https://block.github.io/goose/docs/getting-started/installation/) Desktop을 설치합니다.

### Ollama와 함께 사용하기

1. Goose에서 **Settings** → **Configure Provider**를 엽니다.
2. **Ollama**를 찾아 **Configure**를 클릭합니다.
3. **API Host**가 `http://localhost:11434`인지 확인하고 Submit을 클릭합니다.

### ollama.com에 연결하기

1. ollama.com에서 [API 키](https://ollama.com/settings/keys)를 생성하고 `.env`에 저장합니다.
2. Goose에서 **API Host**를 `https://ollama.com`으로 설정합니다.

## Goose CLI

[Goose](https://block.github.io/goose/docs/getting-started/installation/) CLI를 설치합니다.

### Ollama와 함께 사용하기

1. `goose configure`를 실행합니다.
2. **Configure Providers**를 선택한 뒤 **Ollama**를 선택합니다.
3. 모델 이름을 입력합니다(예: `qwen3`).

### ollama.com에 연결하기

1. ollama.com에서 [API 키](https://ollama.com/settings/keys)를 생성하고 `.env`에 저장합니다.
2. `goose configure`를 실행합니다.
3. **Configure Providers**를 선택한 뒤 **Ollama**를 선택합니다.
4. **OLLAMA_HOST**를 `https://ollama.com`으로 변경합니다.

> 원문: https://docs.ollama.com/integrations/goose

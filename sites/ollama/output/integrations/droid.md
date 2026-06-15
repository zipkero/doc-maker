# Droid

## 설치

[Droid CLI](https://factory.ai/)를 설치합니다.

```bash
curl -fsSL https://app.factory.ai/cli | sh
```

> 참고: Droid는 더 큰 컨텍스트 윈도우가 필요합니다. 최소 64k 토큰 이상의 컨텍스트 윈도우를 권장합니다.
> 자세한 내용은 [컨텍스트 길이](/context-length) 문서를 참고하세요.

## Ollama와 함께 사용하기

### 빠른 설정

```bash
ollama launch droid
```

실행하지 않고 설정만 하려면 다음을 사용합니다.

```shell
ollama launch droid --config
```

### 수동 설정

`~/.factory/config.json`에 로컬 설정 블록을 추가합니다.

```json
{
  "custom_models": [
    {
      "model_display_name": "qwen3-coder [Ollama]",
      "model": "qwen3-coder",
      "base_url": "http://localhost:11434/v1/",
      "api_key": "not-needed",
      "provider": "generic-chat-completion-api",
      "max_tokens": 32000 
    }
  ]
}
```

## 클라우드 모델

Droid에는 `qwen3-coder:480b-cloud` 모델이 권장됩니다.

`~/.factory/config.json`에 클라우드 설정 블록을 추가합니다.

```json
{
  "custom_models": [
    {
      "model_display_name": "qwen3-coder [Ollama Cloud]",
      "model": "qwen3-coder:480b-cloud",
      "base_url": "http://localhost:11434/v1/",
      "api_key": "not-needed",
      "provider": "generic-chat-completion-api",
      "max_tokens": 128000
    }
  ]
}
```

## ollama.com에 연결하기

1. ollama.com에서 [API 키](https://ollama.com/settings/keys)를 생성하고 `OLLAMA_API_KEY`로
   내보냅니다.
2. `~/.factory/config.json`에 클라우드 설정 블록을 추가합니다.

   ```json
   {
     "custom_models": [
       {
         "model_display_name": "qwen3-coder [Ollama Cloud]",
         "model": "qwen3-coder:480b",
         "base_url": "https://ollama.com/v1/",
         "api_key": "OLLAMA_API_KEY",
         "provider": "generic-chat-completion-api",
         "max_tokens": 128000
       }
     ]
   }
   ```

새 터미널에서 `droid`를 실행하면 변경한 설정이 적용됩니다.

> 원문: https://docs.ollama.com/integrations/droid

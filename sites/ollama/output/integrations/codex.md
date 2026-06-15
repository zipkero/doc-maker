# Codex CLI

## 설치

[Codex CLI](https://developers.openai.com/codex/cli/)를 설치합니다. 데스크톱 앱은
[Codex App](/integrations/codex-app) 문서를 참고하세요.

```
npm install -g @openai/codex
```

## Ollama와 함께 사용하기

> 참고: Codex는 더 큰 컨텍스트 윈도우가 필요합니다. 최소 64k 토큰 이상의 컨텍스트 윈도우를 권장합니다.

### 빠른 설정

```
ollama launch codex
```

`ollama launch codex`로 실행하면 Ollama가 모델 카탈로그를 갱신하고, 해당 세션에 전용 Codex 프로필을
사용합니다.

실행하지 않고 설정만 하려면 다음을 사용합니다.

```shell
ollama launch codex --config
```

Ollama launch 프로필과 생성된 모델 카탈로그를 제거하려면 다음을 사용합니다.

```shell
ollama launch codex --restore
```

### 수동 설정

`codex`를 Ollama와 함께 쓰려면 `--oss` 플래그를 사용합니다.

```
codex --oss
```

특정 모델을 지정하려면 `-m` 플래그를 전달합니다.

```
codex --oss -m gpt-oss:120b
```

클라우드 모델을 쓰려면 다음과 같이 합니다.

```
codex --oss -m gpt-oss:120b-cloud
```

### 프로필 기반 설정

Codex CLI 설정을 지속적으로 유지하려면 `~/.codex/ollama-launch.config.toml` 파일을 만듭니다.

```toml
model = "gpt-oss:120b"
model_provider = "ollama-launch"
model_catalog_json = "/Users/you/.codex/model.json"

[model_providers.ollama-launch]
name = "Ollama"
base_url = "http://localhost:11434/v1/"
wire_api = "responses"
```

그런 다음 다음을 실행합니다.

```
codex --profile ollama-launch
```

> 원문: https://docs.ollama.com/integrations/codex

# Copilot CLI

GitHub Copilot CLI는 터미널에서 동작하는 GitHub의 AI 코딩 에이전트입니다. 코드베이스를 이해하고,
코드를 편집하고, 명령을 실행해 더 빠른 소프트웨어 개발을 돕습니다.

Ollama를 통해 오픈 모델을 Copilot CLI와 함께 사용할 수 있습니다. 예를 들어 `qwen3.5`,
`glm-5.1:cloud`, `kimi-k2.5:cloud` 같은 모델을 쓸 수 있습니다.

## 설치

[Copilot CLI](https://github.com/features/copilot/cli/)를 설치합니다.

macOS / Linux (Homebrew):

```shell
brew install copilot-cli
```

npm (모든 플랫폼):

```shell
npm install -g @github/copilot
```

macOS / Linux (스크립트):

```shell
curl -fsSL https://gh.io/copilot-install | bash
```

Windows (WinGet):

```powershell
winget install GitHub.Copilot
```

## Ollama와 함께 사용하기

### 빠른 설정

```shell
ollama launch copilot
```

### 특정 모델로 바로 실행

```shell
ollama launch copilot --model kimi-k2.5:cloud
```

## 추천 모델

* `kimi-k2.5:cloud`
* `glm-5:cloud`
* `minimax-m2.7:cloud`
* `qwen3.5:cloud`
* `glm-4.7-flash`
* `qwen3.5`

클라우드 모델은 [ollama.com/search?c=cloud](https://ollama.com/search?c=cloud)에서도 확인할 수 있습니다.

## 비대화형(헤드리스) 모드

Docker, CI/CD, 스크립트 등에서 상호작용 없이 Copilot CLI를 실행할 수 있습니다.

```shell
ollama launch copilot --model kimi-k2.5:cloud --yes -- -p "how does this repository work?"
```

`--yes` 플래그는 모델을 자동으로 pull하고 선택 절차를 건너뜁니다. 이 플래그를 쓰려면 `--model`을 함께
지정해야 합니다. `--` 뒤의 인자는 Copilot CLI에 그대로 전달됩니다.

## 수동 설정

Copilot CLI는 환경변수를 통해 OpenAI 호환 API로 Ollama에 연결합니다.

1. 환경변수를 설정합니다.

```shell
export COPILOT_PROVIDER_BASE_URL=http://localhost:11434/v1
export COPILOT_PROVIDER_API_KEY=
export COPILOT_PROVIDER_WIRE_API=responses
export COPILOT_MODEL=qwen3.5
```

2. Copilot CLI를 실행합니다.

```shell
copilot
```

또는 환경변수를 인라인으로 지정해 실행할 수도 있습니다.

```shell
COPILOT_PROVIDER_BASE_URL=http://localhost:11434/v1 COPILOT_PROVIDER_API_KEY= COPILOT_PROVIDER_WIRE_API=responses COPILOT_MODEL=glm-5:cloud copilot
```

> 참고: Copilot은 큰 컨텍스트 윈도우가 필요합니다. 최소 64k 토큰을 권장합니다. Ollama에서 컨텍스트
> 길이를 조정하는 방법은 [컨텍스트 길이 문서](/context-length)를 참고하세요.

> 원문: https://docs.ollama.com/integrations/copilot-cli

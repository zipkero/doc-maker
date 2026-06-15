# OpenClaw

OpenClaw는 사용자의 기기에서 직접 동작하는 개인용 AI 어시스턴트입니다. 메시징 서비스(WhatsApp,
Telegram, Slack, Discord, iMessage 등)를 중앙 게이트웨이를 통해 AI 코딩 에이전트와 연결합니다.

## 빠른 시작

```bash
ollama launch openclaw
```

이 명령 하나로 Ollama가 다음을 자동으로 처리합니다.

1. **설치** — OpenClaw가 설치되어 있지 않으면 npm으로 설치하도록 안내합니다.
2. **보안** — 첫 실행 시 도구 접근에 따르는 위험을 설명하는 보안 안내가 표시됩니다.
3. **모델** — 선택기에서 로컬 또는 클라우드 모델을 고릅니다.
4. **온보딩** — 제공자를 설정하고, 게이트웨이 데몬을 설치하고, 선택한 모델을 기본 모델로 지정한 뒤, OpenClaw에 내장된 Ollama 웹 검색을 활성화합니다.
5. **게이트웨이** — 백그라운드에서 시작되며 OpenClaw TUI를 엽니다.

> 참고: OpenClaw는 비교적 큰 컨텍스트 창을 요구합니다. 로컬 모델을 사용할 경우 최소 64k 토큰의
> 컨텍스트 창을 권장합니다. 자세한 내용은 [컨텍스트 길이](/context-length) 문서를 참고하세요.

> 참고: 이전 이름은 Clawdbot이었습니다. `ollama launch clawdbot`도 별칭으로 계속 동작합니다.

## 웹 검색 및 가져오기

OpenClaw에는 Ollama `web_search` 제공자가 내장되어 있어, 로컬이나 클라우드 기반 Ollama 환경에서
설정된 Ollama 호스트를 통해 웹을 검색할 수 있습니다.

```bash
ollama launch openclaw
```

Ollama를 통해 OpenClaw를 실행하면 Ollama 웹 검색이 자동으로 활성화됩니다. 수동으로 설정하려면 다음을
실행합니다.

```bash
openclaw configure --section web
```

> 참고: 로컬 모델에서 Ollama 웹 검색을 사용하려면 `ollama signin`이 필요합니다.

## 실행 없이 설정만 변경하기

게이트웨이와 TUI를 시작하지 않고 모델만 변경하려면 다음을 실행합니다.

```bash
ollama launch openclaw --config
```

특정 모델을 바로 사용하려면 다음을 실행합니다.

```bash
ollama launch openclaw --model kimi-k2.5:cloud
```

게이트웨이가 이미 실행 중이면 새 모델을 반영하기 위해 자동으로 재시작됩니다.

## 권장 모델

**클라우드 모델**

* `kimi-k2.5:cloud` — 서브에이전트를 갖춘 멀티모달 추론
* `qwen3.5:cloud` — 추론, 코딩, 비전 기반 에이전트형 도구 활용
* `glm-5.1:cloud` — 추론 및 코드 생성
* `minimax-m2.7:cloud` — 빠르고 효율적인 코딩 및 실무 생산성

**로컬 모델**

* `gemma4` — 로컬 환경에서 추론 및 코드 생성(VRAM 약 16 GB)
* `qwen3.5` — 로컬 환경에서 추론, 코딩, 시각 이해(VRAM 약 11 GB)

더 많은 모델은 [ollama.com/search](https://ollama.com/search?c=cloud)에서 확인할 수 있습니다.

## 비대화형(헤드리스) 모드

Docker, CI/CD, 스크립트에서 사용할 수 있도록 상호작용 없이 OpenClaw를 실행합니다.

```bash
ollama launch openclaw --model kimi-k2.5:cloud --yes
```

`--yes` 플래그는 모델을 자동으로 받고, 선택기를 건너뛰며, `--model` 지정을 필수로 요구합니다.

## 메시징 앱 연결

```bash
openclaw configure --section channels
```

WhatsApp, Telegram, Slack, Discord, iMessage를 연결하면 어디서나 로컬 모델과 대화할 수 있습니다.

## 게이트웨이 중지

```bash
openclaw gateway stop
```

> 원문: https://docs.ollama.com/integrations/openclaw

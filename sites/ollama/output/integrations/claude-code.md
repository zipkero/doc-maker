# Claude Code

Claude Code는 Anthropic이 만든 에이전트형 코딩 도구로, 작업 디렉터리의 코드를 읽고 수정하고 실행할 수 있습니다.

Ollama의 Anthropic 호환 API를 통해 오픈 모델을 Claude Code와 함께 사용할 수 있습니다. 예를 들어
`qwen3.5`, `glm-5:cloud`, `kimi-k2.5:cloud` 같은 모델을 쓸 수 있습니다.

## 설치

[Claude Code](https://code.claude.com/docs/en/overview)를 설치합니다.

macOS / Linux:

```shell
curl -fsSL https://claude.ai/install.sh | bash
```

Windows:

```powershell
irm https://claude.ai/install.ps1 | iex
```

## Ollama와 함께 사용하기

### 빠른 설정

```shell
ollama launch claude
```

### 특정 모델로 바로 실행

```shell
ollama launch claude --model kimi-k2.5:cloud
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

Docker, CI/CD, 스크립트 등에서 상호작용 없이 Claude Code를 실행할 수 있습니다.

```shell
ollama launch claude --model kimi-k2.5:cloud --yes -- -p "how does this repository work?"
```

`--yes` 플래그는 모델을 자동으로 pull하고 선택 절차를 건너뜁니다. 이 플래그를 쓰려면 `--model`을 함께
지정해야 합니다. `--` 뒤의 인자는 Claude Code에 그대로 전달됩니다.

## 웹 검색

Claude Code는 Ollama의 웹 검색 API를 통해 웹을 검색할 수 있습니다. 설정과 사용법은
[웹 검색 문서](/capabilities/web-search)를 참고하세요.

## `/loop`로 예약 작업 실행하기

`/loop` 명령은 Claude Code 안에서 프롬프트나 슬래시 명령을 정해진 주기로 반복 실행합니다. PR 점검,
리서치, 리마인더 설정처럼 반복적인 작업을 자동화할 때 유용합니다.

```
/loop <interval> <prompt or /command>
```

### 예시

내 PR 챙기기:

```
/loop 30m Check my open PRs and summarize their status
```

리서치 작업 자동화:

```
/loop 1h Research the latest AI news and summarize key developments
```

버그 리포트 및 분류 자동화:

```
/loop 15m Check for new GitHub issues and triage by priority
```

리마인더 설정:

```
/loop 1h Remind me to review the deploy status
```

## Telegram

봇을 세션에 연결하면 Telegram에서 Claude Code와 대화할 수 있습니다.
[Telegram 플러그인](https://github.com/anthropics/claude-plugins-official)을 설치하고,
[@BotFather](https://t.me/BotFather)로 봇을 만든 다음, 채널 플래그를 붙여 실행합니다.

```shell
ollama launch claude -- --channels plugin:telegram@claude-plugins-official
```

Claude Code는 대부분의 작업에서 권한을 묻습니다. 봇이 자율적으로 동작하게 하려면
[권한 규칙](https://code.claude.com/docs/en/permissions)을 설정하거나, 격리된 환경에서
`--dangerously-skip-permissions`를 전달하세요.

페어링과 접근 제어를 포함한 전체 설정 방법은
[플러그인 README](https://github.com/anthropics/claude-plugins-official/tree/main/external_plugins/telegram)를 참고하세요.

## 수동 설정

Claude Code는 Anthropic 호환 API로 Ollama에 연결합니다.

1. 환경변수를 설정합니다.

```shell
export ANTHROPIC_AUTH_TOKEN=ollama
export ANTHROPIC_API_KEY=""
export ANTHROPIC_BASE_URL=http://localhost:11434
```

2. Ollama 모델로 Claude Code를 실행합니다.

```shell
claude --model qwen3.5
```

또는 환경변수를 인라인으로 지정해 실행할 수도 있습니다.

```shell
ANTHROPIC_AUTH_TOKEN=ollama ANTHROPIC_BASE_URL=http://localhost:11434 ANTHROPIC_API_KEY="" claude --model glm-5:cloud
```

> 참고: Claude Code는 큰 컨텍스트 윈도우가 필요합니다. 최소 64k 토큰을 권장합니다. Ollama에서
> 컨텍스트 길이를 조정하는 방법은 [컨텍스트 길이 문서](/context-length)를 참고하세요.

> 원문: https://docs.ollama.com/integrations/claude-code

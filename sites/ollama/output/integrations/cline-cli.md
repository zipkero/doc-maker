# Cline CLI

Cline CLI는 대화형 터미널 세션에서 동작하는 자율 코딩 에이전트입니다.

## 설치

[Cline CLI](https://docs.cline.bot/usage/cli-overview)를 설치합니다. IDE 확장 기능은
[Cline](/integrations/cline) 문서를 참고하세요.

```bash
npm install -g cline
```

> 참고: Cline CLI가 설치되어 있지 않고 `npm`을 사용할 수 있는 경우, `ollama launch cline`을 실행하면
> `cline@latest` 설치를 묻습니다.

## Ollama와 함께 사용하기

### 빠른 설정

```bash
ollama launch cline
```

`ollama launch cline`으로 실행하면 Ollama가 Cline의 공급자를 Ollama로 설정하고, 로컬 Ollama
엔드포인트를 가리키도록 한 뒤, 선택한 모델을 적용합니다.

실행하지 않고 설정만 하려면 다음을 사용합니다.

```shell
ollama launch cline --config
```

### 특정 모델로 바로 실행

```shell
ollama launch cline --model qwen3.5
```

클라우드 모델을 쓰려면 다음과 같이 합니다.

```shell
ollama launch cline --model kimi-k2.6:cloud
```

### Cline에 프롬프트 전달하기

`--` 뒤의 인자는 Cline에 그대로 전달됩니다.

```shell
ollama launch cline -- "summarize this repository"
```

Cline의 칸반 보드를 열려면 다음과 같이 합니다.

```shell
ollama launch cline -- kanban
```

### 수동 설정

Cline CLI를 수동으로 설정하려면, 먼저 Ollama가 실행 중인지, 그리고 사용할 모델이 준비되어 있는지
확인합니다.

```shell
ollama pull qwen3.5
```

그런 다음 Cline의 대화형 인증 절차를 실행합니다.

```shell
cline auth
```

공급자로 Ollama를 선택하고, 기본 URL을 묻는다면 `http://localhost:11434`를 입력한 뒤,
`qwen3.5`나 `kimi-k2.6:cloud` 같은 모델을 고릅니다.

현재 Cline 설정을 확인하려면 다음을 실행합니다.

```shell
cline config
```

대화형 세션을 시작하려면 다음을 실행합니다.

```shell
cline
```

> 원문: https://docs.ollama.com/integrations/cline-cli

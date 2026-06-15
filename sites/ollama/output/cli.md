# CLI 레퍼런스

Ollama를 명령줄에서 다루는 주요 명령을 정리했습니다.

## 모델 실행

모델을 실행하려면 다음과 같이 입력합니다.

```
ollama run gemma4
```

### 여러 줄 입력

여러 줄을 한 번에 입력하려면 텍스트를 `"""`로 감쌉니다.

```
>>> """Hello,
... world!
... """
I'm a basic program that prints the famous "Hello, world!" message to the console.
```

### 멀티모달 모델

이미지를 함께 전달하려면 프롬프트에 파일 경로를 포함합니다.

```
ollama run gemma4 "What's in this image? /Users/jmorgan/Desktop/smile.png"
```

## 통합(integration) 실행

`ollama launch`는 Ollama 모델을 사용하는 외부 애플리케이션을 설정하고 실행합니다.
지원 앱과의 통합을 대화형으로 구성하고 시작할 수 있습니다.

```
ollama launch
```

### 지원하는 통합

- **OpenCode** — 오픈소스 코딩 어시스턴트
- **Claude Code** — Anthropic의 에이전트형 코딩 도구
- **Codex** — OpenAI의 코딩 어시스턴트
- **VS Code** — AI 채팅이 내장된 Microsoft IDE
- **Droid** — Factory의 AI 코딩 에이전트

### 사용 예시

통합을 대화형으로 실행:

```
ollama launch
```

특정 통합을 지정해 실행:

```
ollama launch claude
```

사용할 모델을 지정해 실행:

```
ollama launch claude --model qwen3.5
```

실행하지 않고 설정만 진행:

```
ollama launch droid --config
```

## 임베딩 생성

임베딩 모델을 실행해 임베딩을 생성합니다.

```
ollama run embeddinggemma "Hello world"
```

표준 입력으로 텍스트를 전달할 수도 있으며, 출력은 JSON 배열로 반환됩니다.

```
echo "Hello world" | ollama run nomic-embed-text
```

## 모델 내려받기

```
ollama pull gemma4
```

## 모델 삭제

```
ollama rm gemma4
```

## 모델 목록 보기

```
ollama ls
```

## Ollama 로그인

```
ollama signin
```

## Ollama 로그아웃

```
ollama signout
```

## 커스텀 모델 만들기

먼저 `Modelfile`을 작성합니다.

```
FROM gemma4
SYSTEM """You are a happy cat."""
```

그다음 `ollama create`를 실행합니다.

```
ollama create -f Modelfile
```

## 실행 중인 모델 목록 보기

```
ollama ps
```

## 실행 중인 모델 중지

```
ollama stop gemma4
```

## Ollama 시작

```
ollama serve
```

설정할 수 있는 환경 변수 목록을 보려면 `ollama serve --help`를 실행합니다.

> 원문: https://docs.ollama.com/cli

# 빠른 시작

Ollama는 macOS, Windows, Linux에서 사용할 수 있습니다. [Ollama 다운로드 페이지](https://ollama.com/download)에서
설치 파일을 받을 수 있습니다.

## 시작하기

터미널에서 `ollama`를 실행하면 대화형 메뉴가 열립니다.

```sh
ollama
```

`↑/↓`로 이동하고, `enter`로 실행, `→`로 모델 변경, `esc`로 종료합니다.

메뉴에서는 다음에 빠르게 접근할 수 있습니다.

- **Run a model** — 대화형 채팅 시작
- **Launch tools** — Claude Code, Codex, OpenClaw 등 실행
- **Additional integrations** — "More..." 아래에서 확인

## 어시스턴트

100개 이상의 스킬을 갖춘 개인용 AI [OpenClaw](/integrations/openclaw)를 실행합니다.

```sh
ollama launch openclaw
```

## 코딩

Ollama 모델로 [Claude Code](/integrations/claude-code)를 비롯한 여러 코딩 도구를 실행합니다.

```sh
ollama launch claude
```

```sh
ollama launch codex
```

```sh
ollama launch opencode
```

지원되는 모든 도구는 [통합(integrations)](/integrations) 문서를 참고하세요.

## API

[API](/api)를 사용해 Ollama를 애플리케이션에 통합할 수 있습니다.

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "gemma4",
  "messages": [{ "role": "user", "content": "Hello!" }]
}'
```

Python, JavaScript 등 다른 언어로의 통합은 [API 문서](/api)를 참고하세요.

> 원문: https://docs.ollama.com/quickstart

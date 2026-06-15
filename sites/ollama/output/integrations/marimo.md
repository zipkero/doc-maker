# marimo

## 설치

[marimo](https://marimo.io)를 설치합니다. `pip`나 `uv`를 사용할 수 있습니다. 다음과 같이 `uv`로
marimo용 샌드박스 환경을 만들 수도 있습니다.

```
uvx marimo edit --sandbox notebook.py
```

## Ollama와 함께 사용하기

1. marimo에서 사용자 설정(user settings)으로 들어가 AI 탭을 엽니다. 여기서 Ollama를 AI 제공자로
   찾아 설정할 수 있습니다. 로컬에서 사용할 때는 보통 base URL을 `http://localhost:11434/v1`로
   지정합니다.
2. AI 제공자를 설정하고 나면, 사용할 특정 AI 모델을 켜거나 끌 수 있습니다.
3. 설정 화면 맨 아래로 스크롤해 UI에서 사용 가능한 모델 목록에 새 모델을 추가할 수도 있습니다.
4. 설정이 끝나면 marimo에서 Ollama로 AI 채팅을 사용할 수 있습니다.
5. 또는 marimo에서 Ollama로 **인라인 코드 자동완성**을 사용할 수도 있습니다. 이 기능은 "AI Features"
   탭에서 설정합니다.

## ollama.com에 연결하기

1. `ollama signin`으로 Ollama 클라우드에 로그인합니다.
2. Ollama 모델 설정에서 Ollama가 호스팅하는 모델(예: `gpt-oss:120b`)을 추가합니다.
3. 이제 marimo에서 해당 모델을 참조할 수 있습니다.

> 원문: https://docs.ollama.com/integrations/marimo

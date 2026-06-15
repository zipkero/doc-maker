# Cline

## 설치

IDE에 [Cline](https://docs.cline.bot/getting-started/installing-cline)을 설치합니다.

## Ollama와 함께 사용하기

1. Cline 설정에서 `API Configuration`을 열고 `API Provider`를 `Ollama`로 설정합니다.
2. `Model`에서 모델을 선택하거나 직접 입력합니다(예: `qwen3`).
3. `Context Window`에서 컨텍스트 윈도우를 최소 32K 토큰 이상으로 설정합니다.

> 참고: 코딩 도구는 더 큰 컨텍스트 윈도우가 필요합니다. 최소 32K 토큰 이상의 컨텍스트 윈도우를
> 권장합니다. 자세한 내용은 [컨텍스트 길이](/context-length) 문서를 참고하세요.

## ollama.com에 연결하기

1. ollama.com에서 [API 키](https://ollama.com/settings/keys)를 생성합니다.
2. `Use custom base URL`을 클릭하고 `https://ollama.com`으로 설정합니다.
3. 발급받은 **Ollama API Key**를 입력합니다.
4. 목록에서 모델을 선택합니다.

### 추천 모델

* `qwen3-coder:480b`
* `deepseek-v3.1:671b`

> 원문: https://docs.ollama.com/integrations/cline

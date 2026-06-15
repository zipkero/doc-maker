# Roo Code

## 설치

VS Code 마켓플레이스에서 [Roo Code](https://marketplace.visualstudio.com/items?itemName=RooVeterinaryInc.roo-cline)를
설치합니다.

## Ollama와 함께 사용하기

1. VS Code에서 Roo Code를 열고, Roo Code 창 오른쪽 위의 **기어 아이콘**을 클릭해
   **Provider Settings**를 엽니다.
2. `API Provider`를 `Ollama`로 설정합니다.
3. (선택) Ollama 인스턴스를 원격에서 실행 중이라면 `Base URL`을 수정합니다. 기본값은
   `http://localhost:11434`입니다.
4. 유효한 `Model ID`를 입력합니다(예: `qwen3` 또는 `qwen3-coder:480b-cloud`).
5. 코딩 작업을 위해 `Context Window`를 최소 32K 토큰 이상으로 조정합니다.

> 참고: 코딩 도구는 비교적 큰 컨텍스트 창을 필요로 합니다. 최소 32K 토큰 이상의 컨텍스트 창 사용을
> 권장합니다. 자세한 내용은 [컨텍스트 길이](/context-length) 문서를 참고하세요.

## ollama.com에 연결하기

1. ollama.com에서 [API 키](https://ollama.com/settings/keys)를 생성합니다.
2. `Use custom base URL`을 켜고 값을 `https://ollama.com`으로 설정합니다.
3. 발급받은 **Ollama API Key**를 입력합니다.
4. 목록에서 모델을 선택합니다.

### 추천 모델

* `qwen3-coder:480b`
* `deepseek-v3.1:671b`

> 원문: https://docs.ollama.com/integrations/roo-code

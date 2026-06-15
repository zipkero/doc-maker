# 컨텍스트 길이

컨텍스트 길이는 모델이 메모리에서 접근할 수 있는 최대 토큰 수를 말합니다.

> 참고: Ollama는 VRAM 용량에 따라 다음 컨텍스트 길이를 기본값으로 사용합니다.
>
> - VRAM 24 GiB 미만: 4k 컨텍스트
> - VRAM 24~48 GiB: 32k 컨텍스트
> - VRAM 48 GiB 이상: 256k 컨텍스트

웹 검색, 에이전트, 코딩 도구처럼 큰 컨텍스트가 필요한 작업은 최소 64000 토큰 이상으로 설정하는 것이 좋습니다.

## 컨텍스트 길이 설정하기

컨텍스트 길이를 크게 잡으면 모델 실행에 필요한 메모리도 늘어납니다.
컨텍스트 길이를 늘리기 전에 충분한 VRAM이 확보되어 있는지 확인하세요.

클라우드 모델은 기본적으로 해당 모델의 최대 컨텍스트 길이로 설정됩니다.

### 앱에서 설정

Ollama 앱의 설정에서 슬라이더를 원하는 컨텍스트 길이로 옮깁니다.

![Ollama 앱의 컨텍스트 길이 설정](https://mintcdn.com/ollama-9269c548/SjntZZpXgbN5v4M5/images/ollama-settings.png?fit=max&auto=format&n=SjntZZpXgbN5v4M5&q=85&s=e8a7ccd30fd9cee5e93662db05b43dc7)

### CLI에서 설정

앱에서 컨텍스트 길이를 바꿀 수 없는 경우, Ollama를 서비스로 띄울 때 컨텍스트 길이를 지정할 수도 있습니다.

```
OLLAMA_CONTEXT_LENGTH=64000 ollama serve
```

### 할당된 컨텍스트 길이와 모델 오프로딩 확인

최상의 성능을 내려면 모델이 지원하는 최대 컨텍스트 길이를 사용하고, 모델이 CPU로 오프로딩되지 않게 하는 것이 좋습니다.
`ollama ps`를 실행해 `PROCESSOR` 항목의 분배 상태를 확인하세요.

```
ollama ps
```

```
NAME             ID              SIZE      PROCESSOR    CONTEXT    UNTIL
gemma4:latest    c6eb396dbd59    9.6 GB    100% GPU     131072     2 minutes from now
```

> 원문: https://docs.ollama.com/context-length

# OpenCode

OpenCode는 터미널에서 실행되는 오픈소스 AI 코딩 어시스턴트입니다.

## 설치

[OpenCode CLI](https://opencode.ai)를 설치합니다.

```bash
curl -fsSL https://opencode.ai/install | bash
```

> 참고: OpenCode는 비교적 큰 컨텍스트 창을 필요로 합니다. 최소 64k 토큰 이상의 컨텍스트 창 사용을
> 권장합니다. 자세한 내용은 [컨텍스트 길이](/context-length) 문서를 참고하세요.

## Ollama와 함께 사용하기

### 빠른 설정

다음 명령으로 바로 실행할 수 있습니다.

```bash
ollama launch opencode
```

실행하지 않고 설정만 하려면 `--config` 옵션을 사용합니다.

```shell
ollama launch opencode --config
```

> 참고: `ollama launch opencode`는 `OPENCODE_CONFIG_CONTENT` 환경변수를 통해 설정을 OpenCode에
> 인라인으로 전달합니다. OpenCode는 시작 시 여러 설정 소스를 깊은 병합(deep-merge)하므로,
> `~/.config/opencode/opencode.json`에 선언한 내용도 그대로 존중되어 OpenCode 안에서 사용할 수
> 있습니다. 다만 `opencode.json`에만 선언된 모델은 `ollama launch`의 모델 선택 메뉴에는 나타나지
> 않습니다.

> 원문: https://docs.ollama.com/integrations/opencode

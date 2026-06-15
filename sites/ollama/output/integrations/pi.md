# Pi

Pi는 간결하면서도 확장 가능한 코딩 에이전트입니다.

## 빠른 설정

다음 명령 한 줄이면 필요할 경우 Pi를 설치하고, Ollama를 웹 도구를 포함한 프로바이더로 구성한 뒤
대화형 세션으로 진입합니다.

```bash
ollama launch pi
```

실행하지 않고 설정만 하려면 `--config` 옵션을 사용합니다.

```shell
ollama launch pi --config
```

### 특정 모델로 바로 실행

`--model` 옵션으로 사용할 모델을 지정해 바로 실행할 수 있습니다.

```shell
ollama launch pi --model qwen3.5:cloud
```

클라우드 모델은 [ollama.com](https://ollama.com/search?c=cloud)에서도 확인할 수 있습니다.

## 확장(Extensions)

Pi는 `read`, `write`, `edit`, `bash` 네 가지 핵심 도구를 기본 제공하며, 그 밖의 모든 기능은 확장
시스템을 통해 추가됩니다. 확장은 필요할 때 `/skill:name` 명령으로 호출하는 기능 패키지입니다.

npm이나 git에서 설치할 수 있습니다.

```bash
pi install npm:@foo/some-tools
pi install git:github.com/user/repo@v1
```

전체 패키지 목록은 [pi.dev](https://pi.dev/packages)에서 볼 수 있습니다.

### 웹 검색

Pi는 `@ollama/pi-web-search` 패키지를 통해 웹 검색 및 페이지 가져오기 도구를 사용할 수 있습니다.
Ollama를 통해 Pi를 실행하면 패키지 설치와 업데이트가 자동으로 관리됩니다. 수동으로 설치하려면 다음과
같이 합니다.

```bash
pi install npm:@ollama/pi-web-search
```

### `pi-autoresearch`를 이용한 자동 연구

[pi-autoresearch](https://github.com/davebcn87/pi-autoresearch)는 자율적인 실험 루프를 Pi에
더해 줍니다. Karpathy의 autoresearch에서 영감을 받아, 측정 가능한 모든 지표(테스트 속도, 번들 크기,
빌드 시간, 모델 학습 손실, Lighthouse 점수 등)를 최적화 대상으로 삼습니다.

```bash
pi install https://github.com/davebcn87/pi-autoresearch
```

최적화할 대상을 Pi에 알려 주기만 하면, Pi가 실험을 돌리고 각각을 벤치마크해 개선은 남기고 퇴보는
되돌리는 과정을 자율적으로 반복합니다. 내장 대시보드는 모든 실행을 추적하며, 신뢰도 점수를 통해 실제
개선과 벤치마크 노이즈를 구분합니다.

```bash
/autoresearch optimize unit test runtime
```

채택된 실험은 자동으로 커밋되고, 실패한 실험은 되돌려집니다. 작업이 끝나면 Pi가 개선 사항들을 독립적인
브랜치로 묶어 깔끔하게 리뷰하고 병합할 수 있도록 정리해 줍니다.

## 수동 설정

### 설치

[Pi](https://github.com/earendil-works/pi)를 설치합니다.

```bash
npm install -g @earendil-works/pi-coding-agent
```

`~/.pi/agent/models.json`에 다음 설정 블록을 추가합니다.

```json
{
  "providers": {
    "ollama": {
      "baseUrl": "http://localhost:11434/v1",
      "api": "openai-completions",
      "apiKey": "ollama",
      "models": [
        {
          "id": "qwen3-coder"
        }
      ]
    }
  }
}
```

`~/.pi/agent/settings.json`을 수정해 기본 프로바이더를 지정합니다.

```json
{
  "defaultProvider": "ollama",
  "defaultModel": "qwen3-coder"
}
```

> 원문: https://docs.ollama.com/integrations/pi

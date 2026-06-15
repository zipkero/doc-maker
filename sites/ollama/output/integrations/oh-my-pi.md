# Oh My Pi

Oh My Pi(OMP)는 IDE 수준의 도구를 내장한 터미널 코딩 에이전트입니다. 채팅, 프로젝트 컨텍스트,
구조화된 코드 편집, 언어 서버 지원, 디버깅 도구, 브라우저 접근, 플러그인, 서브에이전트를 하나의 터미널
워크플로에 모았습니다.

Ollama는 OMP가 Ollama를 모델 제공자로 사용하도록 구성하고 대화형 세션을 실행할 수 있습니다.

## 빠른 설정

```bash
ollama launch omp
```

이 명령은 Ollama를 제공자로 구성하고, 웹 검색 도구를 설정한 뒤 OMP를 시작합니다.

### 특정 모델로 바로 실행

```shell
ollama launch omp --model <model>
```

## 플러그인

OMP는 추가 도구와 기능을 위한 플러그인을 지원합니다. Ollama를 통해 OMP를 실행하면 Ollama 웹 검색
플러그인이 자동으로 관리됩니다.

## 수동 설정

[omp.sh](https://omp.sh)에서 OMP를 설치한 뒤 다음을 실행합니다.

```bash
ollama launch omp --config
```

> 원문: https://docs.ollama.com/integrations/oh-my-pi

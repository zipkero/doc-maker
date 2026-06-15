# Pool

Pool은 Poolside가 만든 터미널용 소프트웨어 에이전트로, 엔터프라이즈 개발 워크플로를 위해 설계되었습니다.

## 설치

[Pool](https://github.com/poolsideai/pool)을 설치합니다.

## Ollama와 함께 사용하기

### 빠른 설정

```shell
ollama launch pool
```

### 특정 모델로 바로 실행

`--model` 옵션으로 사용할 모델을 지정할 수 있습니다.

```shell
ollama launch pool --model kimi-k2.6:cloud
```

### Pool로 인자 전달하기

`--` 뒤에 오는 인자는 Pool에 그대로 전달됩니다.

```shell
ollama launch pool -- --help
```

## 수동 설정

Pool은 환경변수를 통해 OpenAI 호환 API로 Ollama에 연결합니다.

1. 환경변수를 설정합니다.

```shell
export POOLSIDE_STANDALONE_BASE_URL=http://localhost:11434/v1
export POOLSIDE_API_KEY=ollama
```

2. Ollama 모델로 Pool을 실행합니다.

```shell
pool -m kimi-k2.6:cloud
```

또는 환경변수를 인라인으로 함께 지정해 실행할 수도 있습니다.

```shell
POOLSIDE_STANDALONE_BASE_URL=http://localhost:11434/v1 POOLSIDE_API_KEY=ollama pool -m kimi-k2.6:cloud
```

> 원문: https://docs.ollama.com/integrations/pool

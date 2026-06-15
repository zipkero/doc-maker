# 소개

Ollama API를 사용하면 모델을 프로그래밍 방식으로 실행하고 다룰 수 있습니다.

## 시작하기

Ollama를 처음 쓴다면 [quickstart](/quickstart) 문서를 따라 API 환경을 먼저 갖춰 보세요.

## 기본 URL

설치 후 Ollama API는 기본적으로 다음 주소에서 제공됩니다.

```
http://localhost:11434/api
```

**ollama.com**의 클라우드 모델을 사용할 때도 동일한 API를 다음 기본 URL로 호출할 수 있습니다.

```
https://ollama.com/api
```

## 요청 예시

Ollama가 실행 중이면 API가 자동으로 활성화되며 `curl`로 바로 호출할 수 있습니다.

```shell
curl http://localhost:11434/api/generate -d '{
  "model": "gemma4",
  "prompt": "Why is the sky blue?"
}'
```

## 라이브러리

Ollama는 Python과 JavaScript용 공식 라이브러리를 제공합니다.

- [Python](https://github.com/ollama/ollama-python)
- [JavaScript](https://github.com/ollama/ollama-js)

이 외에도 커뮤니티가 관리하는 라이브러리가 여럿 있습니다. 전체 목록은
[Ollama GitHub 저장소](https://github.com/ollama/ollama?tab=readme-ov-file#libraries-1)에서 확인할 수 있습니다.

## 버전 관리

Ollama API는 엄격한 버전 체계를 따르지는 않지만, 안정적이고 하위 호환을 유지하도록
설계되어 있습니다. 기능 폐기(deprecation)는 드물게 일어나며, 발생할 경우
[릴리스 노트](https://github.com/ollama/ollama/releases)를 통해 공지됩니다.

> 원문: https://docs.ollama.com/api/introduction

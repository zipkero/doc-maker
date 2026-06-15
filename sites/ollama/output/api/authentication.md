# 인증

`http://localhost:11434`로 로컬에서 Ollama API에 접근할 때는 인증이 필요하지 않습니다.

다음 경우에는 인증이 필요합니다.

- ollama.com을 통해 클라우드 모델을 실행할 때
- 모델을 게시(publish)할 때
- 비공개 모델을 다운로드할 때

Ollama는 두 가지 인증 방식을 지원합니다.

- **로그인**: 로컬 설치본에서 로그인해 두면, 명령을 실행할 때 ollama.com에 대한 요청을 Ollama가
  자동으로 인증합니다.
- **API 키**: ollama.com API에 프로그램으로 접근하기 위한 API 키입니다.

## 로그인

로컬에 설치한 Ollama에서 ollama.com에 로그인하려면 다음을 실행합니다.

```
ollama signin
```

로그인하면 이후 필요한 명령을 Ollama가 자동으로 인증합니다.

```
ollama run gpt-oss:120b-cloud
```

마찬가지로, 클라우드 접근이 필요한 로컬 API 엔드포인트에 접근할 때도 Ollama가 요청을 자동으로
인증합니다.

```shell
curl http://localhost:11434/api/generate -d '{
  "model": "gpt-oss:120b-cloud",
  "prompt": "Why is the sky blue?"
}'
```

## API 키

`https://ollama.com/api`에서 제공되는 ollama.com API에 직접 접근하려면 API 키를 통한 인증이
필요합니다.

먼저 [API 키](https://ollama.com/settings/keys)를 발급한 뒤, `OLLAMA_API_KEY` 환경 변수를
설정합니다.

```shell
export OLLAMA_API_KEY=your_api_key
```

그런 다음 Authorization 헤더에 API 키를 넣어 사용합니다.

```shell
curl https://ollama.com/api/generate \
  -H "Authorization: Bearer $OLLAMA_API_KEY" \
  -d '{
    "model": "gpt-oss:120b",
    "prompt": "Why is the sky blue?",
    "stream": false
  }'
```

API 키는 현재 만료되지 않지만,
[API 키 설정](https://ollama.com/settings/keys)에서 언제든지 폐기할 수 있습니다.

> 원문: https://docs.ollama.com/api/authentication

# n8n

## 설치

[n8n](https://docs.n8n.io/choose-n8n/)을 설치합니다.

## 로컬 Ollama 사용하기

1. 오른쪽 위 모서리에서 드롭다운을 클릭한 뒤 **Create Credential**을 선택합니다.
2. **Add new credential**에서 **Ollama**를 선택합니다.
3. Base URL을 확인합니다. 로컬에서 실행 중이면 `http://localhost:11434`, Docker로 실행 중이면
   `http://host.docker.internal:11434`로 설정한 뒤 **Save**를 클릭합니다.

   > 참고: Docker Desktop을 쓰지 않는 환경(예: Linux 서버 설치)에서는 `host.docker.internal`이
   > 자동으로 추가되지 않습니다. n8n을 Docker로 실행할 때 `--add-host=host.docker.internal:host-gateway`
   > 옵션을 주거나, docker compose 파일에 다음을 추가하세요.
   >
   > ```yaml
   > extra_hosts:
   >   - "host.docker.internal:host-gateway"
   > ```

   `Connection tested successfully` 메시지가 보여야 합니다.

4. 새 워크플로를 만들 때 **Add a first step**를 선택하고 **Ollama node**를 선택합니다.
5. 원하는 모델(예: `qwen3-coder`)을 선택합니다.

## ollama.com에 연결하기

1. **ollama.com**에서 [API 키](https://ollama.com/settings/keys)를 생성합니다.
2. n8n에서 **Create Credential**을 클릭하고 **Ollama**를 선택합니다.
3. **API URL**을 `https://ollama.com`으로 설정합니다.
4. **API Key**를 입력하고 **Save**를 클릭합니다.

> 원문: https://docs.ollama.com/integrations/n8n

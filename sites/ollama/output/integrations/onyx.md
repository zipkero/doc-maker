# Onyx

## 개요

[Onyx](http://onyx.app/)는 모든 Ollama 모델과 연동되는 자체 호스팅 가능한 채팅 UI입니다. 주요 기능은
다음과 같습니다.

* 사용자 정의 에이전트 생성
* 웹 검색
* 심층 리서치(Deep Research)
* 업로드한 문서와 연결된 앱을 대상으로 한 RAG
* Google Drive, Email, Slack 등 애플리케이션 커넥터
* MCP 및 OpenAPI Actions 지원
* 이미지 생성
* 사용자/그룹 관리, RBAC, SSO 등

Onyx는 단일 사용자부터 대규모 조직까지 다양한 규모로 배포할 수 있습니다.

## Onyx 설치

[퀵스타트 가이드](https://docs.onyx.app/deployment/getting_started/quickstart)를 따라 Onyx를
배포합니다.

> 참고: 리소스 산정과 스케일링 관련 문서는
> [여기](https://docs.onyx.app/deployment/getting_started/resourcing)에서 확인할 수 있습니다.

## Ollama와 함께 사용하기

1. Onyx 배포 환경에 로그인합니다(먼저 계정을 생성하세요).
2. 설정 과정에서 LLM 제공자로 `Ollama`를 선택합니다.
3. **Ollama API URL**을 입력하고 사용할 모델을 선택합니다.

   > 참고: Onyx를 Docker로 실행 중이라면 컴퓨터의 로컬 네트워크에 접근하기 위해 `http://127.0.0.1`
   > 대신 `http://host.docker.internal`을 사용하세요.

설정의 `Ollama Cloud` 탭에서 Onyx Cloud를 손쉽게 연결할 수도 있습니다.

## 첫 질의 보내기

설정이 끝나면 Onyx에서 모델에 바로 질의를 보낼 수 있습니다.

> 원문: https://docs.ollama.com/integrations/onyx

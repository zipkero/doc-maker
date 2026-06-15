# Hermes Desktop

Hermes Desktop은 Nous Research가 만든 네이티브 AI 어시스턴트 앱입니다. Hermes Agent를 위한 데스크톱
채팅 인터페이스를 제공하며, 이 에이전트는 모델을 다루고, 도구를 실행하고, 프로젝트를 관리하고, 메모리와
스킬을 활용하고, 메시징 게이트웨이에 연결할 수 있습니다.

## 빠른 시작

```bash
ollama launch hermes-desktop
```

Ollama가 설정 과정을 자동으로 처리합니다.

1. **설치** - Hermes Desktop이 설치되어 있지 않으면 Ollama가 설치를 묻습니다.
2. **모델** - 선택기에서 모델을 고릅니다.
3. **구성** - Ollama가 선택한 Ollama 모델을 사용하도록 Hermes Desktop을 설정합니다.
4. **실행** - Ollama가 Hermes Desktop을 엽니다.

## 특정 모델로 바로 실행

```bash
ollama launch hermes-desktop --model <model>
```

나중에 모델을 바꾸려면 `ollama launch hermes-desktop`을 다시 실행하면 됩니다.

> 원문: https://docs.ollama.com/integrations/hermes-desktop

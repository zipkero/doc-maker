# Hermes Agent

Hermes Agent는 Nous Research가 만든 자기 개선형 AI 에이전트입니다. 스킬을 자동으로 만들어내고,
세션 간에 메모리를 유지하며, 기본으로 70개 이상의 스킬을 제공합니다.

## 빠른 시작

```bash
ollama launch hermes
```

이 명령 하나로 Ollama가 다음을 자동으로 처리합니다.

1. **설치** — Hermes가 설치되어 있지 않으면 Nous Research 설치 스크립트로 설치하도록 안내합니다.
2. **모델** — 선택기에서 로컬 또는 클라우드 모델을 고릅니다.
3. **온보딩** — Ollama 제공자를 설정하고 Hermes가 `http://127.0.0.1:11434/v1`을 바라보도록 한 뒤, 선택한 모델을 기본 모델로 지정합니다.
4. **게이트웨이** — 선택적으로 메시징 플랫폼(Telegram, Discord, Slack, WhatsApp, Signal, Email)을 연결하고 Hermes 채팅을 실행합니다.

## 권장 모델

**클라우드 모델**

* `kimi-k2.5:cloud` — 서브에이전트를 갖춘 멀티모달 추론
* `glm-5.1:cloud` — 추론 및 코드 생성
* `qwen3.5:cloud` — 추론, 코딩, 비전 기반 에이전트형 도구 활용
* `minimax-m2.7:cloud` — 빠르고 효율적인 코딩 및 실무 생산성

**로컬 모델**

* `gemma4` — 로컬 환경에서 추론 및 코드 생성(VRAM 약 16 GB)
* `qwen3.6` — 로컬 환경에서 추론, 코딩, 시각 이해(VRAM 약 24 GB)

더 많은 모델은 [ollama.com/search](https://ollama.com/search?c=cloud)에서 확인할 수 있습니다.

## 메시징 앱 연결

Telegram, Discord, Slack, WhatsApp, Signal, Email을 연결하면 어디서나 모델과 대화할 수 있습니다.

```bash
hermes gateway setup
```

## 재설정

언제든 전체 설정 마법사를 다시 실행할 수 있습니다.

```bash
hermes setup
```

## 수동 설정

`ollama launch hermes` 대신 Hermes 자체 마법사로 진행하고 싶다면, Hermes를 직접 설치합니다.

```bash
curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash
```

설치하면 Hermes가 설정 마법사를 자동으로 띄웁니다. **Quick setup**을 선택하세요.

```
How would you like to set up Hermes?

 →  Quick setup — provider, model & messaging (recommended)
    Full setup — configure everything
```

### Ollama 연결

1. **More providers...**를 선택합니다.

2. **Custom endpoint (enter URL manually)**를 선택합니다.

3. API base URL을 Ollama의 OpenAI 호환 엔드포인트로 지정합니다.

   ```
   API base URL [e.g. https://api.example.com/v1]: http://127.0.0.1:11434/v1
   ```

4. 로컬 Ollama에서는 API 키가 필요 없으므로 비워 둡니다.

   ```
   API key [optional]:
   ```

5. Hermes가 내려받은 모델을 자동으로 감지하므로, 사용할 모델을 확인합니다.

   ```
   Verified endpoint via http://127.0.0.1:11434/v1/models (1 model(s) visible)
     Detected model: kimi-k2.5:cloud
     Use this model? [Y/n]:
   ```

6. 컨텍스트 길이는 비워 두면 자동으로 감지됩니다.

   ```
   Context length in tokens [leave blank for auto-detect]:
   ```

### 메시징 연결

설정 도중 메시징 플랫폼을 선택적으로 연결할 수 있습니다.

```
Connect a messaging platform? (Telegram, Discord, etc.)

 →  Set up messaging now (recommended)
    Skip — set up later with 'hermes setup gateway'
```

### 실행

```
Launch hermes chat now? [Y/n]: Y
```

> 원문: https://docs.ollama.com/integrations/hermes

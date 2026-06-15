# NemoClaw

NemoClaw는 [OpenClaw](/integrations/openclaw)를 위한 NVIDIA의 오픈소스 보안 스택입니다.
OpenClaw를 NVIDIA OpenShell 런타임으로 감싸, AI 에이전트에 커널 수준 샌드박싱, 네트워크 정책 제어,
감사 로그를 제공합니다.

## 빠른 시작

모델을 받습니다.

```bash
ollama pull nemotron-3-nano:30b
```

설치 스크립트를 실행합니다.

```bash
curl -fsSL https://www.nvidia.com/nemoclaw.sh | \
  NEMOCLAW_NON_INTERACTIVE=1 \
  NEMOCLAW_PROVIDER=ollama \
  NEMOCLAW_MODEL=nemotron-3-nano:30b \
  bash
```

샌드박스에 연결합니다.

```bash
nemoclaw my-assistant connect
```

TUI를 엽니다.

```bash
openclaw tui
```

> 참고: NemoClaw의 Ollama 지원은 아직 실험 단계입니다.

## 플랫폼 지원

| 플랫폼 | 런타임 | 상태 |
| --- | --- | --- |
| Linux (Ubuntu 22.04+) | Docker | 주력(Primary) |
| macOS (Apple Silicon) | Colima 또는 Docker Desktop | 지원 |
| Windows | WSL2 + Docker Desktop | 지원 |

Windows에서는 CMD와 PowerShell이 지원되지 않으며 WSL2가 필요합니다.

> 참고: 설치 스크립트를 실행하기 전에 Ollama가 설치되어 실행 중이어야 합니다. WSL2나 컨테이너 안에서
> 실행할 때는 샌드박스에서 Ollama에 접근할 수 있도록 설정하세요(예: `OLLAMA_HOST=0.0.0.0`).

## 시스템 요구 사항

* CPU: 최소 4 vCPU
* RAM: 최소 8 GB(16 GB 권장)
* 디스크: 여유 공간 20 GB(로컬 모델 사용 시 40 GB 권장)
* Node.js 20+ 및 npm 10+
* 컨테이너 런타임(Docker 권장)

## 권장 모델

* `nemotron-3-super:cloud` — 뛰어난 추론과 코딩
* `qwen3.5:cloud` — 397B, 추론 및 코드 생성
* `nemotron-3-nano:30b` — 권장 로컬 모델, VRAM 24 GB에 적합
* `qwen3.5:27b` — 빠른 로컬 추론(VRAM 약 18 GB)
* `glm-4.7-flash` — 추론 및 코드 생성(VRAM 약 25 GB)

더 많은 모델은 [ollama.com/search](https://ollama.com/search)에서 확인할 수 있습니다.

> 원문: https://docs.ollama.com/integrations/nemoclaw

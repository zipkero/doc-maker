# 자주 묻는 질문

## Ollama를 어떻게 업그레이드하나요?

macOS와 Windows에서는 Ollama가 업데이트를 자동으로 내려받습니다. 작업 표시줄이나 메뉴 막대의 아이콘을
클릭한 뒤 "Restart to update"를 누르면 업데이트가 적용됩니다. [다운로드 페이지](https://ollama.com/download/)에서
최신 버전을 직접 받아 설치할 수도 있습니다.

Linux에서는 설치 스크립트를 다시 실행하면 됩니다.

```shell
curl -fsSL https://ollama.com/install.sh | sh
```

## 로그는 어떻게 확인하나요?

로그 활용 방법은 [문제 해결](./troubleshooting.mdx) 문서를 참고하세요.

## 내 GPU가 Ollama와 호환되나요?

[GPU 문서](./gpu.mdx)를 참고하세요.

## 컨텍스트 창 크기는 어떻게 지정하나요?

기본적으로 Ollama는 4096 토큰의 컨텍스트 창 크기를 사용합니다.

이 값은 `OLLAMA_CONTEXT_LENGTH` 환경 변수로 재정의할 수 있습니다. 예를 들어 기본 컨텍스트 창을 8K로
설정하려면 다음과 같이 실행합니다.

```shell
OLLAMA_CONTEXT_LENGTH=8192 ollama serve
```

`ollama run` 사용 중에 바꾸려면 `/set parameter`를 씁니다.

```shell
/set parameter num_ctx 4096
```

API를 사용할 때는 `num_ctx` 파라미터를 지정합니다.

```shell
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Why is the sky blue?",
  "options": {
    "num_ctx": 4096
  }
}'
```

## 모델이 GPU에 올라갔는지 어떻게 확인하나요?

`ollama ps` 명령으로 현재 메모리에 적재된 모델을 확인할 수 있습니다.

```shell
ollama ps
```

출력 예시:

```
NAME        ID            SIZE    PROCESSOR   UNTIL
llama3:70b  bcfb190ca3a7  42 GB   100% GPU    4 minutes from now
```

`Processor` 열에 모델이 어느 메모리에 적재됐는지 표시됩니다.

- `100% GPU` — 모델 전체가 GPU에 적재됨
- `100% CPU` — 모델 전체가 시스템 메모리에 적재됨
- `48%/52% CPU/GPU` — GPU와 시스템 메모리에 나뉘어 적재됨

## Ollama 서버는 어떻게 설정하나요?

Ollama 서버는 환경 변수로 설정합니다.

### macOS에서 환경 변수 설정

Ollama를 macOS 애플리케이션으로 실행하는 경우 `launchctl`로 환경 변수를 설정합니다.

1. 설정할 환경 변수마다 `launchctl setenv`를 호출합니다.

   ```bash
   launchctl setenv OLLAMA_HOST "0.0.0.0:11434"
   ```

2. Ollama 애플리케이션을 다시 시작합니다.

### Linux에서 환경 변수 설정

Ollama를 systemd 서비스로 실행하는 경우 `systemctl`로 환경 변수를 설정합니다.

1. `systemctl edit ollama.service`를 실행해 systemd 서비스를 편집합니다. 편집기가 열립니다.

2. 설정할 환경 변수마다 `[Service]` 섹션 아래에 `Environment` 줄을 추가합니다.

   ```ini
   [Service]
   Environment="OLLAMA_HOST=0.0.0.0:11434"
   ```

3. 저장하고 종료합니다.

4. `systemd`를 다시 로드하고 Ollama를 재시작합니다.

   ```shell
   systemctl daemon-reload
   systemctl restart ollama
   ```

### Windows에서 환경 변수 설정

Windows에서 Ollama는 사용자 및 시스템 환경 변수를 그대로 상속합니다.

1. 먼저 작업 표시줄의 Ollama 아이콘을 클릭해 종료합니다.

2. 설정(Windows 11) 또는 제어판(Windows 10)을 열고 *environment variables*(환경 변수)를 검색합니다.

3. *Edit environment variables for your account*(내 계정의 환경 변수 편집)를 클릭합니다.

4. 사용자 계정에 대해 `OLLAMA_HOST`, `OLLAMA_MODELS` 등의 변수를 편집하거나 새로 만듭니다.

5. OK/적용을 눌러 저장합니다.

6. Windows 시작 메뉴에서 Ollama 애플리케이션을 시작합니다.

## 프록시 뒤에서 Ollama를 어떻게 사용하나요?

Ollama는 인터넷에서 모델을 받아오므로 모델 접근에 프록시 서버가 필요할 수 있습니다. `HTTPS_PROXY`를
사용해 외부로 나가는 요청을 프록시로 우회시키세요. 프록시 인증서는 시스템 인증서로 설치되어 있어야 합니다.
플랫폼별 환경 변수 설정 방법은 위 섹션을 참고하세요.

> 참고: `HTTP_PROXY`는 설정하지 마세요. Ollama는 모델을 받아올 때 HTTP가 아닌 HTTPS만 사용합니다.
> `HTTP_PROXY`를 설정하면 클라이언트와 서버 간 연결이 끊길 수 있습니다.

### Docker에서 프록시 뒤로 Ollama를 사용하려면?

Ollama Docker 컨테이너 이미지는 컨테이너를 시작할 때 `-e HTTPS_PROXY=https://proxy.example.com`을
전달해 프록시를 사용하도록 설정할 수 있습니다.

또는 Docker 데몬 자체를 프록시를 사용하도록 설정할 수도 있습니다. 설정 방법은 Docker Desktop의
[macOS](https://docs.docker.com/desktop/settings/mac/#proxies),
[Windows](https://docs.docker.com/desktop/settings/windows/#proxies),
[Linux](https://docs.docker.com/desktop/settings/linux/#proxies) 문서와
[systemd 기반 Docker 데몬](https://docs.docker.com/config/daemon/systemd/#httphttps-proxy) 문서에
나와 있습니다.

HTTPS를 사용할 때는 인증서를 시스템 인증서로 설치해야 합니다. 자체 서명 인증서를 사용한다면 새 Docker
이미지를 만들어야 할 수 있습니다.

```dockerfile
FROM ollama/ollama
COPY my-ca.pem /usr/local/share/ca-certificates/my-ca.crt
RUN update-ca-certificates
```

이 이미지를 빌드하고 실행합니다.

```shell
docker build -t ollama-with-ca .
docker run -d -e HTTPS_PROXY=https://my.proxy.example.com -p 11434:11434 ollama-with-ca
```

## Ollama가 내 프롬프트와 답변을 ollama.com으로 보내나요?

Ollama는 로컬에서 실행됩니다. 로컬로 실행하는 동안에는 여러분의 프롬프트나 데이터를 저희가 볼 수 없습니다.
클라우드 호스팅 모델을 사용할 때는 서비스 제공을 위해 프롬프트와 응답을 처리하지만, 그 내용을 저장하거나
로그로 남기지 않으며 학습에도 절대 사용하지 않습니다. 서비스 제공에 필요한 기본 계정 정보와 제한된 사용
메타데이터를 수집하지만, 여기에 프롬프트나 응답 내용은 포함되지 않습니다. 데이터를 판매하지 않으며, 계정은
언제든 삭제할 수 있습니다.

## Ollama의 클라우드 기능을 어떻게 끄나요?

Ollama는 클라우드 기능을 비활성화해 로컬 전용 모드로 실행할 수 있습니다. 클라우드 기능을 끄면 Ollama의
클라우드 모델과 웹 검색은 사용할 수 없게 됩니다.

`~/.ollama/server.json`에 `disable_ollama_cloud`를 설정합니다.

```json
{
  "disable_ollama_cloud": true
}
```

환경 변수로 설정할 수도 있습니다.

```shell
OLLAMA_NO_CLOUD=1
```

설정을 바꾼 뒤에는 Ollama를 재시작하세요. 비활성화되면 Ollama 로그에 `Ollama cloud disabled: true`가
표시됩니다.

## Ollama를 네트워크에 노출하려면 어떻게 하나요?

Ollama는 기본적으로 `127.0.0.1`의 11434 포트에 바인딩됩니다. 바인딩 주소는 `OLLAMA_HOST` 환경 변수로
바꿀 수 있습니다.

플랫폼별 환경 변수 설정 방법은 [위 섹션](#ollama-서버는-어떻게-설정하나요)을 참고하세요.

## 프록시 서버와 함께 Ollama를 사용하려면?

Ollama는 HTTP 서버를 실행하므로 Nginx 같은 프록시 서버로 노출할 수 있습니다. 프록시가 요청을 전달하도록
설정하고, (네트워크에 노출하지 않는 경우) 필요한 헤더를 함께 지정하면 됩니다. 예를 들어 Nginx에서는 다음과
같이 설정합니다.

```nginx
server {
    listen 80;
    server_name example.com;  # Replace with your domain or IP
    location / {
        proxy_pass http://localhost:11434;
        proxy_set_header Host localhost:11434;
    }
}
```

## ngrok과 함께 Ollama를 사용하려면?

Ollama는 다양한 터널링 앱으로 접근할 수 있습니다. 예를 들어 Ngrok을 사용하면 다음과 같습니다.

```shell
ngrok http 11434 --host-header="localhost:11434"
```

## Cloudflare Tunnel과 함께 Ollama를 사용하려면?

Cloudflare Tunnel과 함께 Ollama를 사용하려면 `--url`과 `--http-host-header` 플래그를 사용합니다.

```shell
cloudflared tunnel --url http://localhost:11434 --http-host-header="localhost:11434"
```

## 추가 웹 출처(origin)에서 Ollama에 접근하도록 허용하려면?

Ollama는 기본적으로 `127.0.0.1`과 `0.0.0.0`에서 오는 교차 출처(cross-origin) 요청을 허용합니다. 추가
출처는 `OLLAMA_ORIGINS`로 설정할 수 있습니다.

브라우저 확장 프로그램의 경우 해당 확장의 출처 패턴을 명시적으로 허용해야 합니다. 모든 브라우저 확장의
접근을 허용하려면 `OLLAMA_ORIGINS`에 `chrome-extension://*`, `moz-extension://*`,
`safari-web-extension://*`를 포함시키고, 특정 확장만 허용하려면 해당 출처를 지정하세요.

```
# Allow all Chrome, Firefox, and Safari extensions
OLLAMA_ORIGINS=chrome-extension://*,moz-extension://*,safari-web-extension://* ollama serve
```

플랫폼별 환경 변수 설정 방법은 [위 섹션](#ollama-서버는-어떻게-설정하나요)을 참고하세요.

## 모델은 어디에 저장되나요?

- macOS: `~/.ollama/models`
- Linux: `/usr/share/ollama/.ollama/models`
- Windows: `C:\Users\%username%\.ollama\models`

### 다른 위치로 바꾸려면?

다른 디렉터리를 사용하려면 `OLLAMA_MODELS` 환경 변수를 원하는 디렉터리로 설정하세요.

> 참고: 표준 설치 프로그램으로 설치한 Linux에서는 `ollama` 사용자가 지정한 디렉터리에 대해 읽기·쓰기
> 권한을 가져야 합니다. 해당 디렉터리를 `ollama` 사용자에게 할당하려면
> `sudo chown -R ollama:ollama <directory>`를 실행하세요.

플랫폼별 환경 변수 설정 방법은 [위 섹션](#ollama-서버는-어떻게-설정하나요)을 참고하세요.

## Visual Studio Code에서 Ollama를 사용하려면?

VS Code를 비롯한 여러 편집기에는 이미 Ollama를 활용하는 플러그인이 많이 나와 있습니다. 메인 저장소
readme 하단의 [확장 프로그램 및 플러그인 목록](https://github.com/ollama/ollama#extensions--plugins)을
참고하세요.

## Docker에서 GPU 가속과 함께 Ollama를 사용하려면?

Ollama Docker 컨테이너는 Linux 또는 Windows(WSL2 사용)에서 GPU 가속을 사용하도록 설정할 수 있습니다.
이를 위해서는 [nvidia-container-toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)이
필요합니다. 자세한 내용은 [ollama/ollama](https://hub.docker.com/r/ollama/ollama)를 참고하세요.

macOS의 Docker Desktop에서는 GPU 패스스루와 에뮬레이션이 지원되지 않아 GPU 가속을 사용할 수 없습니다.

## Windows 10의 WSL2에서 네트워크가 느린 이유는?

이 문제는 Ollama 설치와 모델 다운로드 모두에 영향을 줄 수 있습니다.

`제어판 > 네트워크 및 인터넷 > 네트워크 상태 및 작업 보기`를 열고 왼쪽 패널에서 `어댑터 설정 변경`을
클릭합니다. `vEthernet (WSL)` 어댑터를 찾아 마우스 오른쪽 버튼을 누르고 `속성`을 선택합니다. `구성`을
클릭한 뒤 `고급` 탭을 엽니다. 속성 목록에서 `Large Send Offload Version 2 (IPv4)`와
`Large Send Offload Version 2 (IPv6)`를 찾아 둘 다 *사용 안 함*으로 설정합니다.

## 더 빠른 응답을 위해 모델을 미리 적재(preload)하려면?

API를 사용한다면 Ollama 서버에 빈 요청을 보내 모델을 미리 적재할 수 있습니다. `/api/generate`와
`/api/chat` 엔드포인트 모두에서 동작합니다.

generate 엔드포인트로 mistral 모델을 미리 적재하려면 다음과 같이 합니다.

```shell
curl http://localhost:11434/api/generate -d '{"model": "mistral"}'
```

채팅 완성 엔드포인트를 사용하려면 다음과 같이 합니다.

```shell
curl http://localhost:11434/api/chat -d '{"model": "mistral"}'
```

CLI로 모델을 미리 적재하려면 다음 명령을 사용합니다.

```shell
ollama run llama3.2 ""
```

## 모델을 메모리에 유지하거나 즉시 내리려면?

기본적으로 모델은 5분간 메모리에 유지된 뒤 내려갑니다. LLM에 여러 요청을 연이어 보낼 때 응답 속도를
높여 줍니다. 모델을 메모리에서 즉시 내리려면 `ollama stop` 명령을 사용합니다.

```shell
ollama stop llama3.2
```

API를 사용한다면 `/api/generate`와 `/api/chat` 엔드포인트에서 `keep_alive` 파라미터로 모델이 메모리에
머무는 시간을 설정할 수 있습니다. `keep_alive`에는 다음 값을 지정할 수 있습니다.

- 기간 문자열(예: "10m", "24h")
- 초 단위 숫자(예: 3600)
- 임의의 음수 — 모델을 메모리에 계속 유지(예: -1 또는 "-1m")
- '0' — 응답 생성 직후 모델을 즉시 내림

예를 들어 모델을 미리 적재하고 메모리에 계속 두려면 다음과 같이 합니다.

```shell
curl http://localhost:11434/api/generate -d '{"model": "llama3.2", "keep_alive": -1}'
```

모델을 내리고 메모리를 해제하려면 다음과 같이 합니다.

```shell
curl http://localhost:11434/api/generate -d '{"model": "llama3.2", "keep_alive": 0}'
```

또는 Ollama 서버를 시작할 때 `OLLAMA_KEEP_ALIVE` 환경 변수를 설정해 모든 모델이 메모리에 유지되는
시간을 바꿀 수도 있습니다. `OLLAMA_KEEP_ALIVE`는 위에서 설명한 `keep_alive`와 같은 형식의 값을 받습니다.
환경 변수 설정 방법은 [Ollama 서버 설정 방법](#ollama-서버는-어떻게-설정하나요) 섹션을 참고하세요.

`/api/generate`와 `/api/chat` 엔드포인트의 `keep_alive` 파라미터는 `OLLAMA_KEEP_ALIVE` 설정보다
우선합니다.

## Ollama 서버가 대기열에 담을 수 있는 최대 요청 수를 관리하려면?

서버에 너무 많은 요청이 들어오면 서버 과부하를 알리는 503 오류로 응답합니다. 대기열에 담을 수 있는 요청
수는 `OLLAMA_MAX_QUEUE`로 조정할 수 있습니다.

## Ollama는 동시 요청을 어떻게 처리하나요?

Ollama는 두 가지 수준의 동시 처리를 지원합니다. 시스템에 메모리가 충분하면(CPU 추론 시 시스템 메모리,
GPU 추론 시 VRAM) 여러 모델을 동시에 적재할 수 있습니다. 특정 모델에 대해 적재 시점에 메모리가 충분하면
병렬 요청 처리가 가능하도록 구성됩니다.

이미 하나 이상의 모델이 적재된 상태에서 새 모델을 적재할 메모리가 부족하면, 새 모델을 적재할 수 있을 때까지
모든 새 요청이 대기열에 들어갑니다. 기존 모델이 유휴 상태가 되면 새 모델 공간을 확보하기 위해 하나 이상이
내려가고, 대기 중인 요청은 순서대로 처리됩니다. GPU 추론에서는 동시 모델 적재를 허용하려면 새 모델이
VRAM에 완전히 들어가야 합니다.

특정 모델에 대한 병렬 요청 처리는 병렬 요청 수만큼 컨텍스트 크기를 늘립니다. 예를 들어 2K 컨텍스트에 4개의
병렬 요청이 있으면 컨텍스트가 8K로 늘어나고 그만큼 메모리가 추가로 할당됩니다.

대부분의 플랫폼에서 다음 서버 설정으로 동시 요청 처리 방식을 조정할 수 있습니다.

- `OLLAMA_MAX_LOADED_MODELS` — 가용 메모리에 들어가는 한 동시에 적재할 수 있는 최대 모델 수. 기본값은
  GPU 수의 3배, CPU 추론 시에는 3.
- `OLLAMA_NUM_PARALLEL` — 각 모델이 동시에 처리하는 최대 병렬 요청 수. 기본값 1. 필요한 RAM은
  `OLLAMA_NUM_PARALLEL` × `OLLAMA_CONTEXT_LENGTH`에 비례해 늘어납니다.
- `OLLAMA_MAX_QUEUE` — 서버가 바쁠 때 추가 요청을 거부하기 전까지 대기열에 담는 최대 요청 수. 기본값 512.

참고: Radeon GPU를 사용하는 Windows는 현재 ROCm v5.7의 가용 VRAM 보고 제약 때문에 기본적으로 최대 1개
모델로 제한됩니다. ROCm v6.2가 제공되면 Windows Radeon도 위 기본값을 따릅니다. Windows의 Radeon에서
동시 모델 적재를 직접 켤 수 있지만, GPU의 VRAM에 들어가는 것보다 많은 모델을 적재하지 않도록 주의하세요.

## Ollama는 여러 GPU에 모델을 어떻게 적재하나요?

새 모델을 적재할 때 Ollama는 모델에 필요한 VRAM을 현재 가용 VRAM과 비교합니다. 모델이 단일 GPU 하나에
모두 들어가면 해당 GPU에 적재합니다. 추론 중 PCI 버스를 통한 데이터 전송량을 줄여 주므로 보통 가장 좋은
성능을 냅니다. 모델이 GPU 하나에 모두 들어가지 않으면 사용 가능한 모든 GPU에 나눠 적재합니다.

## Flash Attention은 어떻게 활성화하나요?

Flash Attention은 대부분의 최신 모델이 갖춘 기능으로, 컨텍스트 크기가 커질수록 메모리 사용량을 크게 줄여
줍니다. 활성화하려면 Ollama 서버를 시작할 때 `OLLAMA_FLASH_ATTENTION` 환경 변수를 `1`로 설정하세요.

## K/V 캐시의 양자화 유형은 어떻게 설정하나요?

Flash Attention이 활성화된 상태에서는 K/V 컨텍스트 캐시를 양자화해 메모리 사용량을 크게 줄일 수 있습니다.

양자화된 K/V 캐시를 사용하려면 다음 환경 변수를 설정합니다.

- `OLLAMA_KV_CACHE_TYPE` — K/V 캐시의 양자화 유형. 기본값은 `f16`.

> 참고: 현재 이 설정은 전역 옵션이라 모든 모델이 지정한 양자화 유형으로 실행됩니다.

현재 사용할 수 있는 K/V 캐시 양자화 유형은 다음과 같습니다.

- `f16` — 높은 정밀도와 메모리 사용량(기본값).
- `q8_0` — 8비트 양자화. `f16`의 약 1/2 메모리를 사용하며 정밀도 손실이 매우 작아 보통 모델 품질에
  눈에 띄는 영향이 없습니다(f16을 쓰지 않는다면 권장).
- `q4_0` — 4비트 양자화. `f16`의 약 1/4 메모리를 사용하지만 정밀도 손실이 중소 수준이며, 컨텍스트가
  커질수록 더 두드러질 수 있습니다.

캐시 양자화가 모델 응답 품질에 미치는 영향은 모델과 작업에 따라 다릅니다. GQA 수가 높은 모델(예: Qwen2)은
GQA 수가 낮은 모델보다 양자화로 인한 정밀도 영향이 더 클 수 있습니다.

메모리 사용량과 품질 사이의 최적 균형을 찾으려면 여러 양자화 유형을 직접 시험해 봐야 할 수 있습니다.

## Ollama 공개 키는 어디서 찾나요?

**Ollama 공개 키**는 로컬 Ollama 인스턴스가 [ollama.com](https://ollama.com)과 통신할 수 있게 해 주는
키 쌍의 공개 부분입니다.

다음 작업에 필요합니다.

- Ollama에 모델 push
- Ollama에서 비공개 모델을 내 컴퓨터로 pull
- [Ollama Cloud](https://ollama.com/cloud)에 호스팅된 모델 실행

### 키를 추가하는 방법

- **Mac 및 Windows 앱의 설정 페이지에서 로그인**

- **CLI로 로그인**

```shell
ollama signin
```

- **Ollama Keys 페이지에서 직접 복사·붙여넣기**:
  [https://ollama.com/settings/keys](https://ollama.com/settings/keys)

### Ollama 공개 키의 위치

| OS      | `id_ed25519.pub` 경로                        |
| :------ | :------------------------------------------- |
| macOS   | `~/.ollama/id_ed25519.pub`                   |
| Linux   | `/usr/share/ollama/.ollama/id_ed25519.pub`   |
| Windows | `C:\Users\<username>\.ollama\id_ed25519.pub` |

> 참고: `<username>`은 실제 Windows 사용자 이름으로 바꾸세요.

## 컴퓨터에 로그인할 때 Ollama가 자동 시작되지 않게 하려면?

Windows와 macOS용 Ollama는 설치 시 로그인 항목으로 등록됩니다. 자동 시작을 원하지 않으면 이를 비활성화할
수 있습니다. 애플리케이션을 제거하지 않는 한, Ollama는 업그레이드 후에도 이 설정을 유지합니다.

**Windows**

- `작업 관리자`의 `시작 프로그램` 탭에서 `ollama`를 찾아 `사용 안 함`을 클릭합니다.

**macOS**

- `설정`을 열고 "로그인 항목"을 검색한 뒤, `백그라운드에서 허용` 아래의 `Ollama` 항목을 찾아 슬라이더를
  꺼서 비활성화합니다.

> 원문: https://docs.ollama.com/faq

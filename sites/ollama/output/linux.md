# Linux

Linux에서 Ollama를 설치하고 운영하는 방법을 설명합니다.

## 설치

다음 명령으로 Ollama를 설치합니다.

```shell
curl -fsSL https://ollama.com/install.sh | sh
```

## 수동 설치

> 참고: 이전 버전에서 업그레이드하는 경우, 먼저 `sudo rm -rf /usr/lib/ollama`로 오래된 라이브러리를 제거하세요.

패키지를 내려받아 압축을 풉니다.

```shell
curl -fsSL https://ollama.com/download/ollama-linux-amd64.tar.zst \
    | sudo tar x -C /usr
```

Ollama를 실행합니다.

```shell
ollama serve
```

다른 터미널에서 Ollama가 실행 중인지 확인합니다.

```shell
ollama -v
```

### AMD GPU 설치

AMD GPU를 사용한다면 추가로 ROCm 패키지도 내려받아 압축을 풉니다.

```shell
curl -fsSL https://ollama.com/download/ollama-linux-amd64-rocm.tar.zst \
    | sudo tar x -C /usr
```

### ARM64 설치

ARM64 전용 패키지를 내려받아 압축을 풉니다.

```shell
curl -fsSL https://ollama.com/download/ollama-linux-arm64.tar.zst \
    | sudo tar x -C /usr
```

### 시작 서비스로 등록 (권장)

Ollama 전용 사용자와 그룹을 만듭니다.

```shell
sudo useradd -r -s /bin/false -U -m -d /usr/share/ollama ollama
sudo usermod -a -G ollama $(whoami)
```

`/etc/systemd/system/ollama.service` 위치에 서비스 파일을 만듭니다.

```ini
[Unit]
Description=Ollama Service
After=network-online.target

[Service]
ExecStart=/usr/bin/ollama serve
User=ollama
Group=ollama
Restart=always
RestartSec=3
Environment="PATH=$PATH"

[Install]
WantedBy=multi-user.target
```

그런 다음 서비스를 시작합니다.

```shell
sudo systemctl daemon-reload
sudo systemctl enable ollama
```

### CUDA 드라이버 설치 (선택)

[CUDA를 내려받아 설치](https://developer.nvidia.com/cuda-downloads)합니다.

다음 명령으로 드라이버가 설치되었는지 확인합니다. 정상이라면 GPU 정보가 출력됩니다.

```shell
nvidia-smi
```

### AMD ROCm 드라이버 설치 (선택)

[ROCm v7을 내려받아 설치](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/tutorial/quick-start.html)합니다.

### Ollama 시작

Ollama를 시작하고 실행 상태를 확인합니다.

```shell
sudo systemctl start ollama
sudo systemctl status ollama
```

> 참고: AMD가 공식 Linux 커널 소스에 `amdgpu` 드라이버를 기여하긴 했지만, 그 버전은 오래되어
> 모든 ROCm 기능을 지원하지 않을 수 있습니다. Radeon GPU를 가장 잘 지원하려면
> [AMD 공식 페이지](https://www.amd.com/en/support/linux-drivers)에서 최신 드라이버를
> 설치하는 것을 권장합니다.

## 설정 변경

설치된 Ollama를 사용자화하려면 systemd 서비스 파일이나 환경변수를 편집합니다.

```shell
sudo systemctl edit ollama
```

또는 `/etc/systemd/system/ollama.service.d/override.conf`에 오버라이드 파일을 직접 만듭니다.

```ini
[Service]
Environment="OLLAMA_DEBUG=1"
```

## 업데이트

설치 스크립트를 다시 실행하면 Ollama가 업데이트됩니다.

```shell
curl -fsSL https://ollama.com/install.sh | sh
```

또는 패키지를 다시 내려받아도 됩니다.

```shell
curl -fsSL https://ollama.com/download/ollama-linux-amd64.tar.zst \
    | sudo tar x -C /usr
```

## 특정 버전 설치

설치 스크립트에 `OLLAMA_VERSION` 환경변수를 지정하면 사전 릴리스를 포함한 특정 버전을 설치할 수 있습니다.
버전 번호는 [릴리스 페이지](https://github.com/ollama/ollama/releases)에서 확인할 수 있습니다.

예시:

```shell
curl -fsSL https://ollama.com/install.sh | OLLAMA_VERSION=0.5.7 sh
```

## 로그 보기

시작 서비스로 실행 중인 Ollama의 로그를 보려면 다음을 실행합니다.

```shell
journalctl -e -u ollama
```

## 제거

Ollama 서비스를 제거합니다.

```shell
sudo systemctl stop ollama
sudo systemctl disable ollama
sudo rm /etc/systemd/system/ollama.service
```

라이브러리 디렉터리(`/usr/local/lib`, `/usr/lib`, `/lib` 중 하나)에서 Ollama 라이브러리를 제거합니다.

```shell
sudo rm -r $(which ollama | tr 'bin' 'lib')
```

bin 디렉터리(`/usr/local/bin`, `/usr/bin`, `/bin` 중 하나)에서 Ollama 바이너리를 제거합니다.

```shell
sudo rm $(which ollama)
```

내려받은 모델과 Ollama 서비스 사용자·그룹을 제거합니다.

```shell
sudo userdel ollama
sudo groupdel ollama
sudo rm -r /usr/share/ollama
```

> 원문: https://docs.ollama.com/linux

# 문제 해결

Ollama가 예상대로 동작하지 않을 때 원인을 파악하고 해결하는 방법을 설명합니다.

문제가 생겼을 때 무슨 일이 있었는지 알아내는 가장 좋은 방법 중 하나는 로그를 확인하는 것입니다.

**Mac**에서는 다음 명령으로 로그를 봅니다.

```shell
cat ~/.ollama/logs/server.log
```

systemd를 사용하는 **Linux**에서는 다음 명령으로 로그를 찾습니다.

```shell
journalctl -u ollama --no-pager --follow --pager-end
```

**컨테이너**에서 Ollama를 실행하면 로그가 컨테이너의 stdout/stderr로 나갑니다.

```shell
docker logs <container-name>
```

(컨테이너 이름은 `docker ps`로 확인합니다.)

터미널에서 `ollama serve`를 직접 실행 중이라면 로그가 그 터미널에 표시됩니다.

**Windows**에서는 로그가 여러 위치에 저장됩니다. `<cmd>+R`을 눌러 다음을 입력하면 탐색기 창에서
확인할 수 있습니다.

- `explorer %LOCALAPPDATA%\Ollama` — 로그 위치. 최근 서버 로그는 `server.log`에, 이전 로그는 `server-#.log`에 있습니다.
- `explorer %LOCALAPPDATA%\Programs\Ollama` — 바이너리 위치(설치 시 사용자 PATH에 추가됩니다)
- `explorer %HOMEPATH%\.ollama` — 모델과 설정 저장 위치
- `explorer %TEMP%` — 임시 실행 파일이 하나 이상의 `ollama*` 디렉터리에 저장되는 위치

추가 디버그 로깅을 켜서 문제를 진단하려면, 먼저 **트레이 메뉴에서 실행 중인 앱을 종료**한 뒤
PowerShell 터미널에서 다음을 실행합니다.

```powershell
$env:OLLAMA_DEBUG="1"
& "ollama app.exe"
```

로그 해석에 도움이 필요하면 [Discord](https://discord.gg/ollama)에 참여하세요.

## LLM 라이브러리

Ollama에는 서로 다른 GPU와 CPU 벡터 기능에 맞게 컴파일된 여러 LLM 라이브러리가 포함되어 있습니다.
Ollama는 시스템 성능에 따라 가장 적합한 것을 자동으로 선택합니다. 자동 감지에 문제가 있거나 다른
문제(예: GPU 크래시)가 발생하면 특정 LLM 라이브러리를 강제로 지정해 우회할 수 있습니다. 성능은
`cpu_avx2`가 가장 좋고, 그다음이 `cpu_avx`이며, 가장 느리지만 호환성이 가장 높은 것은 `cpu`입니다.
macOS의 Rosetta 에뮬레이션에서는 `cpu` 라이브러리가 동작합니다.

서버 로그에서 다음과 비슷한 메시지를 볼 수 있습니다(릴리스마다 다름).

```
Dynamic LLM libraries [rocm_v6 cpu cpu_avx cpu_avx2 cuda_v11 rocm_v5]
```

**실험적 LLM 라이브러리 오버라이드**

`OLLAMA_LLM_LIBRARY`를 사용 가능한 LLM 라이브러리 중 하나로 설정하면 자동 감지를 건너뛸 수
있습니다. 예를 들어 CUDA 카드가 있지만 AVX2 벡터를 지원하는 CPU LLM 라이브러리를 강제로 쓰려면
다음과 같이 합니다.

```shell
OLLAMA_LLM_LIBRARY="cpu_avx2" ollama serve
```

CPU가 지원하는 기능은 다음으로 확인할 수 있습니다.

```shell
cat /proc/cpuinfo| grep flags | head -1
```

## Linux에서 이전·사전 릴리스 버전 설치

Linux에서 문제가 생겨 이전 버전을 설치하고 싶거나, 정식 릴리스 전에 사전 릴리스를 시험해 보고 싶다면
설치 스크립트에 설치할 버전을 지정할 수 있습니다.

```shell
curl -fsSL https://ollama.com/install.sh | OLLAMA_VERSION=0.5.7 sh
```

## Linux tmp noexec

Ollama가 임시 실행 파일을 저장하는 위치에 "noexec" 플래그가 설정되어 있다면, `OLLAMA_TMPDIR`을
Ollama 실행 사용자가 쓸 수 있는 위치로 지정해 대체 경로를 쓸 수 있습니다. 예: `OLLAMA_TMPDIR=/usr/share/ollama/`

## Linux docker

도커 컨테이너에서 Ollama가 처음에는 GPU로 동작하다가 일정 시간이 지난 뒤 CPU로 전환되고 서버
로그에 GPU 발견 실패가 기록된다면, 도커의 systemd cgroup 관리를 비활성화해 해결할 수 있습니다.
호스트의 `/etc/docker/daemon.json`을 편집해 도커 설정에 `"exec-opts": ["native.cgroupdriver=cgroupfs"]`를
추가하세요.

## NVIDIA GPU 발견

Ollama는 시작할 때 시스템에 있는 GPU를 조사해 호환성과 사용 가능한 VRAM 양을 판단합니다. 때때로
이 발견 과정이 GPU를 찾지 못할 수 있습니다. 일반적으로 최신 드라이버를 사용하면 가장 좋은 결과를
얻습니다.

### Linux NVIDIA 문제 해결

컨테이너로 Ollama를 실행한다면, [docker](./docker) 문서에 설명된 대로 컨테이너 런타임을 먼저
설정했는지 확인하세요.

Ollama가 GPU를 초기화하는 데 어려움을 겪을 때가 있습니다. 서버 로그를 확인하면 "3"(미초기화),
"46"(장치 사용 불가), "100"(장치 없음), "999"(알 수 없음) 등 다양한 오류 코드로 나타날 수 있습니다.
다음 기법들이 문제 해결에 도움이 될 수 있습니다.

- 컨테이너를 쓴다면 컨테이너 런타임이 동작하나요? `docker run --gpus all ubuntu nvidia-smi`를 시도해 보세요. 이게 안 되면 Ollama도 NVIDIA GPU를 볼 수 없습니다.
- uvm 드라이버가 로드되어 있나요? `sudo nvidia-modprobe -u`
- `nvidia_uvm` 드라이버를 다시 로드해 보세요 — `sudo rmmod nvidia_uvm` 후 `sudo modprobe nvidia_uvm`
- 재부팅해 보세요.
- 최신 nvidia 드라이버를 사용하고 있는지 확인하세요.

이 중 어느 것으로도 해결되지 않으면 추가 정보를 모아 이슈를 등록하세요.

- `CUDA_ERROR_LEVEL=50`을 설정하고 다시 시도해 더 자세한 진단 로그를 얻으세요.
- dmesg에서 오류를 확인하세요 — `sudo dmesg | grep -i nvrm` 및 `sudo dmesg | grep -i nvidia`

## AMD GPU 발견

Linux에서 AMD GPU에 접근하려면 보통 `/dev/kfd` 장치에 접근하기 위해 `video` 또는 `render` 그룹
멤버십이 필요합니다. 권한이 올바르게 설정되어 있지 않으면 Ollama가 이를 감지해 서버 로그에 오류를
보고합니다.

컨테이너에서 실행할 때, 일부 Linux 배포판과 컨테이너 런타임에서는 ollama 프로세스가 GPU에 접근하지
못할 수 있습니다. 호스트에서 `ls -lnd /dev/kfd /dev/dri /dev/dri/*`로 시스템의 **숫자형** 그룹 ID를
확인하고, 컨테이너가 필요한 장치에 접근할 수 있도록 `--group-add ...` 인자를 추가로 전달하세요.
예를 들어 다음 출력에서 `crw-rw---- 1 0  44 226,   0 Sep 16 16:55 /dev/dri/card0`의 그룹 ID 열은 `44`입니다.

Ollama가 GPU를 올바르게 발견하거나 추론에 사용하지 못하는 문제가 있다면, 다음이 원인 파악에 도움이
될 수 있습니다.

- `AMD_LOG_LEVEL=3` — AMD HIP/ROCm 라이브러리의 정보 로그 레벨을 켭니다. 더 자세한 오류 코드를 보여 줘 문제 해결에 도움이 됩니다.
- `OLLAMA_DEBUG=1` — GPU 발견 과정에서 추가 정보가 보고됩니다.
- dmesg에서 amdgpu·kfd 드라이버 오류를 확인하세요 — `sudo dmesg | grep -i amdgpu` 및 `sudo dmesg | grep -i kfd`

### AMD 드라이버 버전 불일치

Linux에서 AMD GPU가 감지되지 않고 서버 로그에 다음과 같은 메시지가 있다면,

```
msg="failure during GPU discovery" ... error="failed to finish discovery before timeout"
msg="bootstrap discovery took" duration=30s ...
```

보통 시스템의 AMD GPU 드라이버가 너무 오래되었다는 뜻입니다. Ollama는 ROCm 7 Linux
라이브러리를 번들하는데, 이는 호환되는 ROCm 7 커널 드라이버를 요구합니다. 시스템이 더 오래된
드라이버(ROCm 6.x 이하)를 사용 중이면 장치 발견 단계에서 GPU 초기화가 멈춰 결국 타임아웃되고,
Ollama는 CPU로 폴백합니다.

해결하려면 [AMD의 ROCm 문서](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/)에 있는
`amdgpu-install` 유틸리티로 ROCm v7 드라이버로 업그레이드하세요. 업그레이드 후 재부팅하고 Ollama를
다시 시작하세요.

## 여러 AMD GPU

Linux에서 여러 AMD GPU에 걸쳐 모델이 로드될 때 의미 없는 응답이 나온다면, 다음 가이드를 참고하세요.

- [https://rocm.docs.amd.com/projects/radeon/en/latest/docs/install/native_linux/mgpu.html#mgpu-known-issues-and-limitations](https://rocm.docs.amd.com/projects/radeon/en/latest/docs/install/native_linux/mgpu.html#mgpu-known-issues-and-limitations)

## Windows 터미널 오류

오래된 버전의 Windows 10(예: 21H1)에는 표준 터미널 프로그램이 제어 문자를 올바르게 표시하지 못하는
버그가 있는 것으로 알려져 있습니다. 이로 인해 `←[?25h←[?25l` 같은 긴 문자열이 표시되고 때로는
`The parameter is incorrect` 오류가 발생할 수 있습니다. 이 문제를 해결하려면 Win 10 22H1 이상으로
업데이트하세요.

> 원문: https://docs.ollama.com/troubleshooting

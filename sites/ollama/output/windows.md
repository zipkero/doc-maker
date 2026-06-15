# Windows

Ollama는 NVIDIA·AMD Radeon GPU를 지원하는 네이티브 Windows 애플리케이션으로 동작합니다.
Windows용 Ollama를 설치하면 Ollama가 백그라운드에서 실행되고, `cmd`, `powershell` 또는 원하는
터미널에서 `ollama` 명령줄을 사용할 수 있습니다. 평소처럼 Ollama [API](/api)는
`http://localhost:11434`에서 제공됩니다.

## 시스템 요구 사항

- Windows 10 22H2 이상, Home 또는 Pro
- NVIDIA 카드를 사용한다면 NVIDIA 452.39 이상 드라이버
- ROCm 가속을 위한 AMD ROCm v7 / HIP7 지원 드라이버 스택, 또는 Vulkan 가속을 위한 Vulkan 지원 AMD Radeon 드라이버

Ollama는 진행 상황 표시에 유니코드 문자를 사용하는데, Windows 10의 일부 오래된 터미널 폰트에서는
알 수 없는 사각형으로 보일 수 있습니다. 이런 경우 터미널 폰트 설정을 바꿔 보세요.

> 참고: RX 6800급 카드를 포함한 일부 RDNA2 / Radeon RX 6000 시스템은 현재 Windows AMD
> 드라이버에서 ROCm v7을 노출하지 않을 수 있습니다. 이런 시스템에서는 Vulkan이 기본으로
> 활성화되며 권장 폴백입니다. iGPU/dGPU 혼합 시스템이 불안정한 Vulkan iGPU를 선택한다면,
> `GGML_VK_VISIBLE_DEVICES`를 외장 GPU 인덱스로 설정하세요.

## 파일 시스템 요구 사항

Ollama 설치에는 관리자 권한이 필요하지 않으며, 기본적으로 홈 디렉터리에 설치됩니다. 바이너리 설치에는
최소 4GB의 공간이 필요합니다. 설치 후에는 대형 언어 모델을 저장할 추가 공간이 필요한데, 모델 크기는
수십에서 수백 GB에 이를 수 있습니다. 홈 디렉터리에 공간이 부족하다면 바이너리 설치 위치와 모델 저장
위치를 변경할 수 있습니다.

### 설치 위치 변경

Ollama 애플리케이션을 홈 디렉터리가 아닌 다른 위치에 설치하려면, 다음 플래그와 함께 설치 프로그램을
시작합니다.

```powershell
OllamaSetup.exe /DIR="d:\some\location"
```

### 모델 저장 위치 변경

내려받은 모델을 홈 디렉터리 대신 다른 곳에 저장하려면, 사용자 계정에 환경변수 `OLLAMA_MODELS`를
설정합니다.

1. 설정(Windows 11) 또는 제어판(Windows 10)을 열고 *environment variables*(환경 변수)를 검색합니다.
2. *Edit environment variables for your account*(내 계정의 환경 변수 편집)를 클릭합니다.
3. 모델을 저장할 위치로 `OLLAMA_MODELS` 변수를 편집하거나 새로 만듭니다.
4. OK/적용을 눌러 저장합니다.

Ollama가 이미 실행 중이라면, 트레이 애플리케이션을 종료한 뒤 시작 메뉴에서 다시 실행하거나,
환경변수를 저장한 후 새로 연 터미널에서 실행합니다.

## API 접근

다음은 `powershell`에서 API에 접근하는 간단한 예시입니다.

```powershell
(Invoke-WebRequest -method POST -Body '{"model":"llama3.2", "prompt":"Why is the sky blue?", "stream": false}' -uri http://localhost:11434/api/generate ).Content | ConvertFrom-json
```

## 문제 해결

Windows의 Ollama는 파일을 여러 위치에 저장합니다. `<Ctrl>+R`을 눌러 다음을 입력하면 탐색기 창에서
확인할 수 있습니다.

- `explorer %LOCALAPPDATA%\Ollama` — 로그와 내려받은 업데이트
  - `app.log` — GUI 애플리케이션의 최근 로그
  - `server.log` — 최근 서버 로그
  - `upgrade.log` — 업그레이드 로그 출력
- `explorer %LOCALAPPDATA%\Programs\Ollama` — 바이너리 위치(설치 시 사용자 PATH에 추가됩니다)
- `explorer %HOMEPATH%\.ollama` — 모델과 설정
- `explorer %TEMP%` — 하나 이상의 `ollama*` 디렉터리에 있는 임시 실행 파일

## 제거

Ollama Windows 설치 프로그램은 제거 프로그램을 등록합니다. Windows 설정의 `앱 및 기능`(Add or
remove programs)에서 Ollama를 제거할 수 있습니다.

> 참고: [모델 저장 위치를 변경](#모델-저장-위치-변경)했다면, 설치 제거 프로그램은 내려받은 모델을
> 삭제하지 않습니다.

## 독립형(Standalone) CLI

Windows에서 Ollama를 설치하는 가장 쉬운 방법은 `OllamaSetup.exe` 설치 프로그램을 사용하는
것입니다. 관리자 권한 없이 사용자 계정에 설치됩니다. Ollama는 최신 모델을 지원하도록 정기적으로
업데이트되며, 이 설치 프로그램이 최신 상태를 유지하도록 도와줍니다.

Ollama를 서비스로 설치하거나 통합하려면, Ollama CLI와 Nvidia용 GPU 라이브러리 의존성만 담긴
독립형 `ollama-windows-amd64.zip` 파일을 사용할 수 있습니다. 하드웨어에 따라 추가 패키지를 같은
디렉터리에 내려받아 압축을 풀어야 할 수도 있습니다.

- **AMD GPU**: `ollama-windows-amd64-rocm.zip`
- **MLX (CUDA)**: `ollama-windows-amd64-mlx.zip`

이를 통해 기존 애플리케이션에 Ollama를 임베딩하거나, [NSSM](https://nssm.cc/) 같은 도구로
`ollama serve`를 시스템 서비스로 실행할 수 있습니다.

> 참고: 이전 버전에서 업그레이드하는 경우, 먼저 오래된 디렉터리를 제거해야 합니다.

> 원문: https://docs.ollama.com/windows

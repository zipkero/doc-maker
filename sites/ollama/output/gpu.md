# 하드웨어 지원

## Nvidia

Ollama는 컴퓨트 능력(compute capability) 5.0 이상, 드라이버 버전 531 이상인 Nvidia GPU를 지원합니다.
컴퓨트 능력 5.0부터 6.2까지인 Nvidia GPU는 드라이버 버전 570 이상이 필요합니다.

내 카드가 지원되는지는 컴퓨트 호환성에서 확인하세요:
[https://developer.nvidia.com/cuda-gpus](https://developer.nvidia.com/cuda-gpus)

| 컴퓨트 능력 | 제품군              | 카드                                                                                                                          |
| ---------- | ------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| 12.1       | NVIDIA              | `GB10 (DGX Spark)`                                                                                                            |
| 12.0       | GeForce RTX 50xx    | `RTX 5060` `RTX 5060 Ti` `RTX 5070` `RTX 5070 Ti` `RTX 5080` `RTX 5090`                                                       |
|            | NVIDIA Professional | `RTX PRO 4000 Blackwell` `RTX PRO 4500 Blackwell` `RTX PRO 5000 Blackwell` `RTX PRO 6000 Blackwell`                           |
| 9.0        | NVIDIA              | `H200` `H100`                                                                                                                 |
| 8.9        | GeForce RTX 40xx    | `RTX 4090` `RTX 4080 SUPER` `RTX 4080` `RTX 4070 Ti SUPER` `RTX 4070 Ti` `RTX 4070 SUPER` `RTX 4070` `RTX 4060 Ti` `RTX 4060` |
|            | NVIDIA Professional | `L4` `L40` `RTX 6000`                                                                                                         |
| 8.6        | GeForce RTX 30xx    | `RTX 3090 Ti` `RTX 3090` `RTX 3080 Ti` `RTX 3080` `RTX 3070 Ti` `RTX 3070` `RTX 3060 Ti` `RTX 3060` `RTX 3050 Ti` `RTX 3050`  |
|            | NVIDIA Professional | `A40` `RTX A6000` `RTX A5000` `RTX A4000` `RTX A3000` `RTX A2000` `A10` `A16` `A2`                                            |
| 8.0        | NVIDIA              | `A100` `A30`                                                                                                                  |
| 7.5        | GeForce GTX/RTX     | `GTX 1650 Ti` `TITAN RTX` `RTX 2080 Ti` `RTX 2080` `RTX 2070` `RTX 2060`                                                      |
|            | NVIDIA Professional | `T4` `RTX 5000` `RTX 4000` `RTX 3000` `T2000` `T1200` `T1000` `T600` `T500`                                                   |
|            | Quadro              | `RTX 8000` `RTX 6000` `RTX 5000` `RTX 4000`                                                                                   |
| 7.0        | NVIDIA              | `TITAN V` `V100` `Quadro GV100`                                                                                               |
| 6.1        | NVIDIA TITAN        | `TITAN Xp` `TITAN X`                                                                                                          |
|            | GeForce GTX         | `GTX 1080 Ti` `GTX 1080` `GTX 1070 Ti` `GTX 1070` `GTX 1060` `GTX 1050 Ti` `GTX 1050`                                         |
|            | Quadro              | `P6000` `P5200` `P4200` `P3200` `P5000` `P4000` `P3000` `P2200` `P2000` `P1000` `P620` `P600` `P500` `P520`                   |
|            | Tesla               | `P40` `P4`                                                                                                                    |
| 6.0        | NVIDIA              | `Tesla P100` `Quadro GP100`                                                                                                   |
| 5.2        | GeForce GTX         | `GTX TITAN X` `GTX 980 Ti` `GTX 980` `GTX 970` `GTX 960` `GTX 950`                                                            |
|            | Quadro              | `M6000 24GB` `M6000` `M5000` `M5500M` `M4000` `M2200` `M2000` `M620`                                                          |
|            | Tesla               | `M60` `M40`                                                                                                                   |
| 5.0        | GeForce GTX         | `GTX 750 Ti` `GTX 750` `NVS 810`                                                                                              |
|            | Quadro              | `K2200` `K1200` `K620` `M1200` `M520` `M5000M` `M4000M` `M3000M` `M2000M` `M1000M` `K620M` `M600M` `M500M`                    |

더 오래된 GPU를 지원하도록 로컬에서 빌드하는 방법은 [개발자 문서](./development#linux-cuda-nvidia)를 참고하세요.

### GPU 선택

시스템에 NVIDIA GPU가 여러 개 있고 Ollama가 일부만 사용하도록 제한하려면 `CUDA_VISIBLE_DEVICES`에 GPU
목록을 쉼표로 구분해 설정하면 됩니다. 숫자 ID를 쓸 수도 있지만 순서가 달라질 수 있어 UUID가 더 안정적입니다.
GPU의 UUID는 `nvidia-smi -L`로 확인할 수 있습니다. GPU를 무시하고 CPU 사용을 강제하려면 잘못된 GPU
ID(예: "-1")를 지정하세요.

### Linux의 절전·재개(Suspend/Resume)

Linux에서는 절전·재개 사이클 이후 Ollama가 NVIDIA GPU를 찾지 못하고 CPU로 폴백(fallback)하는 경우가
있습니다. 이 드라이버 버그는 NVIDIA UVM 드라이버를 다시 로드해 우회할 수 있습니다:
`sudo rmmod nvidia_uvm && sudo modprobe nvidia_uvm`

## AMD Radeon

Ollama는 ROCm 라이브러리를 통해 다음 AMD GPU를 지원합니다.

> 참고: 추가적인 AMD GPU 지원은 Vulkan 라이브러리를 통해 제공됩니다 — 아래를 참고하세요.

### Linux 지원

Linux에서 Ollama는 AMD ROCm v7 드라이버가 필요합니다.
[AMD ROCm 문서](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/)의 `amdgpu-install`
유틸리티로 설치하거나 업그레이드할 수 있습니다.

| 제품군            | 카드 및 가속기                                                                                                                                                                                                          |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| AMD Radeon RX     | `9070 XT` `9070 GRE` `9070` `9060 XT` `9060 XT LP` `9060` `7900 XTX` `7900 XT` `7900 GRE` `7800 XT` `7700 XT` `7700` `7600 XT` `7600` `6950 XT` `6900 XTX` `6900XT` `6800 XT` `6800` `5700 XT` `5700` `5600 XT` `5500 XT` |
| AMD Radeon AI PRO | `R9700` `R9600D`                                                                                                                                                                                                          |
| AMD Radeon PRO    | `W7900` `W7800` `W7700` `W7600` `W7500` `W6900X` `W6800X Duo` `W6800X` `W6800` `V620`                                                                                                                                     |
| AMD Ryzen AI      | `Ryzen AI Max+ 395` `Ryzen AI Max 390` `Ryzen AI Max 385` `Ryzen AI 9 HX 475` `Ryzen AI 9 HX 470` `Ryzen AI 9 465` `Ryzen AI 9 HX 375` `Ryzen AI 9 HX 370` `Ryzen AI 9 365`                                               |
| AMD Instinct      | `MI350X` `MI300X` `MI300A` `MI250X` `MI250` `MI210` `MI100`                                                                                                                                                               |

### Windows 지원

Windows에서 Ollama는 AMD ROCm v7 / HIP7 지원 드라이버 스택이 필요합니다.

| 제품군         | 카드 및 가속기                                                                                                     |
| -------------- | ------------------------------------------------------------------------------------------------------------------- |
| AMD Radeon RX  | `7900 XTX` `7900 XT` `7900 GRE` `7800 XT` `7700 XT` `7600 XT` `7600` `6950 XT` `6900 XTX` `6900XT` `6800 XT` `6800` |
| AMD Radeon PRO | `W7900` `W7800` `W7700` `W7600` `W7500` `W6900X` `W6800X Duo` `W6800X` `W6800` `V620`                               |

### Linux에서의 재정의(Override)

Ollama는 AMD ROCm 라이브러리를 활용하는데, 이 라이브러리는 모든 AMD GPU를 지원하지는 않습니다. 일부
경우에는 시스템이 비슷한 LLVM 타깃을 사용하도록 강제할 수 있습니다. 예를 들어 Radeon RX 5400은
`gfx1034`(10.3.4)인데 ROCm은 현재 이 타깃을 지원하지 않습니다. 가장 가까운 지원 타깃은 `gfx1030`입니다.
`x.y.z` 형식으로 `HSA_OVERRIDE_GFX_VERSION` 환경 변수를 사용할 수 있습니다. 예를 들어 시스템이 RX 5400에서
실행되도록 강제하려면 서버에 `HSA_OVERRIDE_GFX_VERSION="10.3.0"`을 환경 변수로 설정합니다. 지원되지 않는
AMD GPU가 있다면 아래 지원 타깃 목록을 보고 실험해 볼 수 있습니다.

GFX 버전이 서로 다른 GPU가 여러 개 있으면, 환경 변수에 숫자 장치 번호를 덧붙여 개별 설정할 수 있습니다.
예를 들어 `HSA_OVERRIDE_GFX_VERSION_0=10.3.0`, `HSA_OVERRIDE_GFX_VERSION_1=11.0.0`처럼 지정합니다.

현재 Linux에서 알려진 지원 GPU 유형은 다음 LLVM 타깃들입니다. 아래 표는 각 LLVM 타깃에 매핑되는 GPU
예시를 보여 줍니다.

| **LLVM 타깃**   | **GPU 예시**                  |
| --------------- | ----------------------------- |
| gfx908          | Radeon Instinct MI100         |
| gfx90a          | Radeon Instinct MI210/MI250   |
| gfx942          | Radeon Instinct MI300X/MI300A |
| gfx950          | Radeon Instinct MI350X        |
| gfx1010         | Radeon RX 5700 XT             |
| gfx1012         | Radeon RX 5500 XT             |
| gfx1030         | Radeon PRO V620               |
| gfx1100         | Radeon PRO W7900              |
| gfx1101         | Radeon PRO W7700              |
| gfx1102         | Radeon RX 7600                |
| gfx1103         | Radeon 780M                   |
| gfx1150         | Ryzen AI 9 HX 375             |
| gfx1151         | Ryzen AI Max+ 395             |
| gfx1200         | Radeon RX 9070                |
| gfx1201         | Radeon RX 9070 XT             |

추가 도움이 필요하면 [Discord](https://discord.gg/ollama)로 문의하거나
[이슈](https://github.com/ollama/ollama/issues)를 등록하세요.

### GPU 선택

시스템에 AMD GPU가 여러 개 있고 Ollama가 일부만 사용하도록 제한하려면 `ROCR_VISIBLE_DEVICES`에 GPU
목록을 쉼표로 구분해 설정하면 됩니다. 장치 목록은 `rocminfo`로 확인할 수 있습니다. GPU를 무시하고 CPU
사용을 강제하려면 잘못된 GPU ID(예: "-1")를 지정하세요. 가능하면 숫자 값 대신 `Uuid`를 사용해 장치를
고유하게 식별하세요.

### 컨테이너 권한

일부 Linux 배포판에서는 SELinux가 컨테이너의 AMD GPU 장치 접근을 막을 수 있습니다. 호스트 시스템에서
`sudo setsebool container_use_devices=1`을 실행하면 컨테이너가 장치를 사용하도록 허용할 수 있습니다.

## Metal (Apple GPU)

Ollama는 Metal API를 통해 Apple 기기에서 GPU 가속을 지원합니다.

## Vulkan GPU 지원

Windows와 Linux에서는 [Vulkan](https://www.vulkan.org/)을 통해 추가 GPU 지원이 제공됩니다. 백엔드가
설치되면 Vulkan이 기본적으로 활성화됩니다. Windows에서는 대부분의 GPU 벤더 드라이버에 Vulkan 지원이
함께 포함되어 별도 설정이 필요 없습니다. 반면 대부분의 Linux 배포판은 추가 구성 요소를 설치해야 하며,
Mesa와 GPU 벤더 전용 패키지 중에서 Vulkan 드라이버를 선택할 수 있습니다.

- Linux Intel GPU 안내 — [https://dgpu-docs.intel.com/driver/client/overview.html](https://dgpu-docs.intel.com/driver/client/overview.html)
- Linux AMD GPU 안내 — [https://amdgpu-install.readthedocs.io/en/latest/install-script.html#specifying-a-vulkan-implementation](https://amdgpu-install.readthedocs.io/en/latest/install-script.html#specifying-a-vulkan-implementation)

일부 Linux 배포판에서는 AMD GPU를 사용하려면 `ollama` 사용자를 `render` 그룹에 추가해야 할 수 있습니다.

Ollama 스케줄러는 GPU 라이브러리가 보고하는 가용 VRAM 데이터를 활용해 최적의 스케줄링을 결정합니다.
Vulkan에서 이 가용 VRAM 데이터를 노출하려면 추가 권한이 있거나 root로 실행해야 합니다. root 접근 권한도
이 권한도 없으면 Ollama는 모델의 대략적인 크기를 사용해 최선의 스케줄링을 시도합니다.

```bash
sudo setcap cap_perfmon+ep /usr/local/bin/ollama
```

### GPU 선택

특정 Vulkan GPU를 선택하려면 Ollama 서버에서 `GGML_VK_VISIBLE_DEVICES` 환경 변수에 하나 이상의 숫자
ID를 설정하면 됩니다. 설정 방법은 [FAQ](faq#ollama-서버는-어떻게-설정하나요)를 참고하세요. Vulkan 기반
GPU에서 문제가 발생하면 `OLLAMA_VULKAN=0` 또는 `GGML_VK_VISIBLE_DEVICES=-1`로 설정해 모든 Vulkan GPU를
비활성화할 수 있습니다.

iGPU와 dGPU가 섞인 시스템에서 Vulkan iGPU가 불안정하다면, Vulkan은 켜 둔 채
`GGML_VK_VISIBLE_DEVICES`를 외장 GPU 인덱스로 설정하세요. 예를 들어 `Vulkan1`이 외장 GPU라면
`GGML_VK_VISIBLE_DEVICES=1`을 사용합니다.

> 원문: https://docs.ollama.com/gpu

# Docker

Docker로 Ollama를 실행하는 방법을 환경별로 정리했습니다.

## CPU 전용

GPU 없이 CPU만으로 실행하려면 다음 명령을 사용합니다.

```shell
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

## NVIDIA GPU

먼저 [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html#installation)을 설치합니다.

### Apt로 설치

1. 저장소를 구성합니다.

   ```shell
   curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey \
       | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
   curl -fsSL https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list \
       | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' \
       | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
   sudo apt-get update
   ```

2. NVIDIA Container Toolkit 패키지를 설치합니다.

   ```shell
   sudo apt-get install -y nvidia-container-toolkit
   ```

### Yum 또는 Dnf로 설치

1. 저장소를 구성합니다.

   ```shell
   curl -fsSL https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo \
       | sudo tee /etc/yum.repos.d/nvidia-container-toolkit.repo
   ```

2. NVIDIA Container Toolkit 패키지를 설치합니다.

   ```shell
   sudo yum install -y nvidia-container-toolkit
   ```

### Docker가 NVIDIA 드라이버를 사용하도록 설정

```shell
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

### 컨테이너 시작

```shell
docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

> 참고: NVIDIA JetPack 시스템에서 실행하는 경우, Ollama가 올바른 JetPack 버전을 자동으로 감지하지 못합니다.
> 컨테이너에 `JETSON_JETPACK=5` 또는 `JETSON_JETPACK=6` 환경 변수를 전달해 버전 5 또는 6을 선택하세요.

## AMD GPU

AMD GPU에서 Docker로 Ollama를 실행하려면 `rocm` 태그와 다음 명령을 사용합니다.

```shell
docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama:rocm
```

## Vulkan 지원

Vulkan은 `ollama/ollama` 이미지에 포함되어 있으며, 컨테이너가 GPU 장치에 접근할 수 있을 때 기본으로 활성화됩니다.

```shell
docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

Vulkan을 비활성화하려면 `OLLAMA_VULKAN=0`을, 특정 Vulkan 장치를 선택하려면 `GGML_VK_VISIBLE_DEVICES=<ids>`를 사용합니다.

## 로컬에서 모델 실행

이제 모델을 실행할 수 있습니다.

```shell
docker exec -it ollama ollama run llama3.2
```

## 다른 모델 사용해 보기

더 많은 모델은 [Ollama 라이브러리](https://ollama.com/library)에서 찾을 수 있습니다.

> 원문: https://docs.ollama.com/docker

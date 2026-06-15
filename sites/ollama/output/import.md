# 모델 가져오기

## 목차

- [Safetensors 어댑터 가져오기](#safetensors-가중치로-파인튜닝한-어댑터-가져오기)
- [Safetensors 모델 가져오기](#safetensors-가중치로-모델-가져오기)
- [GGUF 파일 가져오기](#gguf-기반-모델-또는-어댑터-가져오기)
- [ollama.com에 모델 공유하기](#ollamacom에-모델-공유하기)

## Safetensors 가중치로 파인튜닝한 어댑터 가져오기

먼저 `Modelfile`을 만들고, 파인튜닝에 사용한 기반 모델을 가리키는 `FROM` 명령과 Safetensors 어댑터가 있는
디렉터리를 가리키는 `ADAPTER` 명령을 작성합니다.

```dockerfile
FROM <base model name>
ADAPTER /path/to/safetensors/adapter/directory
```

`FROM` 명령에는 반드시 어댑터를 만들 때 사용한 것과 동일한 기반 모델을 지정해야 합니다. 그렇지 않으면
결과가 불안정해집니다. 대부분의 프레임워크가 서로 다른 양자화 방식을 사용하므로, 양자화되지 않은(즉
QLoRA가 아닌) 어댑터를 사용하는 것이 가장 좋습니다. 어댑터가 `Modelfile`과 같은 디렉터리에 있다면
`ADAPTER .`로 어댑터 경로를 지정할 수 있습니다.

이제 `Modelfile`을 만든 디렉터리에서 `ollama create`를 실행합니다.

```shell
ollama create my-model
```

마지막으로 모델을 테스트합니다.

```shell
ollama run my-model
```

Ollama가 어댑터 가져오기를 지원하는 모델 아키텍처는 다음과 같습니다.

- Llama (Llama 2, Llama 3, Llama 3.1, Llama 3.2 포함)
- Mistral (Mistral 1, Mistral 2, Mixtral 포함)
- Gemma (Gemma 1, Gemma 2 포함)

어댑터는 Safetensors 형식으로 출력할 수 있는 파인튜닝 프레임워크나 도구로 만들 수 있습니다. 예를 들면
다음과 같습니다.

- Hugging Face [파인튜닝 프레임워크](https://huggingface.co/docs/transformers/en/training)
- [Unsloth](https://github.com/unslothai/unsloth)
- [MLX](https://github.com/ml-explore/mlx)

## Safetensors 가중치로 모델 가져오기

먼저 `Modelfile`을 만들고, Safetensors 가중치가 들어 있는 디렉터리를 가리키는 `FROM` 명령을 작성합니다.

```dockerfile
FROM /path/to/safetensors/directory
```

가중치와 같은 디렉터리에 Modelfile을 만든다면 `FROM .` 명령을 사용할 수 있습니다.

이제 `Modelfile`을 만든 디렉터리에서 `ollama create` 명령을 실행합니다.

```shell
ollama create my-model
```

마지막으로 모델을 테스트합니다.

```shell
ollama run my-model
```

Ollama가 모델 가져오기를 지원하는 아키텍처는 다음과 같습니다.

- Llama (Llama 2, Llama 3, Llama 3.1, Llama 3.2 포함)
- Mistral (Mistral 1, Mistral 2, Mixtral 포함)
- Gemma (Gemma 1, Gemma 2 포함)
- Phi3

여기에는 파운데이션 모델은 물론, 파운데이션 모델과 *융합(fused)*된 파인튜닝 모델도 포함됩니다.

## GGUF 기반 모델 또는 어댑터 가져오기

GGUF 기반 모델이나 어댑터가 있다면 Ollama로 가져올 수 있습니다. GGUF 모델이나 어댑터는 다음 방법으로
얻을 수 있습니다.

- Llama.cpp의 `convert_hf_to_gguf.py`로 Safetensors 모델 변환
- Llama.cpp의 `convert_lora_to_gguf.py`로 Safetensors 어댑터 변환
- HuggingFace 같은 곳에서 모델이나 어댑터 다운로드

GGUF 모델을 가져오려면 다음 내용을 담은 `Modelfile`을 만듭니다.

```dockerfile
FROM /path/to/file.gguf
```

GGUF 어댑터라면 다음과 같이 `Modelfile`을 만듭니다.

```dockerfile
FROM <model name>
ADAPTER /path/to/file.gguf
```

GGUF 어댑터를 가져올 때는 그 어댑터를 만들 때 사용한 것과 동일한 기반 모델을 사용하는 것이 중요합니다.
기반 모델로는 다음을 사용할 수 있습니다.

- Ollama의 모델
- GGUF 파일
- Safetensors 기반 모델

`Modelfile`을 만든 뒤 `ollama create` 명령으로 모델을 빌드합니다.

```shell
ollama create my-model
```

## 모델 양자화하기

모델을 양자화하면 더 빠르게, 더 적은 메모리로 실행할 수 있지만 정확도는 낮아집니다. 덕분에 더 평범한
하드웨어에서도 모델을 실행할 수 있습니다.

Ollama는 `ollama create` 명령에 `-q/--quantize` 플래그를 사용해 FP16 및 FP32 기반 모델을 다양한 양자화
수준으로 변환할 수 있습니다.

먼저 양자화하려는 FP16 또는 FP32 기반 모델로 Modelfile을 만듭니다.

```dockerfile
FROM /path/to/my/gemma/f16/model
```

그런 다음 `ollama create`로 양자화된 모델을 만듭니다.

```shell
$ ollama create --quantize q4_K_M mymodel
transferring model data
quantizing F16 model to Q4_K_M
creating new layer sha256:735e246cc1abfd06e9cdcf95504d6789a6cd1ad7577108a70d9902fef503c1bd
creating new layer sha256:0853f0ad24e5865173bbf9ffcc7b0f5d56b66fd690ab1009867e45e7d2c4db0f
writing manifest
success
```

### 지원되는 양자화

- `q8_0`

#### K-means 양자화

- `q4_K_S`
- `q4_K_M`

## ollama.com에 모델 공유하기

직접 만든 모델은 [ollama.com](https://ollama.com)에 push해 다른 사용자가 사용해 볼 수 있도록 공유할 수
있습니다.

먼저 브라우저로 [Ollama 가입](https://ollama.com/signup) 페이지에 접속합니다. 이미 계정이 있다면 이
단계는 건너뛰어도 됩니다.

`Username` 필드는 모델 이름의 일부로 사용되므로(예: `jmorganca/mymodel`), 선택한 사용자 이름이 마음에
드는지 확인하세요.

계정을 만들고 로그인했다면 [Ollama Keys 설정](https://ollama.com/settings/keys) 페이지로 이동합니다.

페이지의 안내에 따라 Ollama 공개 키의 위치를 확인합니다.

`Add Ollama Public Key` 버튼을 클릭하고, Ollama 공개 키의 내용을 복사해 텍스트 필드에 붙여 넣습니다.

[ollama.com](https://ollama.com)에 모델을 push하려면 먼저 모델 이름이 사용자 이름을 포함해 올바르게
지정되어 있어야 합니다. 올바른 이름을 붙이려면 `ollama cp` 명령으로 모델을 복사해야 할 수도 있습니다.
모델 이름이 마음에 들면 `ollama push` 명령으로 [ollama.com](https://ollama.com)에 push합니다.

```shell
ollama cp mymodel myuser/mymodel
ollama push myuser/mymodel
```

모델을 push하고 나면 다른 사용자가 다음 명령으로 모델을 pull해 실행할 수 있습니다.

```shell
ollama run myuser/mymodel
```

> 원문: https://docs.ollama.com/import

# VS Code

VS Code는 GitHub Copilot Chat을 통해 AI 채팅 기능을 기본 제공합니다. Copilot Chat의 모델 선택기에서
Ollama 모델을 바로 사용할 수 있습니다.

## 사전 준비

* Ollama v0.18.3 이상
* [VS Code 1.113 이상](https://code.visualstudio.com/download)
* [GitHub Copilot Chat 확장 0.41.0 이상](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat)

> 참고: VS Code는 모델 선택기를 사용하려면 커스텀 모델을 쓰는 경우라도 로그인이 필요합니다. 유료
> GitHub Copilot 계정이 필요하지는 않으며, GitHub Copilot Free만으로도 커스텀 모델 선택 기능이
> 활성화됩니다.

## 빠른 설정

```shell
ollama launch vscode
```

명령을 실행하면 추천 모델이 표시됩니다. 최신 모델은 [ollama.com](https://ollama.com/search?c=tools)에서
확인할 수 있습니다.

Ollama 모델을 사용하려면 Copilot Chat 패널 하단에서 **Local**이 선택되어 있는지 확인하세요.

## 특정 모델로 바로 실행

`--model` 옵션으로 사용할 모델을 지정할 수 있습니다.

```shell
ollama launch vscode --model qwen3.5:cloud
```

클라우드 모델은 [ollama.com](https://ollama.com/search?c=cloud)에서도 확인할 수 있습니다.

## 수동 설정

`ollama launch` 없이 Ollama를 직접 구성하려면 다음 순서를 따릅니다.

1. 오른쪽 위 모서리에서 **Copilot Chat** 사이드바를 엽니다.
2. **설정 기어 아이콘**을 클릭해 Language Models 창을 띄웁니다.
3. **Add Models**를 클릭하고 **Ollama**를 선택하면 모든 Ollama 모델이 VS Code로 로드됩니다.
4. 모델 선택기에서 **Unhide** 버튼을 클릭해 Ollama 모델을 표시합니다.

> 원문: https://docs.ollama.com/integrations/vscode

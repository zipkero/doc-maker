# Codex App

Codex App은 OpenAI의 데스크톱 코딩 에이전트로, macOS와 Windows에서 동작합니다. Ollama는 Codex App이
Ollama의 OpenAI 호환 엔드포인트를 사용하도록 설정하므로, 데스크톱 앱에서 로컬 모델과 Ollama 클라우드
모델을 함께 쓸 수 있습니다.

## 설치

macOS 또는 Windows용 [Codex App](https://developers.openai.com/codex/quickstart/)을 설치합니다.

> 참고: Codex App 지원은 Ollama v0.24.0 이상에서 제공됩니다.

## 빠른 설정

```shell
ollama launch codex-app
```

Codex App이 열리면 평소처럼 작업을 시작하거나 저장소를 엽니다.

## 내장 브라우저

Codex App은 내장 브라우저로 로컬 서버나 사이트를 열 수 있습니다. 페이지에 직접 주석을 달아 변경을
요청할 수도 있습니다.

## 리뷰 모드

리뷰 모드를 사용하면 앱을 벗어나지 않고도 코드 변경 사항을 살펴보고, 코멘트를 남기고, 수정을 반복할
수 있습니다.

### 특정 모델로 바로 실행

```shell
ollama launch codex-app --model kimi-k2.6:cloud
```

모델 이름을 전달해 로컬 모델을 사용할 수도 있습니다.

```shell
ollama launch codex-app --model gemma4:31b
```

`ollama launch codex-app`로 지정한 설정은 유지되므로, 다음에 Codex를 열 때도 해당 모델이 선택된
상태가 됩니다.

### Codex App 복원하기

`ollama launch codex-app`을 실행하기 전에 쓰던 프로필로 Codex App을 되돌리려면 다음을 실행합니다.

```shell
ollama launch codex-app --restore
```

Ollama가 Codex App의 설정과 구성을 복원합니다. Codex App이 열려 있으면 재시작 전에 확인을 묻습니다.

`ollama launch codex`로 관리하는 Codex CLI 프로필은 Codex App 프로필과 분리되어 따로 유지됩니다.

Ollama Launch는 Codex App 구성 파일을 덮어쓰기 전에 `~/.ollama/backup/codex-app/`에 백업을
저장합니다. Windows에서는 `~`가 사용자 프로필 디렉터리로 해석됩니다.

## 문제 해결

설정 후에도 Codex App이 열리지 않으면, Codex를 수동으로 한 번 연 다음 `ollama launch codex-app`을
다시 실행하세요.

Codex App이 이미 실행 중인데 모델이 전환되지 않으면, 안내가 나올 때 Ollama가 앱을 재시작하도록
허용하거나, Codex App을 종료한 뒤 `ollama launch codex-app`을 다시 실행하세요.

> 원문: https://docs.ollama.com/integrations/codex-app

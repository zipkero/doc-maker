# macOS

macOS에서 Ollama를 설치·관리하는 방법을 설명합니다.

## 시스템 요구 사항

- macOS Sonoma(v14) 이상
- Apple M 시리즈(CPU·GPU 지원) 또는 x86(CPU 전용)

## 파일 시스템 요구 사항

권장 설치 방법은 `ollama.dmg`를 마운트한 뒤 Ollama 애플리케이션을 시스템 전역
`Applications` 폴더로 드래그 앤 드롭하는 것입니다. 시작 시 Ollama 앱은 PATH에 `ollama` CLI가
있는지 확인하고, 없으면 `/usr/local/bin`에 링크를 만들 권한을 요청합니다.

설치 후에는 대형 언어 모델을 저장할 추가 공간이 필요합니다. 모델 크기는 수십에서 수백 GB에 이를 수
있습니다. 홈 디렉터리에 공간이 부족하다면 바이너리 설치 위치와 모델 저장 위치를 변경할 수 있습니다.

### 설치 위치 변경

Ollama 애플리케이션을 `Applications`가 아닌 다른 위치에 설치하려면, 원하는 위치에 앱을 두고
CLI(`Ollama.app/Contents/Resources/ollama`) 또는 그 심볼릭 링크가 PATH에서 검색되도록 합니다.
처음 실행할 때 나오는 "Move to Applications?" 요청은 거절하세요.

## 문제 해결

macOS의 Ollama는 여러 위치에 파일을 저장합니다.

- `~/.ollama` — 모델과 설정
- `~/.ollama/logs` — 로그
  - `app.log` — GUI 애플리케이션의 최근 로그
  - `server.log` — 최근 서버 로그
- `<설치 위치>/Ollama.app/Contents/Resources/ollama` — CLI 바이너리

## 제거

Ollama를 시스템에서 완전히 제거하려면 다음 파일과 폴더를 삭제합니다.

```shell
sudo rm -rf /Applications/Ollama.app
sudo rm /usr/local/bin/ollama
rm -rf "~/Library/Application Support/Ollama"
rm -rf "~/Library/Saved Application State/com.electron.ollama.savedState"
rm -rf ~/Library/Caches/com.electron.ollama/
rm -rf ~/Library/Caches/ollama
rm -rf ~/Library/WebKit/com.electron.ollama
rm -rf ~/.ollama
```

> 원문: https://docs.ollama.com/macos

# doc-maker

외부 기술 문서 사이트를 설정 기반으로 수집해, 사람이 읽기 좋은 한국어 문서로 재구성·번역하는 도구.
수집은 Go 도구가 결정적으로 수행하고, 번역(재구성)은 Claude가 수행한다.

## 사전 요구

- Go 1.26+ (win32 우선, 단일 실행 파일로 동작)
- 번역 단계는 Claude Code 세션에서 `/translate` 커맨드로 수행한다(외부 번역 API를 쓰지 않는다).

## 디렉터리 구조

한 사이트의 모든 것이 사이트 폴더 하나에 모인다. 폴더명이 곧 사이트 식별자다.

```
sites/<siteID>/
├─ config.json     # 입력: 수집 설정 (사람이 작성)
├─ glossary.json   # 입력: 용어집 (사람이 작성)
├─ raw/            # 생성: 수집된 원문 (원본 경로 구조 보존)
├─ output/         # 생성: 한국어 번역 문서
└─ manifest.json   # 생성: 증분 상태(원문/번역 해시)
```

`config.json`·`glossary.json`만 저장소에 추적하고, `raw/`·`output/`·`manifest.json`은 재생성 가능한
생성물이라 `.gitignore`로 제외한다.

## 새 사이트 추가하기

코드 수정 없이 사이트 폴더만 추가하면 된다.

1. `sites/<siteID>/` 폴더를 만든다(`<siteID>`가 식별자가 된다, 예: `ollama`).
2. `config.json`을 작성한다(아래 참조).
3. `glossary.json`을 작성한다(비어 있어도 됨: `{}`).
4. 수집·번역을 실행한다.

## config.json

```json
{
  "base_url": "https://docs.ollama.com",
  "source_type": "llms.txt",
  "include_patterns": [],
  "exclude_patterns": []
}
```

| 필드 | 설명 |
|---|---|
| `base_url` | 사이트 기준 주소. 페이지 목록은 `base_url + /llms.txt`에서 읽는다 |
| `source_type` | 페이지 목록 확보 방식. **현재 `llms.txt`만 동작**(아래 참조) |
| `include_patterns` | 수집할 경로 glob 목록. 비우면 전체 포함. 예: `["api/**"]` |
| `exclude_patterns` | 제외할 경로 glob 목록. 비우면 제외 없음. 예: `["api/blog/**"]` |

패턴은 URL 경로(선행 `/` 제외)에 적용한다. `*`는 한 세그먼트, `**`는 임의 깊이를 매칭한다.

### source_type 고르는 법

| 값 | 상태 |
|---|---|
| `llms.txt` | ✅ 동작 — 사이트가 `/llms.txt`(페이지 목록 파일)를 제공할 때 |
| `sitemap` | ⛔ 미구현 (선택 시 명시적 오류로 중단) |
| `crawl` | ⛔ 미구현 (선택 시 명시적 오류로 중단) |

확인법: 브라우저나 `curl`로 `<base_url>/llms.txt`를 열어 마크다운 링크 목록(`- [제목](url)`)이 나오면
`llms.txt`를 쓴다. 사실상 현재 선택지는 `llms.txt` 하나다.

## glossary.json

번역어를 고정해 페이지·재실행 간 일관성을 유지한다. `원어 → 한국어` 단순 매핑이다.

```json
{
  "model": "모델",
  "embedding": "임베딩",
  "pull": "pull",
  "ollama": "Ollama"
}
```

- 정의된 원어는 번역 시 항상 그 번역어로 반영된다.
- 번역하고 싶지 않은 고유명사·명령어는 원어를 그대로 값으로 둔다(예: `"pull": "pull"`).
- 비어 있어도(`{}`) 동작한다.

## 사용법

### 1. 수집 (Go)

```bash
go run ./cmd/collect <siteID>      # 예: go run ./cmd/collect ollama
```

- `config.json`을 읽어 `llms.txt`에서 페이지 목록을 확보하고, 패턴으로 거른 뒤 각 원문을 내려받아
  `sites/<siteID>/raw/`에 원본 경로 구조로 저장한다.
- 외부 사이트에 과도한 요청을 피하도록 요청 간격·재시도를 둔다.
- 재실행하면 변경된 원문만 갱신한다(증분). 출력: `갱신=N, 스킵=N, 실패=N`.
- 플래그: `-sites <dir>` (사이트 폴더 루트, 기본 `./sites`).

### 2. 번역 (Claude)

번역 텍스트 생성은 Claude가 하고, 대상 선별·완료 기록은 CLI가 한다.

```bash
go run ./cmd/translate plan <siteID>     # 1) 번역 대상(미번역·변경분) 목록 출력
# 2) Claude Code에서 /translate <siteID> — 각 원문을 한국어 문서로 재구성해 output/에 작성
go run ./cmd/translate commit <siteID>   # 3) output/에 있는 페이지의 번역 해시를 매니페스트에 기록
```

- `commit`은 `output/`에 파일이 실제로 있는 페이지만 기록하고, 없으면 건너뛴다.
- `commit` 후 그 페이지는 다음 `plan`에서 빠진다(증분).
- 원문이 OpenAPI 스펙 등 기계친화 포맷이면, 그대로 옮기지 않고 표·설명으로 재구성한다.

## 증분 동작

`manifest.json`이 페이지별로 원문 해시(`source_hash`)와 번역한 원문 해시(`translated_hash`)를 들고 있다.

- 수집: 새 원문 해시 ≠ `source_hash` → 갱신, 같으면 스킵.
- 번역: `translated_hash` ≠ 현재 `source_hash` → 번역 대상, 같으면 스킵.
- 원문이 바뀐 페이지만 다시 수집·번역되므로, 대량 문서도 변경분만 처리한다.

## 빌드·테스트

```bash
go build ./...
go test ./...
```

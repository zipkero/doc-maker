# translate-api-automation — 분석·설계

## 승인 전 확인
- 부모 doc-localization SPEC §3의 "외부 번역 API 미호출" 제약을, 에이전트 경로 공존을 전제로 API 경로가
  대체·확장하는 것이 맞는지 — spec.md §승인 전 확인이 이 판단을 main(사용자)에게 위임하고 있다. 본 분석은
  "공존·추가"를 전제로 §1–§5를 구성했으나, 이 전제가 부정되면 §1 구조와 §4 영향 범위가 무효가 된다.
  관련 본문: §1, §4
- Anthropic Go SDK의 Message Batches 바인딩 정확 심볼(타입·메서드명)이 현 시점 미확정이다. 동기 경로의
  심볼(`anthropic.NewClient`, `client.Messages.New`, `MessageNewParams`, `System` 캐시 컨트롤)은 확인된
  반면, 배치 제출/폴링/결과 수집 심볼은 추정이다. implement 단계 진입 전 SDK 소스로 확정이 필요하다.
  관련 본문: §2, §3, §5 D-d

## 근거

읽은 spec 범위: spec.md 전체(§1 범위 ~ §5 완료 조건 8항). 범위는 §1로 한정하며 요구사항 신규 추가는 하지 않는다.

코드베이스에서 확인한 사실:
- `go.mod`는 `module doc-maker` + `go 1.26.4`만 있고 외부 의존성이 전혀 없다. `github.com/anthropics/anthropic-sdk-go`는
  미존재 — 이 feature에서 처음 추가된다.
- `cmd/translate/main.go`: 최상위 `-sites` 플래그, `plan`/`commit` 두 서브커맨드, `runPlan`/`runCommit`/`loadSite`
  구조. plan은 `translator.SelectTargets`로 대상 목록만 출력하고, commit은 `output/`에 파일이 존재하는 페이지만
  `translator.TranslatedPath`로 확인 후 `translator.CommitTranslation`으로 TranslatedHash를 기록한다. 번역 텍스트
  생성은 이 CLI에 전혀 없다(현재는 Claude 에이전트 절차가 담당).
- `internal/translator/translator.go`: `TranslationTarget{PageURL, LocalPath, SourceHash}`와
  `SelectTargets(site) ([]TranslationTarget, error)`. 선별 기준은 `manifest.NeedsTranslation`(TranslatedHash !=
  현재 SourceHash). raw 디렉터리 부재 시 빈 목록.
- `internal/translator/store.go`: `TranslatedPath(site, pageURL)`(URL 경로 → output 경로, 확장자 .md 정규화),
  `SaveTranslation(site, pageURL, translatedPath, content, sourceHash)`(파일 저장 + TranslatedHash 기록 동시 수행),
  `CommitTranslation(site, pageURL, sourceHash)`(파일은 안 건드리고 TranslatedHash만 기록). 두 함수 모두
  매니페스트를 로드·갱신·저장한다.
- `internal/translator/glossary.go`: `Glossary map[string]string`(원어→한국어), `LoadGlossary(path)`. 파일 부재·빈
  파일은 빈 매핑(오류 아님), 파싱 실패만 오류. `site.GlossaryPath()`는 `<siteDir>/glossary.json`로 고정.
- `internal/config/config.go`: `Site{ID, SiteDir, BaseURL, SourceType, IncludePatterns, ExcludePatterns}`. JSON
  태그는 `base_url`/`source_type`/`include_patterns`/`exclude_patterns`. **model 필드는 현재 없다.** `RawDir/OutputDir/
  GlossaryPath/ManifestDir` 헬퍼는 모두 사이트 폴더 규약으로 경로를 도출. `Load`는 BaseURL·SourceType 누락 시 오류.
- `internal/manifest/manifest.go`: `Entry{SourceHash, SourcePath, TranslatedHash}`, `Load/Save/Get/Set/
  NeedsTranslation/HashContent`. `HashContent`는 SHA-256 hex.
- `internal/fetcher/fetcher.go`: 기존 외부호출 관례 — `Config{MaxRetries, Delay, BackoffFactor}`, `DefaultConfig`
  (1초 간격·3회·배율 2), 지수 backoff, 4xx=`ErrPermanent`(즉시 실패)·5xx/네트워크=재시도. `HTTPDoer` 인터페이스로
  테스트 주입. 외부 호출의 rate limit/retry 관례가 이 패키지에 이미 자리잡혀 있다.
- 번역기 함수 호출자: 프로덕션 코드에서는 `cmd/translate/main.go` 한 곳뿐(나머지는 모두 테스트). 새 서브커맨드를
  추가해도 기존 plan/commit 호출 경로에 간섭하지 않는다.

확인된 Anthropic API 사실(외부 자료로 확정, 본 분석에서 그대로 사용): Sonnet 4.6 모델 ID `claude-sonnet-4-6`. Go SDK는
`github.com/anthropics/anthropic-sdk-go`, `anthropic.NewClient()`가 `ANTHROPIC_API_KEY`를 자동 사용. 동기 단건은
`client.Messages.New(ctx, MessageNewParams{Model, MaxTokens, System([]TextBlockParam, 마지막 블록에 CacheControl),
Messages})`, 큰 출력은 `client.Messages.NewStreaming`. 프롬프트 캐싱은 System 마지막 블록에
`anthropic.NewCacheControlEphemeralParam()`, 적중은 `resp.Usage.CacheReadInputTokens`. 프롬프트 캐싱 최소 토큰:
Sonnet 4.6는 2048 토큰(미만이면 캐시 없이 정상 동작).

추정(확정 아님): Go SDK의 Message Batches 바인딩 정확 심볼. Batch 엔드포인트(`POST /v1/messages/batches`)는 존재가
확인되나, Go SDK의 제출/폴링/결과 다운로드 타입·메서드명은 본 분석에서 "SDK 배치 API"로 추상 기술하며 implement 진입
전 SDK 소스 확인을 전제로 한다.

## 1. 구조

경계 배치의 핵심은 "결정적 부분(선별·경로·증분·매니페스트)은 기존 그대로, 비결정적 부분(재구성 텍스트 생성)만 API
호출로 대체"다(SPEC §1, §3).

- **새 패키지 — API 호출 경계.** Anthropic SDK 의존과 호출 로직은 `internal/translator`에 직접 넣지 않고 별도 패키지
  (예: `internal/llmtranslate` 또는 `internal/apitranslate`)에 둔다. 근거: ① `internal/translator`는 현재 외부 의존
  0이고 결정적 로직(선별·저장·해시)만 담는 경계인데, 여기에 SDK·네트워크·키 처리를 섞으면 그 경계가 무너진다. ②
  `internal/fetcher`가 이미 "외부 호출은 별도 패키지" 선례를 세웠다. ③ 이 경계가 동기/Batch 두 실행 방식, 프롬프트
  조립, 캐싱, 재시도를 모두 가진다. 이 패키지는 입력으로 `[]TranslationTarget`(또는 원문+규칙+용어집)을 받아 재구성
  텍스트(또는 페이지별 결과/실패)를 산출하는 것으로 경계를 닫는다.
- **`internal/translator`는 그대로 재사용.** `SelectTargets`(대상 선별), `TranslatedPath`(출력 경로),
  `SaveTranslation`/`CommitTranslation`(저장·기록), `LoadGlossary`(용어집)는 변경 없이 호출만 한다. 새 패키지는
  translator에 의존하되 그 역은 없다(단방향).
- **`cmd/translate` 새 서브커맨드 경계.** 기존 `plan`/`commit` 옆에 API 실행 진입점을 추가한다. 이 진입점이 site
  로드 → SelectTargets → 용어집/규칙 로드 → 새 API 패키지 호출 → 결과 저장 → 매니페스트 기록의 오케스트레이션을
  맡는다. plan/commit 코드는 손대지 않는다. (동기/Batch를 한 서브커맨드의 모드로 둘지 둘로 나눌지는 §5 D-a.)
- **규칙 텍스트(재구성 규칙·스타일 기준)의 위치**도 이 새 패키지 또는 그 인접에 둔다(§5 D-b). 프롬프트의 공유
  프리픽스를 구성하는 자산이므로 API 경계 안에 두는 것이 응집도가 높다.

## 2. 데이터 흐름

```
[CLI 서브커맨드]
   loadSite(siteID) ──→ config.Site (+ model 필드)
        │
   translator.SelectTargets(site) ──→ []TranslationTarget  (미번역·변경분만; 기존 증분 판정 그대로)
        │
   translator.LoadGlossary(site.GlossaryPath()) ──→ Glossary
   규칙·스타일 기준 텍스트 로드 ───────────────────→ 공유 프리픽스 자산
        │
   ┌────────────── 프롬프트 구성 ──────────────┐
   │ System: [재구성 규칙 + 용어집 + 스타일]   │  ← 페이지 불변 → 마지막 블록에 CacheControl
   │ Messages(user): 페이지 원문(target.LocalPath 읽기)│  ← 페이지별 가변
   └───────────────────────────────────────────┘
        │
        ├─[동기 단건 경로] target마다:
        │     client.Messages.New / NewStreaming(model, System(cache), 원문)
        │        └→ 재구성 텍스트 → translator.Save/Commit → (다음 target)
        │
        └─[Batch 경로] 전체 target을:
              제출(요청 N건) → batchID
                 ↓ 폴링(상태: 제출됨 → 진행 중 → 완료/부분완료)
              완료 시 결과 수집 → 페이지별로 매핑 → 성공분만 저장
        │
   저장: translator.TranslatedPath(site, pageURL) 경로에 재구성 텍스트 기록
   기록: TranslatedHash = target.SourceHash  (저장 시 즉시 vs commit 재사용은 §5 D-g)
```

상태 전이(Batch): 제출 직후 batchID 확보 → 폴링 루프에서 진행 상태 관찰(SPEC §5.2 "제출·진행·완료 상태가
관찰된다") → 완료 시 각 요청 결과를 페이지(custom_id 등 페이지 식별자)로 역매핑. 부분 실패 시 완료된 페이지 결과는
수집·저장하고 실패 페이지는 사유와 함께 보고한다(SPEC §5.7). Batch는 비동기라 요청 타임아웃 제약을 받지 않는다
(SPEC §3).

에러·실패 경로:
- **API 키 없음**: 호출 전(클라이언트 구성 또는 첫 호출 시점)에 명확한 오류로 중단하고 조용히 실패하지 않는다
  (SPEC §5.7). 종료 코드 1.
- **페이지 단위 실패**: 동기 경로는 한 페이지 호출 실패가 전체를 중단시키지 않고, 성공한 페이지는 저장하며 실패
  페이지는 URL·사유를 모아 보고한다. Batch 경로는 부분 실패를 같은 원칙으로 처리한다(SPEC §5.7).
- **rate limit / 일시 오류**: 동기 경로는 재시도를 고려한다(SPEC §3). SDK가 429/5xx 자동 재시도를 제공하면 그에
  맡기고, 추가 backoff가 필요하면 fetcher의 관례(지수 backoff)를 참고한다.

동시성: 동기 경로는 순차 처리를 기본으로 둔다(rate limit·결과 보고 단순성). 병렬화는 spec 범위 밖이므로 도입하지
않는다. Batch는 본질적으로 일괄 제출이므로 페이지별 동시성 이슈가 없다.

## 3. 인터페이스

경계를 가로지르는 계약만 기술한다.

- **새 CLI 서브커맨드(들).** 기존 형태 `translate [플래그] <서브커맨드> <siteID>`를 따른다. API 실행 진입점은
  미번역·변경분(SelectTargets 결과)을 입력으로 받아 재구성·저장·기록까지 수행하고, 처리 건수·실패 페이지·출력
  경로를 보고한다. 동기/Batch 선택은 모드 플래그 또는 별도 서브커맨드명으로 노출한다(§5 D-a 결정에 종속).
- **번역 생성 함수의 입출력 계약(새 API 패키지).** 입력: 원문 + 재구성 규칙 + 용어집(+ 모델 ID). 출력: 재구성된
  한국어 마크다운 텍스트. 페이지 다건을 다룰 때는 페이지별 (성공 텍스트 | 실패 사유)를 구분해 반환해 부분 실패
  보고(SPEC §5.7)와 성공분 저장이 가능하게 한다. 동기 단건과 Batch 일괄 두 진입 형태가 같은 입력 의미를
  공유한다.
- **`config.json`의 model 필드 추가.** `Site`에 모델 지정 필드(JSON 키 예: `"model"`)를 추가한다. 비어 있으면
  기본값 `claude-sonnet-4-6`로 호출하고, 지정 시 그 값으로 호출한다(SPEC §5.4). 기존 필수 검증(BaseURL·SourceType)
  에는 포함하지 않는다 — model은 선택 필드이며 부재 시 기본값으로 동작해야 한다.
- **기존 함수와의 접점.** `translator.SelectTargets`(입력 선별), `translator.TranslatedPath`(출력 경로),
  `translator.SaveTranslation` 또는 `translator.CommitTranslation`(저장·기록), `translator.LoadGlossary`(용어집).
  새 경로는 이 계약들을 그대로 호출하고 시그니처를 바꾸지 않는다.

내부 helper 시그니처는 명시하지 않는다(implement 단계 소관).

## 4. 영향 범위

직접 건드리는 기존 모듈:
- **`internal/config`**: `Site`에 model 필드 추가. 직접 의존(호출자)은 `cmd/*`의 `loadSite`와 translator 함수들
  (site를 인자로 받음)이지만, **필드 추가만으로는 기존 동작이 깨지지 않는다** — JSON에 키가 없으면 zero value(빈
  문자열)로 언마샬되고, 새 경로가 빈 값을 기본 모델로 해석한다. 기존 plan/commit은 model을 읽지 않으므로 무영향.
  하위호환: 기존 `config.json`(model 부재)은 그대로 유효하며 기본값 처리로 흡수된다(§5 D-c).
- **`cmd/translate`**: 새 서브커맨드 추가(plan/commit 코드 불변). `main`의 switch에 케이스가 늘고, 새 run 함수와
  새 API 패키지 호출이 추가된다.
- **`internal/translator`**: 저장·commit·선별·용어집·경로 함수를 **재사용만** 한다(수정 없음). 새 API 패키지가 이
  패키지에 의존하는 신규 의존 방향이 생긴다.
- **`go.mod`**: `github.com/anthropics/anthropic-sdk-go` 의존 추가(+ `go.sum`). 현재 의존 0이므로 첫 외부 의존이다.
- **새 패키지**(예: `internal/llmtranslate`): 신규 생성. 기존 코드에 대한 호출자는 새 CLI 서브커맨드뿐.

간접 의존 확인: translator 함수의 프로덕션 호출자는 `cmd/translate/main.go` 한 곳뿐(`SelectTargets`/`TranslatedPath`/
`CommitTranslation`/`SaveTranslation`/`LoadGlossary` grep 결과 나머지는 모두 테스트 파일). 따라서 새 경로 추가가
기존 호출 경로를 깨뜨리지 않음을 코드로 확인했다.

마이그레이션: 위 config 기본값 처리 외에는 해당 없음. 기존 에이전트 기반 `/translate` 절차는 변경 없이 공존하며,
출력 경로·매니페스트·증분 규약을 동일하게 쓰므로 두 경로가 같은 산출물 구조를 공유한다(SPEC §1, §5.8).

## 5. Decision Points

### D-a. 동기/Batch 경로의 CLI 노출 방식
- 고려 옵션: (1) 한 서브커맨드 + 모드 플래그(예: `--batch`), (2) 별도 서브커맨드 두 개.
- 트레이드오프: (1)은 공통 오케스트레이션(site 로드→선별→프롬프트 조립→저장·기록)을 한 진입점에 모아 중복이
  적고 사용 표면이 단순; 실행 방식만 분기. (2)는 두 경로의 수명주기 차이(Batch는 제출→폴링→수집의 비동기 다단계,
  동기는 단일 실행)를 서브커맨드 의미로 명확히 드러내지만, plan/commit과 합쳐 서브커맨드가 늘어 사용 표면이 커진다.
- 채택: (1) 한 서브커맨드 + 모드 플래그, 기본은 Batch(SPEC §1이 "Batch 기본 + 동기 단건"을 명시).
- 근거: 두 경로가 입력(SelectTargets 결과)·출력(저장·기록)·프롬프트 조립을 공유하므로 한 진입점이 응집도가 높고,
  spec이 둘을 같은 feature의 두 실행 방식으로 규정한다.

### D-b. 재구성 규칙·용어집·스타일 기준의 조립·캐시 출처
- 고려 옵션: 규칙 텍스트의 출처를 (1) `.claude/commands/translate.md` 재사용, (2) Go 내 상수, (3) 별도 파일(사이트
  무관 공용 자산).
- 트레이드오프: (1)은 에이전트 절차와 단일 출처(규칙 표류 방지)이나 그 파일은 커맨드 문서 포맷이라 프롬프트로
  바로 쓰기엔 잡음(사용법·경로표 등)이 섞여 정제가 필요. (2)는 빌드에 고정돼 캐시 프리픽스가 바이트 안정적이라
  프롬프트 캐싱에 유리하고 외부 파일 의존이 없으나, 규칙 수정이 코드 변경이 된다. (3)은 수정 용이하나 파일 부재·표류
  관리 부담.
- 채택: (2) Go 내 상수로 규칙·스타일 기준을 보유하고, 용어집은 사이트별 `LoadGlossary`로 합류. 용어집은 사이트마다
  다르므로 공유 프리픽스(System)의 안정 부분(규칙·스타일)과 가변 부분(용어집)을 구분해 조립한다. 규칙 텍스트의
  내용 기준은 translate.md와 일치시킨다(원문 병기 없음·보일러플레이트 제거·문서 끝 출처 한 줄·기계친화 포맷 재구성).
- 근거: 프롬프트 캐싱은 프리픽스 바이트 안정성에 민감하고(타임스탬프·비결정 직렬화 금지), 규칙은 페이지·사이트
  불변이라 상수가 캐시에 최적. spec은 규칙을 "기존 translate 절차와 동일하게 유지"로 규정(SPEC §3)하므로 내용은
  translate.md를 권위 출처로 삼되 형식만 프롬프트용으로 정제한다.

### D-c. 모델 설정 위치와 기본값
- 고려 옵션: (1) `config.json`의 `model` 필드(기본 `claude-sonnet-4-6`), (2) CLI 플래그, (3) 환경변수.
- 트레이드오프: (1)은 사이트 단위 설정으로 spec(SPEC §3 "사이트 설정으로 지정")과 정확히 일치하고 영속적. (2)는
  일회성·디버깅엔 편하나 사이트 영속 설정이 아님. (3)은 키와 혼동 우려.
- 채택: (1) `config.json` model 필드, 미지정 시 `claude-sonnet-4-6`.
- 근거: spec이 모델 지정을 "사이트 설정"으로 명시(SPEC §3, §5.4). 필드 부재는 기본값으로 흡수되어 기존 config와
  하위호환(§4).

### D-d. Batch 폴링/결과 매핑과 부분 실패 처리
- 고려 옵션: 결과 매핑 키로 (1) 페이지 URL을 요청 식별자(custom_id 등)에 직접 사용, (2) 인덱스→페이지 별도
  매핑 테이블. 폴링은 (a) 고정 간격, (b) backoff.
- 트레이드오프: (1)은 결과를 받아 곧장 pageURL로 저장·기록할 수 있어 매핑 단순. (2)는 식별자 제약이 있을 때
  우회 가능하나 상태 보관 필요. 폴링 (a)는 단순, (b)는 장시간 작업에 호출량 절감.
- 채택: 요청 식별자에 페이지 식별자를 실어 결과를 pageURL로 직접 역매핑(1), 폴링은 합리적 간격의 backoff(b).
  부분 실패는 완료 페이지만 저장하고 실패 페이지는 URL·사유 수집·보고(SPEC §5.7).
- 근거: 저장·기록 계약이 pageURL 키이므로 결과를 pageURL로 받으면 접점이 매끄럽다. **단, Go SDK 배치 바인딩의
  정확 심볼(custom_id 필드명·폴링/다운로드 메서드)은 미확정이므로 implement 진입 전 SDK 소스 확인이 전제다(승인 전
  확인 항목과 동일).**

### D-e. 동기 경로 큰 출력 스트리밍
- 고려 옵션: (1) 항상 `Messages.New`(비스트리밍), (2) 항상 `NewStreaming`, (3) 출력 크기 예상에 따라 분기.
- 트레이드오프: 큰 출력 비스트리밍은 SDK HTTP 타임아웃 위험(잘림). 스트리밍은 타임아웃을 피하고 max_tokens를
  크게(~64K) 줄 수 있으나 누적 처리 코드가 필요.
- 채택: 동기 경로는 스트리밍 기본(`NewStreaming` + 최종 메시지 누적), max_tokens는 충분히 크게.
- 근거: 재구성 문서는 단일 페이지라도 길 수 있고(OpenAPI 표 가공 등), spec이 "동기 경로에서 큰 페이지는 잘림을
  피하도록 처리(충분한 출력 한도/스트리밍)"를 요구(SPEC §3). Batch는 비동기라 타임아웃 무관.

### D-f. 캐시 프리픽스가 최소 토큰 미만일 때
- 고려 옵션: (1) 항상 CacheControl 부착, (2) 프리픽스가 모델 최소 토큰 이상일 때만 부착, (3) 캐시 비활성 옵션.
- 트레이드오프: Sonnet 4.6 최소 캐시 토큰은 2048. 미만이면 CacheControl을 붙여도 조용히 캐시되지 않고 정상
  동작(오류 아님). (1)은 단순하고 무해, (2)는 불필요한 마커를 피하나 토큰 추정 로직 필요.
- 채택: (1) 공유 프리픽스(규칙+용어집+스타일)에 CacheControl을 항상 부착하되, 미충족 시 캐시 없이 정상 동작함을
  설계 전제로 둔다. 적중은 `resp.Usage.CacheReadInputTokens`로 관찰(SPEC §5.5).
- 근거: spec이 "최소 캐시 토큰 충족 시에만 적용되며 미충족 시 캐시 없이 정상 동작"을 명시(SPEC §3, §5.5). 마커
  부착은 무해하므로 별도 분기 없이 단순화한다.

### D-g. 저장 시 기존 commit 재사용 vs 즉시 기록
- 고려 옵션: (1) `SaveTranslation`(파일 저장 + TranslatedHash 즉시 기록)을 페이지마다 호출, (2) 모든 파일 저장 후
  `CommitTranslation`을 일괄 호출(기존 에이전트 절차 방식), (3) 파일은 `SaveTranslation`로 저장하되 기록은 별도
  `commit` 서브커맨드 실행에 맡김.
- 트레이드오프: (1)은 페이지 단위 원자성이 높아 부분 실패 시 성공분만 정확히 기록되고 재실행이 깔끔(증분 일관).
  (2)는 에이전트 절차와 동일 패턴이나 저장과 기록 사이 실패 시 불일치 여지. (3)은 책임 분리가 명확하나 사용자가
  두 명령을 실행해야 함.
- 채택: (1) 페이지별 `SaveTranslation`로 저장과 TranslatedHash 기록을 함께 수행. 부분 실패 시 성공 페이지만
  기록되어 재실행 시 실패분만 대상이 된다(SPEC §5.3, §5.7).
- 근거: API 경로는 한 명령으로 생성·저장·기록을 끝내는 자동화가 목표(SPEC §2)이므로 페이지 단위 즉시 기록이 부분
  실패·증분 동작과 가장 잘 맞는다. `CommitTranslation`은 기존 에이전트 절차(외부에서 파일 작성 후 기록)용으로 그대로
  공존시킨다.

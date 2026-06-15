# doc-localization

## 요약
외부 기술 문서 사이트를 설정 기반으로 수집해 한국어로 번역·미러링하는 도구. 올라마를 첫 타깃 겸 레퍼런스로
삼고, 다른 문서 사이트도 같은 틀로 추가할 수 있게 한다.

## 상태
- [x] SPEC
- [x] ANALYSIS
- [x] IMPLEMENT

## 문서
- [spec.md](./spec.md)
- [analysis.md](./analysis.md) (ANALYSIS 단계에서 생성)
- [implement.md](./implement.md) (IMPLEMENT 단계에서 생성)

## 작업 히스토리
- 2026-06-15: SPEC 작성
- 2026-06-15: ANALYSIS 작성
- 2026-06-15: IMPLEMENT 체크리스트 작성
- 2026-06-15: IMPLEMENT 완료
- 2026-06-15: 요구사항 변경 — 번역을 "원문 충실 번역"에서 "읽기 좋은 한국어 문서로 재구성"으로
  (SPEC §1·§2·§5.4, ANALYSIS D11 추가·D1 갱신, implement task-010, translate 커맨드 갱신)
- 2026-06-15: 저장소 레이아웃 변경 — 종류별 폴더(configs/raw/manifests/output/glossary)에서
  사이트별 폴더(`sites/<siteID>/`)로. 식별자=폴더명, 설정은 3종값(출력·용어집·매니페스트는 폴더 규약).
  (SPEC §5.1·§5.2·§5.4·§5.8, ANALYSIS D3·D4·D5·D10·§1·§3, implement task-001/002/007/008/010/012)
- 2026-06-15: 번역 진입점 CLI(`cmd/translate` plan/commit) 추가 — 번역 대상 선별·완료 기록을
  손편집 없이 처리(task-013, SPEC §5.4·§5.5)

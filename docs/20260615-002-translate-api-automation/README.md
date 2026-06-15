# translate-api-automation

## 요약
번역 재구성(2단계) 생성을 Anthropic Messages API로 자동화하는 경로를 추가한다. 기본 모델 Sonnet 4.6, 대량은
Batch + 프롬프트 캐싱, 단건·증분은 동기 호출로 처리하며, 기존 에이전트 기반 `/translate`와 공존한다.

## 상태
- [x] SPEC
- [x] ANALYSIS
- [ ] IMPLEMENT

## 문서
- [spec.md](./spec.md)
- [analysis.md](./analysis.md) (ANALYSIS 단계에서 생성)
- [implement.md](./implement.md) (IMPLEMENT 단계에서 생성)

## 작업 히스토리
- 2026-06-15: SPEC 작성
- 2026-06-15: ANALYSIS 작성
- 2026-06-15: IMPLEMENT 체크리스트 작성

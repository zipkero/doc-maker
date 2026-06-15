# 스트리밍

`/api/generate`처럼 일부 엔드포인트는 기본적으로 응답을 스트리밍으로 전송합니다.
이때 응답은 줄바꿈으로 구분된 JSON 형식(NDJSON, `application/x-ndjson` 콘텐츠 타입)으로
전달되며, 각 줄이 하나의 부분 응답입니다. 예시는 다음과 같습니다.

```json
{"model":"gemma4","created_at":"2025-10-26T17:15:24.097767Z","response":"That","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.109172Z","response":"'","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.121485Z","response":"s","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.132802Z","response":" a","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.143931Z","response":" fantastic","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.155176Z","response":" question","done":false}
{"model":"gemma4","created_at":"2025-10-26T17:15:24.166576Z","response":"!","done":true, "done_reason": "stop"}
```

마지막 줄에서 `done`이 `true`가 되며, 여기에 `done_reason` 같은 종료 정보가 포함됩니다.

## 스트리밍 비활성화

스트리밍을 지원하는 엔드포인트라면 요청 본문에 `{"stream": false}`를 넣어 비활성화할 수
있습니다. 이 경우 응답은 여러 줄로 나뉘지 않고 `application/json` 형식의 단일 객체로
한 번에 반환됩니다.

```json
{"model":"gemma4","created_at":"2025-10-26T17:15:24.166576Z","response":"That's a fantastic question!","done":true}
```

## 스트리밍과 비스트리밍, 언제 쓸까

**스트리밍(기본값)**

- 응답을 실시간으로 생성해 보여줄 때
- 체감 지연을 낮추고 싶을 때
- 긴 응답을 생성할 때

**비스트리밍**

- 응답 처리를 단순하게 하고 싶을 때
- 짧은 응답이나 구조화된 출력을 다룰 때
- 일부 애플리케이션에서 다루기가 더 쉬울 때

> 원문: https://docs.ollama.com/api/streaming

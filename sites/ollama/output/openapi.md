# Ollama API 레퍼런스 개요

Ollama HTTP API 전체를 담은 마스터 명세 요약입니다. 각 엔드포인트의 자세한 요청/응답 예시는 `api/`
하위 문서를 참고하고, 이 문서는 전체 엔드포인트와 공용 스키마를 한눈에 보는 레퍼런스로 사용합니다.

- **버전**: 0.1.0 (OpenAPI 3.1.0)
- **라이선스**: MIT
- **서버 base URL**: `http://localhost:11434`
- **인증**: 기본은 인증 없음. 일부 환경에서 Bearer 토큰(`Authorization: Bearer <API Key>`)을 사용할 수 있습니다.

## 엔드포인트 목록

| 메서드 | 경로 | operationId | 설명 |
|---|---|---|---|
| `POST` | `/api/generate` | `generate` | 주어진 프롬프트에 대한 응답을 생성합니다 |
| `POST` | `/api/chat` | `chat` | 사용자와 어시스턴트 대화에서 다음 채팅 메시지를 생성합니다 |
| `POST` | `/api/embed` | `embed` | 입력 텍스트를 나타내는 벡터 임베딩을 생성합니다 |
| `GET` | `/api/tags` | `list` | 로컬 모델 목록과 상세 정보를 조회합니다 |
| `GET` | `/api/ps` | `ps` | 현재 실행 중(메모리 적재)인 모델 목록을 조회합니다 |
| `POST` | `/api/show` | `show` | 특정 모델의 상세 정보를 조회합니다 |
| `POST` | `/api/create` | `create` | 기존 모델을 기반으로 새 모델을 생성합니다 |
| `POST` | `/api/copy` | `copy` | 모델을 다른 이름으로 복사합니다 |
| `POST` | `/api/pull` | `pull` | 레지스트리에서 모델을 내려받습니다(pull) |
| `POST` | `/api/push` | `push` | 모델을 레지스트리에 게시합니다(push) |
| `DELETE` | `/api/delete` | `delete` | 모델을 삭제합니다 |
| `GET` | `/api/version` | `version` | Ollama 버전 정보를 조회합니다 |

스트리밍을 지원하는 엔드포인트(`generate`, `chat`, `create`, `pull`, `push`)는 `application/x-ndjson`
형식으로 부분 이벤트를 여러 번 전송하며, 마지막 이벤트의 `done`(또는 최종 `status`)으로 완료를 표시합니다.

상세 요청/응답 예시(curl, 구조화 출력, 이미지 입력 등)는 각 엔드포인트 문서를 참고하세요.

## 공용 스키마 (components/schemas)

### 생성 (Generate)

#### `ModelOptions`

생성을 제어하는 런타임 옵션입니다. 정의되지 않은 추가 속성도 허용합니다(`additionalProperties: true`).

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `seed` | integer |  | 재현 가능한 출력을 위한 난수 시드 |
| `temperature` | number(float) |  | 생성의 무작위성 제어(높을수록 무작위) |
| `top_k` | integer |  | 다음 토큰 후보를 가장 가능성 높은 K개로 제한 |
| `top_p` | number(float) |  | 뉴클리어스 샘플링의 누적 확률 임계값 |
| `min_p` | number(float) |  | 토큰 선택의 최소 확률 임계값 |
| `stop` | string \| string[] |  | 생성을 중단시킬 정지 시퀀스 |
| `num_ctx` | integer |  | 컨텍스트 길이(토큰 수) |
| `num_predict` | integer |  | 생성할 최대 토큰 수 |

#### `GenerateRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `prompt` | string |  | 모델이 응답을 생성할 입력 텍스트 |
| `suffix` | string |  | fill-in-the-middle 모델용. 프롬프트 뒤·모델 응답 앞에 오는 텍스트 |
| `images` | string[] |  | 이미지 입력 지원 모델용 base64 인코딩 이미지 목록 |
| `format` | string \| object |  | 구조화 출력 형식. `"json"` 문자열 또는 JSON 스키마 객체 |
| `system` | string |  | 모델에 적용할 시스템 프롬프트 |
| `stream` | boolean |  | true면 부분 응답을 스트리밍. 기본값 `true` |
| `think` | boolean \| string |  | 사고 과정 출력 여부. 불리언 또는 `high`/`medium`/`low`(지원 모델 한정) |
| `raw` | boolean |  | true면 프롬프트 템플릿 없이 모델의 원시 응답 반환 |
| `keep_alive` | string \| number |  | 모델 메모리 유지 시간(예: `5m`, 즉시 해제는 `0`) |
| `options` | `ModelOptions` |  | 생성을 제어하는 런타임 옵션 |
| `logprobs` | boolean |  | 출력 토큰의 로그 확률 반환 여부 |
| `top_logprobs` | integer |  | `logprobs` 활성화 시 각 토큰 위치에서 반환할 상위 토큰 수 |

#### `GenerateResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 모델 이름 |
| `created_at` | string | 응답 생성 시각(ISO 8601) |
| `response` | string | 모델이 생성한 텍스트 응답 |
| `thinking` | string | 모델이 생성한 사고 과정 출력 |
| `done` | boolean | 생성 완료 여부 |
| `done_reason` | string | 생성이 종료된 이유 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 프롬프트의 입력 토큰 수 |
| `prompt_eval_duration` | integer | 프롬프트 평가 시간(나노초) |
| `eval_count` | integer | 응답으로 생성된 토큰 수 |
| `eval_duration` | integer | 토큰 생성 시간(나노초) |
| `logprobs` | `Logprob[]` | `logprobs` 활성화 시 생성 토큰의 로그 확률 정보 |

#### `GenerateStreamEvent`

스트리밍(`application/x-ndjson`) 모드에서 전송되는 부분 이벤트입니다. 필드 구성은 `GenerateResponse`와
동일하며, `response`·`thinking`은 해당 청크 분량만 담고 `done`이 `true`인 이벤트가 마지막입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 모델 이름 |
| `created_at` | string | 응답 생성 시각(ISO 8601) |
| `response` | string | 이 청크의 생성 텍스트 |
| `thinking` | string | 이 청크의 사고 과정 출력 |
| `done` | boolean | 스트림 종료 여부 |
| `done_reason` | string | 스트리밍 종료 이유 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 프롬프트의 입력 토큰 수 |
| `prompt_eval_duration` | integer | 프롬프트 평가 시간(나노초) |
| `eval_count` | integer | 생성된 출력 토큰 수 |
| `eval_duration` | integer | 토큰 생성 시간(나노초) |

### 채팅 (Chat)

#### `ChatMessage`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `role` | string(`system`/`user`/`assistant`/`tool`) | ✅ | 메시지 작성 주체 |
| `content` | string | ✅ | 메시지 본문 텍스트 |
| `images` | string[] |  | 멀티모달 모델용 인라인 이미지(base64) 목록 |
| `tool_calls` | `ToolCall[]` |  | 모델이 생성한 도구 호출 요청 |

#### `ToolCall`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `function` | object |  | 호출 대상 함수 정보 |
| `function.name` | string | ✅ | 호출할 함수 이름 |
| `function.description` | string |  | 함수가 수행하는 작업 설명 |
| `function.arguments` | object |  | 함수에 전달할 인자(JSON 객체) |

#### `ToolDefinition`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `type` | string(`function`) | ✅ | 도구 유형(항상 `function`) |
| `function` | object | ✅ | 모델에 노출할 함수 정의 |
| `function.name` | string | ✅ | 모델에 노출되는 함수 이름 |
| `function.description` | string |  | 사람이 읽을 수 있는 함수 설명 |
| `function.parameters` | object | ✅ | 함수 매개변수의 JSON 스키마 |

#### `ChatRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `messages` | `ChatMessage[]` | ✅ | 채팅 기록(각 항목은 `role`과 `content`를 가짐) |
| `tools` | `ToolDefinition[]` |  | 채팅 중 모델이 호출할 수 있는 함수 도구 목록 |
| `format` | string(`json`) \| object |  | 응답 형식. `json` 또는 JSON 스키마 |
| `options` | `ModelOptions` |  | 생성을 제어하는 런타임 옵션 |
| `stream` | boolean |  | 스트리밍 여부. 기본값 `true` |
| `think` | boolean \| string |  | 사고 과정 출력 여부. 불리언 또는 `high`/`medium`/`low`(지원 모델 한정) |
| `keep_alive` | string \| number |  | 모델 메모리 유지 시간(예: `5m`, 즉시 해제는 `0`) |
| `logprobs` | boolean |  | 출력 토큰의 로그 확률 반환 여부 |
| `top_logprobs` | integer |  | `logprobs` 활성화 시 각 토큰 위치에서 반환할 상위 토큰 수 |

#### `ChatResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 이 메시지를 생성한 모델 이름 |
| `created_at` | string(date-time) | 응답 생성 시각(ISO 8601) |
| `message` | object | 생성된 메시지 객체(아래 하위 필드 참조) |
| `message.role` | string(`assistant`) | 모델 응답은 항상 `assistant` |
| `message.content` | string | 어시스턴트 메시지 텍스트 |
| `message.thinking` | string | `think` 활성화 시의 사고 과정 추적 |
| `message.tool_calls` | `ToolCall[]` | 어시스턴트가 요청한 도구 호출 |
| `message.images` | string[] | 응답에 포함된 base64 이미지(선택) |
| `done` | boolean | 채팅 응답 완료 여부 |
| `done_reason` | string | 응답이 종료된 이유 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 프롬프트의 토큰 수 |
| `prompt_eval_duration` | integer | 프롬프트 평가 시간(나노초) |
| `eval_count` | integer | 응답으로 생성된 토큰 수 |
| `eval_duration` | integer | 토큰 생성 시간(나노초) |
| `logprobs` | `Logprob[]` | `logprobs` 활성화 시 생성 토큰의 로그 확률 정보 |

#### `ChatStreamEvent`

채팅 스트리밍 모드에서 전송되는 부분 이벤트입니다. `message`의 각 필드는 해당 청크 분량만 담으며,
스트림의 마지막 이벤트에서 `done`이 `true`가 됩니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 이 스트림 이벤트에 사용된 모델 이름 |
| `created_at` | string(date-time) | 이 청크 생성 시각(ISO 8601) |
| `message` | object | 부분 메시지 객체(아래 하위 필드 참조) |
| `message.role` | string | 이 청크 메시지의 역할 |
| `message.content` | string | 부분 어시스턴트 메시지 텍스트 |
| `message.thinking` | string | `think` 활성화 시의 부분 사고 텍스트 |
| `message.tool_calls` | `ToolCall[]` | 부분 도구 호출(있는 경우) |
| `message.images` | string[] | 부분 base64 이미지(있는 경우) |
| `done` | boolean | 스트림의 마지막 이벤트면 `true` |

### 임베딩 (Embed)

#### `EmbedRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 모델 이름 |
| `input` | string \| string[] | ✅ | 임베딩을 생성할 텍스트 또는 텍스트 배열 |
| `truncate` | boolean |  | true면 컨텍스트 창을 초과한 입력을 잘라냄. false면 오류 반환. 기본값 `true` |
| `dimensions` | integer |  | 생성할 임베딩 차원 수 |
| `keep_alive` | string |  | 모델 메모리 유지 시간 |
| `options` | `ModelOptions` |  | 생성을 제어하는 런타임 옵션 |

#### `EmbedResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `model` | string | 임베딩을 생성한 모델 |
| `embeddings` | number[][] | 벡터 임베딩 배열 |
| `total_duration` | integer | 총 생성 시간(나노초) |
| `load_duration` | integer | 모델 로드 시간(나노초) |
| `prompt_eval_count` | integer | 임베딩 생성을 위해 처리한 입력 토큰 수 |

### 모델 관리 (Create / Copy / Delete / Pull / Push)

#### `CreateRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 생성할 모델 이름 |
| `from` | string |  | 기반으로 삼을 기존 모델 |
| `template` | string |  | 모델이 사용할 프롬프트 템플릿 |
| `license` | string \| string[] |  | 모델 라이선스 문자열 또는 목록 |
| `system` | string |  | 모델에 내장할 시스템 프롬프트 |
| `parameters` | object |  | 모델의 키-값 파라미터 |
| `messages` | `ChatMessage[]` |  | 모델에 사용할 메시지 기록 |
| `quantize` | string |  | 적용할 양자화 수준(예: `q4_K_M`, `q8_0`) |
| `stream` | boolean |  | 상태 업데이트 스트리밍 여부. 기본값 `true` |

#### `CopyRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `source` | string | ✅ | 복사할 원본 모델 이름 |
| `destination` | string | ✅ | 생성할 새 모델 이름 |

#### `DeleteRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 삭제할 모델 이름 |

#### `PullRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 내려받을 모델 이름 |
| `insecure` | boolean |  | 비보안 연결을 통한 다운로드 허용 |
| `stream` | boolean |  | 진행 상태 스트리밍 여부. 기본값 `true` |

#### `PushRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 게시할 모델 이름 |
| `insecure` | boolean |  | 비보안 연결을 통한 게시 허용 |
| `stream` | boolean |  | 진행 상태 스트리밍 여부. 기본값 `true` |

### 모델 조회 (Show / Tags / Ps)

#### `ShowRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `model` | string | ✅ | 조회할 모델 이름 |
| `verbose` | boolean |  | true면 응답에 대용량 상세 필드 포함 |

#### `ShowResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `parameters` | string | 텍스트로 직렬화된 모델 파라미터 설정 |
| `license` | string | 모델 라이선스 |
| `modified_at` | string | 최종 수정 시각(ISO 8601) |
| `details` | object | 모델 상위 수준 상세 정보 |
| `template` | string | 모델이 프롬프트 렌더링에 사용하는 템플릿 |
| `capabilities` | string[] | 지원 기능 목록 |
| `model_info` | object | 추가 모델 메타데이터 |

#### `ModelSummary`

로컬에서 사용 가능한 모델의 요약 정보입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `name` | string | 모델 이름 |
| `model` | string | 모델 이름 |
| `remote_model` | string | 원격 모델인 경우 상위(업스트림) 모델 이름 |
| `remote_host` | string | 원격 모델인 경우 업스트림 Ollama 호스트 URL |
| `modified_at` | string | 최종 수정 시각(ISO 8601) |
| `size` | integer | 디스크상 모델 총 크기(바이트) |
| `digest` | string | 모델 콘텐츠의 SHA256 다이제스트 식별자 |
| `details` | object | 모델 형식·계열 상세 정보(아래 하위 필드 참조) |
| `details.format` | string | 모델 파일 형식(예: `gguf`) |
| `details.family` | string | 주 모델 계열(예: `llama`) |
| `details.families` | string[] | 모델이 속한 모든 계열(해당 시) |
| `details.parameter_size` | string | 대략적 파라미터 수 레이블(예: `7B`, `13B`) |
| `details.quantization_level` | string | 사용된 양자화 수준(예: `Q4_0`) |

#### `ListResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `models` | `ModelSummary[]` | 사용 가능한 모델 목록 |

#### `Ps`

실행 중인(메모리 적재) 모델 하나의 정보입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `name` | string | 실행 중인 모델 이름 |
| `model` | string | 실행 중인 모델 이름 |
| `size` | integer | 모델 크기(바이트) |
| `digest` | string | 모델의 SHA256 다이제스트 |
| `details` | object | 형식·계열 등 모델 상세 정보 |
| `expires_at` | string | 모델이 메모리에서 해제될 예정 시각 |
| `size_vram` | integer | VRAM 사용량(바이트) |
| `context_length` | integer | 실행 중인 모델의 컨텍스트 길이 |

#### `PsResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `models` | `Ps[]` | 현재 실행 중인 모델 목록 |

### 상태·버전 (Status / Version)

#### `StatusEvent`

스트리밍 작업(create/pull/push)의 진행 상태 이벤트입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `status` | string | 사람이 읽을 수 있는 상태 메시지 |
| `digest` | string | 해당 상태와 연관된 콘텐츠 다이제스트(있는 경우) |
| `total` | integer | 작업에 예상되는 총 바이트 수 |
| `completed` | integer | 현재까지 전송된 바이트 수 |

#### `StatusResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `status` | string | 현재 상태 메시지 |

#### `VersionResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `version` | string | Ollama 버전 |

### 로그 확률 (Logprobs)

#### `Logprob`

생성된 토큰 하나에 대한 로그 확률 정보입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `token` | string | 토큰의 텍스트 표현 |
| `logprob` | number | 이 토큰의 로그 확률 |
| `bytes` | integer[] | 토큰의 원시 바이트 표현 |
| `top_logprobs` | `TokenLogprob[]` | 이 위치에서 가장 가능성 높은 토큰들과 각 로그 확률 |

#### `TokenLogprob`

토큰 후보 하나에 대한 로그 확률 정보입니다.

| 필드 | 타입 | 설명 |
|---|---|---|
| `token` | string | 토큰의 텍스트 표현 |
| `logprob` | number | 이 토큰의 로그 확률 |
| `bytes` | integer[] | 토큰의 원시 바이트 표현 |

### 웹 검색·페치 (Web Search / Fetch)

> 아래 스키마는 명세의 `components/schemas`에 정의돼 있으나, 이 마스터 명세의 `paths`에는
> 대응하는 엔드포인트가 포함돼 있지 않습니다.

#### `WebSearchRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `query` | string | ✅ | 검색 질의 문자열 |
| `max_results` | integer |  | 반환할 최대 결과 수(1~10, 기본값 `5`) |

#### `WebSearchResult`

| 필드 | 타입 | 설명 |
|---|---|---|
| `title` | string | 결과 페이지 제목 |
| `url` | string(uri) | 결과의 확정 URL |
| `content` | string | 추출된 텍스트 콘텐츠 일부 |

#### `WebSearchResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `results` | `WebSearchResult[]` | 일치하는 검색 결과 배열 |

#### `WebFetchRequest`

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `url` | string(uri) | ✅ | 가져올 URL |

#### `WebFetchResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `title` | string | 가져온 페이지 제목 |
| `content` | string | 추출된 페이지 콘텐츠 |
| `links` | string(uri)[] | 페이지에서 발견된 링크 목록 |

### 오류 (Error)

#### `ErrorResponse`

| 필드 | 타입 | 설명 |
|---|---|---|
| `error` | string | 무엇이 잘못되었는지 설명하는 오류 메시지 |

> 원문: https://docs.ollama.com/openapi.yaml

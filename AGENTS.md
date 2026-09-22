# OmniRelay — Agent 작업 지침

제품 개요·API 사용법은 `README.md` 참고. 이 파일에는 README에 없거나 README가 틀린, 세션에서 반드시 알아야 할 내용만 담는다.

## 명령어 (cwd 주의)

`backend/`에서:

```bash
go run ./cmd/server/              # 개발 서버 (:8080)
go vet ./... && go test ./...     # 검증. 단일 테스트: go test ./internal/proxy/ -run TestXyz
```

`frontend/`에서 (bun 전용 — npm/yarn/pnpm 금지):

```bash
bun install                       # 의존성 설치 + bun.lock 갱신
bun run build                     # 유일한 검증: vue-tsc --noEmit 타입체크 후 vite build
bun run dev                       # :5173
```

- 테스트 러너·lint·format 스크립트가 없다 (prettier는 의존성에만 있고 실행 스크립트·설정 파일 없음). 프론트 검증은 `bun run build`뿐.
- `Dockerfile`는 `bun install --frozen-lockfile`을 쓴다 — `package.json` 변경 후 `bun install`로 lockfile을 갱신하지 않으면 이미지 빌드가 실패한다.
- 완료 전: 백엔드 `go vet ./... && go test ./...`, 프론트 `bun run build` 둘 다 실제 통과 출력을 확인.
- CI(`.github/workflows/docker.yml`)는 Docker 이미지만 빌드하고 테스트·린트를 돌리지 않는다 — 로컬 검증이 유일한 방어선이다.

## README 불일치 (코드가 정답)

README는 패스스루(URL 릴레이) 기능을 아직 상세히 설명하지만 commit `13c5423`에서 제거됐다. 아래는 전부 존재하지 않는다:

- `internal/passthrough/`, `PassthroughView.vue`, `passthrough_logs` 테이블, `PASSTHROUGH_*` 환경변수
- `OpenAPI-Specification/` 디렉터리
- README의 마이그레이션 버전 표기(v9/v14 — 실제 최신은 v17), 뷰 개수 표기

실제 환경변수는 `backend/internal/config/config.go`에 있는 `LISTEN_ADDR`, `DATABASE_PATH`, `JWT_SECRET`, `ENCRYPT_KEY`, `CORS_ORIGINS` 다섯 개뿐. 설정·동작은 항상 코드/설정 파일을 신뢰할 것.

## 아키텍처 배선 (파일명으로 못 찾는 사실)

- 단일 컨테이너: Bun 빌드 → Go 정적 바이너리 → `caddy:2-alpine`. `Dockerfile` ENTRYPOINT가 Go 서버와 Caddy를 셸에서 함께 띄운다.
- 새 API 경로는 **세 곳**에 등록해야 완전하다:
  1. `backend/cmd/server/main.go` — 실제 Gin 라우트
  2. `frontend/Caddyfile` — `:80`에서 `:8080`으로 프록시할 패턴. 빠지면 SPA fallback이 응답한다.
  3. `frontend/vite.config.ts` — 개발 프록시. 현재 `/v1`, `/admin`만 프록시하므로, **path-routed 요청(`/openai/v1/...`)은 개발 중 `:5173`이 아니라 `:8080`으로 직접 호출**해야 한다.
- `frontend/Caddyfile` 규칙:
  - 매처는 `path`/`path_regexp`만 사용. `path_prefix` 매처는 stock `caddy:2-alpine`에 없어 빌드/기동이 죽는다.
  - 새로 추가하는 reverse_proxy 블록에는 반드시 `flush_interval -1` — 없으면 SSE 스트리밍이 버퍼링된다.

## 데이터베이스 마이그레이션

`backend/internal/database/migrations.go` — 시작 시 자동 실행, append-only.

- 컬럼 추가는 반드시 멱등하게: SQLite에 `ADD COLUMN IF NOT EXISTS`가 없으므로 `hasColumn()`(PRAGMA table_info)로 감싼다. v1부터 쓰인 테이블의 새 컬럼도 새 마이그레이션이 필요하다.
- 버전 번호를 재사용·재정렬하지 말 것. 목록에 v15·v16이 없는 것은 의도적(구 passthrough 마이그레이션 제거). 이미 찍힌 DB가 새 구조를 건너뛰지 않도록 v14와 v17이 같은 `ensureProviderAPIKeys`를 공유한다. 새 마이그레이션은 항상 맨 뒤에 현재 최대 버전 + 1로 추가한다.

## 운영·보안

- `GIN_MODE=release`에서 `JWT_SECRET`/`ENCRYPT_KEY`가 기본값이면 기동 즉시 Fatal이다 (`config.go`). 프로덕션 배포 시 `ENCRYPT_KEY`는 64자 hex(`openssl rand -hex 32`).
- API 키 인증은 `Authorization: Bearer om-ni-...`와 `x-api-key` 헤더 둘 다 받는다 (`middleware/apikey_auth.go`).

## 프록시/스트리밍 (가장 오류가 잦은 영역)

- OpenAI 계열(openai/lmstudio/ollama) 스트리밍은 사용량을 받으려면 요청 본문에 `stream_options: {include_usage: true}` 주입이 필요하다. `isOpenAICompat()`이 게이트이며 **주입 지점은 셋 다** 챙겨야 한다: `chat_handler.go`의 `executeChat`, `proxy.go`의 `executeMessages`, `proxy.go`의 `handlePathRoutedProxy`.
- 캐시 토큰 필드가 제공자마다 다르다 — Anthropic `cache_creation_input_tokens`/`cache_read_input_tokens`, OpenAI `prompt_tokens_details.cached_tokens`, Gemini `usageMetadata.cached_content_token_count`. 추출은 `extractCacheTokens()`(`cache.go`)와 비스트리밍 종합 `extractUsageFromRawResponse()`(`cost.go` — 흔히 오해하는 `http_helpers.go`가 아니다). provider 분기마다 캐시 토큰을 빠뜨리면 과거에 실제로 누락 버그가 난다.
- 스트림 핸들러는 네 개다: `handleStreamResponse`(OpenAI SSE), `handleMessagesStreamResponse`(Anthropic SSE), `handleResponsesStream`(`/v1/responses` 전용), `handleRawStreamResponse`(미확인 포맷 — 지연만 기록). path-routed 요청도 알려진 어댑터면 앞의 둘 중 하나로 dispatch해야 토큰·비용이 기록된다.
- 다중 업스트림 키 failover: `tryKeys`(`upstream.go`)가 클라이언트가 아무 바이트도 받기 전에 네트워크 오류/401/403/429/5xx에서 다음 키로 재시도하고, 401/403은 해당 키를 자동 비활성화한다. 라운드로빈 커서는 인메모리다 (멀티프로세스가 되면 저장 필요 — 코드에 `// ponytail:` 주석). `usage_logs`에는 최종 시도 1건만 기록한다.
- 모델 추가 시 가격을 `https://models.dev/api.json`에서 자동 채운다 (`service/modelsdev.go`, 실행 중 네트워크 fetch·10s 타임아웃). 실패해도 가격 0으로 진행하니 장애가 아닌 것으로 오인하지 말 것.

## 프론트엔드

- 새 UI 문자열은 `src/locales/{en,ja,ko}.ts` 세 곳에 동시에 추가한다 (i18n 누락이 기본 상태).
- Pinia 스토어: 소문자 파일명 + `useXStore` export (예: `stores/providers.ts` → `useProvidersStore`), 뷰는 `PascalCaseView.vue`.

## 제약

- 셸에서 `rg` 사용 금지.
- 사용자가 명시적으로 요청하지 않는 한 커밋 금지.
- 큰 변경 전 `docs/superpowers/specs/`·`docs/superpowers/plans/`에 기존 설계/계획이 있는지 확인하고, 해당 plan 파일의 Global Constraints를 지킬 것.
- 테스트는 구현 옆에 `*_test.go`로. httptest·임시 SQLite를 쓰므로 외부 서비스가 필요 없다.

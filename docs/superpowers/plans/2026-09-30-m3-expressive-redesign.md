# OmniRelay M3 Expressive 리디자인 — 구현 계획

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** OmniRelay 프론트엔드를 Vuetify 없이 자체 Material 3 Expressive 디자인 시스템으로 전환하고, 라이트/다크 듀얼 테마와 형태 모핑·스프링 모션을 전 화면에 적용한다.

**Architecture:** `src/styles/tokens/`가 공식 M3/M3E 토큰(색 role·타입·형태·모션·엘리베이션)의 유일한 원천이고, `src/components/ui/`의 `Or*` 프리미티브가それを 소비한다. Vuetify는 마지막 단계에서 제거된다 — 새 시스템이 전 화면을 커버한 뒤에 빠진다.

**Tech Stack:** Vue 3 · Vite 5 · Pinia · vue-i18n(en/ja/ko) · chart.js/vue-chartjs · Roboto(@fontsource) · Material Symbols

**Spec:** `docs/superpowers/specs/2026-09-30-m3-expressive-redesign-design.md`

## Global Constraints

- 백엔드 API·데이터 모델·라우트·스토어·i18n 키 구조 **변경 금지** (새 문자열 키 추가만 허용).
- 새 UI 문자열은 `src/locales/en.ts`, `src/locales/ja.ts`, `src/locales/ko.ts` **세 곳 동시** 추가.
- 프론트 유일한 검증은 `cd frontend && bun run build` (`vue-tsc --noEmit` + `vite build`).
- `package.json` 변경 후 반드시 `cd frontend && bun install`로 `bun.lock` 갱신 (`Dockerfile`이 `--frozen-lockfile`).
- 패키지 매니저는 **bun만**. npm/yarn/pnpm 금지.
- 커밋 금지 (사용자 명시 요청 없이).
- 셸에서 `rg` 사용 금지.
- 하드코딩 색/라디우스/지속시간은 `src/styles/tokens/` 외 금지.
- `frontend/Caddyfile`와 `frontend/vite.config.ts`는 수정하지 않는다 (새 API 경로 없음).

## Review Focus

1. **테마 FOUC** — 첫 페인트 전에 `data-theme`가 없으면 밝은 배경이 번쩍인다. `index.html` 인라인 스크립트가 `localStorage` → `prefers-color-scheme` 순으로 즉시 설정해야 한다. (Task 2)
2. **듀얼 테마 누락 색** — 다크 블록에 없는 role이 있으면 라이트 hex가 다크에서 튄다. 두 블록의 키 집합이 정확히 같은지 검산해야 한다. (Task 2)
3. **다이얼로그 접근성** — 포커스트랩 이탈, ESC 미동작, 스크림 클릭 시 상태 불일치, 스크롤 잠금 해제 누락. (Task 6)
4. **`OrDataTable` 슬롯 호환** — 5개 뷰의 `#item.<key>` 40여 개와 `#no-data`가 그대로 동작하지 않으면 화면이 깨진다. (Task 7)
5. **i18n 누락** — 테마 토글/새 라벨이 en에만 추가되고 ja/ko가 빠지면 기본 상태가 깨진다. (Task 3)
6. **차트 테마 전환** — 테마를 바꿔도 차트 색이 갱신되지 않으면 라이트에서 가독성이 죽는다. (Task 9)
7. **폰트/CJK** — Roboto에 없는 한글·일본어 글리프가 깨지지 않도록 스택 폴백이 있어야 한다. (Task 1)

---

### Task 1: 스프링 생성기와 모션 토큰

**Files:**
- Create: `frontend/scripts/gen-springs.mjs`
- Create: `frontend/src/styles/tokens/motion-springs.css` (생성 산출물, 커밋)
- Create: `frontend/src/styles/tokens/motion.css`

**Interfaces:**
- Produces: `--m3-motion-spatial-fast` / `-spatial-default` / `-spatial-slow` / `-effects-fast` / `-effects-default` / `-effects-slow` — 값은 `linear(...)` 곡선.
- Produces: `--m3-duration-*` (16), `--m3-easing-*` (6), `--or-motion-entrance|settle|press|route-enter|route-exit`

- [ ] **Step 1: 스프링 생성기 작성**

`gen-springs.mjs`는 질량 1 스프링 `x(t) = 1 - e^{-rωt}(cos ω_d t + (rω/ω_d) sin ω_d t)`
(ω=√k, ω_d=ω√(1-r²))을 응답이 `|x-1| < 0.001`이 될 때까지 샘플링해
`linear(0, ..., 1)` 문자열을 출력한다. 임계 감쇠(r=1)는 `x = 1-(1+ωt)e^{-ωt}`.

표 (Expressive 스킴 — 기본값):

| 이름 | k | r |
|---|---|---|
| spatial-fast | 800 | 0.6 |
| spatial-default | 380 | 0.8 |
| spatial-slow | 200 | 0.8 |
| effects-fast | 3800 | 1 |
| effects-default | 1600 | 1 |
| effects-slow | 800 | 1 |

Standard 스킴(k/r: 1400/0.9, 700/0.9, 300/0.9, 나머지 동일)도 같은 방식으로 출력하되
`-standard` 접미사를 붙인다. 각 샘플 수는 최소 24개, 소수점 4자리.

- [ ] **Step 2: 생성기를 실행하고 산출물을 검증**

Run: `cd frontend && node scripts/gen-springs.mjs`
Expected: `src/styles/tokens/motion-springs.css` 생성, 각 `--m3-motion-spatial-*` 값이 `linear(`로 시작하고 마지막 인자가 `1`로 끝남. 오버슈트가 있는 spatial 값은 중간에 `1`보다 큰 숫자를 포함해야 한다 (r<1 확인).

- [ ] **Step 3: `motion.css` 작성**

공식 M3 duration 16종·easing 6종(스펙 3.7 표 그대로) + `--or-motion-*` 5종을 정의하고,
`motion-springs.css`를 import한다. `linear()` 미지원 대비 `@supports not (transition-timing-function: linear(0,1))` 블록에서 각 스프링을 `cubic-bezier` 폴백으로 재정의한다. 파일 끝에 reduced-motion 게이트(스펙 3.7 코드 블록 그대로)를 둔다.

- [ ] **Step 4: 검증**

Run: `cd frontend && bun run build`
Expected: PASS (이 단계에서는 아직 import되지 않아도 빌드는 통과해야 한다)

---

### Task 2: 컬러·타입·형태·엘리베이션·간격·상태 토큰 + 테마 인프라

**Files:**
- Create: `frontend/src/styles/tokens/{palette,color,type,shape,elevation,space,state,index}.css`
- Create: `frontend/src/styles/base.css`
- Create: `frontend/src/stores/theme.ts`
- Modify: `frontend/index.html` (인라인 테마 스크립트 + font preconnect)
- Modify: `frontend/src/main.ts` (`tokens.css` → `tokens/index.css`, `base.css` 추가)
- Delete: `frontend/src/styles/tokens.css`

**Interfaces:**
- Produces: `--m3-color-<role>` 49종 × 2테마, `--m3-palette-*`, `--m3-typescale-*`/`--m3-typescale-emphasized-*` 각 15, `--m3-shape-*` 10, `--m3-elevation-*`, `--m3-space-*`, `--m3-state-*`
- Produces: `useThemeStore()` → `{ mode: Ref<"light"|"dark">, resolved: ComputedRef<"light"|"dark">, toggle(): void }`
- Consumes: Task 1의 `--m3-motion-*`

- [ ] **Step 1: `palette.css` 작성**

시드 `#6750A4`에서 파생된 M3 baseline palette 전 톤을 `--m3-palette-primary-0..100`,
`-secondary-*`, `-tertiary-*`, `-error-*`, `-neutral-*`(0,4,6,10,12,17,20,22,24,30,40,50,60,70,80,87,90,92,94,95,96,98,99,100),
`-neutral-variant-*`으로 정의한다. hex 값은 스펙 3.3 표의 주석에 표기된 원소를 따른다.

- [ ] **Step 2: `color.css` 작성**

`:root[data-theme="light"]`와 `:root[data-theme="dark"]` 두 블록에 role 49종을 정의하고
전부 `var(--m3-palette-*)`에 위임한다. **두 블록의 키 집합은 정확히 일치**해야 한다.
이어서 앱 별칭 `--or-success*`(3), `--or-warning*`(3)를 스펙 3.4 표 값으로 정의한다.

- [ ] **Step 3: `type.css`, `shape.css`, `elevation.css`, `space.css`, `state.css` 작성**

스펙 3.5 / 3.6 / 3.8 / 3.9 표를 그대로 CSS 변수로 옮긴다. `type.css`에는 서체 스택
(`--m3-typeface-plain`, `--m3-typeface-mono`)과 `.m3-typescale-<role>` 헬퍼 클래스도 정의한다.

- [ ] **Step 4: `index.css`와 `base.css` 작성**

`index.css`는 9개 토큰 파일을 순서대로 import하는 유일 진입점.
`base.css`는 reset(`box-sizing`, `margin:0`, `font: inherit`) + `html/body`가
`--m3-color-background`/`--m3-color-on-background`/`--m3-typeface-plain`을 쓰고
`color-scheme`을 테마에 맞게 설정 + `:focus-visible` 아웃라인 규칙 + `::selection`을 정의한다.

- [ ] **Step 5: 테마 스토어와 `index.html` 인라인 스크립트 작성**

`stores/theme.ts`: `mode`는 `"system" | "light" | "dark"`를 보유하고
`localStorage["omnirelay.theme"]`(없으면 `"system"`)에서 초기화한다. `resolved`는
`mode==="system"`이면 `matchMedia("(prefers-color-scheme: dark)")` 결과, 아니면 `mode`.
`resolved`가 변할 때 `document.documentElement.dataset.theme = resolved`를 쓴다.
`toggle()`은 light↔dark를 왕복한다(시스템 선택은 별도 `setMode`).

`index.html` `<head>` 최상단 인라인 스크립트: `localStorage` → `matchMedia` 순으로 값을
읽어 `document.documentElement.dataset.theme`를 즉시 설정한다. FOUC 방지가 목적.

- [ ] **Step 6: `main.ts` import 교체 및 빌드 검증**

`import "./styles/tokens.css"` → `import "./styles/tokens/index.css"` + `import "./styles/base.css"`.
`styles/tokens.css` 삭제.
Run: `cd frontend && bun run build`
Expected: PASS

- [ ] **Step 7: 드얼 테마 키 집합 검산**

Run: `grep -c '^  --m3-color-' frontend/src/styles/tokens/color.css`
Expected: 라이트 블록 키 수와 다크 블록 키 수가 정확히 같음 (둘 다 49).
차이가 있으면 한쪽 블록을 보완한 뒤 다시 실행.

---

### Task 3: 아이콘·폰트 의존성과 라우트/엔트리 정돈

**Files:**
- Modify: `frontend/package.json` (add `@fontsource/roboto`, `@fontsource-variable/material-symbols`; del `vuetify`, `@mdi/font`)
- Create: `frontend/src/plugins/components.ts` (전역 컴포넌트 등록)
- Create: `frontend/src/components/ui/OrIcon.vue`
- Modify: `frontend/src/main.ts` (plugins 교체, 폰트 import)
- Delete: `frontend/src/plugins/vuetify.ts`
- Modify: `frontend/src/locales/{en,ja,ko}.ts` (새 키)

**Interfaces:**
- Produces: `<OrIcon name="search" :size="20" :weight="400" :fill="false" />` — ligature 렌더
- Produces: `useThemeStore` 는 Task 2

- [ ] **Step 1: `bun add` / `bun remove` 실행 후 lockfile 갱신 확인**

Run: `cd frontend && bun add @fontsource/roboto @fontsource-variable/material-symbols && bun remove vuetify @mdi/font && bun install`
Expected: `bun.lock` 변경, `node_modules/vuetify` 없음.

- [ ] **Step 2: `OrIcon.vue` 작성**

Material Symbols 변수 폰트 ligature를 렌더하는 얇은 래퍼.
props: `name: string`, `size?: number|string` (기본 24), `weight?: 100..700` (기본 400),
`fill?: boolean` (기본 false), `color?: string`.
`aria-hidden` 기본 true, `name`이 접근성 정보일 때만 `role="img"` + `aria-label`.

- [ ] **Step 3: `plugins/components.ts`와 `main.ts` 정돈**

`plugins/vuetify.ts` 삭제. `main.ts`는 `components`(지금은 `OrIcon`만 등록, 이후 단계에서 추가)와
`i18n`만 사용한다. `@fontsource/roboto/400.css`·`500.css`·`700.css`, `@fontsource/roboto-mono` 대체로
`@fontsource/roboto` 내 mono는 없으므로 `--m3-typeface-mono`는 시스템 모노에 위임한다(폰트 패키지 추가는 이 스텝에서 결정된 것대로).

- [ ] **Step 4: locales 세 곳에 키 추가**

`common.theme.light` / `common.theme.dark` / `common.theme.system`을 en·ja·ko에 동시 추가한다.

- [ ] **Step 5: 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 4: 핵심 프리미티브 (버튼·필드·카드·칩)

**Files:**
- Create: `frontend/src/components/ui/OrButton.vue`, `OrIconButton.vue`, `OrCard.vue`,
  `OrChip.vue`, `OrSegmentedButton.vue`, `OrTextField.vue`, `OrSelect.vue`,
  `OrCheckbox.vue`, `OrSwitch.vue`, `OrSpinner.vue`, `OrProgressBar.vue`, `OrDivider.vue`
- Create: `frontend/src/composables/useRipple.ts`
- Modify: `frontend/src/plugins/components.ts` (전체 등록)

**Interfaces:**
- `OrButton`: props `variant?: "filled"|"tonal"|"outlined"|"text"|"elevated"` (기본 `filled`),
  `tone?: "primary"|"danger"` (기본 `primary`), `size?: "sm"|"md"|"lg"` (기본 `md`),
  `loading?: boolean`, `disabled?: boolean`, `type?`, `block?: boolean`
  emits `click`. 슬롯: default, `leading`, `trailing`.
  **로딩 중에는 disabled + spinner로 label을 대체하고, 너비는 고정한다(`min-width` 유지).**
- `OrIconButton`: props `icon: string`, `label: string`(필수, `aria-label`), `tone?`,
  `variant?: "standard"|"filled"|"tonal"|"outlined"`, `size?`, `disabled?`
- `OrCard`: props `variant?: "elevated"|"filled"|"outlined"` (기본 `elevated` = tonal),
  `interactive?: boolean` (pressed 시 `corner-2xl → lg` 모핑)
- `OrChip`: props `variant?: "assist"|"filter"|"input"`, `selected?: boolean`,
  `removable?: boolean`, `disabled?`, `icon?`
- `OrSegmentedButton`: props `options: {value: string; label?: string; icon?: string}[]`,
  `modelValue: string`, `multiple?: boolean`; emits `update:modelValue`
- `OrTextField`: props `modelValue`, `label`, `hint?`, `error?`, `type?`(기본 `"text"`),
  `disabled?`, `readonly?`, `leadingIcon?`, `clearable?`; emits `update:modelValue`
  구조: `<label>`+`<input>`+플로팅 라벨(span, `:empty` 가에서 위로 이동)+hint/error 라인.
  M3 outlined: `border-radius`가 focus 시 `sm → xs` 모핑.
- `OrSelect`: 동일 props에 `options: {value:string; label:string}[]`. 네이티브 `<select>` + M3 셸.
- `OrCheckbox`: props `modelValue: boolean|unknown[]`, `value?`, `label?`, `hint?`; emits `update:modelValue`
- `OrSwitch`: props `modelValue: boolean`, `label?`, `disabled?`; emits `update:modelValue`
- `OrSpinner`: props `size?` (기본 20), `label?`
- `OrProgressBar`: props `value?: number`(0..1, 미정의면 indeterminate), `tone?`
- `useRipple`: `useRipple(el)` → `bind(el)` 이벤트 핸들러. `pointerdown` 시 원형 state layer를
  `--m3-state-press` 불투명도로 생성하고 `--or-motion-press` 동안 확산 후 제거.
  `prefers-reduced-motion`이면 ripple를 생략하고 state layer만 남긴다.

- [ ] **Step 1: `useRipple`과 `OrIcon` 외 12개 컴포넌트를 토큰만으로 작성**

하드코딩 색·라디우스·지속시간 금지 — 전부 `var(--m3-color-*)`, `var(--m3-shape-*)`,
`var(--m3-motion-*)`, `var(--m3-space-*)`, `var(--m3-typescale-*)`만 사용.
모든 인터랙티브 요소는 최소 48×48 탭 타깃(시각 크기는 작아도 hit area 보장).

- [ ] **Step 2: `plugins/components.ts`에 전량 등록**

- [ ] **Step 3: 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 5: 톤다운 재작성 — `page-shared.css` / `auth.css` / 앱 컴포넌트

**Files:**
- Rewrite: `frontend/src/styles/page-shared.css` (토큰 기반, M3 역할로)
- Rewrite: `frontend/src/styles/auth.css`
- Rewrite: `frontend/src/components/{AppAlert,AuthShell,EmptyState,MobileDataCard,MonoTag,PageHeader,StatCard,StatusChip}.vue`
  — **prop/slot 시그니처는 불변**, 내부 마크업·스타일만 교체.

**Interfaces:**
- Consumes: Task 2 토큰, Task 4 프리미티브
- Produces: `.page`, `.page-header`, `.filter-bar`, `.filter-col`, `.table-card`, `.field-group`,
  `.mono-tag`, `.status-chip--{on,off,warning,neutral}`, `.alert--{success,error,warning,page}`,
  `.empty-state`, `.btn-spinner` 등 기존 클래스명 **유지** (호출 7개 뷰가 그대로 참조)

- [ ] **Step 1: `page-shared.css`를 토큰 기반으로 재작성**

클래스명과 선택자 구조는 유지하고 값만 교체한다. `.btn-primary` 등은 내부를
`OrButton`의 시각 규칙과 동일하게 맞춘다(이후 뷰 마이그레이션에서 실제 컴포넌트로 교체).

- [ ] **Step 2: `auth.css` 재작성**

M3 카드·outlined 필드 규칙으로. 배경 그리드/글로우 제거, `AuthShell` 마크업에서
`.auth-grid`/`.auth-glow` div도 함께 제거.

- [ ] **Step 3: 앱 컴포넌트 8개 재작성**

`<v-icon>`은 `<OrIcon>`으로 교체. `StatusChip` variants는 `--or-success/-warning`과
M3 `error` role을 쓴다.

- [ ] **Step 4: 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 6: `OrDialog`, `OrSnackbar`, 포커스트랩

**Files:**
- Create: `frontend/src/components/ui/OrDialog.vue`, `OrSnackbar.vue`
- Create: `frontend/src/composables/useFocusTrap.ts`
- Create: `frontend/src/stores/snackbar.ts`
- Modify: `frontend/src/plugins/components.ts`

**Interfaces:**
- `OrDialog`: v-model `boolean`, props `width?: number|string` (기본 480),
  `fullscreen?: boolean`, `dismissible?` (기본 true), `title?: string`.
  슬롯: default, `header`, `footer`. emits `update:modelValue`.
  접근성: `role="dialog"`, `aria-modal="true"`, 라벨은 `title` 또는 `aria-label`,
  열림 시 첫 포커서블로 이동, ESC 닫기, 닫힘 시 트리거로 복귀, 스크림 클릭 시 닫기,
  `document.body` 스크롤 잠금(중첩 대응: 카운터), `Teleport to="body"`.
  모션: 스크림 fade는 `--m3-motion-effects-fast`, 컨테이너는
  scale 0.92→1 + `corner-4xl → 2xl` 모핑, `--m3-motion-spatial-fast`.
- `OrSnackbar`: props `message`, `variant?: "info"|"success"|"error"`,
  `timeout?` (기본 4000, 0이면 영구), `action?: {label,onClick}`.
  `stores/snackbar.ts`의 `useSnackbarStore()`가 `show(message, opts)`를 노출.
  화면 최하단 중앙(모바일은 하단 네비 위 88px), `role="status" aria-live="polite"`.
  슬롯을 왼쪽에서 `spatial-fast`로 등장.

- [ ] **Step 1: `useFocusTrap` 작성**

`useFocusTrap(container, active)` → Tab/Shift+Tab 순환, 포커서블 목록 재계산, 외부 포커스 강제 복귀.

- [ ] **Step 2: `OrDialog` 작성** — 위 접근성 요구 전부 구현.

- [ ] **Step 3: `OrSnackbar` + 스토어 작성**, 등록, locales 세 곳 `common.close` 키 추가.

- [ ] **Step 4: 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 7: `OrDataTable`

**Files:**
- Create: `frontend/src/components/ui/OrDataTable.vue`
- Modify: `frontend/src/plugins/components.ts`

**Interfaces:** 스펙 4.2 계약을 그대로 따른다. 추가로:
- `density="compact"`: 셀 패딩 `8px 16px`, `comfortable`: `12px 16px`, 행 높이 최소 48/56px.
- `loading` true: tbody 오버레이에 `OrSpinner`, `aria-busy="true"`.
- 머리글 셀: `sortable`이면 `<button>` 내부 + `aria-sort` + 아이콘 방향 표시.
- 빈 상태: `#no-data` 슬롯 없이 `items.length === 0`이면 기본 메시지.

- [ ] **Step 1: `OrDataTable` 작성** (스펙 4.2 + 위 밀도/정렬/로딩 규칙)

- [ ] **Step 2: 등록 및 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 8: 레이아웃/IA — NavRail, NavigationBar, DefaultLayout

**Files:**
- Create: `frontend/src/components/ui/OrNavRail.vue`, `OrNavigationBar.vue`
- Rewrite: `frontend/src/layouts/DefaultLayout.vue`
- Modify: `frontend/src/App.vue` (제거되는 `v-app` 교체)
- Modify: `frontend/src/composables/useMobile.ts` (768 → 600px)
- Modify: `frontend/src/locales/{en,ja,ko}.ts`

**Interfaces:**
- `OrNavRail`: props `expanded: boolean`, `items: {to:string; icon:string; label:string; adminOnly?:boolean}[]`,
  emits `toggle`. expanded(256px) ↔ rail(80px) 전환, 활성 항목은 M3 내비 규칙.
- `OrNavigationBar`: props `items` (모바일 하단 80px + safe-area, 활성 pill).
- `DefaultLayout`: `expanded`는 viewport ≥1100px에서만 true. 600–1099px은 rail,
  <600px은 `OrNavigationBar`. 하단에 테마 토글(3분할 segmented: system/light/dark)과
  로케일(EN/JA/KO)·로그아웃.
- `App.vue`: `v-app` → plain `<div id="app">`.

- [ ] **Step 1: `OrNavRail` / `OrNavigationBar` 작성**

- [ ] **Step 2: `DefaultLayout` 재작성 + `App.vue` + `useMobile` 수정**

- [ ] **Step 3: locales 세 곳 새 키 추가** (`nav.theme`, `nav.language` 등 실제 사용 키)

- [ ] **Step 4: 검증**

Run: `cd frontend && bun run build`
Expected: PASS

---

### Task 9: 뷰 마이그레이션 (10개)

**Files:**
- Modify: `frontend/src/views/{Dashboard,Logs,Providers,Models,ApiKeys,Usage,Performance,Users,Login,Register}View.vue`

**Interfaces:**
- Consumes: Task 4/6/7/8의 모든 `Or*` 컴포넌트, Task 5의 공용 클래스

**매핑 규칙 (모든 뷰 공통):**

| 기존 | 교체 |
|---|---|
| `<v-icon size="N">mdi-x</v-icon>` | `<OrIcon name="x" :size="N" />` |
| `<v-data-table :headers :items :loading density hide-default-footer :items-per-page="-1" fixed-header>` | `<OrDataTable :headers :items :loading :density>` + 필요 시 `#footer` 슬롯 |
| `#item.<key>` / `#no-data` | 그대로 유지 (슬롯명 불변) |
| `<v-dialog v-model :max-width="n" :fullscreen="isMobile">` | `<OrDialog v-model :width="n" :fullscreen="isMobile">` |
| `.btn-primary` / `.btn-tonal` / `.btn-ghost` / `.btn-danger` | `<OrButton variant="filled|tonal|text|outlined" tone="primary|danger">` |
| `.row-btn`, `.dialog-close`, `.pw-toggle`, `.copy-btn` | `<OrIconButton icon="…" label="…" />` |
| `.field-group` + `.field-input` | `<OrTextField v-model label hint error />` |
| `.field-select` + `<select>` | `<OrSelect v-model :options label />` |
| `.checkbox-row` + `<input type=checkbox>` | `<OrCheckbox v-model label hint />` |
| `.preset-chip` 그룹 (Performance) | `<OrSegmentedButton v-model :options />` |
| `.status-chip--*` | `StatusChip` 유지(내부 재작성됨) |
| `.btn-spinner` / `.loading-center` | `<OrSpinner />` |
| `window.alert(...)` | `useSnackbarStore().show(msg, {variant:"error"})` |
| `window.confirm(...)` | **유지** (네이티브 확인이 접근성·비용 면에서 우수) |

- [ ] **Step 1: Dashboard + Usage (차트·stat 카드)** — 차트 옵션 색을 role 변수에서 읽고
  `useThemeStore` 구독으로 테마 전환 시 갱신.
- [ ] **Step 2: Logs (표·필터·페이저)**
- [ ] **Step 3: Providers + Models (표·다이얼로그·체크리스트)**
- [ ] **Step 4: ApiKeys + Users (다이얼로그 4개·복사 피드백·스위치)**
- [ ] **Step 5: Performance (세그먼트·native 표 3개·차트 2개)**
- [ ] **Step 6: Login + Register (AuthShell 폼)**
- [ ] **Step 7: 각 스텝 후 검증**

Run: `cd frontend && bun run build`
Expected: PASS (10회 모두)

---

### Task 10: 잔여 `v-*` 제거 및 최종 검증

**Files:**
- Modify: `frontend/src/**/*` (남은 `v-icon`/`v-dialog`/`v-data-table`/`v-btn`/`v-app`/`v-main`/`v-container`/`v-navigation-drawer`)
- Modify: `frontend/src/styles/page-shared.css` (`.v-data-table` 미디어 규칙 제거)

- [ ] **Step 1: 잔재 검색**

Run: `grep -rn 'v-icon\|v-data-table\|v-dialog\|v-btn\|v-app\|v-main\|v-container\|v-navigation-drawer' frontend/src`
Expected: 결과 없음.

- [ ] **Step 2: 잔여 하드코딩 토큰 검색**

Run: `grep -rn '#[0-9a-fA-F]\{6\}' frontend/src --include=*.vue | grep -v 'styles/tokens'`
Expected: 허용 목록(차트 데이터 URI 등 문서화된 것) 외 결과 없음.

- [ ] **Step 3: 백엔드 회귀 확인**

Run: `cd backend && go vet ./... && go test ./...`
Expected: PASS (프론트 작업이 백엔드를 건드리지 않았는지의 회귀 게이트)

- [ ] **Step 4: 프론트 전체 검증**

Run: `cd frontend && bun run build`
Expected: `vue-tsc --noEmit` 0 error + `vite build` 성공.

- [ ] **Step 5: lockfile 확인**

Run: `grep -c 'vuetify\|@mdi/font' frontend/bun.lock`
Expected: 0

- [ ] **Step 6: 수동 체크리스트 이행 (스펙 8)**

라이트/다크 토글 유지·시스템 연동, 10개 화면 렌더, 모바일 600px 하단 내비,
reduced-motion, 대비 검산 5개 항목을 브라우저로 확인한다.

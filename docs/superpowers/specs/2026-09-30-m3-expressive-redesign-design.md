# OmniRelay 프론트엔드 Material 3 Expressive 전면 리디자인 — 설계 스펙

날짜: 2026-09-30
상태: 승인됨 (사용자가 설계 방향 5개 결정과 A안을 승인하고, 이후 "보여주지 말고 알아서 다 해놔"로 임상 승인 위임)

## 1. 목표

프론트엔드 전체를 **진짜 Material 3 Expressive** 인상으로 다시 디자인한다. 시각 스타일과
레이아웃/정보구조(IA)까지 전면 재설계하되, **기능·API·i18n 로직은 변하지 않는다**.

성공 기준:

- Vuetify 의존성이 제거되고, 디자인 시스템은 저장소 자체의 M3 Expressive 토큰 + 자체 Vue
  프리미티브로만 구성된다.
- 라이트/다크 듀얼 테마가 시스템 설정을 따르고 앱 내 토글로 강제 가능하다.
- 상태 변화 시 형태 모핑(shape morphing)과 스프링 모션이 실제로 동작한다.
- `prefers-reduced-motion: reduce`에서 모든 애니메이션이 즉시 종료 상태로 점프한다.
- 백엔드 `go vet ./... && go test ./...`, 프론트 `bun run build` 모두 통과한다.

## 2. 확정된 결정

| # | 항목 | 결정 |
|---|---|---|
| 1 | 출실도 | 진짜 M3 Expressive 인상 (Roboto 기반, 톤 서피스, 형태 모핑, 스프링) |
| 2 | 테마 | 라이트 + 다크 듀얼, `prefers-color-scheme` 연동 + 앱 내 토글, `localStorage` 지속 |
| 3 | 시드 | M3 baseline 보라색 `#6750A4` (`primary40`) |
| 4 | 강도 | 풀 Expressive — 형태 모핑 + 스프링. 앱 내 모션 토글 없음(시스템 설정만) |
| 5 | 깊이 | 시각 + 레이아웃/IA 전면 재설계 |
| 6 | 구현 | A안 — 커스텀 M3 시스템 구축, Vuetify 제거 |

제약(AGENTS.md):

- 새 UI 문자열은 `src/locales/{en,ja,ko}.ts` **세 곳 동시** 추가.
- 프론트 검증은 `bun run build`뿐 (lint/format/test 스크립트 없음).
- `package.json` 변경 후 **`bun install`로 `bun.lock` 갱신 필수** — `Dockerfile`이
  `bun install --frozen-lockfile`을 쓴다.
- 렌더링 없이 검증할 수 없으므로, 각 단계는 `bun run build`(vue-tsc 타입체크 포함)로 확인한다.
- 커밋하지 않는다(사용자 명시 요청 없이).

## 3. 토큰 아키텍처

`src/styles/tokens.css`는 디렉터리 `src/styles/tokens/`로 대체한다. `src/styles/tokens/index.css`가
유일한 진입점이고 `main.ts`에서 한 번 import한다. 모든 하드코딩 색/라디우스/지속시간은 이곳에
있다. 소스는 아래 표의 공식 M3 v0.192 토큰(material-components/material-web)과
M3 Expressive 토큰이다.

### 3.1 파일

| 파일 | 내용 |
|---|---|
| `tokens/palette.css` | `--m3-palette-*` 원시 톤 (0–100). 단일 시드 팔레트. |
| `tokens/color.css` | `:root[data-theme="light"]` / `[data-theme="dark"]` 블록의 role 49종. |
| `tokens/type.css` | M3 타입 스케일 15종 + emphasized 15종 + 서체 스택. |
| `tokens/shape.css` | corner ladder 10종. |
| `tokens/motion.css` | duration 16종, easing 6종, 스프링 6종(`linear()`), reduced-motion 게이트. |
| `tokens/elevation.css` | tonal elevation 매핑 + transient overlay 그림자. |
| `tokens/state.css` | state layer 불투명도. |
| `tokens/space.css` | 8dp 격자 간격 스케일. |
| `tokens/index.css` | 위 전부 import. |

### 3.2 테마 해석 규칙

JS가 항상 `<html data-theme="light|dark">`를 설정한다. FOUC 방지를 위해 `index.html` 인라인
스크립트가 첫 페인트 전에 `localStorage("omnirelay.theme")` → 없으면
`matchMedia("(prefers-color-scheme: dark)")`로 값을 읽어 `data-theme`를 쓴다.
`stores/theme.ts`(Pinia)가 토글 상태를 보유하고 동일 규칙을 재계산한다.
CSS에는 두 블록만 존재한다 — 중복 없음.

### 3.3 색 role (공식 M3 baseline, 시드 `#6750A4`)

라이트 (괄호 안은 palette 원소):

```
primary #6750A4(primary40) · on-primary #FFFFFF · primary-container #EADDFF(primary90)
on-primary-container #21005D(primary10)
secondary #625B71(secondary40) · on-secondary #FFFFFF · secondary-container #E8DEF8(secondary90)
on-secondary-container #1D192B(secondary10)
tertiary #7D5260(tertiary40) · on-tertiary #FFFFFF · tertiary-container #FFD8E4(tertiary90)
on-tertiary-container #31111D(tertiary10)
error #B3261E(error40) · on-error #FFFFFF · error-container #F9DEDC(error90)
on-error-container #410E0B(error10)
background #FEF7FF(neutral98) · on-background #1D1B20(neutral10)
surface #FEF7FF(neutral98) · on-surface #1D1B20(neutral10)
surface-variant #E7E0EC(nv90) · on-surface-variant #49454F(nv30)
outline #79747E(nv50) · outline-variant #CAC4D0(nv80)
surface-container-lowest #FFFFFF(neutral100) · -low #F7F2FA(neutral96)
-surface-container #F3EDF7(neutral94) · -high #ECE6F0(neutral92) · -highest #E6E0E9(neutral90)
surface-dim #DED8E1(neutral87) · surface-bright #FEF7FF(neutral98)
inverse-surface #322F35(neutral20) · inverse-on-surface #F5EFF7(neutral95)
inverse-primary #D0BCFF(primary80) · surface-tint #6750A4 · scrim/shadow #000000
```

다크:

```
primary #D0BCFF(primary80) · on-primary #381E72(primary20) · primary-container #4F378B(primary30)
on-primary-container #EADDFF(primary90)
secondary #CCC2DC(secondary80) · on-secondary #332D41(secondary20)
secondary-container #4A4458(secondary30) · on-secondary-container #E8DEF8(secondary90)
tertiary #EFB8C8(tertiary80) · on-tertiary #492532(tertiary20)
tertiary-container #633B48(tertiary30) · on-tertiary-container #FFD8E4(tertiary90)
error #F2B8B5(error80) · on-error #601410(error20) · error-container #8C1D18(error30)
on-error-container #F9DEDC(error90)
background #141218(neutral6) · on-background #E6E0E9(neutral90)
surface #141218(neutral6) · on-surface #E6E0E9(neutral90)
surface-variant #49454F(nv30) · on-surface-variant #CAC4D0(nv80)
outline #938F99(nv60) · outline-variant #49454F(nv30)
surface-container-lowest #0F0D13(neutral4) · -low #1D1B20(neutral10)
-surface-container #211F26(neutral12) · -high #2B2930(neutral17) · -highest #36343B(neutral22)
surface-dim #141218(neutral6) · surface-bright #3B383E(neutral24)
inverse-surface #E6E0E9(neutral90) · inverse-on-surface #322F35(neutral20)
inverse-primary #6750A4(primary40) · surface-tint #D0BCFF · scrim/shadow #000000
```

`palette.css`에는 시드에서 파생된 primary/secondary/tertiary/error/neutral/neutral-variant
전 톤을 위 hex로 보관한다(role은 palette에 `var()`로 위임 — palette 교체만으로 재생성 가능).

### 3.4 앱 의미 별칭 (M3에 없는 role)

M3 core에는 success/warning이 없다. container/on-container 쌍만 정의한다.
danger는 M3 `error` role을 그대로 쓴다.

| 토큰 | 라이트 | 다크 |
|---|---|---|
| `--or-success` / `--or-success-container` / `--or-on-success-container` | `#1B6B29` / `#C8E6C9` / `#0D3B13` | `#7CD98D` / `#1B4D22` / `#B7EAC1` |
| `--or-warning` / `--or-warning-container` / `--or-on-warning-container` | `#B26A00` / `#FBE3B0` / `#422B00` | `#FFDEA1` / `#5B4100` / `#FFDEA1` |

성공/경고 대비는WCAG AA (본문 4.5:1, 큰 글씨 3:1)를 만족해야 한다. 구현 단계에서
`--or-on-*-container` 대비를 검산한다.

### 3.5 타입

- 서체: `--m3-typeface-plain: Roboto` (본문/라벨), 표시용 동일.
  **CJK 대응**: `Roboto, "Noto Sans JP", "Noto Sans KR", "Apple SD Gothic Neo",
  "Hiragino Sans", system-ui, sans-serif` — CJK 웹폰트는 배포하지 않고 OS 폰트에 위임한다.
  코드/식별자: `"Roboto Mono", "JetBrains Mono", ui-monospace, monospace`.
  로딩은 `@fontsource/roboto` + `@fontsource/roboto-mono` 셀프호스트(CDN 의존 없음).
  Google Sans는 라이선스상 사용 불가 — M3 스펙 기본 서체인 Roboto로 통일한다.
- M3 타입 스케일 (공식, size/line-height/tracking):

| role | size | line-height | tracking | weight |
|---|---|---|---|---|
| display-large | 3.5625rem | 4rem | -0.015625rem | 400 |
| display-medium | 2.8125rem | 3.25rem | 0 | 400 |
| display-small | 2.25rem | 2.75rem | 0 | 400 |
| headline-large | 2rem | 2.5rem | 0 | 400 |
| headline-medium | 1.75rem | 2.25rem | 0 | 400 |
| headline-small | 1.5rem | 2rem | 0 | 400 |
| title-large | 1.375rem | 1.75rem | 0 | 400 |
| title-medium | 1rem | 1.5rem | 0.009375rem | 500 |
| title-small | 0.875rem | 1.25rem | 0.00625rem | 500 |
| body-large | 1rem | 1.5rem | 0.03125rem | 400 |
| body-medium | 0.875rem | 1.25rem | 0.015625rem | 400 |
| body-small | 0.75rem | 1rem | 0.025rem | 400 |
| label-large | 0.875rem | 1.25rem | 0.00625rem | 500 |
| label-medium | 0.75rem | 1rem | 0.03125rem | 500 |
| label-small | 0.6875rem | 1rem | 0.03125rem | 500 |

- **emphasized 세트** (`--m3-typescale-emphasized-*`): 크기/라인하이트는 baseline 그대로,
  weight는 한 단계 가늘게/두껍게 올리고 tracking만 조정한다.

| role | emphasized weight | emphasized tracking |
|---|---|---|
| display·headline·title-large | 500 | 0 (baseline 유지) |
| body-large | 500 | 0.15px |
| title-medium | 700 | 0.15px |
| title-small | 700 | 0.1px |
| body-medium | 500 | 0.25px |
| body-small | 500 | 0.4px |
| label-large | 700 | 0.1px |
| label-medium · label-small | 700 | 0.5px |

사용 규칙: 화면당 강조 요소 하나(핵심 제목, 활성 라벨, CTA)에만 사용한다.

### 3.6 형태 (shape)

공식 M3 corner ladder + M3E 확장:

| 토큰 | px |
|---|---|
| `--m3-shape-none` | 0 |
| `--m3-shape-xs` | 4 |
| `--m3-shape-sm` | 8 |
| `--m3-shape-md` | 12 |
| `--m3-shape-lg` | 16 |
| `--m3-shape-xl` | 20 (M3E large-increased) |
| `--m3-shape-2xl` | 28 (M3 extra-large) |
| `--m3-shape-3xl` | 32 (M3E) |
| `--m3-shape-4xl` | 48 (M3E) |
| `--m3-shape-full` | 9999px |

**형태 모핑**: 상태 변화 시 `border-radius`를 전환하고 모션 토큰으로 애니메이트한다.
핵심 규칙:

- 버튼 pressed: `full` ↔ `xs` 로 모핑 (pill → squarer)
- 카드 pressed: `2xl` → `lg`
- 칩 selected: `full` ↔ `md`
- 다이얼로그 enter: `4xl` → `2xl`
- 입력 필드 focus: `sm` → `xs`

35종 폴리곤 라이브러리(cookie/sunny/clover 등)는 **범위 밖**(YAGNI) — 로딩 인디케이터와
장식 요소에 필요한 최소 세트(`full`, `pill`, squircle 근사)만 구현한다.

### 3.7 모션

**전환 (durations/easing)** — 공식 M3 토큰:

```
duration: short1 50 / short2 100 / short3 150 / short4 200
          medium1 250 / medium2 300 / medium3 350 / medium4 400
          long1 450 / long2 500 / long3 550 / long4 600
          extra-long1 700 / extra-long2 800 / extra-long3 900 / extra-long4 1000
easing-standard            cubic-bezier(0.2, 0, 0, 1)
easing-emphasized          cubic-bezier(0.2, 0, 0, 1)
easing-emphasized-decel    cubic-bezier(0.05, 0.7, 0.1, 1)   // 등장
easing-emphasized-accel    cubic-bezier(0.3, 0, 0.8, 0.15)   // 퇴장
easing-linear              linear
```

앱 수준 전환 토큰 (Expressive tuning):

```
--or-motion-entrance: 360ms   // 화면/섹션 등장
--or-motion-settle:   240ms   // 제자리 상태 변화 (칩 선택, 확장/축소)
--or-motion-press:    140ms   // press 피드백
--or-motion-route-enter: 400ms  + emphasized-decelerate
--or-motion-route-exit:  200ms  + emphasized-accelerate
```

**스프링 (M3 motion-physics matrix)** — 질량 1, 강성(stiffness)/감쇠비(damping ratio):

| 토큰 | expressive | standard | 용도 |
|---|---|---|---|
| `spatial-fast` | k=800, r=0.6 | k=1400, r=0.9 | 소형(칩, 스위치, press) |
| `spatial-default` | k=380, r=0.8 | k=700, r=0.9 | 중형(카드, 시트) |
| `spatial-slow` | k=200, r=0.8 | k=300, r=0.9 | 전환/전체 화면 |
| `effects-fast` | k=3800, r=1 | 동일 | 색/불투명(무반동) |
| `effects-default` | k=1600, r=1 | 동일 | 색/불투명 |
| `effects-slow` | k=800, r=1 | 동일 | 색/불투명 |

- `spatial`은 과잉 감쇠 미만(r<1) → 약간의 오버슈트. `effects`는 임계 감쇠(r=1) → 오버슈트 없음.
- CSS에는 스프링 응답 곡선을 `linear()`로 샘플링해 넣는다. 생성기
  `frontend/scripts/gen-springs.mjs`가 샘플 값을 산출하고, 산출물은 `tokens/motion.css`에
  커밋한다(재현 가능성 유지).
- `linear()` 미지원 브라우저를 위해 `cubic-bezier` 폴백을 `@supports not` 블록에 제공한다.
- 기본 스킴은 **expressive**. `standard`는 토큰만 제공하고 기본으로 쓰지 않는다.

**reduced-motion 게이트** — 단일 규칙으로 전체 적용:

```css
@media (prefers-reduced-motion: reduce) {
  *, ::before, ::after {
    animation-duration: 1ms !important;
    animation-delay: 0ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 1ms !important;
    transition-delay: 0ms !important;
    scroll-behavior: auto !important;
  }
}
```

### 3.8 엘리베이션 & 간격

- **톤 우선**: 깊이는 그림자가 아니라 `surface-container-*` 톤으로 읽는다.
  level0 → `surface`, 1 → `-low`, 2 → `-container`, 3 → `-high`, 4/5 → `-highest`.
- 그림자는 일시적 오버레이(다이얼로그, 메뉴)에만: `--m3-elevation-3` =
  `0 4px 8px 3px rgba(0,0,0,.15), 0 1px 3px rgba(0,0,0,.3)`.
- 8dp 격자: 2 / 4 / 8 / 12 / 16 / 24 / 32 / 48. 최소 탭 타깃 48×48.

### 3.9 상태 레이어

`--m3-state-hover: 0.08`, `--m3-state-focus: 0.10`, `--m3-state-press: 0.10`,
`--m3-state-drag: 0.16`. 대상 색(`on-surface`/`on-primary`/…)의 `color-mix()`로 합성하고,
호버·포커스·프레스 상태를 `::after` 오버레이로 구현한다.

## 4. 컴포넌트 레이어

Vuetify 제거. 자체 프리미티브는 `src/components/ui/`에 `Or` 접두사로 둔다
(기존 앱 컴포넌트 `PageHeader` 등과 구분). 전역 등록은 `src/plugins/components.ts`.

### 4.1 프리미티브 (`src/components/ui/`)

| 컴포넌트 | 필요한 이유 / 대체 대상 |
|---|---|
| `OrIcon.vue` | `v-icon` 58개 대체. Material Symbols(변수 폰트) 렌더. |
| `OrButton.vue` | filled / tonal / outlined / text / elevated + danger 톤, sm/md/lg, loading. `.btn-primary`·`.btn-tonal`·`.btn-ghost`·`.btn-danger`·`.auth-submit` 대체. |
| `OrIconButton.vue` | `.row-btn`, `.dialog-close`, `.pw-toggle`, `.copy-btn` 대체. |
| `OrCard.vue` | elevated(톤) / filled / outlined. `.table-card`, `.chart-panel`, `.stat-card` 기반. |
| `OrChip.vue` | assist / filter / input. `.preset-chip`, `.type-chip`, `.status-chip` 기반. |
| `OrSegmentedButton.vue` | PerformanceView의 1H/24H/7D/30D 프리셋 그룹. |
| `OrTextField.vue` | outlined/filled, label, hint, error, leading/trailing 슬롯. `.field-input`·`.field-group` 대체. `type` 패스스루(text/date/password/number). |
| `OrSelect.vue` | M3 outlined select(네이티브 `<select>` 캡). `.field-select` 대체. |
| `OrCheckbox.vue` | `.checkbox-row`·`.checkbox` 대체. 체크리스트/다중 선택. |
| `OrSwitch.vue` | 활성/비활성 토글(API 키 활성, 사용자 역할). |
| `OrDialog.vue` | `v-model`, `width`, 모바일 fullscreen. **포커스트랩 + ESC + `aria-modal` + 스크림**. |
| `OrDataTable.vue` | 아래 4.2 참조. `v-data-table` 5개 대체. |
| `OrSnackbar.vue` + `stores/snackbar.ts` | copy 피드백, 저장/삭제 알림. `window.alert` 대체(경고 아님). |
| `OrSpinner.vue` | `.btn-spinner`, `.loading-center`, `.submit-spinner` 대체. |
| `OrProgressBar.vue` | 대시보드 stat bar, 다운로드/작업 표시. |
| `OrDivider.vue` | `.nav-divider`. |
| `OrNavRail.vue` / `OrNavigationBar.vue` | 아래 5장. |

전역 유틸: `composables/useRipple.ts`(클릭 시 state layer ripple)과
`composables/useFocusTrap.ts`(다이얼로그용).

### 4.2 `OrDataTable` 계약

`v-data-table`의 실제 사용만 만족하도록 최소 계약을 고정한다.

```
props:
  headers: { key: string; title: string; sortable?: boolean;
             align?: "start"|"end"; minWidth?: number; width?: number }[]
  items: Record<string, any>[]
  loading?: boolean          // true면 오버레이 스피너
  density?: "compact" | "comfortable"   // 기본 comfortable
  hideHeader?: boolean
emits:
  sort?: (key: string, dir: "asc" | "desc" | null) => void   // sortable 헤더 클릭 시
slots:
  #item.<key>     — 셀 커스터마이즈 (13개 로그 컬럼 등)
  #no-data        — 빈 상태
  #footer         — 커스터마이즈 페이저 (LogsView 전용)
  #bottom         — 페이저 아래 보조 영역
구조:
  <table role="grid">, <th aria-sort="ascending|descending|none">
  모바일(<768px)에서는 스크롤 컨테이너 대신 호출 쪽이 MobileDataCard를 쓰므로
  컴포넌트는 스크롤만 제공한다(반드시 `overflow-x:auto` + `tabindex=0` + aria-label).
정렬:
  서버 정렬이 아니다. sort 상태는 컴포넌트 내부에 두고 `sort` 이벤트로 상위에 통보만 한다.
  상위가 무시하면 현재 정렬 상태를 그대로 표시한다.
```

### 4.3 앱 컴포넌트 (재스타일 유지)

`AppAlert`, `AuthShell`, `EmptyState`, `MobileDataCard`, `MonoTag`, `PageHeader`,
`StatCard`, `StatusChip` — API(prop/slot)는 **변경하지 않고** 내부만 M3로 다시 쓴다.
이름·경로 변경 없음(호출부 7개 뷰가 그대로 참조).

## 5. 레이아웃 / IA

### 5.1 데스크톱: M3 navigation rail → drawer

현재 `v-navigation-drawer` + rail 토글(220px)을 자체 `OrNavRail`로 대체.

- **`expanded` (≥1200px)**: 256px 고정 drawer. 서피스는 `surface-container-low`.
  상단 브랜드 영역, 그 아래 48px 간격 리스트, 하단에 로케일·테마 토글·로그아웃.
  활성 항목은 `primary-container` 배경 + `on-primary-container` 텍스트 + leading 아이콘
  pill (M3 활성 내비 규칙). 상태 변화는 `spatial-fast` 스프링으로 모핑.
- **`rail` (600–1199px)**: 80px 아이콘 전용 rail. 활성 항목은 아이콘 뒤 `primary-container`
  원형 배경 + 아이콘 하단 4px dot. hover 시 툴팁(label).
- **모바일 (<600px)**: 하단 `OrNavigationBar` (80px + safe-area). 활성 탭은 아이콘 뒤
  `secondary-container` pill + 아래 label. 현재 56px 탭바를 대체.

브레이크포인트는 M3 window size class를 따른다: compact <600, medium 600–839,
expanded ≥840. rail/drawer 전환은 **1100px**을 경계로 둔다(현재 `useMobile()`의 768px과
별개로, 탭/모바일 경계는 600px로 상향).

### 5.2 페이지 셸

- 최대 폭 1440px, 좌우 패딩 24px(모바일 16px), 상단 32px.
- 페이지 헤더: `headline-small` emphasized 제목 + `body-medium` `on-surface-variant`
  부제 + 우측 액션 영역. 상단에 1px `outline-variant` 구분선 없음(여백으로 분리).
- 섹션 간격 32px, 카드 간격 16px.

### 5.3 화면별 IA

- **Dashboard**: stat 카드 5 → 2행 그리드(lg: 3열 + 2, md: 2열, sm: 1열). 카드는
  `surface-container-low` + `corner-2xl`, 값은 `display-small` emphasized. 하단은
  차트 2/3 + 정보 패널 1/3.
- **Logs**: 필터 바를 카드 안으로 집어넣고 Apply를 오른쪽 끝에 붙인다. 표는
  `surface-container-lowest` 표면 + `outline-variant` 행 구분. 페이저는 M3 표준
  "1–50 / 전체" + 이전/다음 아이콘 버튼.
- **Providers / Models / ApiKeys / Users**: 표 + 우측 액션. 액션은 아이콘 버튼 3–4개
  (overflow menu 대신 노출 — 현재 상태 유지, YAGNI). 헤더 우측에 FAB가 아닌 `OrButton`
  filled로 "추가".
- **Performance**: 프리셋 칩을 `OrSegmentedButton`으로, 필터를 한 줄 `filter-bar`로.
  native `<table>` 3개는 그대로 두되 표 카드 스타일과 sticky header만 M3화.
- **Usage**: stat 카드 + 차트. 차트 색은 role 기반 팔레트로 교체.
- **Login / Register**: `AuthShell`을 M3 카드로. 배경 그리드/글로우 제거하고
  `surface-container-low` 전체 배경 + 중앙 `corner-4xl` 카드 + 미묘한 `primary` 액센트만.
  필드는 M3 outlined text field(라벨 플로팅)로 교체.

### 5.4 차트

`chart.js` 유지. 색은 하드코딩 hex를 role에서 읽는다:
`--m3-color-primary`(라인1), `--m3-color-tertiary`(라인2),
grid `outline-variant` 40%, 눈금 글자 `on-surface-variant`, 툴팁
`surface-container-highest` + `on-surface`. 테마 전환 시 차트를 파괴/재생성하지 않고
색만 갱신한다(`stores/theme.ts`를 구독).

## 6. 아이콘

`@mdi/font` 제거 → **Material Symbols**(변수 폰트, `@fontsource-variable/material-symbols`
또는 `index.html`의 Google Fonts 링크). `OrIcon`은 `ligature` 문자열을 그대로 렌더한다.

- 58개 `mdi-*` 호출부를 Material Symbols 이름으로 매핑한다. 매핑표는 구현 단계에서
  `frontend/scripts/mdi-to-symbols.json`으로 보관한다.
- `weight`는 `400` 기본, `fill`은 활성/강조에만. 크기·색은 prop/상속.

## 7. 범위 밖

- 백엔드 API·데이터 모델·마이그레이션 변경 없음.
- 라우트/스토어/i18n 키 구조 변경 없음 (새 문자열 키 추가만 허용).
- 35종 폴리곤 shape 라이브러리, FAB 메뉴, split button, toolbar, wavy progress 등
  M3E 신규 컴포넌트 중 현재 화면이 쓰지 않는 것.
- 다크 모드 외 추가 테마(고대비), 앱 내 모션 토글.
- 퍼포먼스 최적화·코드 스플리팅 재설계.

## 8. 검증

각 단계 완료 시 (실패해도 진행하지 않고 되돌린다):

```bash
cd frontend && bun run build      # vue-tsc --noEmit + vite build — 유일한 프론트 검증
cd backend  && go vet ./... && go test ./...   # 레이아웃/라우트가 백엔드를 건드리지 않았는지
```

수동 확인 체크리스트 (구현 완료 시):

1. 라이트/다크 토글이 즉시 반영되고, 새로고침 후에도 유지되며, 시스템 설정을 따른다.
2. 로그인 → 대시보드 → 로그 → 성능 → 사용량 → 제공자 → 모델 → API키 → 사용자 전 화면이
   렌더링되고 빈 상태/에러/로딩이 보인다.
3. 모바일(≤599px)에서 하단 네비게이션 바와 카드형 목록이 동작한다.
4. 브라우저 개발자도구에서 `prefers-reduced-motion: reduce`를 켜면 전환이 즉시 끝난다.
5. 대비 검산: 본문 텍스트 4.5:1 이상, 상태 칩 라벨 4.5:1 이상.

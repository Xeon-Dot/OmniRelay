# M3 Expressive 뷰 마이그레이션 가이드 (Task 9)

대상: `frontend/src/views/*.vue` 10개. 각 서브에이전트는 지정된 뷰만 수정한다.
완료 후 검증은 `cd frontend && bunx vue-tsc --noEmit` — **`bun run build`는 금지**
(병렬 실행 시 `dist/` 충돌). 전체 빌드는 본인이 병렬로 수행한다.

## 절대 규칙

1. **로직·API·i18n 키 구조는 건드리지 않는다.** 템플릿과 `<style>`만 바꾼다.
   script의 데이터 흐름, store 호출, 라우팅, 필터/페이지네이션 로직은 그대로.
2. **하드코딩 색·라디우스·지속시간 금지.** 전부 `var(--m3-*)` / `var(--or-*)`.
3. **새 UI 문자열**은 `src/locales/{en,ja,ko}.ts` 세 곳에 동시에 추가 (없으면 만들지 말 것 —
   기존 키를 재사용).
4. **아이콘**: `<v-icon size="N">mdi-xxx</v-icon>` → `<OrIcon name="..." :size="N" />`.
   이름은 `scripts/mdi-to-symbols.json`의 매핑을 따른다. JSON에 없으면 임의 변환하지 말고
   표기 그대로(`mdi-...`) 남겨두고 보고한다.
5. `EmptyState`의 `icon` prop도 같은 매핑을 따른다 (`icon="mdi-x"` → `icon="x"`).
6. 커밋하지 않는다.

## 컴포넌트 API (정확한 시그니처)

전역 등록돼 있어 import 불필요.

```
OrIcon        name* string, size=24, color?, fill?=false, label?
OrButton      variant?=filled|tonal|outlined|text|elevated, tone?=primary|danger,
              size?=sm|md|lg, loading?, disabled?, type?=button|submit|reset, block?
              slots: default / leading / trailing    emits: click
              로딩 중에는 label을 감추고 spinner를 중앙에 놓인다 (너비 유지)
OrIconButton  icon* string, label* string(aria-label), variant?=standard|filled|tonal|outlined,
              tone?=default|primary|danger, size?=sm|md|lg, disabled?, fill?   emits: click
OrCard        variant?=elevated|filled|outlined, interactive?, flush?          emits: click
OrChip        variant?=assist|filter, selected?, disabled?, icon?              emits: click
OrSegmentedButton  options* {value,label?,icon?}[], modelValue?, multiple?, disabled?, size?=md|sm
              emits: update:modelValue (string | string[])
OrTextField   modelValue?, label* , hint?, error?, type?="text", disabled?, readonly?,
              leadingIcon?, clearable?, placeholder?, autocomplete?, name?, min?, max?, step?
              emits: update:modelValue(string), clear
              `error`가 비어있지 않으면 헌트 대신 에러 표시. type="date"는 라벨 항상 떠 있음.
OrSelect      modelValue?, label*, options* {value,label}[], hint?, error?, disabled?, placeholder?
              emits: update:modelValue(string)
OrCheckbox    modelValue?: boolean|unknown[], value?, label?, hint?, disabled?, error?
              emits: update:modelValue  — `value`를 주고 modelValue에 배열을 넘기면 다중선택
OrSwitch      modelValue?: boolean, label?, disabled?   emits: update:modelValue
OrDialog      modelValue*, width?=480, fullscreen?, dismissible?=true, title?, ariaLabel?
              slots: default / header / footer        emits: update:modelValue
              포커스트랩·ESC·스크림·스크롤잠금 내장. 슬롯 내용은 기존 다이얼로그 카드 마크업을 그대로
              default 슬롯에 넣고, `.dialog-card/.dialog-header/...` 클래스는 제거한다.
              footer 버튼은 `#footer` 슬롯에 → `<template #footer>`.
OrDataTable   headers* {key,title,sortable?,align?:start|end,minWidth?,width?}[],
              items* Record<string,unknown>[], loading?, density?=comfortable|compact,
              hideHeader?, label?
              slots: #item.<key> (scope: item, value, index) / #no-data / #footer / #bottom
              emits: sort(key, dir)
              `#item.<key>`·`#no-data` 슬롯명과 props는 v-data-table과 동일 → 그대로 유지.
              `hide-default-footer`, `:items-per-page="-1"`, `fixed-header` prop은 제거한다
              (내장 스티키 헤더가 처리). 페이지네이션은 `#footer` 슬롯에 그대로 넣는다.
OrSpinner     size?=20, label?
OrProgressBar value?(0..1), tone?=primary|success|warning|danger, label?
OrDivider     vertical?
```

앱 컴포넌트(시그니처 불변, 내부는 이미 재작성됨):
`PageHeader{title,subtitle}+slot`, `StatCard{label,hint,sub,valueClass}+slots default/sub`,
`StatusChip{variant?:on|off|warning}+slot`, `MonoTag{slot}`, `EmptyState{icon,text,iconSize,small}+slot`,
`MobileDataCard{items:{label,value}[]}+slots actions`, `AppAlert{variant?:success|error|warning,page?}+slot`.

스토어:
- `import { useSnackbarStore } from "../stores/snackbar"` → `show(msg,{variant})`, `success(msg)`, `error(msg)`, `hide()`.
  `window.alert(...)`는 이걸로 교체. **`window.confirm(...)`은 유지** (의도적 결정).
- `import { useChartTheme } from "../composables/useChartTheme"` → `{ palette }` (ComputedRef).

## 매핑 표

| 기존 | 교체 |
|---|---|
| `<v-icon size="N">mdi-x</v-icon>` | `<OrIcon name="x" :size="N" />` |
| `<v-data-table :headers :items :loading density hide-default-footer :items-per-page="-1" fixed-header>` | `<OrDataTable :headers :items :loading :density label="...">` |
| `#item.<key>` / `#no-data` 슬롯 | **그대로 유지** |
| `<v-dialog v-model :max-width="n" :fullscreen="isMobile">` | `<OrDialog v-model :width="n" :fullscreen="isMobile" title="...">` |
| `.dialog-card > .dialog-header/.dialog-body/.dialog-footer` | 다이얼로그 default/header/footer 슬롯으로 분리, 클래스 제거 |
| `.btn-primary` `<button>` | `<OrButton variant="filled">` |
| `.btn-tonal` `<button>` | `<OrButton variant="tonal">` |
| `.btn-ghost` `<button>` | `<OrButton variant="text">` |
| `.btn-danger` `<button>` | `<OrButton variant="filled" tone="danger">` |
| `.btn-secondary` `<button>` | `<OrButton variant="outlined">` |
| `.row-btn`, `.dialog-close`, `.pw-toggle`, `.copy-btn` | `<OrIconButton icon="..." label="...">` (label은 i18n 키) |
| `.filter-bar > .filter-col > .field-label + .field-input` | `<div class="filter-col"><OrTextField v-model label hint /></div>` |
| `.field-group > .field-label + .field-input` | `<OrTextField v-model label hint error />` |
| `.field-select` + `<select>` | `<OrSelect v-model :options label />` |
| `.checkbox-row` + `<input type=checkbox>` | `<OrCheckbox v-model label hint />` |
| `.preset-chip` 그룹 (Performance 1H/24H/7D/30D) | `<OrSegmentedButton v-model="preset" :options="[...]" />` |
| `.status-chip--*` / `StatusChip` | 유지 |
| `.btn-spinner` / `.submit-spinner` / `.loading-center` | `<OrSpinner />` |
| `.btn-spinner`가 동작 중인 버튼 안 | `<OrButton :loading="saving">` 로 교체 |
| `.alert alert--error` `<p>` | `<AppAlert variant="error">` |
| `.stat-card` 로컬 마크업 (Dashboard) | `<StatCard label>` 또는 OrCard로 재구성 |
| `.type-chip` | `<OrChip>` |
| `window.alert(msg)` | `useSnackbarStore().show(msg, { variant:"info"|"success"|"error" })` |
| `window.confirm(...)` | **유지** |

## 토큰 (자주 쓰는 것)

```
색   --m3-color-{primary,on-primary,primary-container,on-primary-container,
       secondary-container,on-secondary-container,tertiary,
       error,error-container,on-error-container,
       surface,surface-container{-lowest,-low,,-high,-highest},surface-variant,
       on-surface,on-surface-variant,outline,outline-variant,background,on-background,
       inverse-surface,inverse-on-surface,inverse-primary}
앱   --or-success, --or-success-container, --or-on-success-container,
     --or-warning, --or-warning-container, --or-on-warning-container
형태 --m3-shape-{none,xs(4),sm(8),md(12),lg(16),xl(20),2xl(28),3xl(32),4xl(48),full}
     --or-morph-{button,card,chip,field,dialog,nav}-*
타입 --m3-typescale-{display,headline,title,body,label}-{large,medium,small}
     --m3-typescale-emphasized-*  /  --m3-typeface-{plain,mono}
     헬퍼 클래스: .m3-typescale-title-large 등 15종
간격 --m3-space-{25(2),50(4),75(6),100(8),150(12),200(16),300(24),400(32),600(48)}
     --or-gap-page(32) --or-gap-card(16) --or-pad-card(24) --or-gap-filter(12)
     --m3-min-touch-target(48px)  --or-layout-max-width(1440px)
모션 --or-motion-spatial-{fast,default,slow}  (duration+spring, `transition:` 축약용)
     --or-motion-effects-{fast,default,slow}
     --or-motion-entrance(360ms) --or-motion-settle(240ms) --or-motion-press(140ms)
     --m3-easing-{standard,emphasized-decelerate,emphasized-accelerate}
그림자 --m3-shadow-{1..5}  스크림 --m3-scrim
상태 --m3-state-{hover,press,focus,drag}  포커스 --m3-focus-ring-{color,width,offset}
```

`transition: border-radius var(--or-motion-spatial-fast);` 형태로 쓴다 (duration+curve 축약).

## 차트 (Dashboard, Usage, Performance)

```ts
import { useChartTheme } from "../composables/useChartTheme";
const { palette } = useChartTheme();
// 기존 `const options = {...}` 를 computed로 바꾼다:
const options = computed(() => ({
  ...기존 구조,
  scales: {
    x: { ticks: { color: palette.value.text }, grid: { color: palette.value.grid } },
    y: { ticks: { color: palette.value.text }, grid: { color: palette.value.grid } },
  },
  plugins: { tooltip: { backgroundColor: palette.value.surface, titleColor: palette.value.onSurface, bodyColor: palette.value.onSurface } },
}));
```
vue-chartjs는 `options` prop을 watch하므로 테마 전환 시 자동 갱신된다.
`<Line :options="options" :data="..." />` — data의 series 색도 `palette.value.*`에서 읽는다.

## `<style>` 규칙

- 기존 로컬 `<style scoped>` 안의 `#hex`, `rgba(...)`, `--clr-*`, `--font-display`,
  `--radius-*`, `--space-*`, `--dur-*` 전부 제거 → 토큰으로 교체.
- `font-family: "Fraunces"` / `"DM Sans"` / `"JetBrains Mono"` 참조 제거
  (`var(--m3-typeface-plain)` / `var(--m3-typeface-mono)`).
- `@media (max-width: 768px)` 기준의 `.v-data-table{display:none}` + `.mobile-cards`
  규칙은 **600px**로 통일하고, 테이블은 `OrDataTable`로 대체되므로
  `display:none` 대신 `OrDataTable`을 감싸는 스크롤 컨테이너를 그대로 둔다.
  모바일 카드(`.mobile-cards`)는 `display:none` 기본 + 600px 이하에서만 flex.
- `overflow-x: auto` 컨테이너는 `tabindex="0"`과 `aria-label`을 갖는다 (OrDataTable 내장).
- 사용하지 않는 클래스는 제거한다 (드리프트 방지).

## 완료 조건

1. `cd frontend && bunx vue-tsc --noEmit` → 에러 0.
2. `grep -n 'v-icon\|v-data-table\|v-dialog' <수정한 파일>` → 결과 없음.
3. `grep -nE '#[0-9a-fA-F]{6}|rgba\(' <수정한 파일>` → 문서화된 예외 외 결과 없음.
4. 로직이 바뀌지 않았는지 자기 점검: 필터/페이지네이션/저장/삭제 흐름 동일.

<script setup lang="ts">
/**
 * Modal dialog. Content stays in the default slot, exactly as the legacy Vuetify dialog
 * call sites had it — this component owns only the surface, the scrim and the
 * accessibility contract:
 *
 *   role="dialog" + aria-modal, focus moves in on open, Tab is trapped,
 *   Escape closes, the scrim closes (unless `dismissible` is false), focus
 *   returns to the trigger on close, and body scroll is locked (counted, so
 *   nested dialogs unlock correctly).
 *
 * Motion: the scrim fades on an effects spring; the container scales up and
 * morphs from corner-4xl down to corner-2xl — M3 Expressive's entrance.
 */
import { computed, onBeforeUnmount, ref, useId, watch } from "vue";
import { lockBodyScroll, unlockBodyScroll, useFocusTrap } from "../../composables/useFocusTrap";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    width?: number | string;
    /** Fills the viewport — call sites pass this when `isMobile`. */
    fullscreen?: boolean;
    /** Bottom-sheet presentation: anchored to the bottom edge, full width. */
    sheet?: boolean;
    /** Escape / scrim / close button behaviour. */
    dismissible?: boolean;
    title?: string;
    /** Rendered as the accessible name when `title` is absent. */
    ariaLabel?: string;
  }>(),
  {
    width: 480,
    fullscreen: false,
    sheet: false,
    dismissible: true,
    title: undefined,
    ariaLabel: undefined,
  },
);

const emit = defineEmits<{ "update:modelValue": [value: boolean] }>();
const panel = ref<HTMLElement | null>(null);
/** Unique per instance so two dialogs never share a title id. */
const titleId = useId();

const dialogOpen = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit("update:modelValue", value),
});

useFocusTrap(panel, dialogOpen);

let locked = false;
watch(
  () => props.modelValue,
  (open) => {
    if (open && !locked) {
      lockBodyScroll();
      locked = true;
    } else if (!open && locked) {
      unlockBodyScroll();
      locked = false;
    }
  },
);

// A dialog torn down while open (route change, v-if) must release its lock,
// or the page stays unscrollable for the rest of the session.
onBeforeUnmount(() => {
  if (locked) {
    unlockBodyScroll();
    locked = false;
  }
});

function requestClose() {
  if (props.dismissible) emit("update:modelValue", false);
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === "Escape") {
    event.stopPropagation();
    requestClose();
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="or-dialog">
      <div
        v-if="modelValue"
        class="or-dialog"
        :class="{ 'or-dialog--fullscreen': fullscreen, 'or-dialog--sheet': sheet }"
        @keydown="onKeydown"
      >
        <div
          class="or-dialog__scrim"
          aria-hidden="true"
          @click="requestClose"
        />
        <div
          ref="panel"
          class="or-dialog__panel"
          :style="{ maxWidth: fullscreen || sheet ? 'none' : typeof width === 'number' ? `${width}px` : width }"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="title ? `or-dialog-title-${titleId}` : undefined"
          :aria-label="title ? undefined : ariaLabel"
          tabindex="-1"
        >
          <div v-if="title || $slots.header || dismissible" class="or-dialog__header">
            <slot name="header">
              <h2 v-if="title" :id="`or-dialog-title-${titleId}`" class="or-dialog__title">{{ title }}</h2>
            </slot>
            <button
              v-if="dismissible"
              class="or-dialog__close"
              type="button"
              :aria-label="$t('common.close')"
              @click="requestClose"
            >
              <OrIcon name="close" :size="24" />
            </button>
          </div>

          <div class="or-dialog__body">
            <slot />
          </div>

          <div v-if="$slots.footer" class="or-dialog__footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.or-dialog {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--m3-space-300);
}

.or-dialog--fullscreen {
  padding: 0;
}

.or-dialog__scrim {
  position: absolute;
  inset: 0;
  background: var(--m3-scrim);
}

.or-dialog__panel {
  position: relative;
  display: flex;
  flex-direction: column;
  inline-size: 100%;
  max-block-size: min(90dvh, 800px);
  background: var(--m3-elevation-3-surface);
  color: var(--m3-color-on-surface);
  border-radius: var(--or-morph-dialog-rest);
  box-shadow: var(--m3-shadow-4);
  outline: none;
}

.or-dialog--fullscreen .or-dialog__panel {
  max-block-size: 100dvh;
  block-size: 100dvh;
  max-inline-size: none !important;
  border-radius: 0;
}

/* ── Sheet presentation (mobile "More" navigation) ── */
.or-dialog--sheet {
  align-items: flex-end;
  padding: 0;
}

.or-dialog--sheet .or-dialog__panel {
  inline-size: 100%;
  max-block-size: min(85dvh, 720px);
  border-end-start-radius: 0;
  border-end-end-radius: 0;
}

.or-dialog__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--m3-space-200);
  padding: var(--m3-space-300) var(--m3-space-300) var(--m3-space-150);
}

.or-dialog__title {
  margin: 0;
  font: var(--m3-typescale-emphasized-title-large);
  letter-spacing: var(--m3-typescale-emphasized-title-large-tracking);
}

.or-dialog__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 40px;
  block-size: 40px;
  margin: calc(var(--m3-space-100) * -1) calc(var(--m3-space-75, 6px) * -1) 0 0;
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
  cursor: pointer;
  transition: background-color var(--or-motion-effects-fast);
}

.or-dialog__close:hover {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
}

.or-dialog__body {
  padding: var(--m3-space-150) var(--m3-space-300) var(--m3-space-300);
  overflow-y: auto;
}

.or-dialog__footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--m3-space-150);
  padding: 0 var(--m3-space-300) var(--m3-space-300);
  flex-wrap: wrap;
}

/* ── Motion ── */
.or-dialog-enter-active,
.or-dialog-leave-active {
  transition: opacity var(--or-motion-effects-fast);
}
.or-dialog-enter-active .or-dialog__panel,
.or-dialog-leave-active .or-dialog__panel {
  transition:
    transform var(--or-motion-spatial-fast),
    border-radius var(--or-motion-spatial-fast),
    opacity var(--or-motion-effects-fast);
}

.or-dialog-enter-from,
.or-dialog-leave-to {
  opacity: 0;
}
.or-dialog-enter-from .or-dialog__panel,
.or-dialog-leave-to .or-dialog__panel {
  opacity: 0;
  transform: scale(0.92);
  border-radius: var(--or-morph-dialog-enter);
}

/* The sheet slides up from the bottom edge instead of scaling in. */
.or-dialog--sheet.or-dialog-enter-from .or-dialog__panel,
.or-dialog--sheet.or-dialog-leave-to .or-dialog__panel {
  transform: translateY(100%);
  border-radius: var(--or-morph-dialog-rest);
}
</style>

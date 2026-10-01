<script setup lang="ts">
/**
 * M3 button: filled / tonal / outlined / text / elevated, primary or danger
 * tone, three sizes. Rests as a pill and morphs to a tight corner while
 * pressed — the signature M3 Expressive state change.
 *
 * Loading keeps the button's rendered width: the label is hidden, not removed,
 * so the surrounding layout never jumps.
 */
import OrSpinner from "./OrSpinner.vue";

withDefaults(
  defineProps<{
    variant?: "filled" | "tonal" | "outlined" | "text" | "elevated";
    tone?: "primary" | "danger";
    size?: "sm" | "md" | "lg";
    loading?: boolean;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    /** Stretch to the container width. */
    block?: boolean;
  }>(),
  {
    variant: "filled",
    tone: "primary",
    size: "md",
    loading: false,
    disabled: false,
    type: "button",
    block: false,
  },
);
</script>

<template>
  <button
    v-ripple
    class="or-btn or-interactive"
    :class="[
      `or-btn--${variant}`,
      `or-btn--${tone}`,
      `or-btn--${size}`,
      { 'or-btn--block': block, 'is-loading': loading },
    ]"
    :type="type"
    :disabled="disabled || loading"
    :aria-busy="loading ? 'true' : undefined"
  >
    <span class="or-btn__content">
      <span v-if="$slots.leading" class="or-btn__affix">
        <slot name="leading" />
      </span>
      <span class="or-btn__label"><slot /></span>
      <span v-if="$slots.trailing" class="or-btn__affix">
        <slot name="trailing" />
      </span>
    </span>
    <OrSpinner v-if="loading" class="or-btn__spinner" :size="size === 'lg' ? 22 : 18" />
  </button>
</template>

<style scoped>
.or-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--m3-space-100);
  position: relative;
  padding-inline: var(--m3-space-300);
  border: none;
  border-radius: var(--or-morph-button-rest);
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
  text-decoration: none;
  white-space: nowrap;
  transition:
    border-radius var(--or-motion-spatial-fast),
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast),
    box-shadow var(--or-motion-effects-fast);
}

.or-btn:hover:not(:disabled) {
  text-decoration: none;
}

/* M3 Expressive press: the pill goes squarer. */
.or-btn:active:not(:disabled) {
  border-radius: var(--or-morph-button-active);
}

.or-btn:disabled {
  opacity: 0.38;
}

.or-btn--block {
  display: flex;
  width: 100%;
}

/* ── Sizes (sm is 40px: M3's own standard button height) ── */
.or-btn--sm {
  min-height: 40px;
  padding-inline: var(--m3-space-200);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}
.or-btn--md {
  min-height: 48px;
}
.or-btn--lg {
  min-height: 56px;
  padding-inline: var(--m3-space-400);
  font: var(--m3-typescale-title-small);
  letter-spacing: var(--m3-typescale-title-small-tracking);
}

/* ── Variants × tones ── */
.or-btn--filled.or-btn--primary {
  background: var(--m3-color-primary);
  color: var(--m3-color-on-primary);
}
.or-btn--filled.or-btn--danger {
  background: var(--m3-color-error);
  color: var(--m3-color-on-error);
}

.or-btn--tonal.or-btn--primary {
  background: var(--m3-color-secondary-container);
  color: var(--m3-color-on-secondary-container);
}
.or-btn--tonal.or-btn--danger {
  background: var(--m3-color-error-container);
  color: var(--m3-color-on-error-container);
}

.or-btn--elevated.or-btn--primary {
  background: var(--m3-elevation-1-surface);
  color: var(--m3-color-primary);
  box-shadow: var(--m3-shadow-1);
}
.or-btn--elevated.or-btn--danger {
  background: var(--m3-elevation-1-surface);
  color: var(--m3-color-error);
  box-shadow: var(--m3-shadow-1);
}

.or-btn--outlined {
  background: transparent;
  border: 1px solid var(--m3-color-outline);
}
.or-btn--outlined.or-btn--primary {
  color: var(--m3-color-primary);
}
.or-btn--outlined.or-btn--danger {
  color: var(--m3-color-error);
}

.or-btn--text {
  background: transparent;
  padding-inline: var(--m3-space-200);
}
.or-btn--text.or-btn--primary {
  color: var(--m3-color-primary);
}
.or-btn--text.or-btn--danger {
  color: var(--m3-color-error);
}

/* ── Loading: keep the box, hide the content ── */
.or-btn__content {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-100);
  min-width: 0;
}
.or-btn.is-loading .or-btn__content {
  visibility: hidden;
}
.or-btn__label {
  overflow: hidden;
  text-overflow: ellipsis;
}
.or-btn__affix {
  display: inline-flex;
  flex-shrink: 0;
}
.or-btn__spinner {
  position: absolute;
}
</style>

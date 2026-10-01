<script setup lang="ts">
/**
 * Square icon-only button. `label` is required — an icon button with no
 * accessible name is unusable with a screen reader.
 *
 * Sizes match the button ladder (40 / 48 / 56) so every hit target is at
 * least 40px and the default clears M3's 48dp minimum.
 */
withDefaults(
  defineProps<{
    /** Material Symbols ligature name. */
    icon: string;
    /** Accessible name — rendered as aria-label, never as visible text. */
    label: string;
    variant?: "standard" | "filled" | "tonal" | "outlined";
    tone?: "default" | "primary" | "danger";
    size?: "sm" | "md" | "lg";
    disabled?: boolean;
    /** Fills the icon (variation axis; best-effort — see OrIcon). */
    fill?: boolean;
  }>(),
  {
    variant: "standard",
    tone: "default",
    size: "md",
    disabled: false,
    fill: false,
  },
);

defineEmits<{ click: [event: MouseEvent] }>();
</script>

<template>
  <button
    v-ripple
    class="or-icon-btn or-interactive"
    :class="[`or-icon-btn--${variant}`, `or-icon-btn--${tone}`, `or-icon-btn--${size}`]"
    type="button"
    :aria-label="label"
    :title="label"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <OrIcon :name="icon" :size="size === 'sm' ? 20 : size === 'lg' ? 26 : 24" :fill="fill" />
  </button>
</template>

<style scoped>
.or-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
  transition:
    border-radius var(--or-motion-spatial-fast),
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

.or-icon-btn:active:not(:disabled) {
  border-radius: var(--or-morph-button-active);
}

.or-icon-btn:disabled {
  opacity: 0.38;
}

.or-icon-btn--sm {
  width: 40px;
  height: 40px;
}
.or-icon-btn--md {
  width: 48px;
  height: 48px;
}
.or-icon-btn--lg {
  width: 56px;
  height: 56px;
}

.or-icon-btn--standard {
  background: transparent;
}
.or-icon-btn--standard:hover:not(:disabled) {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
}

.or-icon-btn--filled.or-icon-btn--primary {
  background: var(--m3-color-primary);
  color: var(--m3-color-on-primary);
}
.or-icon-btn--filled.or-icon-btn--danger {
  background: var(--m3-color-error);
  color: var(--m3-color-on-error);
}
.or-icon-btn--filled.or-icon-btn--default {
  background: var(--m3-color-primary);
  color: var(--m3-color-on-primary);
}

.or-icon-btn--tonal.or-icon-btn--danger {
  background: var(--m3-color-error-container);
  color: var(--m3-color-on-error-container);
}
.or-icon-btn--tonal.or-icon-btn--default,
.or-icon-btn--tonal.or-icon-btn--primary {
  background: var(--m3-color-secondary-container);
  color: var(--m3-color-on-secondary-container);
}

.or-icon-btn--outlined {
  border: 1px solid var(--m3-color-outline);
}
.or-icon-btn--outlined.or-icon-btn--danger {
  color: var(--m3-color-error);
}
.or-icon-btn--outlined.or-icon-btn--default,
.or-icon-btn--outlined.or-icon-btn--primary {
  color: var(--m3-color-primary);
}
</style>

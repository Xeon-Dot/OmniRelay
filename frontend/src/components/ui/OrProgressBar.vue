<script setup lang="ts">
withDefaults(
  defineProps<{
    /** 0–1. Omit for an indeterminate bar. */
    value?: number;
    tone?: "primary" | "success" | "warning" | "danger";
    /** Accessible name; required so the bar means something to a reader. */
    label?: string;
  }>(),
  { value: undefined, tone: "primary", label: undefined },
);
</script>

<template>
  <div
    class="or-progress"
    :class="`or-progress--${tone}`"
    role="progressbar"
    :aria-label="label"
    :aria-valuenow="value === undefined ? undefined : Math.round(value * 100)"
    :aria-valuemin="value === undefined ? undefined : 0"
    :aria-valuemax="value === undefined ? undefined : 100"
  >
    <div class="or-progress__track">
      <div
        class="or-progress__fill"
        :class="{ 'or-progress__fill--indeterminate': value === undefined }"
        :style="value === undefined ? undefined : { inlineSize: `${Math.min(100, Math.max(0, value * 100))}%` }"
      />
    </div>
  </div>
</template>

<style scoped>
.or-progress {
  inline-size: 100%;
}

.or-progress__track {
  block-size: 4px;
  inline-size: 100%;
  border-radius: var(--m3-shape-full);
  background: color-mix(in srgb, var(--m3-color-on-surface) 12%, transparent);
  overflow: hidden;
}

.or-progress__fill {
  block-size: 100%;
  border-radius: var(--m3-shape-full);
  background: var(--m3-color-primary);
  transition: inline-size var(--or-motion-spatial-default);
}

.or-progress--success .or-progress__fill {
  background: var(--or-success);
}
.or-progress--warning .or-progress__fill {
  background: var(--or-warning);
}
.or-progress--danger .or-progress__fill {
  background: var(--m3-color-error);
}

.or-progress__fill--indeterminate {
  inline-size: 40%;
  animation: or-progress-slide 1400ms var(--m3-easing-standard) infinite;
}

@keyframes or-progress-slide {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(250%);
  }
}
</style>

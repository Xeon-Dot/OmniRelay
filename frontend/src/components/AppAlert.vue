<template>
  <div class="alert" :class="[`alert--${variant}`, { 'alert--page': page }]" role="alert">
    <OrIcon class="alert__icon" :name="icon" :size="18" />
    <span class="alert__body"><slot /></span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    variant?: "success" | "error" | "warning";
    page?: boolean;
  }>(),
  { variant: "error", page: false },
);

const icon = computed(() =>
  props.variant === "success"
    ? "check_circle"
    : props.variant === "warning"
      ? "warning"
      : "error",
);
</script>

<style scoped>
.alert {
  display: flex;
  align-items: flex-start;
  gap: var(--m3-space-150);
  padding: var(--m3-space-150) var(--m3-space-200);
  border-radius: var(--m3-shape-md);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}

.alert__icon {
  margin-block-start: 1px;
}

.alert__body {
  min-inline-size: 0;
}

.alert--error {
  background: var(--m3-color-error-container);
  color: var(--m3-color-on-error-container);
}

.alert--success {
  background: var(--or-success-container);
  color: var(--or-on-success-container);
}

.alert--warning {
  background: var(--or-warning-container);
  color: var(--or-on-warning-container);
}

.alert--page {
  margin-block-end: var(--m3-space-200);
}
</style>

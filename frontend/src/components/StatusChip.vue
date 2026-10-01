<template>
  <span class="status-chip" :class="`status-chip--${variant}`">
    <slot />
  </span>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ variant?: "on" | "off" | "warning" }>(), {
  variant: "on",
});
</script>

<style scoped>
.status-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-75, 6px);
  min-block-size: 28px;
  padding-inline: var(--m3-space-150);
  border-radius: var(--m3-shape-full);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  white-space: nowrap;
}

.status-chip--on {
  background: var(--or-success-container);
  color: var(--or-on-success-container);
}

.status-chip--off {
  background: var(--m3-color-error-container);
  color: var(--m3-color-on-error-container);
}

.status-chip--warning {
  background: var(--or-warning-container);
  color: var(--or-on-warning-container);
}
</style>

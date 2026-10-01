<template>
  <article class="stat-card">
    <p class="stat-label">
      {{ label }}
      <span v-if="hint" class="stat-hint">{{ hint }}</span>
    </p>
    <p class="stat-value" :class="valueClass">
      <slot />
    </p>
    <p v-if="$slots.sub || sub" class="stat-sub">
      <slot name="sub">{{ sub }}</slot>
    </p>
  </article>
</template>

<script setup lang="ts">
defineProps<{
  label: string;
  hint?: string;
  sub?: string;
  valueClass?: string;
}>();
</script>

<style scoped>
.stat-card {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-75, 6px);
  min-inline-size: 0;
  padding: var(--or-pad-card);
  background: var(--m3-elevation-1-surface);
  border: 1px solid var(--m3-color-outline-variant);
  border-radius: var(--m3-shape-2xl);
}

.stat-label {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m3-space-100);
  margin: 0;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}

.stat-hint {
  color: var(--m3-color-outline);
}

.stat-value {
  margin: 0;
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-emphasized-display-small);
  letter-spacing: var(--m3-typescale-emphasized-display-small-tracking);
  font-variant-numeric: tabular-nums;
  /* A metric is one token — "4,567,890" must not split across two lines. */
  white-space: nowrap;
}

.stat-sub {
  margin: 0;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-small);
  letter-spacing: var(--m3-typescale-body-small-tracking);
}
</style>

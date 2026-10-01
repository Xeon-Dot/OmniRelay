<template>
  <div class="mobile-card">
    <div v-for="(item, idx) in items" :key="idx" class="mobile-card__row">
      <span class="mobile-card__label">{{ item.label }}</span>
      <span class="mobile-card__value">{{ item.value }}</span>
    </div>
    <div v-if="$slots.actions" class="mobile-card__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface CardItem {
  label: string;
  value: string;
}

defineProps<{
  items: CardItem[];
}>();
</script>

<style scoped>
.mobile-card {
  background: var(--m3-elevation-2-surface);
  border: 1px solid var(--m3-color-outline-variant);
  border-radius: var(--m3-shape-xl);
  padding: var(--m3-space-200);
}

.mobile-card__row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m3-space-200);
  padding-block: var(--m3-space-75, 6px);
}

.mobile-card__row + .mobile-card__row {
  border-block-start: 1px solid var(--m3-color-outline-variant);
}

.mobile-card__label {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}

.mobile-card__value {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
  text-align: end;
  overflow-wrap: anywhere;
}

.mobile-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--m3-space-100);
  margin-block-start: var(--m3-space-100);
  padding-block-start: var(--m3-space-100);
  border-block-start: 1px solid var(--m3-color-outline-variant);
}
</style>

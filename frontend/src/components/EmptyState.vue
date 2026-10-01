<template>
  <div class="empty-state" :class="{ small }" role="status" aria-live="polite">
    <OrIcon :name="icon" :size="iconSize" class="empty-state__icon" />
    <p>{{ text }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    /** Material Symbols ligature (not an MDI name). */
    icon: string;
    text: string;
    iconSize?: number | string;
    small?: boolean;
  }>(),
  { iconSize: 32, small: false },
);
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--m3-space-150);
  padding: var(--m3-space-600) var(--m3-space-300);
  color: var(--m3-color-on-surface-variant);
  text-align: center;
}

.empty-state__icon {
  color: var(--m3-color-outline);
}

.empty-state p {
  margin: 0;
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}

.empty-state.small {
  padding: var(--m3-space-400) var(--m3-space-200);
}
</style>

<template>
  <header class="page-header">
    <div class="page-header__text">
      <h1 class="page-title">{{ title }}</h1>
      <p class="page-sub">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.default" class="page-header__actions">
      <slot />
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{ title: string; subtitle: string }>();
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--m3-space-200);
  flex-wrap: wrap;
}

.page-header__text {
  min-inline-size: 0;
}

.page-header__actions {
  display: flex;
  align-items: center;
  gap: var(--m3-space-150);
  flex-wrap: wrap;
}

.page-title {
  margin: 0 0 var(--m3-space-75, 6px);
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-emphasized-headline-small);
  letter-spacing: var(--m3-typescale-emphasized-headline-small-tracking);
  text-wrap: balance;
}

.page-sub {
  margin: 0;
  max-inline-size: 72ch;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
  text-wrap: pretty;
}
</style>

<script setup lang="ts">
withDefaults(defineProps<{ vertical?: boolean }>(), { vertical: false });
</script>

<template>
  <hr class="or-divider" :class="{ 'or-divider--vertical': vertical }" />
</template>

<style scoped>
.or-divider {
  border: none;
  flex-shrink: 0;
  background: var(--m3-color-outline-variant);
  margin: 0;
}

.or-divider:not(.or-divider--vertical) {
  width: 100%;
  height: 1px;
}

.or-divider--vertical {
  width: 1px;
  height: 1em;
  min-height: 24px;
  align-self: stretch;
}
</style>

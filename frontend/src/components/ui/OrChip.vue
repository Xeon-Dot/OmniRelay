<script setup lang="ts">
/**
 * Chip: assist / filter. Selected filter chips fill with the primary
 * container and morph from pill to a tighter corner — the M3 Expressive
 * selection cue.
 */
withDefaults(
  defineProps<{
    variant?: "assist" | "filter";
    selected?: boolean;
    disabled?: boolean;
    /** Leading Material Symbols ligature. */
    icon?: string;
  }>(),
  { variant: "assist", selected: false, disabled: false, icon: undefined },
);

defineEmits<{ click: [event: MouseEvent] }>();
</script>

<template>
  <button
    v-ripple
    class="or-chip or-interactive"
    :class="[`or-chip--${variant}`, { 'is-selected': selected }]"
    type="button"
    :disabled="disabled"
    :aria-pressed="variant === 'filter' ? selected : undefined"
    @click="$emit('click', $event)"
  >
    <OrIcon v-if="icon" class="or-chip__icon" :name="icon" :size="18" />
    <span class="or-chip__label"><slot /></span>
  </button>
</template>

<style scoped>
.or-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-100);
  min-height: 32px;
  padding: 0 var(--m3-space-200);
  border: 1px solid var(--m3-color-outline);
  border-radius: var(--or-morph-chip-rest);
  background: transparent;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
  white-space: nowrap;
  transition:
    border-radius var(--or-motion-spatial-fast),
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast),
    border-color var(--or-motion-effects-fast);
}

.or-chip:disabled {
  opacity: 0.38;
}

.or-chip__icon {
  margin-inline-start: calc(var(--m3-space-100) * -1 + 2px);
}

/* Selected: filled + squarer. */
.or-chip.is-selected {
  background: var(--m3-color-primary-container);
  border-color: transparent;
  color: var(--m3-color-on-primary-container);
  border-radius: var(--or-morph-chip-active);
}
</style>

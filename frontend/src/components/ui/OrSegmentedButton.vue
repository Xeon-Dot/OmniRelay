<script setup lang="ts">
/**
 * Segmented button — a connected row of mutually exclusive choices (the
 * Performance view's 1H/24H/7D/30D presets). Renders as one pill split by
 * dividers; the selected segment fills and morphs inward.
 */
import { computed } from "vue";

type SegmentedOption = { value: string; label?: string; icon?: string };

const props = withDefaults(
  defineProps<{
    options: SegmentedOption[];
    modelValue?: string | string[];
    multiple?: boolean;
    disabled?: boolean;
    size?: "sm" | "md";
    /** Stretch to the container and divide the segments evenly (M3 full-width
     *  segmented button). Used by the drawer footer, where the pill would
     *  otherwise sit in a stretched box with the segments bunched at the start. */
    block?: boolean;
  }>(),
  { modelValue: "", multiple: false, disabled: false, size: "md", block: false },
);

const emit = defineEmits<{ "update:modelValue": [value: string | string[]] }>();

const selected = computed(() =>
  Array.isArray(props.modelValue) ? props.modelValue : [props.modelValue],
);

function isSelected(value: string) {
  return selected.value.includes(value);
}

function select(value: string) {
  if (props.disabled) return;
  if (!props.multiple) {
    emit("update:modelValue", value);
    return;
  }
  const next = new Set(selected.value);
  if (next.has(value)) next.delete(value);
  else next.add(value);
  emit("update:modelValue", [...next]);
}
</script>

<template>
  <div
    class="or-segmented"
    :class="[`or-segmented--${size}`, { 'is-disabled': disabled, 'or-segmented--block': block }]"
    role="group"
  >
    <button
      v-for="(option, index) in options"
      :key="option.value"
      v-ripple
      class="or-segmented__item or-interactive"
      :class="{ 'is-selected': isSelected(option.value) }"
      type="button"
      :disabled="disabled"
      :aria-pressed="isSelected(option.value)"
      :data-first="index === 0 ? '' : undefined"
      :data-last="index === options.length - 1 ? '' : undefined"
      @click="select(option.value)"
    >
      <OrIcon v-if="option.icon" :name="option.icon" :size="18" />
      <span v-if="option.label">{{ option.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.or-segmented {
  display: inline-flex;
  align-items: stretch;
  border: 1px solid var(--m3-color-outline);
  border-radius: var(--m3-shape-full);
  overflow: hidden;
  background: transparent;
}

.or-segmented__item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--m3-space-100);
  min-height: 40px;
  padding-inline: var(--m3-space-200);
  background: transparent;
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
  border-radius: 0;
  /* A segment label is a single word pair ("시스템"); never break it. */
  white-space: nowrap;
  transition:
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

/* Full-width variant: fill the container and share the space evenly. */
.or-segmented--block {
  display: flex;
  inline-size: 100%;
}

.or-segmented--block .or-segmented__item {
  flex: 1 1 0;
}

.or-segmented--sm .or-segmented__item {
  min-height: 32px;
  padding-inline: var(--m3-space-150);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}

.or-segmented__item + .or-segmented__item {
  border-inline-start: 1px solid var(--m3-color-outline);
}

.or-segmented__item.is-selected {
  background: var(--m3-color-primary);
  color: var(--m3-color-on-primary);
}

.or-segmented.is-disabled {
  opacity: 0.38;
}
</style>

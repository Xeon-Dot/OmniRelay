<script setup lang="ts">
/**
 * M3 switch. The track is a pill; the thumb shrinks and gains the primary
 * colour as it travels — a shape morph, not just a colour swap.
 */
import { useId } from "vue";

withDefaults(
  defineProps<{
    modelValue?: boolean;
    label?: string;
    disabled?: boolean;
  }>(),
  { modelValue: false, label: undefined, disabled: false },
);

const emit = defineEmits<{ "update:modelValue": [value: boolean] }>();
const id = useId();

function onChange(event: Event) {
  emit("update:modelValue", (event.target as HTMLInputElement).checked);
}
</script>

<template>
  <div class="or-switch" :class="{ 'is-disabled': disabled }">
    <input
      :id="id"
      class="or-switch__native"
      type="checkbox"
      role="switch"
      :checked="modelValue"
      :disabled="disabled"
      :aria-label="label"
      @change="onChange"
    />
    <label class="or-switch__label" :for="id">
      <span class="or-switch__track" aria-hidden="true">
        <span class="or-switch__thumb">
          <OrIcon class="or-switch__on-icon" name="check" :size="14" />
        </span>
      </span>
      <span v-if="label" class="or-switch__text">{{ label }}</span>
    </label>
  </div>
</template>

<style scoped>
.or-switch {
  position: relative;
}

.or-switch__native {
  position: absolute;
  inset-block-start: 50%;
  inset-inline-start: 0;
  inline-size: 52px;
  block-size: 32px;
  margin: 0;
  opacity: 0;
  transform: translateY(-50%);
  cursor: pointer;
}

.or-switch__label {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-200);
  min-block-size: var(--m3-min-touch-target);
  cursor: pointer;
}

.or-switch__track {
  position: relative;
  display: inline-flex;
  align-items: center;
  inline-size: 52px;
  block-size: 32px;
  padding: 4px;
  border: 2px solid var(--m3-color-outline);
  border-radius: var(--m3-shape-full);
  background: var(--m3-color-surface-container-highest);
  transition:
    background-color var(--or-motion-effects-fast),
    border-color var(--or-motion-effects-fast);
}

.or-switch__thumb {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 24px;
  block-size: 24px;
  border-radius: var(--m3-shape-full);
  background: var(--m3-color-outline);
  color: var(--m3-color-on-primary);
  transform: translateX(0);
  transition:
    transform var(--or-motion-spatial-fast),
    inline-size var(--or-motion-spatial-fast),
    block-size var(--or-motion-spatial-fast),
    background-color var(--or-motion-effects-fast);
}

.or-switch__on-icon {
  opacity: 0;
  transition: opacity var(--or-motion-effects-fast);
}

.or-switch__native:checked + .or-switch__label .or-switch__track {
  background: var(--m3-color-primary);
  border-color: var(--m3-color-primary);
}

.or-switch__native:checked + .or-switch__label .or-switch__thumb {
  background: var(--m3-color-on-primary);
  transform: translateX(20px);
  inline-size: 16px;
  block-size: 16px;
}

.or-switch__native:checked + .or-switch__label .or-switch__on-icon {
  opacity: 1;
}

.or-switch__native:focus-visible + .or-switch__label .or-switch__track {
  outline: var(--m3-focus-ring-width) solid var(--m3-focus-ring-color);
  outline-offset: var(--m3-focus-ring-offset);
}

.or-switch__text {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-body-large);
  letter-spacing: var(--m3-typescale-body-large-tracking);
}

.or-switch.is-disabled {
  opacity: 0.38;
  pointer-events: none;
}
</style>

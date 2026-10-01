<script setup lang="ts">
/**
 * M3 outlined select. Keeps a native `<select>` so mobile gets the platform
 * picker and keyboard/screen-reader behaviour for free; the M3 shell and
 * floating label are decoration over it.
 */
import { useId } from "vue";

withDefaults(
  defineProps<{
    modelValue?: string;
    label: string;
    options: { value: string; label: string }[];
    hint?: string;
    error?: string;
    disabled?: boolean;
    /** Adds an empty leading option so nothing appears selected by default. */
    placeholder?: string;
  }>(),
  {
    modelValue: "",
    hint: undefined,
    error: undefined,
    disabled: false,
    placeholder: undefined,
  },
);

const emit = defineEmits<{ "update:modelValue": [value: string] }>();
const id = useId();

function onChange(event: Event) {
  emit("update:modelValue", (event.target as HTMLSelectElement).value);
}
</script>

<template>
  <div class="or-select" :class="{ 'is-disabled': disabled, 'is-error': Boolean(error) }">
    <div class="or-select__box">
      <select
        :id="id"
        class="or-select__control"
        :value="modelValue"
        :disabled="disabled"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="error || hint ? `${id}-support` : undefined"
        @change="onChange"
      >
        <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
        <option v-for="option in options" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
      <label class="or-select__label" :for="id">{{ label }}</label>
      <OrIcon class="or-select__arrow" name="expand_more" :size="24" />
    </div>
    <p :id="`${id}-support`" class="or-select__support" :class="{ 'is-error': Boolean(error) }">
      {{ error || hint }}
    </p>
  </div>
</template>

<style scoped>
.or-select {
  display: flex;
  flex-direction: column;
  min-inline-size: 0;
}

.or-select__box {
  position: relative;
  display: flex;
  align-items: center;
  min-block-size: 56px;
  padding-inline: var(--m3-space-200) var(--m3-space-100);
  border: 1px solid var(--m3-color-outline);
  border-radius: var(--or-morph-field-rest);
  background: var(--or-field-surface, var(--m3-color-surface));
  transition:
    border-radius var(--or-motion-spatial-fast),
    border-color var(--or-motion-effects-fast);
}

.or-select__box:focus-within {
  border-color: var(--m3-color-primary);
  border-width: 2px;
  border-radius: var(--or-morph-field-active);
}

.or-select.is-error .or-select__box:focus-within {
  border-color: var(--m3-color-error);
}

.or-select__control {
  flex: 1;
  min-inline-size: 0;
  appearance: none;
  border: none;
  outline: none;
  background: transparent;
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-body-large);
  letter-spacing: var(--m3-typescale-body-large-tracking);
  padding-block: var(--m3-space-300);
  cursor: pointer;
}

/* The label floats permanently — a select always shows a value. */
.or-select__label {
  position: absolute;
  inset-inline-start: calc(var(--m3-space-200) + 4px);
  inset-block-start: 0;
  padding-inline: var(--m3-space-25, 2px);
  background: var(--or-field-surface, var(--m3-color-surface));
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  transform: translateY(-50%);
  pointer-events: none;
}

.or-select__box:focus-within .or-select__label {
  color: var(--m3-color-primary);
}

.or-select__arrow {
  color: var(--m3-color-on-surface-variant);
  pointer-events: none;
}

.or-select__support {
  min-block-size: 1.25rem;
  margin: var(--m3-space-75, 6px) var(--m3-space-200) 0;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-small);
  letter-spacing: var(--m3-typescale-body-small-tracking);
}

.or-select__support.is-error {
  color: var(--m3-color-error);
}

.or-select.is-disabled {
  opacity: 0.38;
}
</style>

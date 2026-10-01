<script setup lang="ts">
/**
 * Checkbox with an M3 filled container. Supports both the plain boolean case
 * and the multi-select case (pass `value` and an array as `modelValue`), which
 * is what the provider source-model list and the user provider checklist use.
 */
import { computed, useId } from "vue";

const props = withDefaults(
  defineProps<{
    /** boolean for a standalone box, array of values when `value` is given. */
    modelValue?: boolean | unknown[];
    /** Value contributed to the array in multi-select mode. */
    value?: unknown;
    label?: string;
    hint?: string;
    disabled?: boolean;
    /** Renders the box in its error colour (used by form validation). */
    error?: boolean;
  }>(),
  {
    modelValue: false,
    value: undefined,
    label: undefined,
    hint: undefined,
    disabled: false,
    error: false,
  },
);

const emit = defineEmits<{ "update:modelValue": [value: boolean | unknown[]] }>();
const id = useId();

const checked = computed(() =>
  Array.isArray(props.modelValue) && props.value !== undefined
    ? props.modelValue.includes(props.value)
    : Boolean(props.modelValue),
);

function toggle() {
  if (props.disabled) return;
  if (Array.isArray(props.modelValue) && props.value !== undefined) {
    const next = props.modelValue.filter((item) => item !== props.value);
    if (!checked.value) next.push(props.value);
    emit("update:modelValue", next);
    return;
  }
  emit("update:modelValue", !checked.value);
}
</script>

<template>
  <div class="or-checkbox" :class="{ 'is-disabled': disabled, 'is-error': error }">
    <input
      :id="id"
      class="or-checkbox__native"
      type="checkbox"
      :checked="checked"
      :disabled="disabled"
      :value="value === undefined ? undefined : String(value)"
      @change="toggle"
    />
    <label class="or-checkbox__label" :for="id">
      <span class="or-checkbox__box" aria-hidden="true">
        <OrIcon name="check" :size="18" class="or-checkbox__mark" />
      </span>
      <span class="or-checkbox__text">
        <span v-if="label" class="or-checkbox__title">{{ label }}</span>
        <span v-if="hint" class="or-checkbox__hint">{{ hint }}</span>
      </span>
    </label>
  </div>
</template>

<style scoped>
.or-checkbox {
  position: relative;
}

/* Native input stays in the DOM for form semantics and keyboard support;
     it is visually replaced by `.or-checkbox__box`. */
.or-checkbox__native {
  position: absolute;
  inset-block-start: 50%;
  inset-inline-start: 0;
  inline-size: 18px;
  block-size: 18px;
  margin: 0;
  opacity: 0;
  transform: translateY(-50%);
  cursor: pointer;
}

.or-checkbox__label {
  display: flex;
  align-items: flex-start;
  gap: var(--m3-space-200);
  min-block-size: var(--m3-min-touch-target);
  padding-block: var(--m3-space-150);
  cursor: pointer;
}

.or-checkbox__box {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  inline-size: 18px;
  block-size: 18px;
  margin-block-start: 2px;
  border: 2px solid var(--m3-color-on-surface-variant);
  border-radius: var(--m3-shape-xs);
  background: transparent;
  color: var(--m3-color-on-primary);
  transition:
    background-color var(--or-motion-effects-fast),
    border-color var(--or-motion-effects-fast),
    border-radius var(--or-motion-spatial-fast);
}

.or-checkbox__mark {
  opacity: 0;
  transform: scale(0.4);
  transition:
    opacity var(--or-motion-effects-fast),
    transform var(--or-motion-spatial-fast);
}

.or-checkbox__native:checked + .or-checkbox__label .or-checkbox__box {
  background: var(--m3-color-primary);
  border-color: var(--m3-color-primary);
}

.or-checkbox__native:checked + .or-checkbox__label .or-checkbox__mark {
  opacity: 1;
  transform: scale(1);
}

.or-checkbox__native:focus-visible + .or-checkbox__label .or-checkbox__box {
  outline: var(--m3-focus-ring-width) solid var(--m3-focus-ring-color);
  outline-offset: var(--m3-focus-ring-offset);
}

.or-checkbox.is-error .or-checkbox__box {
  border-color: var(--m3-color-error);
}

.or-checkbox__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-inline-size: 0;
}

.or-checkbox__title {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-body-large);
  letter-spacing: var(--m3-typescale-body-large-tracking);
}

.or-checkbox__hint {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-small);
  letter-spacing: var(--m3-typescale-body-small-tracking);
}

.or-checkbox.is-disabled {
  opacity: 0.38;
  pointer-events: none;
}
</style>

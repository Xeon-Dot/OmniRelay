<script setup lang="ts">
/**
 * M3 outlined text field with a floating label.
 *
 * `type` passes straight through, so date/password/number fields keep the
 * browser's own picker and validation. Support text sits below the box and
 * reserves its height, so showing an error never shifts the form.
 *
 * The floating label's notch paints `--or-field-surface`, which defaults to
 * `--m3-color-surface`. Set it on a card (or any tinted container) so the
 * notch matches the surface the field actually sits on.
 */
import { computed, useId } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue?: string | number;
    label: string;
    hint?: string;
    /** Non-empty switches the field to its error state and replaces `hint`. */
    error?: string;
    type?: string;
    disabled?: boolean;
    readonly?: boolean;
    required?: boolean;
    leadingIcon?: string;
    clearable?: boolean;
    placeholder?: string;
    autocomplete?: string;
    name?: string;
    min?: number | string;
    max?: number | string;
    step?: number | string;
  }>(),
  {
    modelValue: "",
    hint: undefined,
    error: undefined,
    type: "text",
    disabled: false,
    readonly: false,
    required: false,
    leadingIcon: undefined,
    clearable: false,
    placeholder: undefined,
    autocomplete: undefined,
    name: undefined,
    min: undefined,
    max: undefined,
    step: undefined,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
  clear: [];
}>();

const id = useId();

const filled = computed(
  () => String(props.modelValue ?? "").length > 0 || props.type === "date",
);

function onInput(event: Event) {
  emit("update:modelValue", (event.target as HTMLInputElement).value);
}

function onClear() {
  emit("update:modelValue", "");
  emit("clear");
}
</script>

<template>
  <div
    class="or-field"
    :class="{ 'is-disabled': disabled, 'is-error': Boolean(error) }"
  >
    <div class="or-field__box" :class="{ 'is-filled': filled }">
      <OrIcon v-if="leadingIcon" class="or-field__leading" :name="leadingIcon" :size="20" />
      <input
        :id="id"
        class="or-field__input"
        :type="type"
        :value="modelValue"
        :name="name"
        :placeholder="placeholder"
        :autocomplete="autocomplete"
        :disabled="disabled"
        :readonly="readonly"
        :required="required"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="error || hint ? `${id}-support` : undefined"
        :min="min"
        :max="max"
        :step="step"
        @input="onInput"
      />
      <label class="or-field__label" :for="id">{{ label }}</label>
      <button
        v-if="clearable && !disabled"
        class="or-field__clear"
        type="button"
        :aria-label="`${label} — ${$t('common.clear')}`"
        @click="onClear"
      >
        <OrIcon name="close" :size="18" />
      </button>
      <span v-if="$slots.trailing" class="or-field__trailing">
        <slot name="trailing" />
      </span>
    </div>
    <p :id="`${id}-support`" class="or-field__support" :class="{ 'is-error': Boolean(error) }">
      {{ error || hint }}
    </p>
  </div>
</template>

<style scoped>
.or-field {
  display: flex;
  flex-direction: column;
  min-inline-size: 0;
}

.or-field__box {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
  min-block-size: 56px;
  padding-inline: var(--m3-space-200);
  border: 1px solid var(--m3-color-outline);
  border-radius: var(--or-morph-field-rest);
  background: var(--or-field-surface, var(--m3-color-surface));
  color: var(--m3-color-on-surface);
  transition:
    border-radius var(--or-motion-spatial-fast),
    border-color var(--or-motion-effects-fast),
    box-shadow var(--or-motion-effects-fast);
}

.or-field__box:focus-within {
  border-color: var(--m3-color-primary);
  border-width: 2px;
  border-radius: var(--or-morph-field-active);
  padding-inline: calc(var(--m3-space-200) - 1px);
}

.or-field.is-error .or-field__box:focus-within {
  border-color: var(--m3-color-error);
}

.or-field__leading {
  color: var(--m3-color-on-surface-variant);
}

.or-field__input {
  flex: 1;
  min-inline-size: 0;
  border: none;
  outline: none;
  background: transparent;
  color: inherit;
  font: var(--m3-typescale-body-large);
  letter-spacing: var(--m3-typescale-body-large-tracking);
  padding-block: var(--m3-space-200);
}

.or-field__input::placeholder {
  color: transparent;
}

/* ── Floating label ── */
.or-field__label {
  position: absolute;
  inset-inline-start: calc(var(--m3-space-200) + 4px);
  inset-block-start: 50%;
  padding-inline: var(--m3-space-25, 2px);
  background: var(--or-field-surface, var(--m3-color-surface));
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-large);
  letter-spacing: var(--m3-typescale-body-large-tracking);
  transform: translateY(-50%);
  transform-origin: left center;
  pointer-events: none;
  transition:
    inset-block-start var(--or-motion-effects-fast),
    transform var(--or-motion-effects-fast),
    font-size var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

.or-field__leading ~ .or-field__label {
  inset-inline-start: calc(var(--m3-space-200) + 28px);
}

.or-field__box.is-filled .or-field__label,
.or-field__box:focus-within .or-field__label {
  inset-block-start: 0;
  transform: translateY(-50%);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}

.or-field__box:focus-within .or-field__label {
  color: var(--m3-color-primary);
}

.or-field.is-error .or-field__box:focus-within .or-field__label {
  color: var(--m3-color-error);
}

.or-field__clear {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 32px;
  block-size: 32px;
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
}

.or-field__clear:hover {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
}

.or-field__trailing {
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
}

.or-field__support {
  min-block-size: 1.25rem;
  margin: var(--m3-space-75, 6px) var(--m3-space-200) 0;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-small);
  letter-spacing: var(--m3-typescale-body-small-tracking);
}

.or-field__support.is-error {
  color: var(--m3-color-error);
}

.or-field.is-disabled {
  opacity: 0.38;
}
</style>

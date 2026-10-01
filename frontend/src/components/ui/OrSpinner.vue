<script setup lang="ts">
import { computed } from "vue";

/**
 * Indeterminate circular progress. The border morphs between a circle and a
 * squircle while it spins — M3 Expressive's shape-morphing loader. The global
 * reduced-motion gate stops it at its end state.
 */
const props = withDefaults(
  defineProps<{
    /** Diameter in px. */
    size?: number;
    /** Accessible name; omit when the spinner sits inside a labelled control. */
    label?: string;
  }>(),
  { size: 20 },
);

const style = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  borderWidth: `${Math.max(2, Math.round(props.size / 10))}px`,
}));
</script>

<template>
  <span
    class="or-spinner"
    :style="style"
    :role="label ? 'status' : undefined"
    :aria-label="label"
    :aria-hidden="label ? undefined : 'true'"
  />
</template>

<style scoped>
.or-spinner {
  display: inline-block;
  flex-shrink: 0;
  box-sizing: border-box;
  border-style: solid;
  border-color: color-mix(in srgb, currentColor 24%, transparent);
  border-top-color: currentColor;
  animation:
    or-spin var(--m3-duration-extra-long1) var(--m3-easing-linear) infinite,
    or-spin-morph 2400ms var(--m3-easing-standard) infinite;
}

@keyframes or-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes or-spin-morph {
  0%,
  100% {
    border-radius: 50%;
  }
  30% {
    border-radius: 38% 50% 50% 38%;
  }
  60% {
    border-radius: 50% 38% 50% 50%;
  }
}
</style>

<script setup lang="ts">
/**
 * Material Symbols icon (Rounded — matches the M3 Expressive shape language).
 * Replaces every legacy MDI icon element in the app.
 *
 * Only one weight is bundled (`@material-symbols/font-400`), so emphasis is
 * carried by size and colour, not weight. `fill` is best-effort: it sets the
 * FILL variation axis, which the bundled font may not expose — active states
 * never depend on it, they also paint a container behind the icon.
 */
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    /** Ligature name, e.g. "settings", "arrow_back". */
    name: string;
    /** CSS length or number (px). */
    size?: number | string;
    /** Any CSS colour. Defaults to inherited text colour. */
    color?: string;
    /** Filled variant (variation axis; best-effort). */
    fill?: boolean;
    /**
     * Accessible name. Omit for decorative icons — they stay aria-hidden.
     * Provide it only when the icon is the sole carrier of meaning.
     */
    label?: string;
  }>(),
  { size: 24, fill: false },
);

const style = computed(() => ({
  fontSize: typeof props.size === "number" ? `${props.size}px` : props.size,
  ...(props.color ? { color: props.color } : {}),
  ...(props.fill ? { fontVariationSettings: "'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 48" } : {}),
}));

const decorative = computed(() => !props.label);
</script>

<template>
  <span
    class="or-icon material-symbols-rounded"
    :style="style"
    :aria-hidden="decorative ? 'true' : undefined"
    :role="decorative ? undefined : 'img'"
    :aria-label="label || undefined"
  >
    {{ name }}
  </span>
</template>

<style scoped>
.or-icon {
  display: inline-block;
  flex-shrink: 0;
  line-height: 1;
  user-select: none;
  font-weight: normal;
  font-style: normal;
  letter-spacing: normal;
  text-transform: none;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  -webkit-font-smoothing: antialiased;
}
</style>

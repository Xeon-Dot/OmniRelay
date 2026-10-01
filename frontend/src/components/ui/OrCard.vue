<script setup lang="ts">
/**
 * Surface container. Elevation is expressed as *tone* (which surface-container
 * role it sits on), not a drop shadow — shadows are reserved for transient
 * overlays. `interactive` cards morph from corner-2xl to corner-lg when pressed.
 */
withDefaults(
  defineProps<{
    variant?: "elevated" | "filled" | "outlined";
    interactive?: boolean;
    /** Removes internal padding, for cards that manage their own layout. */
    flush?: boolean;
  }>(),
  { variant: "elevated", interactive: false, flush: false },
);

defineEmits<{ click: [event: MouseEvent] }>();
</script>

<template>
  <component
    :is="interactive ? 'button' : 'div'"
    v-ripple="interactive"
    class="or-card"
    :class="[
      `or-card--${variant}`,
      { 'or-card--interactive or-interactive': interactive, 'or-card--flush': flush },
    ]"
    :type="interactive ? 'button' : undefined"
    @click="$emit('click', $event)"
  >
    <slot />
  </component>
</template>

<style scoped>
.or-card {
  display: block;
  width: 100%;
  text-align: start;
  border-radius: var(--or-morph-card-rest);
  background: var(--m3-elevation-2-surface);
  color: var(--m3-color-on-surface);
  transition:
    border-radius var(--or-motion-spatial-default),
    background-color var(--or-motion-effects-default),
    box-shadow var(--or-motion-effects-default);
}

.or-card:not(.or-card--flush) {
  padding: var(--or-pad-card);
}

.or-card--elevated {
  background: var(--m3-elevation-2-surface);
}

.or-card--filled {
  background: var(--m3-color-surface-container-highest);
}

.or-card--outlined {
  background: var(--m3-color-surface);
  border: 1px solid var(--m3-color-outline-variant);
}

.or-card--interactive {
  cursor: pointer;
  border: none;
  font: inherit;
}

.or-card--interactive:hover:not(:disabled) {
  background: var(--m3-elevation-3-surface);
}

.or-card--interactive:active:not(:disabled) {
  border-radius: var(--or-morph-card-active);
}
</style>

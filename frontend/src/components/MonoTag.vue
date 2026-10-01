<template>
  <code class="mono-tag"><slot /></code>
</template>

<style scoped>
.mono-tag {
  display: inline-block;
  padding: 2px var(--m3-space-100);
  border-radius: var(--m3-shape-xs);
  background: color-mix(in srgb, var(--m3-color-on-surface) 6%, transparent);
  color: var(--m3-color-on-surface);
  font-family: var(--m3-typeface-mono);
  font-size: var(--m3-space-200);
  line-height: 1.45;
  /* Model ids stay on one line; wide tables pan horizontally. */
  white-space: nowrap;
}
</style>

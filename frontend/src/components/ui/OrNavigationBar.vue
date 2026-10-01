<script setup lang="ts">
/**
 * Mobile bottom navigation bar (M3 navigation bar), shown below 600px.
 *
 * 80px tall plus the device safe-area inset. The active destination fills
 * with `secondary-container` behind the icon — the M3 pattern — and the bar
 * keeps every destination reachable with one thumb.
 */
import { computed } from "vue";
import { useRoute } from "vue-router";

type NavItem = { to: string; icon: string; label: string };

const props = withDefaults(
  defineProps<{ items: NavItem[]; maxItems?: number; moreLabel?: string }>(),
  { maxItems: 5, moreLabel: "More" },
);
const emit = defineEmits<{ more: [] }>();
const route = useRoute();

/**
 * M3 navigation bars hold 3–5 destinations. Beyond that the extras move behind
 * a "More" destination that opens the navigation sheet; otherwise the labels
 * truncate to "대시보…" at phone widths.
 */
const overflow = computed(() => props.items.length > props.maxItems);
const shown = computed(() =>
  overflow.value ? props.items.slice(0, props.maxItems - 1) : props.items,
);

function isActive(to: string) {
  return to === "/" ? route.path === "/" : route.path.startsWith(to);
}
</script>

<template>
  <nav class="or-navbar" :aria-label="$t('nav.mainNavigation')">
    <router-link
      v-for="item in shown"
      :key="item.to"
      :to="item.to"
      class="or-navbar__item"
      :class="{ 'is-active': isActive(item.to) }"
      :aria-current="isActive(item.to) ? 'page' : undefined"
    >
      <span class="or-navbar__icon-wrap">
        <OrIcon :name="item.icon" :size="24" :fill="isActive(item.to)" />
      </span>
      <span class="or-navbar__label">{{ item.label }}</span>
    </router-link>

    <button
      v-if="overflow"
      class="or-navbar__item"
      type="button"
      :aria-haspopup="true"
      @click="emit('more')"
    >
      <span class="or-navbar__icon-wrap">
        <OrIcon name="more_horiz" :size="24" />
      </span>
      <span class="or-navbar__label">{{ moreLabel }}</span>
    </button>
  </nav>
</template>

<style scoped>
.or-navbar {
  position: fixed;
  inset-inline: 0;
  inset-block-end: 0;
  z-index: 1000;
  display: flex;
  align-items: stretch;
  justify-content: space-around;
  block-size: calc(var(--or-nav-bar-height) + env(safe-area-inset-bottom, 0px));
  padding-block-end: env(safe-area-inset-bottom, 0px);
  padding-inline: var(--m3-space-75, 6px);
  background: var(--m3-color-surface-container);
  border-block-start: 1px solid var(--m3-color-outline-variant);
}

.or-navbar__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--m3-space-75, 6px);
  flex: 1;
  min-inline-size: 0;
  min-block-size: 48px;
  color: var(--m3-color-on-surface-variant);
  text-decoration: none;
  transition: color var(--or-motion-effects-fast);
}

.or-navbar__icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 56px;
  block-size: 32px;
  border-radius: var(--m3-shape-full);
  background: transparent;
  transition:
    background-color var(--or-motion-effects-fast),
    border-radius var(--or-motion-spatial-fast);
}

.or-navbar__item.is-active {
  color: var(--m3-color-on-surface);
}

.or-navbar__item.is-active .or-navbar__icon-wrap {
  background: var(--m3-color-secondary-container);
  color: var(--m3-color-on-secondary-container);
}

.or-navbar__label {
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  max-inline-size: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>

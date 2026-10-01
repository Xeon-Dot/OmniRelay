<script setup lang="ts">
/**
 * Desktop navigation: a 256px drawer when `expanded`, an 80px icon rail
 * otherwise (M3 navigation rail ↔ navigation drawer, spec §5.1).
 *
 * Active state follows M3: `primary-container` behind the leading icon with
 * `on-primary-container` content, morphed in with the fast spatial spring.
 * In rail mode the container collapses to a circle and a dot marks the
 * destination below the icon, with the label available as a tooltip.
 */
import { computed } from "vue";
import { useRoute } from "vue-router";

type NavItem = { to: string; icon: string; label: string };

const props = withDefaults(
  defineProps<{
    items: NavItem[];
    expanded: boolean;
    /** Brand block shown in the header — rendered by the caller. */
    brandSub?: string;
  }>(),
  { brandSub: undefined },
);

const emit = defineEmits<{ toggle: [] }>();
const route = useRoute();

function isActive(to: string) {
  return to === "/" ? route.path === "/" : route.path.startsWith(to);
}

const activeItem = computed(() => props.items.find((item) => isActive(item.to))?.label ?? "");
</script>

<template>
  <aside
    class="or-rail"
    :class="{ 'or-rail--expanded': expanded }"
    :aria-label="$t('nav.mainNavigation')"
  >
    <div class="or-rail__header">
      <slot name="brand" />
      <button
        class="or-rail__toggle"
        type="button"
        :aria-label="expanded ? $t('nav.collapse') : $t('nav.expand')"
        :aria-expanded="expanded"
        @click="emit('toggle')"
      >
        <OrIcon :name="expanded ? 'left_panel_close' : 'left_panel_open'" :size="20" />
      </button>
    </div>

    <nav class="or-rail__list">
      <router-link
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        class="or-rail__item"
        :class="{ 'is-active': isActive(item.to) }"
        :title="expanded ? undefined : item.label"
        :aria-current="isActive(item.to) ? 'page' : undefined"
      >
        <span class="or-rail__icon-wrap">
          <OrIcon class="or-rail__icon" :name="item.icon" :size="22" :fill="isActive(item.to)" />
        </span>
        <span v-if="expanded" class="or-rail__label">{{ item.label }}</span>
        <span v-else class="or-rail__dot" aria-hidden="true" />
      </router-link>
    </nav>

    <div class="or-rail__footer">
      <slot name="footer" />
    </div>

    <!-- Rail mode hides labels; expose the current destination to readers. -->
    <span class="or-visually-hidden">{{ activeItem }}</span>
  </aside>
</template>

<style scoped>
.or-rail {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  inline-size: var(--or-rail-width);
  block-size: 100dvh;
  position: sticky;
  inset-block-start: 0;
  background: var(--m3-color-surface-container-low);
  border-inline-end: 1px solid var(--m3-color-outline-variant);
  transition: inline-size var(--or-motion-spatial-default);
  overflow: hidden;
}

.or-rail--expanded {
  inline-size: var(--or-drawer-width);
}

.or-rail__header {
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
  padding: var(--m3-space-300) var(--m3-space-150);
  min-block-size: 72px;
}

.or-rail--expanded .or-rail__header {
  padding-inline: var(--m3-space-300);
}

.or-rail__toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 40px;
  block-size: 40px;
  margin-inline-start: auto;
  flex-shrink: 0;
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
  cursor: pointer;
  transition: background-color var(--or-motion-effects-fast);
}

.or-rail__toggle:hover {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
}

.or-rail__list {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-75, 6px);
  padding: var(--m3-space-100) var(--m3-space-150);
  overflow-y: auto;
  flex: 1;
}

.or-rail--expanded .or-rail__list {
  padding-inline: var(--m3-space-200);
}

.or-rail__item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  min-block-size: 56px;
  padding-block: var(--m3-space-75, 6px);
  padding-inline: var(--m3-space-100);
  border-radius: var(--or-morph-nav-active);
  color: var(--m3-color-on-surface-variant);
  text-decoration: none;
  white-space: nowrap;
  transition:
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast),
    border-radius var(--or-motion-spatial-fast);
}

.or-rail--expanded .or-rail__item {
  flex-direction: row;
  justify-content: flex-start;
  gap: var(--m3-space-200);
  min-block-size: 48px;
  padding-block: 0;
  padding-inline-start: var(--m3-space-100);
}

.or-rail__item:hover {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
  color: var(--m3-color-on-surface);
  text-decoration: none;
}

.or-rail__item.is-active {
  background: var(--m3-color-primary-container);
  color: var(--m3-color-on-primary-container);
  font-weight: 500;
}

.or-rail__icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 40px;
  block-size: 40px;
  flex-shrink: 0;
  border-radius: var(--m3-shape-full);
  transition: background-color var(--or-motion-effects-fast);
}

.or-rail__item.is-active .or-rail__icon-wrap {
  background: color-mix(in srgb, var(--m3-color-on-primary-container) 14%, transparent);
}

.or-rail__label {
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
}

/* Rail mode: icon only, plus a dot under it for the active destination. */
.or-rail__dot {
  position: absolute;
  inset-block-end: 4px;
  inline-size: 4px;
  block-size: 4px;
  border-radius: var(--m3-shape-full);
  background: currentColor;
  opacity: 0;
}

.or-rail__item.is-active .or-rail__dot {
  opacity: 1;
}

.or-rail--expanded .or-rail__dot {
  display: none;
}

.or-rail__footer {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-100);
  padding: var(--m3-space-200) var(--m3-space-150) var(--m3-space-300);
  border-block-start: 1px solid var(--m3-color-outline-variant);
}

.or-rail--expanded .or-rail__footer {
  padding-inline: var(--m3-space-200);
}
</style>

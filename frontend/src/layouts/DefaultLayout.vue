<template>
  <!-- Desktop / tablet: navigation rail, expanding to a drawer above 1200px -->
  <OrNavRail
    v-if="!isCompact"
    class="layout-rail"
    :items="visibleItems"
    :expanded="drawerExpanded"
    :brand-sub="auth.user?.username || $t('nav.gateway')"
    @toggle="drawerExpanded = !drawerExpanded"
  >
    <template #brand>
      <img class="brand-logo" :src="logoUrl" alt="OmniRelay" />
      <div v-if="drawerExpanded" class="brand-text">
        <div class="brand-name">OmniRelay</div>
        <div class="brand-sub">{{ auth.user?.username || $t("nav.gateway") }}</div>
      </div>
    </template>

    <template #footer>
      <!-- Theme: a three-way choice reads better as labels when there is room -->
      <OrSegmentedButton
        v-if="drawerExpanded"
        size="sm"
        block
        :options="themeOptions"
        :model-value="theme.mode"
        @update:model-value="onThemeMode"
      />
      <OrIconButton
        v-else
        :icon="theme.resolved === 'dark' ? 'dark_mode' : 'light_mode'"
        :label="$t('common.theme')"
        size="sm"
        @click="theme.toggle()"
      />

      <!-- Language -->
      <OrSegmentedButton
        v-if="drawerExpanded"
        size="sm"
        block
        :options="localeOptions"
        :model-value="locale"
        @update:model-value="onLocaleMode"
      />
      <div v-else class="locale-mini">
        <button
          v-for="option in localeOptions"
          :key="option.value"
          class="locale-mini__btn"
          :class="{ 'is-active': locale === option.value }"
          type="button"
          @click="switchLocale(option.value)"
        >
          {{ option.label }}
        </button>
      </div>

      <button class="logout-btn" type="button" @click="handleLogout">
        <OrIcon name="logout" :size="20" />
        <span v-if="drawerExpanded" class="logout-btn__label">{{ $t("common.signOut") }}</span>
      </button>
    </template>
  </OrNavRail>

  <!-- Mobile: bottom navigation bar; destinations past the fifth live behind
       "More", which opens the sheet below (M3 caps a nav bar at 5). -->
  <OrNavigationBar
    v-else
    :items="visibleItems"
    :more-label="$t('nav.more')"
    @more="navSheet = true"
  />

  <main class="layout-main" :class="{ 'layout-main--compact': isCompact }">
    <div class="layout-container">
      <router-view />
    </div>
  </main>

  <!-- Mobile overflow menu: every destination in a bottom sheet -->
  <OrDialog
    v-model="navSheet"
    sheet
    :title="$t('nav.mainNavigation')"
  >
    <nav class="nav-sheet">
      <router-link
        v-for="item in visibleItems"
        :key="item.to"
        :to="item.to"
        class="nav-sheet__item"
        :class="{ 'is-active': isActive(item.to) }"
        :aria-current="isActive(item.to) ? 'page' : undefined"
        @click="navSheet = false"
      >
        <OrIcon :name="item.icon" :size="22" />
        <span>{{ item.label }}</span>
      </router-link>
    </nav>
  </OrDialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth";
import { useThemeStore, type ThemeMode } from "../stores/theme";
import { setLocale } from "../plugins/i18n";
import logoUrl from "../assets/omnirelay-logo.svg";

const auth = useAuthStore();
const theme = useThemeStore();
const router = useRouter();
const route = useRoute();
const { locale, t } = useI18n();

/** Open state of the mobile "More" navigation sheet. */
const navSheet = ref(false);

function isActive(to: string) {
  return to === "/" ? route.path === "/" : route.path.startsWith(to);
}

const menuItems = [
  { i18nKey: "nav.dashboard", icon: "space_dashboard", to: "/" },
  { i18nKey: "nav.providers", icon: "dns", to: "/providers" },
  { i18nKey: "nav.models", icon: "deployed_code", to: "/models" },
  { i18nKey: "nav.apiKeys", icon: "key", to: "/api-keys" },
  { i18nKey: "nav.usage", icon: "monitoring", to: "/usage" },
  { i18nKey: "nav.logs", icon: "receipt_long", to: "/logs" },
  { i18nKey: "nav.performance", icon: "speed", to: "/performance" },
  { i18nKey: "nav.users", icon: "group", to: "/users", adminOnly: true },
] as const;

const visibleItems = computed(() =>
  menuItems
    .filter((item) => !("adminOnly" in item) || auth.user?.is_admin)
    .map((item) => ({ to: item.to, icon: item.icon, label: tOf(item.i18nKey) })),
);

/** i18n at call time, not at module load — the locale can change. */
function tOf(key: string) {
  return t(key);
}

// Labels only: the drawer is 256px, and an icon plus a 3-character CJK label
// does not fit in a third of it (see the `block` segmented below).
const themeOptions = computed(() => [
  { value: "system", label: t("common.themeSystem") },
  { value: "light", label: t("common.themeLight") },
  { value: "dark", label: t("common.themeDark") },
]);

const localeOptions = [
  { value: "en", label: "EN" },
  { value: "ja", label: "JA" },
  { value: "ko", label: "KO" },
];

function switchLocale(next: string) {
  setLocale(next);
}

/** Segmented buttons emit `string | string[]`; both uses here are single-select. */
function onThemeMode(value: string | string[]) {
  if (typeof value === "string") theme.setMode(value as ThemeMode);
}

function onLocaleMode(value: string | string[]) {
  if (typeof value === "string") switchLocale(value);
}

const drawerExpanded = ref(true);

/**
 * Three navigation shells (spec §5.1):
 *   ≥1200px drawer, 600–1199px icon rail, <600px bottom bar.
 */
const DRAWER_BREAKPOINT = 1200;
const viewportWidth = ref(typeof window === "undefined" ? 1440 : window.innerWidth);
function onResize() {
  const next = window.innerWidth;
  const wasDrawer = viewportWidth.value >= DRAWER_BREAKPOINT;
  viewportWidth.value = next;
  // Only crossing the boundary re-evaluates the default; a manual toggle
  // inside a band survives incidental resizes.
  if (wasDrawer !== (next >= DRAWER_BREAKPOINT)) {
    drawerExpanded.value = next >= DRAWER_BREAKPOINT;
  }
}
const isCompact = computed(() => viewportWidth.value < 600);

onMounted(() => {
  window.addEventListener("resize", onResize, { passive: true });
  drawerExpanded.value = viewportWidth.value >= DRAWER_BREAKPOINT;
});
onBeforeUnmount(() => window.removeEventListener("resize", onResize));

function handleLogout() {
  auth.logout();
  router.push("/login");
}
</script>

<style scoped>
.brand-logo {
  inline-size: 32px;
  block-size: 32px;
  flex-shrink: 0;
}

.brand-text {
  min-inline-size: 0;
  overflow: hidden;
}

.brand-name {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-emphasized-title-medium);
  letter-spacing: var(--m3-typescale-emphasized-title-medium-tracking);
  white-space: nowrap;
}

.brand-sub {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  text-transform: uppercase;
  white-space: nowrap;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: var(--m3-space-150);
  min-block-size: 48px;
  padding-inline: var(--m3-space-100);
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
  cursor: pointer;
  transition:
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

.logout-btn:hover {
  background: var(--m3-color-error-container);
  color: var(--m3-color-on-error-container);
}

.logout-btn__label {
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
}

.locale-mini {
  display: flex;
  justify-content: center;
  gap: 2px;
}

.locale-mini__btn {
  flex: 1;
  min-block-size: 32px;
  border: 1px solid var(--m3-color-outline-variant);
  border-radius: var(--m3-shape-xs);
  background: transparent;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  cursor: pointer;
  transition:
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

.locale-mini__btn.is-active {
  background: var(--m3-color-primary);
  border-color: var(--m3-color-primary);
  color: var(--m3-color-on-primary);
}

/* ── Content ── */
.layout-main {
  flex: 1;
  min-inline-size: 0;
  min-block-size: 100dvh;
}

.layout-container {
  max-inline-size: var(--or-layout-max-width);
  margin-inline: auto;
  padding: var(--or-layout-padding-block-start) var(--or-layout-padding-inline);
}

.layout-main--compact .layout-container {
  padding-inline: var(--or-layout-padding-inline-compact);
  padding-block-end: calc(
    var(--or-nav-bar-height) + env(safe-area-inset-bottom, 0px) + var(--m3-space-400)
  );
}

/* ── Mobile "More" navigation sheet ── */
.nav-sheet {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-75, 6px);
}

.nav-sheet__item {
  display: flex;
  align-items: center;
  gap: var(--m3-space-200);
  min-block-size: 56px;
  padding-inline: var(--m3-space-200);
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-on-surface-variant);
  text-decoration: none;
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
  transition:
    background-color var(--or-motion-effects-fast),
    color var(--or-motion-effects-fast);
}

.nav-sheet__item:hover {
  background: color-mix(in srgb, var(--m3-color-on-surface) 8%, transparent);
  color: var(--m3-color-on-surface);
}

.nav-sheet__item.is-active {
  background: var(--m3-color-secondary-container);
  color: var(--m3-color-on-secondary-container);
}

</style>

<script setup lang="ts">
/**
 * Snackbar host. Mount exactly one instance (App.vue) — it reads the store,
 * so any component can call `useSnackbarStore().show(...)`.
 *
 * `role="status"` + `aria-live="polite"` announces the message without
 * interrupting whatever the user is doing. Sits above the mobile navigation
 * bar so it never covers the tabs.
 */
import { computed } from "vue";
import { useSnackbarStore } from "../../stores/snackbar";

const store = useSnackbarStore();

const toneClass = computed(() => `or-snackbar--${store.variant}`);

function onAction() {
  store.action?.onClick();
  store.hide();
}
</script>

<template>
  <Teleport to="body">
    <Transition name="or-snackbar">
      <div
        v-if="store.visible"
        class="or-snackbar"
        :class="toneClass"
        role="status"
        aria-live="polite"
      >
        <OrIcon
          :name="
            store.variant === 'success'
              ? 'check_circle'
              : store.variant === 'error'
                ? 'error'
                : 'info'
          "
          :size="20"
        />
        <span class="or-snackbar__message">{{ store.message }}</span>
        <button v-if="store.action" class="or-snackbar__action" type="button" @click="onAction">
          {{ store.action.label }}
        </button>
        <button class="or-snackbar__close" type="button" :aria-label="$t('common.dismiss')" @click="store.hide()">
          <OrIcon name="close" :size="18" />
        </button>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.or-snackbar {
  position: fixed;
  inset-inline: 0;
  inset-block-end: calc(var(--or-nav-bar-height) + 32px + env(safe-area-inset-bottom, 0px));
  z-index: 1100;
  display: flex;
  align-items: center;
  gap: var(--m3-space-150);
  inline-size: fit-content;
  max-inline-size: min(560px, calc(100vw - 32px));
  margin-inline: auto;
  padding: var(--m3-space-150) var(--m3-space-100) var(--m3-space-150) var(--m3-space-200);
  border-radius: var(--m3-shape-xs);
  background: var(--m3-color-inverse-surface);
  color: var(--m3-color-inverse-on-surface);
  box-shadow: var(--m3-shadow-3);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}

.or-snackbar--success .or-snackbar__message {
  color: var(--m3-color-inverse-on-surface);
}

.or-snackbar__message {
  flex: 1;
  min-inline-size: 0;
  overflow-wrap: anywhere;
}

.or-snackbar__action {
  flex-shrink: 0;
  min-block-size: 32px;
  padding-inline: var(--m3-space-150);
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-inverse-primary);
  font: var(--m3-typescale-label-large);
  letter-spacing: var(--m3-typescale-label-large-tracking);
  cursor: pointer;
}

.or-snackbar__action:hover {
  background: color-mix(in srgb, var(--m3-color-inverse-primary) 16%, transparent);
}

.or-snackbar__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  inline-size: 32px;
  block-size: 32px;
  flex-shrink: 0;
  border-radius: var(--m3-shape-full);
  color: var(--m3-color-inverse-on-surface);
  cursor: pointer;
}

.or-snackbar__close:hover {
  background: color-mix(in srgb, var(--m3-color-inverse-on-surface) 12%, transparent);
}

.or-snackbar-enter-active,
.or-snackbar-leave-active {
  transition:
    transform var(--or-motion-spatial-fast),
    opacity var(--or-motion-effects-fast);
}

.or-snackbar-enter-from,
.or-snackbar-leave-to {
  opacity: 0;
  transform: translateY(24px) scale(0.96);
}
</style>

import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";

export type ThemeMode = "system" | "light" | "dark";
export type ResolvedTheme = "light" | "dark";

const STORAGE_KEY = "omnirelay.theme";

function prefersDark(): boolean {
  return (
    typeof window !== "undefined" &&
    typeof window.matchMedia === "function" &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  );
}

function readStoredMode(): ThemeMode {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw === "light" || raw === "dark" || raw === "system") return raw;
  } catch {
    /* private mode / disabled storage — fall through to system */
  }
  return "system";
}

/**
 * Theme resolution happens in two places: the inline script in index.html
 * (before first paint, so there is no flash) and here (so the toggle stays in
 * sync afterwards). Both read the same key with the same fallback order.
 */
export const useThemeStore = defineStore("theme", () => {
  const mode = ref<ThemeMode>(readStoredMode());
  const systemResolved = ref<ResolvedTheme>(prefersDark() ? "dark" : "light");

  if (typeof window !== "undefined" && typeof window.matchMedia === "function") {
    const query = window.matchMedia("(prefers-color-scheme: dark)");
    const onChange = () => {
      systemResolved.value = query.matches ? "dark" : "light";
    };
    // Modern form first; older Safari only ships addListener.
    if (typeof query.addEventListener === "function") {
      query.addEventListener("change", onChange);
    } else {
      query.addListener(onChange);
    }
  }

  const resolved = computed<ResolvedTheme>(() =>
    mode.value === "system" ? systemResolved.value : mode.value,
  );

  watch(
    resolved,
    (value) => {
      document.documentElement.dataset.theme = value;
    },
    { immediate: true },
  );

  function setMode(next: ThemeMode) {
    mode.value = next;
    try {
      localStorage.setItem(STORAGE_KEY, next);
    } catch {
      /* storage unavailable — the choice simply won't persist */
    }
    // Keep the DOM in sync even when resolved() doesn't change (e.g. switching
    // from "system" to an explicit mode that matches the current system theme).
    document.documentElement.dataset.theme = resolved.value;
  }

  /** Cycles light ↔ dark. Leaves "system" for an explicit choice. */
  function toggle() {
    setMode(resolved.value === "dark" ? "light" : "dark");
  }

  return { mode, resolved, setMode, toggle };
});

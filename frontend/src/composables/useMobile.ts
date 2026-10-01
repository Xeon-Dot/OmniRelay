import { ref, onMounted, onUnmounted } from "vue";

// Must match DefaultLayout's compact shell boundary (spec §5.1: <600px = bottom bar).
const MOBILE_BREAKPOINT = 600;

export function useMobile() {
  const isMobile = ref(false);
  let mql: MediaQueryList | null = null;

  function handleChange(e: MediaQueryListEvent) {
    isMobile.value = e.matches;
  }

  onMounted(() => {
    mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT}px)`);
    isMobile.value = mql.matches;
    mql.addEventListener("change", handleChange);
  });

  onUnmounted(() => {
    mql?.removeEventListener("change", handleChange);
  });

  return { isMobile };
}
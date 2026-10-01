import { computed } from "vue";
import { useThemeStore } from "../stores/theme";

/**
 * Theme-reactive Chart.js colours.
 *
 * CSS custom properties are not observable, so the palette is keyed off the
 * theme store: when `resolved` flips, the computed re-reads `getComputedStyle`
 * and every chart bound to it redraws. Pass a `computed()` built from this
 * straight to `<Line :options="...">` — vue-chartjs watches the prop and
 * calls `update()` itself.
 *
 *   const { palette } = useChartTheme();
 *   const options = computed(() => ({
 *     scales: { y: { ticks: { color: palette.value.text } } },
 *   }));
 */
export interface ChartPalette {
  /** Primary series. */
  primary: string;
  /** Secondary series. */
  tertiary: string;
  /** Tertiary series. */
  secondary: string;
  /** Grid lines and axis borders. */
  grid: string;
  /** Tick labels and legend text. */
  text: string;
  /** Tooltip / legend background. */
  surface: string;
  /** Text on that background. */
  onSurface: string;
  danger: string;
  success: string;
}

function readToken(name: string): string {
  if (typeof window === "undefined") return "";
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}

export function useChartTheme() {
  const theme = useThemeStore();

  const palette = computed<ChartPalette>(() => {
    // Reactive dependency: re-read the CSS variables whenever the theme flips.
    void theme.resolved;
    return {
      primary: readToken("--m3-color-primary"),
      tertiary: readToken("--m3-color-tertiary"),
      secondary: readToken("--m3-color-secondary"),
      grid: readToken("--m3-color-outline-variant"),
      text: readToken("--m3-color-on-surface-variant"),
      surface: readToken("--m3-color-surface-container-highest"),
      onSurface: readToken("--m3-color-on-surface"),
      danger: readToken("--m3-color-error"),
      success: readToken("--or-success"),
    };
  });

  return { palette, theme };
}

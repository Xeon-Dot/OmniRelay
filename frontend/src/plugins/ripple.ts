import type { Directive } from "vue";

/**
 * `v-ripple` — M3 touch feedback.
 *
 * Spawns a growing circle at the pointer position on pointerdown and removes
 * it when its animation finishes. Skipped for keyboard activation (no
 * pointerdown) and when the user has asked for reduced motion, where the CSS
 * gate would collapse the animation anyway.
 */
function prefersReducedMotion(): boolean {
  return (
    typeof window !== "undefined" &&
    typeof window.matchMedia === "function" &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches
  );
}

export const ripple: Directive<HTMLElement, boolean | undefined> = {
  mounted(el, binding) {
    if (binding.value === false) return;

    el.addEventListener("pointerdown", (event: PointerEvent) => {
      if (prefersReducedMotion()) return;
      // Right/middle click and keyboard-triggered clicks carry no useful origin.
      if (event.button !== 0) return;

      const rect = el.getBoundingClientRect();
      const diameter = Math.max(rect.width, rect.height) * 2;
      const span = document.createElement("span");
      span.className = "or-ripple";
      span.style.width = span.style.height = `${diameter}px`;
      span.style.left = `${event.clientX - rect.left - diameter / 2}px`;
      span.style.top = `${event.clientY - rect.top - diameter / 2}px`;
      span.addEventListener("animationend", () => span.remove(), { once: true });

      el.appendChild(span);
      // Guard against an animation that never fires (element hidden / detached).
      window.setTimeout(() => span.remove(), 1000);
    });
  },
};

export default ripple;

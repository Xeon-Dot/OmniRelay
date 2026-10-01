import { onBeforeUnmount, watch, type Ref } from "vue";

const FOCUSABLE = [
  "a[href]",
  "button:not([disabled])",
  "textarea:not([disabled])",
  "input:not([disabled])",
  "select:not([disabled])",
  "[tabindex]:not([tabindex='-1'])",
].join(", ");

function focusableWithin(root: HTMLElement): HTMLElement[] {
  return [...root.querySelectorAll<HTMLElement>(FOCUSABLE)].filter(
    (el) => el.getAttribute("aria-hidden") !== "true" && !el.hasAttribute("inert"),
  );
}

/**
 * Traps keyboard focus inside `container` while `active` is true.
 *
 * Restores focus to whatever was focused before the trap opened — closing a
 * dialog must return the user to the control that opened it, or keyboard
 * navigation restarts from the top of the document.
 */
export function useFocusTrap(container: Ref<HTMLElement | null>, active: Ref<boolean>) {
  let previouslyFocused: HTMLElement | null = null;
  let restoreFrame = 0;

  function onKeydown(event: KeyboardEvent) {
    if (event.key !== "Tab") return;
    const root = container.value;
    if (!root) return;

    const items = focusableWithin(root);
    if (items.length === 0) {
      event.preventDefault();
      root.focus();
      return;
    }
    const first = items[0];
    const last = items[items.length - 1];
    const current = document.activeElement as HTMLElement | null;

    if (event.shiftKey) {
      if (current === first || !root.contains(current)) {
        event.preventDefault();
        last.focus();
      }
    } else if (current === last || !root.contains(current)) {
      event.preventDefault();
      first.focus();
    }
  }

  function onFocusIn(event: FocusEvent) {
    const root = container.value;
    if (!root) return;
    if (root.contains(event.target as Node)) return;
    const items = focusableWithin(root);
    (items[0] ?? root).focus();
  }

  function attach() {
    previouslyFocused = document.activeElement as HTMLElement | null;
    document.addEventListener("keydown", onKeydown, true);
    document.addEventListener("focusin", onFocusIn, true);
    restoreFrame = requestAnimationFrame(() => {
      const root = container.value;
      if (!root) return;
      const items = focusableWithin(root);
      (items[0] ?? root).focus();
    });
  }

  function detach() {
    document.removeEventListener("keydown", onKeydown, true);
    document.removeEventListener("focusin", onFocusIn, true);
    cancelAnimationFrame(restoreFrame);
    const target = previouslyFocused;
    previouslyFocused = null;
    if (target && document.contains(target)) target.focus();
  }

  watch(
    active,
    (on) => {
      if (on) attach();
      else detach();
    },
    { flush: "post" },
  );

  onBeforeUnmount(detach);
}

let scrollLocks = 0;

/**
 * Body scroll lock with a nesting counter: two dialogs open at once (the Users
 * view does this) must both be closed before scrolling comes back.
 */
export function lockBodyScroll() {
  if (scrollLocks++ === 0) {
    document.body.style.overflow = "hidden";
  }
}

export function unlockBodyScroll() {
  scrollLocks = Math.max(0, scrollLocks - 1);
  if (scrollLocks === 0) {
    document.body.style.overflow = "";
  }
}

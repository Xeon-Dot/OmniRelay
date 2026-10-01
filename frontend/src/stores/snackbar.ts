import { defineStore } from "pinia";
import { ref } from "vue";

export type SnackbarVariant = "info" | "success" | "error";

export interface SnackbarOptions {
  variant?: SnackbarVariant;
  /** Auto-dismiss in ms. 0 keeps it until dismissed. Defaults to 4000. */
  timeout?: number;
  action?: { label: string; onClick: () => void };
}

/**
 * Transient feedback. Replaces the `window.alert(...)` calls that the Users
 * and ApiKeys views used for "copied", "saved" and "reset password" notices —
 * those are confirmations of success, not blocking errors, and an alert box
 * stops the user dead for one that.
 *
 * Destructive *confirmations* deliberately keep `window.confirm`: a native
 * confirm is accessible, cannot be dismissed by accident, and costs nothing
 * to keep correct.
 */
export const useSnackbarStore = defineStore("snackbar", () => {
  const visible = ref(false);
  const message = ref("");
  const variant = ref<SnackbarVariant>("info");
  const action = ref<SnackbarOptions["action"]>(undefined);
  let timer = 0;

  function hide() {
    visible.value = false;
    window.clearTimeout(timer);
    timer = 0;
  }

  function show(text: string, options: SnackbarOptions = {}) {
    message.value = text;
    variant.value = options.variant ?? "info";
    action.value = options.action;
    visible.value = true;
    window.clearTimeout(timer);
    const timeout = options.timeout ?? 4000;
    timer = timeout > 0 ? window.setTimeout(hide, timeout) : 0;
  }

  function success(text: string, options: Omit<SnackbarOptions, "variant"> = {}) {
    show(text, { ...options, variant: "success" });
  }

  function error(text: string, options: Omit<SnackbarOptions, "variant"> = {}) {
    show(text, { ...options, variant: "error" });
  }

  return { visible, message, variant, action, show, hide, success, error };
});

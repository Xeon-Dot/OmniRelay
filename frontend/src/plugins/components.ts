import type { App } from "vue";

import OrButton from "../components/ui/OrButton.vue";
import OrCard from "../components/ui/OrCard.vue";
import OrCheckbox from "../components/ui/OrCheckbox.vue";
import OrChip from "../components/ui/OrChip.vue";
import OrDivider from "../components/ui/OrDivider.vue";
import OrDataTable from "../components/ui/OrDataTable.vue";
import OrDialog from "../components/ui/OrDialog.vue";
import OrIconButton from "../components/ui/OrIconButton.vue";
import OrIcon from "../components/ui/OrIcon.vue";
import OrNavigationBar from "../components/ui/OrNavigationBar.vue";
import OrNavRail from "../components/ui/OrNavRail.vue";
import OrProgressBar from "../components/ui/OrProgressBar.vue";
import OrSegmentedButton from "../components/ui/OrSegmentedButton.vue";
import OrSelect from "../components/ui/OrSelect.vue";
import OrSnackbar from "../components/ui/OrSnackbar.vue";
import OrSpinner from "../components/ui/OrSpinner.vue";
import OrSwitch from "../components/ui/OrSwitch.vue";
import OrTextField from "../components/ui/OrTextField.vue";

/**
 * Global registration for the M3 Expressive primitives.
 * Every task that adds a component appends ONE import line and ONE register
 * line — never a wholesale rewrite, so parallel edits do not clobber each other.
 */
export function install(app: App) {
  app.component("OrIcon", OrIcon);
  app.component("OrSpinner", OrSpinner);
  app.component("OrProgressBar", OrProgressBar);
  app.component("OrButton", OrButton);
  app.component("OrIconButton", OrIconButton);
  app.component("OrCard", OrCard);
  app.component("OrChip", OrChip);
  app.component("OrSegmentedButton", OrSegmentedButton);
  app.component("OrTextField", OrTextField);
  app.component("OrSelect", OrSelect);
  app.component("OrCheckbox", OrCheckbox);
  app.component("OrSwitch", OrSwitch);
  app.component("OrDivider", OrDivider);
  app.component("OrDialog", OrDialog);
  app.component("OrDataTable", OrDataTable);
  app.component("OrSnackbar", OrSnackbar);
  app.component("OrNavRail", OrNavRail);
  app.component("OrNavigationBar", OrNavigationBar);
}

export default { install };

import type OrButton from "./components/ui/OrButton.vue";
import type OrCard from "./components/ui/OrCard.vue";
import type OrCheckbox from "./components/ui/OrCheckbox.vue";
import type OrChip from "./components/ui/OrChip.vue";
import type OrDataTable from "./components/ui/OrDataTable.vue";
import type OrDialog from "./components/ui/OrDialog.vue";
import type OrDivider from "./components/ui/OrDivider.vue";
import type OrIcon from "./components/ui/OrIcon.vue";
import type OrIconButton from "./components/ui/OrIconButton.vue";
import type OrNavRail from "./components/ui/OrNavRail.vue";
import type OrNavigationBar from "./components/ui/OrNavigationBar.vue";
import type OrProgressBar from "./components/ui/OrProgressBar.vue";
import type OrSegmentedButton from "./components/ui/OrSegmentedButton.vue";
import type OrSelect from "./components/ui/OrSelect.vue";
import type OrSnackbar from "./components/ui/OrSnackbar.vue";
import type OrSpinner from "./components/ui/OrSpinner.vue";
import type OrSwitch from "./components/ui/OrSwitch.vue";
import type OrTextField from "./components/ui/OrTextField.vue";

declare module "vue" {
  export interface GlobalComponents {
    OrButton: typeof OrButton;
    OrCard: typeof OrCard;
    OrCheckbox: typeof OrCheckbox;
    OrChip: typeof OrChip;
    OrDataTable: typeof OrDataTable;
    OrDialog: typeof OrDialog;
    OrDivider: typeof OrDivider;
    OrIcon: typeof OrIcon;
    OrIconButton: typeof OrIconButton;
    OrNavRail: typeof OrNavRail;
    OrNavigationBar: typeof OrNavigationBar;
    OrProgressBar: typeof OrProgressBar;
    OrSegmentedButton: typeof OrSegmentedButton;
    OrSelect: typeof OrSelect;
    OrSnackbar: typeof OrSnackbar;
    OrSpinner: typeof OrSpinner;
    OrSwitch: typeof OrSwitch;
    OrTextField: typeof OrTextField;
  }
}

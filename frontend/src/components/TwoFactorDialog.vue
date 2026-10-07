<template>
  <OrDialog v-model="open" :width="440" :title="$t('twofa.title')">
    <div class="dialog-stack">
      <p class="dim">{{ $t("twofa.desc") }}</p>

      <AppAlert v-if="error" variant="error">{{ error }}</AppAlert>

      <template v-if="enabled">
        <StatusChip variant="on">{{ $t("twofa.enabled") }}</StatusChip>
        <OrTextField v-model="disablePassword" type="password" :label="$t('twofa.password')" />
        <OrTextField v-model="disableCode" type="text" :label="$t('twofa.code')" />
        <OrButton :loading="busy" @click="disable">{{ $t("twofa.disable") }}</OrButton>
      </template>

      <template v-else-if="!setup">
        <StatusChip variant="off">{{ $t("twofa.disabled") }}</StatusChip>
        <OrButton :loading="busy" @click="startSetup">{{ $t("twofa.enable") }}</OrButton>
      </template>

      <template v-else-if="!recoveryCodes.length">
        <p class="dim">{{ $t("twofa.setupHint") }}</p>
        <img v-if="qrImage" class="qr" :src="qrImage" :alt="$t('twofa.title')" />
        <code class="secret">{{ setup.secret }}</code>
        <code class="secret wrap">{{ setup.otpauth_url }}</code>
        <OrTextField v-model="enableCode" type="text" :label="$t('twofa.code')" />
        <OrButton :loading="busy" @click="confirmEnable">{{ $t("twofa.confirm") }}</OrButton>
      </template>

      <template v-else>
        <AppAlert variant="success">{{ $t("twofa.enabledOk") }}</AppAlert>
        <p class="dim">{{ $t("twofa.recoveryHint") }}</p>
        <div class="recovery">
          <code v-for="c in recoveryCodes" :key="c">{{ c }}</code>
        </div>
        <OrButton @click="open = false">{{ $t("twofa.done") }}</OrButton>
      </template>
    </div>
  </OrDialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import api from "../api/client";
import { useAuthStore } from "../stores/auth";
import AppAlert from "./AppAlert.vue";
import StatusChip from "./StatusChip.vue";
import QRCode from "qrcode";

const open = defineModel<boolean>({ default: false });
const auth = useAuthStore();

const enabled = ref(false);
const setup = ref<{ secret: string; otpauth_url: string } | null>(null);
const qrImage = ref("");
const recoveryCodes = ref<string[]>([]);
const enableCode = ref("");
const disablePassword = ref("");
const disableCode = ref("");
const error = ref("");
const busy = ref(false);

watch(open, async (v) => {
  if (!v) return;
  error.value = "";
  setup.value = null;
  qrImage.value = "";
  recoveryCodes.value = [];
  try {
    const { data } = await api.get("/auth/2fa/status");
    enabled.value = !!data.enabled;
  } catch (e: any) {
    error.value = e.response?.data?.error || "Failed to load status";
  }
});

async function startSetup() {
  busy.value = true;
  error.value = "";
  try {
    const { data } = await api.post("/auth/2fa/setup");
    setup.value = data;
    qrImage.value = await QRCode.toDataURL(data.otpauth_url, { margin: 1, width: 220 });
  } catch (e: any) {
    error.value = e.response?.data?.error || "Setup failed";
  } finally {
    busy.value = false;
  }
}

async function confirmEnable() {
  busy.value = true;
  error.value = "";
  try {
    const { data } = await api.post("/auth/2fa/enable", { code: enableCode.value.trim() });
    recoveryCodes.value = data.recovery_codes || [];
    enabled.value = true;
    setup.value = null;
    if (auth.user) auth.user.totp_enabled = true;
  } catch (e: any) {
    error.value = e.response?.data?.error || "Enable failed";
  } finally {
    busy.value = false;
  }
}

async function disable() {
  busy.value = true;
  error.value = "";
  try {
    await api.post("/auth/2fa/disable", { password: disablePassword.value, code: disableCode.value.trim() });
    enabled.value = false;
    disablePassword.value = "";
    disableCode.value = "";
    if (auth.user) auth.user.totp_enabled = false;
  } catch (e: any) {
    error.value = e.response?.data?.error || "Disable failed";
  } finally {
    busy.value = false;
  }
}
</script>

<style scoped>
.dialog-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.dim {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-medium);
  margin: 0;
}
.secret {
  background: var(--m3-color-surface-container-high);
  padding: 8px 10px;
  border-radius: 8px;
  font-family: monospace;
  word-break: break-all;
}
.qr {
  align-self: center;
  border-radius: 8px;
  background: #fff;
}
.recovery {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4px 12px;
  font-family: monospace;
}
</style>

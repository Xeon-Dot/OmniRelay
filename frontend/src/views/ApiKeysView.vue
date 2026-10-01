<template>
  <div class="page">
    <PageHeader :title="$t('apiKeys.title')" :subtitle="$t('apiKeys.subtitle')">
      <OrButton variant="filled" @click="openCreateDialog">
        <template #leading>
          <OrIcon name="add" :size="18" />
        </template>
        {{ $t("apiKeys.issueKey") }}
      </OrButton>
    </PageHeader>

    <OrDataTable
      :headers="headers"
      :items="store.apiKeys"
      :loading="store.loading"
      density="comfortable"
      :label="$t('apiKeys.title')"
    >
      <template #item.key_prefix="{ item }">
        <MonoTag>{{ item.key_prefix }}</MonoTag>
      </template>
      <template #item.is_active="{ item }">
        <StatusChip :variant="item.is_active ? 'on' : 'off'">
          {{ item.is_active ? $t("apiKeys.active") : $t("apiKeys.revoked") }}
        </StatusChip>
      </template>
      <template #item.last_used_at="{ item }">
        <span class="dim-text">
          {{
            item.last_used_at
              ? new Date(item.last_used_at).toLocaleString()
              : $t("apiKeys.never")
          }}
        </span>
      </template>
      <template #item.rate_limit_rpm="{ item }">
        <span class="mono-val">{{
          item.rate_limit_rpm === 0 ? "∞" : item.rate_limit_rpm
        }}</span>
      </template>
      <template #item.created_at="{ item }">
        <span class="dim-text">{{
          new Date(item.created_at).toLocaleDateString()
        }}</span>
      </template>
      <template #item.actions="{ item }">
        <div class="row-actions">
          <OrIconButton
            v-if="item.is_active"
            icon="block"
            size="sm"
            tone="danger"
            :label="$t('apiKeys.revoke')"
            @click="handleDelete(item.id)"
          />
        </div>
      </template>
      <template #no-data>
        <EmptyState icon="key_off" :text="$t('apiKeys.noKeys')" />
      </template>
    </OrDataTable>

    <!-- Create dialog -->
    <OrDialog
      v-model="createDialog"
      :width="isMobile ? undefined : 460"
      :fullscreen="isMobile"
      :title="$t('apiKeys.issueNew')"
    >
      <div class="dialog-stack">
        <OrTextField
          v-model="form.name"
          :label="$t('apiKeys.keyName')"
          placeholder="My App Key"
        />
        <OrTextField
          v-model.number="form.rate_limit_rpm"
          type="number"
          :label="$t('apiKeys.rateLimitRpm')"
          :hint="$t('apiKeys.rateLimitHint')"
          placeholder="0 = unlimited"
        />
        <AppAlert v-if="dialogError" variant="error">{{
          dialogError
        }}</AppAlert>
      </div>
      <template #footer>
        <OrButton variant="text" @click="createDialog = false">
          {{ $t("common.cancel") }}
        </OrButton>
        <OrButton :loading="creating" @click="handleCreate">
          {{ $t("common.create") }}
        </OrButton>
      </template>
    </OrDialog>

    <!-- Show key dialog -->
    <OrDialog
      v-model="showKey"
      :width="isMobile ? undefined : 500"
      :fullscreen="isMobile"
      :title="$t('apiKeys.keyCreated')"
    >
      <div class="dialog-stack">
        <p class="reveal-note">{{ $t("apiKeys.saveKeyNote") }}</p>
        <div class="key-reveal">
          <code class="key-value">{{ newKey }}</code>
          <OrIconButton
            :icon="copied ? 'check' : 'content_copy'"
            :variant="copied ? 'tonal' : 'standard'"
            :label="$t('common.copy')"
            @click="copyKey"
          />
        </div>
      </div>
      <template #footer>
        <OrButton
          @click="
            showKey = false;
            createDialog = false;
          "
        >
          {{ $t("common.done") }}
        </OrButton>
      </template>
    </OrDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useApiKeysStore } from "../stores/apikeys";
import PageHeader from "../components/PageHeader.vue";
import EmptyState from "../components/EmptyState.vue";
import StatusChip from "../components/StatusChip.vue";
import MonoTag from "../components/MonoTag.vue";
import AppAlert from "../components/AppAlert.vue";
import { useMobile } from "../composables/useMobile";

const { t } = useI18n();
const store = useApiKeysStore();
const createDialog = ref(false);
const showKey = ref(false);
const newKey = ref("");
const creating = ref(false);
const dialogError = ref("");
const copied = ref(false);

const { isMobile } = useMobile();

const form = ref({ name: "", rate_limit_rpm: 0 });

const headers = computed(() => [
  { title: t("apiKeys.name"), key: "name" },
  { title: t("apiKeys.keyPrefix"), key: "key_prefix" },
  { title: t("apiKeys.status"), key: "is_active" },
  { title: t("apiKeys.rateLimit"), key: "rate_limit_rpm" },
  { title: t("apiKeys.lastUsed"), key: "last_used_at" },
  { title: t("apiKeys.created"), key: "created_at" },
  { title: "", key: "actions", sortable: false, width: 60 },
]);

function openCreateDialog() {
  dialogError.value = "";
  form.value = { name: "", rate_limit_rpm: 0 };
  createDialog.value = true;
}

async function handleCreate() {
  creating.value = true;
  dialogError.value = "";
  try {
    const result = await store.create(
      form.value.name,
      form.value.rate_limit_rpm,
    );
    newKey.value = result.plain_key;
    // Close the form before revealing the key: two focus-trapping dialogs open
    // at once fight over focus, and each would render the same title id.
    createDialog.value = false;
    showKey.value = true;
    await store.fetch();
  } catch (e: any) {
    dialogError.value = e.response?.data?.error || t("apiKeys.createFailed");
  } finally {
    creating.value = false;
  }
}

function copyKey() {
  navigator.clipboard.writeText(newKey.value);
  copied.value = true;
  setTimeout(() => {
    copied.value = false;
  }, 2000);
}

async function handleDelete(id: number) {
  if (!confirm(t("apiKeys.revokeConfirm"))) return;
  await store.remove(id);
}

onMounted(() => {
  store.fetch();
});
</script>

<style scoped>
@import "../styles/page-shared.css";

.dialog-stack {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-150);
}

.reveal-note {
  margin: 0;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-small);
  letter-spacing: var(--m3-typescale-body-small-tracking);
}
.key-reveal {
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
  padding: var(--m3-space-150) var(--m3-space-200);
  background: var(--m3-color-surface-container-high);
  border: 1px solid
    color-mix(in srgb, var(--m3-color-primary) 25%, transparent);
  border-radius: var(--m3-shape-md);
}
.key-value {
  flex: 1;
  margin: 0;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--m3-color-primary);
  font-family: var(--m3-typeface-mono);
  font-size: var(--m3-typescale-body-small);
  word-break: break-all;
}
</style>

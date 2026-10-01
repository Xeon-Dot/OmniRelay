<template>
  <div class="page">
    <PageHeader :title="$t('providers.title')" :subtitle="$t('providers.subtitle')">
      <OrButton v-if="isAdmin" variant="filled" @click="openDialog()">
        <template #leading>
          <OrIcon name="add" :size="15" />
        </template>
        {{ $t("providers.addProvider") }}
      </OrButton>
    </PageHeader>

    <AppAlert v-if="testResult && testResult.ok" variant="success" page>
      {{ $t("providers.testSuccess", { latency: testResult.latency_ms }) }}
    </AppAlert>
    <AppAlert v-if="testResult && !testResult.ok" variant="error" page>
      {{ testResult.error || $t("providers.testFailed") }}
    </AppAlert>
    <AppAlert v-if="syncResult" variant="success" page>{{ syncResult }}</AppAlert>
    <AppAlert v-if="syncError" variant="error" page>{{ syncError }}</AppAlert>

    <div class="table-card">
      <OrDataTable
        :headers="headers as any"
        :items="store.providers as unknown as Record<string, unknown>[]"
        :loading="store.loading"
        density="comfortable"
        :label="$t('providers.title')"
      >
        <template #item.provider_key="{ item }">
          <MonoTag>{{ item.provider_key }}</MonoTag>
        </template>
        <template #item.is_active="{ item }">
          <StatusChip :variant="item.is_active ? 'on' : 'off'">
            {{ item.is_active ? $t("providers.active") : $t("providers.inactive") }}
          </StatusChip>
        </template>
        <template #item.provider_type="{ item }">
          <OrChip>{{ item.provider_type }}</OrChip>
        </template>
        <template #item.actions="{ item }">
          <div class="row-actions">
            <OrSpinner
              v-if="testingId === item.id"
              :size="20"
              :label="$t('common.loading')"
            />
            <OrIconButton
              v-else
              icon="cable"
              :label="$t('providers.test')"
              @click="handleTest(item.id as number)"
            />
            <OrIconButton
              icon="edit"
              :label="$t('common.edit')"
              @click="openDialog(item)"
            />
            <OrSpinner
              v-if="syncingId === item.id"
              :size="20"
              :label="$t('common.loading')"
            />
            <OrIconButton
              v-else
              icon="sync"
              :label="$t('providers.sync')"
              @click="handleSync(item.id as number)"
            />
            <OrIconButton
              icon="delete"
              tone="danger"
              :label="$t('common.delete')"
              @click="handleDelete(item.id as number)"
            />
          </div>
        </template>
        <template #no-data>
          <EmptyState icon="cloud_off" :text="$t('providers.noProviders')" />
        </template>
      </OrDataTable>
    </div>

    <!-- Mobile cards -->
    <div class="mobile-cards">
      <MobileDataCard
        v-for="p in store.providers"
        :key="p.id"
        :items="[
          { label: $t('providers.key'), value: p.provider_key },
          { label: $t('providers.name'), value: p.name },
          { label: $t('providers.type'), value: p.provider_type },
          {
            label: $t('providers.status'),
            value: p.is_active
              ? $t('providers.active')
              : $t('providers.inactive'),
          },
        ]"
      >
        <template v-if="isAdmin" #actions>
          <OrSpinner
            v-if="testingId === p.id"
            :size="20"
            :label="$t('common.loading')"
          />
          <OrIconButton
            v-else
            icon="cable"
            :label="$t('providers.test')"
            @click="handleTest(p.id)"
          />
          <OrIconButton
            icon="edit"
            :label="$t('common.edit')"
            @click="openDialog(p)"
          />
          <OrSpinner
            v-if="syncingId === p.id"
            :size="20"
            :label="$t('common.loading')"
          />
          <OrIconButton
            v-else
            icon="sync"
            :label="$t('providers.sync')"
            @click="handleSync(p.id)"
          />
          <OrIconButton
            icon="delete"
            tone="danger"
            :label="$t('common.delete')"
            @click="handleDelete(p.id)"
          />
        </template>
      </MobileDataCard>
      <EmptyState
        v-if="!store.providers.length"
        icon="cloud_off"
        :text="$t('providers.noProviders')"
      />
    </div>

    <!-- Dialog -->
    <OrDialog
      v-model="dialog"
      :width="520"
      :fullscreen="isMobile"
      :title="
        editing
          ? $t('providers.editProvider')
          : $t('providers.addProvider')
      "
    >
      <div class="form-stack">
        <OrTextField
          v-model="form.provider_key"
          :label="$t('providers.providerKey')"
          placeholder="e.g. openai, my-llama"
          :disabled="!!editing"
        />
        <OrTextField
          v-model="form.name"
          :label="$t('providers.displayName')"
          placeholder="e.g. OpenAI"
        />
        <OrTextField
          v-if="form.provider_type !== 'custom'"
          v-model="form.api_base_url"
          :label="$t('providers.apiBaseUrl')"
          placeholder="https://api.openai.com/v1"
        />
        <OrTextField
          v-if="form.provider_type !== 'custom' && !editing"
          v-model="form.api_key"
          type="password"
          :label="$t('providers.apiKey')"
          placeholder="sk-..."
        />
        <div
          v-if="form.provider_type !== 'custom' && editing"
          class="field-group"
        >
          <span class="field-label">{{ $t("providers.apiKeys") }}</span>
          <div
            v-for="key in editing.api_keys"
            :key="key.id"
            class="key-row"
          >
            <span class="key-prefix" :title="$t('providers.keyPrefix')">{{
              key.key_prefix
            }}</span>
            <div
              class="key-active"
              @click.prevent="handleSetKeyActive(key.id, !key.is_active)"
            >
              <OrCheckbox
                :model-value="key.is_active"
                :label="$t('providers.active')"
              />
            </div>
            <OrIconButton
              icon="delete"
              tone="danger"
              :label="$t('common.delete')"
              @click="handleRemoveKey(key.id)"
            />
          </div>
          <div class="endpoint-row">
            <OrTextField
              v-model="newKey"
              type="password"
              class="endpoint-field"
              :label="$t('providers.apiKey')"
              placeholder="sk-..."
            />
            <OrButton
              variant="outlined"
              :disabled="!newKey"
              @click="handleAddKey"
            >
              {{ $t("providers.addKey") }}
            </OrButton>
          </div>
        </div>
        <OrSelect
          v-model="form.provider_type"
          :label="$t('providers.providerType')"
          :options="providerTypes.map((type) => ({ value: type, label: type }))"
        />
        <div
          v-if="form.provider_type !== 'custom'"
          class="field-group"
        >
          <span class="field-label">{{
            $t("providers.additionalFormats")
          }}</span>
          <p class="field-hint">{{ $t("providers.additionalFormatsHint") }}</p>
          <div
            v-for="(ep, i) in form.endpoints"
            :key="i"
            class="endpoint-row"
          >
            <OrSelect
              v-model="ep.api_type"
              class="endpoint-select"
              :label="$t('providers.type')"
              :options="
                endpointTypes.map((type) => ({ value: type, label: type }))
              "
            />
            <OrTextField
              v-model="ep.base_url"
              class="endpoint-field"
              :label="$t('providers.apiBaseUrl')"
              placeholder="https://..."
            />
            <OrIconButton
              icon="close"
              tone="danger"
              :label="$t('common.delete')"
              @click="form.endpoints.splice(i, 1)"
            />
          </div>
          <OrButton
            variant="outlined"
            @click="form.endpoints.push({ api_type: 'openai', base_url: '' })"
          >
            <template #leading>
              <OrIcon name="add" :size="14" />
            </template>
            {{ $t("providers.addFormat") }}
          </OrButton>
        </div>
        <div v-if="form.provider_type === 'custom'" class="field-group">
          <span class="field-label">{{
            $t("providers.sourceModels")
          }}</span>
          <p class="field-hint">{{ $t("providers.selectModelsHint") }}</p>
          <div class="source-models-list">
            <div
              v-for="group in store.sourceModels"
              :key="group.provider_key"
              class="source-model-group"
            >
              <div class="source-model-group-header">
                <span class="source-model-group-name">{{ group.name }}</span>
                <span class="source-model-group-type">{{
                  group.provider_type
                }}</span>
                <span class="source-model-group-key">{{
                  group.provider_key
                }}</span>
              </div>
              <div
                v-for="model in group.models"
                :key="model.model_id"
                class="source-model-item"
              >
                <OrCheckbox
                  v-model="form.source_models"
                  :value="`${group.provider_key}/${model.model_id}`"
                  :label="model.model_id"
                />
              </div>
            </div>
            <div
              v-if="!store.sourceModels.length"
              class="empty-source-models"
            >
              {{ $t("providers.noSourceModels") }}
            </div>
          </div>
        </div>
        <OrCheckbox
          v-model="form.show_in_model_list"
          :label="$t('providers.showInModelList')"
          :hint="$t('providers.showInModelListHint')"
        />
        <OrCheckbox
          v-if="form.provider_type !== 'custom'"
          v-model="form.auto_sync"
          :label="$t('providers.autoSync')"
          :hint="$t('providers.autoSyncHint')"
        />

        <AppAlert v-if="syncResult" variant="success">{{ syncResult }}</AppAlert>
        <AppAlert v-if="dialogError || store.error" variant="error">{{
          dialogError || store.error
        }}</AppAlert>
      </div>

      <template #footer>
        <OrButton variant="text" @click="dialog = false">
          {{ $t("common.cancel") }}
        </OrButton>
        <OrButton variant="filled" :loading="saving" @click="handleSave">
          {{ editing ? $t("common.update") : $t("common.create") }}
        </OrButton>
      </template>
    </OrDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useProvidersStore } from "../stores/providers";
import { useAuthStore } from "../stores/auth";
import MobileDataCard from "../components/MobileDataCard.vue";
import PageHeader from "../components/PageHeader.vue";
import EmptyState from "../components/EmptyState.vue";
import StatusChip from "../components/StatusChip.vue";
import MonoTag from "../components/MonoTag.vue";
import AppAlert from "../components/AppAlert.vue";
import { useMobile } from "../composables/useMobile";

const { t } = useI18n();
const store = useProvidersStore();
const auth = useAuthStore();
const isAdmin = computed(() => !!auth.user?.is_admin);
const { isMobile } = useMobile();

onMounted(() => {
  store.fetch();
  store.fetchSourceModels();
});

const dialog = ref(false);
const editing = ref<any>(null);
const saving = ref(false);
const dialogError = ref("");
const newKey = ref("");
const syncResult = ref("");
const syncError = ref("");
const syncingId = ref<number | null>(null);
const testingId = ref<number | null>(null);
const testResult = ref<{ ok: boolean; latency_ms: number; error?: string } | null>(null);
const providerTypes = [
  "custom",
  "openai",
  "anthropic",
  "lmstudio",
  "ollama",
  "gemini",
];
const endpointTypes = ["openai", "anthropic", "lmstudio", "ollama", "gemini"];

const form = ref({
  provider_key: "",
  name: "",
  api_base_url: "",
  api_key: "",
  provider_type: "openai",
  auto_sync: true,
  show_in_model_list: true,
  source_models: [] as string[],
  endpoints: [] as { api_type: string; base_url: string }[],
});

const headers = computed(() => {
  const cols: Record<string, unknown>[] = [
    { title: t("providers.key"), key: "provider_key", sortable: true },
    { title: t("providers.name"), key: "name", sortable: true },
    { title: t("providers.type"), key: "provider_type" },
    { title: t("providers.status"), key: "is_active" },
  ];
  if (isAdmin.value) {
    cols.push({
      title: t("providers.actions"),
      key: "actions",
      sortable: false,
      align: "end",
    });
  }
  return cols;
});

function openDialog(provider?: any) {
  dialogError.value = "";
  syncResult.value = "";
  syncError.value = "";
  newKey.value = "";
  store.clearError();
  if (provider) {
    editing.value = provider;
    form.value = {
      ...provider,
      api_key: "",
      auto_sync: false,
      source_models: provider.source_models ?? [],
      endpoints: provider.endpoints ? provider.endpoints.map((e: any) => ({ ...e })) : [],
      show_in_model_list: provider.show_in_model_list ?? true,
    };
  } else {
    editing.value = null;
    form.value = {
      provider_key: "",
      name: "",
      api_base_url: "",
      api_key: "",
      provider_type: "openai",
      auto_sync: true,
      show_in_model_list: true,
      source_models: [],
      endpoints: [],
    };
  }
  dialog.value = true;
}

async function handleSave() {
  saving.value = true;
  dialogError.value = "";
  syncResult.value = "";
  try {
    if (editing.value) {
      // For custom providers, include source_models so model selections persist.
      // For non-custom providers, omit it — models are managed via sync.
      const rest = form.value.provider_type === "custom"
        ? (({ auto_sync: _s, api_key: _k, ...r }) => r)(form.value)
        : (({ source_models: _o, auto_sync: _s, api_key: _k, ...r }) => r)(form.value);
      await store.update(editing.value.id, rest);
      if (form.value.auto_sync && form.value.provider_type !== "custom") {
        const { data } = await store.syncModels(editing.value.id);
        syncResult.value = t("providers.syncedModels", {
          count: data.model_count,
        });
      }
    } else {
      const created = await store.create(form.value);
      if (form.value.auto_sync && form.value.provider_type !== "custom") {
        const { data } = await store.syncModels(created.id);
        syncResult.value = t("providers.syncedModels", {
          count: data.model_count,
        });
      }
    }
    dialog.value = false;
  } catch (e: any) {
    dialogError.value = e.response?.data?.error || t("providers.saveFailed");
  } finally {
    saving.value = false;
  }
}

  async function handleSync(id: number) {
    syncResult.value = "";
    syncError.value = "";
    testResult.value = null;
    syncingId.value = id;
    try {
      const { data } = await store.syncModels(id);
      syncResult.value = t("providers.syncedModels", { count: data.model_count });
    } catch (e: any) {
      syncError.value = e.response?.data?.error || t("providers.syncFailed");
    } finally {
      syncingId.value = null;
    }
  }

  async function handleTest(id: number) {
    testResult.value = null;
    syncResult.value = "";
    syncError.value = "";
    testingId.value = id;
    try {
      const result = await store.testProvider(id);
      testResult.value = result;
    } catch (e: any) {
      dialogError.value = e.response?.data?.error || t("providers.testFailed");
    } finally {
      testingId.value = null;
    }
  }

function refreshEditing() {
  if (!editing.value) return;
  const p = store.providers.find((x) => x.id === editing.value.id);
  if (p) editing.value = p;
}

async function handleAddKey() {
  if (!editing.value || !newKey.value) return;
  dialogError.value = "";
  try {
    await store.addKey(editing.value.id, newKey.value);
    newKey.value = "";
    refreshEditing();
  } catch {
    dialogError.value = store.error || t("providers.saveFailed");
  }
}

async function handleSetKeyActive(keyId: number, is_active: boolean) {
  if (!editing.value) return;
  dialogError.value = "";
  try {
    await store.setKeyActive(editing.value.id, keyId, is_active);
    refreshEditing();
  } catch {
    dialogError.value = store.error || t("providers.cannotDeleteLastKey");
  }
}

async function handleRemoveKey(keyId: number) {
  if (!editing.value) return;
  dialogError.value = "";
  try {
    await store.removeKey(editing.value.id, keyId);
    refreshEditing();
  } catch {
    dialogError.value = store.error || t("providers.cannotDeleteLastKey");
  }
}

async function handleDelete(id: number) {
  if (!confirm(t("providers.deleteConfirm"))) return;
  await store.remove(id);
}
</script>

<style scoped>
@import "../styles/page-shared.css";

.form-stack {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-200);
}
.field-hint {
  margin-block-end: var(--m3-space-50);
}
.source-models-list {
  max-block-size: 300px;
  overflow-y: auto;
  border: 1px solid var(--m3-color-outline-variant);
  border-radius: var(--m3-shape-md);
  padding: var(--m3-space-200);
  margin-block-start: var(--m3-space-50);
}
.source-model-group {
  margin-block-end: var(--m3-space-200);
}
.source-model-group-header {
  display: flex;
  align-items: center;
  gap: var(--m3-space-75);
  padding-block: var(--m3-space-50);
}
.source-model-group-name {
  font: var(--m3-typescale-title-small);
  letter-spacing: var(--m3-typescale-title-small-tracking);
}
.source-model-group-type {
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  background: var(--m3-color-secondary-container);
  color: var(--m3-color-on-secondary-container);
  padding: var(--m3-space-25) var(--m3-space-75);
  border-radius: var(--m3-shape-xs);
}
.source-model-group-key {
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  color: var(--m3-color-on-surface-variant);
}
.source-model-item {
  padding: var(--m3-space-25) 0 var(--m3-space-25) var(--m3-space-150);
}
.empty-source-models {
  padding: var(--m3-space-300);
  text-align: center;
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
  color: var(--m3-color-on-surface-variant);
}
.endpoint-row {
  display: flex;
  align-items: center;
  gap: var(--m3-space-75);
  margin-block-end: var(--m3-space-75);
}
.endpoint-row .endpoint-select {
  flex: 0 0 120px;
  min-inline-size: 0;
}
.endpoint-row .endpoint-field {
  flex: 1 1 auto;
  min-inline-size: 0;
}
.key-row {
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
  margin-block-end: var(--m3-space-75);
}
.key-prefix {
  font: var(--m3-typescale-body-medium);
  font-family: var(--m3-typeface-mono);
  flex: 1;
  min-inline-size: 0;
}
.key-active {
  flex: 0 0 auto;
}
@media (max-width: 600px) {
  .table-card {
    display: none;
  }
}
</style>

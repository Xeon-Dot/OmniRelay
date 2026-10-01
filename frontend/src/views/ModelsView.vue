<template>
  <div class="page">
    <PageHeader :title="$t('models.title')" :subtitle="$t('models.subtitle')">
      <OrButton v-if="isAdmin" variant="filled" @click="openDialog()">
        <template #leading>
          <OrIcon name="add" :size="15" />
        </template>
        {{ $t("models.addModel") }}
      </OrButton>
    </PageHeader>

    <OrDataTable
      :headers="headers as any"
      :items="store.models as unknown as Record<string, unknown>[]"
      :loading="store.loading"
      density="comfortable"
      :label="$t('models.title')"
    >
      <template #item.full_id="{ item }">
        <MonoTag>{{ item.provider_key }}/{{ item.model_id }}</MonoTag>
      </template>
      <template #item.is_manual="{ item }">
        <StatusChip :variant="item.is_manual ? 'warning' : 'on'">
          {{ item.is_manual ? $t("models.manual") : $t("models.auto") }}
        </StatusChip>
      </template>
      <template #item.pricing="{ item }">
        <div class="pricing-cell">
          <span class="pricing-pair">
            <span class="pricing-key">{{ $t("models.pricingIn") }}</span>
            <span class="pricing-val">${{ item.input_price_per_1mtok }}</span>
          </span>
          <span class="pricing-sep">·</span>
          <span class="pricing-pair">
            <span class="pricing-key">{{ $t("models.pricingOut") }}</span>
            <span class="pricing-val">${{ item.output_price_per_1mtok }}</span>
          </span>
          <span v-if="item.cache_read_price_per_1mtok" class="pricing-sep"
            >·</span
          >
          <span v-if="item.cache_read_price_per_1mtok" class="pricing-pair">
            <span class="pricing-key">{{ $t("models.pricingCache") }}</span>
            <span class="pricing-val"
              >${{ item.cache_read_price_per_1mtok }}</span
            >
          </span>
        </div>
      </template>
      <template #item.context_window="{ item }">
        <span class="mono-val">{{
          item.context_window
            ? ((item.context_window as number) / 1000).toFixed(0) + "k"
            : "—"
        }}</span>
      </template>
      <template #item.actions="{ item }">
        <div class="row-actions">
          <OrIconButton
            icon="edit"
            :label="$t('common.edit')"
            @click="openEditDialog(item)"
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
        <EmptyState icon="deployed_code_off" :text="$t('models.noModels')" />
      </template>
    </OrDataTable>

    <OrDialog
      v-model="dialog"
      :width="500"
      :fullscreen="isMobile"
      :title="editMode ? $t('models.editModel') : $t('models.addModel')"
    >
      <div class="form-stack">
        <OrSelect
          v-if="!editMode"
          :model-value="
            form.provider_id === null ? '' : String(form.provider_id)
          "
          :options="
            providerOptions.map((p) => ({ value: String(p.value), label: p.text }))
          "
          :label="$t('models.provider')"
          @update:model-value="form.provider_id = Number($event)"
        />
        <OrTextField
          v-model="form.model_id"
          :label="$t('models.modelId')"
          placeholder="gpt-4o"
        />
        <OrTextField
          v-model="form.display_name"
          :label="$t('models.displayName')"
          placeholder="GPT-4 Omni"
        />
        <div class="price-grid">
          <OrTextField
            :model-value="form.input_price_per_1mtok"
            type="number"
            step="0.01"
            :label="$t('models.inputPrice')"
            @update:model-value="
              form.input_price_per_1mtok =
                ($event === '' ? '' : Number($event)) as any
            "
          />
          <OrTextField
            :model-value="form.output_price_per_1mtok"
            type="number"
            step="0.01"
            :label="$t('models.outputPrice')"
            @update:model-value="
              form.output_price_per_1mtok =
                ($event === '' ? '' : Number($event)) as any
            "
          />
          <OrTextField
            :model-value="form.cache_write_5m_price_per_1mtok"
            type="number"
            step="0.01"
            :label="$t('models.cacheWrite5m')"
            @update:model-value="
              form.cache_write_5m_price_per_1mtok =
                ($event === '' ? '' : Number($event)) as any
            "
          />
          <OrTextField
            :model-value="form.cache_write_1h_price_per_1mtok"
            type="number"
            step="0.01"
            :label="$t('models.cacheWrite1h')"
            @update:model-value="
              form.cache_write_1h_price_per_1mtok =
                ($event === '' ? '' : Number($event)) as any
            "
          />
          <OrTextField
            :model-value="form.cache_read_price_per_1mtok"
            type="number"
            step="0.01"
            :label="$t('models.cacheRead')"
            @update:model-value="
              form.cache_read_price_per_1mtok =
                ($event === '' ? '' : Number($event)) as any
            "
          />
          <OrTextField
            :model-value="form.context_window"
            type="number"
            :label="$t('models.contextWindow')"
            placeholder="128000"
            @update:model-value="
              form.context_window =
                ($event === '' ? '' : Number($event)) as any
            "
          />
        </div>
        <AppAlert v-if="dialogError" variant="error">{{ dialogError }}</AppAlert>
      </div>

      <template #footer>
        <OrButton variant="text" @click="dialog = false">
          {{ $t("common.cancel") }}
        </OrButton>
        <OrButton variant="filled" :loading="saving" @click="handleSave">
          {{ editMode ? $t("common.update") : $t("common.create") }}
        </OrButton>
      </template>
    </OrDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useModelsStore } from "../stores/models";
import { useProvidersStore } from "../stores/providers";
import { useAuthStore } from "../stores/auth";
import PageHeader from "../components/PageHeader.vue";
import EmptyState from "../components/EmptyState.vue";
import StatusChip from "../components/StatusChip.vue";
import MonoTag from "../components/MonoTag.vue";
import AppAlert from "../components/AppAlert.vue";
import { useMobile } from "../composables/useMobile";

const { t } = useI18n();
const store = useModelsStore();
const providersStore = useProvidersStore();
const auth = useAuthStore();
const isAdmin = computed(() => !!auth.user?.is_admin);
const dialog = ref(false);
const editMode = ref(false);
const editingId = ref<number | null>(null);
const saving = ref(false);
const dialogError = ref("");

const { isMobile } = useMobile();

const form = ref({
  provider_id: null as number | null,
  model_id: "",
  display_name: "",
  input_price_per_1mtok: 0,
  output_price_per_1mtok: 0,
  cache_write_5m_price_per_1mtok: 0,
  cache_write_1h_price_per_1mtok: 0,
  cache_read_price_per_1mtok: 0,
  context_window: 0,
});

const providerOptions = computed(() =>
  providersStore.providers.map((p) => ({ text: p.name, value: p.id })),
);

const headers = computed(() => {
  const cols: Record<string, unknown>[] = [
    { title: t("models.model"), key: "full_id", sortable: false },
    { title: t("models.provider"), key: "provider_key" },
    { title: t("models.source"), key: "is_manual" },
    { title: t("models.pricing"), key: "pricing", sortable: false },
    { title: t("models.context"), key: "context_window" },
  ];
  if (isAdmin.value) {
    cols.push({ title: "", key: "actions", sortable: false, width: 80 });
  }
  return cols;
});

function openDialog() {
  editMode.value = false;
  editingId.value = null;
  dialogError.value = "";
  form.value = {
    provider_id: providersStore.providers[0]?.id || null,
    model_id: "",
    display_name: "",
    input_price_per_1mtok: 0,
    output_price_per_1mtok: 0,
    cache_write_5m_price_per_1mtok: 0,
    cache_write_1h_price_per_1mtok: 0,
    cache_read_price_per_1mtok: 0,
    context_window: 0,
  };
  dialog.value = true;
}

function openEditDialog(model: any) {
  editMode.value = true;
  editingId.value = model.id;
  dialogError.value = "";
  form.value = {
    provider_id: model.provider_id,
    model_id: model.model_id,
    display_name: model.display_name,
    input_price_per_1mtok: model.input_price_per_1mtok || 0,
    output_price_per_1mtok: model.output_price_per_1mtok || 0,
    cache_write_5m_price_per_1mtok: model.cache_write_5m_price_per_1mtok || 0,
    cache_write_1h_price_per_1mtok: model.cache_write_1h_price_per_1mtok || 0,
    cache_read_price_per_1mtok: model.cache_read_price_per_1mtok || 0,
    context_window: model.context_window || 0,
  };
  dialog.value = true;
}

async function handleSave() {
  saving.value = true;
  dialogError.value = "";
  try {
    if (editMode.value && editingId.value) {
      await store.update(editingId.value, {
        display_name: form.value.display_name,
        input_price_per_1mtok: form.value.input_price_per_1mtok,
        output_price_per_1mtok: form.value.output_price_per_1mtok,
        cache_write_5m_price_per_1mtok:
          form.value.cache_write_5m_price_per_1mtok,
        cache_write_1h_price_per_1mtok:
          form.value.cache_write_1h_price_per_1mtok,
        cache_read_price_per_1mtok: form.value.cache_read_price_per_1mtok,
        context_window: form.value.context_window,
      });
    } else {
      if (!form.value.provider_id) {
        dialogError.value = t("models.selectProvider");
        saving.value = false;
        return;
      }
      await store.create({
        model_id: form.value.model_id,
        display_name: form.value.display_name || form.value.model_id,
        provider_id: form.value.provider_id,
        input_price_per_1mtok: form.value.input_price_per_1mtok,
        output_price_per_1mtok: form.value.output_price_per_1mtok,
        cache_write_5m_price_per_1mtok:
          form.value.cache_write_5m_price_per_1mtok,
        cache_write_1h_price_per_1mtok:
          form.value.cache_write_1h_price_per_1mtok,
        cache_read_price_per_1mtok: form.value.cache_read_price_per_1mtok,
        context_window: form.value.context_window,
      });
    }
    dialog.value = false;
  } catch (e: any) {
    dialogError.value = e.response?.data?.error || t("models.saveFailed");
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  if (!confirm(t("models.deleteConfirm"))) return;
  await store.remove(id);
}

onMounted(async () => {
  await providersStore.fetch();
  await store.fetch();
});
</script>

<style scoped>
@import "../styles/page-shared.css";

.form-stack {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-200);
}
.mono-val {
  font: var(--m3-typescale-body-small);
  font-family: var(--m3-typeface-mono);
  letter-spacing: var(--m3-typescale-body-small-tracking);
  color: var(--m3-color-on-surface-variant);
}
.pricing-cell {
  display: flex;
  align-items: center;
  gap: var(--m3-space-75);
  flex-wrap: wrap;
}
.pricing-pair {
  display: flex;
  align-items: center;
  gap: var(--m3-space-50);
}
.pricing-key {
  font: var(--m3-typescale-label-small);
  font-family: var(--m3-typeface-mono);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  text-transform: uppercase;
  color: var(--m3-color-on-surface-variant);
}
.pricing-val {
  font: var(--m3-typescale-label-medium);
  font-family: var(--m3-typeface-mono);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  color: var(--m3-color-on-surface);
}
.pricing-sep {
  color: var(--m3-color-outline);
}
.price-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m3-space-150);
}
@media (max-width: 600px) {
  .price-grid {
    grid-template-columns: 1fr;
  }
}
</style>

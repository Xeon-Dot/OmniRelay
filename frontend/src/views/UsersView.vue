<template>
  <div class="page">
    <PageHeader :title="$t('users.title')" :subtitle="$t('users.subtitle')" />

    <div class="table-card">
      <OrDataTable
        :headers="headers"
        :items="store.users"
        :loading="store.loading"
        density="comfortable"
        :label="$t('users.title')"
      >
        <template #item.username="{ item }">
          <span class="username-cell">
            <OrIcon name="account_circle" :size="16" class="user-icon" />
            {{ item.username }}
          </span>
        </template>
        <template #item.email="{ item }">
          <span class="dim-text">{{ item.email }}</span>
        </template>
        <template #item.is_admin="{ item }">
          <StatusChip :variant="item.is_admin ? 'on' : 'off'">
            {{ item.is_admin ? $t("users.admin") : $t("users.member") }}
          </StatusChip>
        </template>
        <template #item.created_at="{ item }">
          <span class="dim-text">{{
            new Date(item.created_at).toLocaleDateString()
          }}</span>
        </template>
        <template #item.actions="{ item }">
          <div class="row-actions">
            <OrIconButton
              size="sm"
              icon="admin_panel_settings"
              :label="$t('users.toggleRole')"
              @click="handleToggleRole(item)"
            />
            <OrIconButton
              size="sm"
              icon="restart_alt"
              :label="$t('users.resetPassword')"
              @click="handleResetPassword(item)"
            />
            <OrIconButton
              size="sm"
              icon="hub"
              :label="$t('users.providers')"
              @click="openProvidersDialog(item)"
            />
            <OrIconButton
              size="sm"
              icon="delete"
              tone="danger"
              :label="$t('users.delete')"
              @click="handleDelete(item)"
            />
          </div>
        </template>
        <template #no-data>
          <EmptyState icon="group" :text="$t('users.noUsers')" />
        </template>
      </OrDataTable>
    </div>

    <!-- Mobile cards -->
    <div class="mobile-cards">
      <MobileDataCard
        v-for="u in store.users"
        :key="u.id"
        :items="[
          { label: $t('users.username'), value: u.username },
          { label: $t('users.email'), value: u.email },
          {
            label: $t('users.role'),
            value: u.is_admin ? $t('users.admin') : $t('users.member'),
          },
          {
            label: $t('users.created'),
            value: new Date(u.created_at).toLocaleDateString(),
          },
        ]"
      >
        <template #actions>
          <OrIconButton
            size="sm"
            icon="admin_panel_settings"
            :label="$t('users.toggleRole')"
            @click="handleToggleRole(u)"
          />
          <OrIconButton
            size="sm"
            icon="restart_alt"
            :label="$t('users.resetPassword')"
            @click="handleResetPassword(u)"
          />
          <OrIconButton
            size="sm"
            icon="hub"
            :label="$t('users.providers')"
            @click="openProvidersDialog(u)"
          />
          <OrIconButton
            size="sm"
            icon="delete"
            tone="danger"
            :label="$t('users.delete')"
            @click="handleDelete(u)"
          />
        </template>
      </MobileDataCard>
      <EmptyState
        v-if="!store.users.length"
        icon="group"
        :text="$t('users.noUsers')"
      />
    </div>

    <!-- Confirm delete dialog -->
    <OrDialog
      v-model="deleteDialog"
      :width="400"
      :title="$t('users.deleteTitle')"
    >
      <div class="dialog-stack">
        <p class="confirm-text">
          {{ $t("users.deleteConfirm", { name: targetUser?.username }) }}
        </p>
        <AppAlert v-if="dialogError" variant="error">{{
          dialogError
        }}</AppAlert>
      </div>
      <template #footer>
        <OrButton variant="text" @click="deleteDialog = false">
          {{ $t("common.cancel") }}
        </OrButton>
        <OrButton variant="filled" tone="danger" :loading="busy" @click="confirmDelete">
          {{ $t("common.delete") }}
        </OrButton>
      </template>
    </OrDialog>

    <!-- Reset password result dialog -->
    <OrDialog v-model="resetDialog" :width="420" :title="$t('users.resetTitle')">
      <div class="dialog-stack">
        <p class="reveal-note">{{ $t("users.resetNote") }}</p>
        <div class="key-reveal">
          <code class="key-value">{{ resetCode }}</code>
          <OrIconButton
            :icon="copied ? 'check' : 'content_copy'"
            :variant="copied ? 'tonal' : 'standard'"
            :label="$t('common.copy')"
            @click="copyCode"
          />
        </div>
        <AppAlert v-if="dialogError" variant="error">{{
          dialogError
        }}</AppAlert>
      </div>
      <template #footer>
        <OrButton @click="resetDialog = false">
          {{ $t("common.close") }}
        </OrButton>
      </template>
    </OrDialog>

    <!-- Providers dialog -->
    <OrDialog
      v-model="providersDialog"
      :width="isMobile ? undefined : 460"
      :fullscreen="isMobile"
      :title="$t('users.providersTitle')"
    >
      <div class="dialog-stack">
        <p class="field-hint">{{ $t("users.providersHint") }}</p>
        <div v-if="providersLoading" class="loading-center">
          <OrSpinner :size="28" />
        </div>
        <div v-else class="provider-checklist">
          <div
            v-for="p in providersStore.providers"
            :key="p.id"
            class="provider-check-item"
          >
            <OrCheckbox
              v-model="selectedProviderIds"
              :value="p.id"
              :label="p.name"
            />
            <MonoTag>{{ p.provider_key }}</MonoTag>
          </div>
          <p v-if="!providersStore.providers.length" class="dim-text">
            {{ $t("users.noProviders") }}
          </p>
        </div>
        <AppAlert v-if="dialogError" variant="error">{{
          dialogError
        }}</AppAlert>
      </div>
      <template #footer>
        <OrButton variant="text" @click="providersDialog = false">
          {{ $t("common.cancel") }}
        </OrButton>
        <OrButton :loading="busy" @click="saveProviders">
          {{ $t("common.save") }}
        </OrButton>
      </template>
    </OrDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useUsersStore, type User } from "../stores/users";
import { useProvidersStore } from "../stores/providers";
import { useMobile } from "../composables/useMobile";
import { useSnackbarStore } from "../stores/snackbar";
import PageHeader from "../components/PageHeader.vue";
import StatusChip from "../components/StatusChip.vue";
import MonoTag from "../components/MonoTag.vue";
import EmptyState from "../components/EmptyState.vue";
import MobileDataCard from "../components/MobileDataCard.vue";
import AppAlert from "../components/AppAlert.vue";

const { t } = useI18n();
const store = useUsersStore();
const providersStore = useProvidersStore();
const { isMobile } = useMobile();

const headers = computed(() => [
  { title: t("users.username"), key: "username" },
  { title: t("users.email"), key: "email" },
  { title: t("users.role"), key: "is_admin" },
  { title: t("users.created"), key: "created_at" },
  { title: "", key: "actions", sortable: false, align: "end" as const },
]);

const targetUser = ref<User | null>(null);
const deleteDialog = ref(false);
const resetDialog = ref(false);
const providersDialog = ref(false);
const resetCode = ref("");
const copied = ref(false);
const busy = ref(false);
const dialogError = ref("");
const providersLoading = ref(false);
const selectedProviderIds = ref<number[]>([]);

async function handleToggleRole(u: User) {
  try {
    dialogError.value = "";
    await store.setRole(u.id, !u.is_admin);
  } catch (err: any) {
    dialogError.value = err?.response?.data?.error || err?.message;
    useSnackbarStore().error(dialogError.value);
  }
}

function handleDelete(u: User) {
  targetUser.value = u;
  dialogError.value = "";
  deleteDialog.value = true;
}

async function confirmDelete() {
  if (!targetUser.value) return;
  busy.value = true;
  dialogError.value = "";
  try {
    await store.remove(targetUser.value.id);
    deleteDialog.value = false;
  } catch (err: any) {
    dialogError.value = err?.response?.data?.error || err?.message;
  } finally {
    busy.value = false;
  }
}

async function handleResetPassword(u: User) {
  dialogError.value = "";
  try {
    resetCode.value = await store.resetPassword(u.id);
    copied.value = false;
    resetDialog.value = true;
  } catch (err: any) {
    useSnackbarStore().error(err?.response?.data?.error || err?.message);
  }
}

function copyCode() {
  navigator.clipboard.writeText(resetCode.value);
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
}

async function openProvidersDialog(u: User) {
  targetUser.value = u;
  dialogError.value = "";
  providersLoading.value = true;
  providersDialog.value = true;
  try {
    if (!providersStore.providers.length) {
      await providersStore.fetch();
    }
    selectedProviderIds.value = await store.getProviders(u.id);
  } catch (err: any) {
    dialogError.value = err?.response?.data?.error || err?.message;
  } finally {
    providersLoading.value = false;
  }
}

async function saveProviders() {
  if (!targetUser.value) return;
  busy.value = true;
  dialogError.value = "";
  try {
    await store.setProviders(targetUser.value.id, selectedProviderIds.value);
    providersDialog.value = false;
  } catch (err: any) {
    dialogError.value = err?.response?.data?.error || err?.message;
  } finally {
    busy.value = false;
  }
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

.username-cell {
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
  color: var(--m3-color-on-surface);
  font-weight: 500;
}
.user-icon {
  color: var(--m3-color-on-surface-variant);
}

.provider-checklist {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-100);
}
.provider-check-item {
  display: flex;
  align-items: center;
  gap: var(--m3-space-150);
}

.confirm-text {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}

.loading-center {
  display: flex;
  justify-content: center;
  padding: var(--m3-space-300) 0;
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

@media (max-width: 600px) {
  /* MobileDataCard replaces the table below 600px — same as ProvidersView. */
  .table-card {
    display: none;
  }
}
</style>

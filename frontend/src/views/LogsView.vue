<template>
  <div class="page">
    <PageHeader :title="$t('logs.title')" :subtitle="$t('logs.subtitle')">
      <OrButton variant="tonal" :loading="store.loading" @click="loadLogs">
        <template #leading>
          <OrIcon name="refresh" :size="15" />
        </template>
        {{ $t("common.refresh") }}
      </OrButton>
    </PageHeader>

    <!-- Filter bar -->
    <div class="filter-bar">
      <div class="filter-col">
        <OrTextField
          v-model="filters.model"
          :label="$t('logs.model')"
          :hint="$t('logs.modelFilter')"
        />
      </div>
      <div class="filter-col">
        <OrTextField
          v-model="filters.provider"
          :label="$t('logs.provider')"
          :hint="$t('logs.providerFilter')"
        />
      </div>
      <div class="filter-col">
        <OrTextField v-model="filters.from" type="date" :label="$t('logs.from')" />
      </div>
      <div class="filter-col">
        <OrTextField v-model="filters.to" type="date" :label="$t('logs.to')" />
      </div>
      <OrButton class="filter-submit" :loading="store.loading" @click="loadLogs">
        {{ $t("common.apply") }}
      </OrButton>
    </div>

    <div class="table-card">
      <OrDataTable
        :headers="headers"
        :items="store.logs"
        :loading="store.loading"
        density="compact"
        :label="$t('logs.title')"
      >
        <template #item.started_at="{ item }">
          <span class="dim-text">{{
            item.started_at ? formatTime(item.started_at) : "-"
          }}</span>
        </template>
        <template #item.completed_at="{ item }">
          <span class="dim-text">{{
            item.completed_at ? formatTime(item.completed_at) : "-"
          }}</span>
        </template>
        <template #item.duration_sec="{ item }">
          <span class="mono-val">
            {{ (item.latency_ms / 1000).toFixed(2) }}s
          </span>
        </template>
        <template #item.provider_name="{ item }">
          <span class="dim-text">{{ item.provider_name || "-" }}</span>
        </template>
        <template #item.model="{ item }">
          <MonoTag>{{ getModelName(item.model) }}</MonoTag>
        </template>
        <template #item.request_tokens="{ item }">
          <span class="mono-val">{{
            item.request_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.response_tokens="{ item }">
          <span class="mono-val">{{
            item.response_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.cache_write_5m_tokens="{ item }">
          <span class="mono-val">{{
            item.cache_write_5m_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.cache_write_1h_tokens="{ item }">
          <span class="mono-val">{{
            item.cache_write_1h_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.cache_read_tokens="{ item }">
          <span class="mono-val">{{
            item.cache_read_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.total_tokens="{ item }">
          <span class="mono-val mono-val--accent">{{
            item.total_tokens.toLocaleString()
          }}</span>
        </template>
        <template #item.cost="{ item }">
          <span class="cost-val">${{ item.cost.toFixed(6) }}</span>
        </template>
        <template #item.is_error="{ item }">
          <StatusChip :variant="item.is_error ? 'off' : 'on'">
            {{ item.is_error ? $t("logs.error") : $t("logs.ok") }}
          </StatusChip>
        </template>
        <template #no-data>
          <EmptyState icon="manage_search" :text="$t('logs.noRecords')" />
        </template>

        <template #footer>
          <div class="pager">
            <span class="mono-val">{{ store.logs.length }}</span>
            <span class="dim-text"> {{ $t("logs.of") }} </span>
            <span class="mono-val">{{ store.total }}</span>
            <span class="dim-text"> {{ $t("logs.records") }}</span>
            <div class="pagination-btns">
              <OrIconButton
                icon="chevron_left"
                :label="$t('common.previous')"
                size="sm"
                :disabled="offset <= 0"
                @click="prevPage"
              />
              <OrIconButton
                icon="chevron_right"
                :label="$t('common.next')"
                size="sm"
                :disabled="offset + limit >= store.total"
                @click="nextPage"
              />
            </div>
          </div>
        </template>
      </OrDataTable>

      <!-- Mobile cards -->
      <div class="mobile-cards">
        <MobileDataCard
          v-for="log in store.logs"
          :key="log.id"
          :items="[
            {
              label: $t('logs.startedAt'),
              value: log.started_at ? formatTime(log.started_at) : '-',
            },
            {
              label: $t('logs.completedAt'),
              value: log.completed_at ? formatTime(log.completed_at) : '-',
            },
            {
              label: $t('logs.duration'),
              value: (log.latency_ms / 1000).toFixed(2) + 's',
            },
            { label: $t('logs.provider'), value: log.provider_name || '-' },
            { label: $t('logs.model'), value: getModelName(log.model) },
            {
              label: $t('logs.totalTokens'),
              value: log.total_tokens.toLocaleString(),
            },
            { label: $t('logs.totalCost'), value: '$' + log.cost.toFixed(6) },
            {
              label: $t('logs.status'),
              value: log.is_error ? $t('logs.error') : $t('logs.ok'),
            },
          ]"
        />
        <EmptyState v-if="!store.logs.length" icon="manage_search" :text="$t('logs.noRecords')" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import { useUsageStore } from "../stores/usage";
import { useMobile } from "../composables/useMobile";
import PageHeader from "../components/PageHeader.vue";
import MonoTag from "../components/MonoTag.vue";
import StatusChip from "../components/StatusChip.vue";
import EmptyState from "../components/EmptyState.vue";
import MobileDataCard from "../components/MobileDataCard.vue";

const { t } = useI18n();
const store = useUsageStore();
const { isMobile } = useMobile();

const limit = 50;
const offset = ref(0);

const filters = ref({ model: "", provider: "", from: "", to: "" });

const headers = computed(() => [
  { title: t("logs.startedAt"), key: "started_at", minWidth: 140 },
  { title: t("logs.completedAt"), key: "completed_at", minWidth: 140 },
  { title: t("logs.duration"), key: "duration_sec", minWidth: 90 },
  { title: t("logs.provider"), key: "provider_name", minWidth: 100 },
  { title: t("logs.model"), key: "model", minWidth: 140 },
  { title: t("logs.inputTokens"), key: "request_tokens", minWidth: 90 },
  { title: t("logs.outputTokens"), key: "response_tokens", minWidth: 90 },
  {
    title: t("logs.cacheWrite5m"),
    key: "cache_write_5m_tokens",
    minWidth: 90,
  },
  {
    title: t("logs.cacheWrite1h"),
    key: "cache_write_1h_tokens",
    minWidth: 90,
  },
  { title: t("logs.cacheRead"), key: "cache_read_tokens", minWidth: 90 },
  { title: t("logs.totalTokens"), key: "total_tokens", minWidth: 90 },
  { title: t("logs.totalCost"), key: "cost", minWidth: 100 },
  { title: t("logs.status"), key: "is_error", minWidth: 80 },
]);

function formatTime(ts: string): string {
  return new Date(ts).toLocaleString();
}

function getModelName(model: string): string {
  const idx = model.indexOf("/");
  return idx >= 0 ? model.substring(idx + 1) : model;
}

async function loadLogs() {
  const params: Record<string, any> = { limit, offset: offset.value };
  if (filters.value.model) params.model = filters.value.model;
  if (filters.value.from) params.from = filters.value.from;
  if (filters.value.to) params.to = filters.value.to;
  await store.fetchLogs(params);
}

function nextPage() {
  offset.value += limit;
  loadLogs();
}

function prevPage() {
  offset.value = Math.max(0, offset.value - limit);
  loadLogs();
}

onMounted(() => {
  loadLogs();
  store.connect();
});

onUnmounted(() => {
  store.disconnect();
});
</script>

<style scoped>
@import "../styles/page-shared.css";

/* page-shared supplies the token colours for .dim-text / .mono-val /
   .cost-val / .mono-val--accent; only the type treatment differs here. */
.dim-text {
  font-size: 0.78rem;
}
.mono-val {
  font-family: var(--m3-typeface-mono);
  font-size: 0.78rem;
}
.cost-val {
  font-family: var(--m3-typeface-mono);
  font-size: 0.75rem;
}

/* The Or fields reserve a support line below the box, so lift the button
   to sit flush with the field boxes instead of the hint baseline. */
.filter-submit {
  align-self: flex-end;
  margin-block-end: calc(var(--m3-space-75, 6px) + 1.25rem);
}

/* Rendered inside OrDataTable's .table-footer slot wrapper. */
.pager {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: flex-end;
  gap: var(--m3-space-100);
}

.pagination-btns {
  display: flex;
  align-items: center;
  gap: var(--m3-space-50, 4px);
  margin-inline-start: var(--m3-space-150);
}

@media (max-width: 600px) {
  /* MobileDataCard replaces the table below 600px; the pager stays. */
  .table-card :deep(.or-table-scroll) {
    display: none;
  }
  .pager {
    justify-content: center;
  }
}
</style>

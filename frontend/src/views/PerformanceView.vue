<template>
  <div class="page">
    <PageHeader :title="$t('performance.title')" :subtitle="$t('performance.subtitle')">
      <OrButton variant="tonal" :loading="store.loading" @click="load">
        <template #leading>
          <OrIcon name="refresh" :size="15" />
        </template>
        {{ $t("common.refresh") }}
      </OrButton>
    </PageHeader>

    <!-- Filter bar -->
    <div class="filter-bar">
      <div class="filter-col">
        <OrSelect
          v-model="filters.provider_id"
          :label="$t('performance.provider')"
          :options="[
            { value: '', label: $t('performance.allProviders') },
            ...providersStore.providers.map((p) => ({
              value: String(p.id),
              label: p.name,
            })),
          ]"
        />
      </div>
      <div class="filter-col">
        <OrTextField v-model="filters.from" type="date" :label="$t('performance.from')" />
      </div>
      <div class="filter-col">
        <OrTextField v-model="filters.to" type="date" :label="$t('performance.to')" />
      </div>
      <div class="filter-col filter-col--narrow">
        <OrSelect
          v-model="filters.granularity"
          :label="$t('performance.granularity')"
          :options="[
            { value: '', label: $t('performance.auto') },
            { value: 'minute', label: 'Minute' },
            { value: 'hour', label: 'Hour' },
            { value: 'day', label: 'Day' },
          ]"
        />
      </div>
      <div class="filter-col filter-col--narrow filter-col--preset">
        <OrSegmentedButton
          v-model="preset"
          size="sm"
          :options="presets.map((p) => ({ value: p.key, label: p.label }))"
          @update:model-value="applyPreset"
        />
      </div>
      <OrButton class="filter-submit" :loading="store.loading" @click="load">
        {{ $t("common.apply") }}
      </OrButton>
    </div>

    <AppAlert v-if="store.error" variant="error" page>{{ store.error }}</AppAlert>

    <template v-if="store.data">
      <!-- Summary cards -->
      <div class="stats-grid">
        <StatCard :label="$t('performance.rpm')" :sub="`${summary.total_requests.toLocaleString()} req`">
          {{ fmtNum(summary.rpm) }}
        </StatCard>
        <StatCard :label="$t('performance.tpm')" value-class="stat-value--accent" sub="&nbsp;">
          {{ fmtNum(summary.tpm) }}
        </StatCard>
        <StatCard :label="$t('performance.avgLatency')">
          {{ fmtMs(summary.avg_latency_ms) }}
          <template #sub>{{ $t("performance.p50") }} {{ fmtMs(summary.p50_ms) }}</template>
        </StatCard>
        <StatCard :label="$t('performance.ttft')" :hint="$t('performance.ttftHint')" :value-class="summary.avg_ttft_ms === null ? 'stat-value--dim' : undefined">
          {{ summary.avg_ttft_ms === null ? "-" : fmtMs(summary.avg_ttft_ms) }}
          <template #sub>
            {{
              summary.avg_ttft_ms === null
                ? $t("performance.noStreaming")
                : summary.ttft_count.toLocaleString() + " streams"
            }}
          </template>
        </StatCard>
        <StatCard :label="`${$t('performance.p95')} / ${$t('performance.p99')}`" :sub="`${$t('performance.p99')} ${fmtMs(summary.p99_ms)}`">
          {{ fmtMs(summary.p95_ms) }}
        </StatCard>
        <StatCard :label="$t('performance.errorRate')" :value-class="summary.error_rate > 0.05 ? 'stat-value--error' : undefined" sub="&nbsp;">
          {{ (summary.error_rate * 100).toFixed(2) }}%
        </StatCard>
        <StatCard :label="$t('performance.cacheHitRate')" value-class="stat-value--cache" sub="&nbsp;">
          {{ (summary.cache_hit_rate * 100).toFixed(1) }}%
        </StatCard>
      </div>

      <!-- Timeseries charts -->
      <div class="charts-row">
        <div class="table-card chart-section">
          <h2 class="chart-heading">{{ $t("performance.throughput") }}</h2>
          <div class="chart-area">
            <Line v-if="store.data.timeseries.length" :data="throughputChart" :options="throughputOptions" />
            <EmptyState v-else icon="show_chart" :text="$t('performance.noData')" />
          </div>
        </div>
        <div class="table-card chart-section">
          <h2 class="chart-heading">{{ $t("performance.latency") }}</h2>
          <div class="chart-area">
            <Line v-if="store.data.timeseries.length" :data="latencyChart" :options="latencyOptions" />
            <EmptyState v-else icon="show_chart" :text="$t('performance.noData')" />
          </div>
        </div>
      </div>

      <!-- Breakdowns -->
      <div class="breakdown-row">
        <div class="table-card breakdown-card">
          <h2 class="chart-heading">{{ $t("performance.byProvider") }}</h2>
          <div class="table-scroll" tabindex="0" :aria-label="$t('performance.byProvider')">
            <table v-if="store.data.by_provider.length" class="perf-table">
              <thead>
                <tr>
                  <th>{{ $t("performance.provider") }}</th>
                  <th class="num">{{ $t("performance.requests") }}</th>
                  <th class="num">{{ $t("performance.tokens") }}</th>
                  <th class="num">DAVG</th>
                  <th class="num">{{ $t("performance.cost") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in store.data.by_provider" :key="r.provider_id ?? 'none'">
                  <td>{{ r.provider_name || "-" }}</td>
                  <td class="num mono-val">{{ r.requests.toLocaleString() }}</td>
                  <td class="num mono-val">{{ r.tokens.toLocaleString() }}</td>
                  <td class="num mono-val">{{ fmtMs(r.avg_latency_ms) }}</td>
                  <td class="num cost-val">${{ r.cost.toFixed(4) }}</td>
                </tr>
              </tbody>
            </table>
            <EmptyState v-else icon="show_chart" :text="$t('performance.noData')" small />
          </div>
        </div>

        <div class="table-card breakdown-card">
          <h2 class="chart-heading">{{ $t("performance.byModel") }}</h2>
          <div class="table-scroll" tabindex="0" :aria-label="$t('performance.byModel')">
            <table v-if="store.data.by_model.length" class="perf-table">
              <thead>
                <tr>
                  <th>{{ $t("performance.model") }}</th>
                  <th class="num">{{ $t("performance.requests") }}</th>
                  <th class="num">{{ $t("performance.tokens") }}</th>
                  <th class="num">DAVG</th>
                  <th class="num">{{ $t("performance.cost") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in store.data.by_model" :key="r.model">
                  <td><MonoTag>{{ r.model }}</MonoTag></td>
                  <td class="num mono-val">{{ r.requests.toLocaleString() }}</td>
                  <td class="num mono-val">{{ r.tokens.toLocaleString() }}</td>
                  <td class="num mono-val">{{ fmtMs(r.avg_latency_ms) }}</td>
                  <td class="num cost-val">${{ r.cost.toFixed(4) }}</td>
                </tr>
              </tbody>
            </table>
            <EmptyState v-else icon="show_chart" :text="$t('performance.noData')" small />
          </div>
        </div>
      </div>

      <div class="table-card breakdown-card">
        <h2 class="chart-heading">{{ $t("performance.topModels") }}</h2>
        <div class="table-scroll" tabindex="0" :aria-label="$t('performance.topModels')">
          <table v-if="store.data.top_models_by_cost.length" class="perf-table">
            <thead>
              <tr>
                <th>{{ $t("performance.model") }}</th>
                <th class="num">{{ $t("performance.requests") }}</th>
                <th class="num">{{ $t("performance.tokens") }}</th>
                <th class="num">{{ $t("performance.cost") }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in store.data.top_models_by_cost" :key="r.model">
                <td><MonoTag>{{ r.model }}</MonoTag></td>
                <td class="num mono-val">{{ r.requests.toLocaleString() }}</td>
                <td class="num mono-val">{{ r.tokens.toLocaleString() }}</td>
                <td class="num cost-val">${{ r.cost.toFixed(4) }}</td>
              </tr>
            </tbody>
          </table>
          <EmptyState v-else icon="show_chart" :text="$t('performance.noData')" small />
        </div>
      </div>
    </template>

    <EmptyState v-else-if="!store.loading && !store.error" icon="monitoring" :text="$t('performance.noData')" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from "chart.js";
import { usePerformanceStore } from "../stores/performance";
import { useProvidersStore } from "../stores/providers";
import { useChartTheme } from "../composables/useChartTheme";
import PageHeader from "../components/PageHeader.vue";
import StatCard from "../components/StatCard.vue";
import MonoTag from "../components/MonoTag.vue";
import EmptyState from "../components/EmptyState.vue";
import AppAlert from "../components/AppAlert.vue";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
);

const { t } = useI18n();
const store = usePerformanceStore();
const providersStore = useProvidersStore();
const { palette } = useChartTheme();

const filters = reactive({
  provider_id: "",
  from: "",
  to: "",
  granularity: "",
});

const preset = ref("");
let applyingPreset = false;
const presets = [
  { key: "1h", label: "1H" },
  { key: "24h", label: "24H" },
  { key: "7d", label: "7D" },
  { key: "30d", label: "30D" },
];

const summary = computed(
  () =>
    store.data?.summary ?? {
      total_requests: 0,
      rpm: 0,
      tpm: 0,
      avg_latency_ms: 0,
      p50_ms: 0,
      p95_ms: 0,
      p99_ms: 0,
      avg_ttft_ms: null,
      ttft_count: 0,
      error_rate: 0,
      cache_hit_rate: 0,
    },
);

function fmtNum(v: number): string {
  if (v >= 1000000) return (v / 1000000).toFixed(2) + "M";
  if (v >= 1000) return (v / 1000).toFixed(1) + "k";
  return v.toFixed(v < 10 && v % 1 !== 0 ? 1 : 0);
}

function fmtMs(v: number): string {
  if (v >= 60000) return (v / 60000).toFixed(1) + "m";
  if (v >= 1000) return (v / 1000).toFixed(2) + "s";
  return Math.round(v) + "ms";
}

function fmtBucket(b: string): string {
  if (b.length === 10) return b.slice(5);
  return b.slice(5, 16);
}

/** Segmented buttons emit `string | string[]`; this group is single-select. */
function applyPreset(value: string | string[]) {
  const key = Array.isArray(value) ? value[0] : value;
  if (!key) return;
  applyingPreset = true;
  preset.value = key;
  const now = new Date();
  const ms: Record<string, number> = {
    "1h": 60 * 60 * 1000,
    "24h": 24 * 60 * 60 * 1000,
    "7d": 7 * 24 * 60 * 60 * 1000,
    "30d": 30 * 24 * 60 * 60 * 1000,
  };
  const fromDate = new Date(now.getTime() - ms[key]);
  const iso = (d: Date) => d.toISOString().slice(0, 10);
  filters.from = iso(fromDate);
  filters.to = iso(now);
  filters.granularity = key === "1h" ? "minute" : key === "30d" ? "day" : "hour";
  load();
}

function load() {
  const params: Record<string, any> = {};
  if (filters.provider_id !== "") params.provider_id = filters.provider_id;
  if (filters.from) params.from = filters.from;
  if (filters.to) params.to = filters.to + " 23:59:59";
  if (filters.granularity) params.granularity = filters.granularity;
  store.fetchPerformance(params);
}

// Manual filter edits clear the preset highlight
watch(filters, () => {
  if (applyingPreset) {
    applyingPreset = false;
    return;
  }
  preset.value = "";
});

/**
 * Blend a token colour into a translucent hex (#RRGGBBAA) for chart fills.
 * Canvas cannot resolve CSS custom properties or `color-mix()`, and every
 * value here still originates from a token — nothing is hard-coded.
 */
function withAlpha(color: string, alpha: number): string {
  const hex = color.trim().replace(/^#/, "");
  if (!/^[0-9a-fA-F]{6}$/.test(hex)) return color;
  const channel = Math.round(alpha * 255)
    .toString(16)
    .padStart(2, "0");
  return `#${hex}${channel}`;
}

/** Canvas cannot read CSS custom properties — resolve the typeface tokens once. */
function typefaceToken(name: string): string {
  if (typeof window === "undefined") return "";
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}

const plainTypeface = typefaceToken("--m3-typeface-plain");
const monoTypeface = typefaceToken("--m3-typeface-mono");

const bucketLabels = computed(() =>
  (store.data?.timeseries ?? []).map((b) => fmtBucket(b.bucket)),
);

const throughputChart = computed(() => ({
  labels: bucketLabels.value,
  datasets: [
    {
      label: t("performance.rpm"),
      data: store.data?.timeseries.map((b) => b.rpm) ?? [],
      borderColor: palette.value.primary,
      backgroundColor: withAlpha(palette.value.primary, 0.06),
      fill: true,
      tension: 0.35,
      pointRadius: 2,
      pointHoverRadius: 4,
      borderWidth: 1.5,
      yAxisID: "y",
    },
    {
      label: t("performance.tpm"),
      data: store.data?.timeseries.map((b) => b.tpm) ?? [],
      borderColor: palette.value.tertiary,
      backgroundColor: withAlpha(palette.value.tertiary, 0.04),
      fill: true,
      tension: 0.35,
      pointRadius: 2,
      pointHoverRadius: 4,
      borderWidth: 1.5,
      yAxisID: "y1",
    },
  ],
}));

const latencyChart = computed(() => ({
  labels: bucketLabels.value,
  datasets: [
    {
      label: "DAVG",
      data: store.data?.timeseries.map((b) => b.avg_latency_ms) ?? [],
      borderColor: palette.value.primary,
      backgroundColor: withAlpha(palette.value.primary, 0.06),
      fill: true,
      tension: 0.35,
      pointRadius: 2,
      pointHoverRadius: 4,
      borderWidth: 1.5,
    },
    {
      label: t("performance.ttft"),
      data: store.data?.timeseries.map((b) => b.avg_ttft_ms) ?? [],
      borderColor: palette.value.tertiary,
      backgroundColor: "transparent",
      borderDash: [5, 4],
      fill: false,
      tension: 0.35,
      pointRadius: 2,
      pointHoverRadius: 4,
      borderWidth: 1.5,
      spanGaps: true,
    },
  ],
}));

const baseChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: "index" as const },
  plugins: {
    legend: {
      labels: {
        color: palette.value.text,
        font: { family: plainTypeface, size: 11 },
        usePointStyle: true,
        pointStyleWidth: 8,
        boxHeight: 6,
      },
    },
    tooltip: {
      backgroundColor: palette.value.surface,
      borderColor: palette.value.grid,
      borderWidth: 1,
      titleColor: palette.value.onSurface,
      bodyColor: palette.value.onSurface,
      titleFont: { family: plainTypeface, size: 12 },
      bodyFont: { family: monoTypeface, size: 11 },
      padding: 10,
    },
  },
  scales: {
    x: {
      grid: { color: withAlpha(palette.value.grid, 0.4) },
      ticks: {
        color: palette.value.text,
        font: { family: monoTypeface, size: 10 },
        maxRotation: 0,
        autoSkip: true,
      },
    },
  },
}));

const throughputOptions = computed(() => ({
  ...baseChartOptions.value,
  scales: {
    ...baseChartOptions.value.scales,
    y: {
      position: "left" as const,
      grid: { color: withAlpha(palette.value.grid, 0.4) },
      ticks: { color: palette.value.primary, font: { family: monoTypeface, size: 10 } },
    },
    y1: {
      position: "right" as const,
      grid: { drawOnChartArea: false },
      ticks: { color: palette.value.tertiary, font: { family: monoTypeface, size: 10 } },
    },
  },
}));

const latencyOptions = computed(() => ({
  ...baseChartOptions.value,
  scales: {
    ...baseChartOptions.value.scales,
    y: {
      grid: { color: withAlpha(palette.value.grid, 0.4) },
      ticks: {
        color: palette.value.text,
        font: { family: monoTypeface, size: 10 },
        callback: (v: number | string) => fmtMs(Number(v)),
      },
    },
  },
}));

onMounted(() => {
  providersStore.fetch();
  load();
});
</script>

<style scoped>
@import "../styles/page-shared.css";

/* ── Stats grid ── */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--m3-space-150);
}
@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

/* value-class lands on StatCard's inner element, which carries StatCard's
   scope id — reach it through :deep() so the accents actually paint. */
.stats-grid :deep(.stat-value--accent) { color: var(--m3-color-primary); }
.stats-grid :deep(.stat-value--cache) { color: var(--or-success); }
.stats-grid :deep(.stat-value--error) { color: var(--m3-color-error); }
.stats-grid :deep(.stat-value--dim) { color: var(--m3-color-on-surface-variant); }

/* ── Filter bar extras ── */
.filter-col--narrow {
  flex: 0 1 auto;
  min-inline-size: 0;
}
/* Or fields reserve a support line below the box; lift these to sit flush
   with the field boxes instead of the hint baseline. */
.filter-col--preset,
.filter-submit {
  margin-block-end: calc(var(--m3-space-75, 6px) + 1.25rem);
}

/* ── Charts ── */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m3-space-150);
}
@media (max-width: 900px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}
.chart-section {
  padding: var(--m3-space-200) var(--m3-space-300);
}
.chart-heading {
  font: var(--m3-typescale-title-small);
  letter-spacing: var(--m3-typescale-title-small-tracking);
  color: var(--m3-color-on-surface);
  margin: 0 0 var(--m3-space-150);
}
.chart-area {
  height: 240px;
  position: relative;
}

/* ── Breakdown tables ── */
.breakdown-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m3-space-150);
}
@media (max-width: 900px) {
  .breakdown-row {
    grid-template-columns: 1fr;
  }
}
.breakdown-card {
  padding: var(--m3-space-200) var(--m3-space-300);
}
.perf-table {
  width: 100%;
  border-collapse: collapse;
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}
.perf-table th {
  position: sticky;
  inset-block-start: 0;
  z-index: 1;
  text-align: left;
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
  text-transform: uppercase;
  padding: var(--m3-space-75, 6px) var(--m3-space-150);
  border-bottom: 1px solid var(--m3-color-outline-variant);
  white-space: nowrap;
  background: var(--m3-color-surface-container-low);
}
.perf-table td {
  padding: var(--m3-space-100) var(--m3-space-150);
  border-bottom: 1px solid var(--m3-color-outline-variant);
  color: var(--m3-color-on-surface);
  white-space: nowrap;
}
.perf-table tr:last-child td {
  border-bottom: none;
}
.perf-table .num {
  text-align: right;
}
.mono-val {
  font-family: var(--m3-typeface-mono);
  font-size: 0.78rem;
}
.cost-val {
  font-family: var(--m3-typeface-mono);
  font-size: 0.78rem;
}
</style>

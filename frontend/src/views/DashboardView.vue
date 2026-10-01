<template>
  <div class="dash">
    <!-- Page header -->
    <PageHeader
      :title="$t('dashboard.title')"
      :subtitle="$t('dashboard.subtitle')"
    >
      <div class="dash-time">{{ currentTime }}</div>
    </PageHeader>

    <!-- Stat cards -->
    <div class="stat-grid">
      <StatCard :label="$t('dashboard.todayCost')">
        ${{ todayCost.toFixed(4) }}
        <template #sub>
          {{ $t("dashboard.total") }}: ${{
            (stats?.total_cost ?? 0).toFixed(4)
          }}
          <OrProgressBar
            class="stat-bar"
            :value="todayCostBarPct / 100"
            tone="warning"
            :label="$t('dashboard.todayCost')"
          />
        </template>
      </StatCard>

      <StatCard :label="$t('dashboard.todayRequests')">
        {{ todayRequests.toLocaleString() }}
        <template #sub>
          {{ $t("dashboard.total") }}:
          {{ (stats?.total_requests ?? 0).toLocaleString() }}
          <OrProgressBar
            class="stat-bar"
            :value="todayRequestsBarPct / 100"
            tone="success"
            :label="$t('dashboard.todayRequests')"
          />
        </template>
      </StatCard>

      <StatCard :label="$t('dashboard.todayTokens')">
        {{ todayTokens.toLocaleString() }}
        <template #sub>
          {{ $t("dashboard.total") }}:
          {{ (stats?.total_tokens ?? 0).toLocaleString() }}
          <OrProgressBar
            class="stat-bar"
            :value="todayTokensBarPct / 100"
            tone="primary"
            :label="$t('dashboard.todayTokens')"
          />
        </template>
      </StatCard>

      <StatCard :label="$t('dashboard.performance')">
        <span class="perf-stack">
          <span>RPM: {{ rpm.toFixed(0) }}</span>
          <span>TPM: {{ tpm.toFixed(0) }}</span>
        </span>
        <template #sub>
          &nbsp;
          <OrProgressBar
            class="stat-bar"
            :value="perfBarPct / 100"
            tone="danger"
            :label="$t('dashboard.performance')"
          />
        </template>
      </StatCard>

      <div class="stat-cell">
        <StatCard :label="$t('dashboard.avgLatency')">
          {{ latencySec.toFixed(2) }} s
          <template #sub>
            &nbsp;
            <OrProgressBar
              class="stat-bar"
              :value="latencyBarPct / 100"
              :tone="latencySec > 2 ? 'danger' : 'warning'"
              :label="$t('dashboard.avgLatency')"
            />
          </template>
        </StatCard>
        <OrIcon
          class="stat-cell__icon"
          name="speed"
          :size="16"
          :color="
            latencySec > 2 ? 'var(--m3-color-error)' : 'var(--or-warning)'
          "
        />
      </div>
    </div>

    <!-- Charts row -->
    <div class="charts-row">
      <OrCard>
        <div class="panel-header">
          <span class="panel-title">{{ $t("dashboard.usage30Days") }}</span>
          <div class="legend">
            <span class="legend-dot legend-dot--tokens" />
            <span class="legend-label">{{ $t("dashboard.tokens") }}</span>
            <span class="legend-dot legend-dot--cost" />
            <span class="legend-label">{{ $t("dashboard.cost") }}</span>
          </div>
        </div>
        <div class="chart-wrap">
          <Line v-if="chartData" :data="chartData" :options="chartOptions" />
          <div v-else class="chart-empty">
            <OrIcon class="chart-empty__icon" name="show_chart" :size="28" />
            <span>{{ $t("dashboard.noUsageData") }}</span>
          </div>
        </div>
      </OrCard>

      <OrCard>
        <div class="panel-header">
          <span class="panel-title">{{ $t("dashboard.system") }}</span>
          <span class="status-badge">
            <span class="status-dot" />
            {{ $t("dashboard.operational") }}
          </span>
        </div>
        <div class="info-rows">
          <div v-for="row in infoRows" :key="row.label" class="info-row">
            <div class="info-row__left">
              <OrIcon class="info-row__icon" :name="row.icon" :size="14" />
              <span class="info-row__label">{{ row.label }}</span>
            </div>
            <span class="info-row__value">{{ row.value }}</span>
          </div>
        </div>
      </OrCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";
import { useUsageStore } from "../stores/usage";
import { useChartTheme } from "../composables/useChartTheme";
import PageHeader from "../components/PageHeader.vue";
import StatCard from "../components/StatCard.vue";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
);

const { t } = useI18n();
const usageStore = useUsageStore();
const stats = computed(() => usageStore.stats);

const currentTime = ref("");
function tick() {
  currentTime.value = new Date().toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}
tick();
setInterval(tick, 1000);

const todayCost = computed(() => stats.value?.today_cost ?? 0);
const todayRequests = computed(() => stats.value?.today_requests ?? 0);
const todayTokens = computed(() => stats.value?.today_tokens ?? 0);
const rpm = computed(() => stats.value?.rpm ?? 0);
const tpm = computed(() => stats.value?.tpm ?? 0);
const latency = computed(() => stats.value?.avg_latency_ms ?? 0);
const latencySec = computed(() => latency.value / 1000);

const todayCostBarPct = computed(() => {
  const total = stats.value?.total_cost ?? 0;
  return total > 0 ? Math.min(100, (todayCost.value / total) * 100) : 0;
});
const todayRequestsBarPct = computed(() => {
  const total = stats.value?.total_requests ?? 0;
  return total > 0 ? Math.min(100, (todayRequests.value / total) * 100) : 0;
});
const todayTokensBarPct = computed(() => {
  const total = stats.value?.total_tokens ?? 0;
  return total > 0 ? Math.min(100, (todayTokens.value / total) * 100) : 0;
});
const perfBarPct = computed(() => Math.min(100, (rpm.value / 1000) * 100));
const latencyBarPct = computed(() =>
  Math.min(100, (latency.value / 5000) * 100),
);

const infoRows = computed(() => [
  {
    label: t("dashboard.activeProviders"),
    value: String(stats.value?.providers_count ?? 0),
    icon: "dns",
  },
  {
    label: t("dashboard.registeredModels"),
    value: String(stats.value?.models_count ?? 0),
    icon: "deployed_code",
  },
  {
    label: t("dashboard.activeApiKeys"),
    value: String(stats.value?.active_keys ?? 0),
    icon: "key",
  },
  {
    label: t("logs.totalRequests"),
    value: (stats.value?.total_requests ?? 0).toLocaleString(),
    icon: "send",
  },
]);

const { palette } = useChartTheme();

const chartData = computed(() => {
  if (!stats.value?.daily_usage?.length) return null;
  return {
    labels: stats.value.daily_usage.map((d) => d.date.slice(5)),
    datasets: [
      {
        label: t("dashboard.tokens"),
        data: stats.value.daily_usage.map((d) => d.total_tokens),
        borderColor: palette.value.primary,
        backgroundColor: `${palette.value.primary}0f`,
        fill: true,
        tension: 0.4,
        pointRadius: 3,
        pointHoverRadius: 5,
        pointBackgroundColor: palette.value.primary,
        borderWidth: 1.5,
        yAxisID: "y",
      },
      {
        label: t("dashboard.costAxis"),
        data: stats.value.daily_usage.map((d) => d.total_cost),
        borderColor: palette.value.tertiary,
        backgroundColor: `${palette.value.tertiary}0a`,
        fill: true,
        tension: 0.4,
        pointRadius: 3,
        pointHoverRadius: 5,
        pointBackgroundColor: palette.value.tertiary,
        borderWidth: 1.5,
        yAxisID: "y1",
      },
    ],
  };
});

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: "index" as const },
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: palette.value.surface,
      borderColor: palette.value.grid,
      borderWidth: 1,
      titleColor: palette.value.onSurface,
      bodyColor: palette.value.text,
      titleFont: { size: 12 },
      bodyFont: { size: 11 },
      padding: 10,
    },
  },
  scales: {
    x: {
      grid: { color: palette.value.grid, drawBorder: false },
      ticks: { color: palette.value.text, font: { size: 10 } },
    },
    y: {
      type: "linear" as const,
      position: "left" as const,
      grid: { color: palette.value.grid, drawBorder: false },
      ticks: { color: palette.value.text, font: { size: 10 } },
    },
    y1: {
      type: "linear" as const,
      position: "right" as const,
      grid: {
        color: palette.value.grid,
        drawOnChartArea: false,
        drawBorder: false,
      },
      ticks: { color: palette.value.text, font: { size: 10 } },
    },
  },
}));

onMounted(() => {
  usageStore.fetchStats();
  usageStore.connect();
});

onUnmounted(() => {
  usageStore.disconnect();
});
</script>

<style scoped>
.dash {
  display: flex;
  flex-direction: column;
  gap: var(--or-gap-page);
  max-inline-size: var(--or-layout-max-width);
}

/* ── Header ── */
.dash-time {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  font-family: var(--m3-typeface-mono);
  font-variant-numeric: tabular-nums;
}

/* ── Stat grid ── */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--or-gap-card);
}

.stat-cell {
  position: relative;
  display: flex;
}
.stat-cell__icon {
  position: absolute;
  inset-block-start: var(--or-pad-card);
  inset-inline-end: var(--or-pad-card);
  padding: var(--m3-space-75);
  border-radius: var(--m3-shape-sm);
  background: color-mix(in srgb, currentColor 14%, transparent);
}

.stat-bar {
  margin-block-start: var(--m3-space-75);
}

.perf-stack {
  display: flex;
  flex-direction: column;
  gap: var(--m3-space-25);
}

/* ── Charts row ── */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: var(--or-gap-card);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m3-space-150);
  margin-block-end: var(--m3-space-200);
}
.panel-title {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-title-medium);
  letter-spacing: var(--m3-typescale-title-medium-tracking);
}

.legend {
  display: flex;
  align-items: center;
  gap: var(--m3-space-150);
}
.legend-dot {
  inline-size: var(--m3-space-75);
  block-size: var(--m3-space-75);
  border-radius: var(--m3-shape-full);
}
.legend-dot--tokens {
  background: var(--m3-color-primary);
}
.legend-dot--cost {
  background: var(--m3-color-tertiary);
}
.legend-label {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
}

.chart-wrap {
  block-size: 240px;
  position: relative;
}
.chart-empty {
  block-size: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--m3-space-150);
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}
.chart-empty__icon {
  color: var(--m3-color-outline);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-75);
  padding: var(--m3-space-25) var(--m3-space-150);
  border-radius: var(--m3-shape-full);
  background: var(--or-success-container);
  color: var(--or-on-success-container);
  font: var(--m3-typescale-label-small);
  letter-spacing: var(--m3-typescale-label-small-tracking);
}
.status-dot {
  inline-size: var(--m3-space-75);
  block-size: var(--m3-space-75);
  border-radius: var(--m3-shape-full);
  background: var(--or-success);
  animation: pulse calc(var(--m3-duration-extra-long4) * 2)
    var(--m3-easing-standard) infinite;
}
@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.info-rows {
  display: flex;
  flex-direction: column;
}
.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m3-space-150);
  padding-block: var(--m3-space-150);
  border-block-end: 1px solid var(--m3-color-outline-variant);
}
.info-row:last-child {
  border-block-end: none;
}
.info-row__left {
  display: flex;
  align-items: center;
  gap: var(--m3-space-100);
}
.info-row__icon {
  color: var(--m3-color-outline);
}
.info-row__label {
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
}
.info-row__value {
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-title-small);
  letter-spacing: var(--m3-typescale-title-small-tracking);
  font-family: var(--m3-typeface-mono);
  font-variant-numeric: tabular-nums;
}

@media (max-width: 839px) {
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .charts-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
  .dash {
    gap: var(--m3-space-300);
  }
  .chart-wrap {
    block-size: 220px;
  }
  .charts-row {
    gap: var(--m3-space-150);
  }
}
</style>

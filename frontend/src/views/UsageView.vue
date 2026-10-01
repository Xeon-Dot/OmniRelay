<template>
  <div class="page">
    <PageHeader :title="$t('usage.title')" :subtitle="$t('usage.subtitle')">
      <OrButton variant="tonal" @click="store.fetchStats()">
        <template #leading>
          <OrIcon name="refresh" :size="15" />
        </template>
        {{ $t("common.refresh") }}
      </OrButton>
    </PageHeader>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <StatCard :label="$t('dashboard.totalRequests')">
        {{ stats?.total_requests?.toLocaleString() ?? "-" }}
      </StatCard>
      <StatCard :label="$t('dashboard.totalTokens')" value-class="stat-value--accent">
        {{ stats?.total_tokens?.toLocaleString() ?? "-" }}
      </StatCard>
      <StatCard :label="$t('dashboard.totalCost')" value-class="stat-value--cost">
        ${{ stats?.total_cost?.toFixed(4) ?? "-" }}
      </StatCard>
      <StatCard :label="$t('dashboard.avgLatency')">
        {{ stats ? (stats.avg_latency_ms / 1000).toFixed(2) + "s" : "-" }}
      </StatCard>
      <StatCard :label="$t('usage.cacheWrite5m')">
        {{ stats?.total_cache_write_5m?.toLocaleString() ?? "-" }}
      </StatCard>
      <StatCard :label="$t('usage.cacheWrite1h')">
        {{ stats?.total_cache_write_1h?.toLocaleString() ?? "-" }}
      </StatCard>
      <StatCard :label="$t('usage.cacheRead')" value-class="stat-value--cache-read">
        {{ stats?.total_cache_read?.toLocaleString() ?? "-" }}
      </StatCard>
    </div>

    <!-- 30-day Chart -->
    <OrCard>
      <h2 class="chart-heading">{{ $t("dashboard.usage30Days") }}</h2>
      <div class="chart-area">
        <template v-if="stats?.daily_usage?.length">
          <Line :data="chartData" :options="chartOptions" />
        </template>
        <EmptyState v-else icon="bar_chart" :text="$t('dashboard.noUsageData')" />
      </div>
    </OrCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from "vue";
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
import { useUsageStore } from "../stores/usage";
import { useChartTheme } from "../composables/useChartTheme";
import PageHeader from "../components/PageHeader.vue";
import StatCard from "../components/StatCard.vue";
import EmptyState from "../components/EmptyState.vue";

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
const store = useUsageStore();
const stats = computed(() => store.stats);
const { palette } = useChartTheme();

const chartData = computed(() => ({
  labels: stats.value?.daily_usage.map((d) => d.date.slice(5)) ?? [],
  datasets: [
    {
      label: t("dashboard.tokens"),
      data: stats.value?.daily_usage.map((d) => d.total_tokens) ?? [],
      backgroundColor: `${palette.value.primary}1a`,
      borderColor: palette.value.primary,
      borderWidth: 2,
      fill: true,
      tension: 0.3,
      pointRadius: 3,
      pointHoverRadius: 5,
      yAxisID: "y",
    },
    {
      label: t("dashboard.cost"),
      data: stats.value?.daily_usage.map((d) => d.total_cost) ?? [],
      backgroundColor: `${palette.value.tertiary}1a`,
      borderColor: palette.value.tertiary,
      borderWidth: 2,
      fill: true,
      tension: 0.3,
      pointRadius: 3,
      pointHoverRadius: 5,
      yAxisID: "y1",
    },
  ],
}));

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { mode: "index" as const, intersect: false },
  plugins: {
    legend: {
      labels: {
        color: palette.value.text,
        font: { size: 12 },
        usePointStyle: true,
        pointStyleWidth: 8,
      },
    },
    tooltip: {
      backgroundColor: palette.value.surface,
      titleColor: palette.value.onSurface,
      bodyColor: palette.value.onSurface,
      borderColor: palette.value.grid,
      borderWidth: 1,
      padding: 12,
    },
  },
  scales: {
    x: {
      ticks: { color: palette.value.text, font: { size: 10 } },
      grid: { color: palette.value.grid },
    },
    y: {
      type: "linear" as const,
      display: true,
      position: "left" as const,
      title: {
        display: true,
        text: t("dashboard.tokens"),
        color: palette.value.primary,
      },
      ticks: { color: palette.value.primary, font: { size: 10 } },
      grid: { color: palette.value.grid },
    },
    y1: {
      type: "linear" as const,
      display: true,
      position: "right" as const,
      title: {
        display: true,
        text: t("dashboard.costAxis"),
        color: palette.value.tertiary,
      },
      ticks: { color: palette.value.tertiary, font: { size: 10 } },
      grid: { color: palette.value.grid, drawOnChartArea: false },
    },
  },
}));

onMounted(() => {
  store.fetchStats();
  store.connect();
});

onUnmounted(() => {
  store.disconnect();
});
</script>

<style scoped>
@import "../styles/page-shared.css";

.stats-grid {
  display: grid;
  /* 220px keeps a 9-digit token count on one line at the title-large size;
     160px squeezed six cards per row and split the numbers. */
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--m3-space-150);
}

.stat-value--accent {
  color: var(--m3-color-primary);
}

.stat-value--cost {
  color: var(--m3-color-tertiary);
}

.stat-value--cache-read {
  color: var(--m3-color-tertiary);
}

.chart-heading {
  margin: 0 0 var(--m3-space-200);
  color: var(--m3-color-on-surface);
  font: var(--m3-typescale-title-medium);
  letter-spacing: var(--m3-typescale-title-medium-tracking);
}

.chart-area {
  inline-size: 100%;
  block-size: 320px;
}

@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--m3-space-100);
  }
  .chart-area {
    block-size: 240px;
  }
}
@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  .chart-area {
    block-size: 200px;
  }
}
</style>

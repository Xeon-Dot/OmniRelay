<script setup lang="ts">
/**
 * Data table — replacement for the legacy Vuetify data table.
 *
 * Contract (spec §4.2) is deliberately the subset the five call sites use:
 * `headers` + `items` + `loading` + `density`, per-column `#item.<key>` slots
 * and `#no-data`. No built-in pagination (LogsView rolls its own into
 * `#footer`), no row selection, no server-side sorting.
 *
 * Sorting is client-side and internal, matching what the legacy table did for the
 * Providers and Models screens; a `sort` event is emitted so a parent can take
 * over if it ever needs to. `aria-sort` is kept on the header so the state is
 * announced, and the scroll container is focusable so keyboard users can pan
 * wide tables.
 */
import { computed, ref } from "vue";
import OrSpinner from "./OrSpinner.vue";

type TableColumn = {
  key: string;
  title: string;
  sortable?: boolean;
  align?: "start" | "end";
  minWidth?: number;
  width?: number;
};

const props = withDefaults(
  defineProps<{
    headers: TableColumn[];
    /**
     * Row objects. Column keys are dynamic (`#item.<key>`) and the five call
     * sites pass different row shapes, so a precise row type cannot be
     * expressed without generics; `any[]` keeps call sites from casting every
     * cell. Matches the legacy data table's `any[]`.
     */
    items: any[];
    loading?: boolean;
    density?: "compact" | "comfortable";
    hideHeader?: boolean;
    /** Accessible name for the scroll region and table. */
    label?: string;
  }>(),
  {
    loading: false,
    density: "comfortable",
    hideHeader: false,
    label: undefined,
  },
);

const emit = defineEmits<{ sort: [key: string, direction: "asc" | "desc"] }>();

const sortKey = ref<string | null>(null);
const sortDirection = ref<"asc" | "desc">("asc");

function toggleSort(header: TableColumn) {
  if (!header.sortable) return;
  if (sortKey.value === header.key) {
    sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
  } else {
    sortKey.value = header.key;
    sortDirection.value = "asc";
  }
  emit("sort", header.key, sortDirection.value);
}

function ariaSort(header: TableColumn): "ascending" | "descending" | "none" | undefined {
  if (!header.sortable) return undefined;
  if (sortKey.value !== header.key) return "none";
  return sortDirection.value === "asc" ? "ascending" : "descending";
}

function compare(a: unknown, b: unknown): number {
  if (a == null && b == null) return 0;
  if (a == null) return -1;
  if (b == null) return 1;
  if (typeof a === "number" && typeof b === "number") return a - b;
  return String(a).localeCompare(String(b), undefined, { numeric: true });
}

const rows = computed(() => {
  const key = sortKey.value;
  if (!key) return props.items;
  const direction = sortDirection.value === "asc" ? 1 : -1;
  // Copy before sorting: mutating the prop would fight the parent's store.
  return [...props.items].sort(
    (left, right) => direction * compare(left[key], right[key]),
  );
});

const isEmpty = computed(() => !props.loading && rows.value.length === 0);
</script>

<template>
  <div
    class="or-table-scroll"
    role="region"
    :aria-label="label"
    :aria-busy="loading ? 'true' : undefined"
    tabindex="0"
  >
    <table class="or-table" :class="`or-table--${density}`">
      <thead v-if="!hideHeader">
        <tr>
          <th
            v-for="header in headers"
            :key="header.key"
            scope="col"
            :style="{
              minWidth: header.minWidth ? `${header.minWidth}px` : undefined,
              width: header.width ? `${header.width}px` : undefined,
              textAlign: header.align ?? 'start',
            }"
            :aria-sort="ariaSort(header)"
          >
            <button
              v-if="header.sortable"
              class="or-table__sort"
              type="button"
              @click="toggleSort(header)"
            >
              <span>{{ header.title }}</span>
              <OrIcon
                :name="ariaSort(header) === 'descending' ? 'arrow_downward' : 'arrow_upward'"
                :size="16"
                :class="[
                  'or-table__sort-icon',
                  { 'is-active': sortKey === header.key },
                ]"
              />
            </button>
            <span v-else>{{ header.title }}</span>
          </th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="(item, index) in rows" :key="item.id ?? index" class="or-table__row">
          <td
            v-for="header in headers"
            :key="header.key"
            :style="{ textAlign: header.align ?? 'start' }"
          >
            <slot
              :name="`item.${header.key}`"
              :item="item"
              :value="item[header.key]"
              :index="index"
            >
              {{ item[header.key] }}
            </slot>
          </td>
        </tr>

        <tr v-if="isEmpty">
          <td class="or-table__empty" :colspan="headers.length">
            <slot name="no-data">
              <span class="or-table__empty-text">—</span>
            </slot>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="loading" class="or-table__loading" aria-hidden="true">
      <OrSpinner :size="28" />
    </div>
  </div>

  <div v-if="$slots.footer" class="table-footer">
    <slot name="footer" />
  </div>
  <div v-if="$slots.bottom">
    <slot name="bottom" />
  </div>
</template>

<style scoped>
.or-table-scroll {
  position: relative;
  overflow-x: auto;
  border-radius: var(--m3-shape-2xl);
}

.or-table-scroll:focus-visible {
  outline: var(--m3-focus-ring-width) solid var(--m3-focus-ring-color);
  outline-offset: calc(var(--m3-focus-ring-offset) * -1);
}

.or-table {
  inline-size: 100%;
  border-collapse: collapse;
  background: var(--m3-color-surface-container-lowest);
  color: var(--m3-color-on-surface);
}

.or-table thead th {
  position: sticky;
  inset-block-start: 0;
  z-index: 1;
  padding: var(--m3-space-200) var(--m3-space-300);
  background: var(--m3-color-surface-container-low);
  color: var(--m3-color-on-surface-variant);
  font: var(--m3-typescale-label-medium);
  letter-spacing: var(--m3-typescale-label-medium-tracking);
  text-align: start;
  white-space: nowrap;
  border-block-end: 1px solid var(--m3-color-outline-variant);
}

.or-table__sort {
  display: inline-flex;
  align-items: center;
  gap: var(--m3-space-75, 6px);
  color: inherit;
  font: inherit;
  letter-spacing: inherit;
  cursor: pointer;
}

.or-table__sort:hover {
  color: var(--m3-color-on-surface);
}

.or-table__sort-icon {
  opacity: 0.35;
}

.or-table__sort-icon.is-active {
  opacity: 1;
  color: var(--m3-color-primary);
}

.or-table tbody td {
  padding: var(--m3-space-200) var(--m3-space-300);
  border-block-end: 1px solid var(--m3-color-outline-variant);
  font: var(--m3-typescale-body-medium);
  letter-spacing: var(--m3-typescale-body-medium-tracking);
  vertical-align: middle;
  /* Cells are atomic: wide tables pan horizontally (the scroll region above)
     rather than wrapping numbers, dates and model names mid-token. */
  white-space: nowrap;
}

.or-table--compact tbody td {
  padding: var(--m3-space-100) var(--m3-space-300);
}

.or-table__row {
  transition: background-color var(--or-motion-effects-fast);
}

.or-table__row:hover td {
  background: color-mix(in srgb, var(--m3-color-primary) 6%, transparent);
}

.or-table__row:last-child td {
  border-block-end: none;
}

.or-table__empty {
  padding: var(--m3-space-600) var(--m3-space-300) !important;
  text-align: center;
  /* The empty slot may hold a sentence; let it wrap like normal prose. */
  white-space: normal;
}

.or-table__empty-text {
  color: var(--m3-color-on-surface-variant);
}

.or-table__loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--m3-color-surface) 72%, transparent);
  color: var(--m3-color-primary);
}
</style>

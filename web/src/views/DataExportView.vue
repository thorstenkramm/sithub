<template>
  <div class="page-container" data-cy="data-export-page">
    <PageHeader
      :title="$t('dataExport.title')"
      :subtitle="$t('dataExport.subtitle')"
      :breadcrumbs="[{ text: $t('common.home'), to: '/' }, { text: $t('dataExport.title') }]"
    />

    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-column flex-sm-row flex-wrap align-start ga-4">
          <DatePickerField
            v-model="fromDate"
            :label="$t('dataExport.fromDate')"
            :max="toDate"
            data-cy="export-from"
            class="export-filter-field"
          />
          <DatePickerField
            v-model="toDate"
            :label="$t('dataExport.toDate')"
            :min="fromDate"
            data-cy="export-to"
            class="export-filter-field"
          />
          <v-select
            v-model="area"
            :items="areaOptions"
            :label="$t('dataExport.area')"
            class="export-filter-field"
            data-cy="export-area"
          />
        </div>
      </v-card-text>
    </v-card>

    <v-card>
      <!-- Order is fixed by the server (date, then name), so no column is sortable. -->
      <v-data-table-server
        :headers="headers"
        :items="rows"
        :items-length="total"
        :items-per-page="PAGE_SIZE"
        :page="page"
        :loading="loading"
        :no-data-text="$t('dataExport.emptyTable')"
        :row-props="rowProps"
        item-value="id"
        density="comfortable"
        hide-default-footer
        data-cy="export-table"
      >
        <template #[`item.name`]="{ item }">
          <span data-cy="export-name">{{ item.name || $t('dataExport.unknownUser') }}</span>
        </template>
        <template #[`item.email`]="{ item }">
          <span data-cy="export-email">{{ item.email }}</span>
        </template>
        <template #[`item.bookedItem`]="{ item }">
          <span data-cy="export-item">{{ item.bookedItem }}</span>
        </template>
        <template #[`item.bookingDate`]="{ item }">
          <span data-cy="export-date">{{ item.bookingDate }}</span>
        </template>
      </v-data-table-server>
    </v-card>

    <v-pagination
      v-model="page"
      :length="pageCount"
      class="mt-4"
      data-cy="export-pagination"
    />

    <v-snackbar
      v-model="errorOpen"
      :timeout="6000"
      color="error"
      location="bottom"
      closable
      data-cy="export-error-snackbar"
    >
      <span class="api-error-text">{{ errorMessage }}</span>
    </v-snackbar>
    <v-btn
      class="mt-4"
      prepend-icon="$download"
      :loading="isExporting"
      :disabled="isExporting || !filtersValid"
      data-cy="data-export-csv-button"
      @click="exportCsv"
    >
      {{ $t('dataExport.exportCsv') }}
    </v-btn>
    <v-snackbar v-model="csvSuccess" color="success" :timeout="3000" location="bottom"
      data-cy="data-export-csv-success">
      {{ $t('dataExport.exportCsvSuccess') }}
    </v-snackbar>
    <v-snackbar v-model="csvErrorOpen" color="error" :timeout="6000" location="bottom"
      data-cy="data-export-csv-error">
      <span class="api-error-text">{{ csvError }}</span>
      <template #actions>
        <v-btn variant="text" data-cy="data-export-csv-error-close" @click="csvErrorOpen = false">
          {{ $t('common.close') }}
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { fetchAreas } from '../api/areas';
import { downloadBookingsCsv, fetchBookingExport, type BookingExportAttributes } from '../api/bookings';
import type { JsonApiResource } from '../api/types';
import { formatApiErrorMessage } from '../api/client';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { formatMediumDate, toLocalISODate } from '../utils/date';
import { DatePickerField, PageHeader } from '../components';
import { saveBlob } from '../utils/download';

const PAGE_SIZE = 20;
/** The default range covers the last 5 days including today. */
const DEFAULT_RANGE_DAYS = 5;

const { t, locale } = useI18n();
const { loading, run } = useApi();
const { handleAuthError } = useAuthErrorHandler();

const today = new Date();
const rangeStart = new Date(today);
rangeStart.setDate(rangeStart.getDate() - (DEFAULT_RANGE_DAYS - 1));

const fromDate = ref(toLocalISODate(rangeStart));
const toDate = ref(toLocalISODate(today));
const area = ref('');
const areas = ref<{ id: string; name: string }[]>([]);
const bookings = ref<JsonApiResource<BookingExportAttributes>[]>([]);
const total = ref(0);
const page = ref(1);
const errorOpen = ref(false);
const errorMessage = ref('');
const isExporting = ref(false);
const csvSuccess = ref(false);
const csvErrorOpen = ref(false);
const csvError = ref('');
const filtersValid = computed(() => validDate(fromDate.value) && validDate(toDate.value) &&
  fromDate.value <= toDate.value && (!area.value || areas.value.some((entry) => entry.id === area.value)));

function validDate(value: string) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;
  const date = new Date(`${value}T00:00:00Z`);
  return !Number.isNaN(date.getTime()) && date.toISOString().slice(0, 10) === value;
}

async function exportCsv() {
  if (isExporting.value || !filtersValid.value) return;
  isExporting.value = true;
  csvSuccess.value = false;
  csvErrorOpen.value = false;
  try {
    const result = await downloadBookingsCsv({ from: fromDate.value, to: toDate.value, area: area.value });
    saveBlob(result.blob, result.filename);
    csvSuccess.value = true;
  } catch (err) {
    if (await handleAuthError(err)) return;
    csvError.value = formatApiErrorMessage(t('dataExport.exportCsvFailed'), err);
    csvErrorOpen.value = true;
  } finally {
    isExporting.value = false;
  }
}

const areaOptions = computed(() => [
  { title: t('dataExport.allAreas'), value: '' },
  ...areas.value.map((entry) => ({ title: entry.name, value: entry.id }))
]);

const rows = computed(() =>
  bookings.value.map((resource) => ({
    id: resource.id,
    name: resource.attributes.user_name,
    email: resource.attributes.user_email,
    bookedItem: resource.attributes.booked_item,
    bookingDate: formatMediumDate(resource.attributes.booking_date, locale.value || undefined)
  }))
);

const pageCount = computed(() => Math.max(1, Math.ceil(total.value / PAGE_SIZE)));

const headers = computed(() => [
  { title: t('dataExport.colName'), key: 'name', sortable: false },
  { title: t('dataExport.colEmail'), key: 'email', sortable: false },
  { title: t('dataExport.colItem'), key: 'bookedItem', sortable: false },
  { title: t('dataExport.colDate'), key: 'bookingDate', sortable: false }
]);

// Gives every table row a stable selector for end-to-end tests.
function rowProps() {
  return { 'data-cy': 'export-row' };
}

function showError(err: unknown) {
  errorMessage.value = formatApiErrorMessage(t('dataExport.loadFailedPrefix'), err);
  errorOpen.value = true;
}

async function loadAreas() {
  try {
    const resp = await fetchAreas();
    areas.value = resp.data.map((resource) => ({ id: resource.id, name: resource.attributes.name }));
  } catch (err) {
    if (await handleAuthError(err)) return;
    showError(err);
  }
}

// Guards against out-of-order responses when filters or the page change quickly.
let latestRequest = 0;

async function loadBookings() {
  const requestId = ++latestRequest;
  try {
    const resp = await run(() => fetchBookingExport({
      from: fromDate.value,
      to: toDate.value,
      area: area.value,
      limit: PAGE_SIZE,
      offset: (page.value - 1) * PAGE_SIZE
    }));
    if (requestId !== latestRequest) return;
    total.value = resp.meta.total;
    bookings.value = resp.data;
  } catch (err) {
    if (requestId !== latestRequest) return;
    if (await handleAuthError(err)) return;
    showError(err);
  }
}

// A changed filter always restarts at page 1; the page watcher then performs the reload.
watch([fromDate, toDate, area], () => {
  if (page.value === 1) {
    void loadBookings();
  } else {
    page.value = 1;
  }
});
watch(page, loadBookings);

onMounted(() => {
  void loadAreas();
  void loadBookings();
});
</script>

<style scoped>
.export-filter-field {
  min-width: 180px;
  max-width: 280px;
  width: 100%;
}
</style>

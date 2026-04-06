<template>
  <div class="page-container">
    <PageHeader
      :title="$t('history.title')"
      :subtitle="$t('history.subtitle')"
      :breadcrumbs="[{ text: $t('common.home'), to: '/' }, { text: $t('history.title') }]"
    />

    <!-- Date Filter Card -->
    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-column flex-sm-row flex-wrap align-end ga-4">
          <DatePickerField
            v-model="fromDate"
            :label="$t('history.fromDate')"
            data-cy="from-date"
            class="history-date-field"
          />
          <DatePickerField
            v-model="toDate"
            :label="$t('history.toDate')"
            data-cy="to-date"
            class="history-date-field"
          />
          <v-btn
            color="primary"
            variant="tonal"
            :loading="historyLoading"
            data-cy="filter-btn"
            @click="loadHistory"
          >
            <v-icon start>$search</v-icon>
            {{ $t('history.filter') }}
          </v-btn>
        </div>
      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="historyLoading" type="list" :count="5" data-cy="history-loading" />

    <!-- Error State -->
    <v-alert v-else-if="historyError" type="error" class="mb-4" data-cy="history-error">
      {{ historyError }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!bookings.length"
      :title="$t('history.emptyTitle')"
      :message="$t('history.emptyMessage')"
      icon="$calendar"
      data-cy="history-empty"
    />

    <!-- History List -->
    <v-card v-else data-cy="history-list">
      <v-list lines="two">
        <v-list-item
          v-for="booking in bookings"
          :key="booking.id"
          data-cy="history-item"
          :data-cy-booking-id="booking.id"
        >
          <template #prepend>
            <v-avatar :color="getBookingColor(booking)" variant="tonal" size="40">
              <v-icon size="20">$desk</v-icon>
            </v-avatar>
          </template>
          <v-list-item-title class="d-flex align-center flex-wrap ga-2">
            {{ booking.attributes.item_name }}
            <StatusChip
              v-if="booking.attributes.is_guest"
              status="guest"
              size="x-small"
            />
          </v-list-item-title>
          <v-list-item-subtitle>
            <span>{{ booking.attributes.item_group_name }} &bull; {{ booking.attributes.area_name }}</span>
            <br />
            <span class="text-primary">{{ formatDate(booking.attributes.booking_date) }}</span>
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { fetchBookingHistory, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import { isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, DatePickerField, StatusChip } from '../components';

const authStore = useAuthStore();
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const { loading: historyLoading, error: historyError, run: runHistory } = useApi();
const { handleAuthError } = useAuthErrorHandler();
const { locale } = useI18n();

// Default: last 30 days
const today = new Date();
const thirtyDaysAgo = new Date(today);
thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

const fromDate = ref(formatDateISO(thirtyDaysAgo));
const toDate = ref(formatDateISO(new Date(today.getTime() - 24 * 60 * 60 * 1000))); // Yesterday

function formatDateISO(date: Date) {
  return date.toISOString().slice(0, 10);
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr + 'T00:00:00');
  return date.toLocaleDateString(locale.value || undefined, {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

const getBookingColor = (booking: JsonApiResource<MyBookingAttributes>) => {
  if (booking.attributes.is_guest) return 'warning';
  return 'primary';
};

const loadHistory = async () => {
  try {
    const resp = await runHistory(() =>
      fetchBookingHistory({ from: fromDate.value, to: toDate.value })
    );
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
  }
};

onMounted(async () => {
  try {
    const resp = await fetchMe();
    authStore.userName = resp.data.attributes.display_name;
    authStore.isAdmin = resp.data.attributes.is_admin;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (isConnectionError(err)) {
      historyError.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    throw err;
  }

  await loadHistory();
});
</script>

<style scoped>
.history-date-field {
  min-width: 160px;
  max-width: 280px;
  width: 100%;
}

@media (min-width: 600px) {
  .history-date-field {
    max-width: 200px;
  }
}
</style>

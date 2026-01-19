<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="history-title">
            Booking History
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
          <v-card-text>
            <div class="d-flex flex-wrap gap-2 mb-4">
              <div>
                <label class="text-caption font-weight-medium" for="from-date">From</label>
                <input
                  id="from-date"
                  v-model="fromDate"
                  class="d-block mt-1"
                  type="date"
                  data-cy="from-date"
                />
              </div>
              <div>
                <label class="text-caption font-weight-medium" for="to-date">To</label>
                <input
                  id="to-date"
                  v-model="toDate"
                  class="d-block mt-1"
                  type="date"
                  data-cy="to-date"
                />
              </div>
              <div class="d-flex align-end">
                <v-btn
                  color="primary"
                  size="small"
                  variant="tonal"
                  :loading="historyLoading"
                  data-cy="filter-btn"
                  @click="loadHistory"
                >
                  Filter
                </v-btn>
              </div>
            </div>
            <v-progress-linear
              v-if="historyLoading"
              class="mb-3"
              indeterminate
              data-cy="history-loading"
            />
            <v-alert v-else-if="historyError" type="error" variant="tonal" data-cy="history-error">
              Unable to load booking history.
            </v-alert>
            <div v-else>
              <v-list v-if="bookings.length" data-cy="history-list">
                <v-list-item
                  v-for="booking in bookings"
                  :key="booking.id"
                  data-cy="history-item"
                  :data-cy-booking-id="booking.id"
                >
                  <v-list-item-title>
                    {{ booking.attributes.desk_name }}
                    <v-chip
                      v-if="booking.attributes.is_guest"
                      size="x-small"
                      color="warning"
                      variant="tonal"
                      class="ml-2"
                    >
                      Guest
                    </v-chip>
                  </v-list-item-title>
                  <v-list-item-subtitle>
                    <div>{{ booking.attributes.room_name }} - {{ booking.attributes.area_name }}</div>
                    <div>{{ formatDate(booking.attributes.booking_date) }}</div>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="history-empty">
                No bookings found in this date range.
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchBookingHistory, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const router = useRouter();
const { loading: historyLoading, error: historyError, run: runHistory } = useApi();

// Default: last 30 days
const today = new Date();
const thirtyDaysAgo = new Date(today);
thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

const fromDate = ref(formatDateISO(thirtyDaysAgo));
const toDate = ref(formatDateISO(new Date(today.getTime() - 24 * 60 * 60 * 1000))); // Yesterday

function formatDateISO(date: Date) {
  return date.toISOString().split('T')[0];
}

const handleAuthError = async (err: unknown) => {
  if (err instanceof ApiError && err.status === 401) {
    window.location.href = '/oauth/login';
    return true;
  }
  if (err instanceof ApiError && err.status === 403) {
    await router.push('/access-denied');
    return true;
  }
  return false;
};

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr + 'T00:00:00');
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
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
    userName.value = resp.data.attributes.display_name;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    throw err;
  }

  await loadHistory();
});
</script>

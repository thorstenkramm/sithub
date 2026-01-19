<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="my-bookings-title">
            My Bookings
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
          <v-card-text>
            <v-progress-linear
              v-if="bookingsLoading"
              class="mb-3"
              indeterminate
              data-cy="bookings-loading"
              aria-label="Loading bookings"
            />
            <v-alert v-else-if="bookingsError" type="error" variant="tonal" data-cy="bookings-error">
              Unable to load bookings.
            </v-alert>
            <div v-else>
              <v-list v-if="bookings.length" data-cy="bookings-list">
                <v-list-item
                  v-for="booking in bookings"
                  :key="booking.id"
                  data-cy="booking-item"
                >
                  <v-list-item-title>
                    {{ booking.attributes.desk_name }}
                  </v-list-item-title>
                  <v-list-item-subtitle>
                    <div data-cy="booking-location">
                      {{ booking.attributes.room_name }} - {{ booking.attributes.area_name }}
                    </div>
                    <div data-cy="booking-date">
                      {{ formatDate(booking.attributes.booking_date) }}
                    </div>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="bookings-empty">
                No upcoming bookings.
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
import { fetchMyBookings, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const router = useRouter();
const { loading: bookingsLoading, error: bookingsError, run: runBookings } = useApi();

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

  try {
    const resp = await runBookings(() => fetchMyBookings());
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
  }
});
</script>

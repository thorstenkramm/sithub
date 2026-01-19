<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="room-bookings-title">
            Room Bookings
          </v-card-title>
          <v-card-text>
            <div class="mb-4">
              <label class="text-caption font-weight-medium" for="bookings-date">Date</label>
              <input
                id="bookings-date"
                v-model="selectedDate"
                class="d-block mt-1"
                type="date"
                data-cy="bookings-date"
                aria-label="Select date"
              />
            </div>
            <v-progress-linear
              v-if="loading"
              class="mb-3"
              indeterminate
              data-cy="bookings-loading"
              aria-label="Loading bookings"
            />
            <v-alert v-else-if="errorMessage" type="error" variant="tonal" data-cy="bookings-error">
              {{ errorMessage }}
            </v-alert>
            <div v-else>
              <v-list v-if="bookings.length" data-cy="bookings-list">
                <v-list-item
                  v-for="booking in bookings"
                  :key="booking.id"
                  data-cy="booking-item"
                >
                  <v-list-item-title>{{ booking.attributes.desk_name }}</v-list-item-title>
                  <v-list-item-subtitle>
                    Booked by: {{ booking.attributes.user_name || 'Unknown' }}
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="bookings-empty">
                No bookings for this date.
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchRoomBookings } from '../api/roomBookings';
import type { RoomBookingAttributes } from '../api/roomBookings';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const bookings = ref<JsonApiResource<RoomBookingAttributes>[]>([]);
const errorMessage = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const route = useRoute();
const router = useRouter();
const { loading, run } = useApi();
const activeRoomId = ref<string | null>(null);

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

const loadBookings = async (roomId: string, date: string) => {
  errorMessage.value = null;
  try {
    const resp = await run(() => fetchRoomBookings(roomId, date));
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      errorMessage.value = 'Room not found.';
      return;
    }
    errorMessage.value = 'Unable to load bookings.';
  }
};

onMounted(async () => {
  const roomId = route.params.roomId;
  if (typeof roomId !== 'string' || roomId.trim() === '') {
    errorMessage.value = 'Room not found.';
    return;
  }

  activeRoomId.value = roomId;
  await loadBookings(roomId, selectedDate.value);
});

watch(
  selectedDate,
  async (value) => {
    if (!activeRoomId.value) {
      return;
    }
    await loadBookings(activeRoomId.value, value);
  },
  { flush: 'post' }
);

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

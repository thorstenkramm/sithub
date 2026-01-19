<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="desks-title">
            Desks
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
          <v-card-text>
            <div class="mb-4">
              <label class="text-caption font-weight-medium" for="desks-date">Date</label>
              <input
                id="desks-date"
                v-model="selectedDate"
                class="d-block mt-1"
                type="date"
                data-cy="desks-date"
                aria-label="Select booking date"
              />
            </div>
            <v-alert
              v-if="bookingSuccessMessage"
              type="success"
              variant="tonal"
              class="mb-3"
              closable
              data-cy="booking-success"
              @click:close="bookingSuccessMessage = null"
            >
              {{ bookingSuccessMessage }}
            </v-alert>
            <v-alert
              v-if="bookingErrorMessage"
              type="error"
              variant="tonal"
              class="mb-3"
              closable
              data-cy="booking-error"
              @click:close="bookingErrorMessage = null"
            >
              {{ bookingErrorMessage }}
            </v-alert>
            <v-progress-linear
              v-if="desksLoading"
              class="mb-3"
              indeterminate
              data-cy="desks-loading"
              aria-label="Loading desks"
            />
            <v-alert v-else-if="desksErrorMessage" type="error" variant="tonal" data-cy="desks-error">
              {{ desksErrorMessage }}
            </v-alert>
            <div v-else>
              <v-list v-if="desks.length" data-cy="desks-list">
                <v-list-item
                  v-for="desk in desks"
                  :key="desk.id"
                  data-cy="desk-item"
                  :data-cy-desk-id="desk.id"
                  :data-cy-availability="desk.attributes.availability"
                >
                  <v-list-item-title>{{ desk.attributes.name }}</v-list-item-title>
                  <v-list-item-subtitle>
                    <ul class="pl-4" data-cy="desk-equipment">
                      <li v-for="item in desk.attributes.equipment" :key="item">{{ item }}</li>
                    </ul>
                    <div v-if="desk.attributes.warning" class="text-caption mt-1" data-cy="desk-warning">
                      {{ desk.attributes.warning }}
                    </div>
                    <div class="text-caption mt-1" data-cy="desk-status">
                      Status: {{ desk.attributes.availability === 'occupied' ? 'Occupied' : 'Available' }}
                    </div>
                  </v-list-item-subtitle>
                  <template #append>
                    <v-btn
                      v-if="desk.attributes.availability === 'available'"
                      color="primary"
                      size="small"
                      variant="tonal"
                      :loading="bookingDeskId === desk.id"
                      :disabled="bookingDeskId !== null"
                      data-cy="book-desk-btn"
                      @click="bookDesk(desk.id)"
                    >
                      Book
                    </v-btn>
                  </template>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="desks-empty">No desks available.</div>
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
import { createBooking } from '../api/bookings';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import type { DeskAttributes } from '../api/desks';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const desks = ref<JsonApiResource<DeskAttributes>[]>([]);
const desksErrorMessage = ref<string | null>(null);
const bookingSuccessMessage = ref<string | null>(null);
const bookingErrorMessage = ref<string | null>(null);
const bookingDeskId = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const route = useRoute();
const router = useRouter();
const { loading: desksLoading, run: runDesks } = useApi();
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

const ensureDate = (value: string) => {
  if (value.trim() !== '') {
    return value;
  }
  const today = formatDate(new Date());
  if (selectedDate.value !== today) {
    selectedDate.value = today;
  }
  return today;
};

const loadDesks = async (roomId: string, date: string) => {
  desksErrorMessage.value = null;
  try {
    const normalizedDate = ensureDate(date);
    const resp = await runDesks(() => fetchDesks(roomId, normalizedDate));
    desks.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      desksErrorMessage.value = 'Room not found.';
      return;
    }
    desksErrorMessage.value = 'Unable to load desks.';
  }
};

const bookDesk = async (deskId: string) => {
  bookingSuccessMessage.value = null;
  bookingErrorMessage.value = null;
  bookingDeskId.value = deskId;

  try {
    await createBooking(deskId, selectedDate.value);
    bookingSuccessMessage.value = 'Desk booked successfully!';

    // Reload desks to reflect updated availability
    if (activeRoomId.value) {
      await loadDesks(activeRoomId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 409) {
      bookingErrorMessage.value = 'This desk is already booked for the selected date.';
    } else if (err instanceof ApiError && err.status === 404) {
      bookingErrorMessage.value = 'Desk not found.';
    } else {
      bookingErrorMessage.value = 'Unable to book desk. Please try again.';
    }
  } finally {
    bookingDeskId.value = null;
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

  const roomId = route.params.roomId;
  if (typeof roomId !== 'string' || roomId.trim() === '') {
    desksErrorMessage.value = 'Room not found.';
    return;
  }

  activeRoomId.value = roomId;
  await loadDesks(roomId, selectedDate.value);
});

watch(
  selectedDate,
  async (value) => {
    if (!activeRoomId.value) {
      return;
    }
    await loadDesks(activeRoomId.value, value);
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

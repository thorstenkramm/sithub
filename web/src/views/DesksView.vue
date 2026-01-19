<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="desks-title">
            Desks
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
            <router-link
              v-if="activeRoomId"
              :to="`/rooms/${activeRoomId}/bookings`"
              class="text-caption ml-4"
              data-cy="view-room-bookings"
            >
              View Room Bookings
            </router-link>
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
            <v-radio-group v-model="bookingType" inline density="compact" class="mb-2">
              <v-radio label="Book for myself" value="self" data-cy="book-self-radio" />
              <v-radio label="Book for a colleague" value="colleague" data-cy="book-colleague-radio" />
              <v-radio label="Book for a guest" value="guest" data-cy="book-guest-radio" />
            </v-radio-group>
            <div v-if="bookingType === 'colleague'" class="mb-4 pl-4">
              <v-text-field
                v-model="colleagueId"
                label="Colleague ID (email)"
                density="compact"
                variant="outlined"
                class="mb-2"
                data-cy="colleague-id-input"
                placeholder="e.g. jane.doe@example.com"
              />
              <v-text-field
                v-model="colleagueName"
                label="Colleague Name"
                density="compact"
                variant="outlined"
                data-cy="colleague-name-input"
                placeholder="e.g. Jane Doe"
              />
            </div>
            <div v-if="bookingType === 'guest'" class="mb-4 pl-4">
              <v-text-field
                v-model="guestName"
                label="Guest Name"
                density="compact"
                variant="outlined"
                class="mb-2"
                data-cy="guest-name-input"
                placeholder="e.g. John Visitor"
              />
              <v-text-field
                v-model="guestEmail"
                label="Guest Email (optional)"
                density="compact"
                variant="outlined"
                data-cy="guest-email-input"
                placeholder="e.g. visitor@example.com"
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
                      <span
                        v-if="
                          authStore.isAdmin &&
                          desk.attributes.availability === 'occupied' &&
                          desk.attributes.booking_id
                        "
                        data-cy="desk-booker"
                      >
                        (Booked)
                      </span>
                    </div>
                  </v-list-item-subtitle>
                  <template #append>
                    <v-btn
                      v-if="desk.attributes.availability === 'available'"
                      color="primary"
                      size="small"
                      variant="tonal"
                      :loading="bookingDeskId === desk.id"
                      :disabled="bookingDeskId !== null || cancelingBookingId !== null"
                      data-cy="book-desk-btn"
                      @click="bookDesk(desk.id)"
                    >
                      Book
                    </v-btn>
                    <v-btn
                      v-if="
                        authStore.isAdmin &&
                        desk.attributes.availability === 'occupied' &&
                        desk.attributes.booking_id
                      "
                      color="error"
                      size="small"
                      variant="tonal"
                      :loading="cancelingBookingId === desk.attributes.booking_id"
                      :disabled="bookingDeskId !== null || cancelingBookingId !== null"
                      data-cy="admin-cancel-btn"
                      @click="adminCancelBooking(desk.attributes.booking_id!)"
                    >
                      Cancel
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
import {
  createBooking,
  cancelBooking,
  type BookOnBehalfOptions,
  type GuestBookingOptions
} from '../api/bookings';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import type { DeskAttributes } from '../api/desks';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthStore } from '../stores/useAuthStore';

const authStore = useAuthStore();
const userName = ref('');
const desks = ref<JsonApiResource<DeskAttributes>[]>([]);
const desksErrorMessage = ref<string | null>(null);
const bookingSuccessMessage = ref<string | null>(null);
const bookingErrorMessage = ref<string | null>(null);
const bookingDeskId = ref<string | null>(null);
const cancelingBookingId = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const route = useRoute();
const router = useRouter();
const { loading: desksLoading, run: runDesks } = useApi();
const activeRoomId = ref<string | null>(null);
const bookingType = ref<'self' | 'colleague' | 'guest'>('self');
const colleagueId = ref('');
const colleagueName = ref('');
const guestName = ref('');
const guestEmail = ref('');

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

  // Validate colleague fields if booking on behalf
  if (bookingType.value === 'colleague') {
    if (!colleagueId.value.trim() || !colleagueName.value.trim()) {
      bookingErrorMessage.value = 'Please enter both colleague ID and name.';
      return;
    }
  }

  // Validate guest fields
  if (bookingType.value === 'guest') {
    if (!guestName.value.trim()) {
      bookingErrorMessage.value = 'Please enter the guest name.';
      return;
    }
  }

  bookingDeskId.value = deskId;

  try {
    const onBehalf: BookOnBehalfOptions | undefined =
      bookingType.value === 'colleague'
        ? { forUserId: colleagueId.value.trim(), forUserName: colleagueName.value.trim() }
        : undefined;

    const guest: GuestBookingOptions | undefined =
      bookingType.value === 'guest'
        ? { guestName: guestName.value.trim(), guestEmail: guestEmail.value.trim() || undefined }
        : undefined;

    await createBooking(deskId, selectedDate.value, onBehalf, guest);

    if (bookingType.value === 'colleague') {
      bookingSuccessMessage.value = `Desk booked successfully for ${colleagueName.value}!`;
      // Reset colleague fields after successful booking
      colleagueId.value = '';
      colleagueName.value = '';
      bookingType.value = 'self';
    } else if (bookingType.value === 'guest') {
      bookingSuccessMessage.value = `Desk booked successfully for guest ${guestName.value}!`;
      // Reset guest fields after successful booking
      guestName.value = '';
      guestEmail.value = '';
      bookingType.value = 'self';
    } else {
      bookingSuccessMessage.value = 'Desk booked successfully!';
    }

    // Reload desks to reflect updated availability
    if (activeRoomId.value) {
      await loadDesks(activeRoomId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 409) {
      // Use backend's detail message if available, otherwise a generic message
      const detail = err.detail || 'This desk is no longer available for the selected date.';
      bookingErrorMessage.value = `${detail} Please choose another desk.`;

      // Refresh desk list so user sees updated availability
      if (activeRoomId.value) {
        await loadDesks(activeRoomId.value, selectedDate.value);
      }
    } else if (err instanceof ApiError && err.status === 404) {
      bookingErrorMessage.value = 'Desk not found.';
    } else {
      bookingErrorMessage.value = 'Unable to book desk. Please try again.';
    }
  } finally {
    bookingDeskId.value = null;
  }
};

const adminCancelBooking = async (bookingId: string) => {
  bookingSuccessMessage.value = null;
  bookingErrorMessage.value = null;
  cancelingBookingId.value = bookingId;

  try {
    await cancelBooking(bookingId);
    bookingSuccessMessage.value = 'Booking cancelled successfully.';

    // Reload desks to reflect updated availability
    if (activeRoomId.value) {
      await loadDesks(activeRoomId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      bookingErrorMessage.value = 'Booking not found or already cancelled.';
    } else {
      bookingErrorMessage.value = 'Unable to cancel booking. Please try again.';
    }
  } finally {
    cancelingBookingId.value = null;
  }
};

onMounted(async () => {
  try {
    const resp = await fetchMe();
    userName.value = resp.data.attributes.display_name;
    authStore.userName = resp.data.attributes.display_name;
    authStore.isAdmin = resp.data.attributes.is_admin;
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

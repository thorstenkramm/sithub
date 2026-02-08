<template>
  <div class="page-container">
    <PageHeader
      title="Desks"
      subtitle="Select a desk to book for your chosen date"
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <v-btn
          v-if="activeRoomId"
          variant="text"
          size="small"
          :to="`/rooms/${activeRoomId}/bookings`"
          data-cy="view-room-bookings"
        >
          View Room Bookings
        </v-btn>
      </template>
    </PageHeader>

    <!-- Date Selection & Booking Options -->
    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-wrap align-end ga-4 mb-4">
          <DatePickerField
            v-model="selectedDate"
            label="Booking Date"
            :min="todayDate"
            data-cy="desks-date"
            style="max-width: 280px;"
          />
        </div>

        <!-- Booking Type Selection -->
        <v-radio-group v-model="bookingType" inline density="compact" class="mb-2" hide-details>
          <v-radio label="Book for myself" value="self" data-cy="book-self-radio" />
          <v-radio label="Book for colleague" value="colleague" data-cy="book-colleague-radio" />
          <v-radio label="Book for guest" value="guest" data-cy="book-guest-radio" />
        </v-radio-group>

        <!-- Colleague Fields -->
        <v-expand-transition>
          <div v-if="bookingType === 'colleague'" class="mt-4 d-flex flex-wrap ga-4">
            <v-text-field
              v-model="colleagueId"
              label="Colleague Email"
              density="compact"
              data-cy="colleague-id-input"
              placeholder="jane.doe@example.com"
              style="max-width: 280px;"
            />
            <v-text-field
              v-model="colleagueName"
              label="Colleague Name"
              density="compact"
              data-cy="colleague-name-input"
              placeholder="Jane Doe"
              style="max-width: 280px;"
            />
          </div>
        </v-expand-transition>

        <!-- Guest Fields -->
        <v-expand-transition>
          <div v-if="bookingType === 'guest'" class="mt-4 d-flex flex-wrap ga-4">
            <v-text-field
              v-model="guestName"
              label="Guest Name"
              density="compact"
              data-cy="guest-name-input"
              placeholder="John Visitor"
              style="max-width: 280px;"
            />
            <v-text-field
              v-model="guestEmail"
              label="Guest Email (optional)"
              density="compact"
              data-cy="guest-email-input"
              placeholder="visitor@example.com"
              style="max-width: 280px;"
            />
          </div>
        </v-expand-transition>

        <!-- Multi-day booking -->
        <v-checkbox
          v-model="multiDayBooking"
          label="Book multiple days"
          density="compact"
          hide-details
          class="mt-2"
          data-cy="multi-day-checkbox"
        />
        <v-expand-transition>
          <div v-if="multiDayBooking" class="mt-2">
            <v-text-field
              v-model="additionalDates"
              label="Additional Dates (comma-separated)"
              density="compact"
              data-cy="additional-dates-input"
              placeholder="2026-01-21, 2026-01-22"
              hint="Format: YYYY-MM-DD. Selected date above will be included."
              persistent-hint
              style="max-width: 400px;"
            />
          </div>
        </v-expand-transition>
      </v-card-text>
    </v-card>

    <!-- Success/Error Messages -->
    <v-alert
      v-if="bookingSuccessMessage"
      type="success"
      class="mb-4"
      closable
      data-cy="booking-success"
      @click:close="bookingSuccessMessage = null"
    >
      {{ bookingSuccessMessage }}
    </v-alert>
    <v-alert
      v-if="bookingErrorMessage"
      type="error"
      class="mb-4"
      closable
      data-cy="booking-error"
      @click:close="bookingErrorMessage = null"
    >
      {{ bookingErrorMessage }}
    </v-alert>

    <!-- Loading State -->
    <LoadingState v-if="desksLoading" type="cards" :count="6" data-cy="desks-loading" />

    <!-- Error State -->
    <v-alert v-else-if="desksErrorMessage" type="error" class="mb-4" data-cy="desks-error">
      {{ desksErrorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!desks.length"
      title="No desks available"
      message="This room doesn't have any desks configured yet."
      icon="$desk"
      data-cy="desks-empty"
    />

    <!-- Desks Grid -->
    <div v-else class="card-grid" data-cy="desks-list">
      <v-card
        v-for="desk in desks"
        :key="desk.id"
        :class="[
          'desk-card',
          { 'desk-available': desk.attributes.availability === 'available' },
          { 'desk-occupied': desk.attributes.availability === 'occupied' }
        ]"
        data-cy="desk-item"
        :data-cy-desk-id="desk.id"
        :data-cy-availability="desk.attributes.availability"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar
              :color="desk.attributes.availability === 'available' ? 'success' : 'warning'"
              variant="tonal"
              size="48"
            >
              <v-icon size="24">$desk</v-icon>
            </v-avatar>
          </template>
          <v-card-title class="d-flex align-center">
            {{ desk.attributes.name }}
            <StatusChip
              :status="desk.attributes.availability === 'available' ? 'available' : 'booked'"
              size="x-small"
              class="ml-2"
              data-cy="desk-status"
            />
          </v-card-title>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Equipment -->
          <div v-if="desk.attributes.equipment?.length" class="mb-2" data-cy="desk-equipment">
            <div class="text-caption text-medium-emphasis mb-1">Equipment</div>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="item in desk.attributes.equipment"
                :key="item"
                size="x-small"
                variant="outlined"
              >
                {{ item }}
              </v-chip>
            </div>
          </div>

          <!-- Warning -->
          <v-alert
            v-if="desk.attributes.warning"
            type="warning"
            variant="tonal"
            density="compact"
            class="mt-2"
            data-cy="desk-warning"
          >
            {{ desk.attributes.warning }}
          </v-alert>

          <!-- Booked by (admin only) -->
          <div
            v-if="authStore.isAdmin && desk.attributes.availability === 'occupied' && desk.attributes.booking_id"
            class="text-caption text-medium-emphasis mt-2"
            data-cy="desk-booker"
          >
            Booked for this date
          </div>
        </v-card-text>

        <v-card-actions class="px-4 pb-4">
          <v-btn
            v-if="desk.attributes.availability === 'available'"
            color="primary"
            variant="flat"
            block
            :loading="bookingDeskId === desk.id"
            :disabled="bookingDeskId !== null || cancelingBookingId !== null"
            data-cy="book-desk-btn"
            @click="bookDesk(desk.id)"
          >
            Book This Desk
          </v-btn>
          <v-btn
            v-else-if="authStore.isAdmin && desk.attributes.booking_id"
            color="error"
            variant="tonal"
            block
            :loading="cancelingBookingId === desk.attributes.booking_id"
            :disabled="bookingDeskId !== null || cancelingBookingId !== null"
            data-cy="admin-cancel-btn"
            @click="adminCancelBooking(desk.attributes.booking_id!)"
          >
            Cancel Booking
          </v-btn>
          <div v-else class="text-center w-100 text-caption text-medium-emphasis py-2">
            Not available for {{ formattedDate }}
          </div>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue';
import { useRoute } from 'vue-router';
import { ApiError } from '../api/client';
import {
  createBooking,
  createMultiDayBooking,
  cancelBooking,
  type BookOnBehalfOptions,
  type GuestBookingOptions
} from '../api/bookings';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import { fetchRooms } from '../api/rooms';
import { fetchAreas } from '../api/areas';
import type { DeskAttributes } from '../api/desks';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, StatusChip, DatePickerField } from '../components';

const authStore = useAuthStore();
const desks = ref<JsonApiResource<DeskAttributes>[]>([]);
const desksErrorMessage = ref<string | null>(null);
const bookingSuccessMessage = ref<string | null>(null);
const bookingErrorMessage = ref<string | null>(null);
const bookingDeskId = ref<string | null>(null);
const cancelingBookingId = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const todayDate = formatDate(new Date());
const route = useRoute();
const { loading: desksLoading, run: runDesks } = useApi();
const activeRoomId = ref<string | null>(null);
const areaName = ref('');
const roomName = ref('');
const bookingType = ref<'self' | 'colleague' | 'guest'>('self');
const colleagueId = ref('');
const colleagueName = ref('');
const guestName = ref('');
const guestEmail = ref('');
const multiDayBooking = ref(false);
const additionalDates = ref('');

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area', to: areaName.value ? undefined : '/' },
  { text: roomName.value || 'Room' }
]);

const formattedDate = computed(() => {
  const date = new Date(selectedDate.value);
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
});

const { handleAuthError } = useAuthErrorHandler();

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

    // Handle multi-day booking
    if (multiDayBooking.value && additionalDates.value.trim()) {
      const dates = [selectedDate.value];
      additionalDates.value.split(',').forEach((d) => {
        const trimmed = d.trim();
        if (trimmed && !dates.includes(trimmed)) {
          dates.push(trimmed);
        }
      });

      const result = await createMultiDayBooking(deskId, dates, onBehalf, guest);
      const createdCount = result.created.length;
      const conflictCount = result.conflicts?.length || 0;

      if (conflictCount > 0) {
        bookingSuccessMessage.value = `Created ${createdCount} booking(s). ${conflictCount} date(s) had conflicts.`;
        bookingErrorMessage.value = result.conflicts?.join('; ') || null;
      } else {
        bookingSuccessMessage.value = `Successfully booked ${createdCount} day(s)!`;
      }

      // Reset multi-day fields
      multiDayBooking.value = false;
      additionalDates.value = '';
    } else {
      await createBooking(deskId, selectedDate.value, onBehalf, guest);
      bookingSuccessMessage.value = 'Desk booked successfully!';
    }

    // Reset booking type fields
    if (bookingType.value === 'colleague') {
      colleagueId.value = '';
      colleagueName.value = '';
      bookingType.value = 'self';
    } else if (bookingType.value === 'guest') {
      guestName.value = '';
      guestEmail.value = '';
      bookingType.value = 'self';
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

  // Fetch area and room names for breadcrumbs
  try {
    const areasResp = await fetchAreas();
    for (const area of areasResp.data) {
      const roomsResp = await fetchRooms(area.id);
      const room = roomsResp.data.find(r => r.id === roomId);
      if (room) {
        areaName.value = area.attributes.name;
        roomName.value = room.attributes.name;
        break;
      }
    }
  } catch {
    // Ignore - breadcrumbs will just show generic names
  }

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

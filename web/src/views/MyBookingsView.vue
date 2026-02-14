<template>
  <div class="page-container">
    <PageHeader
      title="My Bookings"
      subtitle="View and manage your upcoming reservations"
      :breadcrumbs="[{ text: 'Home', to: '/' }, { text: 'My Bookings' }]"
    />

    <!-- Success/Error Messages -->
    <v-alert
      v-if="cancelSuccessMessage"
      type="success"
      class="mb-4"
      closable
      data-cy="cancel-success"
      @click:close="cancelSuccessMessage = null"
    >
      {{ cancelSuccessMessage }}
    </v-alert>
    <v-alert
      v-if="cancelErrorMessage"
      type="error"
      class="mb-4"
      closable
      data-cy="cancel-error"
      @click:close="cancelErrorMessage = null"
    >
      {{ cancelErrorMessage }}
    </v-alert>

    <!-- Loading State -->
    <LoadingState v-if="bookingsLoading" type="cards" :count="3" data-cy="bookings-loading" />

    <!-- Error State -->
    <v-alert v-else-if="bookingsError" type="error" class="mb-4" data-cy="bookings-error">
      Unable to load bookings. Please try again later.
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!bookings.length"
      title="No upcoming bookings"
      message="You don't have any reservations scheduled. Browse available items to make a booking."
      icon="$calendar"
      action-text="Find an Item"
      action-to="/"
      data-cy="bookings-empty"
    />

    <!-- Bookings Grid -->
    <div v-else class="card-grid" data-cy="bookings-list">
      <BookingCard
        v-for="booking in bookings"
        :key="booking.id"
        :booking="booking"
        :show-cancel="true"
        :cancelling="cancellingBookingId === booking.id"
        data-cy="booking-item"
        @cancel="handleCancelBooking"
        @note-updated="handleNoteUpdated"
      />
    </div>

    <!-- Confirm Cancel Dialog -->
    <ConfirmDialog
      v-model="showCancelDialog"
      title="Cancel Booking"
      message="Are you sure you want to cancel this booking? This action cannot be undone."
      confirm-text="Cancel Booking"
      confirm-color="error"
      @confirm="confirmCancelBooking"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ApiError } from '../api/client';
import { cancelBooking, fetchMyBookings, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, BookingCard, ConfirmDialog } from '../components';

const authStore = useAuthStore();
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const cancelSuccessMessage = ref<string | null>(null);
const cancelErrorMessage = ref<string | null>(null);
const cancellingBookingId = ref<string | null>(null);
const showCancelDialog = ref(false);
const pendingCancelId = ref<string | null>(null);
const { loading: bookingsLoading, error: bookingsError, run: runBookings } = useApi();
const { handleAuthError } = useAuthErrorHandler();

const loadBookings = async () => {
  try {
    const resp = await runBookings(() => fetchMyBookings());
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
  }
};

const handleNoteUpdated = (bookingId: string, note: string) => {
  const booking = bookings.value.find(b => b.id === bookingId);
  if (booking) {
    booking.attributes.note = note;
  }
};

const handleCancelBooking = (bookingId: string) => {
  pendingCancelId.value = bookingId;
  showCancelDialog.value = true;
};

const confirmCancelBooking = async () => {
  if (!pendingCancelId.value) return;

  const bookingId = pendingCancelId.value;
  cancelSuccessMessage.value = null;
  cancelErrorMessage.value = null;
  cancellingBookingId.value = bookingId;

  try {
    await cancelBooking(bookingId);
    cancelSuccessMessage.value = 'Booking cancelled successfully.';
    await loadBookings();
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      cancelErrorMessage.value = 'Booking not found or already cancelled.';
    } else {
      cancelErrorMessage.value = 'Unable to cancel booking. Please try again.';
    }
  } finally {
    cancellingBookingId.value = null;
    pendingCancelId.value = null;
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

  await loadBookings();
});
</script>

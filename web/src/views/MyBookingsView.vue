<template>
  <div class="page-container">
    <PageHeader
      :title="$t('bookings.title')"
      :subtitle="$t('bookings.subtitle')"
      :breadcrumbs="[{ text: $t('common.home'), to: '/' }, { text: $t('bookings.title') }]"
    />

    <!-- Error Messages -->
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
      {{ bookingsError }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!bookings.length"
      :title="$t('bookings.emptyTitle')"
      :message="$t('bookings.emptyMessage')"
      icon="$calendar"
      :action-text="$t('bookings.findAnItem')"
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
      :title="$t('bookings.cancelTitle')"
      :message="$t('bookings.cancelMessage')"
      :confirm-text="$t('bookings.cancelButton')"
      confirm-color="error"
      @confirm="confirmCancelBooking"
    />

    <v-snackbar v-model="showCancelSuccess" :timeout="3000" location="bottom" color="success" data-cy="cancel-success">
      {{ cancelSuccessMessage }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { cancelBooking, fetchMyBookings, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, BookingCard, ConfirmDialog } from '../components';

const { t } = useI18n();
const authStore = useAuthStore();
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const cancelSuccessMessage = ref<string | null>(null);
const showCancelSuccess = computed({
  get: () => cancelSuccessMessage.value !== null,
  set: (v: boolean) => { if (!v) cancelSuccessMessage.value = null; }
});
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
    cancelSuccessMessage.value = t('bookings.cancelledSuccessfully');
    await loadBookings();
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      cancelErrorMessage.value = t('bookings.notFoundOrCancelled');
    } else {
      cancelErrorMessage.value = t('bookings.unableToCancel');
    }
  } finally {
    showCancelDialog.value = false;
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
    if (isConnectionError(err)) {
      bookingsError.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    throw err;
  }

  await loadBookings();
});
</script>

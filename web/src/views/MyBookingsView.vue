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

    <!-- View switch: Tiles / Table -->
    <div class="d-flex justify-end mb-4" data-cy="view-switch-container">
      <div class="d-flex align-center">
        <span
          class="text-button mr-1"
          :class="activeView === 'cards' ? 'text-primary font-weight-bold' : 'text-medium-emphasis'"
        >{{ $t('itemGroups.viewTiles') }}</span>
        <v-tooltip v-if="isCompactViewport" location="top">
          <template #activator="{ props: tooltipProps }">
            <div v-bind="tooltipProps" data-cy="view-switch-disabled-wrapper">
              <v-switch
                :model-value="activeView === 'table'"
                :disabled="true"
                hide-details
                inline
                inset
                density="compact"
                color="primary"
                base-color="primary"
                data-cy="view-switch"
                class="view-switch"
              />
            </div>
          </template>
          <span data-cy="view-switch-tooltip">{{ $t('bookings.viewTableDesktopOnly') }}</span>
        </v-tooltip>
        <v-switch
          v-else
          :model-value="activeView === 'table'"
          :disabled="false"
          hide-details
          inline
          inset
          density="compact"
          color="primary"
          base-color="primary"
          data-cy="view-switch"
          class="view-switch"
          @update:model-value="toggleView"
        />
        <span
          class="text-button ml-1"
          :class="activeView === 'table' ? 'text-primary font-weight-bold' : 'text-medium-emphasis'"
        >{{ $t('itemGroups.viewTable') }}</span>
      </div>
    </div>

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

    <!-- Bookings Grid (tiles) -->
    <div v-else-if="activeView === 'cards'" class="card-grid" data-cy="bookings-list">
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

    <!-- Bookings Table -->
    <v-data-table
      v-else
      :headers="tableHeaders"
      :items="tableItems"
      item-value="id"
      density="comfortable"
      class="elevation-1"
      data-cy="bookings-table"
    >
      <template #[`item.status`]="{ item }">
        <StatusChip
          v-if="item.status"
          :status="item.status"
          size="x-small"
          :data-cy="`status-chip-${item.id}`"
        />
      </template>
      <template #[`item.onBehalf`]="{ item }">
        <span v-if="item.guestName" data-cy="table-guest">
          {{ $t('bookings.guest', { name: item.guestName }) }}
        </span>
        <span v-else-if="item.forUserName" data-cy="table-on-behalf-of">
          {{ $t('bookings.onBehalfOf', { name: item.forUserName }) }}
        </span>
        <span v-else-if="item.bookedForMe && item.bookedByUserName" data-cy="table-booked-by">
          {{ $t('bookings.bookedBy', { name: item.bookedByUserName }) }}
        </span>
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn
          color="error"
          variant="tonal"
          size="small"
          :loading="cancellingBookingId === item.id"
          :disabled="cancellingBookingId === item.id"
          data-cy="cancel-btn"
          @click="handleCancelBooking(item.id)"
        >
          {{ $t('bookings.cancelBooking') }}
        </v-btn>
      </template>
    </v-data-table>

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
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { cancelBooking, fetchMyBookings, type MyBookingAttributes } from '../api/bookings';
import { fetchMe } from '../api/me';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useMyBookingsViewPreference } from '../composables/useMyBookingsViewPreference';
import { deriveBookingStatus } from '../utils/bookingStatus';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, BookingCard, ConfirmDialog, StatusChip } from '../components';

const { t, locale } = useI18n();
const authStore = useAuthStore();
const bookings = ref<JsonApiResource<MyBookingAttributes>[]>([]);
const { activeView, load: loadViewPref, save: saveViewPref } = useMyBookingsViewPreference();
const isCompactViewport = ref(false);

const updateViewport = () => {
  if (typeof window.matchMedia !== 'function') {
    isCompactViewport.value = false;
    return;
  }
  isCompactViewport.value = window.matchMedia('(max-width: 768px)').matches;
};

const handleResize = () => {
  updateViewport();
};

const toggleView = (val: boolean | null) => {
  saveViewPref(val ? 'table' : 'cards');
};

const formatBookingDate = (bookingDate: string): string => {
  const date = new Date(bookingDate + 'T00:00:00');
  return date.toLocaleDateString(locale.value || undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

const tableHeaders = computed(() => [
  { title: t('bookings.colDate'), key: 'date', sortable: true },
  { title: t('bookings.colItem'), key: 'itemName', sortable: true },
  { title: t('bookings.colArea'), key: 'area', sortable: true },
  { title: t('bookings.colStatus'), key: 'status', sortable: false },
  { title: t('bookings.colOnBehalf'), key: 'onBehalf', sortable: false },
  { title: t('bookings.colActions'), key: 'actions', sortable: false, align: 'end' as const }
]);

const tableItems = computed(() =>
  bookings.value.map((b) => ({
    id: b.id,
    date: formatBookingDate(b.attributes.booking_date),
    itemName: b.attributes.item_name,
    area: `${b.attributes.item_group_name} · ${b.attributes.area_name}`,
    status: deriveBookingStatus(b.attributes),
    forUserName: b.attributes.for_user_name ?? '',
    bookedForMe: b.attributes.booked_for_me,
    bookedByUserName: b.attributes.booked_by_user_name,
    guestName: b.attributes.guest_name ?? ''
  }))
);
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
  updateViewport();
  loadViewPref(!isCompactViewport.value);
  window.addEventListener('resize', handleResize);

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

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

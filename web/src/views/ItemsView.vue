<template>
  <div class="page-container">
    <PageHeader
      title="Items"
      subtitle="Select an item to book for your chosen date"
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <v-btn
          v-if="activeItemGroupId"
          variant="text"
          size="small"
          :to="{
            name: 'item-group-bookings',
            params: { itemGroupId: activeItemGroupId! },
            query: breadcrumbAreaId ? { areaId: breadcrumbAreaId } : {}
          }"
          data-cy="view-item-group-bookings"
        >
          View Item Group Bookings
        </v-btn>
      </template>
    </PageHeader>

    <!-- Date Selection & Booking Options -->
    <v-card class="mb-6">
      <v-card-text>
        <!-- Booking Mode Toggle -->
        <v-btn-toggle
          v-model="bookingMode"
          mandatory
          density="compact"
          class="mb-4"
          data-cy="booking-mode-toggle"
        >
          <v-btn value="day" data-cy="mode-day-btn">Day</v-btn>
          <v-btn value="week" data-cy="mode-week-btn">Week</v-btn>
        </v-btn-toggle>

        <div class="d-flex flex-wrap align-end ga-4 mb-4">
          <!-- Day mode: date picker -->
          <DatePickerField
            v-if="bookingMode === 'day'"
            v-model="selectedDate"
            label="Booking Date"
            :min="todayDate"
            data-cy="items-date"
            style="max-width: 280px;"
          />

          <!-- Week mode: week selector -->
          <v-select
            v-if="bookingMode === 'week'"
            v-model="selectedWeek"
            :items="weekOptions"
            item-title="label"
            item-value="value"
            label="Calendar Week"
            density="compact"
            hide-details
            data-cy="week-selector"
            style="max-width: 320px;"
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
          v-if="bookingMode === 'day'"
          v-model="multiDayBooking"
          label="Book multiple days"
          density="compact"
          hide-details
          class="mt-2"
          data-cy="multi-day-checkbox"
        />
        <v-expand-transition>
          <div v-if="bookingMode === 'day' && multiDayBooking" class="mt-2">
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
      @click:close="closeSuccessMessage"
    >
      {{ bookingSuccessMessage }}
      <v-btn
        v-if="lastBookingId"
        variant="text"
        size="small"
        class="ml-2"
        data-cy="add-note-after-booking"
        @click="openPostBookingNoteDialog"
      >
        Add note
      </v-btn>
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
    <LoadingState v-if="itemsLoading || weekDataLoading" type="cards" :count="6" data-cy="items-loading" />

    <!-- Error State -->
    <v-alert v-else-if="itemsErrorMessage" type="error" class="mb-4" data-cy="items-error">
      {{ itemsErrorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="bookingMode === 'day' && !items.length"
      title="No items available"
      message="This item group doesn't have any items configured yet."
      icon="$desk"
      data-cy="items-empty"
    />
    <EmptyState
      v-else-if="bookingMode === 'week' && !weekItems.length"
      title="No items available"
      message="This item group doesn't have any items configured yet."
      icon="$desk"
      data-cy="items-empty"
    />

    <!-- Items Grid (Day mode) -->
    <div v-else-if="bookingMode === 'day'" class="card-grid" data-cy="items-list">
      <v-card
        v-for="entry in items"
        :key="entry.id"
        :class="[
          'item-card',
          { 'item-available': entry.attributes.availability === 'available' },
          { 'item-occupied': entry.attributes.availability === 'occupied' }
        ]"
        data-cy="item-entry"
        :data-cy-item-id="entry.id"
        :data-cy-availability="entry.attributes.availability"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar
              :color="entry.attributes.availability === 'available' ? 'success' : 'warning'"
              variant="tonal"
              size="48"
            >
              <v-icon size="24">$desk</v-icon>
            </v-avatar>
          </template>
          <v-card-title class="d-flex align-center">
            {{ entry.attributes.name }}
            <StatusChip
              :status="entry.attributes.availability === 'available' ? 'available' : 'booked'"
              size="x-small"
              class="ml-2"
              data-cy="item-status"
            />
          </v-card-title>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Equipment -->
          <div v-if="entry.attributes.equipment?.length" class="mb-2" data-cy="item-equipment">
            <div class="text-caption text-medium-emphasis mb-1">Equipment</div>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="equip in entry.attributes.equipment"
                :key="equip"
                size="x-small"
                variant="outlined"
              >
                {{ equip }}
              </v-chip>
            </div>
          </div>

          <!-- Warning -->
          <v-alert
            v-if="entry.attributes.warning"
            type="warning"
            variant="tonal"
            density="compact"
            class="mt-2"
            data-cy="item-warning"
          >
            {{ entry.attributes.warning }}
          </v-alert>

          <!-- Booker name -->
          <div
            v-if="entry.attributes.availability === 'occupied' && entry.attributes.booker_name"
            class="text-caption text-medium-emphasis mt-2"
            data-cy="item-booker"
          >
            <v-icon size="14" class="mr-1">$user</v-icon>
            {{ entry.attributes.booker_name }}
          </div>

          <!-- Booking note -->
          <div
            v-if="entry.attributes.availability === 'occupied' && entry.attributes.note"
            class="d-flex align-center ga-1 mt-1 text-caption text-medium-emphasis"
            data-cy="item-note"
          >
            <v-icon size="14">mdi-text-box-outline</v-icon>
            <span :ref="setNoteRef(entry.id)" class="note-text">{{ entry.attributes.note }}</span>
            <v-btn
              v-if="noteTruncatedMap[entry.id]"
              icon
              size="x-small"
              variant="text"
              data-cy="item-note-expand"
              @click="expandedNote = entry.attributes.note"
            >
              <v-icon size="14">mdi-arrow-expand</v-icon>
            </v-btn>
          </div>
        </v-card-text>

        <v-card-actions class="px-4 pb-4">
          <v-btn
            v-if="entry.attributes.availability === 'available'"
            color="primary"
            variant="flat"
            block
            :loading="bookingItemId === entry.id"
            :disabled="bookingItemId !== null || cancelingBookingId !== null"
            data-cy="book-item-btn"
            @click="bookItem(entry.id)"
          >
            Book This Item
          </v-btn>
          <v-btn
            v-else-if="authStore.isAdmin && entry.attributes.booking_id"
            color="error"
            variant="tonal"
            block
            :loading="cancelingBookingId === entry.attributes.booking_id"
            :disabled="bookingItemId !== null || cancelingBookingId !== null"
            data-cy="admin-cancel-btn"
            @click="adminCancelBooking(entry.attributes.booking_id!)"
          >
            Cancel Booking
          </v-btn>
          <div v-else class="text-center w-100 text-caption text-medium-emphasis py-2">
            Not available for {{ formattedDate }}
          </div>
        </v-card-actions>
      </v-card>
    </div>

    <!-- Items Grid (Week mode) -->
    <div v-else-if="bookingMode === 'week' && weekItems.length" class="card-grid" data-cy="week-items-list">
      <v-card
        v-for="item in weekItems"
        :key="item.id"
        class="item-card"
        data-cy="week-item-entry"
        :data-cy-item-name="item.name"
        :data-cy-item-id="item.id"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar color="primary" variant="tonal" size="48">
              <v-icon size="24">$desk</v-icon>
            </v-avatar>
          </template>
          <v-card-title>{{ item.name }}</v-card-title>
        </v-card-item>

        <v-card-text class="pt-0">
          <div :class="isMobile ? 'week-days-compact' : 'week-days'" data-cy="week-days">
            <div
              v-for="(date, dayIdx) in selectedWeekDates"
              :key="date"
              class="week-day-slot"
              :data-cy-weekday="getWeekdayLabel(dayIdx)"
            >
              <span class="week-day-label text-caption font-weight-medium">
                {{ getWeekdayLabel(dayIdx, isMobile) }}
              </span>
              <v-checkbox
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :model-value="isWeekDaySelected(item.id, date)"
                hide-details
                density="compact"
                color="success"
                class="week-day-checkbox"
                data-cy="week-day-checkbox"
                @update:model-value="toggleWeekDay(item.id, date)"
              />
              <v-checkbox
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :model-value="true"
                hide-details
                density="compact"
                color="primary"
                disabled
                class="week-day-checkbox"
                data-cy="week-day-mine"
              />
              <v-checkbox
                v-else
                :model-value="true"
                hide-details
                density="compact"
                disabled
                class="week-day-checkbox"
                :data-cy="getWeekDayStatus(item.id, date) === 'unavailable' ? 'week-day-unavailable' : 'week-day-other'"
              />
              <span
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                class="week-day-status text-caption text-success"
              >free</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                class="week-day-status text-caption text-primary"
              >{{ authStore.userName || 'Me' }}</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'"
                class="week-day-status text-caption text-error"
              >{{ getWeekDayBooker(item.id, date) }}</span>
              <span
                v-else
                class="week-day-status text-caption text-medium-emphasis"
              >n/a</span>
            </div>
          </div>
        </v-card-text>
      </v-card>
    </div>

    <!-- Confirm Booking Button (Week mode) -->
    <div v-if="bookingMode === 'week' && weekSelectionCount > 0" class="mt-4" data-cy="week-confirm-section">
      <v-btn
        color="primary"
        variant="flat"
        size="large"
        block
        :loading="weekBookingInProgress"
        data-cy="week-confirm-btn"
        @click="submitWeekBookings"
      >
        Confirm My Booking ({{ weekSelectionCount }} {{ weekSelectionCount === 1 ? 'day' : 'days' }})
      </v-btn>
    </div>

    <!-- Week Booking Results -->
    <v-card v-if="weekBookingResults.length" class="mt-4" data-cy="week-booking-results">
      <v-card-title>Booking Results</v-card-title>
      <v-card-text>
        <div v-for="result in weekBookingResults" :key="result.date + result.itemName" class="d-flex align-center ga-2 mb-1">
          <v-icon :color="result.success ? 'success' : 'error'" size="18">
            {{ result.success ? 'mdi-check-circle' : 'mdi-close-circle' }}
          </v-icon>
          <span class="text-body-2">
            {{ result.itemName }} - {{ result.dayLabel }}:
            {{ result.success ? 'Booked' : result.error }}
          </span>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" data-cy="week-results-close" @click="weekBookingResults = []">Close</v-btn>
      </v-card-actions>
    </v-card>

    <!-- Add Note Dialog (after booking) -->
    <v-dialog v-model="showPostBookingNoteDialog" max-width="500">
      <v-card>
        <v-card-title>Add Note</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="noteText"
            label="Note"
            :counter="500"
            :maxlength="500"
            rows="3"
            auto-grow
            data-cy="post-booking-note-input"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showPostBookingNoteDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :loading="savingNote"
            data-cy="post-booking-note-save"
            @click="saveNoteAfterBooking"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note view dialog (desktop) -->
    <v-dialog v-if="!useBottomSheet" v-model="showItemNoteDialog" max-width="500">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="item-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showItemNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note view bottom sheet (mobile) -->
    <v-bottom-sheet v-else v-model="showItemNoteDialog">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="item-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showItemNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import type { ComponentPublicInstance } from 'vue';
import { useRoute } from 'vue-router';
import { ApiError } from '../api/client';
import {
  createBooking,
  createMultiDayBooking,
  cancelBooking,
  updateBookingNote,
  fetchMyBookings,
  type BookOnBehalfOptions,
  type GuestBookingOptions
} from '../api/bookings';
import { fetchItems } from '../api/items';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import type { ItemAttributes } from '../api/items';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useWeekSelector, getWeekdayLabel } from '../composables/useWeekSelector';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, StatusChip, DatePickerField } from '../components';

const authStore = useAuthStore();
const items = ref<JsonApiResource<ItemAttributes>[]>([]);
const itemsErrorMessage = ref<string | null>(null);
const bookingSuccessMessage = ref<string | null>(null);
const bookingErrorMessage = ref<string | null>(null);
const bookingItemId = ref<string | null>(null);
const cancelingBookingId = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const todayDate = formatDate(new Date());
const route = useRoute();
const { loading: itemsLoading, run: runItems } = useApi();
const activeItemGroupId = ref<string | null>(null);
const areaName = ref('');
const itemGroupName = ref('');
const bookingType = ref<'self' | 'colleague' | 'guest'>('self');
const colleagueId = ref('');
const colleagueName = ref('');
const guestName = ref('');
const guestEmail = ref('');
const multiDayBooking = ref(false);
const additionalDates = ref('');
const lastBookingId = ref<string | null>(null);
const showPostBookingNoteDialog = ref(false);
const noteText = ref('');
const savingNote = ref(false);
const expandedNote = ref('');
const noteTruncatedMap = ref<Record<string, boolean>>({});
const noteElements = new Map<string, HTMLElement>();
const isMobile = ref(false);
const useBottomSheet = computed(() => isMobile.value);
const showItemNoteDialog = computed({
  get: () => expandedNote.value !== '',
  set: (v: boolean) => { if (!v) expandedNote.value = ''; }
});

// Week booking mode
const bookingMode = ref<'day' | 'week'>(
  (localStorage.getItem('sithub_booking_mode') as 'day' | 'week') || 'day'
);
const { weekOptions, selectedWeek, selectedWeekDates } = useWeekSelector();

// Per-day data for week mode: map of date -> items array
const weekData = ref<Record<string, JsonApiResource<ItemAttributes>[]>>({});
const weekDataLoading = ref(false);
const myWeekBookings = ref<Set<string>>(new Set());

// Week day selections: Set of "itemId::date" keys
const weekSelections = ref<Set<string>>(new Set());
const weekBookingInProgress = ref(false);

interface WeekBookingResult {
  itemName: string;
  date: string;
  dayLabel: string;
  success: boolean;
  error?: string;
}
const weekBookingResults = ref<WeekBookingResult[]>([]);

// Unique items across all days in week mode
const weekItems = computed(() => {
  const itemsMap = new Map<string, string>();
  for (const dayItems of Object.values(weekData.value)) {
    for (const item of dayItems) {
      itemsMap.set(item.id, item.attributes.name);
    }
  }
  return Array.from(itemsMap.entries())
    .map(([id, name]) => ({ id, name }))
    .sort((a, b) => a.name.localeCompare(b.name));
});

const weekSelectionCount = computed(() => weekSelections.value.size);

const getWeekSelectionKey = (itemId: string, date: string) => `${itemId}::${date}`;

const getWeekDayStatus = (
  itemId: string,
  date: string
): 'free' | 'booked-by-me' | 'booked-by-other' | 'unavailable' => {
  const dayItems = weekData.value[date];
  if (!dayItems) return 'unavailable';
  const item = dayItems.find(i => i.id === itemId);
  if (!item) return 'unavailable';
  if (item.attributes.availability === 'available') return 'free';
  if (isBookedByMe(itemId, date)) return 'booked-by-me';
  return 'booked-by-other';
};

const getWeekDayBooker = (itemId: string, date: string): string => {
  const dayItems = weekData.value[date];
  if (!dayItems) return 'Booked';
  const item = dayItems.find(i => i.id === itemId);
  return item?.attributes.booker_name || 'Booked';
};

const isWeekDaySelected = (itemId: string, date: string) =>
  weekSelections.value.has(getWeekSelectionKey(itemId, date));

const isBookedByMe = (itemId: string, date: string) =>
  myWeekBookings.value.has(getWeekSelectionKey(itemId, date));

const toggleWeekDay = (itemId: string, date: string) => {
  if (getWeekDayStatus(itemId, date) !== 'free') return;
  const key = getWeekSelectionKey(itemId, date);
  const next = new Set(weekSelections.value);
  if (next.has(key)) {
    next.delete(key);
  } else {
    next.add(key);
  }
  weekSelections.value = next;
};

const loadWeekData = async (itemGroupId: string, keepResults = false) => {
  weekDataLoading.value = true;
  weekSelections.value = new Set();
  itemsErrorMessage.value = null;
  if (!keepResults) weekBookingResults.value = [];
  try {
    const dates = selectedWeekDates.value;
    const results = await Promise.all(
      dates.map(date => fetchItems(itemGroupId, date).then(resp => ({ date, items: resp.data })))
    );
    const data: Record<string, JsonApiResource<ItemAttributes>[]> = {};
    for (const { date, items: dayItems } of results) {
      data[date] = dayItems;
    }
    weekData.value = data;

    const bookingsResp = await fetchMyBookings().catch(() => ({ data: [] }));
    const bookedSet = new Set<string>();
    for (const booking of bookingsResp.data) {
      const bookingDate = booking.attributes.booking_date;
      if (dates.includes(bookingDate)) {
        bookedSet.add(getWeekSelectionKey(booking.attributes.item_id, bookingDate));
      }
    }
    myWeekBookings.value = bookedSet;
  } catch {
    weekData.value = {};
    myWeekBookings.value = new Set();
    itemsErrorMessage.value = 'Unable to load weekly items.';
  } finally {
    weekDataLoading.value = false;
  }
};

const submitWeekBookings = async () => {
  if (!activeItemGroupId.value || weekSelections.value.size === 0) return;

  bookingErrorMessage.value = null;
  if (bookingType.value === 'colleague') {
    if (!colleagueId.value.trim() || !colleagueName.value.trim()) {
      bookingErrorMessage.value = 'Please enter both colleague ID and name.';
      return;
    }
  }
  if (bookingType.value === 'guest') {
    if (!guestName.value.trim()) {
      bookingErrorMessage.value = 'Please enter the guest name.';
      return;
    }
  }

  weekBookingInProgress.value = true;
  weekBookingResults.value = [];

  const entries = Array.from(weekSelections.value).map(key => {
    const sep = key.indexOf('::');
    const itemId = key.substring(0, sep);
    const date = key.substring(sep + 2);
    const itemName = weekItems.value.find(item => item.id === itemId)?.name || 'Item';
    return { itemId, itemName, date };
  });

  const onBehalf: BookOnBehalfOptions | undefined =
    bookingType.value === 'colleague'
      ? { forUserId: colleagueId.value.trim(), forUserName: colleagueName.value.trim() }
      : undefined;

  const guest: GuestBookingOptions | undefined =
    bookingType.value === 'guest'
      ? { guestName: guestName.value.trim(), guestEmail: guestEmail.value.trim() || undefined }
      : undefined;

  const promises = entries.map(async ({ itemId, itemName, date }) => {
    try {
      await createBooking(itemId, date, onBehalf, guest);
      return { itemName, date, success: true };
    } catch (err) {
      const msg = err instanceof ApiError && err.detail ? err.detail : 'Booking failed';
      return { itemName, date, success: false, error: msg };
    }
  });

  const results = await Promise.allSettled(promises);
  const dayLabels = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];

  weekBookingResults.value = results.map(r => {
    const val = r.status === 'fulfilled'
      ? r.value
      : { itemName: '', date: '', success: false as const, error: 'Unexpected error' };
    const dayIdx = selectedWeekDates.value.indexOf(val.date);
    const dayLabel = dayIdx >= 0 ? (dayLabels[dayIdx] ?? val.date) : val.date;
    return { itemName: val.itemName, date: val.date, dayLabel, success: val.success, error: val.error };
  });

  weekSelections.value = new Set();

  // Refresh week data (keep results visible)
  if (activeItemGroupId.value) {
    await loadWeekData(activeItemGroupId.value, true);
  }

  weekBookingInProgress.value = false;
};

const queryAreaId = computed(() => {
  const value = route.query.areaId;
  return typeof value === 'string' ? value : undefined;
});
const resolvedAreaId = ref<string | null>(null);
const breadcrumbAreaId = computed(() =>
  resolvedAreaId.value ? resolvedAreaId.value : areaName.value ? undefined : queryAreaId.value
);

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  {
    text: areaName.value || 'Area',
    to: breadcrumbAreaId.value ? `/areas/${breadcrumbAreaId.value}/item-groups` : undefined
  },
  { text: itemGroupName.value || 'Item Group' }
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

const loadItems = async (itemGroupId: string, date: string) => {
  itemsErrorMessage.value = null;
  try {
    const normalizedDate = ensureDate(date);
    const resp = await runItems(() => fetchItems(itemGroupId, normalizedDate));
    items.value = resp.data;
    await nextTick();
    updateNoteTruncation();
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      itemsErrorMessage.value = 'Item group not found.';
      return;
    }
    itemsErrorMessage.value = 'Unable to load items.';
  }
};

const bookItem = async (itemId: string) => {
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

  bookingItemId.value = itemId;

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

      const result = await createMultiDayBooking(itemId, dates, onBehalf, guest);
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
      const result = await createBooking(itemId, selectedDate.value, onBehalf, guest);
      lastBookingId.value = result.data.id;
      bookingSuccessMessage.value = 'Item booked successfully!';
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

    // Reload items to reflect updated availability
    if (activeItemGroupId.value) {
      await loadItems(activeItemGroupId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 409) {
      // Use backend's detail message if available, otherwise a generic message
      const detail = err.detail || 'This item is no longer available for the selected date.';
      bookingErrorMessage.value = `${detail} Please choose another item.`;

      // Refresh item list so user sees updated availability
      if (activeItemGroupId.value) {
        await loadItems(activeItemGroupId.value, selectedDate.value);
      }
    } else if (err instanceof ApiError && err.status === 404) {
      bookingErrorMessage.value = 'Item not found.';
    } else {
      bookingErrorMessage.value = 'Unable to book item. Please try again.';
    }
  } finally {
    bookingItemId.value = null;
  }
};

const adminCancelBooking = async (bookingId: string) => {
  bookingSuccessMessage.value = null;
  bookingErrorMessage.value = null;
  cancelingBookingId.value = bookingId;

  try {
    await cancelBooking(bookingId);
    bookingSuccessMessage.value = 'Booking cancelled successfully.';

    // Reload items to reflect updated availability
    if (activeItemGroupId.value) {
      await loadItems(activeItemGroupId.value, selectedDate.value);
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

const closeSuccessMessage = () => {
  bookingSuccessMessage.value = null;
  lastBookingId.value = null;
};

const openPostBookingNoteDialog = () => {
  noteText.value = '';
  showPostBookingNoteDialog.value = true;
};

const saveNoteAfterBooking = async () => {
  if (!lastBookingId.value) return;
  savingNote.value = true;
  try {
    await updateBookingNote(lastBookingId.value, noteText.value);
    showPostBookingNoteDialog.value = false;
    bookingSuccessMessage.value = 'Booking created with note!';
    lastBookingId.value = null;
    if (activeItemGroupId.value) {
      await loadItems(activeItemGroupId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    bookingErrorMessage.value = 'Unable to save note. Please try again.';
  } finally {
    savingNote.value = false;
  }
};

const setNoteRef = (id: string) => (el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    noteElements.set(id, el);
    return;
  }
  if (el && '$el' in el && (el.$el instanceof HTMLElement)) {
    noteElements.set(id, el.$el);
    return;
  }
  noteElements.delete(id);
};

const updateNoteTruncation = () => {
  const map: Record<string, boolean> = {};
  for (const entry of items.value) {
    const el = noteElements.get(entry.id);
    if (el) {
      map[entry.id] = el.scrollWidth > el.clientWidth;
    }
  }
  noteTruncatedMap.value = map;
};

onMounted(async () => {
  updateViewport();
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

  const itemGroupId = route.params.itemGroupId;
  if (typeof itemGroupId !== 'string' || itemGroupId.trim() === '') {
    itemsErrorMessage.value = 'Item group not found.';
    return;
  }

  activeItemGroupId.value = itemGroupId;

  // Fetch area and item group names for breadcrumbs
  try {
    const areasResp = await fetchAreas();
    for (const area of areasResp.data) {
      const igResp = await fetchItemGroups(area.id);
      const ig = igResp.data.find(ig => ig.id === itemGroupId);
      if (ig) {
        areaName.value = area.attributes.name;
        itemGroupName.value = ig.attributes.name;
        resolvedAreaId.value = area.id;
        break;
      }
    }
  } catch {
    // Ignore - breadcrumbs will just show generic names
  }

  if (bookingMode.value === 'week') {
    await loadWeekData(itemGroupId);
  } else {
    await loadItems(itemGroupId, selectedDate.value);
  }
});

watch(
  selectedDate,
  async (value) => {
    if (!activeItemGroupId.value || bookingMode.value !== 'day') {
      return;
    }
    await loadItems(activeItemGroupId.value, value);
  },
  { flush: 'post' }
);

watch(bookingMode, async (mode) => {
  localStorage.setItem('sithub_booking_mode', mode);
  if (!activeItemGroupId.value) return;
  if (mode === 'week') {
    multiDayBooking.value = false;
    additionalDates.value = '';
    await loadWeekData(activeItemGroupId.value);
  } else {
    weekData.value = {};
    weekSelections.value = new Set();
    weekBookingResults.value = [];
    await loadItems(activeItemGroupId.value, selectedDate.value);
  }
});

watch(selectedWeek, async () => {
  if (!activeItemGroupId.value || bookingMode.value !== 'week') return;
  await loadWeekData(activeItemGroupId.value);
});

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

function updateViewport() {
  if (typeof window.matchMedia === 'function') {
    isMobile.value = window.matchMedia('(max-width: 600px)').matches;
    return;
  }
  isMobile.value = false;
}

function handleResize() {
  updateViewport();
  updateNoteTruncation();
}

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

<style scoped>
.note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.week-days {
  display: flex;
  gap: 8px;
  justify-content: space-between;
}

.week-days-compact {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 4px;
}

.week-day-slot {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 44px;
  padding: 4px;
}

.week-day-label {
  margin-bottom: 2px;
}

.week-day-checkbox {
  min-height: 44px;
}

.week-day-status {
  font-size: 0.7rem;
  line-height: 1.2;
  text-align: center;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>

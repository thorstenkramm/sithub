<template>
  <div class="page-container">
    <PageHeader
      title=""
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
        </v-radio-group>

        <!-- Colleague Fields -->
        <v-expand-transition>
          <div v-if="bookingType === 'colleague'" class="mt-4">
            <v-autocomplete
              v-model="selectedColleagueId"
              :items="usersList"
              item-title="displayName"
              item-value="id"
              label="Select colleague"
              density="compact"
              :loading="usersLoading"
              clearable
              data-cy="colleague-select"
              style="max-width: 360px;"
            />
          </div>
        </v-expand-transition>

      </v-card-text>
    </v-card>

    <!-- Success/Error Messages -->
    <v-alert
      v-if="bookingSuccessMessage || bookingSuccessDetails"
      type="success"
      class="mb-4"
      closable
      data-cy="booking-success"
      @click:close="closeSuccessMessage"
    >
      <div class="d-flex align-center ga-2">
        <v-icon color="success" size="18">mdi-check-circle</v-icon>
        <span class="text-body-2" data-cy="booking-success-text">
          {{ bookingSuccessDetails
            ? `${bookingSuccessDetails.itemName} - ${formatDisplayDate(bookingSuccessDetails.date)}`
            : bookingSuccessMessage
          }}
        </span>
      </div>
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
      v-if="bookingErrorMessage || bookingErrorDetails"
      type="error"
      class="mb-4"
      closable
      data-cy="booking-error"
      @click:close="clearErrorMessage"
    >
      <div class="d-flex align-center ga-2">
        <v-icon color="error" size="18">mdi-alert-circle</v-icon>
        <span class="text-body-2" data-cy="booking-error-text">
          {{ bookingErrorDetails
            ? `${bookingErrorDetails.itemName} - ${formatDisplayDate(bookingErrorDetails.date)}: ${bookingErrorDetails.error}`
            : bookingErrorMessage
          }}
        </span>
      </div>
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
          <template #append>
            <div v-if="entry.attributes.availability === 'occupied'" class="d-flex align-center">
              <v-tooltip
                v-if="!expandedDayTiles.has(entry.id) && entry.attributes.warning"
                location="top"
              >
                <template #activator="{ props: tooltipProps }">
                  <v-btn
                    v-bind="tooltipProps"
                    icon
                    variant="text"
                    size="x-small"
                    color="warning"
                    class="mr-1"
                    :aria-label="`View warning for ${entry.attributes.name}`"
                    data-cy="folded-warning-icon"
                  >
                    <v-icon size="18">mdi-alert</v-icon>
                  </v-btn>
                </template>
                {{ entry.attributes.warning }}
              </v-tooltip>
              <v-btn
                icon
                variant="text"
                size="small"
                data-cy="day-tile-chevron"
                :aria-label="`Toggle details for ${entry.attributes.name}`"
                :aria-expanded="expandedDayTiles.has(entry.id)"
                @click="toggleDayTileExpansion(entry.id)"
              >
                <v-icon>
                  {{ expandedDayTiles.has(entry.id) ? 'mdi-chevron-down' : 'mdi-chevron-left' }}
                </v-icon>
              </v-btn>
            </div>
          </template>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Equipment (hidden on folded booked tiles) -->
          <div
            v-if="entry.attributes.equipment?.length
              && (entry.attributes.availability === 'available' || expandedDayTiles.has(entry.id))"
            class="mb-2"
            data-cy="item-equipment"
          >
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

          <!-- Warning (hidden on folded booked tiles) -->
          <v-alert
            v-if="entry.attributes.warning
              && (entry.attributes.availability === 'available' || expandedDayTiles.has(entry.id))"
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
            class="text-body-2 text-medium-emphasis mt-2"
            data-cy="item-booker"
          >
            <v-icon size="14" class="mr-1">$user</v-icon>
            {{ entry.attributes.booker_name }}
          </div>

          <!-- Booking note -->
          <div
            v-if="entry.attributes.availability === 'occupied' && entry.attributes.note"
            class="d-flex align-center ga-1 mt-1 text-body-2 text-medium-emphasis"
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
            Book
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
          <div v-else class="py-2" />
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
          <template #append>
            <div class="d-flex align-center">
              <v-tooltip
                v-if="!expandedWeekTiles.has(item.id) && getWeekItemAttributes(item.id).warning"
                location="top"
              >
                <template #activator="{ props: tooltipProps }">
                  <v-btn
                    v-bind="tooltipProps"
                    icon
                    variant="text"
                    size="x-small"
                    color="warning"
                    class="mr-1"
                    :aria-label="`View warning for ${item.name}`"
                    data-cy="week-folded-warning-icon"
                  >
                    <v-icon size="18">mdi-alert</v-icon>
                  </v-btn>
                </template>
                {{ getWeekItemAttributes(item.id).warning }}
              </v-tooltip>
              <v-btn
                icon
                variant="text"
                size="small"
                data-cy="week-tile-chevron"
                :aria-label="`Toggle details for ${item.name}`"
                :aria-expanded="expandedWeekTiles.has(item.id)"
                @click="toggleWeekTileExpansion(item.id)"
              >
                <v-icon>
                  {{ expandedWeekTiles.has(item.id) ? 'mdi-chevron-down' : 'mdi-chevron-left' }}
                </v-icon>
              </v-btn>
            </div>
          </template>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Folded view: compact M-F row -->
          <div
            v-if="!expandedWeekTiles.has(item.id)"
            :class="isMobile ? 'week-days-compact' : 'week-days'"
            :style="isMobile ? { gridTemplateColumns: `repeat(${selectedWeekDates.length}, 1fr)` } : undefined"
            data-cy="week-days"
          >
            <div
              v-for="(date, dayIdx) in selectedWeekDates"
              :key="date"
              :class="['week-day-slot', { 'week-day-past': isDateInPast(date) }]"
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
                :disabled="isDateInPast(date)"
                class="week-day-checkbox"
                :data-cy="isDateInPast(date) ? 'week-day-checkbox-past' : 'week-day-checkbox'"
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
                :data-cy="getWeekDayStatus(item.id, date) === 'unavailable'
                  ? 'week-day-unavailable' : 'week-day-other'"
              />
              <span
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :class="['week-day-status', 'text-caption', isDateInPast(date) ? 'text-medium-emphasis' : 'text-success']"
              >free</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :class="['week-day-status', 'text-caption', isDateInPast(date) ? 'text-medium-emphasis' : 'text-primary']"
              >{{ authStore.userName || 'Me' }}</span>
              <template v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'">
                <v-tooltip location="top" :disabled="!shouldShowWeekNameTooltip(getWeekDayBooker(item.id, date))">
                  <template #activator="{ props: tooltipProps }">
                    <span
                      v-bind="tooltipProps"
                      class="week-day-status week-day-status-truncated text-caption text-error"
                    >{{ getWeekDayBooker(item.id, date) }}</span>
                  </template>
                  {{ getWeekDayBooker(item.id, date) }}
                </v-tooltip>
              </template>
              <span
                v-else
                :class="['week-day-status', 'text-caption', 'text-medium-emphasis']"
              >n/a</span>
            </div>
          </div>
          <!-- Expanded view: one line per day -->
          <div v-else data-cy="week-days-expanded">
            <div
              v-for="(date, dayIdx) in selectedWeekDates"
              :key="date"
              :class="['week-day-expanded', { 'week-day-past': isDateInPast(date) }]"
              :data-cy-weekday="getWeekdayLabel(dayIdx)"
            >
              <span class="text-body-2 font-weight-medium week-day-expanded-label">
                {{ getFullDayLabel(date, dayIdx) }}
              </span>
              <v-checkbox
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :model-value="isWeekDaySelected(item.id, date)"
                hide-details
                density="compact"
                color="success"
                :disabled="isDateInPast(date)"
                class="week-day-checkbox"
                :data-cy="isDateInPast(date) ? 'week-day-checkbox-past' : 'week-day-checkbox'"
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
                :data-cy="getWeekDayStatus(item.id, date) === 'unavailable'
                  ? 'week-day-unavailable' : 'week-day-other'"
              />
              <span
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-success']"
              >free</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-primary']"
              >{{ authStore.userName || 'Me' }}</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-error']"
              >{{ getWeekDayBooker(item.id, date) }}</span>
              <span
                v-else
                :class="['text-body-2', 'text-medium-emphasis']"
              >n/a</span>
            </div>
            <!-- Equipment -->
            <div
              v-if="getWeekItemAttributes(item.id).equipment.length"
              class="mt-3"
              data-cy="week-item-equipment"
            >
              <div class="text-caption text-medium-emphasis mb-1">Equipment</div>
              <div class="d-flex flex-wrap ga-1">
                <v-chip
                  v-for="equip in getWeekItemAttributes(item.id).equipment"
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
              v-if="getWeekItemAttributes(item.id).warning"
              type="warning"
              variant="tonal"
              density="compact"
              class="mt-2"
              data-cy="week-item-warning"
            >
              {{ getWeekItemAttributes(item.id).warning }}
            </v-alert>
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
            {{ result.success ? 'mdi-check-circle' : 'mdi-alert-circle' }}
          </v-icon>
          <span class="text-body-2">
            {{ result.itemName }} - {{ result.dayLabel }}{{ result.success ? '' : ': ' + result.error }}
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
  cancelBooking,
  updateBookingNote,
  fetchMyBookings,
  type BookOnBehalfOptions
} from '../api/bookings';
import { fetchItems } from '../api/items';
import { fetchUsers } from '../api/users';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import type { ItemAttributes } from '../api/items';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useWeekSelector, getWeekdayLabel } from '../composables/useWeekSelector';
import { useWeekendPreference } from '../composables/useWeekendPreference';
import { getSafeLocalStorage } from '../composables/storage';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState, StatusChip, DatePickerField } from '../components';

const authStore = useAuthStore();
const items = ref<JsonApiResource<ItemAttributes>[]>([]);
const itemsErrorMessage = ref<string | null>(null);
const bookingSuccessMessage = ref<string | null>(null);
const bookingErrorMessage = ref<string | null>(null);
const bookingSuccessDetails = ref<{ itemName: string; date: string } | null>(null);
const bookingErrorDetails = ref<{ itemName: string; date: string; error: string } | null>(null);
const lastBookingDetails = ref<{ itemName: string; date: string } | null>(null);
const bookingItemId = ref<string | null>(null);
const cancelingBookingId = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const todayDate = formatDate(new Date());
const route = useRoute();
const { loading: itemsLoading, run: runItems } = useApi();
const activeItemGroupId = ref<string | null>(null);
const areaName = ref('');
const itemGroupName = ref('');
const bookingType = ref<'self' | 'colleague'>('self');
const selectedColleagueId = ref<string | null>(null);
const usersList = ref<Array<{ id: string; displayName: string }>>([]);
const usersLoading = ref(false);
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
const storage = getSafeLocalStorage();
const bookingMode = ref<'day' | 'week'>(
  (storage?.getItem('sithub_booking_mode') as 'day' | 'week') || 'day'
);
const { showWeekends } = useWeekendPreference();
const { weekOptions, selectedWeek, selectedWeekDates } = useWeekSelector(showWeekends);

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

// Story 15-1: Week tile expansion
const expandedWeekTiles = ref<Set<string>>(new Set());
const WEEKDAY_LONG_FORMATTER = new Intl.DateTimeFormat(undefined, { weekday: 'long' });
const WEEK_NAME_TRUNCATE_LIMIT = 12;

const getFullDayLabel = (date: string, fallbackIndex: number): string => {
  const parsed = new Date(`${date}T00:00:00`);
  if (!Number.isNaN(parsed.getTime())) {
    return WEEKDAY_LONG_FORMATTER.format(parsed);
  }
  const fallback = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
  return fallback[fallbackIndex] ?? date;
};

const shouldShowWeekNameTooltip = (name: string): boolean => name.length > WEEK_NAME_TRUNCATE_LIMIT;

const toggleWeekTileExpansion = (itemId: string) => {
  const next = new Set(expandedWeekTiles.value);
  if (next.has(itemId)) {
    next.delete(itemId);
  } else {
    next.add(itemId);
  }
  expandedWeekTiles.value = next;
};

const weekItemAttributesMap = computed(() => {
  const map = new Map<string, { equipment: string[]; warning?: string }>();
  for (const dayItems of Object.values(weekData.value)) {
    for (const item of dayItems) {
      if (!map.has(item.id)) {
        map.set(item.id, {
          equipment: item.attributes.equipment || [],
          warning: item.attributes.warning
        });
      }
    }
  }
  return map;
});

const getWeekItemAttributes = (itemId: string): { equipment: string[]; warning?: string } =>
  weekItemAttributesMap.value.get(itemId) ?? { equipment: [] };

// Story 15-2: Day tile expansion
const expandedDayTiles = ref<Set<string>>(new Set());

const toggleDayTileExpansion = (itemId: string) => {
  const next = new Set(expandedDayTiles.value);
  if (next.has(itemId)) {
    next.delete(itemId);
  } else {
    next.add(itemId);
  }
  expandedDayTiles.value = next;
};

// Story 15-4: Past date detection
const isDateInPast = (date: string): boolean => date < todayDate;

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
  if (isDateInPast(date)) return;
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
  expandedWeekTiles.value = new Set();
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

const loadUsers = async () => {
  usersLoading.value = true;
  try {
    const resp = await fetchUsers();
    usersList.value = resp.data.map(u => ({
      id: u.id,
      displayName: u.attributes.display_name
    }));
  } catch {
    // Silently fail — colleague dropdown will just be empty
  } finally {
    usersLoading.value = false;
  }
};

const resolveColleagueName = (userId: string): string | undefined => {
  const user = usersList.value.find(u => u.id === userId);
  return user?.displayName;
};

const submitWeekBookings = async () => {
  if (!activeItemGroupId.value || weekSelections.value.size === 0) return;

  bookingErrorMessage.value = null;
  bookingErrorDetails.value = null;
  if (bookingType.value === 'colleague') {
    if (!selectedColleagueId.value) {
      bookingErrorMessage.value = 'Please select a colleague.';
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
    bookingType.value === 'colleague' && selectedColleagueId.value
      ? { forUserId: selectedColleagueId.value, forUserName: resolveColleagueName(selectedColleagueId.value) }
      : undefined;

  const promises = entries.map(async ({ itemId, itemName, date }) => {
    try {
      await createBooking(itemId, date, onBehalf);
      return { itemName, date, success: true };
    } catch (err) {
      const msg = err instanceof ApiError && err.detail ? err.detail : 'Booking failed';
      return { itemName, date, success: false, error: msg };
    }
  });

  const results = await Promise.allSettled(promises);

  weekBookingResults.value = results.map(r => {
    const val = r.status === 'fulfilled'
      ? r.value
      : { itemName: '', date: '', success: false as const, error: 'Unexpected error' };
    const dayIdx = selectedWeekDates.value.indexOf(val.date);
    const dayLabel = dayIdx >= 0 ? getFullDayLabel(val.date, dayIdx) : val.date;
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
  expandedDayTiles.value = new Set();
  expandedDayTiles.value = new Set();
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
  bookingSuccessDetails.value = null;
  bookingErrorMessage.value = null;
  bookingErrorDetails.value = null;
  lastBookingDetails.value = null;

  // Validate colleague selection
  if (bookingType.value === 'colleague') {
    if (!selectedColleagueId.value) {
      bookingErrorMessage.value = 'Please select a colleague.';
      return;
    }
  }

  bookingItemId.value = itemId;

  try {
    const onBehalf: BookOnBehalfOptions | undefined =
      bookingType.value === 'colleague' && selectedColleagueId.value
        ? { forUserId: selectedColleagueId.value, forUserName: resolveColleagueName(selectedColleagueId.value) }
        : undefined;

    const result = await createBooking(itemId, selectedDate.value, onBehalf);
    lastBookingId.value = result.data.id;

    const itemName = items.value.find(entry => entry.id === itemId)?.attributes.name || 'Item';
    const details = { itemName, date: selectedDate.value };
    bookingSuccessDetails.value = details;
    lastBookingDetails.value = details;

    // Reset booking type fields
    if (bookingType.value === 'colleague') {
      selectedColleagueId.value = null;
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

    const itemName = items.value.find(entry => entry.id === itemId)?.attributes.name || 'Item';
    let detail = 'Unable to book item. Please try again.';

    if (err instanceof ApiError && err.status === 409) {
      detail = err.detail || 'This item is no longer available for the selected date.';

      // Refresh item list so user sees updated availability
      if (activeItemGroupId.value) {
        await loadItems(activeItemGroupId.value, selectedDate.value);
      }
    } else if (err instanceof ApiError && err.status === 404) {
      detail = 'Item not found.';
    }

    bookingErrorDetails.value = { itemName, date: selectedDate.value, error: detail };
  } finally {
    bookingItemId.value = null;
  }
};

const adminCancelBooking = async (bookingId: string) => {
  bookingSuccessMessage.value = null;
  bookingSuccessDetails.value = null;
  bookingErrorMessage.value = null;
  bookingErrorDetails.value = null;
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
  bookingSuccessDetails.value = null;
  lastBookingId.value = null;
  lastBookingDetails.value = null;
};

const clearErrorMessage = () => {
  bookingErrorMessage.value = null;
  bookingErrorDetails.value = null;
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
    bookingSuccessMessage.value = null;
    bookingSuccessDetails.value = lastBookingDetails.value;
    lastBookingId.value = null;
    lastBookingDetails.value = null;
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

  // Load users list for colleague dropdown (non-blocking)
  loadUsers();

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
  if (storage) {
    storage.setItem('sithub_booking_mode', mode);
  }
  if (!activeItemGroupId.value) return;
  if (mode === 'week') {
    await loadWeekData(activeItemGroupId.value);
  } else {
    weekData.value = {};
    weekSelections.value = new Set();
    weekBookingResults.value = [];
    await loadItems(activeItemGroupId.value, selectedDate.value);
  }
});

watch([selectedWeek, showWeekends], async () => {
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

function formatDisplayDate(dateStr: string) {
  if (!dateStr) return '';
  const date = new Date(`${dateStr}T00:00:00`);
  if (Number.isNaN(date.getTime())) return dateStr;
  return new Intl.DateTimeFormat(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  }).format(date);
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
  grid-template-columns: repeat(5, 1fr); /* overridden by inline style when weekends enabled */
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

.week-day-status-truncated {
  display: inline-block;
}

.week-day-past {
  opacity: 0.5;
}

.week-day-expanded {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
}

.week-day-expanded-label {
  min-width: 90px;
}
</style>

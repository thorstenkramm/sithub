<template>
  <div class="area-weekly-matrix" data-cy="area-weekly-matrix">
    <!-- Loading -->
    <LoadingState v-if="loading" type="cards" :count="2" data-cy="matrix-loading" />

    <!-- Error -->
    <v-alert v-else-if="errorMessage" type="error" class="mb-4" data-cy="matrix-error">
      {{ errorMessage }}
    </v-alert>

    <!-- Matrix -->
    <div v-else-if="matrixData.length > 0" class="matrix-scroll-container" data-cy="matrix-container">
      <table class="matrix-table">
        <thead>
          <tr class="matrix-header-row" data-cy="matrix-header-row">
            <th class="matrix-corner sticky-col sticky-header" />
            <th
              v-for="day in days"
              :key="day.date"
              class="matrix-day-header sticky-header"
              :class="{ 'matrix-past-day': isPastDay(day.date) }"
              :data-cy="`matrix-day-${day.weekday}`"
            >
              <div class="day-header-weekday">{{ localizeWeekday(day.weekday, t) }}</div>
              <div class="day-header-date text-caption">{{ formatShortDate(day.date) }}</div>
            </th>
          </tr>
        </thead>
        <tbody>
          <template v-for="group in matrixData" :key="group.id">
            <AreaWeeklyMatrixRoomSection
              :group="group"
              :days="days"
              :collapsed="isCollapsed(group.id)"
              :current-user-id="currentUserId"
              :is-admin="isAdmin"
              :today="today"
              @toggle-collapse="toggleCollapse(group.id)"
            />
          </template>
        </tbody>
      </table>
    </div>

    <!-- Booking popover -->
    <MatrixBookingPopover
      v-if="activeBookItem && activeBookCell"
      v-model="showBookPopover"
      :activator-el="popoverActivator"
      :item="activeBookItem"
      :cell="activeBookCell"
      @booked="onBooked"
      @booking-conflict="onBookingConflict"
    />

    <!-- Cancel popover -->
    <MatrixCancelPopover
      v-if="activeCancelItem && activeCancelCell"
      v-model="showCancelPopover"
      :activator-el="popoverActivator"
      :item="activeCancelItem"
      :cell="activeCancelCell"
      @cancelled="onCancelled"
    />

    <!-- Snackbar feedback -->
    <v-snackbar
      v-model="showSnackbar"
      :color="snackbarColor"
      location="bottom"
      :timeout="snackbarTimeout"
      data-cy="matrix-snackbar"
    >
      {{ snackbarMessage }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, provide, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { fetchWeeklyMatrix } from '../../api/itemGroupMatrix';
import type { ItemGroupMatrixAttributes, MatrixDayMeta, MatrixCell, MatrixItem } from '../../api/itemGroupMatrix';
import type { JsonApiResource } from '../../api/types';
import { localizeWeekday } from '../../composables/useWeekSelector';
import { getSafeLocalStorage } from '../../composables/storage';
import { useLiveBookingRefresh } from '../../composables/useLiveBookingRefresh';
import { useAuthStore } from '../../stores/useAuthStore';
import { LoadingState } from '../../components';
import AreaWeeklyMatrixRoomSection from './AreaWeeklyMatrixRoomSection.vue';
import MatrixBookingPopover from './MatrixBookingPopover.vue';
import MatrixCancelPopover from './MatrixCancelPopover.vue';
import type { MatrixCellClickEvent } from './matrixTypes';

const props = defineProps<{
  areaId: string;
  week: string;
  showWeekends: boolean;
}>();

const { t } = useI18n();
const authStore = useAuthStore();
const loading = ref(false);
const errorMessage = ref<string | null>(null);
const matrixData = ref<JsonApiResource<ItemGroupMatrixAttributes>[]>([]);
const days = ref<MatrixDayMeta[]>([]);
const today = computed(() => new Date().toISOString().slice(0, 10));
const currentUserId = computed(() => authStore.userId ?? '');
const isAdmin = computed(() => authStore.isAdmin);

// Collapse state
const COLLAPSE_KEY_PREFIX = 'sithub_matrix_collapsed_';

function getCollapseKey(areaId: string): string {
  return COLLAPSE_KEY_PREFIX + areaId;
}

const collapsedRooms = ref<Set<string>>(new Set());

function loadCollapsedState() {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  try {
    const raw = storage.getItem(getCollapseKey(props.areaId));
    if (raw) {
      const ids = JSON.parse(raw) as string[];
      collapsedRooms.value = new Set(ids);
    } else {
      collapsedRooms.value = new Set();
    }
  } catch {
    collapsedRooms.value = new Set();
  }
}

function saveCollapsedState() {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  try {
    storage.setItem(getCollapseKey(props.areaId), JSON.stringify([...collapsedRooms.value]));
  } catch {
    // Storage full
  }
}

function isCollapsed(groupId: string): boolean {
  return collapsedRooms.value.has(groupId);
}

function toggleCollapse(groupId: string) {
  const next = new Set(collapsedRooms.value);
  if (next.has(groupId)) {
    next.delete(groupId);
  } else {
    next.add(groupId);
  }
  collapsedRooms.value = next;
  saveCollapsedState();
}

function isPastDay(dateStr: string): boolean {
  return dateStr < today.value;
}

function formatShortDate(dateStr: string): string {
  const [, m, d] = dateStr.split('-');
  return `${d}.${m}.`;
}

// Popover state
const popoverActivator = ref<HTMLElement>();
const showBookPopover = ref(false);
const activeBookItem = ref<MatrixItem | null>(null);
const activeBookCell = ref<MatrixCell | null>(null);
const showCancelPopover = ref(false);
const activeCancelItem = ref<MatrixItem | null>(null);
const activeCancelCell = ref<MatrixCell | null>(null);

// Snackbar
const showSnackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('success');
const snackbarTimeout = ref(3000);

function showFeedback(message: string, color: string, timeout: number) {
  snackbarMessage.value = message;
  snackbarColor.value = color;
  snackbarTimeout.value = timeout;
  showSnackbar.value = true;
}

// Provide cell click handler to children
provide('matrixCellClick', (event: MatrixCellClickEvent) => {
  // Close any open popover first
  showBookPopover.value = false;
  showCancelPopover.value = false;

  popoverActivator.value = event.el;

  if (event.type === 'book') {
    activeBookItem.value = event.item;
    activeBookCell.value = event.cell;
    // Defer opening to next tick so v-menu picks up the new activator
    requestAnimationFrame(() => {
      showBookPopover.value = true;
    });
  } else {
    activeCancelItem.value = event.item;
    activeCancelCell.value = event.cell;
    requestAnimationFrame(() => {
      showCancelPopover.value = true;
    });
  }
});

function onBooked() {
  showFeedback(t('matrix.bookingConfirmed'), 'success', 3000);
  loadMatrix();
}

function onBookingConflict() {
  loadMatrix();
}

function onCancelled() {
  showFeedback(t('matrix.cancelConfirmed'), 'success', 3000);
  loadMatrix();
}

async function loadMatrix(opts: { silent?: boolean } = {}) {
  if (!opts.silent) {
    loading.value = true;
    errorMessage.value = null;
  }
  try {
    const dayCount = props.showWeekends ? 7 : 5;
    const resp = await fetchWeeklyMatrix(props.areaId, props.week, dayCount);
    matrixData.value = resp.data;
    const first = resp.data[0];
    if (first) {
      days.value = first.attributes.days;
    } else if (!opts.silent) {
      days.value = [];
    }
  } catch {
    if (!opts.silent) {
      errorMessage.value = t('itemGroups.unableToLoad');
      matrixData.value = [];
      days.value = [];
    }
    // Silent refreshes ignore transient errors and keep the last known state.
  } finally {
    if (!opts.silent) {
      loading.value = false;
    }
  }
}

onMounted(() => {
  loadCollapsedState();
  loadMatrix();
});

watch(() => [props.week, props.showWeekends], () => {
  loadMatrix();
});

watch(() => props.areaId, () => {
  loadCollapsedState();
  loadMatrix();
});

useLiveBookingRefresh({
  refresh: () => loadMatrix({ silent: true }),
  isRelevant: (event) => {
    if (!days.value.some((day) => day.date === event.booking_date)) {
      return false;
    }
    return matrixData.value.some((group) =>
      group.attributes.items.some((item) => item.item_id === event.item_id)
    );
  }
});
</script>

<style scoped>
.matrix-scroll-container {
  overflow: auto;
  max-height: calc(100vh - 200px);
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 8px;
}

.matrix-table {
  border-collapse: separate;
  border-spacing: 0;
  width: 100%;
  min-width: 500px;
}

.sticky-header {
  position: sticky;
  top: 0;
  z-index: 3;
  background: rgb(var(--v-theme-surface));
}

.sticky-col {
  position: sticky;
  left: 0;
  z-index: 2;
  background: rgb(var(--v-theme-surface));
}

.sticky-header.sticky-col {
  z-index: 4;
}

.matrix-corner {
  min-width: 140px;
  width: 140px;
}

.matrix-day-header {
  text-align: center;
  padding: 8px 12px;
  min-width: 80px;
  border-bottom: 2px solid rgba(var(--v-border-color), var(--v-border-opacity));
  font-weight: 500;
}

.day-header-weekday {
  font-size: 0.85rem;
  font-weight: 600;
}

.day-header-date {
  opacity: 0.7;
}

.matrix-past-day {
  opacity: 0.5;
}
</style>

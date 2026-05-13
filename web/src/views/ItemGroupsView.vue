<template>
  <div class="page-container">
    <PageHeader
      title=""
      :breadcrumbs="breadcrumbs"
    />

    <!-- Toolbar -->
    <v-card class="mb-6" data-cy="week-selector-card">
      <v-card-text>
        <!-- Single row: week selector, floor plan, view switch, equipment filter (with info icon adjacent) -->
        <div class="ig-controls-row d-flex flex-wrap align-center ga-3">
          <v-select
            v-model="selectedWeek"
            :items="weekOptions"
            item-title="label"
            item-value="value"
            :label="$t('itemGroups.calendarWeek')"
            density="compact"
            hide-details
            data-cy="week-selector"
            class="ig-week-selector"
          />
          <FloorPlanButton
            v-if="areaFloorPlan"
            data-cy="area-floor-plan-btn"
            @click="showFloorPlanDialog = true"
          />

          <!-- View switch: Tiles / Table -->
          <div class="d-flex align-center" data-cy="view-switch-container">
            <span class="text-button mr-1" :class="activeView === 'cards' ? 'text-primary font-weight-bold' : 'text-medium-emphasis'">{{ $t('itemGroups.viewTiles') }}</span>
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
              <span data-cy="view-switch-tooltip">{{ $t('itemGroups.viewTableDesktopOnly') }}</span>
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
            <span class="text-button ml-1" :class="activeView === 'table' ? 'text-primary font-weight-bold' : 'text-medium-emphasis'">{{ $t('itemGroups.viewTable') }}</span>
          </div>

          <!-- Equipment filter cluster (input → info icon adjacent → save/delete on demand) -->
          <div v-if="itemGroups.length > 0" class="d-flex align-center ga-1 ig-equipment-filter-cluster">
            <v-combobox
              v-model="equipmentFilter"
              :items="savedFilterItems"
              :label="$t('itemGroups.filterEquipment')"
              density="compact"
              hide-details
              clearable
              prepend-inner-icon="$filterOutline"
              data-cy="ig-equipment-filter"
            />
            <v-btn
              icon
              variant="text"
              size="small"
              data-cy="ig-equipment-filter-info"
              :aria-label="$t('itemGroups.equipmentFilterHelp')"
              @click="showFilterHelp = true"
            >
              <v-icon>$info</v-icon>
            </v-btn>
            <!-- Save/delete: always rendered to keep the cluster width stable
                 regardless of filter content (no layout shift on type/clear). -->
            <v-tooltip
              :text="isCurrentFilterSaved ? $t('itemGroups.deleteSavedFilter') : $t('itemGroups.saveFilter')"
              location="top"
              :disabled="!equipmentFilter"
            >
              <template #activator="{ props: tooltipProps }">
                <v-btn
                  v-bind="tooltipProps"
                  icon
                  variant="text"
                  size="small"
                  :class="{ 'filter-action-placeholder': !equipmentFilter }"
                  :data-cy="isCurrentFilterSaved ? 'ig-equipment-filter-delete' : 'ig-equipment-filter-save'"
                  :aria-label="isCurrentFilterSaved ? $t('itemGroups.deleteSavedFilter') : $t('itemGroups.saveFilter')"
                  :aria-hidden="!equipmentFilter ? 'true' : undefined"
                  :tabindex="!equipmentFilter ? -1 : undefined"
                  @click="equipmentFilter && toggleSaveFilter()"
                >
                  <v-icon>{{ isCurrentFilterSaved ? '$delete' : '$save' }}</v-icon>
                </v-btn>
              </template>
            </v-tooltip>
          </div>
        </div>
      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="itemGroupsLoading" type="cards" :count="4" data-cy="item-groups-loading" />

    <!-- Error State -->
    <v-alert v-else-if="itemGroupsErrorMessage" type="error" class="mb-4" data-cy="item-groups-error">
      {{ itemGroupsErrorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!itemGroups.length"
      :title="$t('itemGroups.emptyTitle')"
      :message="$t('itemGroups.emptyMessage')"
      icon="$room"
      :action-text="$t('itemGroups.backToAreas')"
      action-to="/"
      data-cy="item-groups-empty"
    />

    <!-- Table View -->
    <AreaWeeklyMatrixView
      v-else-if="activeView === 'table'"
      :area-id="route.params.areaId as string"
      :week="selectedWeek"
      :show-weekends="showWeekends"
      :parsed-equipment-filter="parsedEquipmentFilter"
    />

    <!-- Item Groups Grid (Card View) -->
    <div v-else class="card-grid" data-cy="item-groups-list">
      <div
        v-for="ig in sortedItemGroups"
        :key="ig.id"
        class="item-filter-wrapper"
      >
        <v-card
          :class="['card-hover', { 'item-filtered-out': isItemGroupFilteredOut(ig.id) }]"
          role="button"
          tabindex="0"
          :aria-label="`Select items in ${ig.attributes.name}`"
          data-cy="item-group-item"
          @click="!isItemGroupFilteredOut(ig.id) && goToItems(ig.id)"
          @keydown.enter="!isItemGroupFilteredOut(ig.id) && goToItems(ig.id)"
        >
        <div v-if="isItemGroupFilteredOut(ig.id)" class="item-filtered-overlay">
          <span class="text-body-2 text-medium-emphasis">{{ $t('itemGroups.equipmentNotAvailable') }}</span>
        </div>
        <v-card-item>
          <template #prepend>
            <v-avatar color="secondary" variant="tonal" size="48">
              <v-icon size="24">{{ resolveIcon(ig.attributes.icon, '$room') }}</v-icon>
            </v-avatar>
          </template>
          <v-card-title class="text-h6">{{ ig.attributes.name }}</v-card-title>
          <v-card-subtitle v-if="ig.attributes.description">
            {{ ig.attributes.description }}
          </v-card-subtitle>
        </v-card-item>

        <!-- Weekly Availability Indicators -->
        <v-card-text v-if="availabilityMap[ig.id]" class="pt-0" data-cy="availability-indicators">
          <div class="d-flex ga-2">
            <span
              v-for="day in availabilityMap[ig.id]"
              :key="day.date"
              class="availability-indicator"
              :class="day.available > 0 ? 'available' : 'fully-booked'"
              :aria-label="formatAvailabilityAriaLabel(day)"
              :data-cy-weekday="day.weekday"
            >
              <span
                class="indicator-dot"
                :class="day.available > 0 ? 'dot-available' : 'dot-booked'"
              />
              <span class="indicator-label text-caption">{{ localizeWeekday(day.weekday, t) }}</span>
            </span>
          </div>
        </v-card-text>

        <v-card-actions class="px-4 pb-4">
          <v-btn
            color="primary"
            variant="tonal"
            size="small"
            @click.stop="goToItems(ig.id)"
          >
            {{ $t('itemGroups.select') }}
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            :to="{
              name: 'item-group-bookings',
              params: { itemGroupId: ig.id },
              query: { areaId: route.params.areaId as string }
            }"
            @click.stop
          >
            {{ $t('itemGroups.viewBookings') }}
          </v-btn>
        </v-card-actions>
      </v-card>
      </div>
    </div>

    <v-dialog
      v-model="showFloorPlanDialog"
      persistent
      :fullscreen="isCompactFloorPlanViewport"
      data-cy="floor-plan-dialog"
    >
      <v-card class="floor-plan-dialog-card">
        <v-card-text class="floor-plan-dialog-body">
          <InteractiveFloorPlan
            v-if="areaFloorPlan"
            :floor-plan="areaFloorPlan"
            :title="areaName || 'Floor Plan'"
            :week-label="weekOptions.find(o => o.value === selectedWeek)?.label || ''"
            :week-dates="selectedWeekDates"
            item-group-id=""
            :area-level="true"
            :area-id="route.params.areaId as string"
            @close="showFloorPlanDialog = false"
          />
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showFilterHelp" max-width="500">
      <v-card>
        <v-card-title>{{ $t('itemGroups.equipmentFilterHelp') }}</v-card-title>
        <v-card-text data-cy="ig-equipment-filter-help">
          <p class="mb-3">{{ $t('items.filterSyntaxDescription') }}</p>
          <ul class="mb-3">
            <li>{{ $t('items.filterSyntaxOr') }}</li>
            <li>{{ $t('items.filterSyntaxAnd') }}</li>
            <li>{{ $t('items.filterSyntaxExact') }}</li>
            <li>{{ $t('items.filterSyntaxCase') }}</li>
          </ul>
          <p class="text-caption text-medium-emphasis">{{ $t('items.filterSyntaxExample') }} <code>"27 inch display" + webcam</code></p>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showFilterHelp = false">{{ $t('common.close') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar
      :key="successSnackbarKey"
      v-model="showSuccessSnackbar"
      :timeout="successSnackbarTimeout"
      location="bottom"
      color="success"
      :data-cy="successSnackbarCy"
    >
      {{ successSnackbarMessage }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { fetchWeeklyAvailability } from '../api/itemGroupAvailability';
import type { DayAvailability } from '../api/itemGroupAvailability';
import type { ItemGroupAttributes } from '../api/itemGroups';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useWeekSelector, localizeWeekday } from '../composables/useWeekSelector';
import { useWeekendPreference } from '../composables/useWeekendPreference';
import { fetchItems } from '../api/items';
import { matchesParsedFilter, parseFilter } from '../composables/useEquipmentFilter';
import { useSavedFilters } from '../composables/useSavedFilters';
import { useFavorites } from '../composables/useFavorites';
import { useDateState } from '../composables/useDateState';
import { useLiveBookingRefresh } from '../composables/useLiveBookingRefresh';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '../stores/useAuthStore';
import { resolveConfiguredIcon } from '../utils/icons';
import { fetchSettings } from '../api/settings';
import { PageHeader, LoadingState, EmptyState, FloorPlanButton } from '../components';
import InteractiveFloorPlan from '../components/InteractiveFloorPlan.vue';
import AreaWeeklyMatrixView from '../components/area-weekly-matrix/AreaWeeklyMatrixView.vue';
import { useAreaViewPreference } from '../composables/useAreaViewPreference';

const authStore = useAuthStore();
const { t } = useI18n();
const areaName = ref('');
const areaFloorPlan = ref<string | null>(null);
const areaIcon = ref<string | null>(null);
const showFloorPlanDialog = ref(false);
const isCompactFloorPlanViewport = ref(false);
const isCompactViewport = ref(false);
const { activeView, load: loadViewPref, save: saveViewPref } = useAreaViewPreference();

const toggleView = (val: boolean | null) => {
  const areaId = route.params.areaId as string;
  const next = val ? 'table' : 'cards';
  saveViewPref(areaId, next);
};
const itemGroups = ref<JsonApiResource<ItemGroupAttributes>[]>([]);
const itemGroupsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: itemGroupsLoading, run: runItemGroups } = useApi();
// Favorites are no longer tied to item groups (story 31.2). The composable is
// imported only for side effects: visiting an item-group view used to mark the
// first-load purge of legacy item-group favorites; calling it here keeps that
// behaviour deterministic.
useFavorites();
const successSnackbarMessage = ref<string | null>(null);
const successSnackbarCy = ref('item-groups-success');
const successSnackbarTimeout = ref(3000);
const successSnackbarKey = ref(0);
const showSuccessSnackbar = computed({
  get: () => successSnackbarMessage.value !== null,
  set: (v: boolean) => {
    if (!v) {
      successSnackbarMessage.value = null;
      successSnackbarCy.value = 'item-groups-success';
      successSnackbarTimeout.value = 3000;
    }
  }
});
const availabilityMap = ref<Record<string, DayAvailability[]>>({});
const visibleAreaItemIds = ref(new Set<string>());

const formatAvailabilityAriaLabel = (day: DayAvailability): string => {
  const weekday = localizeWeekday(day.weekday, t);
  const status = day.available > 0
    ? `${day.available} ${t('status.available').toLowerCase()}`
    : t('common.booked');
  return `${weekday}: ${status}`;
};

const equipmentFilter = ref<string | null>('');
watch(equipmentFilter, (value) => {
  if (value === null || value === undefined) {
    equipmentFilter.value = '';
  }
});
const showFilterHelp = ref(false);
const { comboboxItems: savedFilterItems, saveFilter, deleteFilter, isSavedFilter } = useSavedFilters();
const isCurrentFilterSaved = computed(() => !!equipmentFilter.value && isSavedFilter(equipmentFilter.value));
const showSuccessFeedback = (message: string, cy: string) => {
  successSnackbarKey.value += 1;
  successSnackbarMessage.value = message;
  successSnackbarCy.value = cy;
  successSnackbarTimeout.value = 3000;
};
const showFilterFeedback = (message: string) => {
  showSuccessFeedback(message, 'ig-filter-message');
};
const toggleSaveFilter = () => {
  if (!equipmentFilter.value) return;
  if (isCurrentFilterSaved.value) {
    deleteFilter(equipmentFilter.value);
    equipmentFilter.value = '';
    showFilterFeedback(t('itemGroups.savedFilterDeleted'));
  } else {
    if (saveFilter(equipmentFilter.value)) {
      showFilterFeedback(t('itemGroups.filterSaved'));
    }
  }
};
const parsedEquipmentFilter = computed(() => parseFilter(equipmentFilter.value ?? ''));
const itemGroupEquipment = ref<Record<string, string[]>>({});

const isRelevantLiveEvent = (event: { item_id: string; booking_date: string }): boolean => {
  if (activeView.value !== 'cards') {
    return false;
  }
  if (!selectedWeekDates.value.includes(event.booking_date)) {
    return false;
  }
  return visibleAreaItemIds.value.size === 0 || visibleAreaItemIds.value.has(event.item_id);
};

const isItemGroupFilteredOut = (igId: string): boolean => {
  if (!equipmentFilter.value) return false;
  const equipment = itemGroupEquipment.value[igId] ?? [];
  return !matchesParsedFilter(equipment, parsedEquipmentFilter.value);
};

const { showWeekends } = useWeekendPreference();
const weeksInAdvanced = ref(7);
const { weekOptions, selectedWeek, selectedWeekDates } = useWeekSelector(showWeekends, weeksInAdvanced);
const { getWeek, setWeek } = useDateState();

// Restore memorized week on mount
const storedWeek = getWeek();
if (weekOptions.value.some(o => o.value === storedWeek)) {
  selectedWeek.value = storedWeek;
}
const breadcrumbs = computed(() => [
  { text: t('common.home'), to: '/' },
  { text: areaName.value || t('common.area') }
]);

// Item groups render in YAML order. Per story 31.2 we no longer hoist
// favorites to the top of this view — favorites live in the dedicated virtual
// "Favorites" area accessed from the home overview.
const sortedItemGroups = computed(() => itemGroups.value);

const resolveIcon = (icon: string | undefined, fallback: string) => {
  return resolveConfiguredIcon(icon || areaIcon.value, fallback);
};

const goToItems = async (igId: string) => {
  const areaId = route.params.areaId as string;
  await router.push({ name: 'items', params: { itemGroupId: igId }, query: { areaId } });
};

const handleAuthError = async (err: unknown) => {
  if (err instanceof ApiError && err.status === 401) {
    window.location.href = '/login';
    return true;
  }
  if (err instanceof ApiError && err.status === 403) {
    await router.push('/access-denied');
    return true;
  }
  return false;
};

const loadAvailability = async (areaId: string, week: string) => {
  try {
    const days = showWeekends.value ? 7 : undefined;
    const resp = await fetchWeeklyAvailability(areaId, week, days);
    const map: Record<string, DayAvailability[]> = {};
    for (const resource of resp.data) {
      map[resource.attributes.item_group_id] = resource.attributes.days;
    }
    availabilityMap.value = map;
  } catch {
    // Non-critical: availability indicators just won't show
    availabilityMap.value = {};
  }
};

const updateViewport = () => {
  if (typeof window.matchMedia !== 'function') {
    isCompactFloorPlanViewport.value = false;
    isCompactViewport.value = false;
    return;
  }

  const narrow = window.matchMedia('(max-width: 768px)').matches;
  const short = window.matchMedia('(max-height: 500px)').matches;
  isCompactFloorPlanViewport.value = narrow || short;
  isCompactViewport.value = narrow;
};

const handleResize = () => {
  updateViewport();
};

onMounted(async () => {
  updateViewport();
  loadViewPref(route.params.areaId as string, !isCompactViewport.value);
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
      itemGroupsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    throw err;
  }

  // Fetch booking settings (non-blocking, uses default on failure)
  try {
    const settingsResp = await fetchSettings();
    weeksInAdvanced.value = settingsResp.data.attributes.weeks_in_advanced;
  } catch {
    // Non-critical: week selector uses default
  }

  const areaId = route.params.areaId;
  if (typeof areaId !== 'string' || areaId.trim() === '') {
    itemGroupsErrorMessage.value = t('areas.notFound');
    return;
  }

  // Fetch area name for breadcrumb
  try {
    const areasResp = await fetchAreas();
    const area = areasResp.data.find(a => a.id === areaId);
    if (area) {
      areaName.value = area.attributes.name;
      areaFloorPlan.value = area.attributes.floor_plan || null;
      areaIcon.value = area.attributes.icon || null;
    }
  } catch (err) {
    if (isConnectionError(err)) {
      itemGroupsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    // Ignore other errors - breadcrumb will just show "Area"
  }

  try {
    const resp = await runItemGroups(() => fetchItemGroups(areaId));
    itemGroups.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (isConnectionError(err)) {
      itemGroupsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      itemGroupsErrorMessage.value = t('areas.notFound');
      return;
    }
    itemGroupsErrorMessage.value = t('itemGroups.unableToLoad');
  }

  await loadAvailability(areaId, selectedWeek.value);

  // Load equipment per item group for filtering (non-blocking)
  visibleAreaItemIds.value = new Set();
  if (itemGroups.value.length > 0) {
    try {
      const results = await Promise.all(
        itemGroups.value.map(ig => fetchItems(ig.id).then(r => ({ igId: ig.id, items: r.data })))
      );
      const map: Record<string, string[]> = {};
      const nextVisibleItemIds = new Set<string>();
      for (const { igId, items } of results) {
        const allEquipment = new Set<string>();
        for (const item of items) {
          nextVisibleItemIds.add(item.id);
          for (const eq of item.attributes.equipment ?? []) {
            allEquipment.add(eq);
          }
        }
        map[igId] = [...allEquipment];
      }
      itemGroupEquipment.value = map;
      visibleAreaItemIds.value = nextVisibleItemIds;
    } catch {
      // Non-critical: filter just won't work
      visibleAreaItemIds.value = new Set();
    }
  }
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

watch([selectedWeek, showWeekends], async ([week]) => {
  setWeek(week);
  const areaId = route.params.areaId;
  if (typeof areaId === 'string' && areaId.trim() !== '') {
    await loadAvailability(areaId, week);
  }
});

useLiveBookingRefresh({
  refresh: async () => {
    const areaId = route.params.areaId;
    if (typeof areaId !== 'string' || areaId.trim() === '') return;
    await loadAvailability(areaId, selectedWeek.value);
  },
  isRelevant: isRelevantLiveEvent
});


</script>

<style scoped>
.ig-controls-row {
  min-height: 40px;
}

.ig-week-selector {
  flex: 0 0 240px;
  max-width: 240px;
  min-width: 200px;
}

.ig-equipment-filter-cluster {
  flex: 1 1 300px;
  min-width: 240px;
  max-width: 420px;
}

.filter-action-placeholder {
  visibility: hidden;
  pointer-events: none;
}

@media (max-width: 600px) {
  .ig-week-selector,
  .ig-equipment-filter-cluster {
    flex: 1 1 100%;
    max-width: 100%;
  }
}

.availability-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 28px;
}

.indicator-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  display: inline-block;
}

.dot-available {
  background-color: rgb(var(--v-theme-success));
}

.dot-booked {
  background-color: transparent;
  border: 2px solid rgb(var(--v-theme-error));
}

.indicator-label {
  font-size: 0.65rem;
  line-height: 1;
}

.item-filter-wrapper {
  position: relative;
}

.item-filtered-out {
  filter: blur(3px);
  opacity: 0.5;
  pointer-events: none;
}

.item-filtered-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.floor-plan-dialog-card {
  height: 100%;
}

.floor-plan-dialog-body {
  height: 100%;
}

.view-switch {
  flex: none;
}
</style>

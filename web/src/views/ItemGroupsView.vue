<template>
  <div class="page-container">
    <PageHeader
      title=""
      :breadcrumbs="breadcrumbs"
    />

    <!-- Week Selector -->
    <v-card class="mb-6" data-cy="week-selector-card">
      <v-card-text>
        <div class="d-flex flex-wrap align-end ga-4">
          <v-select
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
          <v-btn
            v-if="areaFloorPlan"
            variant="outlined"
            size="small"
            prepend-icon="$map"
            data-cy="area-floor-plan-btn"
            @click="showFloorPlanDialog = true"
          >
            Floor plan
          </v-btn>
        </div>
      </v-card-text>
    </v-card>

    <!-- Equipment Filter -->
    <v-card v-if="itemGroups.length > 0" class="mb-6">
      <v-card-text>
        <div class="d-flex align-center ga-2" style="max-width: 420px;">
          <v-combobox
            v-model="equipmentFilter"
            :items="savedFilterItems"
            label="Filter equipment"
            density="compact"
            hide-details
            clearable
            prepend-inner-icon="$filterOutline"
            data-cy="ig-equipment-filter"
          />
          <v-tooltip :text="isCurrentFilterSaved ? 'Delete saved filter' : 'Save filter'" location="top">
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                icon
                variant="text"
                size="small"
                :data-cy="isCurrentFilterSaved ? 'ig-equipment-filter-delete' : 'ig-equipment-filter-save'"
                :aria-label="isCurrentFilterSaved ? 'Delete saved filter' : 'Save filter'"
                @click="toggleSaveFilter"
              >
                <v-icon>{{ isCurrentFilterSaved ? '$delete' : '$plus' }}</v-icon>
              </v-btn>
            </template>
          </v-tooltip>
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
      title="No item groups available"
      message="This area doesn't have any item groups configured yet."
      icon="$room"
      action-text="Back to Areas"
      action-to="/"
      data-cy="item-groups-empty"
    />

    <!-- Item Groups Grid -->
    <div v-else class="card-grid" data-cy="item-groups-list">
      <!-- Third-level favorites promoted to this view -->
      <div
        v-for="fav in sortedItemGroups.thirdLevelFavs"
        :key="`fav-item-${fav.areaId}-${fav.itemGroupId}-${fav.itemId}`"
        class="item-filter-wrapper"
      >
        <v-card
          :class="['card-hover', { 'item-filtered-out': isItemGroupFilteredOut(fav.itemGroupId) }]"
          role="button"
          tabindex="0"
          data-cy="favorite-item-tile"
          @click="!isItemGroupFilteredOut(fav.itemGroupId) && goToItems(fav.itemGroupId)"
        >
          <div v-if="isItemGroupFilteredOut(fav.itemGroupId)" class="item-filtered-overlay">
            <span class="text-body-2 text-medium-emphasis">equipment not available</span>
          </div>
          <v-card-item>
            <template #prepend>
              <v-avatar color="primary" variant="tonal" size="48">
                <v-icon size="24">$desk</v-icon>
              </v-avatar>
            </template>
            <v-card-title class="text-h6">{{ fav.itemName }}</v-card-title>
            <v-card-subtitle>{{ fav.itemGroupName }}</v-card-subtitle>
          </v-card-item>

          <!-- Weekly Availability Indicators (from parent item group) -->
          <v-card-text v-if="availabilityMap[fav.itemGroupId]" class="pt-0" data-cy="availability-indicators">
            <div class="d-flex ga-2">
              <span
                v-for="day in availabilityMap[fav.itemGroupId]"
                :key="day.date"
                class="availability-indicator"
                :class="day.available > 0 ? 'available' : 'fully-booked'"
                :aria-label="`${day.weekday}: ${day.available > 0 ? day.available + ' available' : 'fully booked'}`"
                :data-cy-weekday="day.weekday"
              >
                <span
                  class="indicator-dot"
                  :class="day.available > 0 ? 'dot-available' : 'dot-booked'"
                />
                <span class="indicator-label text-caption">{{ day.weekday }}</span>
              </span>
            </div>
          </v-card-text>

          <v-card-actions class="px-4 pb-4">
            <v-btn
              color="primary"
              variant="tonal"
              size="small"
              @click.stop="goToItems(fav.itemGroupId)"
            >
              Select
            </v-btn>
            <v-btn
              variant="text"
              size="small"
              :to="{
                name: 'item-group-bookings',
                params: { itemGroupId: fav.itemGroupId },
                query: { areaId: route.params.areaId as string }
              }"
              @click.stop
            >
              View Bookings
            </v-btn>
            <v-spacer />
            <v-btn
              icon
              variant="text"
              size="small"
              data-cy="favorite-item-heart"
              @click.stop="handleToggleItemFavorite(fav)"
            >
              <v-icon color="error">$heart</v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </div>

      <!-- Second-level favorites (sorted A-Z) then rest (YAML order) -->
      <div
        v-for="ig in [...sortedItemGroups.igFavs, ...sortedItemGroups.rest]"
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
          <span class="text-body-2 text-medium-emphasis">equipment not available</span>
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
              :aria-label="`${day.weekday}: ${day.available > 0 ? day.available + ' available' : 'fully booked'}`"
              :data-cy-weekday="day.weekday"
            >
              <span
                class="indicator-dot"
                :class="day.available > 0 ? 'dot-available' : 'dot-booked'"
              />
              <span class="indicator-label text-caption">{{ day.weekday }}</span>
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
            Select
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
            View Bookings
          </v-btn>
          <v-spacer />
          <v-btn
            icon
            variant="text"
            size="small"
            data-cy="ig-favorite-heart"
            @click.stop="handleToggleItemGroupFavorite(ig.id, ig.attributes.name)"
          >
            <v-icon :color="isItemGroupFavorite(route.params.areaId as string, ig.id) ? 'error' : undefined">
              {{ isItemGroupFavorite(route.params.areaId as string, ig.id) ? '$heart' : '$heartOutline' }}
            </v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
      </div>
    </div>

    <v-dialog
      v-model="showFloorPlanDialog"
      max-width="1100"
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
            @close="showFloorPlanDialog = false"
          />
        </v-card-text>
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
import { useWeekSelector } from '../composables/useWeekSelector';
import { useWeekendPreference } from '../composables/useWeekendPreference';
import { fetchItems } from '../api/items';
import { matchesParsedFilter, parseFilter } from '../composables/useEquipmentFilter';
import { useSavedFilters } from '../composables/useSavedFilters';
import { useFavorites } from '../composables/useFavorites';
import { useDateState } from '../composables/useDateState';
import { useAuthStore } from '../stores/useAuthStore';
import { resolveConfiguredIcon } from '../utils/icons';
import { PageHeader, LoadingState, EmptyState } from '../components';
import InteractiveFloorPlan from '../components/InteractiveFloorPlan.vue';

const authStore = useAuthStore();
const areaName = ref('');
const areaFloorPlan = ref<string | null>(null);
const areaIcon = ref<string | null>(null);
const showFloorPlanDialog = ref(false);
const isCompactFloorPlanViewport = ref(false);
const itemGroups = ref<JsonApiResource<ItemGroupAttributes>[]>([]);
const itemGroupsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: itemGroupsLoading, run: runItemGroups } = useApi();
const {
  isItemGroupFavorite, toggleItemGroupFavorite,
  favoriteItems, toggleItemFavorite
} = useFavorites();
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

const equipmentFilter = ref('');
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
    showFilterFeedback('Saved filter deleted.');
  } else {
    if (saveFilter(equipmentFilter.value)) {
      showFilterFeedback('Filter saved.');
    }
  }
};
const parsedEquipmentFilter = computed(() => parseFilter(equipmentFilter.value));
const itemGroupEquipment = ref<Record<string, string[]>>({});

const isItemGroupFilteredOut = (igId: string): boolean => {
  if (!equipmentFilter.value) return false;
  const equipment = itemGroupEquipment.value[igId] ?? [];
  return !matchesParsedFilter(equipment, parsedEquipmentFilter.value);
};

const { showWeekends } = useWeekendPreference();
const { weekOptions, selectedWeek, selectedWeekDates } = useWeekSelector(showWeekends);
const { getWeek, setWeek } = useDateState();

// Restore memorized week on mount
const storedWeek = getWeek();
if (weekOptions.value.some(o => o.value === storedWeek)) {
  selectedWeek.value = storedWeek;
}
const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area' }
]);

const sortedItemGroups = computed(() => {
  const areaId = route.params.areaId as string;
  // Third-level favorites for this area's item groups
  const thirdLevelFavs = favoriteItems.value
    .filter(f => f.areaId === areaId && itemGroups.value.some(ig => ig.id === f.itemGroupId))
    .sort((a, b) => `${a.itemGroupName} ${a.itemName}`.localeCompare(`${b.itemGroupName} ${b.itemName}`));

  // Second-level favorites sorted A-Z
  const igFavs = itemGroups.value
    .filter(ig => isItemGroupFavorite(areaId, ig.id))
    .sort((a, b) => a.attributes.name.localeCompare(b.attributes.name));
  const igFavIds = new Set(igFavs.map(ig => ig.id));

  // Rest in YAML order, minus second-level favorites
  const rest = itemGroups.value.filter(ig => !igFavIds.has(ig.id));

  return { thirdLevelFavs, igFavs, rest, areaId };
});

const handleToggleItemGroupFavorite = (igId: string, igName: string) => {
  const areaId = route.params.areaId as string;
  const { added } = toggleItemGroupFavorite(areaId, igId);
  showSuccessFeedback(
    added ? `${igName} saved as favorite.` : `${igName} removed from favorites.`,
    'favorite-message'
  );
};

const handleToggleItemFavorite = (fav: { itemId: string; itemName: string; itemGroupId: string; itemGroupName: string }) => {
  const areaId = route.params.areaId as string;
  const { added } = toggleItemFavorite({ ...fav, areaId });
  const label = `${fav.itemGroupName} ${fav.itemName}`;
  showSuccessFeedback(
    added ? `${label} saved as favorite.` : `${label} removed from favorites.`,
    'favorite-message'
  );
};

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
    return;
  }

  const narrow = window.matchMedia('(max-width: 768px)').matches;
  const short = window.matchMedia('(max-height: 500px)').matches;
  isCompactFloorPlanViewport.value = narrow || short;
};

const handleResize = () => {
  updateViewport();
};

onMounted(async () => {
  updateViewport();
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

  const areaId = route.params.areaId;
  if (typeof areaId !== 'string' || areaId.trim() === '') {
    itemGroupsErrorMessage.value = 'Area not found.';
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
      itemGroupsErrorMessage.value = 'Area not found.';
      return;
    }
    itemGroupsErrorMessage.value = 'Unable to load item groups.';
  }

  await loadAvailability(areaId, selectedWeek.value);

  // Load equipment per item group for filtering (non-blocking)
  if (itemGroups.value.length > 0) {
    try {
      const results = await Promise.all(
        itemGroups.value.map(ig => fetchItems(ig.id).then(r => ({ igId: ig.id, items: r.data })))
      );
      const map: Record<string, string[]> = {};
      for (const { igId, items } of results) {
        const allEquipment = new Set<string>();
        for (const item of items) {
          for (const eq of item.attributes.equipment ?? []) {
            allEquipment.add(eq);
          }
        }
        map[igId] = [...allEquipment];
      }
      itemGroupEquipment.value = map;
    } catch {
      // Non-critical: filter just won't work
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


</script>

<style scoped>
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
</style>

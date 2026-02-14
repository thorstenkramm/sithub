<template>
  <div class="page-container">
    <PageHeader
      title=""
      :breadcrumbs="breadcrumbs"
    />

    <!-- Week Selector -->
    <v-card class="mb-6" data-cy="week-selector-card">
      <v-card-text>
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
      <v-card
        v-for="ig in itemGroups"
        :key="ig.id"
        class="card-hover"
        role="button"
        tabindex="0"
        :aria-label="`Select items in ${ig.attributes.name}`"
        data-cy="item-group-item"
        @click="goToItems(ig.id)"
        @keydown.enter="goToItems(ig.id)"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar color="secondary" variant="tonal" size="48">
              <v-icon size="24">$room</v-icon>
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
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { fetchWeeklyAvailability } from '../api/itemGroupAvailability';
import type { DayAvailability } from '../api/itemGroupAvailability';
import type { ItemGroupAttributes } from '../api/itemGroups';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useWeekSelector } from '../composables/useWeekSelector';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState } from '../components';

const authStore = useAuthStore();
const areaName = ref('');
const itemGroups = ref<JsonApiResource<ItemGroupAttributes>[]>([]);
const itemGroupsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: itemGroupsLoading, run: runItemGroups } = useApi();
const availabilityMap = ref<Record<string, DayAvailability[]>>({});

const { weekOptions, selectedWeek } = useWeekSelector();

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area' }
]);

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
    const resp = await fetchWeeklyAvailability(areaId, week);
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
    }
  } catch {
    // Ignore - breadcrumb will just show "Area"
  }

  try {
    const resp = await runItemGroups(() => fetchItemGroups(areaId));
    itemGroups.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      itemGroupsErrorMessage.value = 'Area not found.';
      return;
    }
    itemGroupsErrorMessage.value = 'Unable to load item groups.';
  }

  await loadAvailability(areaId, selectedWeek.value);
});

watch(selectedWeek, async (week) => {
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
</style>

<template>
  <div class="page-container">
    <PageHeader
      title="Item Groups"
      subtitle="Select an item group to view available items"
      :breadcrumbs="breadcrumbs"
    />

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
        :aria-label="`View items in ${ig.attributes.name}`"
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
        <v-card-actions class="px-4 pb-4">
          <v-btn
            color="primary"
            variant="tonal"
            size="small"
            @click.stop="goToItems(ig.id)"
          >
            View Items
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            :to="{ name: 'item-group-bookings', params: { itemGroupId: ig.id } }"
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
import { onMounted, ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import type { ItemGroupAttributes } from '../api/itemGroups';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState } from '../components';

const authStore = useAuthStore();
const areaName = ref('');
const itemGroups = ref<JsonApiResource<ItemGroupAttributes>[]>([]);
const itemGroupsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: itemGroupsLoading, run: runItemGroups } = useApi();

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area' }
]);

const goToItems = async (igId: string) => {
  await router.push({ name: 'items', params: { itemGroupId: igId } });
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
});
</script>

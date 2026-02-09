<template>
  <div class="page-container">
    <PageHeader
      title="Areas"
      subtitle="Select an area to view available item groups and items"
      :breadcrumbs="[{ text: 'Home' }]"
    />

    <!-- Loading State -->
    <LoadingState v-if="areasLoading" type="cards" :count="4" data-cy="areas-loading" />

    <!-- Error State -->
    <v-alert v-else-if="areasError" type="error" class="mb-4" data-cy="areas-error">
      Unable to load areas. Please try again later.
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!areas.length"
      title="No areas available"
      message="There are no office areas configured yet. Contact your administrator to set up areas."
      icon="$area"
      data-cy="areas-empty"
    />

    <!-- Areas Grid -->
    <div v-else class="card-grid" data-cy="areas-list">
      <v-card
        v-for="area in areas"
        :key="area.id"
        class="card-hover"
        role="button"
        tabindex="0"
        :aria-label="`View item groups in ${area.attributes.name}`"
        data-cy="area-item"
        @click="goToItemGroups(area.id)"
        @keydown.enter="goToItemGroups(area.id)"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar color="primary" variant="tonal" size="48">
              <v-icon size="24">$area</v-icon>
            </v-avatar>
          </template>
          <v-card-title class="text-h6">{{ area.attributes.name }}</v-card-title>
          <v-card-subtitle v-if="area.attributes.description">
            {{ area.attributes.description }}
          </v-card-subtitle>
        </v-card-item>
        <v-card-actions class="px-4 pb-4">
          <v-btn
            color="primary"
            variant="tonal"
            size="small"
            @click.stop="goToItemGroups(area.id)"
          >
            View Item Groups
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            :to="{ name: 'area-presence', params: { areaId: area.id } }"
            data-cy="area-presence-link"
            @click.stop
          >
            Today's Presence
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>

    <!-- Admin Controls (hidden, kept for compatibility) -->
    <div v-if="isAdmin" class="mt-6 d-none">
      <v-btn data-cy="admin-cancel" size="small" variant="tonal">Cancel booking (admin)</v-btn>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { fetchAreas } from '../api/areas';
import { fetchMe } from '../api/me';
import type { AreaAttributes } from '../api/areas';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState } from '../components';

const authStore = useAuthStore();
const isAdmin = ref(false);
const areas = ref<JsonApiResource<AreaAttributes>[]>([]);
const router = useRouter();
const { loading: areasLoading, error: areasError, run: runAreas } = useApi();
const { handleAuthError } = useAuthErrorHandler();

const goToItemGroups = async (areaId: string) => {
  await router.push({ name: 'item-groups', params: { areaId } });
};

onMounted(async () => {
  try {
    const resp = await fetchMe();
    authStore.userName = resp.data.attributes.display_name;
    authStore.isAdmin = resp.data.attributes.is_admin;
    isAdmin.value = resp.data.attributes.is_admin;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    throw err;
  }

  try {
    const resp = await runAreas(() => fetchAreas());
    areas.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
  }
});
</script>

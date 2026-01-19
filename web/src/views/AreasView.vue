<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="areas-title">
            Areas
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
        <v-card-text>
          <v-progress-linear
            v-if="areasLoading"
            class="mb-3"
            indeterminate
            data-cy="areas-loading"
            aria-label="Loading areas"
          />
          <v-alert v-else-if="areasError" type="error" variant="tonal" data-cy="areas-error">
            Unable to load areas.
          </v-alert>
          <div v-else>
            <v-list v-if="areas.length" data-cy="areas-list">
              <v-list-item
                v-for="area in areas"
                :key="area.id"
                data-cy="area-item"
                @click="goToRooms(area.id)"
              >
                <v-list-item-title>{{ area.attributes.name }}</v-list-item-title>
                <template #append>
                  <router-link
                    :to="{ name: 'area-presence', params: { areaId: area.id } }"
                    class="text-caption mr-2"
                    data-cy="area-presence-link"
                    @click.stop
                  >
                    View Presence
                  </router-link>
                </template>
              </v-list-item>
            </v-list>
            <div v-else class="text-caption" data-cy="areas-empty">No areas available.</div>
          </div>
          <div v-if="isAdmin" class="mt-2">
            <div class="text-caption">Admin-only cancellation controls</div>
            <v-btn data-cy="admin-cancel" size="small" variant="tonal">Cancel booking (admin)</v-btn>
          </div>
        </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchAreas } from '../api/areas';
import { fetchMe } from '../api/me';
import type { AreaAttributes } from '../api/areas';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const isAdmin = ref(false);
const areas = ref<JsonApiResource<AreaAttributes>[]>([]);
const router = useRouter();
const { loading: areasLoading, error: areasError, run: runAreas } = useApi();

const goToRooms = async (areaId: string) => {
  await router.push({ name: 'rooms', params: { areaId } });
};

const handleAuthError = async (err: unknown) => {
  if (err instanceof ApiError && err.status === 401) {
    window.location.href = '/oauth/login';
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
    userName.value = resp.data.attributes.display_name;
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

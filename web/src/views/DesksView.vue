<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="desks-title">
            Desks
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
          <v-card-text>
            <v-progress-linear
              v-if="desksLoading"
              class="mb-3"
              indeterminate
              data-cy="desks-loading"
              aria-label="Loading desks"
            />
            <v-alert v-else-if="desksErrorMessage" type="error" variant="tonal" data-cy="desks-error">
              {{ desksErrorMessage }}
            </v-alert>
            <div v-else>
              <v-list v-if="desks.length" data-cy="desks-list">
                <v-list-item v-for="desk in desks" :key="desk.id" data-cy="desk-item">
                  <v-list-item-title>{{ desk.attributes.name }}</v-list-item-title>
                  <v-list-item-subtitle>
                    <ul class="pl-4" data-cy="desk-equipment">
                      <li v-for="item in desk.attributes.equipment" :key="item">{{ item }}</li>
                    </ul>
                    <div v-if="desk.attributes.warning" class="text-caption mt-1" data-cy="desk-warning">
                      {{ desk.attributes.warning }}
                    </div>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="desks-empty">No desks available.</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import type { DeskAttributes } from '../api/desks';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const desks = ref<JsonApiResource<DeskAttributes>[]>([]);
const desksErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: desksLoading, run: runDesks } = useApi();

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
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    throw err;
  }

  const roomId = route.params.roomId;
  if (typeof roomId !== 'string' || roomId.trim() === '') {
    desksErrorMessage.value = 'Room not found.';
    return;
  }

  try {
    const resp = await runDesks(() => fetchDesks(roomId));
    desks.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      desksErrorMessage.value = 'Room not found.';
      return;
    }
    desksErrorMessage.value = 'Unable to load desks.';
  }
});
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="area-presence-title">
            Today's Presence
          </v-card-title>
          <v-card-text>
            <div class="mb-4">
              <label class="text-caption font-weight-medium" for="presence-date">Date</label>
              <input
                id="presence-date"
                v-model="selectedDate"
                class="d-block mt-1"
                type="date"
                data-cy="presence-date"
                aria-label="Select date"
              />
            </div>
            <v-progress-linear
              v-if="loading"
              class="mb-3"
              indeterminate
              data-cy="presence-loading"
              aria-label="Loading presence"
            />
            <v-alert v-else-if="errorMessage" type="error" variant="tonal" data-cy="presence-error">
              {{ errorMessage }}
            </v-alert>
            <div v-else>
              <v-list v-if="presence.length" data-cy="presence-list">
                <v-list-item
                  v-for="entry in presence"
                  :key="entry.id"
                  data-cy="presence-item"
                >
                  <v-list-item-title>{{ entry.attributes.user_name || 'Unknown' }}</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ entry.attributes.room_name }} - {{ entry.attributes.desk_name }}
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <div v-else class="text-caption" data-cy="presence-empty">
                No one is scheduled for this date.
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchAreaPresence } from '../api/areaPresence';
import type { PresenceAttributes } from '../api/areaPresence';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const presence = ref<JsonApiResource<PresenceAttributes>[]>([]);
const errorMessage = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const route = useRoute();
const router = useRouter();
const { loading, run } = useApi();
const activeAreaId = ref<string | null>(null);

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

const loadPresence = async (areaId: string, date: string) => {
  errorMessage.value = null;
  try {
    const resp = await run(() => fetchAreaPresence(areaId, date));
    presence.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      errorMessage.value = 'Area not found.';
      return;
    }
    errorMessage.value = 'Unable to load presence.';
  }
};

onMounted(async () => {
  const areaId = route.params.areaId;
  if (typeof areaId !== 'string' || areaId.trim() === '') {
    errorMessage.value = 'Area not found.';
    return;
  }

  activeAreaId.value = areaId;
  await loadPresence(areaId, selectedDate.value);
});

watch(
  selectedDate,
  async (value) => {
    if (!activeAreaId.value) {
      return;
    }
    await loadPresence(activeAreaId.value, value);
  },
  { flush: 'post' }
);

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

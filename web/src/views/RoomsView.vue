<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="rooms-title">
            Rooms
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
          <v-card-text>
            <v-progress-linear
              v-if="roomsLoading"
              class="mb-3"
              indeterminate
              data-cy="rooms-loading"
              aria-label="Loading rooms"
            />
            <v-alert v-else-if="roomsErrorMessage" type="error" variant="tonal" data-cy="rooms-error">
              {{ roomsErrorMessage }}
            </v-alert>
            <div v-else>
              <v-list v-if="rooms.length" data-cy="rooms-list">
              <v-list-item
                v-for="room in rooms"
                :key="room.id"
                data-cy="room-item"
                @click="goToDesks(room.id)"
              >
                <v-list-item-title>{{ room.attributes.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
              <div v-else class="text-caption" data-cy="rooms-empty">No rooms available.</div>
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
import { fetchMe } from '../api/me';
import { fetchRooms } from '../api/rooms';
import type { RoomAttributes } from '../api/rooms';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';

const userName = ref('');
const rooms = ref<JsonApiResource<RoomAttributes>[]>([]);
const roomsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: roomsLoading, run: runRooms } = useApi();

const goToDesks = async (roomId: string) => {
  await router.push({ name: 'desks', params: { roomId } });
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
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    throw err;
  }

  const areaId = route.params.areaId;
  if (typeof areaId !== 'string' || areaId.trim() === '') {
    roomsErrorMessage.value = 'Area not found.';
    return;
  }

  try {
    const resp = await runRooms(() => fetchRooms(areaId));
    rooms.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      roomsErrorMessage.value = 'Area not found.';
      return;
    }
    roomsErrorMessage.value = 'Unable to load rooms.';
  }
});
</script>

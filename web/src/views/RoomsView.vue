<template>
  <div class="page-container">
    <PageHeader
      title="Rooms"
      subtitle="Select a room to view available desks"
      :breadcrumbs="breadcrumbs"
    />

    <!-- Loading State -->
    <LoadingState v-if="roomsLoading" type="cards" :count="4" data-cy="rooms-loading" />

    <!-- Error State -->
    <v-alert v-else-if="roomsErrorMessage" type="error" class="mb-4" data-cy="rooms-error">
      {{ roomsErrorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!rooms.length"
      title="No rooms available"
      message="This area doesn't have any rooms configured yet."
      icon="$room"
      action-text="Back to Areas"
      action-to="/"
      data-cy="rooms-empty"
    />

    <!-- Rooms Grid -->
    <div v-else class="card-grid" data-cy="rooms-list">
      <v-card
        v-for="room in rooms"
        :key="room.id"
        class="card-hover"
        data-cy="room-item"
        @click="goToDesks(room.id)"
      >
        <v-card-item>
          <template #prepend>
            <v-avatar color="secondary" variant="tonal" size="48">
              <v-icon size="24">$room</v-icon>
            </v-avatar>
          </template>
          <v-card-title class="text-h6">{{ room.attributes.name }}</v-card-title>
          <v-card-subtitle v-if="room.attributes.description">
            {{ room.attributes.description }}
          </v-card-subtitle>
        </v-card-item>
        <v-card-actions class="px-4 pb-4">
          <v-btn
            color="primary"
            variant="tonal"
            size="small"
            @click.stop="goToDesks(room.id)"
          >
            View Desks
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            :to="{ name: 'room-bookings', params: { roomId: room.id } }"
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
import { fetchRooms } from '../api/rooms';
import { fetchAreas } from '../api/areas';
import type { RoomAttributes } from '../api/rooms';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthStore } from '../stores/useAuthStore';
import { PageHeader, LoadingState, EmptyState } from '../components';

const authStore = useAuthStore();
const areaName = ref('');
const rooms = ref<JsonApiResource<RoomAttributes>[]>([]);
const roomsErrorMessage = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const { loading: roomsLoading, run: runRooms } = useApi();

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area' }
]);

const goToDesks = async (roomId: string) => {
  await router.push({ name: 'desks', params: { roomId } });
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
    roomsErrorMessage.value = 'Area not found.';
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

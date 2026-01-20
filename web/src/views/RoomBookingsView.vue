<template>
  <div class="page-container">
    <PageHeader
      title="Room Bookings"
      :subtitle="`Desk reservations${roomName ? ' in ' + roomName : ''}`"
      :breadcrumbs="breadcrumbs"
    />

    <!-- Date Selection -->
    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-wrap align-end ga-4">
          <DatePickerField
            v-model="selectedDate"
            label="Select Date"
            data-cy="bookings-date"
            style="max-width: 280px;"
          />
        </div>
      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="loading" type="list" :count="5" data-cy="bookings-loading" />

    <!-- Error State -->
    <v-alert v-else-if="errorMessage" type="error" class="mb-4" data-cy="bookings-error">
      {{ errorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!bookings.length"
      title="No bookings"
      message="No desks have been booked in this room for the selected date."
      icon="$calendar"
      data-cy="bookings-empty"
    />

    <!-- Bookings List -->
    <v-card v-else data-cy="bookings-list">
      <v-list lines="two">
        <v-list-item
          v-for="booking in bookings"
          :key="booking.id"
          data-cy="booking-item"
        >
          <template #prepend>
            <v-avatar color="primary" variant="tonal" size="40">
              <v-icon size="20">$desk</v-icon>
            </v-avatar>
          </template>
          <v-list-item-title class="d-flex align-center flex-wrap ga-2">
            {{ booking.attributes.desk_name }}
            <StatusChip
              v-if="booking.attributes.is_guest"
              status="guest"
              size="x-small"
            />
          </v-list-item-title>
          <v-list-item-subtitle>
            <v-icon size="14" class="mr-1">$user</v-icon>
            {{ booking.attributes.user_name || 'Unknown' }}
            <span v-if="booking.attributes.guest_name" class="ml-1">
              (Guest: {{ booking.attributes.guest_name }})
            </span>
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- Summary -->
    <div v-if="bookings.length" class="mt-4 text-body-2 text-medium-emphasis">
      {{ bookings.length }} {{ bookings.length === 1 ? 'desk' : 'desks' }} booked for {{ formattedDate }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchRoomBookings } from '../api/roomBookings';
import { fetchAreas } from '../api/areas';
import { fetchRooms } from '../api/rooms';
import type { RoomBookingAttributes } from '../api/roomBookings';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { PageHeader, LoadingState, EmptyState, DatePickerField, StatusChip } from '../components';

const bookings = ref<JsonApiResource<RoomBookingAttributes>[]>([]);
const errorMessage = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const areaName = ref('');
const roomName = ref('');
const route = useRoute();
const router = useRouter();
const { loading, run } = useApi();
const activeRoomId = ref<string | null>(null);

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area', to: '/' },
  { text: roomName.value || 'Room' },
  { text: 'Bookings' }
]);

const formattedDate = computed(() => {
  const date = new Date(selectedDate.value);
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric'
  });
});

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

const loadBookings = async (roomId: string, date: string) => {
  errorMessage.value = null;
  try {
    const resp = await run(() => fetchRoomBookings(roomId, date));
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      errorMessage.value = 'Room not found.';
      return;
    }
    errorMessage.value = 'Unable to load bookings.';
  }
};

onMounted(async () => {
  const roomId = route.params.roomId;
  if (typeof roomId !== 'string' || roomId.trim() === '') {
    errorMessage.value = 'Room not found.';
    return;
  }

  activeRoomId.value = roomId;

  // Fetch area and room names for breadcrumbs
  try {
    const areasResp = await fetchAreas();
    for (const area of areasResp.data) {
      const roomsResp = await fetchRooms(area.id);
      const room = roomsResp.data.find((r) => r.id === roomId);
      if (room) {
        areaName.value = area.attributes.name;
        roomName.value = room.attributes.name;
        break;
      }
    }
  } catch {
    // Ignore - breadcrumbs will just show generic names
  }

  await loadBookings(roomId, selectedDate.value);
});

watch(
  selectedDate,
  async (value) => {
    if (!activeRoomId.value) {
      return;
    }
    await loadBookings(activeRoomId.value, value);
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

<template>
  <div class="page-container">
    <PageHeader
      title="Today's Presence"
      :subtitle="`Who's in the office${areaName ? ' - ' + areaName : ''}`"
      :breadcrumbs="breadcrumbs"
    />

    <!-- Date Selection -->
    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-wrap align-end ga-4">
          <DatePickerField
            v-model="selectedDate"
            label="Select Date"
            data-cy="presence-date"
            style="max-width: 280px;"
          />
        </div>
      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="loading" type="list" :count="5" data-cy="presence-loading" />

    <!-- Error State -->
    <v-alert v-else-if="errorMessage" type="error" class="mb-4" data-cy="presence-error">
      {{ errorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!presence.length"
      title="No one scheduled"
      message="No one has an item booked for this date in this area."
      icon="$user"
      data-cy="presence-empty"
    />

    <!-- Presence List -->
    <v-card v-else data-cy="presence-list">
      <v-list lines="two">
        <v-list-item
          v-for="entry in presence"
          :key="entry.id"
          data-cy="presence-item"
        >
          <template #prepend>
            <v-avatar color="primary" variant="tonal" size="40">
              <span class="text-body-2 font-weight-medium">
                {{ getInitials(entry.attributes.user_name) }}
              </span>
            </v-avatar>
          </template>
          <v-list-item-title>
            {{ entry.attributes.user_name || 'Unknown' }}
          </v-list-item-title>
          <v-list-item-subtitle>
            <v-icon size="14" class="mr-1">$room</v-icon>
            {{ entry.attributes.item_group_name }}
            <span class="mx-1">&bull;</span>
            <v-icon size="14" class="mr-1">$desk</v-icon>
            {{ entry.attributes.item_name }}
          </v-list-item-subtitle>
          <div
            v-if="entry.attributes.note"
            class="d-flex align-center ga-1 mt-1 text-caption text-medium-emphasis"
            data-cy="presence-note"
          >
            <v-icon size="14">mdi-text-box-outline</v-icon>
            <span :ref="setNoteRef(entry.id)" class="note-text">{{ entry.attributes.note }}</span>
            <v-btn
              v-if="noteTruncatedMap[entry.id]"
              icon
              size="x-small"
              variant="text"
              data-cy="presence-note-expand"
              @click="expandedNote = entry.attributes.note"
            >
              <v-icon size="14">mdi-arrow-expand</v-icon>
            </v-btn>
          </div>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- Summary -->
    <div v-if="presence.length" class="mt-4 text-body-2 text-medium-emphasis">
      {{ presence.length }} {{ presence.length === 1 ? 'person' : 'people' }} scheduled for {{ formattedDate }}
    </div>

    <!-- Note expand dialog (desktop) -->
    <v-dialog v-if="!useBottomSheet" v-model="showNoteDialog" max-width="500">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="presence-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note expand bottom sheet (mobile) -->
    <v-bottom-sheet v-else v-model="showNoteDialog">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="presence-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import type { ComponentPublicInstance } from 'vue';
import { useRoute } from 'vue-router';
import { ApiError } from '../api/client';
import { fetchAreaPresence } from '../api/areaPresence';
import { fetchAreas } from '../api/areas';
import type { PresenceAttributes } from '../api/areaPresence';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { PageHeader, LoadingState, EmptyState, DatePickerField } from '../components';

const presence = ref<JsonApiResource<PresenceAttributes>[]>([]);
const errorMessage = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const areaName = ref('');
const route = useRoute();
const { loading, run } = useApi();
const { handleAuthError } = useAuthErrorHandler();
const activeAreaId = ref<string | null>(null);
const expandedNote = ref('');
const noteTruncatedMap = ref<Record<string, boolean>>({});
const noteElements = new Map<string, HTMLElement>();
const isMobile = ref(false);
const useBottomSheet = computed(() => isMobile.value);
const showNoteDialog = computed({
  get: () => expandedNote.value !== '',
  set: (v: boolean) => { if (!v) expandedNote.value = ''; }
});

const breadcrumbs = computed(() => [
  { text: 'Home', to: '/' },
  { text: areaName.value || 'Area', to: '/' },
  { text: 'Presence' }
]);

const formattedDate = computed(() => {
  const date = new Date(selectedDate.value);
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric'
  });
});

const getInitials = (name: string | undefined) => {
  if (!name) return '?';
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const loadPresence = async (areaId: string, date: string) => {
  errorMessage.value = null;
  try {
    const resp = await run(() => fetchAreaPresence(areaId, date));
    presence.value = resp.data;
    await nextTick();
    updateNoteTruncation();
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

  // Fetch area name for breadcrumb
  try {
    const areasResp = await fetchAreas();
    const area = areasResp.data.find((a) => a.id === areaId);
    if (area) {
      areaName.value = area.attributes.name;
    }
  } catch {
    // Ignore - breadcrumb will just show "Area"
  }

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

const setNoteRef = (id: string) => (el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    noteElements.set(id, el);
    return;
  }
  if (el && '$el' in el && (el.$el instanceof HTMLElement)) {
    noteElements.set(id, el.$el);
    return;
  }
  noteElements.delete(id);
};

const updateNoteTruncation = () => {
  const map: Record<string, boolean> = {};
  for (const entry of presence.value) {
    const el = noteElements.get(entry.id);
    if (el) {
      map[entry.id] = el.scrollWidth > el.clientWidth;
    }
  }
  noteTruncatedMap.value = map;
};

function updateViewport() {
  if (typeof window.matchMedia === 'function') {
    isMobile.value = window.matchMedia('(max-width: 600px)').matches;
    return;
  }
  isMobile.value = false;
}

function handleResize() {
  updateViewport();
  updateNoteTruncation();
}

onMounted(() => {
  updateViewport();
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

<style scoped>
.note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}
</style>

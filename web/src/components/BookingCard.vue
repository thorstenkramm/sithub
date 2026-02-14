<template>
  <v-card
    :class="['booking-card', { 'booking-card--cancellable': showCancel }]"
    :data-cy="dataCy"
    :data-cy-booking-id="booking.id"
  >
    <v-card-item>
      <template #prepend>
        <v-avatar :color="avatarColor" variant="tonal" size="48">
          <v-icon size="24">$desk</v-icon>
        </v-avatar>
      </template>
      <v-card-title class="d-flex align-center flex-wrap ga-2">
        {{ booking.attributes.item_name }}
        <StatusChip
          v-if="booking.attributes.is_guest"
          status="guest"
          size="x-small"
          data-cy="guest-chip"
        />
        <StatusChip
          v-else-if="booking.attributes.booked_for_me"
          status="booked-for-me"
          size="x-small"
          data-cy="booked-for-me-chip"
        />
        <StatusChip
          v-else-if="booking.attributes.booked_by_user_id && !booking.attributes.booked_for_me"
          status="on-behalf"
          size="x-small"
          data-cy="on-behalf-chip"
        />
      </v-card-title>
      <v-card-subtitle>
        {{ booking.attributes.item_group_name }} &bull; {{ booking.attributes.area_name }}
      </v-card-subtitle>
    </v-card-item>

    <v-card-text class="pt-0">
      <div class="d-flex align-center ga-2 text-body-2">
        <v-icon size="16" color="primary">$calendar</v-icon>
        <span data-cy="booking-date">{{ formattedDate }}</span>
      </div>
      <div
        v-if="booking.attributes.booked_for_me && booking.attributes.booked_by_user_name"
        class="text-caption text-medium-emphasis mt-1"
        data-cy="booked-by"
      >
        Booked by {{ booking.attributes.booked_by_user_name }}
      </div>
      <div
        v-if="booking.attributes.guest_name"
        class="text-caption text-medium-emphasis mt-1"
        data-cy="guest-name"
      >
        Guest: {{ booking.attributes.guest_name }}
      </div>

      <!-- Note display -->
      <div
        v-if="displayNote"
        class="d-flex align-center ga-1 mt-2 text-body-2 note-row"
        data-cy="booking-note"
      >
        <v-icon size="14" color="secondary">mdi-text-box-outline</v-icon>
        <span ref="noteTextEl" class="note-text">{{ displayNote }}</span>
        <v-btn
          v-if="isNoteTruncated"
          icon
          size="x-small"
          variant="text"
          data-cy="note-expand-btn"
          @click="showNoteDialog = true"
        >
          <v-icon size="14">mdi-arrow-expand</v-icon>
        </v-btn>
      </div>

      <!-- Add note link (when no note exists) -->
      <div v-if="showCancel && !displayNote" class="mt-2">
        <v-btn
          variant="text"
          size="small"
          color="primary"
          class="px-0"
          data-cy="add-note-btn"
          @click="openEditDialog"
        >
          <v-icon size="14" start>mdi-plus</v-icon>
          Add note
        </v-btn>
      </div>
    </v-card-text>

    <v-card-actions v-if="showCancel" class="px-4 pb-4">
      <v-btn
        v-if="displayNote"
        variant="text"
        size="small"
        color="secondary"
        data-cy="edit-note-btn"
        @click="openEditDialog"
      >
        <v-icon size="14" start>mdi-pencil</v-icon>
        Edit note
      </v-btn>
      <v-spacer />
      <v-btn
        color="error"
        variant="tonal"
        size="small"
        :loading="cancelling"
        :disabled="cancelling"
        data-cy="cancel-btn"
        @click="$emit('cancel', booking.id)"
      >
        Cancel Booking
      </v-btn>
    </v-card-actions>

    <!-- Note view dialog (desktop) -->
    <v-dialog v-if="!useBottomSheet" v-model="showNoteDialog" max-width="500">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="note-dialog-text">{{ displayNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note view bottom sheet (mobile) -->
    <v-bottom-sheet v-else v-model="showNoteDialog">
      <v-card>
        <v-card-title>Booking Note</v-card-title>
        <v-card-text data-cy="note-dialog-text">{{ displayNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showNoteDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>

    <!-- Note edit dialog -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card>
        <v-card-title>{{ displayNote ? 'Edit Note' : 'Add Note' }}</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="editNoteText"
            label="Note"
            :counter="500"
            :maxlength="500"
            rows="3"
            auto-grow
            data-cy="note-edit-input"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showEditDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :loading="savingNote"
            data-cy="note-save-btn"
            @click="saveNote"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, nextTick, ref, watch } from 'vue';
import type { JsonApiResource } from '../api/types';
import type { MyBookingAttributes } from '../api/bookings';
import { updateBookingNote } from '../api/bookings';
import StatusChip from './StatusChip.vue';

const props = withDefaults(defineProps<{
  booking: JsonApiResource<MyBookingAttributes>;
  showCancel?: boolean;
  cancelling?: boolean;
  dataCy?: string;
}>(), {
  showCancel: false,
  cancelling: false,
  dataCy: 'booking-card'
});

const emit = defineEmits<{
  cancel: [bookingId: string];
  'note-updated': [bookingId: string, note: string];
}>();

const noteTextEl = ref<HTMLElement | null>(null);
const showNoteDialog = ref(false);
const showEditDialog = ref(false);
const editNoteText = ref('');
const savingNote = ref(false);
const isNoteTruncated = ref(false);
const isMobile = ref(false);
const useBottomSheet = computed(() => isMobile.value);

const displayNote = computed(() => props.booking.attributes.note || '');

const formattedDate = computed(() => {
  const date = new Date(props.booking.attributes.booking_date + 'T00:00:00');
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
});

const avatarColor = computed(() => {
  if (props.booking.attributes.is_guest) return 'warning';
  if (props.booking.attributes.booked_for_me) return 'info';
  return 'primary';
});

function openEditDialog() {
  editNoteText.value = displayNote.value;
  showEditDialog.value = true;
}

function updateNoteTruncation() {
  const el = noteTextEl.value;
  if (!el) {
    isNoteTruncated.value = false;
    return;
  }
  isNoteTruncated.value = el.scrollWidth > el.clientWidth;
}

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

async function saveNote() {
  savingNote.value = true;
  try {
    await updateBookingNote(props.booking.id, editNoteText.value);
    emit('note-updated', props.booking.id, editNoteText.value);
    showEditDialog.value = false;
  } finally {
    savingNote.value = false;
  }
}

watch(displayNote, async () => {
  await nextTick();
  updateNoteTruncation();
});

onMounted(() => {
  handleResize();
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
.booking-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.booking-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}
</style>

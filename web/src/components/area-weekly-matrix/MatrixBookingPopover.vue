<template>
  <v-menu
    :model-value="modelValue"
    :activator="activatorEl"
    :close-on-content-click="false"
    location="bottom"
    max-width="340"
    data-cy="matrix-booking-popover"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card class="pa-4" data-cy="matrix-booking-card">
      <v-card-title class="text-subtitle-1 pa-0 mb-2">
        {{ item.item_name }} — {{ cell.date }}
      </v-card-title>

      <!-- Booking type radio -->
      <v-radio-group
        v-model="bookingType"
        inline
        density="compact"
        hide-details
        class="mb-2"
      >
        <v-radio
          :label="$t('items.bookForMyself')"
          value="self"
          data-cy="matrix-book-self-radio"
        />
        <v-radio
          :label="$t('items.bookForColleague')"
          value="colleague"
          data-cy="matrix-book-colleague-radio"
        />
      </v-radio-group>

      <!-- Colleague picker -->
      <v-expand-transition>
        <div v-if="bookingType === 'colleague'" class="mt-2">
          <v-autocomplete
            v-model="selectedColleagueId"
            :items="colleagueList"
            item-title="displayName"
            item-value="id"
            :label="$t('items.selectColleague')"
            density="compact"
            :loading="colleaguesLoading"
            clearable
            hide-details
            data-cy="matrix-colleague-select"
          />
        </div>
      </v-expand-transition>

      <!-- Note field -->
      <v-text-field
        v-model="noteText"
        :label="$t('items.noteLabel')"
        density="compact"
        hide-details
        class="mt-3"
        data-cy="matrix-booking-note"
      />

      <!-- Error message -->
      <v-alert
        v-if="errorMessage"
        type="error"
        density="compact"
        class="mt-3"
        data-cy="matrix-booking-error"
      >
        {{ errorMessage }}
      </v-alert>

      <v-card-actions class="pa-0 mt-3">
        <v-spacer />
        <v-btn
          variant="text"
          size="small"
          data-cy="matrix-booking-cancel"
          @click="$emit('update:modelValue', false)"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          size="small"
          :loading="submitting"
          data-cy="matrix-booking-confirm"
          @click="submitBooking"
        >
          {{ $t('items.book') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>

  <WarningConfirmDialog
    v-model="warningShow"
    v-model:dont-show-again="warningDontShowAgain"
    :item-name="warningItemName"
    :message="warningMessage"
    @confirm="warningConfirmAction"
    @cancel="warningCancelAction"
  />
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import type { MatrixCell, MatrixItem } from '../../api/itemGroupMatrix';
import { createBooking, type BookOnBehalfOptions } from '../../api/bookings';
import { fetchColleagues } from '../../api/users';
import { ApiError } from '../../api/client';
import { useWarningConfirmation } from '../../composables/useWarningConfirmation';
import { getSafeLocalStorage } from '../../composables/storage';
import WarningConfirmDialog from '../WarningConfirmDialog.vue';

const LAST_COLLEAGUE_KEY = 'sithub_matrix_last_colleague';

const props = defineProps<{
  modelValue: boolean;
  activatorEl: HTMLElement | undefined;
  item: MatrixItem;
  cell: MatrixCell;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  booked: [];
  bookingConflict: [];
}>();

const { t } = useI18n();
const {
  show: warningShow,
  itemName: warningItemName,
  message: warningMessage,
  dontShowAgain: warningDontShowAgain,
  present: presentWarnings,
  confirm: warningConfirmAction,
  cancel: warningCancelAction,
} = useWarningConfirmation();

const bookingType = ref<'self' | 'colleague'>('self');
const selectedColleagueId = ref<string | null>(null);
const colleagueList = ref<Array<{ id: string; displayName: string }>>([]);
const colleaguesLoading = ref(false);
const noteText = ref('');
const errorMessage = ref<string | null>(null);
const submitting = ref(false);

function loadLastColleague() {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  try {
    const raw = storage.getItem(LAST_COLLEAGUE_KEY);
    if (raw) {
      selectedColleagueId.value = raw;
    }
  } catch {
    // Ignore
  }
}

function saveLastColleague(id: string) {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  try {
    storage.setItem(LAST_COLLEAGUE_KEY, id);
  } catch {
    // Storage full
  }
}

async function loadColleagues() {
  if (colleagueList.value.length > 0) return;
  colleaguesLoading.value = true;
  try {
    const resp = await fetchColleagues();
    colleagueList.value = resp.data.map(r => ({
      id: r.id,
      displayName: r.attributes.display_name
    }));
  } catch {
    colleagueList.value = [];
  } finally {
    colleaguesLoading.value = false;
  }
}

// Reset state when the popover opens. The warning confirmation is an
// independent persistent modal controlled solely by its own CANCEL/CONFIRM —
// we must NOT abort it when the menu closes. Interacting with the dialog (e.g.
// ticking "Don't show again", whose checkbox click propagates to the document)
// closes the underlying v-menu; tying an abort to that close would wrongly
// cancel the booking mid-confirmation.
watch(() => props.modelValue, (open) => {
  if (open) {
    bookingType.value = 'self';
    noteText.value = '';
    errorMessage.value = null;
    submitting.value = false;
  }
});

// Load colleagues when switching to colleague mode
watch(bookingType, (type) => {
  if (type === 'colleague') {
    loadColleagues();
    loadLastColleague();
  }
});

function resolveColleagueName(id: string): string {
  return colleagueList.value.find(c => c.id === id)?.displayName ?? '';
}

function submitBooking() {
  errorMessage.value = null;

  if (bookingType.value === 'colleague' && !selectedColleagueId.value) {
    errorMessage.value = t('items.selectColleagueError');
    return;
  }

  // Uniform pre-booking warning confirmation (skipped/suppressed as needed).
  const warnItems = props.item.warning
    ? [{ itemId: props.item.item_id, itemName: props.item.item_name, warning: props.item.warning }]
    : [];
  presentWarnings(warnItems, () => void doBook());
}

async function doBook() {
  submitting.value = true;
  try {
    const onBehalf: BookOnBehalfOptions | undefined =
      bookingType.value === 'colleague' && selectedColleagueId.value
        ? { forUserId: selectedColleagueId.value, forUserName: resolveColleagueName(selectedColleagueId.value) }
        : undefined;

    const note = noteText.value.trim() || undefined;

    await createBooking(props.item.item_id, props.cell.date, onBehalf, undefined, note);

    if (bookingType.value === 'colleague' && selectedColleagueId.value) {
      saveLastColleague(selectedColleagueId.value);
    }

    emit('update:modelValue', false);
    emit('booked');
  } catch (err) {
    if (err instanceof ApiError && err.status === 409) {
      errorMessage.value = t('matrix.deskNoLongerAvailable');
      emit('bookingConflict');
    } else {
      errorMessage.value = t('items.unableToBook');
    }
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  loadColleagues();
});
</script>

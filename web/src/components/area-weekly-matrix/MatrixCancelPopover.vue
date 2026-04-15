<template>
  <v-menu
    :model-value="modelValue"
    :activator="activatorEl"
    :close-on-content-click="false"
    location="bottom"
    max-width="300"
    data-cy="matrix-cancel-popover"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card class="pa-4" data-cy="matrix-cancel-card">
      <div class="text-body-2 mb-1" data-cy="matrix-cancel-person">
        <strong>{{ $t('matrix.person') }}:</strong> {{ cell.booker_name }}
      </div>
      <div class="text-body-2 mb-1" data-cy="matrix-cancel-desk">
        <strong>{{ $t('matrix.desk') }}:</strong> {{ item.item_name }}
      </div>
      <div class="text-body-2 mb-3" data-cy="matrix-cancel-date">
        <strong>{{ $t('common.date') }}:</strong> {{ cell.date }}
      </div>

      <!-- Error message -->
      <v-alert
        v-if="errorMessage"
        type="error"
        density="compact"
        class="mb-3"
        data-cy="matrix-cancel-error"
      >
        {{ errorMessage }}
      </v-alert>

      <v-card-actions class="pa-0">
        <v-spacer />
        <v-btn
          variant="text"
          size="small"
          data-cy="matrix-cancel-close"
          @click="$emit('update:modelValue', false)"
        >
          {{ $t('common.close') }}
        </v-btn>
        <v-btn
          color="error"
          variant="flat"
          size="small"
          :loading="submitting"
          data-cy="matrix-cancel-confirm"
          @click="submitCancel"
        >
          {{ $t('matrix.cancelBooking') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { MatrixCell, MatrixItem } from '../../api/itemGroupMatrix';
import { cancelBooking } from '../../api/bookings';

const props = defineProps<{
  modelValue: boolean;
  activatorEl: HTMLElement | undefined;
  item: MatrixItem;
  cell: MatrixCell;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  cancelled: [];
}>();

const { t } = useI18n();
const errorMessage = ref<string | null>(null);
const submitting = ref(false);

watch(() => props.modelValue, (open) => {
  if (open) {
    errorMessage.value = null;
    submitting.value = false;
  }
});

async function submitCancel() {
  if (!props.cell.booking_id) return;

  errorMessage.value = null;
  submitting.value = true;
  try {
    await cancelBooking(props.cell.booking_id);
    emit('update:modelValue', false);
    emit('cancelled');
  } catch {
    errorMessage.value = t('matrix.cancelFailed');
  } finally {
    submitting.value = false;
  }
}
</script>

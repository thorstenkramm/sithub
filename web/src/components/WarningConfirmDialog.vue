<script setup lang="ts">
/**
 * WarningConfirmDialog is the single, uniform pre-booking warning confirmation
 * shown from every booking surface (tiles, floor plan, weekly table). The queue
 * and suppression logic live in useWarningConfirmation; this component is
 * presentation + user actions only.
 */
import { useI18n } from 'vue-i18n';
import ItemWarning from './ItemWarning.vue';

defineProps<{
  modelValue: boolean;
  itemName: string;
  message: string;
  dontShowAgain: boolean;
}>();

defineEmits<{
  'update:modelValue': [value: boolean];
  'update:dontShowAgain': [value: boolean];
  confirm: [];
  cancel: [];
}>();

useI18n(); // expose $t in template
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="400"
    persistent
    data-cy="warning-dialog"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>{{ $t('items.warningDialogTitle') }}</v-card-title>
      <v-card-text>
        <div
          data-cy="warning-item-name"
          class="text-subtitle-2 mb-2"
          style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap"
        >
          {{ itemName }}
        </div>
        <ItemWarning mode="inline" :show-icon="false" :warning="message" data-cy="warning-message" />
      </v-card-text>
      <v-card-actions class="flex-column align-start px-4 pb-4">
        <v-checkbox
          :model-value="dontShowAgain"
          :label="$t('items.warningDontShowAgain')"
          density="compact"
          hide-details
          data-cy="warning-dont-show-checkbox"
          class="mb-2"
          @update:model-value="$emit('update:dontShowAgain', $event as boolean)"
        />
        <div class="d-flex w-100 justify-end ga-2">
          <v-btn variant="text" data-cy="warning-cancel-btn" @click.stop="$emit('cancel')">
            {{ $t('items.warningCancel') }}
          </v-btn>
          <v-btn color="primary" variant="flat" data-cy="warning-confirm-btn" @click.stop="$emit('confirm')">
            {{ $t('items.warningConfirm') }}
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

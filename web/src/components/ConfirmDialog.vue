<template>
  <v-dialog
    :model-value="modelValue"
    max-width="400"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6 font-weight-bold">
        {{ title }}
      </v-card-title>
      <v-card-text class="text-body-1">
        {{ message }}
      </v-card-text>
      <v-card-actions class="pa-4 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          :disabled="loading"
          @click="handleCancel"
          data-cy="confirm-dialog-cancel"
        >
          {{ resolvedCancelText }}
        </v-btn>
        <v-btn
          :color="confirmColor"
          :loading="loading"
          variant="flat"
          @click="handleConfirm"
          data-cy="confirm-dialog-confirm"
        >
          {{ resolvedConfirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = withDefaults(defineProps<{
  modelValue: boolean;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  confirmColor?: string;
  loading?: boolean;
}>(), {
  confirmText: '',
  cancelText: '',
  confirmColor: 'primary',
  loading: false
});

const resolvedConfirmText = computed(() => props.confirmText || t('common.confirm'));
const resolvedCancelText = computed(() => props.cancelText || t('common.cancel'));

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  confirm: [];
  cancel: [];
}>();

function handleConfirm() {
  emit('confirm');
}

function handleCancel() {
  emit('update:modelValue', false);
  emit('cancel');
}
</script>

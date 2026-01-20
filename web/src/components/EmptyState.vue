<template>
  <div class="empty-state text-center py-12 px-4">
    <v-icon
      :icon="icon"
      size="64"
      class="text-disabled mb-4"
    />
    <h3 class="text-h6 font-weight-medium mb-2">{{ title }}</h3>
    <p v-if="message" class="text-body-2 text-medium-emphasis mb-4 mx-auto" style="max-width: 400px;">
      {{ message }}
    </p>
    <v-btn
      v-if="actionText"
      :to="actionTo"
      color="primary"
      variant="flat"
      @click="handleAction"
      data-cy="empty-state-action"
    >
      {{ actionText }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
import type { RouteLocationRaw } from 'vue-router';

const props = defineProps<{
  title: string;
  message?: string;
  icon?: string;
  actionText?: string;
  actionTo?: RouteLocationRaw;
}>();

const emit = defineEmits<{
  action: [];
}>();

// Default icon
const icon = props.icon || '$info';

function handleAction() {
  if (!props.actionTo) {
    emit('action');
  }
}
</script>

<style scoped>
.empty-state {
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>

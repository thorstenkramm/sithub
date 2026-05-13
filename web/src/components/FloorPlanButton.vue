<template>
  <v-btn
    variant="outlined"
    density="default"
    prepend-icon="$map"
    class="floor-plan-button"
    v-bind="$attrs"
  >
    {{ label || $t('itemGroups.floorPlan') }}
  </v-btn>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

defineProps<{
  /** Optional override for the button label. Defaults to the shared "Floor plan" i18n key. */
  label?: string;
}>();

useI18n(); // expose $t in template
</script>

<style scoped>
.floor-plan-button {
  text-transform: none;
}
</style>

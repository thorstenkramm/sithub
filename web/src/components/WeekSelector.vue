<template>
  <div class="week-selector d-flex align-center ga-1">
    <v-btn
      icon
      variant="text"
      size="small"
      :disabled="prevDisabled"
      :aria-label="t('weekSelector.previousWeek')"
      data-cy="week-prev"
      @click="step(-1)"
    >
      <v-icon>$chevronLeft</v-icon>
    </v-btn>
    <v-select
      v-model="selectedWeek"
      :items="options"
      item-title="label"
      item-value="value"
      :label="label"
      density="compact"
      hide-details
      :data-cy="selectDataCy"
      class="week-selector-select"
    />
    <v-btn
      icon
      variant="text"
      size="small"
      :disabled="nextDisabled"
      :aria-label="t('weekSelector.nextWeek')"
      data-cy="week-next"
      @click="step(1)"
    >
      <v-icon>$chevronRight</v-icon>
    </v-btn>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import type { WeekOption } from '../composables/useWeekSelector';

// Week dropdown with previous/next arrows. The options list is the single source of the
// allowed range: the arrows only move between neighbouring options and never compute dates.
const props = withDefaults(defineProps<{
  options: WeekOption[];
  label: string;
  selectDataCy?: string;
}>(), {
  selectDataCy: 'week-selector'
});

const selectedWeek = defineModel<string>({ required: true });

const { t } = useI18n();

const selectedIndex = computed(() => props.options.findIndex(option => option.value === selectedWeek.value));
const prevDisabled = computed(() => selectedIndex.value <= 0);
const nextDisabled = computed(() => selectedIndex.value < 0 || selectedIndex.value >= props.options.length - 1);

function step(offset: -1 | 1) {
  const target = props.options[selectedIndex.value + offset];
  if (selectedIndex.value >= 0 && target) {
    selectedWeek.value = target.value;
  }
}
</script>

<style scoped>
.week-selector-select {
  flex: 1 1 auto;
  min-width: 0;
}
</style>

<template>
  <v-menu
    v-model="menuOpen"
    :close-on-content-click="false"
    location="bottom start"
  >
    <template #activator="{ props: menuProps }">
      <v-text-field
        v-bind="menuProps"
        :model-value="displayValue"
        :label="label"
        :disabled="disabled"
        readonly
        prepend-inner-icon="$calendar"
        :data-cy="dataCy"
      />
    </template>
    <v-date-picker
      :model-value="internalDate"
      :min="min"
      :max="max"
      @update:model-value="handleDateChange"
    />
  </v-menu>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = withDefaults(defineProps<{
  modelValue: string; // YYYY-MM-DD format
  label?: string;
  min?: string;
  max?: string;
  disabled?: boolean;
  dataCy?: string;
}>(), {
  label: 'Date'
});

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const menuOpen = ref(false);

// Convert string to Date for v-date-picker
const internalDate = computed(() => {
  if (!props.modelValue) return null;
  return new Date(props.modelValue);
});

// Format for display
const displayValue = computed(() => {
  if (!props.modelValue) return '';
  const date = new Date(props.modelValue);
  return date.toLocaleDateString(undefined, {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
});

function handleDateChange(date: unknown) {
  if (date instanceof Date) {
    // Format as YYYY-MM-DD
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    emit('update:modelValue', `${year}-${month}-${day}`);
    menuOpen.value = false;
  }
}

// Close menu when value changes externally
watch(() => props.modelValue, () => {
  menuOpen.value = false;
});
</script>

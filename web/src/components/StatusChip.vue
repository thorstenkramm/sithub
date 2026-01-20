<template>
  <v-chip
    :color="chipColor"
    :size="size"
    variant="tonal"
    :prepend-icon="chipIcon"
  >
    {{ displayLabel }}
  </v-chip>
</template>

<script setup lang="ts">
import { computed } from 'vue';

type StatusType = 'available' | 'booked' | 'mine' | 'unavailable' | 'guest' | 'pending' | 'booked-for-me' | 'on-behalf';

const props = withDefaults(defineProps<{
  status: StatusType;
  label?: string;
  size?: 'x-small' | 'small' | 'default';
}>(), {
  size: 'small'
});

const statusConfig: Record<StatusType, { color: string; label: string; icon: string }> = {
  available: { color: 'success', label: 'Available', icon: '$success' },
  booked: { color: 'warning', label: 'Booked', icon: '$calendar' },
  mine: { color: 'primary', label: 'My Booking', icon: '$check' },
  unavailable: { color: 'error', label: 'Unavailable', icon: '$close' },
  guest: { color: 'warning', label: 'Guest', icon: '$userPlus' },
  pending: { color: 'warning', label: 'Pending', icon: '$clock' },
  'booked-for-me': { color: 'info', label: 'Booked for you', icon: '$userPlus' },
  'on-behalf': { color: 'secondary', label: 'On behalf', icon: '$user' }
};

const chipColor = computed(() => statusConfig[props.status]?.color || 'default');
const chipIcon = computed(() => statusConfig[props.status]?.icon);
const displayLabel = computed(() => props.label || statusConfig[props.status]?.label || props.status);
</script>

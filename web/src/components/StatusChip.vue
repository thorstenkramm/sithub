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
import { useI18n } from 'vue-i18n';

type StatusType = 'available' | 'booked' | 'mine' | 'unavailable' | 'guest' | 'pending' | 'booked-for-me' | 'on-behalf';

const { t } = useI18n();

const props = withDefaults(defineProps<{
  status: StatusType;
  label?: string;
  size?: 'x-small' | 'small' | 'default';
}>(), {
  size: 'small'
});

const statusConfig: Record<StatusType, { color: string; labelKey: string; icon: string }> = {
  available: { color: 'success', labelKey: 'status.available', icon: '$success' },
  booked: { color: 'warning', labelKey: 'status.booked', icon: '$calendar' },
  mine: { color: 'primary', labelKey: 'status.mine', icon: '$check' },
  unavailable: { color: 'error', labelKey: 'status.unavailable', icon: '$close' },
  guest: { color: 'warning', labelKey: 'status.guest', icon: '$userPlus' },
  pending: { color: 'warning', labelKey: 'status.pending', icon: '$clock' },
  'booked-for-me': { color: 'info', labelKey: 'status.bookedForMe', icon: '$userPlus' },
  'on-behalf': { color: 'secondary', labelKey: 'status.onBehalf', icon: '$user' }
};

const chipColor = computed(() => statusConfig[props.status]?.color || 'default');
const chipIcon = computed(() => statusConfig[props.status]?.icon);
const displayLabel = computed(() => props.label || t(statusConfig[props.status]?.labelKey) || props.status);
</script>

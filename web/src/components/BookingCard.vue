<template>
  <v-card
    :class="['booking-card', { 'booking-card--cancellable': showCancel }]"
    :data-cy="dataCy"
    :data-cy-booking-id="booking.id"
  >
    <v-card-item>
      <template #prepend>
        <v-avatar :color="avatarColor" variant="tonal" size="48">
          <v-icon size="24">$desk</v-icon>
        </v-avatar>
      </template>
      <v-card-title class="d-flex align-center flex-wrap ga-2">
        {{ booking.attributes.item_name }}
        <StatusChip
          v-if="booking.attributes.is_guest"
          status="guest"
          size="x-small"
          data-cy="guest-chip"
        />
        <StatusChip
          v-else-if="booking.attributes.booked_for_me"
          status="booked-for-me"
          size="x-small"
          data-cy="booked-for-me-chip"
        />
        <StatusChip
          v-else-if="booking.attributes.booked_by_user_id && !booking.attributes.booked_for_me"
          status="on-behalf"
          size="x-small"
          data-cy="on-behalf-chip"
        />
      </v-card-title>
      <v-card-subtitle>
        {{ booking.attributes.item_group_name }} &bull; {{ booking.attributes.area_name }}
      </v-card-subtitle>
    </v-card-item>

    <v-card-text class="pt-0">
      <div class="d-flex align-center ga-2 text-body-2">
        <v-icon size="16" color="primary">$calendar</v-icon>
        <span data-cy="booking-date">{{ formattedDate }}</span>
      </div>
      <div
        v-if="booking.attributes.booked_for_me && booking.attributes.booked_by_user_name"
        class="text-caption text-medium-emphasis mt-1"
        data-cy="booked-by"
      >
        Booked by {{ booking.attributes.booked_by_user_name }}
      </div>
      <div
        v-if="booking.attributes.guest_name"
        class="text-caption text-medium-emphasis mt-1"
        data-cy="guest-name"
      >
        Guest: {{ booking.attributes.guest_name }}
      </div>
    </v-card-text>

    <v-card-actions v-if="showCancel" class="px-4 pb-4">
      <v-spacer />
      <v-btn
        color="error"
        variant="tonal"
        size="small"
        :loading="cancelling"
        :disabled="cancelling"
        data-cy="cancel-btn"
        @click="$emit('cancel', booking.id)"
      >
        Cancel Booking
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { JsonApiResource } from '../api/types';
import type { MyBookingAttributes } from '../api/bookings';
import StatusChip from './StatusChip.vue';

const props = withDefaults(defineProps<{
  booking: JsonApiResource<MyBookingAttributes>;
  showCancel?: boolean;
  cancelling?: boolean;
  dataCy?: string;
}>(), {
  showCancel: false,
  cancelling: false,
  dataCy: 'booking-card'
});

defineEmits<{
  cancel: [bookingId: string];
}>();

const formattedDate = computed(() => {
  const date = new Date(props.booking.attributes.booking_date + 'T00:00:00');
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
});

const avatarColor = computed(() => {
  if (props.booking.attributes.is_guest) return 'warning';
  if (props.booking.attributes.booked_for_me) return 'info';
  return 'primary';
});
</script>

<style scoped>
.booking-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.booking-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>

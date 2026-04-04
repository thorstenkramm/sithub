<template>
  <div class="page-container">
    <PageHeader
      :title="$t('itemGroupBookings.title')"
      :subtitle="itemGroupName ? $t('itemGroupBookings.reservations') + ' in ' + itemGroupName : $t('itemGroupBookings.reservations')"
      :breadcrumbs="breadcrumbs"
    />

    <!-- Date Selection -->
    <v-card class="mb-6">
      <v-card-text>
        <div class="d-flex flex-wrap align-end ga-4">
          <DatePickerField
            v-model="selectedDate"
            :label="$t('itemGroupBookings.selectDate')"
            data-cy="bookings-date"
            style="max-width: 280px;"
          />
        </div>
      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="loading" type="list" :count="5" data-cy="bookings-loading" />

    <!-- Error State -->
    <v-alert v-else-if="errorMessage" type="error" class="mb-4" data-cy="bookings-error">
      {{ errorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="!bookings.length"
      :title="$t('itemGroupBookings.emptyTitle')"
      :message="$t('itemGroupBookings.emptyMessage')"
      icon="$calendar"
      data-cy="bookings-empty"
    />

    <!-- Bookings List -->
    <v-card v-else data-cy="bookings-list">
      <v-list lines="two">
        <v-list-item
          v-for="booking in bookings"
          :key="booking.id"
          data-cy="booking-item"
        >
          <template #prepend>
            <v-avatar color="primary" variant="tonal" size="40">
              <v-icon size="20">$desk</v-icon>
            </v-avatar>
          </template>
          <v-list-item-title>
            {{ booking.attributes.item_name }}
          </v-list-item-title>
          <v-list-item-subtitle>
            <v-icon size="14" class="mr-1">$user</v-icon>
            {{ booking.attributes.user_name || $t('itemGroupBookings.unknown') }}
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- Summary -->
    <div v-if="bookings.length" class="mt-4 text-body-2 text-medium-emphasis">
      {{ $t('itemGroupBookings.itemCount', { count: bookings.length, date: formattedDate }, bookings.length) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { fetchItemGroupBookings } from '../api/itemGroupBookings';
import { fetchAreas } from '../api/areas';
import { fetchItemGroups } from '../api/itemGroups';
import type { ItemGroupBookingAttributes } from '../api/itemGroupBookings';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { PageHeader, LoadingState, EmptyState, DatePickerField } from '../components';

const bookings = ref<JsonApiResource<ItemGroupBookingAttributes>[]>([]);
const errorMessage = ref<string | null>(null);
const selectedDate = ref(formatDate(new Date()));
const areaName = ref('');
const itemGroupName = ref('');
const route = useRoute();
const { t, locale } = useI18n();
const { loading, run } = useApi();
const { handleAuthError } = useAuthErrorHandler();
const activeItemGroupId = ref<string | null>(null);

const queryAreaId = computed(() => {
  const value = route.query.areaId;
  return typeof value === 'string' ? value : undefined;
});
const resolvedAreaId = ref<string | null>(null);
const breadcrumbAreaId = computed(() =>
  resolvedAreaId.value ? resolvedAreaId.value : areaName.value ? undefined : queryAreaId.value
);

const breadcrumbs = computed(() => [
  { text: t('common.home'), to: '/' },
  {
    text: areaName.value || t('common.area'),
    to: breadcrumbAreaId.value ? `/areas/${breadcrumbAreaId.value}/item-groups` : undefined
  },
  {
    text: itemGroupName.value || t('common.itemGroup'),
    to: activeItemGroupId.value
      ? {
        name: 'items' as const,
        params: { itemGroupId: activeItemGroupId.value },
        query: breadcrumbAreaId.value ? { areaId: breadcrumbAreaId.value } : {}
      }
      : undefined
  },
  { text: t('common.bookings') }
]);

const formattedDate = computed(() => {
  const date = new Date(selectedDate.value);
  return date.toLocaleDateString(locale.value || undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric'
  });
});

const loadBookings = async (itemGroupId: string, date: string) => {
  errorMessage.value = null;
  try {
    const resp = await run(() => fetchItemGroupBookings(itemGroupId, date));
    bookings.value = resp.data;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (isConnectionError(err)) {
      errorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      errorMessage.value = t('itemGroupBookings.notFound');
      return;
    }
    errorMessage.value = t('itemGroupBookings.unableToLoad');
  }
};

onMounted(async () => {
  const itemGroupId = route.params.itemGroupId;
  if (typeof itemGroupId !== 'string' || itemGroupId.trim() === '') {
    errorMessage.value = t('itemGroupBookings.notFound');
    return;
  }

  activeItemGroupId.value = itemGroupId;

  // Fetch area and item group names for breadcrumbs
  try {
    const areasResp = await fetchAreas();
    for (const area of areasResp.data) {
      const igResp = await fetchItemGroups(area.id);
      const ig = igResp.data.find((ig) => ig.id === itemGroupId);
      if (ig) {
        areaName.value = area.attributes.name;
        itemGroupName.value = ig.attributes.name;
        resolvedAreaId.value = area.id;
        break;
      }
    }
  } catch (err) {
    if (isConnectionError(err)) {
      errorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    // Ignore other errors - breadcrumbs will just show generic names
  }

  await loadBookings(itemGroupId, selectedDate.value);
});

watch(
  selectedDate,
  async (value) => {
    if (!activeItemGroupId.value) {
      return;
    }
    await loadBookings(activeItemGroupId.value, value);
  },
  { flush: 'post' }
);

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

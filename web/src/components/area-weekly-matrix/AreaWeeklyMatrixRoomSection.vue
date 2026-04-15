<template>
  <!-- Room header row -->
  <tr class="matrix-room-header" :data-cy="`matrix-room-${group.id}`">
    <td
      class="room-header-cell sticky-col"
      :colspan="1"
    >
      <div class="d-flex align-center ga-1">
        <v-btn
          icon
          variant="text"
          size="x-small"
          :data-cy="`matrix-room-toggle-${group.id}`"
          :aria-label="collapsed ? 'Expand' : 'Collapse'"
          @click.stop="$emit('toggleCollapse')"
        >
          <v-icon size="18">{{ collapsed ? '$chevronRight' : '$chevronDown' }}</v-icon>
        </v-btn>
        <span class="room-header-name text-subtitle-2">{{ group.attributes.item_group_name }}</span>
      </div>
    </td>
    <!-- Collapsed summary: occupied counts per day -->
    <td
      v-for="day in days"
      :key="`summary-${day.date}`"
      class="room-summary-cell text-center text-caption"
      :class="{ 'matrix-past-day': isPastDay(day.date) }"
    >
      <span v-if="collapsed" class="text-medium-emphasis" :data-cy="`matrix-room-summary-${group.id}-${day.weekday}`">
        {{ occupiedCount(day.date) }}/{{ group.attributes.items.length }}
      </span>
    </td>
  </tr>

  <!-- Desk rows (visible when expanded) -->
  <template v-if="!collapsed">
    <AreaWeeklyMatrixRow
      v-for="item in group.attributes.items"
      :key="item.item_id"
      :item="item"
      :days="days"
      :current-user-id="currentUserId"
      :is-admin="isAdmin"
      :today="today"
    />
  </template>
</template>

<script setup lang="ts">
import type { ItemGroupMatrixAttributes, MatrixDayMeta } from '../../api/itemGroupMatrix';
import type { JsonApiResource } from '../../api/types';
import AreaWeeklyMatrixRow from './AreaWeeklyMatrixRow.vue';

const props = defineProps<{
  group: JsonApiResource<ItemGroupMatrixAttributes>;
  days: MatrixDayMeta[];
  collapsed: boolean;
  currentUserId: string;
  isAdmin: boolean;
  today: string;
}>();

defineEmits<{
  toggleCollapse: [];
}>();

function isPastDay(dateStr: string): boolean {
  return dateStr < props.today;
}

function occupiedCount(date: string): number {
  return props.group.attributes.items.reduce((count, item) => {
    const cell = item.cells.find(c => c.date === date);
    return count + (cell && cell.availability === 'occupied' ? 1 : 0);
  }, 0);
}
</script>

<style scoped>
.matrix-room-header {
  background: rgba(var(--v-theme-primary), 0.06);
}

.room-header-cell {
  padding: 6px 8px;
  font-weight: 500;
  background: rgba(var(--v-theme-primary), 0.06);
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.room-header-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.room-summary-cell {
  padding: 6px 4px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background: rgba(var(--v-theme-primary), 0.06);
}

.matrix-past-day {
  opacity: 0.5;
}
</style>

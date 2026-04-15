<template>
  <tr class="matrix-desk-row" :data-cy="`matrix-row-${item.item_id}`">
    <!-- Sticky desk name column -->
    <td class="matrix-desk-name sticky-col">
      <span class="desk-label" :data-cy="`matrix-desk-label-${item.item_id}`">
        {{ item.item_name }}

        <v-tooltip v-if="item.equipment.length > 0" location="right">
          <template #activator="{ props: eqProps }">
            <v-icon v-bind="eqProps" size="14" class="ml-1" data-cy="matrix-equipment-icon">$equipment</v-icon>
          </template>
          <span data-cy="matrix-equipment-tooltip">{{ item.equipment.join(', ') }}</span>
        </v-tooltip>

        <v-tooltip v-if="item.warning" location="right">
          <template #activator="{ props: warnProps }">
            <v-icon v-bind="warnProps" size="14" color="warning" class="ml-1" data-cy="matrix-warning-icon">$warning</v-icon>
          </template>
          <span data-cy="matrix-warning-tooltip">{{ item.warning }}</span>
        </v-tooltip>
      </span>
    </td>

    <!-- Day cells -->
    <AreaWeeklyMatrixCell
      v-for="cell in item.cells"
      :key="cell.date"
      :cell="cell"
      :item="item"
      :reserved="item.reserved ?? false"
      :is-past="cell.date < today"
      :current-user-id="currentUserId"
      :is-admin="isAdmin"
    />
  </tr>
</template>

<script setup lang="ts">
import type { MatrixItem, MatrixDayMeta } from '../../api/itemGroupMatrix';
import AreaWeeklyMatrixCell from './AreaWeeklyMatrixCell.vue';

defineProps<{
  item: MatrixItem;
  days: MatrixDayMeta[];
  currentUserId: string;
  isAdmin: boolean;
  today: string;
}>();
</script>

<style scoped>
.matrix-desk-row:hover {
  background: rgba(var(--v-theme-on-surface), 0.04);
}

.matrix-desk-name {
  padding: 6px 8px 6px 32px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  border-bottom: 1px solid rgba(var(--v-border-color), 0.08);
  font-size: 0.85rem;
}

.desk-label {
  cursor: default;
  display: inline-flex;
  align-items: center;
}
</style>

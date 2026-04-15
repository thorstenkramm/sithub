<template>
  <td
    class="matrix-cell"
    :class="cellClasses"
    :data-cy="`matrix-cell-${cell.date}`"
  >
    <!-- Past day: muted, non-interactive -->
    <template v-if="isPast">
      <div v-if="cell.availability === 'occupied'" class="cell-content cell-occupied-past">
        <span class="cell-initials text-caption">{{ initials }}</span>
      </div>
      <div v-else class="cell-content cell-muted" />
    </template>

    <!-- Locked free cell (reserved, user cannot book) -->
    <template v-else-if="cell.availability === 'free' && reserved">
      <div class="cell-content cell-locked" data-cy="matrix-cell-locked">
        <v-icon size="14" color="grey">$lock</v-icon>
      </div>
    </template>

    <!-- Free bookable cell -->
    <template v-else-if="cell.availability === 'free'">
      <div
        class="cell-content cell-free"
        data-cy="matrix-cell-free"
        @click="handleFreeClick"
      />
    </template>

    <!-- Occupied cell -->
    <template v-else>
      <v-tooltip location="top">
        <template #activator="{ props: tooltipProps }">
          <div
            v-bind="tooltipProps"
            class="cell-content cell-occupied"
            :class="{
              'cell-booked-by-me': cell.booked_by_me,
              'cell-inert': !canInteract,
              'cell-interactive': canInteract
            }"
            data-cy="matrix-cell-occupied"
            @click="handleOccupiedClick"
          >
            <v-avatar
              v-if="cell.booker_user_id && !avatarFailed"
              size="24"
              class="cell-avatar"
            >
              <v-img
                :src="avatarUrl"
                :alt="cell.booker_name"
                @error="avatarFailed = true"
              />
            </v-avatar>
            <span class="cell-short-name text-caption" data-cy="matrix-cell-initials">{{ shortName }}</span>
          </div>
        </template>
        <span data-cy="matrix-cell-tooltip">{{ cell.booker_name }}</span>
      </v-tooltip>
    </template>
  </td>
</template>

<script setup lang="ts">
import { computed, inject, ref } from 'vue';
import type { MatrixCell, MatrixItem } from '../../api/itemGroupMatrix';
import { getAvatarUrl } from '../../api/avatars';
import { getInitials, getShortName } from '../../utils/text';
import type { MatrixCellClickHandler } from './matrixTypes';

const props = defineProps<{
  cell: MatrixCell;
  item: MatrixItem;
  reserved: boolean;
  isPast: boolean;
  currentUserId: string;
  isAdmin: boolean;
}>();

const onCellClick = inject<MatrixCellClickHandler>('matrixCellClick');

const avatarFailed = ref(false);

const initials = computed(() => getInitials(props.cell.booker_name));
const shortName = computed(() => getShortName(props.cell.booker_name));
const avatarUrl = computed(() =>
  props.cell.booker_user_id ? getAvatarUrl(props.cell.booker_user_id) : ''
);

const canInteract = computed(() => {
  if (props.cell.availability !== 'occupied') return false;
  return props.cell.booked_by_me || props.isAdmin;
});

const cellClasses = computed(() => ({
  'matrix-past-day': props.isPast,
  'cell-non-interactive': props.isPast || (props.cell.availability === 'free' && props.reserved)
}));

function handleFreeClick(event: MouseEvent) {
  const el = (event.currentTarget as HTMLElement).closest('td') as HTMLElement;
  onCellClick?.({ type: 'book', el, item: props.item, cell: props.cell });
}

function handleOccupiedClick(event: MouseEvent) {
  if (!canInteract.value) return;
  const el = (event.currentTarget as HTMLElement).closest('td') as HTMLElement;
  onCellClick?.({ type: 'cancel', el, item: props.item, cell: props.cell });
}
</script>

<style scoped>
.matrix-cell {
  text-align: center;
  vertical-align: middle;
  padding: 4px;
  min-width: 80px;
  height: 40px;
  border-bottom: 1px solid rgba(var(--v-border-color), 0.08);
  border-left: 1px solid rgba(var(--v-border-color), 0.06);
}

.cell-content {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 28px;
  border-radius: 4px;
  padding: 2px 4px;
}

.cell-free {
  background: rgba(var(--v-theme-success), 0.08);
  border: 1px dashed rgba(var(--v-theme-success), 0.3);
  cursor: pointer;
}

.cell-free:hover {
  background: rgba(var(--v-theme-success), 0.16);
}

.cell-locked {
  background: rgba(var(--v-theme-on-surface), 0.04);
  cursor: not-allowed;
}

.cell-occupied {
  background: rgba(var(--v-theme-primary), 0.1);
  cursor: default;
  gap: 4px;
}

.cell-booked-by-me {
  background: rgba(var(--v-theme-primary), 0.2);
  border: 1px solid rgba(var(--v-theme-primary), 0.4);
}

.cell-interactive {
  cursor: pointer;
}

.cell-interactive:hover {
  opacity: 0.8;
}

.cell-inert {
  cursor: default;
}

.cell-occupied-past {
  opacity: 0.5;
}

.cell-muted {
  opacity: 0.3;
}

.matrix-past-day {
  opacity: 0.5;
}

.cell-non-interactive {
  pointer-events: none;
}

.cell-avatar {
  font-size: 0.65rem;
}

.cell-initials {
  font-weight: 600;
  font-size: 0.7rem;
  color: rgb(var(--v-theme-primary));
}

.cell-short-name {
  font-weight: 600;
  font-size: 0.7rem;
  color: rgb(var(--v-theme-primary));
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 60px;
}

</style>

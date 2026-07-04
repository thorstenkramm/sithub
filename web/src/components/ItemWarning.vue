<script setup lang="ts">
/**
 * ItemWarning renders an item's warning in the single, shared visual style used
 * across the application: dark orange text on a light orange background with the
 * orange circular info icon.
 *
 * - mode="icon" (default): the warning icon with an on-hover tooltip. Use on
 *   tiles, floor-plan items, and the weekly-table rows.
 * - mode="inline": a styled message block. Use where the warning text is shown
 *   directly (e.g. an expanded tile).
 *
 * Presentation only — booking and suppression logic live elsewhere.
 */
withDefaults(
  defineProps<{
    warning: string;
    mode?: 'icon' | 'inline';
    iconSize?: number | string;
    /** Optional data-cy applied to the icon button (mode="icon") or block (mode="inline"). */
    dataCy?: string;
  }>(),
  { mode: 'icon', iconSize: 18, dataCy: undefined },
);
</script>

<template>
  <v-tooltip v-if="mode === 'icon'" location="top" content-class="warning-tooltip">
    <template #activator="{ props: tooltipProps }">
      <v-btn
        v-bind="tooltipProps"
        icon
        variant="text"
        size="x-small"
        color="warning"
        :data-cy="dataCy"
      >
        <v-icon :size="iconSize">$warning</v-icon>
      </v-btn>
    </template>
    {{ warning }}
  </v-tooltip>

  <div v-else class="item-warning-inline" :data-cy="dataCy">
    <v-icon size="18" class="item-warning-inline__icon mr-2">$warning</v-icon>
    <span>{{ warning }}</span>
  </div>
</template>

<style scoped>
.item-warning-inline {
  display: flex;
  align-items: flex-start;
  background-color: #fff3e0;
  color: #e65100;
  font-weight: 500;
  border-radius: 4px;
  padding: 8px 12px;
  white-space: pre-line;
}

.item-warning-inline :deep(.item-warning-inline__icon) {
  color: #e65100;
}
</style>

<!--
  The tooltip content is teleported to the Vuetify overlay, so its style must be
  global. This block is the single source of the shared warning tooltip look.
-->
<style>
.warning-tooltip.v-overlay__content {
  background-color: #fff3e0 !important;
  color: #e65100 !important;
  font-weight: 500;
}
</style>

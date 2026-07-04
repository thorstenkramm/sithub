<script setup lang="ts">
/**
 * ItemWarning renders an item's warning in the single, shared visual style used
 * across the application: dark orange text on a light orange background with the
 * orange circular info icon.
 *
 * - mode="icon" (default): the warning icon with an on-hover tooltip. Use on
 *   tiles (iconVariant="button"), floor-plan items, and weekly-table rows
 *   (iconVariant="plain").
 * - mode="inline": a styled message block. Use where the warning text is shown
 *   directly (e.g. an expanded tile).
 *
 * Presentation only — booking and suppression logic live elsewhere.
 */
withDefaults(
  defineProps<{
    warning: string;
    mode?: 'icon' | 'inline';
    /** 'button' matches the tile icon; 'plain' is a bare inline icon (table rows). */
    iconVariant?: 'button' | 'plain';
    /** Tooltip location for icon mode. */
    location?: 'top' | 'bottom' | 'start' | 'end' | 'left' | 'right';
    iconSize?: number | string;
    /** inline mode only: render the leading icon. Off for the confirmation dialog. */
    showIcon?: boolean;
    /** Optional data-cy applied to the icon (mode="icon") or block (mode="inline"). */
    dataCy?: string;
  }>(),
  { mode: 'icon', iconVariant: 'button', location: 'top', iconSize: 18, showIcon: true, dataCy: undefined },
);
</script>

<template>
  <v-tooltip v-if="mode === 'icon'" :location="location" content-class="warning-tooltip">
    <template #activator="{ props: tooltipProps }">
      <v-btn
        v-if="iconVariant === 'button'"
        v-bind="tooltipProps"
        icon
        variant="text"
        size="x-small"
        color="warning"
        :data-cy="dataCy"
      >
        <v-icon :size="iconSize">$warning</v-icon>
      </v-btn>
      <v-icon
        v-else
        v-bind="tooltipProps"
        :size="iconSize"
        color="warning"
        class="ml-1"
        :data-cy="dataCy"
      >$warning</v-icon>
    </template>
    {{ warning }}
  </v-tooltip>

  <div v-else class="item-warning-inline" :data-cy="dataCy">
    <v-icon v-if="showIcon" size="18" class="item-warning-inline__icon mr-2">$warning</v-icon>
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

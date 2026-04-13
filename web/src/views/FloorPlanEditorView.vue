<template>
  <div class="page-container">
    <PageHeader
      :title="$t('floorPlanEditor.title')"
      :breadcrumbs="[{ text: $t('common.home'), to: '/' }, { text: $t('floorPlanEditor.title') }]"
    />

    <v-alert
      v-if="isMobileViewport"
      type="info"
      variant="tonal"
      class="mb-4"
      data-cy="editor-mobile-banner"
    >
      {{ $t('floorPlanEditor.desktopRecommended') }}
    </v-alert>

    <v-row>
      <v-col cols="12">
        <v-card class="mb-4">
          <v-card-text>
            <div class="d-flex flex-wrap align-center ga-3">
              <v-select
                v-model="selectedFloorPlan"
                :items="floorPlanOptions"
                item-title="label"
                item-value="value"
                :label="$t('floorPlanEditor.floorPlanLabel')"
                density="compact"
                hide-details
                :disabled="saving"
                data-cy="floor-plan-selector"
                style="min-width: 220px; max-width: 320px"
              />

              <v-btn-toggle
                v-if="isAreaLevel"
                v-model="activeTab"
                mandatory
                density="compact"
                data-cy="editor-tabs"
              >
                <v-btn value="areas" size="small" data-cy="tab-areas"
                  >{{ $t('floorPlanEditor.areasTab') }}</v-btn
                >
                <v-btn value="items" size="small" data-cy="tab-items"
                  >{{ $t('floorPlanEditor.itemsTab') }}</v-btn
                >
              </v-btn-toggle>

              <v-select
                v-if="isAreaLevel"
                :model-value="selectedSubAreaId"
                :items="subAreas"
                item-title="name"
                item-value="id"
                :label="$t('floorPlanEditor.subArea')"
                density="compact"
                hide-details
                :disabled="subAreas.length === 0"
                data-cy="toolbar-subarea-select"
                style="min-width: 180px; max-width: 240px"
                @update:model-value="onToolbarSubAreaSelect"
              />

              <v-select
                v-if="selectedFloorPlan && !(isAreaLevel && activeTab === 'areas')"
                :model-value="toolbarSelectedItemId"
                :items="toolbarItems"
                item-title="name"
                item-value="id"
                :label="$t('floorPlanEditor.itemsLabel')"
                density="compact"
                hide-details
                clearable
                :disabled="isAreaLevel && !selectedSubAreaId"
                data-cy="toolbar-items-select"
                style="min-width: 200px; max-width: 280px"
                @update:model-value="onToolbarItemSelect"
              />

              <v-select
                v-model="borderWidth"
                :items="[1, 2, 3, 4, 5]"
                :label="$t('floorPlanEditor.lineLabel')"
                density="compact"
                hide-details
                data-cy="border-width-selector"
                style="min-width: 90px; max-width: 100px"
              />

              <div class="d-flex align-center ga-1">
                <v-btn size="x-small" variant="text" data-cy="zoom-out-btn" @click="adjustZoom(-0.1)">-</v-btn>
                <span class="text-caption text-medium-emphasis" data-cy="zoom-label" style="min-width: 36px; text-align: center">
                  {{ Math.round(zoomScale * 100) }}%
                </span>
                <v-btn size="x-small" variant="text" data-cy="zoom-in-btn" @click="adjustZoom(0.1)">+</v-btn>
              </div>

              <v-chip
                v-if="saveState !== 'idle'"
                :color="saveState === 'saving' ? 'warning' : 'success'"
                size="small"
                variant="tonal"
                data-cy="editor-save-indicator"
              >
                <v-progress-circular v-if="saveState === 'saving'" size="14" width="2" indeterminate class="mr-1" />
                {{ saveState === 'saving' ? $t('floorPlanEditor.saving') : $t('floorPlanEditor.saved') }}
              </v-chip>

              <v-spacer />

              <v-btn
                v-if="selectedRectId"
                color="error"
                variant="text"
                :disabled="saving"
                data-cy="delete-rect-btn"
                @click="deleteSelected"
              >
                {{ $t('floorPlanEditor.deleteItem', { name: selectedRectName }) }}
              </v-btn>

            </div>
          </v-card-text>
        </v-card>

        <v-alert
          v-if="!selectedFloorPlan"
          type="info"
          variant="tonal"
          class="mb-4"
          data-cy="editor-no-floor-plan-alert"
        >
          {{ $t('floorPlanEditor.selectFloorPlanToPosition') }}
        </v-alert>

        <v-card v-if="selectedFloorPlan" data-cy="floor-plan-canvas-card">
          <v-card-text class="pa-2">
            <div ref="canvasShellRef" class="editor-shell" @wheel="onWheelZoom">
              <div class="editor-zoom-layer" :style="zoomLayerStyle">
                <div
                  ref="containerRef"
                  class="floor-plan-editor-container"
                  :class="{
                    'draw-mode': drawModeItemId !== null,
                    'floor-plan-editor-container--saving': saving,
                  }"
                  @pointerdown="onCanvasPointerDown"
                  @pointermove="onCanvasPointerMove"
                  @pointerup="onCanvasPointerUp"
                >
                  <img
                    ref="floorPlanImageRef"
                    :src="`/api/v1/floor-plans/${encodeURIComponent(selectedFloorPlan)}`"
                    draggable="false"
                    class="floor-plan-editor-image"
                    @load="onEditorImageLoad"
                  />

                  <div
                    v-for="pos in contextPositions"
                    :key="`ctx-${pos.itemId}`"
                    class="floor-plan-rect rect-context"
                    :style="rectStyle(pos)"
                  >
                    <span class="rect-label">{{ pos.itemName }}</span>
                  </div>

                  <div
                    v-for="pos in activePositions"
                    :key="pos.itemId"
                    class="floor-plan-rect"
                    :class="{
                      'rect-selected': selectedRectId === pos.itemId,
                      'rect-saved': recentlySaved.has(pos.itemId),
                    }"
                    :style="rectStyle(pos)"
                    :data-cy="`rect-${pos.itemId}`"
                    @pointerdown.stop="startMoveRect($event, pos.itemId)"
                    @click.stop="selectedRectId = pos.itemId"
                  >
                    <span class="rect-label">{{
                      pos.label || pos.itemName
                    }}</span>

                    <template v-if="selectedRectId === pos.itemId">
                      <div
                        class="resize-handle resize-handle--nw"
                        @pointerdown.stop="
                          startResizeRect($event, pos.itemId, 'nw')
                        "
                      />
                      <div
                        class="resize-handle resize-handle--ne"
                        @pointerdown.stop="
                          startResizeRect($event, pos.itemId, 'ne')
                        "
                      />
                      <div
                        class="resize-handle resize-handle--sw"
                        @pointerdown.stop="
                          startResizeRect($event, pos.itemId, 'sw')
                        "
                      />
                      <div
                        class="resize-handle resize-handle--se"
                        @pointerdown.stop="
                          startResizeRect($event, pos.itemId, 'se')
                        "
                      />
                    </template>
                  </div>

                  <div
                    v-if="drawPreview"
                    class="floor-plan-rect rect-preview"
                    :style="rectStyle(drawPreview)"
                  />
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar
      v-model="showSnackbar"
      :timeout="3000"
      location="bottom"
      :color="snackbarColor"
      data-cy="editor-snackbar"
    >
      {{ snackbarMessage }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { fetchAreas } from "../api/areas";
import {
  createFloorPlanPosition,
  deleteFloorPlanPosition,
  fetchFloorPlanPositions,
  updateFloorPlanPosition,
} from "../api/floorPlanPositions";
import { fetchItemGroups } from "../api/itemGroups";
import { fetchItems } from "../api/items";
import type { FloorPlanPositionAttributes } from "../api/floorPlanPositions";
import type { JsonApiResource } from "../api/types";
import { PageHeader } from "../components";

interface FloorPlanOption {
  label: string;
  value: string;
  areaId: string;
  itemGroupId?: string;
}

interface LocalPosition {
  id?: string;
  itemId: string;
  itemName: string;
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
  borderWidth: number;
}

interface EditableItem {
  id: string;
  name: string;
  positioned: boolean;
  scope: string;
}

interface SubArea {
  id: string;
  name: string;
}

type SaveState = "idle" | "saving" | "saved";

type ResizeHandle = "nw" | "ne" | "sw" | "se";

const { t } = useI18n();
const isMobileViewport = ref(typeof window !== 'undefined' && window.innerWidth < 768);

const floorPlanOptions = ref<FloorPlanOption[]>([]);
const selectedFloorPlan = ref<string | null>(null);
const allPositions = ref<LocalPosition[]>([]);
const allEditableItems = ref<EditableItem[]>([]);
const deletedPositionIDs = ref<string[]>([]);
const dirtyItemIDs = ref<Set<string>>(new Set());
const subAreas = ref<SubArea[]>([]);
const activeTab = ref<"areas" | "items">("areas");
const selectedSubAreaId = ref<string | null>(null);
const drawModeItemId = ref<string | null>(null);
const selectedRectId = ref<string | null>(null);
const borderWidth = ref(2);
const zoomScale = ref(1);
const saving = ref(false);

const drawPreview = ref<{
  x: number;
  y: number;
  width: number;
  height: number;
} | null>(null);
const drawStart = ref<{ x: number; y: number } | null>(null);
const movingItemId = ref<string | null>(null);
const moveOffset = ref<{ x: number; y: number }>({ x: 0, y: 0 });
const resizeState = ref<{
  itemId: string;
  handle: ResizeHandle;
  start: { x: number; y: number };
  original: LocalPosition;
} | null>(null);

const saveState = ref<SaveState>("idle");
const recentlySaved = ref<Set<string>>(new Set());
const snackbarMessage = ref("");
const snackbarColor = ref<"success" | "error">("success");
let saveStateResetTimeoutId: number | undefined;
let recentlySavedTimeoutId: number | undefined;
const showSnackbar = computed({
  get: () => snackbarMessage.value.length > 0,
  set: (value: boolean) => {
    if (!value) {
      snackbarMessage.value = "";
    }
  },
});

const containerRef = ref<HTMLElement | null>(null);
const floorPlanImageRef = ref<HTMLImageElement | null>(null);
const canvasShellRef = ref<HTMLElement | null>(null);

function onEditorImageLoad() {
  const shell = canvasShellRef.value;
  const img = floorPlanImageRef.value;
  if (!shell || !img) return;
  img.style.width = `${shell.clientWidth}px`;
}

const currentOption = computed(
  () =>
    floorPlanOptions.value.find(
      (option) => option.value === selectedFloorPlan.value,
    ) || null,
);

const isAreaLevel = computed(
  () => !!currentOption.value && !currentOption.value.itemGroupId,
);

const activeScope = computed(() => {
  if (!isAreaLevel.value) {
    return "items";
  }
  return activeTab.value === "areas" ? "area" : selectedSubAreaId.value || "";
});

const scopedItems = computed(() =>
  allEditableItems.value.filter((item) => item.scope === activeScope.value),
);

const toolbarItems = computed(() => {
  if (!isAreaLevel.value) {
    return allEditableItems.value;
  }
  if (!selectedSubAreaId.value) {
    return [];
  }
  return allEditableItems.value.filter(
    (item) => item.scope === selectedSubAreaId.value,
  );
});

const toolbarSelectedItemId = computed(() => {
  const selectedItemId = drawModeItemId.value ?? selectedRectId.value;
  if (!selectedItemId) {
    return null;
  }
  return toolbarItems.value.some((item) => item.id === selectedItemId)
    ? selectedItemId
    : null;
});

const activePositions = computed(() => {
  if (!isAreaLevel.value) {
    return allPositions.value;
  }

  // On Areas tab with a subarea selected, only that subarea's rect is active
  if (activeTab.value === "areas" && selectedSubAreaId.value) {
    return allPositions.value.filter(
      (pos) => pos.itemId === selectedSubAreaId.value,
    );
  }

  const ids = new Set(scopedItems.value.map((item) => item.id));
  return allPositions.value.filter((pos) => ids.has(pos.itemId));
});

const contextPositions = computed(() => {
  if (!isAreaLevel.value) {
    return [];
  }

  const areaIDs = new Set(
    allEditableItems.value
      .filter((item) => item.scope === "area")
      .map((item) => item.id),
  );
  if (activeTab.value === "areas") {
    // When a subarea is selected, all OTHER positions are context (locked)
    if (selectedSubAreaId.value) {
      return allPositions.value.filter(
        (pos) => pos.itemId !== selectedSubAreaId.value,
      );
    }
    return allPositions.value.filter((pos) => !areaIDs.has(pos.itemId));
  }
  return allPositions.value.filter((pos) => areaIDs.has(pos.itemId));
});

const hasUnsavedChanges = computed(
  () => dirtyItemIDs.value.size > 0 || deletedPositionIDs.value.length > 0,
);

const selectedRectName = computed(() => {
  if (!selectedRectId.value) {
    return "";
  }
  return (
    allPositions.value.find((pos) => pos.itemId === selectedRectId.value)
      ?.itemName || ""
  );
});

const zoomLayerStyle = computed(() => ({
  transform: `scale(${zoomScale.value})`,
  transformOrigin: "top left",
}));

function markDirty(itemId: string) {
  const next = new Set(dirtyItemIDs.value);
  next.add(itemId);
  dirtyItemIDs.value = next;
}

function clearDirty() {
  dirtyItemIDs.value = new Set();
  deletedPositionIDs.value = [];
}

function rectStyle(pos: {
  x: number;
  y: number;
  width: number;
  height: number;
  borderWidth?: number;
}) {
  return {
    left: `${pos.x}%`,
    top: `${pos.y}%`,
    width: `${pos.width}%`,
    height: `${pos.height}%`,
    borderWidth: `${pos.borderWidth || 2}px`,
  };
}

function clamp(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value));
}

function toPercent(event: PointerEvent) {
  const element = containerRef.value;
  if (!element) {
    return { x: 0, y: 0 };
  }

  const rect = element.getBoundingClientRect();
  return {
    x: ((event.clientX - rect.left) / rect.width) * 100,
    y: ((event.clientY - rect.top) / rect.height) * 100,
  };
}

function onToolbarSubAreaSelect(subAreaId: string | null) {
  selectedSubAreaId.value = subAreaId;
  if (subAreaId && isAreaLevel.value && activeTab.value === "areas") {
    const item = allEditableItems.value.find(
      (entry) => entry.id === subAreaId && entry.scope === "area",
    );
    if (item) {
      selectSidebarItem(item);
    }
  }
}

function onToolbarItemSelect(itemId: string | null) {
  if (!itemId) {
    drawModeItemId.value = null;
    selectedRectId.value = null;
    return;
  }
  const item = allEditableItems.value.find((entry) => entry.id === itemId);
  if (!item) {
    return;
  }

  if (isAreaLevel.value && activeTab.value !== "items") {
    activeTab.value = "items";
    void nextTick(() => {
      selectSidebarItem(item);
    });
    return;
  }

  selectSidebarItem(item);
}

function selectSidebarItem(item: EditableItem) {
  if (item.positioned) {
    selectedRectId.value = item.id;
    drawModeItemId.value = null;
    return;
  }

  drawModeItemId.value = item.id;
  selectedRectId.value = null;
}

function onCanvasPointerDown(event: PointerEvent) {
  if (saving.value) {
    return;
  }
  if (!drawModeItemId.value) {
    selectedRectId.value = null;
    return;
  }

  const position = toPercent(event);
  drawStart.value = position;
  drawPreview.value = {
    x: position.x,
    y: position.y,
    width: 0,
    height: 0,
  };
  (event.target as HTMLElement).setPointerCapture(event.pointerId);
}

function onCanvasPointerMove(event: PointerEvent) {
  if (drawStart.value && drawPreview.value) {
    const position = toPercent(event);
    drawPreview.value = {
      x: clamp(Math.min(drawStart.value.x, position.x), 0, 99.5),
      y: clamp(Math.min(drawStart.value.y, position.y), 0, 99.5),
      width: clamp(Math.abs(position.x - drawStart.value.x), 0, 100),
      height: clamp(Math.abs(position.y - drawStart.value.y), 0, 100),
    };
    return;
  }

  if (resizeState.value) {
    const position = toPercent(event);
    const target = allPositions.value.find(
      (pos) => pos.itemId === resizeState.value?.itemId,
    );
    if (!target) {
      return;
    }

    const updated = resizeFromHandle(
      resizeState.value.original,
      resizeState.value.handle,
      position,
    );
    Object.assign(target, updated);
    markDirty(target.itemId);
    return;
  }

  if (movingItemId.value) {
    const position = toPercent(event);
    const target = allPositions.value.find(
      (pos) => pos.itemId === movingItemId.value,
    );
    if (!target) {
      return;
    }

    target.x = clamp(position.x - moveOffset.value.x, 0, 100 - target.width);
    target.y = clamp(position.y - moveOffset.value.y, 0, 100 - target.height);
    markDirty(target.itemId);
  }
}

function onCanvasPointerUp() {
  if (drawStart.value && drawPreview.value && drawModeItemId.value) {
    if (drawPreview.value.width > 0.5 && drawPreview.value.height > 0.5) {
      const itemId = drawModeItemId.value;
      const item = allEditableItems.value.find((entry) => entry.id === itemId);
      allPositions.value.push({
        itemId,
        itemName: item?.name || itemId,
        label: "",
        x: clamp(drawPreview.value.x, 0, 100 - drawPreview.value.width),
        y: clamp(drawPreview.value.y, 0, 100 - drawPreview.value.height),
        width: clamp(drawPreview.value.width, 0.5, 100),
        height: clamp(drawPreview.value.height, 0.5, 100),
        borderWidth: borderWidth.value,
      });
      markDirty(itemId);
      updatePositionedState();
      selectedRectId.value = itemId;
      drawModeItemId.value = null;
    }

    drawPreview.value = null;
    drawStart.value = null;
  }

  movingItemId.value = null;
  resizeState.value = null;

  if (hasUnsavedChanges.value && !saving.value) {
    void autoSave();
  }
}

function startMoveRect(event: PointerEvent, itemId: string) {
  if (saving.value) {
    return;
  }
  if (drawModeItemId.value) {
    return;
  }

  selectedRectId.value = itemId;
  const rect = allPositions.value.find((pos) => pos.itemId === itemId);
  if (!rect) {
    return;
  }

  const position = toPercent(event);
  movingItemId.value = itemId;
  moveOffset.value = {
    x: position.x - rect.x,
    y: position.y - rect.y,
  };
  (event.target as HTMLElement).setPointerCapture(event.pointerId);
}

function startResizeRect(
  event: PointerEvent,
  itemId: string,
  handle: ResizeHandle,
) {
  if (saving.value) {
    return;
  }
  selectedRectId.value = itemId;
  const rect = allPositions.value.find((pos) => pos.itemId === itemId);
  if (!rect) {
    return;
  }

  resizeState.value = {
    itemId,
    handle,
    start: toPercent(event),
    original: { ...rect },
  };
  (event.target as HTMLElement).setPointerCapture(event.pointerId);
}

function resizeFromHandle(
  rect: LocalPosition,
  handle: ResizeHandle,
  cursor: { x: number; y: number },
) {
  const minSize = 0.5;
  const left = rect.x;
  const right = rect.x + rect.width;
  const top = rect.y;
  const bottom = rect.y + rect.height;

  let nextLeft = left;
  let nextRight = right;
  let nextTop = top;
  let nextBottom = bottom;

  if (handle.includes("w")) {
    nextLeft = clamp(cursor.x, 0, right - minSize);
  } else {
    nextRight = clamp(cursor.x, left + minSize, 100);
  }

  if (handle.includes("n")) {
    nextTop = clamp(cursor.y, 0, bottom - minSize);
  } else {
    nextBottom = clamp(cursor.y, top + minSize, 100);
  }

  return {
    x: nextLeft,
    y: nextTop,
    width: nextRight - nextLeft,
    height: nextBottom - nextTop,
  };
}

function updatePositionedState() {
  const positionedIDs = new Set(allPositions.value.map((pos) => pos.itemId));
  allEditableItems.value = allEditableItems.value.map((item) => ({
    ...item,
    positioned: positionedIDs.has(item.id),
  }));
}

function deleteByItemId(itemId: string) {
  const target = allPositions.value.find((pos) => pos.itemId === itemId);
  if (!target) {
    return;
  }

  if (target.id) {
    deletedPositionIDs.value = [...deletedPositionIDs.value, target.id];
  }

  allPositions.value = allPositions.value.filter(
    (pos) => pos.itemId !== itemId,
  );
  markDirty(itemId);
  updatePositionedState();
  if (selectedRectId.value === itemId) {
    selectedRectId.value = null;
  }
  if (drawModeItemId.value === itemId) {
    drawModeItemId.value = null;
  }

  if (hasUnsavedChanges.value && !saving.value) {
    void autoSave();
  }
}

function deleteSelected() {
  if (saving.value || !selectedRectId.value) {
    return;
  }
  deleteByItemId(selectedRectId.value);
}

function onKeyDown(event: KeyboardEvent) {
  if (saving.value) {
    return;
  }
  if (event.key === "Escape") {
    drawModeItemId.value = null;
    drawPreview.value = null;
    drawStart.value = null;
    resizeState.value = null;
    movingItemId.value = null;
    return;
  }

  if (event.key === "Delete" && selectedRectId.value) {
    deleteSelected();
  }
}

function adjustZoom(delta: number) {
  zoomScale.value = clamp(zoomScale.value + delta, 0.75, 2);
}

function onWheelZoom(event: WheelEvent) {
  if (!event.ctrlKey) {
    return;
  }

  event.preventDefault();
  adjustZoom(event.deltaY < 0 ? 0.1 : -0.1);
}

async function saveChanges() {
  const floorPlanToSave = selectedFloorPlan.value;
  if (!floorPlanToSave || !hasUnsavedChanges.value) {
    return;
  }

  window.clearTimeout(saveStateResetTimeoutId);
  saving.value = true;
  saveState.value = "saving";
  try {
    for (const id of deletedPositionIDs.value) {
      await deleteFloorPlanPosition(id);
    }

    const savedItemIDs: string[] = [];
    for (const pos of allPositions.value) {
      if (pos.id && dirtyItemIDs.value.has(pos.itemId)) {
        await updateFloorPlanPosition(pos.id, {
          label: pos.label,
          x: pos.x,
          y: pos.y,
          width: pos.width,
          height: pos.height,
          border_width: pos.borderWidth,
        });
        savedItemIDs.push(pos.itemId);
      } else if (!pos.id) {
        const response = await createFloorPlanPosition({
          floor_plan: floorPlanToSave,
          item_id: pos.itemId,
          label: pos.label,
          x: pos.x,
          y: pos.y,
          width: pos.width,
          height: pos.height,
          border_width: pos.borderWidth,
        });
        pos.id = response.data.id;
        savedItemIDs.push(pos.itemId);
      }
    }

    clearDirty();
    recentlySaved.value = new Set(savedItemIDs);
    window.clearTimeout(recentlySavedTimeoutId);
    recentlySavedTimeoutId = window.setTimeout(() => {
      recentlySaved.value = new Set();
    }, 600);
    saveState.value = "saved";
    saveStateResetTimeoutId = window.setTimeout(() => {
      saveState.value = "idle";
    }, 1500);
  } catch {
    window.clearTimeout(recentlySavedTimeoutId);
    saveState.value = "idle";
    snackbarColor.value = "error";
    snackbarMessage.value = t('floorPlanEditor.saveFailed');
  } finally {
    saving.value = false;
  }
}

async function autoSave() {
  await saveChanges();
}

async function loadFloorPlanOptions() {
  try {
    const areasResponse = await fetchAreas();
    const options: FloorPlanOption[] = [];

    const areaResults = await Promise.all(
      areasResponse.data.map(async (area) => {
        const areaOptions: FloorPlanOption[] = [];
        if (area.attributes.floor_plan) {
          areaOptions.push({
            label: area.attributes.name,
            value: area.attributes.floor_plan,
            areaId: area.id,
          });
        }

        try {
          const itemGroupsResponse = await fetchItemGroups(area.id);
          for (const itemGroup of itemGroupsResponse.data) {
            if (itemGroup.attributes.floor_plan) {
              areaOptions.push({
                label: `${area.attributes.name} > ${itemGroup.attributes.name}`,
                value: itemGroup.attributes.floor_plan,
                areaId: area.id,
                itemGroupId: itemGroup.id,
              });
            }
          }
        } catch {
          // Item groups unavailable for this area; skip.
        }
        return areaOptions;
      }),
    );

    for (const areaOptions of areaResults) {
      options.push(...areaOptions);
    }

    floorPlanOptions.value = options;
    if (!selectedFloorPlan.value && options[0]) {
      selectedFloorPlan.value = options[0].value;
    }
  } catch {
    floorPlanOptions.value = [];
  }
}

watch(selectedFloorPlan, async (floorPlan) => {
  if (!floorPlan) {
    allPositions.value = [];
    allEditableItems.value = [];
    subAreas.value = [];
    clearDirty();
    return;
  }

  drawModeItemId.value = null;
  selectedRectId.value = null;
  activeTab.value = "areas";
  selectedSubAreaId.value = null;
  zoomScale.value = 1;

  const option = currentOption.value;
  if (!option) {
    return;
  }

  try {
    const response = await fetchFloorPlanPositions(floorPlan);
    allPositions.value = response.data.map(
      (resource: JsonApiResource<FloorPlanPositionAttributes>) => ({
        id: resource.id,
        itemId: resource.attributes.item_id,
        itemName: resource.attributes.item_id,
        label: resource.attributes.label || "",
        x: resource.attributes.x,
        y: resource.attributes.y,
        width: resource.attributes.width,
        height: resource.attributes.height,
        borderWidth: resource.attributes.border_width || 2,
      }),
    );
  } catch {
    allPositions.value = [];
  }

  const items: EditableItem[] = [];
  const nextSubAreas: SubArea[] = [];

  try {
    if (option.itemGroupId) {
      const response = await fetchItems(option.itemGroupId);
      for (const item of response.data) {
        items.push({
          id: item.id,
          name: item.attributes.name,
          positioned: false,
          scope: "items",
        });
      }
    } else {
      const itemGroupsResponse = await fetchItemGroups(option.areaId);
      for (const itemGroup of itemGroupsResponse.data) {
        nextSubAreas.push({
          id: itemGroup.id,
          name: itemGroup.attributes.name,
        });
        items.push({
          id: itemGroup.id,
          name: itemGroup.attributes.name,
          positioned: false,
          scope: "area",
        });
      }

      const igItemResults = await Promise.all(
        itemGroupsResponse.data.map(async (itemGroup) => {
          const itemsResponse = await fetchItems(itemGroup.id);
          return { itemGroupId: itemGroup.id, items: itemsResponse.data };
        }),
      );
      for (const { itemGroupId, items: igItems } of igItemResults) {
        for (const item of igItems) {
          items.push({
            id: item.id,
            name: item.attributes.name,
            positioned: false,
            scope: itemGroupId,
          });
        }
      }
    }
  } catch {
    // Leave items empty; saving is not possible without loaded data.
  }

  subAreas.value = nextSubAreas;
  allEditableItems.value = items;

  for (const position of allPositions.value) {
    const match = items.find((item) => item.id === position.itemId);
    if (match) {
      position.itemName = match.name;
    }
  }

  clearDirty();
  updatePositionedState();

  if (nextSubAreas[0]) {
    onToolbarSubAreaSelect(nextSubAreas[0].id);
  }
});

watch(activeTab, () => {
  drawModeItemId.value = null;
  selectedRectId.value = null;
});

watch(selectedSubAreaId, () => {
  if (activeTab.value !== "items") {
    return;
  }

  drawModeItemId.value = null;
  selectedRectId.value = null;
});

onMounted(async () => {
  window.addEventListener("keydown", onKeyDown);
  await loadFloorPlanOptions();
});

onUnmounted(() => {
  window.removeEventListener("keydown", onKeyDown);
  window.clearTimeout(saveStateResetTimeoutId);
  window.clearTimeout(recentlySavedTimeoutId);
});
</script>

<style scoped>
.editor-shell {
  overflow: auto;
  max-height: calc(100vh - 130px);
}

.editor-zoom-layer {
  display: inline-block;
}

.floor-plan-editor-container {
  position: relative;
  display: inline-block;
  cursor: default;
}

.floor-plan-editor-container--saving {
  pointer-events: none;
}

.floor-plan-editor-container.draw-mode {
  cursor: crosshair;
}

.floor-plan-editor-image {
  display: block;
  max-width: none;
  height: auto;
  user-select: none;
}

:global(.v-theme--dark) .floor-plan-editor-image {
  filter: brightness(0.85) contrast(1.1);
}

.floor-plan-rect {
  position: absolute;
  border: 2px solid rgb(var(--v-theme-primary));
  background-color: rgba(var(--v-theme-primary), 0.1);
  cursor: move;
  user-select: none;
}

.floor-plan-rect.rect-context {
  border: 1px dashed rgb(var(--v-theme-outline));
  background-color: rgba(var(--v-theme-outline), 0.08);
  pointer-events: none;
  opacity: 0.35;
}

.floor-plan-rect.rect-selected {
  border-color: rgb(var(--v-theme-info));
  box-shadow: 0 0 0 2px rgba(var(--v-theme-info), 0.2);
}

.floor-plan-rect.rect-saved {
  border-color: rgb(var(--v-theme-success));
  background-color: rgba(var(--v-theme-success), 0.2);
  transition:
    border-color 0.25s,
    background-color 0.25s;
}

.floor-plan-rect.rect-preview {
  border-style: dashed;
  border-color: rgb(var(--v-theme-secondary));
  background-color: rgba(var(--v-theme-secondary), 0.1);
  pointer-events: none;
}

.rect-label {
  position: absolute;
  top: 2px;
  left: 4px;
  max-width: calc(100% - 8px);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.72rem;
  font-weight: 600;
  pointer-events: none;
}

.resize-handle {
  position: absolute;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgb(var(--v-theme-surface));
  border: 2px solid rgb(var(--v-theme-info));
}

.resize-handle--nw {
  top: -6px;
  left: -6px;
  cursor: nwse-resize;
}

.resize-handle--ne {
  top: -6px;
  right: -6px;
  cursor: nesw-resize;
}

.resize-handle--sw {
  bottom: -6px;
  left: -6px;
  cursor: nesw-resize;
}

.resize-handle--se {
  right: -6px;
  bottom: -6px;
  cursor: nwse-resize;
}
</style>

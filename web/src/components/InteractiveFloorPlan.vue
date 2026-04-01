<template>
  <div class="fp-root">
    <div class="fp-header d-flex align-center flex-wrap ga-2" data-cy="fp-header">
      <div class="d-flex align-center flex-wrap ga-2">
        <template v-if="drilledInto">
          <v-btn
            variant="text"
            size="small"
            data-cy="fp-breadcrumb-root"
            @click="drillBack"
          >
            {{ title }}
          </v-btn>
          <v-icon size="16">$chevronRight</v-icon>
          <span class="text-h6" data-cy="fp-breadcrumb-current">{{
            drilledInto.name
          }}</span>
        </template>
        <span v-else class="text-h6">{{ title }}</span>
        <span class="text-body-2 text-medium-emphasis">{{ weekLabel }}</span>
      </div>

      <v-spacer />

      <div class="fp-zoom-controls d-flex align-center ga-2 flex-wrap">
        <v-btn
          size="x-small"
          variant="text"
          data-cy="fp-zoom-out"
          @click="adjustZoom(-0.1)"
        >
          -
        </v-btn>
        <span class="text-caption text-medium-emphasis"
          >{{ Math.round(zoomScale * 100) }}%</span
        >
        <v-btn
          size="x-small"
          variant="text"
          data-cy="fp-zoom-in"
          @click="adjustZoom(0.1)"
        >
          +
        </v-btn>
      </div>
    </div>

    <div
      class="fp-weekday-selector"
      :class="
        isCompactViewport
          ? 'fp-weekday-selector--grid'
          : 'fp-weekday-selector--row'
      "
      data-cy="fp-weekday-selector"
    >
      <v-btn
        v-for="(day, idx) in weekdays"
        :key="day.date"
        :variant="idx === selectedDayIndex ? 'flat' : 'outlined'"
        :color="idx === selectedDayIndex ? 'error' : undefined"
        :disabled="day.past"
        size="small"
        :data-cy="`fp-day-${day.label}`"
        @click="selectedDayIndex = idx"
      >
        {{ day.label }}
      </v-btn>
    </div>

    <div v-if="initialLoading" class="text-center pa-8">
      <v-progress-circular indeterminate color="primary" />
    </div>

    <v-alert
      v-else-if="
        !initialLoading &&
        (isAreaView
          ? areaPositions.length === 0
          : enrichedPositions.length === 0)
      "
      type="info"
      variant="tonal"
      data-cy="fp-no-positions"
    >
      No items have been positioned on this floor plan yet. An administrator can
      set them up in the Floor Plan Editor.
    </v-alert>

    <div
      v-else
      class="fp-scroll-shell"
      :class="{ 'fp-scroll-shell--compact': isCompactViewport }"
      @wheel="onWheelZoom"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd"
    >
      <Transition name="fp-drill" mode="out-in">
        <div
          :key="`${activeFloorPlan}:${activeItemGroupId}:${isAreaView ? 'area' : 'items'}`"
          class="fp-zoom-layer"
          :style="zoomLayerStyle"
        >
          <div class="fp-fit-container">
            <div
              class="fp-container"
              :class="{
                'fp-label-hidden': !showLabels,
                'fp-content-loading': availabilityLoading,
              }"
            >
              <img
                :src="`/api/v1/floor-plans/${encodeURIComponent(activeFloorPlan)}`"
                draggable="false"
                class="fp-image-fit"
              />

              <template v-if="isAreaView">
                <div
                  v-for="pos in areaPositions"
                  :key="pos.itemId"
                  class="fp-item fp-item--area"
                  :class="{ 'fp-item--clickable': canDrillInto(pos.itemId) }"
                  :style="{
                    ...rectStyle(pos),
                    borderColor: areaAvailColor(pos.itemId),
                    backgroundColor: areaAvailBg(pos.itemId),
                  }"
                  :data-cy="`fp-area-${pos.itemId}`"
                  @click="
                    canDrillInto(pos.itemId) && handleAreaClick(pos.itemId)
                  "
                >
                  <span class="fp-item-label">{{ pos.displayLabel }}</span>
                  <span class="fp-item-fraction">{{
                    areaFractionLabel(pos.itemId)
                  }}</span>
                </div>

                <div
                  v-for="pos in deskPositions"
                  :key="`desk-${pos.itemId}`"
                  class="fp-item"
                  :class="{
                    'fp-item--free': pos.status === 'free',
                    'fp-item--busy': pos.status === 'busy',
                    'fp-item--clickable':
                      pos.status === 'free' &&
                      shouldDrillIntoItemGroup(itemToGroupMap.get(pos.itemId)),
                  }"
                  :style="{
                    ...rectStyle(pos),
                    pointerEvents: pos.status === 'free' ? 'auto' : 'none',
                  }"
                  :data-cy="`fp-desk-${pos.itemId}`"
                  @click="
                    pos.status === 'free' &&
                    handleDeskClick(pos.itemId, pos.displayLabel)
                  "
                >
                  <span class="fp-item-label">{{ pos.displayLabel }}</span>
                </div>
              </template>

              <template v-else>
                <template v-for="pos in enrichedPositions" :key="pos.itemId">
                  <v-tooltip
                    v-if="pos.status === 'free'"
                    location="top"
                    :text="pos.tooltipText"
                  >
                    <template #activator="{ props: tooltipProps }">
                      <div
                        v-bind="tooltipProps"
                        class="fp-item fp-item--free"
                        :style="rectStyle(pos)"
                        :data-cy="`fp-item-${pos.itemId}`"
                        @click="requestBooking(pos.itemId, pos.displayLabel)"
                      >
                        <span class="fp-item-label">{{
                          pos.displayLabel
                        }}</span>
                      </div>
                    </template>
                  </v-tooltip>

                  <div
                    v-else
                    class="fp-item"
                    :class="
                      'fp-item--busy'
                    "
                    :style="rectStyle(pos)"
                    :data-cy="`fp-item-${pos.itemId}`"
                  >
                    <span class="fp-item-label">{{ pos.displayLabel }}</span>
                  </div>
                </template>
              </template>
            </div>
          </div>
        </div>
      </Transition>
    </div>

    <div class="fp-footer" data-cy="fp-footer">
      <v-checkbox
        v-if="positions.length > 0"
        v-model="showLabels"
        label="Show labels"
        hide-details
        density="compact"
        data-cy="fp-show-labels"
      />
      <div v-else />
      <v-btn
        variant="text"
        data-cy="fp-close-btn"
        class="fp-close-btn"
        @click="handleClose"
      >
        Close
      </v-btn>
    </div>

    <v-dialog
      v-model="showBookingDialog"
      :fullscreen="isCompactViewport"
      max-width="560"
      persistent
      data-cy="fp-booking-dialog"
    >
      <v-card class="fp-booking-dialog-card">
        <v-card-title data-cy="fp-booking-title">
          {{ bookingDialogTitle }}
        </v-card-title>
        <v-card-text>
          <div class="text-body-1 mb-4" data-cy="fp-booking-summary">
            {{ bookingSummary }}
          </div>

          <div
            v-if="bookingDayAvailabilityLoading"
            class="text-body-2 text-medium-emphasis mb-3"
            data-cy="fp-booking-days-loading"
          >
            Checking availability for the selected week...
          </div>

          <div data-cy="fp-booking-days">
            <div
              v-for="day in bookingDayOptions"
              :key="day.date"
              :class="[
                'fp-booking-day-row',
                { 'fp-booking-day-row--past': day.past },
              ]"
              :data-cy="`fp-booking-day-${day.label}`"
            >
              <v-checkbox
                v-if="getBookingDayStatus(day.date).status === 'free'"
                :model-value="bookingDaySelections[day.date] === true"
                hide-details
                density="compact"
                color="success"
                :disabled="day.past || bookingDayAvailabilityLoading"
                class="fp-booking-day-checkbox"
                @update:model-value="toggleBookingDaySelection(day.date)"
              />
              <v-checkbox
                v-else-if="
                  getBookingDayStatus(day.date).status === 'booked-by-me'
                "
                :model-value="true"
                hide-details
                density="compact"
                color="primary"
                disabled
                class="fp-booking-day-checkbox"
              />
              <v-checkbox
                v-else
                :model-value="false"
                hide-details
                density="compact"
                disabled
                class="fp-booking-day-checkbox"
              />
              <span class="text-body-2 font-weight-medium">
                {{ getFullDayLabel(day.date) }}
              </span>
              <span
                :class="[
                  'text-body-2',
                  getBookingDayStatusColor(day.date, day.past),
                ]"
                :title="bookingDayTitle(day.date)"
              >
                {{ getBookingDayStatusText(day.date) }}
              </span>
            </div>
          </div>

          <div
            v-if="bookingItemEquipment.length"
            class="mt-3"
            data-cy="fp-booking-equipment"
          >
            <div class="text-caption text-medium-emphasis mb-1">Equipment</div>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="equip in bookingItemEquipment"
                :key="equip"
                size="x-small"
                variant="outlined"
              >
                {{ equip }}
              </v-chip>
            </div>
          </div>

          <v-alert
            v-if="bookingItemWarning"
            type="warning"
            variant="tonal"
            density="compact"
            class="mt-2"
            data-cy="fp-booking-warning"
          >
            {{ bookingItemWarning }}
          </v-alert>

          <div
            v-if="selectedBookingDates.length === 0"
            class="text-body-2 text-error mt-3"
            data-cy="fp-booking-no-days"
          >
            Select at least one day to continue.
          </div>
        </v-card-text>
        <v-card-actions
          class="fp-booking-dialog-actions"
          data-cy="fp-booking-actions"
        >
          <v-spacer />
          <v-btn
            variant="text"
            data-cy="fp-booking-cancel"
            @click="showBookingDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :disabled="
              selectedBookingDates.length === 0 || bookingDayAvailabilityLoading
            "
            :loading="bookingInProgress"
            data-cy="fp-booking-confirm"
            @click="confirmPendingBooking"
          >
            Book now
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar
      v-model="showBookingSnackbar"
      :timeout="5000"
      location="bottom"
      color="success"
      data-cy="fp-booking-success"
    >
      {{ bookingSnackbarText }}
      <template #actions>
        <v-btn
          v-if="lastBooking"
          variant="text"
          :loading="undoInProgress"
          @click="undoLastBooking"
        >
          Undo
        </v-btn>
      </template>
    </v-snackbar>

    <v-snackbar
      v-model="showErrorSnackbar"
      :timeout="4000"
      location="bottom"
      color="error"
      data-cy="fp-booking-error"
    >
      {{ errorSnackbarText }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { ApiError } from "../api/client";
import {
  cancelBooking,
  createBooking,
  createMultiDayBooking,
} from "../api/bookings";
import { fetchFloorPlanPositions } from "../api/floorPlanPositions";
import { fetchItemGroups } from "../api/itemGroups";
import { fetchItems } from "../api/items";
import type { FloorPlanPositionAttributes } from "../api/floorPlanPositions";
import type { ItemAttributes } from "../api/items";
import type { JsonApiResource } from "../api/types";

const props = defineProps<{
  floorPlan: string;
  title: string;
  weekLabel: string;
  weekDates: string[];
  itemGroupId: string;
  areaLevel?: boolean;
}>();

const emit = defineEmits<{
  close: [];
}>();

const initialLoading = ref(true);
const availabilityLoading = ref(false);
const selectedDayIndex = ref(0);
const isCompactViewport = ref(false);

const LABELS_KEY = "sithub_fp_show_labels";
const showLabels = ref(localStorage.getItem(LABELS_KEY) !== "false");
watch(showLabels, (value) => localStorage.setItem(LABELS_KEY, String(value)));

const zoomScale = ref(1);
const pinchState = ref<{ startDistance: number; startScale: number } | null>(
  null,
);

interface PositionData {
  id: string;
  itemId: string;
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
  borderWidth: number;
}

interface ItemData {
  name: string;
  equipment: string[];
  warning?: string;
  availability: string;
  bookerName?: string;
  bookedByMe: boolean;
}

interface BookingDayInfo {
  status: "free" | "booked-by-me" | "booked-by-other" | "unavailable";
  bookerName?: string;
}

interface DrillTarget {
  itemGroupId: string;
  name: string;
  floorPlan: string;
}

const positions = ref<PositionData[]>([]);
const itemDataMap = ref<Map<string, ItemData>>(new Map());
const itemToGroupMap = ref<Map<string, string>>(new Map());
const drilledInto = ref<DrillTarget | null>(null);
const areaItemGroupAvailability = ref<
  Map<string, { free: number; total: number }>
>(new Map());
const itemGroupMap = ref<Map<string, { name: string; floorPlan: string }>>(
  new Map(),
);

const activeFloorPlan = computed(
  () => drilledInto.value?.floorPlan || props.floorPlan,
);
const activeItemGroupId = computed(
  () => drilledInto.value?.itemGroupId || props.itemGroupId,
);
const isAreaView = computed(() => props.areaLevel && !drilledInto.value);

const weekdays = computed(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return props.weekDates.map((date) => {
    const day = new Date(`${date}T00:00:00`);
    return {
      date,
      label: formatDayLabel(day),
      past: day < today,
    };
  });
});

const selectedDate = computed(
  () => weekdays.value[selectedDayIndex.value]?.date || "",
);

const showBookingDialog = ref(false);
const bookingInProgress = ref(false);
const pendingBooking = ref<{
  itemId: string;
  label: string;
  itemGroupId: string;
} | null>(null);
const bookingDaySelections = ref<Record<string, boolean>>({});
const bookingDayInfoMap = ref<Map<string, BookingDayInfo>>(new Map());
const bookingItemEquipment = ref<string[]>([]);
const bookingItemWarning = ref<string | undefined>();
const bookingDayAvailabilityLoading = ref(false);

const showBookingSnackbar = ref(false);
const bookingSnackbarText = ref("Booking confirmed.");
const undoInProgress = ref(false);
const lastBooking = ref<{
  bookingIds: string[];
  itemId: string;
  label: string;
} | null>(null);
const showErrorSnackbar = ref(false);
const errorSnackbarText = ref("");

const zoomLayerStyle = computed(() => ({
  transform: `scale(${zoomScale.value})`,
  transformOrigin: "top left",
}));

const bookingDayOptions = computed(() => weekdays.value);

const selectedBookingDates = computed(() =>
  bookingDayOptions.value
    .filter((day) => bookingDaySelections.value[day.date] === true)
    .map((day) => day.date),
);

const bookingDialogTitle = computed(() =>
  pendingBooking.value
    ? `Confirm your booking for ${pendingBooking.value.label}`
    : "Confirm booking",
);

const bookingSummary = computed(() => {
  if (!pendingBooking.value) {
    return "";
  }

  const count = selectedBookingDates.value.length;
  const startDate = selectedBookingDates.value[0]
    ? formatReadableDate(selectedBookingDates.value[0])
    : formatReadableDate(selectedDate.value);
  const location = drilledInto.value?.name || props.title;
  const dayLabel = count === 1 ? "1 day" : `${count} days`;

  return `Book ${pendingBooking.value.label} in ${location} for ${dayLabel} starting ${startDate}.`;
});

watch(showBookingSnackbar, (open) => {
  if (!open && !undoInProgress.value) {
    lastBooking.value = null;
  }
});

watch(showBookingDialog, (open) => {
  if (!open && !bookingInProgress.value) {
    pendingBooking.value = null;
    bookingDaySelections.value = {};
    bookingDayInfoMap.value = new Map();
    bookingItemEquipment.value = [];
    bookingItemWarning.value = undefined;
    bookingDayAvailabilityLoading.value = false;
  }
});

function formatDayLabel(day: Date): string {
  const labels = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"];
  return labels[day.getDay()] || "";
}

function formatReadableDate(dateStr: string): string {
  if (!dateStr) {
    return "";
  }

  const date = new Date(`${dateStr}T00:00:00`);
  if (Number.isNaN(date.getTime())) {
    return dateStr;
  }

  return new Intl.DateTimeFormat(undefined, {
    weekday: "long",
    month: "short",
    day: "numeric",
  }).format(date);
}

const WEEKDAY_LONG_FORMATTER = new Intl.DateTimeFormat(undefined, {
  weekday: "long",
});

function getFullDayLabel(dateStr: string): string {
  const parsed = new Date(`${dateStr}T00:00:00`);
  if (!Number.isNaN(parsed.getTime())) {
    const weekday = WEEKDAY_LONG_FORMATTER.format(parsed);
    const dd = String(parsed.getDate()).padStart(2, "0");
    const mm = String(parsed.getMonth() + 1).padStart(2, "0");
    return `${weekday}, ${dd}.${mm}.`;
  }
  return dateStr;
}

function getBookingDayStatus(date: string): BookingDayInfo {
  return bookingDayInfoMap.value.get(date) ?? { status: "unavailable" };
}

function getBookingDayStatusColor(date: string, past: boolean): string {
  if (past) return "text-medium-emphasis";
  const info = getBookingDayStatus(date);
  switch (info.status) {
    case "free":
      return "text-success";
    case "booked-by-me":
      return "text-primary";
    case "booked-by-other":
      return "text-error";
    default:
      return "text-medium-emphasis";
  }
}

function getBookingDayStatusText(date: string): string {
  const info = getBookingDayStatus(date);
  switch (info.status) {
    case "free":
      return "free";
    case "booked-by-me":
      return info.bookerName || "Me";
    case "booked-by-other":
      return info.bookerName || "Booked";
    case "unavailable":
      return "n/a";
  }
}

function preselectDay() {
  const today = new Date().toISOString().slice(0, 10);
  const todayIndex = weekdays.value.findIndex((day) => day.date === today);
  if (todayIndex >= 0) {
    selectedDayIndex.value = todayIndex;
    return;
  }

  const firstFutureIndex = weekdays.value.findIndex((day) => !day.past);
  selectedDayIndex.value = firstFutureIndex >= 0 ? firstFutureIndex : 0;
}

async function loadPositions() {
  try {
    const response = await fetchFloorPlanPositions(activeFloorPlan.value);
    positions.value = response.data.map(
      (resource: JsonApiResource<FloorPlanPositionAttributes>) => ({
        id: resource.id,
        itemId: resource.attributes.item_id,
        label: resource.attributes.label || "",
        x: resource.attributes.x,
        y: resource.attributes.y,
        width: resource.attributes.width,
        height: resource.attributes.height,
        borderWidth: resource.attributes.border_width || 2,
      }),
    );
  } catch {
    positions.value = [];
  }
}

async function loadAvailability() {
  if (!selectedDate.value || !activeItemGroupId.value || isAreaView.value) {
    return;
  }

  availabilityLoading.value = true;
  try {
    const response = await fetchItems(
      activeItemGroupId.value,
      selectedDate.value,
    );
    const map = new Map<string, ItemData>();
    for (const item of response.data as JsonApiResource<ItemAttributes>[]) {
      map.set(item.id, {
        name: item.attributes.name,
        equipment: item.attributes.equipment || [],
        warning: item.attributes.warning,
        availability: item.attributes.availability,
        bookerName: item.attributes.booker_name,
        bookedByMe: item.attributes.booked_by_me === true,
      });
    }
    itemDataMap.value = map;
  } catch {
    itemDataMap.value = new Map();
  } finally {
    availabilityLoading.value = false;
  }
}

async function loadAreaAvailability() {
  if (!selectedDate.value || !isAreaView.value) {
    return;
  }

  availabilityLoading.value = true;
  try {
    const knownItemGroupIDs = new Set(itemGroupMap.value.keys());
    const itemGroupIDs = positions.value
      .map((pos) => pos.itemId)
      .filter((id) => knownItemGroupIDs.has(id));
    const availabilityMap = new Map<string, { free: number; total: number }>();
    const allItemData = new Map<string, ItemData>();
    const itemGroupLookup = new Map<string, string>();

    const results = await Promise.allSettled(
      itemGroupIDs.map(async (itemGroupID) => {
        const response = await fetchItems(itemGroupID, selectedDate.value);
        return { itemGroupID, response };
      }),
    );

    for (const result of results) {
      if (result.status === "rejected") {
        continue;
      }

      const { itemGroupID, response } = result.value;
      const total = response.data.length;
      const free = response.data.filter(
        (item: JsonApiResource<ItemAttributes>) =>
          item.attributes.availability === "available",
      ).length;

      availabilityMap.set(itemGroupID, { free, total });

      for (const item of response.data as JsonApiResource<ItemAttributes>[]) {
        allItemData.set(item.id, {
          name: item.attributes.name,
          equipment: item.attributes.equipment || [],
          warning: item.attributes.warning,
          availability: item.attributes.availability,
          bookerName: item.attributes.booker_name,
          bookedByMe: item.attributes.booked_by_me === true,
        });
        itemGroupLookup.set(item.id, itemGroupID);
      }
    }

    for (const id of itemGroupIDs) {
      if (!availabilityMap.has(id)) {
        availabilityMap.set(id, { free: 0, total: 0 });
      }
    }

    areaItemGroupAvailability.value = availabilityMap;
    itemDataMap.value = allItemData;
    itemToGroupMap.value = itemGroupLookup;
  } finally {
    availabilityLoading.value = false;
  }
}

async function loadItemGroupMap() {
  if (!props.areaLevel) {
    return;
  }

  try {
    const { fetchAreas } = await import("../api/areas");
    const areasResponse = await fetchAreas();
    const map = new Map<string, { name: string; floorPlan: string }>();

    for (const area of areasResponse.data) {
      const itemGroupsResponse = await fetchItemGroups(area.id);
      for (const itemGroup of itemGroupsResponse.data) {
        map.set(itemGroup.id, {
          name: itemGroup.attributes.name,
          floorPlan: itemGroup.attributes.floor_plan || "",
        });
      }
    }

    itemGroupMap.value = map;
  } catch {
    itemGroupMap.value = new Map();
  }
}

async function refreshAvailability() {
  if (isAreaView.value) {
    await loadAreaAvailability();
    return;
  }

  await loadAvailability();
}

async function initialLoad() {
  initialLoading.value = true;
  preselectDay();
  zoomScale.value = 1;
  showBookingDialog.value = false;
  await loadPositions();
  if (props.areaLevel) {
    await loadItemGroupMap();
  }
  await refreshAvailability();
  initialLoading.value = false;
}

watch(
  [
    () => props.floorPlan,
    () => props.itemGroupId,
    () => props.weekDates.join("|"),
    () => props.areaLevel,
  ],
  async () => {
    drilledInto.value = null;
    await initialLoad();
  },
  { immediate: true },
);

watch(selectedDayIndex, async () => {
  await refreshAvailability();
});

watch(drilledInto, async () => {
  initialLoading.value = true;
  zoomScale.value = 1;
  await loadPositions();
  await refreshAvailability();
  initialLoading.value = false;
});

const areaPositions = computed(() => {
  if (!isAreaView.value) {
    return [];
  }

  const knownItemGroupIDs = new Set(itemGroupMap.value.keys());
  return positions.value
    .filter((pos) => knownItemGroupIDs.has(pos.itemId))
    .map((pos) => {
      const itemGroup = itemGroupMap.value.get(pos.itemId);
      return {
        ...pos,
        displayLabel: pos.label || itemGroup?.name || pos.itemId,
      };
    });
});

const deskPositions = computed(() => {
  if (!isAreaView.value) {
    return [];
  }

  const knownItemGroupIDs = new Set(itemGroupMap.value.keys());
  return positions.value
    .filter((pos) => !knownItemGroupIDs.has(pos.itemId))
    .map((pos) => {
      const item = itemDataMap.value.get(pos.itemId);
      const occupied = item?.availability === "occupied";

      return {
        ...pos,
        displayLabel: pos.label || item?.name || pos.itemId,
        status: occupied ? ("busy" as const) : ("free" as const),
      };
    });
});

const enrichedPositions = computed(() => {
  if (isAreaView.value) {
    return [];
  }

  const filteredPositions = positions.value.filter((pos) =>
    itemDataMap.value.has(pos.itemId),
  );
  return filteredPositions.map((pos) => {
    const item = itemDataMap.value.get(pos.itemId);
    const name = item?.name || pos.itemId;
    const equipmentText = item?.equipment?.join(", ") || "";
    const tooltipParts = [name];
    if (equipmentText) {
      tooltipParts.push(equipmentText);
    }
    if (item?.warning) {
      tooltipParts.push(item.warning.trim());
    }

    const occupied = item?.availability === "occupied";

    return {
      ...pos,
      displayLabel: pos.label || name,
      tooltipText: tooltipParts.join("\n"),
      status: occupied ? "busy" : ("free" as "free" | "busy" | "mine"),
    };
  });
});

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

function areaAvailData(itemGroupID: string) {
  return (
    areaItemGroupAvailability.value.get(itemGroupID) || { free: 0, total: 0 }
  );
}

function areaFractionLabel(itemGroupID: string) {
  const { free, total } = areaAvailData(itemGroupID);
  if (total === 0) {
    return "";
  }
  return `${free}/${total} free`;
}

function areaAvailColor(itemGroupID: string) {
  const { free, total } = areaAvailData(itemGroupID);
  if (total === 0) {
    return "rgb(var(--v-theme-outline))";
  }

  const ratio = free / total;
  if (ratio > 0.5) {
    return "rgb(var(--v-theme-success))";
  }
  if (ratio > 0) {
    return "rgb(var(--v-theme-warning))";
  }
  return "rgb(var(--v-theme-error))";
}

function areaAvailBg(itemGroupID: string) {
  const { free, total } = areaAvailData(itemGroupID);
  if (total === 0) {
    return "transparent";
  }
  if (free === 0) {
    return "rgba(var(--v-theme-error), 0.3)";
  }
  if (free / total <= 0.5) {
    return "rgba(var(--v-theme-warning), 0.18)";
  }
  return "rgba(var(--v-theme-success), 0.08)";
}

function canDrillInto(itemGroupID: string) {
  return Boolean(itemGroupMap.value.get(itemGroupID)?.floorPlan);
}

function shouldDrillIntoItemGroup(itemGroupID?: string) {
  return Boolean(itemGroupID && canDrillInto(itemGroupID));
}

function handleAreaClick(itemGroupID: string) {
  const itemGroup = itemGroupMap.value.get(itemGroupID);
  if (!itemGroup || !itemGroup.floorPlan) {
    return;
  }

  drilledInto.value = {
    itemGroupId: itemGroupID,
    name: itemGroup.name,
    floorPlan: itemGroup.floorPlan,
  };
}

function drillBack() {
  drilledInto.value = null;
}

function resolveBookingItemGroupId(itemID: string) {
  return itemToGroupMap.value.get(itemID) || activeItemGroupId.value || "";
}

function initializeBookingSelection() {
  const nextSelections: Record<string, boolean> = {};
  for (const day of bookingDayOptions.value) {
    nextSelections[day.date] = day.date === selectedDate.value && !day.past;
  }
  bookingDaySelections.value = nextSelections;
}

function handleDeskClick(itemID: string, label: string) {
  const parentItemGroupId = itemToGroupMap.value.get(itemID);
  if (shouldDrillIntoItemGroup(parentItemGroupId)) {
    handleAreaClick(parentItemGroupId!);
    return;
  }

  void requestBooking(itemID, label);
}

function handleClose() {
  if (drilledInto.value) {
    drillBack();
    return;
  }

  emit("close");
}

async function loadBookingDayAvailability() {
  if (!pendingBooking.value || !pendingBooking.value.itemGroupId) {
    bookingDayInfoMap.value = new Map();
    return;
  }

  bookingDayAvailabilityLoading.value = true;
  let capturedEquipment: string[] = [];
  let capturedWarning: string | undefined;
  try {
    const results = await Promise.all(
      bookingDayOptions.value.map(async (day) => {
        if (day.past) {
          return { date: day.date, info: { status: "unavailable" as const } };
        }

        const response = await fetchItems(
          pendingBooking.value!.itemGroupId,
          day.date,
        );
        const item = (
          response.data as JsonApiResource<ItemAttributes>[]
        ).find((entry) => entry.id === pendingBooking.value!.itemId);

        if (!item) {
          return { date: day.date, info: { status: "unavailable" as const } };
        }

        if (capturedEquipment.length === 0 && item.attributes.equipment?.length) {
          capturedEquipment = item.attributes.equipment;
        }
        if (capturedWarning === undefined && item.attributes.warning) {
          capturedWarning = item.attributes.warning;
        }

        if (item.attributes.availability === "available") {
          return { date: day.date, info: { status: "free" as const } };
        }
        if (item.attributes.booked_by_me) {
          return {
            date: day.date,
            info: {
              status: "booked-by-me" as const,
              bookerName: item.attributes.booker_name,
            },
          };
        }
        return {
          date: day.date,
          info: {
            status: "booked-by-other" as const,
            bookerName: item.attributes.booker_name,
          },
        };
      }),
    );

    const nextInfoMap = new Map<string, BookingDayInfo>();
    for (const result of results) {
      nextInfoMap.set(result.date, result.info);
    }
    bookingDayInfoMap.value = nextInfoMap;
    bookingItemEquipment.value = capturedEquipment;
    bookingItemWarning.value = capturedWarning;

    bookingDaySelections.value = Object.fromEntries(
      Object.entries(bookingDaySelections.value).map(([date, selected]) => [
        date,
        selected && nextInfoMap.get(date)?.status === "free",
      ]),
    );
  } catch {
    const fallbackMap = new Map<string, BookingDayInfo>();
    for (const day of bookingDayOptions.value) {
      fallbackMap.set(
        day.date,
        day.date === selectedDate.value
          ? { status: "free" }
          : { status: "unavailable" },
      );
    }
    bookingDayInfoMap.value = fallbackMap;
    errorSnackbarText.value = "Unable to check availability for the selected week.";
    showErrorSnackbar.value = true;
  } finally {
    bookingDayAvailabilityLoading.value = false;
  }
}

function isBookingDayDisabled(date: string, past: boolean) {
  if (past || bookingDayAvailabilityLoading.value) return true;
  const info = bookingDayInfoMap.value.get(date);
  return !info || info.status !== "free";
}

function bookingDayTitle(date: string) {
  const info = bookingDayInfoMap.value.get(date);
  if (info?.status === "booked-by-other") {
    return `Booked by ${info.bookerName || "someone else"}`;
  }
  if (info?.status === "booked-by-me") {
    return "Already booked by you";
  }
  if (info?.status === "unavailable") {
    return "Not available";
  }
  return getFullDayLabel(date);
}

async function requestBooking(itemID: string, label: string) {
  if (!selectedDate.value) {
    return;
  }

  pendingBooking.value = {
    itemId: itemID,
    label,
    itemGroupId: resolveBookingItemGroupId(itemID),
  };
  initializeBookingSelection();
  showBookingDialog.value = true;
  await loadBookingDayAvailability();
}

function toggleBookingDaySelection(date: string) {
  const day = bookingDayOptions.value.find((entry) => entry.date === date);
  if (!day || isBookingDayDisabled(date, day.past)) {
    return;
  }

  bookingDaySelections.value = {
    ...bookingDaySelections.value,
    [date]: !bookingDaySelections.value[date],
  };
}

function parseConflictDates(conflicts: string[]) {
  return conflicts
    .map((conflict) => conflict.split(":")[0]?.trim() || "")
    .filter((date) => date !== "");
}

function formatConflictMessage(conflicts: string[], allFailed: boolean) {
  const dates = parseConflictDates(conflicts);
  if (dates.length === 1) {
    return `The selected item is already booked on ${formatReadableDate(
      dates[0]!,
    )}.`;
  }

  if (allFailed) {
    return "The selected item is already booked on the selected days.";
  }

  return "Some selected days were already booked and were skipped.";
}

function formatBookingError(err: unknown, multiDay: boolean) {
  if (err instanceof ApiError && err.status === 409) {
    return multiDay
      ? "Some selected days are already booked."
      : "The selected item is already booked.";
  }

  if (err instanceof ApiError && err.detail) {
    return err.detail;
  }

  return "The booking could not be completed.";
}

async function confirmPendingBooking() {
  if (!pendingBooking.value || bookingInProgress.value) {
    return;
  }

  const bookingDates = selectedBookingDates.value;
  if (bookingDates.length === 0) {
    errorSnackbarText.value = "Select at least one day to continue.";
    showErrorSnackbar.value = true;
    return;
  }

  bookingInProgress.value = true;
  try {
    let createdBookingIds: string[] = [];
    let conflicts: string[] = [];
    const multiDay = bookingDates.length > 1;

    if (!multiDay) {
      const response = await createBooking(
        pendingBooking.value.itemId,
        bookingDates[0]!,
      );
      createdBookingIds = [response.data.id];
    } else {
      const response = await createMultiDayBooking(
        pendingBooking.value.itemId,
        bookingDates,
      );
      createdBookingIds = response.created.map((resource) => resource.id);
      conflicts = response.conflicts || [];
    }

    if (createdBookingIds.length === 0) {
      errorSnackbarText.value =
        conflicts.length > 0
          ? formatConflictMessage(conflicts, true)
          : "The booking could not be completed.";
      showErrorSnackbar.value = true;
      return;
    }

    lastBooking.value = {
      bookingIds: createdBookingIds,
      itemId: pendingBooking.value.itemId,
      label: pendingBooking.value.label,
    };
    bookingSnackbarText.value =
      createdBookingIds.length === 1
        ? `${pendingBooking.value.label} booked successfully.`
        : `${pendingBooking.value.label} booked for ${createdBookingIds.length} days.`;
    showBookingSnackbar.value = true;
    showBookingDialog.value = false;
    pendingBooking.value = null;
    bookingDaySelections.value = {};
    await refreshAvailability();

    if (conflicts.length > 0) {
      errorSnackbarText.value = formatConflictMessage(conflicts, false);
      showErrorSnackbar.value = true;
    }
  } catch (err) {
    errorSnackbarText.value = formatBookingError(err, bookingDates.length > 1);
    showErrorSnackbar.value = true;
  } finally {
    bookingInProgress.value = false;
  }
}

async function undoLastBooking() {
  if (!lastBooking.value) {
    return;
  }

  undoInProgress.value = true;
  const booking = lastBooking.value;
  try {
    for (const bookingId of booking.bookingIds) {
      await cancelBooking(bookingId);
    }
    showBookingSnackbar.value = false;
    lastBooking.value = null;
    await refreshAvailability();
  } catch {
    errorSnackbarText.value = "Undo failed. The booking is still active.";
    showErrorSnackbar.value = true;
  } finally {
    undoInProgress.value = false;
  }
}

function clampZoom(value: number) {
  return Math.min(2.5, Math.max(0.75, value));
}

function adjustZoom(delta: number) {
  zoomScale.value = clampZoom(zoomScale.value + delta);
}

function onWheelZoom(event: WheelEvent) {
  if (!event.ctrlKey) {
    return;
  }

  event.preventDefault();
  const delta = event.deltaY < 0 ? 0.1 : -0.1;
  adjustZoom(delta);
}

function getTouchDistance(touches: TouchList) {
  if (touches.length < 2) {
    return 0;
  }

  const dx = touches[0]!.clientX - touches[1]!.clientX;
  const dy = touches[0]!.clientY - touches[1]!.clientY;
  return Math.hypot(dx, dy);
}

function onTouchStart(event: TouchEvent) {
  if (event.touches.length === 2) {
    pinchState.value = {
      startDistance: getTouchDistance(event.touches),
      startScale: zoomScale.value,
    };
  }
}

function onTouchMove(event: TouchEvent) {
  if (event.touches.length !== 2 || !pinchState.value) {
    return;
  }

  const distance = getTouchDistance(event.touches);
  if (pinchState.value.startDistance <= 0) {
    return;
  }

  const factor = distance / pinchState.value.startDistance;
  zoomScale.value = clampZoom(pinchState.value.startScale * factor);
}

function onTouchEnd() {
  pinchState.value = null;
}

function updateViewport() {
  if (typeof window.matchMedia !== "function") {
    isCompactViewport.value = false;
    return;
  }

  const narrow = window.matchMedia("(max-width: 768px)").matches;
  const short = window.matchMedia("(max-height: 500px)").matches;
  isCompactViewport.value = narrow || short;
}

function handleResize() {
  updateViewport();
}

onMounted(() => {
  updateViewport();
  window.addEventListener("resize", handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
  pinchState.value = null;
});
</script>

<style scoped>
.fp-root {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-height: 0;
  height: 100%;
}

.fp-header {
  align-items: center;
}

.fp-weekday-selector {
  display: flex;
  gap: 0.5rem;
}

.fp-weekday-selector--row {
  flex-wrap: wrap;
}

.fp-weekday-selector--grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(72px, 1fr));
}


.fp-scroll-shell {
  overflow: auto;
  max-height: calc(100vh - 260px);
  touch-action: pan-x pan-y;
}

.fp-scroll-shell--compact {
  flex: 1;
  min-height: 0;
  max-height: none;
}

.fp-zoom-layer {
  display: inline-block;
  min-width: 100%;
}

.fp-fit-container {
  border: 1px solid rgb(var(--v-theme-outline));
  border-radius: 8px;
  display: inline-flex;
  justify-content: center;
  background: rgb(var(--v-theme-surface));
}

.fp-container {
  position: relative;
  display: inline-block;
}

.fp-content-loading .fp-item {
  opacity: 0.4;
}

.fp-image-fit {
  display: block;
  max-width: 100%;
  width: auto;
  height: auto;
  user-select: none;
}

.fp-item {
  position: absolute;
  border: 2px solid transparent;
  transition:
    background-color 0.2s,
    border-color 0.2s,
    opacity 0.2s,
    transform 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fp-item--free {
  border-color: rgb(var(--v-theme-success));
  background-color: transparent;
  cursor: pointer;
}

.fp-item--free:hover {
  background-color: rgba(var(--v-theme-success), 0.12);
}

.fp-item--busy {
  border-color: rgb(var(--v-theme-error));
  background-color: rgba(var(--v-theme-error), 0.3);
}

.fp-item--area {
  flex-direction: column;
}

.fp-item--clickable {
  cursor: pointer;
}

.fp-item--clickable:hover {
  transform: scale(1.01);
}

.fp-item-label,
.fp-item-fraction {
  font-size: 0.85rem;
  font-weight: 700;
  color: rgb(var(--v-theme-on-surface));
  text-shadow:
    -1px -1px 0 white,
    1px -1px 0 white,
    -1px 1px 0 white,
    1px 1px 0 white;
  pointer-events: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: calc(100% - 6px);
}

.fp-label-hidden .fp-item-label,
.fp-label-hidden .fp-item-fraction {
  display: none;
}

.fp-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.25rem 0;
}


.fp-booking-dialog-card {
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.fp-booking-day-row {
  display: grid;
  grid-template-columns: 40px 1fr auto;
  align-items: center;
  padding: 4px 0;
}

.fp-booking-day-row--past {
  opacity: 0.5;
}

.fp-booking-day-checkbox {
  min-height: 44px;
}

.fp-booking-dialog-actions {
  position: sticky;
  bottom: 0;
  z-index: 2;
  background: rgb(var(--v-theme-surface));
  border-top: 1px solid rgba(var(--v-theme-outline), 0.5);
  padding: 1rem;
}

.fp-drill-enter-active,
.fp-drill-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.fp-drill-enter-from,
.fp-drill-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

@media (max-width: 768px), (max-height: 500px) {
  .fp-header {
    align-items: stretch;
  }

  .fp-zoom-controls {
    width: 100%;
    justify-content: space-between;
  }

  .fp-footer {
    position: sticky;
    bottom: 0;
    z-index: 3;
    background: rgb(var(--v-theme-surface));
    border-top: 1px solid rgba(var(--v-theme-outline), 0.5);
    padding: 0.75rem 0 0.25rem;
  }
}
</style>

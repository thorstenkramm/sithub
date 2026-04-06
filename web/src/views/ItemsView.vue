<template>
  <div class="page-container">
    <PageHeader
      title=""
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <v-btn
          v-if="activeItemGroupId"
          variant="text"
          size="small"
          :to="{
            name: 'item-group-bookings',
            params: { itemGroupId: activeItemGroupId! },
            query: breadcrumbAreaId ? { areaId: breadcrumbAreaId } : {}
          }"
          data-cy="view-item-group-bookings"
        >
          {{ $t('items.viewItemGroupBookings') }}
        </v-btn>
      </template>
    </PageHeader>

    <!-- Date Selection & Booking Options -->
    <v-card class="mb-6">
      <v-card-text>
        <!-- Booking Mode Toggle -->
        <v-btn-toggle
          v-model="bookingMode"
          mandatory
          density="compact"
          class="mb-4"
          data-cy="booking-mode-toggle"
        >
          <v-btn value="day" data-cy="mode-day-btn">{{ $t('items.day') }}</v-btn>
          <v-btn value="week" data-cy="mode-week-btn">{{ $t('items.week') }}</v-btn>
        </v-btn-toggle>

        <div class="d-flex flex-wrap align-end ga-4 mb-4">
          <!-- Day mode: date picker -->
          <DatePickerField
            v-if="bookingMode === 'day'"
            v-model="selectedDate"
            :label="$t('items.bookingDate')"
            :min="todayDate"
            :max="maxBookingDate"
            density="compact"
            hide-details
            data-cy="items-date"
            style="max-width: 320px;"
          />

          <!-- Week mode: week selector -->
          <v-select
            v-if="bookingMode === 'week'"
            v-model="selectedWeek"
            :items="weekOptions"
            item-title="label"
            item-value="value"
            :label="$t('items.calendarWeek')"
            density="compact"
            hide-details
            data-cy="week-selector"
            style="max-width: 320px;"
          />

          <v-btn
            v-if="itemGroupFloorPlan"
            variant="outlined"
            density="compact"
            prepend-icon="$map"
            data-cy="item-group-floor-plan-btn"
            @click="showItemGroupFloorPlanDialog = true"
          >
            {{ $t('items.floorPlan') }}
          </v-btn>
        </div>

        <!-- Booking Type Selection -->
        <v-radio-group v-model="bookingType" inline density="compact" class="mb-2" hide-details>
          <v-radio :label="$t('items.bookForMyself')" value="self" data-cy="book-self-radio" />
          <v-radio :label="$t('items.bookForColleague')" value="colleague" data-cy="book-colleague-radio" />
        </v-radio-group>

        <!-- Colleague Fields -->
        <v-expand-transition>
          <div v-if="bookingType === 'colleague'" class="mt-4">
            <v-autocomplete
              v-model="selectedColleagueId"
              :items="usersList"
              item-title="displayName"
              item-value="id"
              :label="$t('items.selectColleague')"
              density="compact"
              :loading="usersLoading"
              clearable
              data-cy="colleague-select"
              style="max-width: 360px;"
            />
          </div>
        </v-expand-transition>

        <!-- Equipment Filter -->
        <div class="d-flex align-center ga-2 mt-4" style="max-width: 420px;">
          <v-combobox
            v-model="equipmentFilter"
            :items="savedFilterItems"
            :label="$t('items.filterEquipment')"
            density="compact"
            hide-details
            clearable
            prepend-inner-icon="$filterOutline"
            data-cy="equipment-filter-input"
          />
          <v-tooltip :text="isCurrentFilterSaved ? $t('items.deleteSavedFilter') : $t('items.saveFilter')" location="top">
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                icon
                variant="text"
                size="small"
                :data-cy="isCurrentFilterSaved ? 'equipment-filter-delete' : 'equipment-filter-save'"
                :aria-label="isCurrentFilterSaved ? $t('items.deleteSavedFilter') : $t('items.saveFilter')"
                @click="toggleSaveFilter"
              >
                <v-icon>{{ isCurrentFilterSaved ? '$delete' : 'mdi-content-save' }}</v-icon>
              </v-btn>
            </template>
          </v-tooltip>
          <v-btn
            icon
            variant="text"
            size="small"
            data-cy="equipment-filter-info"
            :aria-label="$t('items.equipmentFilterHelp')"
            @click="showFilterHelp = true"
          >
            <v-icon>$info</v-icon>
          </v-btn>
        </div>

      </v-card-text>
    </v-card>

    <!-- Loading State -->
    <LoadingState v-if="itemsLoading || weekDataLoading" type="cards" :count="6" data-cy="items-loading" />

    <!-- Error State -->
    <v-alert v-else-if="itemsErrorMessage" type="error" class="mb-4" data-cy="items-error">
      {{ itemsErrorMessage }}
    </v-alert>

    <!-- Empty State -->
    <EmptyState
      v-else-if="bookingMode === 'day' && !items.length"
      :title="$t('items.emptyTitle')"
      :message="$t('items.emptyMessage')"
      icon="$desk"
      data-cy="items-empty"
    />
    <EmptyState
      v-else-if="bookingMode === 'week' && !weekItems.length"
      :title="$t('items.emptyTitle')"
      :message="$t('items.emptyMessage')"
      icon="$desk"
      data-cy="items-empty"
    />

    <!-- Items Grid (Day mode) -->
    <div v-else-if="bookingMode === 'day'" class="card-grid" data-cy="items-list">
      <div
        v-for="entry in items"
        :key="entry.id"
        :class="['item-filter-wrapper', { 'item-expanded': expandedDayTiles.has(entry.id) }]"
      >
        <div
          v-if="isItemFilteredOut(entry.attributes.equipment || []) || entry.attributes.reserved"
          class="item-filtered-overlay"
          :data-cy="entry.attributes.reserved ? 'item-reserved' : 'equipment-not-available'"
        >
          <span class="text-body-2 font-weight-medium">
            {{ getOverlayLabel(entry.attributes.reserved === true) }}
          </span>
        </div>
        <v-card
          :class="[
            'item-card',
            { 'item-available': entry.attributes.availability === 'available' },
            { 'item-occupied': entry.attributes.availability === 'occupied' },
            { 'item-filtered-out': isItemFilteredOut(entry.attributes.equipment || []) || entry.attributes.reserved }
          ]"
          :title="entry.attributes.reserved ? $t('items.reservedTooltip') : undefined"
          data-cy="item-entry"
          :data-cy-item-id="entry.id"
          :data-cy-availability="entry.attributes.availability"
        >
        <v-card-item>
          <template #prepend>
            <v-avatar
              :color="entry.attributes.availability === 'available' ? 'success' : 'error'"
              variant="tonal"
              size="48"
            >
              <v-icon size="24">{{ resolveItemIcon(entry.attributes.icon) }}</v-icon>
            </v-avatar>
          </template>
          <!-- Line 1: Item name -->
          <v-card-title>
            <v-tooltip location="top" :disabled="!dayNameTruncatedMap[entry.id]">
              <template #activator="{ props: nameTooltip }">
                <span v-bind="nameTooltip" class="item-name-shell">
                  <span :ref="setDayNameRef(entry.id)" class="item-name">{{ getDayNameLabel(entry.id, entry.attributes.name) }}</span>
                  <span :ref="setDayNameMeasureRef(entry.id)" class="item-name-measure" aria-hidden="true">{{ entry.attributes.name }}</span>
                </span>
              </template>
              {{ entry.attributes.name }}
            </v-tooltip>
          </v-card-title>
          <!-- Line 2: Status chip + warning + chevron -->
          <div class="d-flex align-center ga-2 mt-1 px-4">
            <StatusChip
              :status="entry.attributes.availability === 'available' ? 'available' : 'booked'"
              size="x-small"
              data-cy="item-status"
            />
            <v-tooltip
              v-if="entry.attributes.warning"
              location="top"
              content-class="warning-tooltip"
            >
              <template #activator="{ props: tooltipProps }">
                <v-btn
                  v-bind="tooltipProps"
                  icon
                  variant="text"
                  size="x-small"
                  color="warning"
                  data-cy="folded-warning-icon"
                >
                  <v-icon size="18">$warning</v-icon>
                </v-btn>
              </template>
              {{ entry.attributes.warning }}
            </v-tooltip>
            <v-spacer />
            <v-btn
              icon
              variant="text"
              size="small"
              data-cy="day-tile-chevron"
              :aria-label="`Toggle details for ${entry.attributes.name}`"
              :aria-expanded="expandedDayTiles.has(entry.id)"
              @click="toggleDayTileExpansion(entry.id)"
            >
              <v-icon>
                {{ expandedDayTiles.has(entry.id) ? '$chevronDown' : '$chevronLeft' }}
              </v-icon>
            </v-btn>
          </div>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Equipment (hidden on folded booked tiles) -->
          <div
            v-if="entry.attributes.equipment?.length
              && (entry.attributes.availability === 'available' || expandedDayTiles.has(entry.id))"
            class="mb-2"
            data-cy="item-equipment"
          >
            <div class="text-caption text-medium-emphasis mb-1">{{ $t('items.equipment') }}</div>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="equip in entry.attributes.equipment"
                :key="equip"
                size="x-small"
                variant="outlined"
              >
                {{ equip }}
              </v-chip>
            </div>
          </div>

          <!-- Warning (only shown when tile is expanded) -->
          <v-alert
            v-if="entry.attributes.warning && expandedDayTiles.has(entry.id)"
            type="warning"
            variant="tonal"
            density="compact"
            class="mt-2"
            data-cy="item-warning"
          >
            {{ entry.attributes.warning }}
          </v-alert>

          <!-- Booker name -->
          <div
            v-if="entry.attributes.availability === 'occupied' && entry.attributes.booker_name"
            class="text-body-2 text-medium-emphasis mt-2"
            data-cy="item-booker"
          >
            <v-icon size="14" class="mr-1">$user</v-icon>
            {{ entry.attributes.booker_name }}
          </div>

          <!-- Booking note -->
          <div
            v-if="entry.attributes.availability === 'occupied' && entry.attributes.note"
            class="d-flex align-center ga-1 mt-1 text-body-2 text-medium-emphasis"
            data-cy="item-note"
          >
            <v-icon size="14">$textBoxOutline</v-icon>
            <span :ref="setNoteRef(entry.id)" class="note-text">{{ entry.attributes.note }}</span>
            <v-btn
              v-if="noteTruncatedMap[entry.id]"
              icon
              size="x-small"
              variant="text"
              data-cy="item-note-expand"
              @click="expandedNote = entry.attributes.note"
            >
              <v-icon size="14">$arrowExpand</v-icon>
            </v-btn>
          </div>
        </v-card-text>

        <v-card-actions class="px-4 pb-4 ga-2" data-cy="day-item-actions">
          <v-btn
            v-if="entry.attributes.availability === 'available'"
            color="primary"
            variant="flat"
            class="flex-grow-1"
            :loading="bookingItemId === entry.id"
            :disabled="bookingItemId !== null || cancelingBookingId !== null"
            data-cy="book-item-btn"
            @click="bookItem(entry.id)"
          >
            {{ $t('items.book') }}
          </v-btn>
          <v-btn
            v-else-if="authStore.isAdmin && entry.attributes.booking_id"
            color="error"
            variant="tonal"
            class="flex-grow-1"
            :loading="cancelingBookingId === entry.attributes.booking_id"
            :disabled="bookingItemId !== null || cancelingBookingId !== null"
            data-cy="admin-cancel-btn"
            @click="adminCancelBooking(entry.attributes.booking_id!)"
          >
            {{ $t('items.cancelBooking') }}
          </v-btn>
          <div v-else class="py-2 flex-grow-1" />
          <v-spacer />
          <v-btn
            icon
            variant="text"
            size="small"
            data-cy="item-favorite-heart"
            @click.stop="handleToggleItemFav(entry.id, entry.attributes.name)"
          >
            <v-icon size="18" :color="isItemFav(entry.id) ? 'error' : undefined">
              {{ isItemFav(entry.id) ? '$heart' : '$heartOutline' }}
            </v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
      </div>
    </div>

    <!-- Items Grid (Week mode) -->
    <div v-else-if="bookingMode === 'week' && weekItems.length" class="card-grid" data-cy="week-items-list">
      <div
        v-for="item in weekItems"
        :key="item.id"
        :class="['item-filter-wrapper', { 'item-expanded': expandedWeekTiles.has(item.id) }]"
      >
        <div
          v-if="isItemFilteredOut(getWeekItemEquipment(item.id)) || isWeekItemReserved(item.id)"
          class="item-filtered-overlay"
          :data-cy="isWeekItemReserved(item.id) ? 'item-reserved' : 'equipment-not-available'"
        >
          <span class="text-body-2 font-weight-medium">{{ getOverlayLabel(isWeekItemReserved(item.id)) }}</span>
        </div>
        <v-card
          :class="['item-card', { 'item-filtered-out': isItemFilteredOut(getWeekItemEquipment(item.id)) || isWeekItemReserved(item.id) }]"
          :title="isWeekItemReserved(item.id) ? $t('items.reservedTooltip') : undefined"
          data-cy="week-item-entry"
          :data-cy-item-name="item.name"
          :data-cy-item-id="item.id"
        >
        <v-card-item>
          <template #prepend>
            <v-avatar :color="getWeekItemAvatarColor(item.id)" variant="tonal" size="48">
              <v-icon size="24">{{ resolveItemIcon(item.icon) }}</v-icon>
            </v-avatar>
          </template>
          <!-- Line 1: Item name -->
          <v-card-title>
            <v-tooltip location="top" :disabled="!weekNameTruncatedMap[item.id]">
              <template #activator="{ props: nameTooltip }">
                <span v-bind="nameTooltip" class="item-name-shell">
                  <span :ref="setWeekNameRef(item.id)" class="item-name">{{ getWeekNameLabel(item.id, item.name) }}</span>
                  <span :ref="setWeekNameMeasureRef(item.id)" class="item-name-measure" aria-hidden="true">{{ item.name }}</span>
                </span>
              </template>
              {{ item.name }}
            </v-tooltip>
          </v-card-title>
          <!-- Line 2: Availability + warning + chevron -->
          <div class="d-flex align-center ga-2 mt-1 px-4">
            <v-chip
              size="x-small"
              :color="getWeekItemStatusColor(item.id)"
              variant="tonal"
              data-cy="week-item-availability"
            >
              <v-icon start size="14">{{ getWeekItemStatusIcon(item.id) }}</v-icon>
              {{ getWeekItemStatusLabel(item.id) }}
            </v-chip>
            <v-tooltip
              v-if="getWeekItemAttributes(item.id).warning"
              location="top"
              content-class="warning-tooltip"
            >
              <template #activator="{ props: tooltipProps }">
                <v-btn
                  v-bind="tooltipProps"
                  icon
                  variant="text"
                  size="x-small"
                  color="warning"
                  data-cy="week-folded-warning-icon"
                >
                  <v-icon size="18">$warning</v-icon>
                </v-btn>
              </template>
              {{ getWeekItemAttributes(item.id).warning }}
            </v-tooltip>
            <v-spacer />
            <v-btn
              icon
              variant="text"
              size="small"
              data-cy="week-tile-chevron"
              :aria-label="`Toggle details for ${item.name}`"
              :aria-expanded="expandedWeekTiles.has(item.id)"
              @click="toggleWeekTileExpansion(item.id)"
            >
              <v-icon>
                {{ expandedWeekTiles.has(item.id) ? '$chevronDown' : '$chevronLeft' }}
              </v-icon>
            </v-btn>
          </div>
        </v-card-item>

        <v-card-text class="pt-0">
          <!-- Line 3: Equipment (always visible) -->
          <div
            v-if="getWeekItemAttributes(item.id).equipment.length"
            class="mb-2"
            data-cy="week-item-equipment-folded"
          >
            <div class="text-caption text-medium-emphasis mb-1">{{ $t('items.equipment') }}</div>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="equip in getWeekItemAttributes(item.id).equipment"
                :key="equip"
                size="x-small"
                variant="outlined"
              >
                {{ equip }}
              </v-chip>
            </div>
          </div>
          <!-- Line 4: Folded view compact M-F row -->
          <div
            v-if="!expandedWeekTiles.has(item.id)"
            :class="isMobile ? 'week-days-compact' : 'week-days'"
            :style="isMobile ? { gridTemplateColumns: `repeat(${selectedWeekDates.length}, 1fr)` } : undefined"
            data-cy="week-days"
          >
            <div
              v-for="(date, dayIdx) in selectedWeekDates"
              :key="date"
              :class="['week-day-slot', { 'week-day-past': isDateInPast(date) }]"
              :data-cy-weekday="getWeekdayLabel(dayIdx)"
            >
              <span class="week-day-label text-caption font-weight-medium">
                {{ getWeekdayLabel(dayIdx, isMobile, t) }}
              </span>
              <v-checkbox
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :model-value="isWeekDaySelected(item.id, date)"
                hide-details
                density="compact"
                color="success"
                :disabled="isDateInPast(date) || isWeekItemReserved(item.id)"
                class="week-day-checkbox"
                :data-cy="isDateInPast(date) ? 'week-day-checkbox-past' : 'week-day-checkbox'"
                @update:model-value="!isWeekItemReserved(item.id) && toggleWeekDay(item.id, date)"
              />
              <v-checkbox
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :model-value="true"
                hide-details
                density="compact"
                color="primary"
                disabled
                class="week-day-checkbox"
                data-cy="week-day-mine"
              />
              <v-checkbox
                v-else
                :model-value="true"
                hide-details
                density="compact"
                disabled
                class="week-day-checkbox"
                :data-cy="getWeekDayStatus(item.id, date) === 'unavailable'
                  ? 'week-day-unavailable' : 'week-day-other'"
              />
              <span
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :class="['week-day-status', 'text-caption', isDateInPast(date) ? 'text-medium-emphasis' : 'text-success']"
              >{{ $t('items.free') }}</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :class="['week-day-status', 'text-caption', isDateInPast(date) ? 'text-medium-emphasis' : 'text-primary']"
              >{{ authStore.userName || $t('items.me') }}</span>
              <v-icon
                v-if="getWeekDayStatus(item.id, date) === 'booked-by-me' && !isDateInPast(date)"
                size="14"
                color="error"
                class="week-cancel-icon"
                data-cy="week-cancel-btn"
                @click.stop="requestWeekCancel(item.id, date)"
              >$cancelCircle</v-icon>
              <template v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'">
                <v-tooltip location="top">
                  <template #activator="{ props: tooltipProps }">
                    <span
                      v-bind="tooltipProps"
                      class="week-day-status text-caption text-error"
                    >{{ getBookerInitials(item.id, date) }}</span>
                  </template>
                  {{ getWeekDayBooker(item.id, date) }}
                </v-tooltip>
              </template>
            </div>
          </div>
          <!-- Expanded view: one line per day -->
          <div v-else-if="expandedWeekTiles.has(item.id)" data-cy="week-days-expanded">
            <div
              v-for="(date, dayIdx) in selectedWeekDates"
              :key="date"
              :class="['week-day-expanded', { 'week-day-past': isDateInPast(date) }]"
              :data-cy-weekday="getWeekdayLabel(dayIdx)"
            >
              <v-checkbox
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :model-value="isWeekDaySelected(item.id, date)"
                hide-details
                density="compact"
                color="success"
                :disabled="isDateInPast(date) || isWeekItemReserved(item.id)"
                class="week-day-checkbox"
                :data-cy="isDateInPast(date) ? 'week-day-checkbox-past' : 'week-day-checkbox'"
                @update:model-value="!isWeekItemReserved(item.id) && toggleWeekDay(item.id, date)"
              />
              <v-checkbox
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :model-value="true"
                hide-details
                density="compact"
                color="primary"
                disabled
                class="week-day-checkbox"
                data-cy="week-day-mine"
              />
              <v-checkbox
                v-else
                :model-value="true"
                hide-details
                density="compact"
                disabled
                class="week-day-checkbox"
                :data-cy="getWeekDayStatus(item.id, date) === 'unavailable'
                  ? 'week-day-unavailable' : 'week-day-other'"
              />
              <span class="text-body-2 font-weight-medium week-day-expanded-label">
                {{ getFullDayLabel(date, dayIdx) }}
              </span>
              <span
                v-if="getWeekDayStatus(item.id, date) === 'free'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-success']"
              >{{ $t('items.free') }}</span>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-me'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-primary']"
              >{{ authStore.userName || $t('items.me') }}</span>
              <v-icon
                v-if="getWeekDayStatus(item.id, date) === 'booked-by-me' && !isDateInPast(date)"
                size="14"
                color="error"
                class="ml-1 week-cancel-icon"
                data-cy="week-cancel-btn"
                @click.stop="requestWeekCancel(item.id, date)"
              >$cancelCircle</v-icon>
              <span
                v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'"
                :class="['text-body-2', isDateInPast(date) ? 'text-medium-emphasis' : 'text-error']"
              >{{ getWeekDayBooker(item.id, date) }}</span>
            </div>
            <!-- Equipment -->
            <div
              v-if="getWeekItemAttributes(item.id).equipment.length"
              class="mt-3"
              data-cy="week-item-equipment"
            >
              <div class="text-caption text-medium-emphasis mb-1">{{ $t('items.equipment') }}</div>
              <div class="d-flex flex-wrap ga-1">
                <v-chip
                  v-for="equip in getWeekItemAttributes(item.id).equipment"
                  :key="equip"
                  size="x-small"
                  variant="outlined"
                >
                  {{ equip }}
                </v-chip>
              </div>
            </div>
            <!-- Warning -->
            <v-alert
              v-if="getWeekItemAttributes(item.id).warning"
              type="warning"
              variant="tonal"
              density="compact"
              class="mt-2"
              data-cy="week-item-warning"
            >
              {{ getWeekItemAttributes(item.id).warning }}
            </v-alert>
          </div>
        </v-card-text>
        <v-card-actions class="px-4 pb-4 pt-0" data-cy="week-item-actions">
          <v-spacer />
          <v-btn
            icon
            variant="text"
            size="small"
            data-cy="week-item-favorite-heart"
            @click.stop="handleToggleItemFav(item.id, item.name)"
          >
            <v-icon size="18" :color="isItemFav(item.id) ? 'error' : undefined">
              {{ isItemFav(item.id) ? '$heart' : '$heartOutline' }}
            </v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
      </div>
    </div>

    <!-- Confirm Booking Button (Week mode) — sticky footer -->
    <div v-if="bookingMode === 'week' && weekSelectionCount > 0" class="week-book-footer" data-cy="week-confirm-section">
      <v-btn
        color="primary"
        variant="flat"
        size="large"
        block
        :loading="weekBookingInProgress"
        data-cy="week-confirm-btn"
        @click="submitWeekBookings"
      >
        {{ $t('items.bookDays', { count: weekSelectionCount }, weekSelectionCount) }}
      </v-btn>
    </div>

    <!-- Week Booking Results -->
    <v-card v-if="weekBookingResults.length" class="mt-4" data-cy="week-booking-results">
      <v-card-title>{{ $t('items.bookingResults') }}</v-card-title>
      <v-card-text>
        <div v-for="result in weekBookingResults" :key="result.date + result.itemName" class="d-flex align-center ga-2 mb-1">
          <v-icon :color="result.success ? 'success' : 'error'" size="18">
            {{ result.success ? '$success' : '$warning' }}
          </v-icon>
          <span class="text-body-2">
            {{ result.itemName }} - {{ result.dayLabel }}{{ result.success ? '' : ': ' + result.error }}
          </span>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" data-cy="week-results-close" @click="weekBookingResults = []">{{ $t('common.close') }}</v-btn>
      </v-card-actions>
    </v-card>

    <!-- Add Note Dialog (after booking) -->
    <v-dialog v-model="showPostBookingNoteDialog" max-width="500">
      <v-card>
        <v-card-title>{{ $t('items.addNoteTitle') }}</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="noteText"
            :label="$t('items.noteLabel')"
            :counter="500"
            :maxlength="500"
            rows="3"
            auto-grow
            data-cy="post-booking-note-input"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showPostBookingNoteDialog = false">{{ $t('common.cancel') }}</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :loading="savingNote"
            data-cy="post-booking-note-save"
            @click="saveNoteAfterBooking"
          >
            {{ $t('common.save') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note view dialog (desktop) -->
    <v-dialog v-if="!useBottomSheet" v-model="showItemNoteDialog" max-width="500">
      <v-card>
        <v-card-title>{{ $t('items.bookingNote') }}</v-card-title>
        <v-card-text data-cy="item-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showItemNoteDialog = false">{{ $t('common.close') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Note view bottom sheet (mobile) -->
    <v-bottom-sheet v-else v-model="showItemNoteDialog">
      <v-card>
        <v-card-title>{{ $t('items.bookingNote') }}</v-card-title>
        <v-card-text data-cy="item-note-dialog-text">{{ expandedNote }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showItemNoteDialog = false">{{ $t('common.close') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-bottom-sheet>

    <!-- Item Group Floor Plan Dialog -->
    <v-dialog
      v-model="showItemGroupFloorPlanDialog"
      max-width="1100"
      persistent
      :fullscreen="isCompactFloorPlanViewport"
      data-cy="item-group-floor-plan-dialog"
    >
      <v-card class="floor-plan-dialog-card">
        <v-card-text class="floor-plan-dialog-body">
          <InteractiveFloorPlan
            v-if="itemGroupFloorPlan"
            :floor-plan="itemGroupFloorPlan"
            :title="itemGroupName || 'Floor Plan'"
            :week-label="weekOptions.find(o => o.value === selectedWeek)?.label || ''"
            :week-dates="selectedWeekDates"
            :item-group-id="activeItemGroupId || ''"
            @close="showItemGroupFloorPlanDialog = false"
          />
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Equipment Filter Help Dialog -->
    <v-dialog v-model="showFilterHelp" max-width="500">
      <v-card>
        <v-card-title>{{ $t('items.filterSyntaxTitle') }}</v-card-title>
        <v-card-text data-cy="equipment-filter-help">
          <p class="mb-3">{{ $t('items.filterSyntaxDescription') }}</p>
          <ul class="mb-3">
            <li>{{ $t('items.filterSyntaxOr') }}</li>
            <li>{{ $t('items.filterSyntaxAnd') }}</li>
            <li>{{ $t('items.filterSyntaxExact') }}</li>
            <li>{{ $t('items.filterSyntaxCase') }}</li>
          </ul>
          <p class="text-caption text-medium-emphasis">{{ $t('items.filterSyntaxExample') }} <code>"27 inch display" + webcam</code></p>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showFilterHelp = false">{{ $t('common.close') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Confirm Cancel Dialog (week view) -->
    <ConfirmDialog
      v-model="showWeekCancelDialog"
      :title="$t('items.cancelBooking')"
      :message="$t('bookings.cancelMessage')"
      :confirm-text="$t('items.cancelBooking')"
      confirm-color="error"
      @confirm="confirmWeekCancel"
    />

    <v-snackbar
      :key="successSnackbarKey"
      v-model="showSuccessSnackbar"
      :timeout="successSnackbarTimeout"
      location="bottom"
      color="success"
      :data-cy="successSnackbarCy"
    >
      <span :data-cy="successSnackbarCy === 'booking-success' ? 'booking-success-text' : undefined">
        {{ successSnackbarMessage }}
      </span>
      <template v-if="successSnackbarActionLabel" #actions>
        <v-btn
          variant="text"
          size="small"
          data-cy="add-note-after-booking"
          @click="handleSuccessSnackbarAction"
        >
          {{ successSnackbarActionLabel }}
        </v-btn>
      </template>
    </v-snackbar>

    <v-snackbar
      v-model="showErrorSnackbar"
      :timeout="6000"
      location="bottom"
      color="error"
      closable
      data-cy="booking-error"
    >
      <span data-cy="booking-error-text">{{ errorSnackbarMessage }}</span>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { ComponentPublicInstance } from 'vue';
import { useRoute } from 'vue-router';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';
import {
  createBooking,
  cancelBooking,
  updateBookingNote,
  fetchMyBookings,
  type BookOnBehalfOptions
} from '../api/bookings';
import { fetchItems } from '../api/items';
import { fetchColleagues } from '../api/users';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import type { ItemAttributes } from '../api/items';
import type { JsonApiResource } from '../api/types';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { useWeekSelector, getWeekdayLabel } from '../composables/useWeekSelector';
import { useWeekendPreference } from '../composables/useWeekendPreference';
import { matchesParsedFilter, parseFilter } from '../composables/useEquipmentFilter';
import { useSavedFilters } from '../composables/useSavedFilters';
import { useDateState } from '../composables/useDateState';
import { useFavorites } from '../composables/useFavorites';
import { getSafeLocalStorage } from '../composables/storage';
import { useAuthStore } from '../stores/useAuthStore';
import { resolveConfiguredIcon } from '../utils/icons';
import { middleTruncate } from '../utils/text';
import { fetchSettings } from '../api/settings';
import { PageHeader, LoadingState, EmptyState, StatusChip, DatePickerField, ConfirmDialog } from '../components';
import InteractiveFloorPlan from '../components/InteractiveFloorPlan.vue';

const { t, locale } = useI18n();
const authStore = useAuthStore();
const items = ref<JsonApiResource<ItemAttributes>[]>([]);
const itemsErrorMessage = ref<string | null>(null);
const errorSnackbarMessage = ref<string | null>(null);
const showErrorSnackbar = computed({
  get: () => errorSnackbarMessage.value !== null,
  set: (v: boolean) => { if (!v) errorSnackbarMessage.value = null; }
});
const lastBookingDetails = ref<{ itemName: string; date: string } | null>(null);
const bookingItemId = ref<string | null>(null);
const cancelingBookingId = ref<string | null>(null);
const { getDay, setDay, resetDayToToday, getWeek, setWeek } = useDateState();
const selectedDate = ref(getDay());
const todayDate = formatDate(new Date());
const maxBookingDate = computed(() => {
  const now = new Date();
  const day = now.getDay();
  const daysUntilMonday = (8 - day) % 7;
  const nextMonday = new Date(now);
  nextMonday.setDate(now.getDate() + daysUntilMonday);
  const maxDate = new Date(nextMonday);
  maxDate.setDate(nextMonday.getDate() + weeksInAdvanced.value * 7 - 1);
  return formatDate(maxDate);
});
const route = useRoute();
const { loading: itemsLoading, run: runItems } = useApi();
const activeItemGroupId = ref<string | null>(null);
const areaName = ref('');
const itemGroupName = ref('');
const bookingType = ref<'self' | 'colleague'>('self');
const selectedColleagueId = ref<string | null>(null);
const usersList = ref<Array<{ id: string; displayName: string }>>([]);
const usersLoading = ref(false);
const lastBookingId = ref<string | null>(null);
const showPostBookingNoteDialog = ref(false);
const noteText = ref('');
const savingNote = ref(false);
const expandedNote = ref('');
const noteTruncatedMap = ref<Record<string, boolean>>({});
const noteElements = new Map<string, HTMLElement>();
const dayNameTruncatedMap = ref<Record<string, boolean>>({});
const weekNameTruncatedMap = ref<Record<string, boolean>>({});
const dayNameElements = new Map<string, HTMLElement>();
const dayNameMeasureElements = new Map<string, HTMLElement>();
const weekNameElements = new Map<string, HTMLElement>();
const weekNameMeasureElements = new Map<string, HTMLElement>();
const isMobile = ref(false);
const isCompactFloorPlanViewport = ref(false);
const useBottomSheet = computed(() => isMobile.value);
const showItemNoteDialog = computed({
  get: () => expandedNote.value !== '',
  set: (v: boolean) => { if (!v) expandedNote.value = ''; }
});

// Equipment filter
const itemGroupFloorPlan = ref<string | null>(null);
const inheritedIcon = ref<string | null>(null);

const resolveItemIcon = (itemIcon: string | undefined) => {
  return resolveConfiguredIcon(itemIcon || inheritedIcon.value, '$desk');
};
const showItemGroupFloorPlanDialog = ref(false);

const equipmentFilter = ref('');
const showFilterHelp = ref(false);
const { comboboxItems: savedFilterItems, saveFilter, deleteFilter, isSavedFilter } = useSavedFilters();
const { isItemFavorite, toggleItemFavorite } = useFavorites();
const successSnackbarMessage = ref<string | null>(null);
const successSnackbarCy = ref('items-success');
const successSnackbarTimeout = ref(3000);
const successSnackbarKey = ref(0);
const successSnackbarActionLabel = ref<string | null>(null);
const successSnackbarActionHandler = ref<(() => void) | null>(null);
const showSuccessSnackbar = computed({
  get: () => successSnackbarMessage.value !== null,
  set: (v: boolean) => {
    if (!v) {
      successSnackbarMessage.value = null;
      successSnackbarCy.value = 'items-success';
      successSnackbarTimeout.value = 3000;
      successSnackbarActionLabel.value = null;
      successSnackbarActionHandler.value = null;
    }
  }
});
const showSuccessFeedback = (
  message: string,
  cy: string,
  options?: { timeout?: number; actionLabel?: string; actionHandler?: () => void }
) => {
  successSnackbarKey.value += 1;
  successSnackbarMessage.value = message;
  successSnackbarCy.value = cy;
  successSnackbarTimeout.value = options?.timeout ?? 3000;
  successSnackbarActionLabel.value = options?.actionLabel ?? null;
  successSnackbarActionHandler.value = options?.actionHandler ?? null;
};
const handleSuccessSnackbarAction = () => {
  successSnackbarActionHandler.value?.();
};
const isItemFav = (itemId: string) =>
  !!activeItemGroupId.value
  && !!getCurrentAreaId()
  && isItemFavorite(getCurrentAreaId(), activeItemGroupId.value, itemId);
const getDayNameLabel = (itemId: string, name: string) =>
  dayNameTruncatedMap.value[itemId] ? middleTruncate(name, 25) : name;
const getWeekNameLabel = (itemId: string, name: string) =>
  weekNameTruncatedMap.value[itemId] ? middleTruncate(name, 25) : name;
const handleToggleItemFav = (itemId: string, itemName: string) => {
  const areaId = getCurrentAreaId();
  const igName = itemGroupName.value || '';
  if (!activeItemGroupId.value || !areaId) {
    return;
  }
  const { added } = toggleItemFavorite({
    areaId,
    itemId,
    itemName,
    itemGroupId: activeItemGroupId.value,
    itemGroupName: igName
  });
  const label = `${igName} ${itemName}`;
  showSuccessFeedback(
    added ? t('items.savedAsFavorite', { name: label }) : t('items.removedFromFavorites', { name: label }),
    'item-favorite-message'
  );
};
const isCurrentFilterSaved = computed(() => !!equipmentFilter.value && isSavedFilter(equipmentFilter.value));
const showFilterFeedback = (message: string) => {
  showSuccessFeedback(message, 'filter-message');
};
const toggleSaveFilter = () => {
  if (!equipmentFilter.value) return;
  if (isCurrentFilterSaved.value) {
    deleteFilter(equipmentFilter.value);
    equipmentFilter.value = '';
    showFilterFeedback(t('items.savedFilterDeleted'));
  } else {
    if (saveFilter(equipmentFilter.value)) {
      showFilterFeedback(t('items.filterSaved'));
    }
  }
};
const parsedEquipmentFilter = computed(() => parseFilter(equipmentFilter.value));

const isItemFilteredOut = (equipment: string[]): boolean => {
  return !matchesParsedFilter(equipment, parsedEquipmentFilter.value);
};

const getWeekItemEquipment = (itemId: string): string[] => {
  for (const dayItems of Object.values(weekData.value)) {
    const item = dayItems.find(i => i.id === itemId);
    if (item?.attributes.equipment?.length) return item.attributes.equipment;
  }
  return [];
};

// Week booking mode
const storage = getSafeLocalStorage();
const bookingMode = ref<'day' | 'week'>(
  (storage?.getItem('sithub_booking_mode') as 'day' | 'week') || 'day'
);
const { showWeekends } = useWeekendPreference();
const weeksInAdvanced = ref(7);
const { weekOptions, selectedWeek, selectedWeekDates } = useWeekSelector(showWeekends, weeksInAdvanced);

// Restore memorized week
const storedWeek = getWeek();
if (weekOptions.value.some(o => o.value === storedWeek)) {
  selectedWeek.value = storedWeek;
}

// Per-day data for week mode: map of date -> items array
const weekData = ref<Record<string, JsonApiResource<ItemAttributes>[]>>({});
const weekDataLoading = ref(false);
const myWeekBookings = ref<Map<string, string>>(new Map());

// Week day selections: Set of "itemId::date" keys
const weekSelections = ref<Set<string>>(new Set());
const weekBookingInProgress = ref(false);

interface WeekBookingResult {
  itemName: string;
  date: string;
  dayLabel: string;
  success: boolean;
  error?: string;
}
const weekBookingResults = ref<WeekBookingResult[]>([]);

// Story 15-1: Week tile expansion
const expandedWeekTiles = ref<Set<string>>(new Set());

const getFullDayLabel = (date: string, fallbackIndex: number): string => {
  const parsed = new Date(`${date}T00:00:00`);
  if (!Number.isNaN(parsed.getTime())) {
    return new Intl.DateTimeFormat(locale.value || undefined, {
      weekday: 'long',
      month: '2-digit',
      day: '2-digit'
    }).format(parsed);
  }
  return getWeekdayLabel(fallbackIndex, false, t) || date;
};


const toggleWeekTileExpansion = (itemId: string) => {
  const next = new Set(expandedWeekTiles.value);
  if (next.has(itemId)) {
    next.delete(itemId);
  } else {
    next.add(itemId);
  }
  expandedWeekTiles.value = next;
};

const getOverlayLabel = (reserved: boolean): string =>
  reserved ? t('items.reserved') : t('items.equipmentNotAvailable');

const weekItemAttributesMap = computed(() => {
  const map = new Map<string, { equipment: string[]; warning?: string; reserved?: boolean }>();
  for (const dayItems of Object.values(weekData.value)) {
    for (const item of dayItems) {
      if (!map.has(item.id)) {
        map.set(item.id, {
          equipment: item.attributes.equipment || [],
          warning: item.attributes.warning,
          reserved: item.attributes.reserved === true
        });
      }
    }
  }
  return map;
});

const getWeekItemAttributes = (itemId: string): { equipment: string[]; warning?: string; reserved?: boolean } =>
  weekItemAttributesMap.value.get(itemId) ?? { equipment: [] };

const isWeekItemReserved = (itemId: string): boolean =>
  getWeekItemAttributes(itemId).reserved === true;

const getWeekItemFreeDays = (itemId: string): number =>
  selectedWeekDates.value.filter(date => getWeekDayStatus(itemId, date) === 'free').length;

const getWeekItemAvatarColor = (itemId: string): string => {
  const free = getWeekItemFreeDays(itemId);
  const total = selectedWeekDates.value.length;
  if (free === total) return 'success';
  if (free === 0) return 'error';
  return 'primary';
};

const getWeekItemStatusColor = (itemId: string): string =>
  isWeekItemReserved(itemId) ? 'warning' : getWeekItemAvatarColor(itemId);

const getWeekItemStatusIcon = (itemId: string): string => {
  if (isWeekItemReserved(itemId)) return '$lock';
  return getWeekItemFreeDays(itemId) === 0 ? '$calendar' : '$success';
};

const getWeekItemStatusLabel = (itemId: string): string => {
  if (isWeekItemReserved(itemId)) return t('items.reserved');
  const freeDays = getWeekItemFreeDays(itemId);
  if (freeDays === 0) return t('status.booked');
  return `${t('status.available')} ${freeDays}/${selectedWeekDates.value.length}`;
};

// Story 15-2: Day tile expansion
const expandedDayTiles = ref<Set<string>>(new Set());

const toggleDayTileExpansion = (itemId: string) => {
  const next = new Set(expandedDayTiles.value);
  if (next.has(itemId)) {
    next.delete(itemId);
  } else {
    next.add(itemId);
  }
  expandedDayTiles.value = next;
};

// Story 15-4: Past date detection
const isDateInPast = (date: string): boolean => date < todayDate;

// Unique items across all days in week mode
const weekItems = computed(() => {
  const itemsMap = new Map<string, { name: string; icon?: string }>();
  for (const dayItems of Object.values(weekData.value)) {
    for (const item of dayItems) {
      if (!itemsMap.has(item.id)) {
        itemsMap.set(item.id, { name: item.attributes.name, icon: item.attributes.icon });
      }
    }
  }
  return Array.from(itemsMap.entries())
    .map(([id, attrs]) => ({ id, name: attrs.name, icon: attrs.icon }))
    .sort((a, b) => a.name.localeCompare(b.name));
});

const weekSelectionCount = computed(() => weekSelections.value.size);

const getWeekSelectionKey = (itemId: string, date: string) => `${itemId}::${date}`;

const getWeekDayStatus = (
  itemId: string,
  date: string
): 'free' | 'booked-by-me' | 'booked-by-other' | 'unavailable' => {
  const dayItems = weekData.value[date];
  if (!dayItems) return 'unavailable';
  const item = dayItems.find(i => i.id === itemId);
  if (!item) return 'unavailable';
  if (item.attributes.availability === 'available') return 'free';
  if (isBookedByMe(itemId, date)) return 'booked-by-me';
  return 'booked-by-other';
};

const getWeekDayBooker = (itemId: string, date: string): string => {
  const dayItems = weekData.value[date];
  if (!dayItems) return t('common.booked');
  const item = dayItems.find(i => i.id === itemId);
  return item?.attributes.booker_name || t('common.booked');
};

const getBookerInitials = (itemId: string, date: string): string => {
  const name = getWeekDayBooker(itemId, date);
  const parts = name.split(' ');
  if (parts.length >= 2 && parts[0] && parts[parts.length - 1]) {
    return (parts[0].charAt(0) + parts[parts.length - 1]!.charAt(0)).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
};

const isWeekDaySelected = (itemId: string, date: string) =>
  weekSelections.value.has(getWeekSelectionKey(itemId, date));

const isBookedByMe = (itemId: string, date: string) =>
  myWeekBookings.value.has(getWeekSelectionKey(itemId, date));

const weekCancellingKey = ref<string | null>(null);
const showWeekCancelDialog = ref(false);
const pendingWeekCancelKey = ref<string | null>(null);

const requestWeekCancel = (itemId: string, date: string) => {
  const key = getWeekSelectionKey(itemId, date);
  if (!myWeekBookings.value.has(key)) return;
  pendingWeekCancelKey.value = key;
  showWeekCancelDialog.value = true;
};

const confirmWeekCancel = async () => {
  if (!pendingWeekCancelKey.value) return;

  const bookingId = myWeekBookings.value.get(pendingWeekCancelKey.value);
  if (!bookingId) return;

  weekCancellingKey.value = pendingWeekCancelKey.value;
  showWeekCancelDialog.value = false;
  try {
    await cancelBooking(bookingId);
    showSuccessFeedback(t('items.bookingCancelledSuccessfully'), 'week-cancel-success');
    if (activeItemGroupId.value) {
      await loadWeekData(activeItemGroupId.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) return;
    errorSnackbarMessage.value =t('items.unableToCancel');
  } finally {
    weekCancellingKey.value = null;
    pendingWeekCancelKey.value = null;
  }
};

const toggleWeekDay = (itemId: string, date: string) => {
  if (getWeekDayStatus(itemId, date) !== 'free') return;
  if (isDateInPast(date)) return;
  const key = getWeekSelectionKey(itemId, date);
  const next = new Set(weekSelections.value);
  if (next.has(key)) {
    next.delete(key);
  } else {
    next.add(key);
  }
  weekSelections.value = next;
};

const loadWeekData = async (itemGroupId: string, keepResults = false) => {
  weekDataLoading.value = true;
  weekSelections.value = new Set();
  expandedWeekTiles.value = new Set();
  itemsErrorMessage.value = null;
  if (!keepResults) weekBookingResults.value = [];
  try {
    const dates = selectedWeekDates.value;
    const results = await Promise.all(
      dates.map(date => fetchItems(itemGroupId, date).then(resp => ({ date, items: resp.data })))
    );
    const data: Record<string, JsonApiResource<ItemAttributes>[]> = {};
    for (const { date, items: dayItems } of results) {
      data[date] = dayItems;
    }
    weekData.value = data;

    const bookingsResp = await fetchMyBookings().catch(() => ({ data: [] }));
    const bookedMap = new Map<string, string>();
    for (const booking of bookingsResp.data) {
      const bookingDate = booking.attributes.booking_date;
      if (dates.includes(bookingDate)) {
        bookedMap.set(getWeekSelectionKey(booking.attributes.item_id, bookingDate), booking.id);
      }
    }
    myWeekBookings.value = bookedMap;
    await nextTick();
    updateNameTruncation();
  } catch (err) {
    weekData.value = {};
    myWeekBookings.value = new Map();
    itemsErrorMessage.value = isConnectionError(err) ? CONNECTION_LOST_MESSAGE : t('items.unableToLoadWeekly');
  } finally {
    weekDataLoading.value = false;
  }
};

const loadUsers = async () => {
  usersLoading.value = true;
  try {
    const resp = await fetchColleagues();
    usersList.value = resp.data.map(u => ({
      id: u.id,
      displayName: u.attributes.display_name
    }));
  } catch {
    // Silently fail — colleague dropdown will just be empty
  } finally {
    usersLoading.value = false;
  }
};

const resolveColleagueName = (userId: string): string | undefined => {
  const user = usersList.value.find(u => u.id === userId);
  return user?.displayName;
};

const localizeItemsBookingConflict = (err: ApiError): string => {
  const detail = err.detail ?? '';
  const lower = detail.toLowerCase();
  if (lower.includes('booking limit exceeded')) {
    // Parse count and scope from backend message for i18n interpolation
    const countMatch = detail.match(/maximum of (\d+)/);
    const scopeMatch = detail.match(/for (.+)$/);
    const count = countMatch?.[1] ?? '?';
    const scope = scopeMatch?.[1] ?? '';
    if (scope) {
      return t('items.bookingLimitExceeded', { count, scope });
    }
    return t('items.bookingLimitExceededGlobal', { count });
  }
  if (lower.includes('already have this item booked')) {
    return t('items.alreadyBookedByYouForDate');
  }
  return t('items.itemAlreadyBookedForDate');
};

const localizeItemsBookingError = (err: unknown, fallback: string): string => {
  if (!(err instanceof ApiError)) {
    return fallback;
  }
  if (err.status === 409) {
    return localizeItemsBookingConflict(err);
  }
  if (err.status === 404) {
    return t('items.itemNotFound');
  }
  return fallback;
};

const submitWeekBookings = async () => {
  if (!activeItemGroupId.value || weekSelections.value.size === 0) return;

  errorSnackbarMessage.value = null;
  if (bookingType.value === 'colleague') {
    if (!selectedColleagueId.value) {
      errorSnackbarMessage.value =t('items.selectColleagueError');
      return;
    }
  }
  weekBookingInProgress.value = true;
  weekBookingResults.value = [];

  const entries = Array.from(weekSelections.value).map(key => {
    const sep = key.indexOf('::');
    const itemId = key.substring(0, sep);
    const date = key.substring(sep + 2);
    const itemName = weekItems.value.find(item => item.id === itemId)?.name || t('common.item');
    return { itemId, itemName, date };
  });

  const onBehalf: BookOnBehalfOptions | undefined =
    bookingType.value === 'colleague' && selectedColleagueId.value
      ? { forUserId: selectedColleagueId.value, forUserName: resolveColleagueName(selectedColleagueId.value) }
      : undefined;

  const promises = entries.map(async ({ itemId, itemName, date }) => {
    try {
      await createBooking(itemId, date, onBehalf);
      return { itemName, date, success: true };
    } catch (err) {
      const msg = localizeItemsBookingError(err, t('items.bookingFailed'));
      return { itemName, date, success: false, error: msg };
    }
  });

  const results = await Promise.allSettled(promises);

  weekBookingResults.value = results.map(r => {
    const val = r.status === 'fulfilled'
      ? r.value
      : { itemName: '', date: '', success: false as const, error: t('items.unexpectedError') };
    const dayIdx = selectedWeekDates.value.indexOf(val.date);
    const dayLabel = dayIdx >= 0 ? getFullDayLabel(val.date, dayIdx) : val.date;
    return { itemName: val.itemName, date: val.date, dayLabel, success: val.success, error: val.error };
  });

  weekSelections.value = new Set();

  // Refresh week data (keep results visible)
  if (activeItemGroupId.value) {
    await loadWeekData(activeItemGroupId.value, true);
  }

  weekBookingInProgress.value = false;
};

const queryAreaId = computed(() => {
  const value = route.query.areaId;
  return typeof value === 'string' ? value : undefined;
});
const resolvedAreaId = ref<string | null>(null);
function getCurrentAreaId(): string {
  return resolvedAreaId.value || queryAreaId.value || '';
}
const breadcrumbAreaId = computed(() =>
  resolvedAreaId.value ? resolvedAreaId.value : areaName.value ? undefined : queryAreaId.value
);

const breadcrumbs = computed(() => [
  { text: t('common.home'), to: '/' },
  {
    text: areaName.value || t('common.area'),
    to: breadcrumbAreaId.value ? `/areas/${breadcrumbAreaId.value}/item-groups` : undefined
  },
  { text: itemGroupName.value || t('common.itemGroup') }
]);

const { handleAuthError } = useAuthErrorHandler();

const ensureDate = (value: string) => {
  if (value.trim() !== '') {
    return value;
  }
  const today = formatDate(new Date());
  if (selectedDate.value !== today) {
    selectedDate.value = today;
  }
  return today;
};

const loadItems = async (itemGroupId: string, date: string) => {
  itemsErrorMessage.value = null;
  expandedDayTiles.value = new Set();
  expandedDayTiles.value = new Set();
  try {
    const normalizedDate = ensureDate(date);
    const resp = await runItems(() => fetchItems(itemGroupId, normalizedDate));
    items.value = resp.data;
    await nextTick();
    updateNoteTruncation();
    updateNameTruncation();
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (isConnectionError(err)) {
      itemsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      itemsErrorMessage.value = t('items.notFound');
      return;
    }
    itemsErrorMessage.value = t('items.unableToLoad');
  }
};

const bookItem = async (itemId: string) => {
  showSuccessSnackbar.value = false;
  errorSnackbarMessage.value = null;
  lastBookingDetails.value = null;

  // Validate colleague selection
  if (bookingType.value === 'colleague') {
    if (!selectedColleagueId.value) {
      errorSnackbarMessage.value =t('items.selectColleagueError');
      return;
    }
  }

  bookingItemId.value = itemId;
  const bookingDate = selectedDate.value;

  try {
    const onBehalf: BookOnBehalfOptions | undefined =
      bookingType.value === 'colleague' && selectedColleagueId.value
        ? { forUserId: selectedColleagueId.value, forUserName: resolveColleagueName(selectedColleagueId.value) }
        : undefined;

    const result = await createBooking(itemId, bookingDate, onBehalf);
    lastBookingId.value = result.data.id;

    const itemName = items.value.find(entry => entry.id === itemId)?.attributes.name || t('common.item');
    const details = { itemName, date: bookingDate };
    lastBookingDetails.value = details;
    showSuccessFeedback(
      formatBookingSuccessMessage(details),
      'booking-success',
      {
        actionLabel: lastBookingId.value ? t('items.addNote') : undefined,
        actionHandler: lastBookingId.value ? openPostBookingNoteDialog : undefined
      }
    );

    // Reset memorized day to today after successful booking
    resetDayToToday();
    const resetDate = getDay();
    const dayChanged = selectedDate.value !== resetDate;
    selectedDate.value = resetDate;

    // Reset booking type fields
    if (bookingType.value === 'colleague') {
      selectedColleagueId.value = null;
      bookingType.value = 'self';
    }

    // Reload items to reflect updated availability
    if (activeItemGroupId.value && !dayChanged) {
      await loadItems(activeItemGroupId.value, resetDate);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }

    const itemName = items.value.find(entry => entry.id === itemId)?.attributes.name || t('common.item');
    let detail = t('items.unableToBook');

    if (err instanceof ApiError && err.status === 409) {
      detail = localizeItemsBookingConflict(err);

      // Refresh item list so user sees updated availability
      if (activeItemGroupId.value) {
        await loadItems(activeItemGroupId.value, selectedDate.value);
      }
    } else if (err instanceof ApiError && err.status === 404) {
      detail = t('items.itemNotFound');
    }

    errorSnackbarMessage.value = `${itemName} - ${formatDisplayDate(selectedDate.value)}: ${detail}`;
  } finally {
    bookingItemId.value = null;
  }
};

const adminCancelBooking = async (bookingId: string) => {
  showSuccessSnackbar.value = false;
  errorSnackbarMessage.value = null;
  cancelingBookingId.value = bookingId;

  try {
    await cancelBooking(bookingId);
    showSuccessFeedback(t('items.bookingCancelledSuccessfully'), 'booking-success');

    // Reload items to reflect updated availability
    if (activeItemGroupId.value) {
      await loadItems(activeItemGroupId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (err instanceof ApiError && err.status === 404) {
      errorSnackbarMessage.value =t('items.bookingNotFound');
    } else {
      errorSnackbarMessage.value =t('items.unableToCancelBooking');
    }
  } finally {
    cancelingBookingId.value = null;
  }
};


const openPostBookingNoteDialog = () => {
  noteText.value = '';
  showPostBookingNoteDialog.value = true;
};

const saveNoteAfterBooking = async () => {
  if (!lastBookingId.value) return;
  savingNote.value = true;
  try {
    await updateBookingNote(lastBookingId.value, noteText.value);
    showPostBookingNoteDialog.value = false;
    showSuccessFeedback(formatBookingSuccessMessage(lastBookingDetails.value), 'booking-success');
    lastBookingId.value = null;
    lastBookingDetails.value = null;
    if (activeItemGroupId.value) {
      await loadItems(activeItemGroupId.value, selectedDate.value);
    }
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    errorSnackbarMessage.value =t('items.unableToSaveNote');
  } finally {
    savingNote.value = false;
  }
};

const setNoteRef = (id: string) => (el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    noteElements.set(id, el);
    return;
  }
  if (el && '$el' in el && (el.$el instanceof HTMLElement)) {
    noteElements.set(id, el.$el);
    return;
  }
  noteElements.delete(id);
};

const setMeasuredRef = (elements: Map<string, HTMLElement>, id: string) =>
  (el: Element | ComponentPublicInstance | null) => {
    if (el instanceof HTMLElement) {
      elements.set(id, el);
      return;
    }
    if (el && '$el' in el && (el.$el instanceof HTMLElement)) {
      elements.set(id, el.$el);
      return;
    }
    elements.delete(id);
  };

const setDayNameRef = (id: string) => setMeasuredRef(dayNameElements, id);
const setDayNameMeasureRef = (id: string) => setMeasuredRef(dayNameMeasureElements, id);
const setWeekNameRef = (id: string) => setMeasuredRef(weekNameElements, id);
const setWeekNameMeasureRef = (id: string) => setMeasuredRef(weekNameMeasureElements, id);

const updateNoteTruncation = () => {
  const map: Record<string, boolean> = {};
  for (const entry of items.value) {
    const el = noteElements.get(entry.id);
    if (el) {
      map[entry.id] = el.scrollWidth > el.clientWidth;
    }
  }
  noteTruncatedMap.value = map;
};

const updateNameTruncation = () => {
  const dayMap: Record<string, boolean> = {};
  for (const entry of items.value) {
    const displayEl = dayNameElements.get(entry.id);
    const measureEl = dayNameMeasureElements.get(entry.id);
    if (displayEl && measureEl) {
      dayMap[entry.id] = measureEl.scrollWidth > displayEl.clientWidth;
    }
  }
  dayNameTruncatedMap.value = dayMap;

  const weekMap: Record<string, boolean> = {};
  for (const item of weekItems.value) {
    const displayEl = weekNameElements.get(item.id);
    const measureEl = weekNameMeasureElements.get(item.id);
    if (displayEl && measureEl) {
      weekMap[item.id] = measureEl.scrollWidth > displayEl.clientWidth;
    }
  }
  weekNameTruncatedMap.value = weekMap;
};

onMounted(async () => {
  updateViewport();
  try {
    const resp = await fetchMe();
    authStore.userName = resp.data.attributes.display_name;
    authStore.isAdmin = resp.data.attributes.is_admin;
  } catch (err) {
    if (await handleAuthError(err)) {
      return;
    }
    if (isConnectionError(err)) {
      itemsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    throw err;
  }

  // Fetch booking settings (non-blocking, uses default on failure)
  try {
    const settingsResp = await fetchSettings();
    weeksInAdvanced.value = settingsResp.data.attributes.weeks_in_advanced;
  } catch {
    // Non-critical: week selector uses default
  }

  // Load users list for colleague dropdown (non-blocking)
  loadUsers();

  const itemGroupId = route.params.itemGroupId;
  if (typeof itemGroupId !== 'string' || itemGroupId.trim() === '') {
    itemsErrorMessage.value = t('items.notFound');
    return;
  }

  activeItemGroupId.value = itemGroupId;

  // Fetch area and item group names for breadcrumbs
  try {
    const areasResp = await fetchAreas();
    for (const area of areasResp.data) {
      const igResp = await fetchItemGroups(area.id);
      const ig = igResp.data.find(ig => ig.id === itemGroupId);
      if (ig) {
        areaName.value = area.attributes.name;
        itemGroupName.value = ig.attributes.name;
        itemGroupFloorPlan.value = ig.attributes.floor_plan || null;
        inheritedIcon.value = ig.attributes.icon || area.attributes.icon || null;
        resolvedAreaId.value = area.id;
        break;
      }
    }
  } catch (err) {
    if (isConnectionError(err)) {
      itemsErrorMessage.value = CONNECTION_LOST_MESSAGE;
      return;
    }
    // Ignore other errors - breadcrumbs will just show generic names
  }

  if (bookingMode.value === 'week') {
    await loadWeekData(itemGroupId);
  } else {
    await loadItems(itemGroupId, selectedDate.value);
  }
});

watch(
  selectedDate,
  async (value) => {
    setDay(value);
    if (!activeItemGroupId.value || bookingMode.value !== 'day') {
      return;
    }
    await loadItems(activeItemGroupId.value, value);
  },
  { flush: 'post' }
);

watch(bookingMode, async (mode) => {
  if (storage) {
    storage.setItem('sithub_booking_mode', mode);
  }
  if (!activeItemGroupId.value) return;
  if (mode === 'week') {
    await loadWeekData(activeItemGroupId.value);
  } else {
    weekData.value = {};
    weekSelections.value = new Set();
    weekBookingResults.value = [];
    await loadItems(activeItemGroupId.value, selectedDate.value);
  }
});

watch([selectedWeek, showWeekends], async ([week]) => {
  setWeek(week);
  if (!activeItemGroupId.value || bookingMode.value !== 'week') return;
  await loadWeekData(activeItemGroupId.value);
});

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

function updateViewport() {
  if (typeof window.matchMedia === 'function') {
    isMobile.value = window.matchMedia('(max-width: 600px)').matches;
    isCompactFloorPlanViewport.value =
      window.matchMedia('(max-width: 768px)').matches
      || window.matchMedia('(max-height: 500px)').matches;
    return;
  }
  isMobile.value = false;
  isCompactFloorPlanViewport.value = false;
}

function handleResize() {
  updateViewport();
  updateNoteTruncation();
  updateNameTruncation();
}

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function formatDisplayDate(dateStr: string) {
  if (!dateStr) return '';
  const date = new Date(`${dateStr}T00:00:00`);
  if (Number.isNaN(date.getTime())) return dateStr;
  return new Intl.DateTimeFormat(locale.value || undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  }).format(date);
}

function formatBookingSuccessMessage(details: { itemName: string; date: string } | null) {
  return details ? `${details.itemName} - ${formatDisplayDate(details.date)}` : t('items.bookingConfirmed');
}
</script>

<style scoped>
.week-book-footer {
  position: sticky;
  bottom: 0;
  z-index: 2;
  background: rgb(var(--v-theme-surface));
  padding: 12px 0;
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.item-name-shell {
  position: relative;
  display: block;
  width: 100%;
  min-width: 0;
}

.item-name {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: clip;
}

.item-name-measure {
  position: absolute;
  top: 0;
  left: 0;
  visibility: hidden;
  pointer-events: none;
  white-space: nowrap;
}

.week-days {
  display: flex;
  gap: 8px;
  justify-content: space-between;
}

.week-days-compact {
  display: grid;
  grid-template-columns: repeat(5, 1fr); /* overridden by inline style when weekends enabled */
  gap: 4px;
}

.week-day-slot {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 44px;
  padding: 4px;
}

.week-day-label {
  margin-bottom: 2px;
}

.week-day-checkbox {
  min-height: 44px;
}

.week-day-status {
  font-size: 0.7rem;
  line-height: 1.2;
  text-align: center;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.week-day-status-truncated {
  display: inline-block;
}

.week-cancel-icon {
  cursor: pointer;
  vertical-align: middle;
}

.week-day-past {
  opacity: 0.5;
}

.week-day-expanded {
  display: grid;
  grid-template-columns: 40px 180px 1fr;
  align-items: center;
  padding: 4px 0;
}

.week-day-expanded-label {
  white-space: nowrap;
}

.item-filter-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
}

.item-filter-wrapper .item-card {
  flex: 1;
}

.item-expanded {
  grid-column: 1 / -1;
}

.item-filtered-out {
  filter: blur(3px);
  opacity: 0.5;
  pointer-events: none;
}

.item-filtered-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  pointer-events: none;
}

.floor-plan-dialog-card {
  height: 100%;
}

.floor-plan-dialog-body {
  height: 100%;
}
</style>

<style>
.warning-tooltip.v-overlay__content {
  background-color: #fff3e0 !important;
  color: #e65100 !important;
  font-weight: 500;
}
</style>

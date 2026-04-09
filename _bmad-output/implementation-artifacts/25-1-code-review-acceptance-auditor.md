# Acceptance Auditor Review Prompt

You are an Acceptance Auditor. Review this diff against the spec and context docs. Check for:
- violations of acceptance criteria
- deviations from spec intent
- missing implementation of specified behavior
- contradictions between spec constraints and actual code

Output findings as a Markdown list. Each finding:
- one-line title
- which AC/constraint it violates
- evidence from the diff

Additional constraints:
- You receive the diff, the spec, and any context docs.
- No additional context docs were loaded for this story.
- If you find no issues, say `No findings`.

Spec to review against:

```md
# Story 25.1: Editor Layout — Sidebar to Toolbar Dropdowns

Status: review

## Story

As an admin,
I want the floor plan editor to use the full page width with controls in the toolbar,
so that I have maximum canvas space for positioning items on the floor plan.

## Acceptance Criteria

1. **Given** I open the floor plan editor as an admin
   **When** the editor loads
   **Then** there is no left-hand sidebar listing subareas and items; the canvas card uses
   the full available width

2. **Given** the editor is loaded
   **When** I look at the toolbar row
   **Then** I see a subarea dropdown that lists all subareas for the selected floor plan

3. **Given** the editor is loaded
   **When** I select a subarea from the toolbar dropdown
   **Then** the editor switches to that subarea, identical to the old sidebar click behavior

4. **Given** the editor is loaded
   **When** I look at the toolbar row
   **Then** I see an items dropdown that lists all items for the current subarea, each
   indicating whether it is positioned or unpositioned (e.g., via icon or chip)

5. **Given** I select an unpositioned item from the items dropdown
   **When** the selection is made
   **Then** the editor enters draw mode for that item, identical to the old sidebar behavior

6. **Given** I select a positioned item from the items dropdown
   **When** the selection is made
   **Then** the editor selects that item's rectangle on the canvas, identical to the old
   sidebar behavior

7. **Given** I have a positioned item selected via the items dropdown
   **When** I look for a way to delete it
   **Then** I see a delete action accessible from the toolbar that removes
   the item's position from the floor plan
```

Diff to review:

```diff
diff --git a/web/src/views/FloorPlanEditorView.vue b/web/src/views/FloorPlanEditorView.vue
index 998d25f..850a884 100644
--- a/web/src/views/FloorPlanEditorView.vue
+++ b/web/src/views/FloorPlanEditorView.vue
@@ -16,65 +16,7 @@
     </v-alert>
 
     <v-row>
-      <v-col cols="12" md="3" order="2" order-md="1">
-        <v-card class="mb-4" data-cy="editor-sidebar">
-          <v-card-title>
-            {{ isAreaLevel && activeTab === "areas" ? $t('floorPlanEditor.subAreas') : $t('floorPlanEditor.itemsLabel') }}
-          </v-card-title>
-          <v-card-text>
-            <v-alert
-              v-if="!selectedFloorPlan"
-              type="info"
-              variant="tonal"
-              density="compact"
-            >
-              {{ $t('floorPlanEditor.selectFloorPlanToPosition') }}
-            </v-alert>
-
-            <v-list v-else density="compact" nav>
-              <v-list-item
-                v-for="item in scopedItems"
-                :key="item.id"
-                :active="
-                  drawModeItemId === item.id || selectedRectId === item.id
-                "
-                :class="{ 'editor-item--positioned': item.positioned }"
-                :data-cy="`sidebar-item-${item.id}`"
-                @click="selectSidebarItem(item)"
-              >
-                <template #prepend>
-                  <v-icon size="small">
-                    {{ item.positioned ? "$success" : "$location" }}
-                  </v-icon>
-                </template>
-                <v-list-item-title>{{ item.name }}</v-list-item-title>
-                <v-list-item-subtitle>
-                  {{
-                    item.positioned
-                      ? $t('floorPlanEditor.positioned')
-                      : drawModeItemId === item.id
-                        ? $t('floorPlanEditor.drawOnPlan')
-                        : $t('floorPlanEditor.unpositioned')
-                  }}
-                </v-list-item-subtitle>
-                <template #append>
-                  <v-btn
-                    v-if="item.positioned"
-                    icon="$delete"
-                    size="x-small"
-                    variant="text"
-                    color="error"
-                    :data-cy="`sidebar-delete-${item.id}`"
-                    @click.stop="deleteByItemId(item.id)"
-                  />
-                </template>
-              </v-list-item>
-            </v-list>
-          </v-card-text>
-        </v-card>
-      </v-col>
-
-      <v-col cols="12" md="9" order="1" order-md="2">
+      <v-col cols="12">
         <v-card class="mb-4">
           <v-card-text>
             <div class="d-flex flex-wrap align-center ga-3">
@@ -114,10 +56,35 @@
                 :label="$t('floorPlanEditor.subArea')"
                 density="compact"
                 hide-details
-                data-cy="subarea-selector"
+                data-cy="toolbar-subarea-select"
                 style="min-width: 180px; max-width: 240px"
               />
 
+              <v-select
+                v-if="selectedFloorPlan"
+                :model-value="drawModeItemId ?? selectedRectId"
+                :items="scopedItems"
+                item-title="name"
+                item-value="id"
+                :label="$t('floorPlanEditor.itemsLabel')"
+                density="compact"
+                hide-details
+                clearable
+                data-cy="toolbar-items-select"
+                style="min-width: 200px; max-width: 280px"
+                @update:model-value="onToolbarItemSelect"
+              >
+                <template #item="{ item: option, props: listProps }">
+                  <v-list-item v-bind="listProps">
+                    <template #prepend>
+                      <v-icon size="small" :color="option.raw.positioned ? 'success' : undefined">
+                        {{ option.raw.positioned ? 'mdi-check-circle' : 'mdi-map-marker' }}
+                      </v-icon>
+                    </template>
+                  </v-list-item>
+                </template>
+              </v-select>
+
               <v-select
                 v-model="borderWidth"
                 :items="[1, 2, 3, 4, 5]"
@@ -536,6 +503,18 @@ function toPercent(event: PointerEvent) {
   };
 }
 
+function onToolbarItemSelect(itemId: string | null) {
+  if (!itemId) {
+    drawModeItemId.value = null;
+    selectedRectId.value = null;
+    return;
+  }
+  const item = allEditableItems.value.find((entry) => entry.id === itemId);
+  if (item) {
+    selectSidebarItem(item);
+  }
+}
+
 function selectSidebarItem(item: EditableItem) {
   if (item.positioned) {
     selectedRectId.value = item.id;
@@ -1137,10 +1116,6 @@ onUnmounted(() => {
   cursor: nwse-resize;
 }
 
-.editor-item--positioned {
-  opacity: 0.78;
-}
-
 .editor-zoom-controls {
   min-width: 220px;
 }
```

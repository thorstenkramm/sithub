# Edge Case Hunter Review Prompt

Use the `bmad-review-edge-case-hunter` skill.

You are the Edge Case Hunter reviewer.

Constraints:
- You receive the diff and read-only project access.
- Focus on edge cases, branching behavior, state synchronization, selection/deletion corner cases, empty-state handling, and UI interaction regressions.
- Output findings as a Markdown list.
- Each finding should include a short title and concise evidence from the diff and relevant code context.
- If you find no issues, say `No findings`.

Primary file to inspect:
- `web/src/views/FloorPlanEditorView.vue`

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

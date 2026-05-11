import { ref } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY_ITEMS = 'sithub_favorite_items';

// Legacy key from the previous favorites implementation, where item groups
// (rooms/areas) could also be favorited. Story 31.2 removes that feature; we
// purge the key once on first load so stale data does not linger.
const LEGACY_STORAGE_KEY_ITEM_GROUPS = 'sithub_favorite_item_groups';

/** Synthetic area id for the virtual "Favorites" tile. */
export const FAVORITES_AREA_ID = '__favorites__';

export interface ItemFavorite {
  areaId: string;
  itemId: string;
  itemName: string;
  itemGroupId: string;
  itemGroupName: string;
}

let legacyPurged = false;
const favoriteItems = ref<ItemFavorite[]>([]);
let favoritesLoaded = false;

function purgeLegacyItemGroupFavorites() {
  if (legacyPurged) return;
  legacyPurged = true;
  const storage = getSafeLocalStorage();
  if (!storage) return;
  try {
    storage.removeItem(LEGACY_STORAGE_KEY_ITEM_GROUPS);
  } catch {
    // Storage may be unavailable; ignore.
  }
}

/** Reset the one-shot legacy-purge flag. Test-only. */
export function __resetLegacyPurgeForTests() {
  legacyPurged = false;
  favoritesLoaded = false;
  favoriteItems.value = [];
}

function loadItemFavorites(): ItemFavorite[] {
  const storage = getSafeLocalStorage();
  if (!storage) return [];
  try {
    const raw = storage.getItem(STORAGE_KEY_ITEMS);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed)
      ? parsed
        .filter((value): value is ItemFavorite & { areaId?: string } => {
          return typeof value === 'object'
            && value !== null
            && typeof value.itemId === 'string'
            && typeof value.itemName === 'string'
            && typeof value.itemGroupId === 'string'
            && typeof value.itemGroupName === 'string'
            && (typeof value.areaId === 'string' || value.areaId === undefined);
        })
        .map(value => ({
          areaId: value.areaId ?? '',
          itemId: value.itemId,
          itemName: value.itemName,
          itemGroupId: value.itemGroupId,
          itemGroupName: value.itemGroupName
        }))
      : [];
  } catch {
    return [];
  }
}

function persistItemFavorites(items: ItemFavorite[]) {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  storage.setItem(STORAGE_KEY_ITEMS, JSON.stringify(items));
}

export function useFavorites() {
  purgeLegacyItemGroupFavorites();

  if (!favoritesLoaded) {
    favoriteItems.value = loadItemFavorites();
    favoritesLoaded = true;
  }

  const sameItem = (a: ItemFavorite, b: { areaId: string; itemGroupId: string; itemId: string }) =>
    a.itemId === b.itemId
    && (
      (a.areaId === b.areaId && a.itemGroupId === b.itemGroupId)
      || (a.areaId === '' && a.itemGroupId === b.itemGroupId)
    );

  const isItemFavorite = (areaId: string, itemGroupId: string, itemId: string) =>
    favoriteItems.value.some(f => sameItem(f, { areaId, itemGroupId, itemId }));

  function toggleItemFavorite(fav: ItemFavorite): { added: boolean } {
    const existing = favoriteItems.value.findIndex(f =>
      sameItem(f, { areaId: fav.areaId, itemGroupId: fav.itemGroupId, itemId: fav.itemId })
    );
    const added = existing === -1;
    if (added) {
      favoriteItems.value = [...favoriteItems.value, fav];
    } else {
      favoriteItems.value = favoriteItems.value.filter((_, index) => index !== existing);
    }
    persistItemFavorites(favoriteItems.value);
    return { added };
  }

  return {
    favoriteItems,
    isItemFavorite,
    toggleItemFavorite
  };
}

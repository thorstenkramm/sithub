import { ref, computed } from 'vue';
import { getSafeLocalStorage } from './storage';

const STORAGE_KEY_IG = 'sithub_favorite_item_groups';
const STORAGE_KEY_ITEMS = 'sithub_favorite_items';

export interface ItemFavorite {
  areaId: string;
  itemId: string;
  itemName: string;
  itemGroupId: string;
  itemGroupName: string;
}

function loadSet(key: string): Set<string> {
  const storage = getSafeLocalStorage();
  if (!storage) return new Set();
  try {
    const raw = storage.getItem(key);
    if (!raw) return new Set();
    const parsed = JSON.parse(raw);
    return new Set(Array.isArray(parsed) ? parsed.filter((v): v is string => typeof v === 'string') : []);
  } catch {
    return new Set();
  }
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

function persistSet(key: string, set: Set<string>) {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  storage.setItem(key, JSON.stringify([...set]));
}

function persistItemFavorites(items: ItemFavorite[]) {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  storage.setItem(STORAGE_KEY_ITEMS, JSON.stringify(items));
}

export function useFavorites() {
  const favoriteItemGroups = ref<Set<string>>(loadSet(STORAGE_KEY_IG));
  const favoriteItems = ref<ItemFavorite[]>(loadItemFavorites());

  const itemGroupKey = (areaId: string, igId: string) => `${areaId}::${igId}`;
  const itemKey = (areaId: string, itemGroupId: string, itemId: string) =>
    `${areaId}::${itemGroupId}::${itemId}`;

  const isItemGroupFavorite = (areaId: string, igId: string) =>
    favoriteItemGroups.value.has(itemGroupKey(areaId, igId))
    || favoriteItemGroups.value.has(igId);

  function toggleItemGroupFavorite(areaId: string, igId: string): { added: boolean } {
    const next = new Set(favoriteItemGroups.value);
    const scopedKey = itemGroupKey(areaId, igId);
    const added = !next.has(scopedKey) && !next.has(igId);
    if (added) {
      next.add(scopedKey);
    } else {
      next.delete(scopedKey);
      next.delete(igId);
    }
    favoriteItemGroups.value = next;
    persistSet(STORAGE_KEY_IG, next);
    return { added };
  }

  const isItemFavorite = (areaId: string, itemGroupId: string, itemId: string) =>
    favoriteItems.value.some(f =>
      f.itemId === itemId
      && (
        itemKey(f.areaId, f.itemGroupId, f.itemId) === itemKey(areaId, itemGroupId, itemId)
        || (f.areaId === '' && f.itemGroupId === itemGroupId)
      )
    );

  function toggleItemFavorite(fav: ItemFavorite): { added: boolean } {
    const existing = favoriteItems.value.findIndex(f =>
      f.itemId === fav.itemId
      && (
        itemKey(f.areaId, f.itemGroupId, f.itemId) === itemKey(fav.areaId, fav.itemGroupId, fav.itemId)
        || (f.areaId === '' && f.itemGroupId === fav.itemGroupId)
      )
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

  const favoriteItemsForArea = computed(() => {
    return favoriteItems.value;
  });

  return {
    favoriteItemGroups,
    favoriteItems,
    isItemGroupFavorite,
    toggleItemGroupFavorite,
    isItemFavorite,
    toggleItemFavorite,
    favoriteItemsForArea
  };
}

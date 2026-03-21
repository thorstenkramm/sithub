import { useFavorites } from './useFavorites';

describe('useFavorites', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('starts with no favorites', () => {
    const { favoriteItemGroups, favoriteItems } = useFavorites();
    expect(favoriteItemGroups.value.size).toBe(0);
    expect(favoriteItems.value).toEqual([]);
  });

  it('toggles item group favorite on and off', () => {
    const { toggleItemGroupFavorite, isItemGroupFavorite } = useFavorites();

    const result1 = toggleItemGroupFavorite('area-1', 'ig-1');
    expect(result1.added).toBe(true);
    expect(isItemGroupFavorite('area-1', 'ig-1')).toBe(true);
    expect(isItemGroupFavorite('area-2', 'ig-1')).toBe(false);

    const result2 = toggleItemGroupFavorite('area-1', 'ig-1');
    expect(result2.added).toBe(false);
    expect(isItemGroupFavorite('area-1', 'ig-1')).toBe(false);
  });

  it('persists item group favorites to local storage', () => {
    const { toggleItemGroupFavorite } = useFavorites();
    toggleItemGroupFavorite('area-1', 'ig-1');

    const stored = JSON.parse(localStorage.getItem('sithub_favorite_item_groups')!);
    expect(stored).toEqual(['area-1::ig-1']);
  });

  it('toggles item favorite on and off', () => {
    const fav = {
      areaId: 'area-1',
      itemId: 'item-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 1'
    };
    const { toggleItemFavorite, isItemFavorite } = useFavorites();

    const result1 = toggleItemFavorite(fav);
    expect(result1.added).toBe(true);
    expect(isItemFavorite('area-1', 'ig-1', 'item-1')).toBe(true);
    expect(isItemFavorite('area-2', 'ig-1', 'item-1')).toBe(false);

    const result2 = toggleItemFavorite(fav);
    expect(result2.added).toBe(false);
    expect(isItemFavorite('area-1', 'ig-1', 'item-1')).toBe(false);
  });

  it('persists item favorites to local storage', () => {
    const fav = {
      areaId: 'area-1',
      itemId: 'item-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 1'
    };
    const { toggleItemFavorite } = useFavorites();
    toggleItemFavorite(fav);

    const stored = JSON.parse(localStorage.getItem('sithub_favorite_items')!);
    expect(stored).toEqual([fav]);
  });

  it('loads favorites from local storage on init', () => {
    localStorage.setItem('sithub_favorite_item_groups', JSON.stringify(['area-1::ig-1', 'area-2::ig-2']));
    localStorage.setItem('sithub_favorite_items', JSON.stringify([
      {
        areaId: 'area-1',
        itemId: 'item-1',
        itemName: 'Desk 1',
        itemGroupId: 'ig-1',
        itemGroupName: 'Room 1'
      }
    ]));

    const { isItemGroupFavorite, isItemFavorite } = useFavorites();
    expect(isItemGroupFavorite('area-1', 'ig-1')).toBe(true);
    expect(isItemGroupFavorite('area-2', 'ig-2')).toBe(true);
    expect(isItemFavorite('area-1', 'ig-1', 'item-1')).toBe(true);
  });

  it('does not collide favorites across areas or item groups with the same child ids', () => {
    const { toggleItemGroupFavorite, toggleItemFavorite, isItemGroupFavorite, isItemFavorite } = useFavorites();

    toggleItemGroupFavorite('area-1', 'shared-group');
    toggleItemFavorite({
      areaId: 'area-1',
      itemId: 'shared-item',
      itemName: 'Desk 1',
      itemGroupId: 'shared-group',
      itemGroupName: 'Room 1'
    });

    expect(isItemGroupFavorite('area-1', 'shared-group')).toBe(true);
    expect(isItemGroupFavorite('area-2', 'shared-group')).toBe(false);
    expect(isItemFavorite('area-1', 'shared-group', 'shared-item')).toBe(true);
    expect(isItemFavorite('area-2', 'shared-group', 'shared-item')).toBe(false);
  });

  it('handles corrupted local storage gracefully', () => {
    localStorage.setItem('sithub_favorite_item_groups', 'broken');
    localStorage.setItem('sithub_favorite_items', '{bad}');

    const { favoriteItemGroups, favoriteItems } = useFavorites();
    expect(favoriteItemGroups.value.size).toBe(0);
    expect(favoriteItems.value).toEqual([]);
  });
});

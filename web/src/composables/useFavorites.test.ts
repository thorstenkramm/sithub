import { __resetLegacyPurgeForTests, useFavorites } from './useFavorites';

describe('useFavorites', () => {
  beforeEach(() => {
    localStorage.clear();
    __resetLegacyPurgeForTests();
  });

  it('starts with no favorites', () => {
    const { favoriteItems } = useFavorites();
    expect(favoriteItems.value).toEqual([]);
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

  it('shares state across composable callers in mounted views', () => {
    const fav = {
      areaId: 'area-1',
      itemId: 'item-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 1'
    };
    const first = useFavorites();
    const second = useFavorites();

    first.toggleItemFavorite(fav);
    expect(second.isItemFavorite('area-1', 'ig-1', 'item-1')).toBe(true);

    second.toggleItemFavorite(fav);
    expect(first.favoriteItems.value).toEqual([]);
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
    localStorage.setItem('sithub_favorite_items', JSON.stringify([
      {
        areaId: 'area-1',
        itemId: 'item-1',
        itemName: 'Desk 1',
        itemGroupId: 'ig-1',
        itemGroupName: 'Room 1'
      }
    ]));

    const { isItemFavorite } = useFavorites();
    expect(isItemFavorite('area-1', 'ig-1', 'item-1')).toBe(true);
  });

  it('does not collide favorites across areas with the same child ids', () => {
    const { toggleItemFavorite, isItemFavorite } = useFavorites();

    toggleItemFavorite({
      areaId: 'area-1',
      itemId: 'shared-item',
      itemName: 'Desk 1',
      itemGroupId: 'shared-group',
      itemGroupName: 'Room 1'
    });

    expect(isItemFavorite('area-1', 'shared-group', 'shared-item')).toBe(true);
    expect(isItemFavorite('area-2', 'shared-group', 'shared-item')).toBe(false);
  });

  it('handles corrupted local storage gracefully', () => {
    localStorage.setItem('sithub_favorite_items', '{bad}');

    const { favoriteItems } = useFavorites();
    expect(favoriteItems.value).toEqual([]);
  });

  it('purges legacy item-group favorites on first load', () => {
    localStorage.setItem('sithub_favorite_item_groups', JSON.stringify(['area-1::ig-1']));

    useFavorites();

    expect(localStorage.getItem('sithub_favorite_item_groups')).toBeNull();
  });
});

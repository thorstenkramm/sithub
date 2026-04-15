import { fetchWeeklyMatrix } from './itemGroupMatrix';

describe('fetchWeeklyMatrix', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetches matrix for area and week', async () => {
    const mockResponse = {
      data: [
        {
          id: 'ig-1',
          type: 'item-group-weekly-matrix',
          attributes: {
            item_group_id: 'ig-1',
            item_group_name: 'Room 101',
            days: [{ date: '2026-04-13', weekday: 'MO' }],
            items: [
              {
                item_id: 'desk-1',
                item_name: 'Desk 1',
                equipment: ['Dock'],
                cells: [
                  {
                    date: '2026-04-13',
                    availability: 'free',
                    booked_by_me: false
                  }
                ]
              }
            ]
          }
        }
      ]
    };

    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    } as Response);

    const result = await fetchWeeklyMatrix('area-1', '2026-W16', 5);
    expect(result.data).toHaveLength(1);
    expect(result.data[0].attributes.item_group_id).toBe('ig-1');
    expect(result.data[0].attributes.items[0].item_id).toBe('desk-1');
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/areas/area-1/item-groups/matrix?week=2026-W16&days=5',
      expect.objectContaining({ headers: expect.any(Headers) })
    );
  });

  it('fetches matrix without optional params', async () => {
    const mockResponse = { data: [] };

    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    } as Response);

    await fetchWeeklyMatrix('area-1');
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/areas/area-1/item-groups/matrix',
      expect.objectContaining({ headers: expect.any(Headers) })
    );
  });
});

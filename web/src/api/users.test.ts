import { deleteUser, fetchColleagues, fetchUserDeletionImpact, fetchUsers, updateUser } from './users';
import { apiRequest, throwApiError } from './client';

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
  throwApiError: vi.fn(),
  ApiError: class ApiError extends Error {
    status: number;
    detail: string | null;
    constructor(message: string, status: number, detail: string | null = null) {
      super(message);
      this.status = status;
      this.detail = detail;
    }
  },
  CONNECTION_LOST_MESSAGE: 'Connection to server lost'
}));

const apiRequestMock = apiRequest as unknown as ReturnType<typeof vi.fn>;
const throwApiErrorMock = throwApiError as unknown as ReturnType<typeof vi.fn>;

describe('fetchUsers', () => {
  beforeEach(() => {
    apiRequestMock.mockReset();
    apiRequestMock.mockResolvedValue({ data: [], meta: { total: 0, limit: 20, offset: 0 } });
  });

  it.each([
    { limit: 20, offset: 0, query: 'page%5Blimit%5D=20&page%5Boffset%5D=0' },
    { limit: 20, offset: 40, query: 'page%5Blimit%5D=20&page%5Boffset%5D=40' },
    { limit: 5, offset: 15, query: 'page%5Blimit%5D=5&page%5Boffset%5D=15' }
  ])('requests page limit $limit offset $offset', async ({ limit, offset, query }) => {
    await fetchUsers({ limit, offset });

    expect(apiRequestMock).toHaveBeenCalledWith(`/api/v1/users?${query}`);
  });
});

describe('fetchUserDeletionImpact', () => {
  it('calls the deletion-impact endpoint', async () => {
    apiRequestMock.mockResolvedValue({ data: { attributes: { mode: 'delete' } } });

    await fetchUserDeletionImpact('user-1');

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/users/user-1/deletion-impact');
  });
});

describe('deleteUser', () => {
  beforeEach(() => {
    throwApiErrorMock.mockReset();
  });

  it('returns null for hard delete 204', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true, status: 204 });
    vi.stubGlobal('fetch', fetchMock);

    await expect(deleteUser('user-1')).resolves.toBeNull();

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/users/user-1', {
      method: 'DELETE',
      headers: { Accept: 'application/vnd.api+json' }
    });
  });

  it('returns the user resource for deactivation 200', async () => {
    const body = { data: { id: 'user-1', type: 'users', attributes: { deleted: true } } };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve(body) }));

    await expect(deleteUser('user-1')).resolves.toEqual(body);
  });

  it('throws the JSON:API error on failure', async () => {
    const response = { ok: false, status: 400 };
    throwApiErrorMock.mockRejectedValue(new Error('bad request'));
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(response));

    await expect(deleteUser('user-1')).rejects.toThrow('bad request');
    expect(throwApiErrorMock).toHaveBeenCalledWith(response);
  });
});

describe('fetchColleagues', () => {
  it('calls the colleagues endpoint', async () => {
    apiRequestMock.mockResolvedValue({ data: [] });

    await fetchColleagues();

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/colleagues');
  });
});

describe('updateUser', () => {
  it('sends a JSON:API PATCH with only the given attributes', async () => {
    apiRequestMock.mockResolvedValue({ data: {} });

    await updateUser('u1', { is_data_exporter: true });

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/users/u1', {
      method: 'PATCH',
      body: JSON.stringify({ data: { type: 'users', id: 'u1', attributes: { is_data_exporter: true } } })
    });
  });
});

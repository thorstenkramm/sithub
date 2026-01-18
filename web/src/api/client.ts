const JSON_API = 'application/vnd.api+json';

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export async function apiRequest<T>(input: RequestInfo, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set('Accept', JSON_API);
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', JSON_API);
  }

  const response = await fetch(input, {
    ...init,
    headers
  });

  if (!response.ok) {
    throw new ApiError(`Request failed: ${response.status}`, response.status);
  }

  return response.json() as Promise<T>;
}

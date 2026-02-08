const JSON_API = 'application/vnd.api+json';

export class ApiError extends Error {
  status: number;
  detail: string | null;

  constructor(message: string, status: number, detail: string | null = null) {
    super(message);
    this.status = status;
    this.detail = detail;
  }
}

interface JsonApiErrorResponse {
  errors?: Array<{
    status?: string;
    title?: string;
    detail?: string;
    code?: string;
  }>;
}

export async function parseErrorDetail(response: Response): Promise<string | null> {
  try {
    const body = (await response.json()) as JsonApiErrorResponse;
    const detail = body.errors?.[0]?.detail;
    if (detail) {
      return detail;
    }
  } catch {
    // Ignore JSON parse errors
  }
  return null;
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
    const detail = await parseErrorDetail(response);
    throw new ApiError(`Request failed: ${response.status}`, response.status, detail);
  }

  return response.json() as Promise<T>;
}

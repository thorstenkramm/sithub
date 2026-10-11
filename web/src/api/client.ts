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

/** Reads all `errors[].detail` values of a JSON:API error body, joined by newlines. */
export async function parseErrorDetail(response: Response): Promise<string | null> {
  try {
    const body: unknown = await response.json();
    if (!body || typeof body !== 'object' || !('errors' in body) || !Array.isArray(body.errors)) {
      return null;
    }
    const details = body.errors.flatMap((entry: unknown) => {
      if (entry && typeof entry === 'object' && 'detail' in entry &&
          typeof entry.detail === 'string' && entry.detail.length > 0) {
        return [entry.detail];
      }
      return [];
    });
    if (details.length > 0) {
      return details.join('\n');
    }
  } catch {
    // Ignore JSON parse errors
  }
  return null;
}

/** Throws an ApiError for a failed response, carrying the backend error detail if present. */
export async function throwApiError(response: Response): Promise<never> {
  const detail = await parseErrorDetail(response);
  throw new ApiError(`Request failed: ${response.status}`, response.status, detail);
}

export const CONNECTION_LOST_MESSAGE = 'Connection to server lost';

/** Optional translator for the connection-lost message; registered by the i18n plugin. */
let translateConnectionLost: (() => string) | null = null;

/** Registers the locale-aware connection-lost message. Called once by the i18n plugin. */
export function setConnectionLostTranslator(translate: () => string) {
  translateConnectionLost = translate;
}

function connectionLostMessage(): string {
  return translateConnectionLost ? translateConnectionLost() : CONNECTION_LOST_MESSAGE;
}

/**
 * Returns the reason of a failed request for display: the untranslated backend detail,
 * the connection lost message, or the HTTP status as fallback.
 */
export function apiErrorReason(err: unknown): string {
  if (err instanceof ApiError) {
    if (err.detail) {
      return err.detail;
    }
    return err.status === 0 ? connectionLostMessage() : `HTTP ${err.status}`;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return String(err);
}

/**
 * Formats a user-facing error message: the already translated prefix, a line break
 * and the reason from apiErrorReason.
 */
export function formatApiErrorMessage(prefix: string, err: unknown): string {
  return `${prefix}\n${apiErrorReason(err)}`;
}

export function isConnectionError(err: unknown): boolean {
  return err instanceof ApiError && err.status === 0;
}

export async function apiRequest<T>(input: RequestInfo, init: RequestInit = {}): Promise<T> {
  const response = await sendRequest(input, init, JSON_API);
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

/** Downloads a file while preserving the same structured errors as JSON requests. */
export async function apiDownload(input: RequestInfo): Promise<{ blob: Blob; filename: string | null }> {
  const response = await sendRequest(input, {}, 'text/csv');
  const disposition = response.headers.get('Content-Disposition') ?? '';
  const match = /(?:^|;)\s*filename\s*=\s*(?:"([^"]*)"|([^;\s]+))/i.exec(disposition);
  try {
    return { blob: await response.blob(), filename: match ? (match[1] ?? match[2]) : null };
  } catch {
    throw new ApiError(connectionLostMessage(), 0);
  }
}

async function sendRequest(input: RequestInfo, init: RequestInit, accept: string): Promise<Response> {
  const headers = new Headers(init.headers);
  headers.set('Accept', accept);
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', JSON_API);
  }

  let response: Response;
  try {
    response = await fetch(input, {
      ...init,
      headers
    });
  } catch {
    throw new ApiError(connectionLostMessage(), 0);
  }

  if (!response.ok) {
    await throwApiError(response);
  }

  return response;
}

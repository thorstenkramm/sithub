import { getSafeLocalStorage } from './storage';

const STORAGE_KEY = 'sithub_warning_suppressed';

function hashWarning(warning: string): string {
  let hash = 0;
  for (let i = 0; i < warning.length; i++) {
    hash = ((hash << 5) - hash + warning.charCodeAt(i)) | 0;
  }
  return hash.toString(36);
}

function makeKey(itemId: string, warning: string): string {
  return `${itemId}::${hashWarning(warning)}`;
}

function loadSuppressed(): Set<string> {
  const storage = getSafeLocalStorage();
  if (!storage) return new Set();
  try {
    const raw = storage.getItem(STORAGE_KEY);
    if (!raw) return new Set();
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed)) return new Set(parsed);
    return new Set();
  } catch {
    return new Set();
  }
}

function saveSuppressed(keys: Set<string>): void {
  const storage = getSafeLocalStorage();
  if (!storage) return;
  storage.setItem(STORAGE_KEY, JSON.stringify([...keys]));
}

/**
 * Composable for managing warning suppression per item.
 * Stores suppressed keys (itemId::warningHash) in localStorage.
 * Suppression auto-resets when the warning text changes.
 */
export function useWarningSuppression() {
  function isWarningSuppressed(itemId: string, warning: string): boolean {
    return loadSuppressed().has(makeKey(itemId, warning));
  }

  function suppressWarning(itemId: string, warning: string): void {
    const keys = loadSuppressed();
    keys.add(makeKey(itemId, warning));
    saveSuppressed(keys);
  }

  return { isWarningSuppressed, suppressWarning };
}

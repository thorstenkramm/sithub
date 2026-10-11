/** Saves a downloaded blob and releases the temporary DOM element and object URL afterwards. */
export function saveBlob(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = filename;
  document.body.appendChild(anchor);
  try {
    anchor.click();
  } finally {
    anchor.remove();
    // Allow the browser to begin consuming the URL before releasing it.
    setTimeout(() => URL.revokeObjectURL(url), 0);
  }
}

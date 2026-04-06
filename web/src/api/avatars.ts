/** Returns the URL for a user's avatar image. */
export function getAvatarUrl(userId: string): string {
  return `/api/v1/avatars/${encodeURIComponent(userId)}`;
}

/** Uploads a new avatar for the current user. */
export async function uploadAvatar(file: File): Promise<void> {
  const formData = new FormData();
  formData.append('avatar', file);

  const response = await fetch('/api/v1/me/avatar', {
    method: 'POST',
    body: formData
  });

  if (!response.ok) {
    throw new Error(`Upload failed: ${response.status}`);
  }
}

/** Deletes the current user's avatar. */
export async function deleteAvatar(): Promise<void> {
  const response = await fetch('/api/v1/me/avatar', {
    method: 'DELETE'
  });

  if (!response.ok) {
    throw new Error(`Delete failed: ${response.status}`);
  }
}

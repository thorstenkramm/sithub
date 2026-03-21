const MDI_ICON_NAME_RE = /^mdi-[a-z0-9-]+$/;

export function isConfiguredIconName(icon: string | null | undefined): icon is string {
  return typeof icon === 'string' && MDI_ICON_NAME_RE.test(icon.trim());
}

export function resolveConfiguredIcon(
  icon: string | null | undefined,
  fallback: string
): string {
  return isConfiguredIconName(icon) ? `mdiFont:${icon.trim()}` : fallback;
}

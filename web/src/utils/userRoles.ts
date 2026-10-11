import type { UserListAttributes } from '../api/users';

/**
 * userRoleLabels returns the i18n keys of the roles a user holds, in display order.
 * It is the single place deciding which role chips the admin users table shows.
 */
export function userRoleLabels(attrs: Pick<UserListAttributes, 'is_admin' | 'is_data_exporter'>): string[] {
  const roles: string[] = [];
  if (attrs.is_admin) roles.push('adminUsers.roleAdmin');
  if (attrs.is_data_exporter) roles.push('adminUsers.roleDataExporter');
  return roles.length > 0 ? roles : ['adminUsers.roleUser'];
}

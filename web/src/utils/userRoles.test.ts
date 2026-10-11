import { userRoleLabels } from './userRoles';

describe('userRoleLabels', () => {
  it.each([
    { name: 'regular user', attrs: { is_admin: false, is_data_exporter: false }, expected: ['adminUsers.roleUser'] },
    { name: 'admin', attrs: { is_admin: true, is_data_exporter: false }, expected: ['adminUsers.roleAdmin'] },
    {
      name: 'data exporter',
      attrs: { is_admin: false, is_data_exporter: true },
      expected: ['adminUsers.roleDataExporter']
    },
    {
      name: 'admin and data exporter',
      attrs: { is_admin: true, is_data_exporter: true },
      expected: ['adminUsers.roleAdmin', 'adminUsers.roleDataExporter']
    }
  ])('returns the role label keys for a $name', ({ attrs, expected }) => {
    expect(userRoleLabels(attrs)).toEqual(expected);
  });
});

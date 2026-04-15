/** Builds a query string for week-based area endpoints (availability, matrix). */
export function buildWeekQuery(week?: string, days?: number): string {
  const params = new URLSearchParams();
  if (week) params.set('week', week);
  if (days) params.set('days', String(days));
  return params.toString();
}

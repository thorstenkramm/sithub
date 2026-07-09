import { ref } from 'vue';
import { fetchColleagues } from '../api/users';

export interface ColleagueOption {
  id: string;
  displayName: string;
}

/**
 * useColleagues provides the shared colleague list used by every on-behalf
 * booking surface (tiles, floor plan, weekly table). It loads the colleague
 * list once, exposes a loading flag, and resolves a colleague's display name.
 *
 * The returned state is instance-scoped; each hosting component owns its own
 * copy, mirroring useWarningConfirmation.
 */
export function useColleagues() {
  const colleagueList = ref<ColleagueOption[]>([]);
  const colleaguesLoading = ref(false);

  async function loadColleagues() {
    if (colleagueList.value.length > 0) return;
    colleaguesLoading.value = true;
    try {
      const resp = await fetchColleagues();
      colleagueList.value = resp.data.map((r) => ({
        id: r.id,
        displayName: r.attributes.display_name
      }));
    } catch {
      colleagueList.value = [];
    } finally {
      colleaguesLoading.value = false;
    }
  }

  function resolveColleagueName(id: string): string {
    return colleagueList.value.find((c) => c.id === id)?.displayName ?? '';
  }

  return { colleagueList, colleaguesLoading, loadColleagues, resolveColleagueName };
}

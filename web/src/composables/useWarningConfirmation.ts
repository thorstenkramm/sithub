import { ref } from 'vue';
import { useWarningSuppression } from './useWarningSuppression';

/** An item whose warning may need confirmation before booking. */
export interface WarnItem {
  itemId: string;
  itemName: string;
  warning: string;
}

/**
 * useWarningConfirmation drives the uniform, sequential pre-booking warning
 * confirmation shared by every booking surface (tiles, floor plan, weekly
 * table). Call `present(items, onConfirmed)` before booking:
 *
 * - items without a warning, or whose warning was dismissed via "Don't show
 *   again" (keyed on item id + warning text — see useWarningSuppression), are
 *   skipped;
 * - remaining warnings are shown one after another; `onConfirmed` runs only
 *   after all are confirmed;
 * - cancelling any one aborts the whole flow (onConfirmed is not called).
 *
 * Bind the returned reactive state to WarningConfirmDialog.
 */
export function useWarningConfirmation() {
  const { isWarningSuppressed, suppressWarning } = useWarningSuppression();

  const show = ref(false);
  const itemName = ref('');
  const message = ref('');
  const dontShowAgain = ref(false);
  const queue = ref<WarnItem[]>([]);
  const currentItemId = ref('');
  let onAllConfirmed: (() => void) | null = null;

  function showCurrent() {
    const first = queue.value[0];
    if (!first) return;
    currentItemId.value = first.itemId;
    itemName.value = first.itemName;
    message.value = first.warning;
    dontShowAgain.value = false;
    show.value = true;
  }

  function reset() {
    show.value = false;
    itemName.value = '';
    message.value = '';
    dontShowAgain.value = false;
    currentItemId.value = '';
  }

  function present(items: WarnItem[], onConfirmed: () => void) {
    if (show.value) return;
    const pending = items.filter(
      (i) => i.warning && !isWarningSuppressed(i.itemId, i.warning),
    );
    if (pending.length === 0) {
      onConfirmed();
      return;
    }
    onAllConfirmed = onConfirmed;
    queue.value = pending;
    showCurrent();
  }

  function confirm() {
    if (!show.value) return;
    if (dontShowAgain.value) {
      suppressWarning(currentItemId.value, message.value);
    }
    queue.value = queue.value.slice(1);
    if (queue.value.length > 0) {
      showCurrent();
      return;
    }
    reset();
    const callback = onAllConfirmed;
    onAllConfirmed = null;
    callback?.();
  }

  function cancel() {
    reset();
    queue.value = [];
    onAllConfirmed = null;
  }

  return { show, itemName, message, dontShowAgain, present, confirm, cancel };
}

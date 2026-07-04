import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useWarningConfirmation } from '../useWarningConfirmation';

describe('useWarningConfirmation', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('runs onConfirmed immediately when no item has a warning', () => {
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present([], done);
    expect(wc.show.value).toBe(false);
    expect(done).toHaveBeenCalledTimes(1);
  });

  it('skips a whitespace-only warning and books directly', () => {
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present([{ itemId: 'd1', itemName: 'Desk 1', warning: '   \n ' }], done);
    expect(wc.show.value).toBe(false);
    expect(done).toHaveBeenCalledTimes(1);
  });

  it('shows a single confirmation and runs onConfirmed after confirm', () => {
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present([{ itemId: 'd1', itemName: 'Desk 1', warning: 'Apple only' }], done);
    expect(wc.show.value).toBe(true);
    expect(wc.itemName.value).toBe('Desk 1');
    expect(wc.message.value).toBe('Apple only');
    expect(done).not.toHaveBeenCalled();

    wc.confirm();
    expect(wc.show.value).toBe(false);
    expect(done).toHaveBeenCalledTimes(1);
  });

  it('shows warnings sequentially and only books after all are confirmed (FR164)', () => {
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present(
      [
        { itemId: 'd1', itemName: 'Desk 1', warning: 'W1' },
        { itemId: 'd2', itemName: 'Desk 2', warning: 'W2' },
      ],
      done,
    );
    expect(wc.message.value).toBe('W1');
    wc.confirm();
    expect(done).not.toHaveBeenCalled();
    expect(wc.show.value).toBe(true);
    expect(wc.message.value).toBe('W2');
    wc.confirm();
    expect(done).toHaveBeenCalledTimes(1);
    expect(wc.show.value).toBe(false);
  });

  it('cancelling any confirmation aborts the whole flow (FR164)', () => {
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present(
      [
        { itemId: 'd1', itemName: 'Desk 1', warning: 'W1' },
        { itemId: 'd2', itemName: 'Desk 2', warning: 'W2' },
      ],
      done,
    );
    wc.cancel();
    expect(wc.show.value).toBe(false);
    expect(done).not.toHaveBeenCalled();
  });

  it('skips items whose warning was dismissed via "Don\'t show again"', () => {
    // First flow: dismiss d1's warning.
    const first = useWarningConfirmation();
    first.present([{ itemId: 'd1', itemName: 'Desk 1', warning: 'W1' }], () => {});
    first.dontShowAgain.value = true;
    first.confirm();

    // Second flow with d1 (suppressed) + d2 (not) -> only d2 is shown.
    const wc = useWarningConfirmation();
    const done = vi.fn();
    wc.present(
      [
        { itemId: 'd1', itemName: 'Desk 1', warning: 'W1' },
        { itemId: 'd2', itemName: 'Desk 2', warning: 'W2' },
      ],
      done,
    );
    expect(wc.message.value).toBe('W2');
    wc.confirm();
    expect(done).toHaveBeenCalledTimes(1);
  });
});

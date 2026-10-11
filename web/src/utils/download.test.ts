import { saveBlob } from './download';

it('saves the blob through a temporary download anchor and releases its URL', () => {
  vi.useFakeTimers();
  const create = vi.fn(() => 'blob:test');
  const revoke = vi.fn();
  vi.stubGlobal('URL', { createObjectURL: create, revokeObjectURL: revoke });
  const click = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(function (this: HTMLAnchorElement) {
    expect(this.download).toBe('export.csv');
    expect(this.href).toBe('blob:test');
    expect(this.isConnected).toBe(true);
  });
  const blob = new Blob(['csv']);
  saveBlob(blob, 'export.csv');
  expect(create).toHaveBeenCalledWith(blob);
  expect(click).toHaveBeenCalledOnce();
  expect(document.querySelector('a')).toBeNull();
  vi.runAllTimers();
  expect(revoke).toHaveBeenCalledWith('blob:test');
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
  vi.useRealTimers();
});

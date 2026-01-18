import { flushPromises } from '@vue/test-utils';
import { ApiError } from '../api/client';

export function mockWindowLocation() {
  const originalLocation = window.location;
  Object.defineProperty(window, 'location', {
    configurable: true,
    value: { href: 'http://localhost/' }
  });

  return () => {
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation
    });
  };
}

export async function expectLoginRedirect(mountView: () => void) {
  const restore = mockWindowLocation();
  mountView();
  await flushPromises();
  expect(window.location.href).toBe('/oauth/login');
  restore();
}

export async function expectAccessDeniedRedirect(
  mountView: () => void,
  pushMock: (...args: unknown[]) => unknown
) {
  mountView();
  await flushPromises();
  expect(pushMock).toHaveBeenCalledWith('/access-denied');
}

export function defineAuthRedirectTests(
  fetchMeMock: { mockRejectedValue: (value: unknown) => void },
  mountView: () => void,
  pushMock: (...args: unknown[]) => unknown
) {
  it('redirects to login on 401', async () => {
    fetchMeMock.mockRejectedValue(new ApiError('Unauthorized', 401));
    await expectLoginRedirect(mountView);
  });

  it('redirects to access denied on 403', async () => {
    fetchMeMock.mockRejectedValue(new ApiError('Forbidden', 403));
    await expectAccessDeniedRedirect(mountView, pushMock);
  });
}

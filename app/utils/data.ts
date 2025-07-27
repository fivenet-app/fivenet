import type { AsyncDataRequestStatus } from '#app';

export function isRequestPending(status: AsyncDataRequestStatus): boolean {
    return ['idle', 'pending'].includes(status);
}

import { createSharedComposable, useIntervalFn, useNow } from '@vueuse/core';

// Shares one second-resolution clock between active consumers and disposes it
// when no consumers remain.
export const useSecondClock = createSharedComposable(() => useNow({ scheduler: (callback) => useIntervalFn(callback, 1_000) }));

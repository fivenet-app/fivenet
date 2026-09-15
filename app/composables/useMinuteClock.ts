import { createSharedComposable, useIntervalFn, useNow } from '@vueuse/core';

// Shares one minute-resolution clock between active consumers and disposes it
// when no consumers remain.
export const useMinuteClock = createSharedComposable(() =>
    useNow({ scheduler: (callback) => useIntervalFn(callback, 60_000) }),
);

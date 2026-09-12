import { createSharedComposable, useNow } from '@vueuse/core';

// Shares one clock between active consumers and disposes it when none remain.
export const useSecondClock = createSharedComposable(() => useNow({ interval: 1_000 }));

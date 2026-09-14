import { createSharedComposable, useNow } from '@vueuse/core';

// Shares one minute clock between active consumers and disposes it when none remain.
export const useMinuteClock = createSharedComposable(() => useNow({ interval: 60_000 }));

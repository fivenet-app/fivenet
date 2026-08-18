import { createSharedComposable } from '@vueuse/core';

export const useCalendarEntryShortcutState = createSharedComposable(() => {
    const isModalOpen = ref(false);
    const isPopoverOpen = ref(false);

    return {
        isModalOpen,
        isPopoverOpen,
    };
});

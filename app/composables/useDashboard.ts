import { LazyHelpSlideover, LazyNotificationsNotificationSlideover } from '#components';
import { createSharedComposable } from '@vueuse/core';

const _useDashboard = () => {
    const route = useRoute();

    const isDashboardSidebarSlideoverOpen = ref(false);
    const isHelpSlideoverOpen = ref(false);
    const isNotificationSlideoverOpen = ref(false);

    const overlay = useOverlay();

    const notificationsSlideover = overlay.create(LazyNotificationsNotificationSlideover);
    const helpSlideover = overlay.create(LazyHelpSlideover);

    defineShortcuts({
        b: () => (isNotificationSlideoverOpen.value = true),
    });

    watch(isNotificationSlideoverOpen, () => {
        if (isNotificationSlideoverOpen.value) {
            notificationsSlideover.open();
        } else {
            notificationsSlideover.close();
        }
    });

    watch(isHelpSlideoverOpen, () => {
        if (isHelpSlideoverOpen.value) {
            helpSlideover.open();
        } else {
            helpSlideover.close();
        }
    });

    watch(
        () => route.fullPath,
        () => {
            isHelpSlideoverOpen.value = false;
            isNotificationSlideoverOpen.value = false;
        },
    );

    return {
        isDashboardSidebarSlideoverOpen,
        isHelpSlideoverOpen,
        isNotificationSlideoverOpen,
    };
};

export const useDashboard = createSharedComposable(_useDashboard);

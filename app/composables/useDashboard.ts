import { LazyHelpSlideover, LazyNotificationsNotificationSlideover } from '#components';

const _useDashboard = () => {
    const route = useRoute();

    const isDashboardSidebarSlideoverOpen = ref(false);
    const isHelpSlideoverOpen = ref(false);
    const isNotificationSlideoverOpen = ref(false);
    const isCommandSearchOpen = ref(false);

    const overlay = useOverlay();

    const notificationsSlideover = overlay.create(LazyNotificationsNotificationSlideover);
    const helpSlideover = overlay.create(LazyHelpSlideover);

    defineShortcuts({
        b: () => (isNotificationSlideoverOpen.value = true),
    });

    watch(isNotificationSlideoverOpen, (value) => {
        if (value) {
            notificationsSlideover.open();
        } else {
            notificationsSlideover.close();
        }
    });

    watch(isHelpSlideoverOpen, (value) => {
        if (value) {
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
            isCommandSearchOpen.value = false;
        },
    );

    defineShortcuts({
        '?': () => (isHelpSlideoverOpen.value = true),
    });

    return {
        isDashboardSidebarSlideoverOpen,
        isHelpSlideoverOpen,
        isNotificationSlideoverOpen,
        isCommandSearchOpen,
    };
};

export const useDashboard = createSharedComposable(_useDashboard);

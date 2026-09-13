import { LazyHelpSlideover } from '#components';

const _useDashboard = () => {
    const route = useRoute();

    const isDashboardSidebarSlideoverOpen = ref<boolean>(false);
    const isHelpSlideoverOpen = ref<boolean>(false);
    const isNotificationSlideoverOpen = ref<boolean>(false);
    const isCommandSearchOpen = ref<boolean>(false);

    const overlay = useOverlay();

    const helpSlideover = overlay.create(LazyHelpSlideover);

    defineShortcuts({
        b: () => (isNotificationSlideoverOpen.value = true),
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

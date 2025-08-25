import { LazyHelpSlideover, LazyNotificationsNotificationSlideover } from '#components';
import { createSharedComposable } from '@vueuse/core';

const _useDashboard = () => {
    const route = useRoute();
    const router = useRouter();

    const isDashboardSidebarSlideoverOpen = ref(false);
    const isHelpSlideoverOpen = ref(false);
    const isNotificationSlideoverOpen = ref(false);
    const isDashboardSearchModalOpen = ref(false);

    const overlay = useOverlay();

    const notificationsSlideover = overlay.create(LazyNotificationsNotificationSlideover);
    const helpSlideover = overlay.create(LazyHelpSlideover);

    defineShortcuts({
        'g-h': () => router.push('/'),
        'g-e': () => router.push('/mail'),
        'g-c': () => router.push('/citizens'),
        'g-v': () => router.push('/vehicles'),
        'g-d': () => router.push('/documents'),
        'g-j': () => router.push('/jobs'),
        'g-k': () => router.push('/calendar'),
        'g-q': () => router.push('/qualifications'),
        'g-m': () => router.push('/livemap'),
        'g-w': () => router.push('/centrum'),
        'g-l': () => router.push('/wiki'),
        'g-p': () => router.push('/settings'),
        '?': () => (isHelpSlideoverOpen.value = true),
        b: () => (isNotificationSlideoverOpen.value = true),
    });

    watch(isNotificationSlideoverOpen, () => isNotificationSlideoverOpen && notificationsSlideover.open());

    watch(isHelpSlideoverOpen, () => isHelpSlideoverOpen && helpSlideover.open());

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
        isDashboardSearchModalOpen,
    };
};

export const useDashboard = createSharedComposable(_useDashboard);

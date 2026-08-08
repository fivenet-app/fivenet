import { focusNUITargets, onFocusHandler, onNUIMessage, toggleTablet } from '~/composables/nui';

export default defineNuxtPlugin({
    name: 'nui',
    parallel: true,

    async setup() {
        const query = useRouter().currentRoute.value.query;
        const settingsStore = useSettingsStore();

        if (query?.nui !== undefined) {
            const nuiQuery = query.nui as string;

            const logger = useLogger('🎮 NUI');
            if (nuiQuery.toLowerCase() !== 'false') {
                settingsStore.setNuiSettings(true, nuiQuery);
                logger.info('Enabled NUI integration, resource:', settingsStore.nuiResourceName);
            } else {
                settingsStore.setNuiSettings(false);
                logger.info('Disabled NUI integration');
            }
        }

        if (query?.refreshApp !== undefined) {
            return reloadNuxtApp({
                persistState: false,
                ttl: 7_500, // 7.5 seconds
            });
        }

        const overlay = useOverlay();
        const { isDashboardSidebarSlideoverOpen } = useDashboard();
        const { nuiEnabled } = storeToRefs(settingsStore);

        if (nuiEnabled.value) {
            useEventListener(window, 'message', onNUIMessage);

            // Close tablet on escape presses outside of inputs and overlays.
            const stopEscapeKey = onKeyStroke('Escape', (event: KeyboardEvent) => {
                if (event.target instanceof HTMLElement && focusNUITargets.includes(event.target.tagName.toLowerCase())) {
                    return;
                }

                if (isDashboardSidebarSlideoverOpen.value || overlay.overlays.some((o) => o.isOpen)) return;

                void toggleTablet(false);
            });

            onScopeDispose(() => stopEscapeKey());

            useEventListener(window, 'focusin', onFocusHandler, true);
            useEventListener(window, 'focusout', onFocusHandler, true);
        }
    },
});

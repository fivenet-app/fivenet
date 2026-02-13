import { useSettingsStore } from '~/stores/settings';

export default defineNuxtPlugin({
    name: 'update',
    dependsOn: ['nuxt-update', 'i18n:plugin'],
    enforce: 'post',
    parallel: true,

    hooks: {
        'custom:update_check:update': async (version) => {
            const logger = useLogger('🆕 Update Check');

            const settingsStore = useSettingsStore();
            logger.info(
                'Detected new version',
                version,
                'current version',
                APP_VERSION,
                'stored version',
                settingsStore.version,
            );

            if (version === 'UNKNOWN') return;

            const t = useNuxtApp().$i18n.t;
            const toast = useToast();
            toast.add({
                title: t('system.update_available.title', { version: version }),
                description: t('system.update_available.content'),
                actions: [
                    {
                        label: t('common.refresh'),
                        onClick: () => reloadNuxtApp({ persistState: false, force: true }),
                    },
                ],
                icon: 'i-mdi-update',
                color: 'primary',
                duration: 20000,
                close: false,
            });
        },
    },
});

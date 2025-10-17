import { useSettingsStore } from '~/stores/settings';

export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();

    const logger = useLogger('🆕 Update Check');

    $update.on('update', async (version) => {
        const settingsStore = useSettingsStore();
        logger.info('Detected new version', version, 'current version', APP_VERSION, 'stored version', settingsStore.version);

        if (version === 'UNKNOWN') return;

        if (APP_VERSION !== version) {
            settingsStore.setUpdateAvailable(version as string);
        }
    });
});

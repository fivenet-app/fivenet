import { useSettingsStore } from '~/store/settings';

export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();

    $update.on('update', async (version) => {
        const settings = useSettingsStore();
        console.info(
            'Update Check: Detected new version',
            version,
            'app version',
            APP_VERSION,
            'settings store version',
            settings.version,
        );

        if (version === 'UNKNOWN') {
            return;
        }

        if (APP_VERSION !== version) {
            settings.setUpdateAvailable(version as string);
        }
    });
});

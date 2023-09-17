import { useConfigStore } from '~/store/config';
import { useSettingsStore } from '~/store/settings';

export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();

    $update.on('update', (version) => {
        const settings = useSettingsStore();
        console.info(
            'Update Check: Detected new version',
            version,
            'app version',
            __APP_VERSION__,
            'settings store version',
            settings.getVersion,
        );
        const config = useConfigStore();
        config.updateAvailable = version;
    });
});

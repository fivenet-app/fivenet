import { useConfigStore } from '~/store/config';
import { useSettingsStore } from '~/store/settings';

export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();

    $update.on('update', (version) => {
        const settings = useSettingsStore();
        console.info('Update Check detected a new version', version, __APP_VERSION__, settings.getVersion);
        const config = useConfigStore();
        config.updateAvailable = version;
    });
});

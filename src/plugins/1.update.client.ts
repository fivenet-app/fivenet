import { useSettingsStore } from '~/store/settings';

export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();

    $update.on('update', (version) => {
        const userSettings = useSettingsStore();

        console.warn('UPDATE RECEIVED', version, __APP_VERSION__, userSettings.getVersion);

        // TODO add translations
        if (confirm(`New FiveNet version ${version} available. Reload?`)) {
            location.reload();
        }
    });
});

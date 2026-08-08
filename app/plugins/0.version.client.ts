import { useClipboardStore } from '~/stores/clipboard';
import { useSearchesStore } from '~/stores/searches';
import { useSettingsStore } from '~/stores/settings';

export default defineNuxtPlugin({
    name: 'version-migrations',

    async setup() {
        const logger = useLogger('⚙️ Settings');
        const settingsStore = useSettingsStore();

        if (APP_VERSION !== settingsStore.version) {
            logger.info('Resetting app data because new version has been detected', settingsStore.version, APP_VERSION);

            useClipboardStore().clear();
            useSearchesStore().clear();
            settingsStore.setVersion(APP_VERSION);
        }

        // Remove legacy dashboard cookies on app start to avoid sending stale state back to the server.
        const cookies = await cookieStore.getAll();
        cookies.forEach((cookie) => {
            if (cookie?.name && cookie.name.startsWith('dashboard-')) {
                logger.info('Removing dashboard cookie:', cookie.name);
                cookieStore.delete(cookie.name);
            }
        });
    },
});

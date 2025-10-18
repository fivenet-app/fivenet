export default defineNuxtPlugin(() => {
    if (import.meta.server) return;

    const query = useRouter().currentRoute.value.query;

    if (query?.nui !== undefined) {
        const nuiQuery = query.nui as string;

        const logger = useLogger('🎮 NUI');
        const settingsStore = useSettingsStore();
        if (nuiQuery.toLowerCase() !== 'false') {
            settingsStore.setNuiSettings(true, nuiQuery);
            logger.info('Enabled NUI integration, resource:', settingsStore.nuiResourceName);
        } else {
            settingsStore.setNuiSettings(false);
            logger.info('Disabled NUI integration');
        }
    }

    if (query?.refreshApp !== undefined) {
        reloadNuxtApp({
            persistState: false,
            ttl: 7_500, // 7.5 seconds
        });
    }
});

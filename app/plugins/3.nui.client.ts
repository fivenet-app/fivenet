export default defineNuxtPlugin(() => {
    if (import.meta.server) return;

    const query = useRouter().currentRoute.value.query;

    if (query?.nui !== undefined) {
        const nuiQuery = query.nui as string;

        const logger = useLogger('🎮 NUI');
        const settings = useSettingsStore();
        if (nuiQuery.toLowerCase() !== 'false') {
            settings.setNuiSettings(true, nuiQuery);
            logger.info('Enabled NUI integration, resource:', settings.nuiResourceName);
        } else {
            settings.setNuiSettings(false);
            logger.info('Disabled NUI integration');
        }
    }

    if (query?.refreshApp !== undefined) {
        reloadNuxtApp({
            persistState: false,
            ttl: 8500, // 8.5 seconds
        });
    }
});

import type { RouteLocationNormalized } from 'vue-router';
import { useSettingsStore } from '~/stores/settings';

const logger = useLogger('🎮 NUI');

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    if (import.meta.server) {
        return;
    }

    const route = from ?? to;

    if (route.query?.nui !== undefined) {
        const nuiQuery = route.query.nui as string;

        const settings = useSettingsStore();
        if (nuiQuery.toLowerCase() !== 'false') {
            settings.setNuiSettings(true, nuiQuery);
            logger.info('Enabled NUI integration, resource:', settings.nuiResourceName);
        } else {
            settings.setNuiSettings(false, undefined);
            logger.info('Disabled NUI integration');
        }
    } else {
        logger.debug('No NUI query param detected');
    }

    if (route.query?.refreshApp !== undefined) {
        reloadNuxtApp({
            persistState: false,
            ttl: 10000,
        });
    }
});

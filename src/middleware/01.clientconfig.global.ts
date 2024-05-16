import { type RouteLocationNormalized } from 'vue-router';
import { useSettingsStore } from '~/store/settings';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    if (import.meta.server) {
        return;
    }

    const route = from ?? to;

    if (route.query?.nui !== undefined) {
        const nuiQuery = route.query.nui as string;

        const settings = useSettingsStore();
        if (nuiQuery.toLowerCase() !== 'false') {
            settings.setNuiDetails(true, nuiQuery);
            console.info('ClientCfg: Enabled NUI integration! Resource Name:', settings.nuiResourceName);
        } else {
            settings.setNuiDetails(false, undefined);
            console.info('ClientCfg: Disabled NUI integration!');
        }
    } else {
        console.debug('ClientCfg: No NUI query param detected.');
    }

    if (route.query?.refreshApp !== undefined) {
        reloadNuxtApp({
            persistState: false,
            ttl: 10000,
        });
    }
});

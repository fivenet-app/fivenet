import { type NavigationGuard, type RouteLocationNormalized } from 'vue-router';
import { useConfigStore } from '~/store/config';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const route = from ?? to;

        if (route.query?.nui !== undefined) {
            const nuiQuery = route.query.nui as string;
            const configStore = useConfigStore();
            const { setNuiDetails } = configStore;

            if (nuiQuery.toLowerCase() !== 'false') {
                setNuiDetails(true, nuiQuery);
                console.info('Enabled NUI integration! Resource Name:', configStore.nuiResourceName);
            } else {
                setNuiDetails(false, undefined);
                console.info('Disabled NUI integration!');
            }
        } else {
            console.debug('No NUI query param detected.');
        }

        if (route.query?.refreshApp !== undefined) {
            reloadNuxtApp({
                persistState: false,
                ttl: 10000,
            });
        }
    },
);

import { type NavigationGuard, type RouteLocationNormalized } from 'vue-router';
import { useConfigStore } from '~/store/config';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const route = from ?? to;
        if (route.query?.nui !== undefined) {
            const nuiQuery = route.query.nui as string;
            const configStore = useConfigStore();
            const { clientConfig } = storeToRefs(configStore);
            if (nuiQuery.toLowerCase() !== 'false') {
                clientConfig.value.nuiEnabled = true;
                clientConfig.value.nuiResourceName = nuiQuery;
                console.info('Enabled NUI integration! Resource Name:', clientConfig.value.nuiResourceName);
            } else {
                clientConfig.value.nuiEnabled = false;
                clientConfig.value.nuiResourceName = undefined;
                console.info('Disabled NUI integration!');
            }
        }

        if (route.query?.refreshApp !== undefined) {
            reloadNuxtApp({
                persistState: false,
                ttl: 10000,
            });
        }
    },
);

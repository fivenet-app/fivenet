import { type NavigationGuard, type RouteLocationNormalized } from 'vue-router';
import { useConfigStore } from '~/store/config';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const route = from ?? to;
        if (route.query?.nui !== undefined) {
            const nuiQuery = route.query.nui as string;
            const configStore = useConfigStore();
            const { nuiEnabled, nuiResourceName } = storeToRefs(configStore);

            if (nuiQuery.toLowerCase() !== 'false') {
                nuiEnabled.value = true;
                nuiResourceName.value = nuiQuery;
                console.info('Enabled NUI integration! Resource Name:', nuiResourceName.value);
            } else {
                nuiEnabled.value = false;
                nuiResourceName.value = undefined;
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

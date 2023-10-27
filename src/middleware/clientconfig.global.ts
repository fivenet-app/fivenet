import { type NavigationGuard, type RouteLocationNormalized } from 'vue-router';
import { useConfigStore } from '~/store/config';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const route = from ?? to;
        if (route.query?.nui) {
            const nuiQuery = route.query?.nui as string;
            const configStore = useConfigStore();
            const { clientConfig } = storeToRefs(configStore);

            if (nuiQuery.toLowerCase() !== 'false') {
                clientConfig.value.NUIEnabled = true;
                clientConfig.value.NUIResourceName = nuiQuery;
                console.info('Enabled NUI integration! Resource Name:', clientConfig.value.NUIResourceName);
            } else {
                clientConfig.value.NUIEnabled = false;
                clientConfig.value.NUIResourceName = undefined;
                console.info('Disabled NUI integration!');
            }
        }
    },
);

import { createWebHistory, createRouter, setupDataFetchingGuard, RouteRecordRaw, Router } from 'vue-router/auto';
import { dispatchNotification } from './components/notification';
import { store } from './store/store';
import slug from './utils/slugify';

export const router = createRouter({
    history: createWebHistory(),
    extendRoutes: (routes: RouteRecordRaw[]) => {
        return routes;
    },
}) as Router;

setupDataFetchingGuard(router);

router.beforeResolve((to, from) => {
    // Default is that a page requires authentication
    if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
        if (to.meta.authOnlyToken && store.state.auth?.accessToken) {
            return true;
        }

        // Check if user has access token
        if (store.state.auth?.accessToken || store.state.auth?.activeChar) {
            // Route has permission attached to it, check if user has required permission
            if (to.meta.permission) {
                const perm = slug(to.meta.permission as string);
                if (store.state.auth?.permissions.includes(perm)) {
                    // User has permission
                    return true;
                } else {
                    dispatchNotification({
                        title: "You don't have permission!",
                        content: 'No permission to navigate to ' + (to.name ? to.name?.toString() : to.path) + '.',
                        type: 'warning',
                    });

                    if (store.state.auth?.accessToken) {
                        return {
                            path: '/',
                        };
                    }
                }
            } else {
                // No route permission check, so we can go ahead and return
                return true;
            }
        }

        return {
            path: '/login',
            // save the location we were at to come back later
            query: { redirect: to.fullPath },
        };
    }
});

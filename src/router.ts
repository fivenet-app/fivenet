import slugify from 'slugify';
import { createWebHistory, createRouter, setupDataFetchingGuard, RouteRecordRaw } from 'vue-router/auto';
import { dispatchNotification } from './components/notification';
import store from './store';

const router = createRouter({
    history: createWebHistory(),
    extendRoutes: (routes: RouteRecordRaw[]) => {
        return routes;
    },
});

setupDataFetchingGuard(router);

router.beforeEach((to, from) => {
    // Default is that a page requires authentication
    if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
        // Check if user has access token
        if (store.state.accessToken && store.state.activeChar) {
            // Route has permission attached to it, check if user has required permission
            if (to.meta.permission) {
                const perm = slugify(to.meta.permission as string);
                if (store.state.permissions.includes(perm)) {
                    // User has permission
                    return;
                } else {
                    dispatchNotification({
                        title: "You don't have permission!",
                        content: "You don't have permission to go to " + (to.name ? to.name?.toString() : to.path) + '.',
                        type: 'warning',
                    });

                    if (store.state.accessToken) {
                        return {
                            path: '/overview',
                        };
                    }
                }
            } else {
                // No route permission check, so we can go ahead and return
                return;
            }
        }

        return {
            path: '/login',
            // save the location we were at to come back later
            query: { redirect: to.fullPath },
        };
    }
});

export default router;

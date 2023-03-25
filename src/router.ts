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

// Check user permission
router.beforeResolve((to, from) => {
    // Default is that a page requires authentication
    if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
        if (to.meta.authOnlyToken && store.state.auth?.accessToken) {
            return true;
        }

        // Check if user has access token
        if (store.state.auth?.accessToken) {
            // If the user has an acitve char, check for perms otherwise, redirect to char selector
            if (store.state.auth?.activeChar) {
                // Route has permission attached to it, check if user has required permission
                if (to.meta.permission) {
                    const perm = slug(to.meta.permission as string);
                    if (store.state.auth?.permissions.includes(perm)) {
                        // User has permission
                        return true;
                    } else {
                        dispatchNotification({
                            title: "You don't have permission!",
                            content: 'No permission to go to ' + (to.name ? to.name?.toString() : to.path) + '.',
                            type: 'warning',
                        });

                        if (store.state.auth?.accessToken) {
                            return {
                                name: 'Overview',
                            };
                        }
                    }
                } else {
                    // No route permission check, so we can go ahead and return
                    return true;
                }
            } else {
                // Only update the redirect query param if it isn't set already
                const redirect = to.query.redirect ?? to.fullPath;
                return {
                    name: 'Character Selector',
                    query: { redirect: redirect },
                };
            }
        }

        // Only update the redirect query param if it isn't set already
        const redirect = to.query.redirect ?? to.fullPath;
        return {
            name: 'Login',
            // save the location we were at to come back later
            query: { redirect: redirect },
        };
    }
});

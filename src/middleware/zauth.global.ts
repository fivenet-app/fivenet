import { NavigationGuard, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import slug from '~/utils/slugify';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        // Default is that a page requires authentication
        if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
            const store = useAuthStore();
            if (to.meta.authOnlyToken && store.$state.accessToken) {
                return true;
            }

            // Check if user has access token
            if (store.$state.accessToken) {
                // If the user has an acitve char, check for perms otherwise, redirect to char selector
                if (store.$state.activeChar) {
                    // Route has permission attached to it, check if user has required permission
                    if (to.meta.permission) {
                        const perm = slug(to.meta.permission as string);
                        if (store.$state.permissions.includes(perm)) {
                            // User has permission
                            return true;
                        } else {
                            useNotificationsStore().dispatchNotification({
                                title: "You don't have permission!",
                                content: 'No permission to go to ' + (to.name ? to.name?.toString() : to.path) + '.',
                                type: 'warning',
                            });

                            if (store.$state.accessToken) {
                                return {
                                    name: 'overview',
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
                        name: 'auth-character-selector',
                        query: { redirect: redirect },
                    };
                }
            }

            // Only update the redirect query param if it isn't set already
            const redirect = to.query.redirect ?? to.fullPath;
            return {
                name: 'auth-login',
                // save the location we were at to come back later
                query: { redirect: redirect },
            };
        }

        return true;
    }
);

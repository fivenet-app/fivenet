import { NavigationGuard, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import slug from '~/utils/slugify';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        // Default is that a page requires authentication
        if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
            const authStore = useAuthStore();
            if (to.meta.authOnlyToken && authStore.getAccessToken !== null) {
                return true;
            }

            // Check if user has access token
            if (authStore.getAccessToken !== null) {
                // If the user has an acitve char, check for perms otherwise, redirect to char selector
                if (authStore.getActiveChar !== null) {
                    // Route has permission attached to it, check if user has required permission
                    if (to.meta.permission) {
                        const perm = slug(to.meta.permission as string);
                        if (authStore.getPermissions.includes(perm)) {
                            // User has permission
                            return true;
                        } else {
                            useNotificationsStore().dispatchNotification({
                                title: "You don't have permission!",
                                content: 'No permission to go to ' + (to.name ? to.name?.toString() : to.path) + '.',
                                type: 'warning',
                            });

                            if (authStore.getAccessToken) {
                                return navigateTo({
                                    name: 'overview',
                                });
                            }
                        }
                    } else {
                        // No route permission check, so we can go ahead and return
                        return true;
                    }
                } else {
                    // Only update the redirect query param if it isn't set already
                    const redirect = to.query.redirect ?? to.fullPath;
                    return navigateTo({
                        name: 'auth-character-selector',
                        query: { redirect: redirect },
                    });
                }
            }

            // Only update the redirect query param if it isn't set already
            const redirect = to.query.redirect ?? to.fullPath;
            return navigateTo({
                name: 'auth-login',
                // save the location we were at to come back later
                query: { redirect: redirect },
            });
        }

        return true;
    }
);

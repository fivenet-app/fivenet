import { type RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, _: RouteLocationNormalized) => {
    const authStore = useAuthStore();
    const { activeChar, accessToken, lastCharID } = storeToRefs(authStore);

    // Default is that a page requires authentication, but if it doesn't exit quickly
    if (to.meta.requiresAuth === false) {
        return true;
    }

    // Auth token is not null and only needed
    if (to.meta.authOnlyToken && accessToken.value !== null) {
        return true;
    }

    // Check if user has access token
    if (accessToken.value !== null) {
        // If the user has an acitve char, check for perms otherwise, redirect to char selector
        if (activeChar.value === null) {
            // If we don't have an active char, but a last char ID set, try to choose it and immidiately continue
            if (lastCharID.value > 0) {
                const { setActiveChar, setPermissions, setJobProps } = authStore;
                try {
                    await authStore.chooseCharacter(authStore.lastCharID);
                } catch (e) {
                    setActiveChar(null);
                    setPermissions([]);
                    setJobProps(undefined);
                }
            }

            if (activeChar.value === null) {
                // Only update the redirect query param if it isn't set already
                const redirect = to.query.redirect ?? to.fullPath;
                return navigateTo({
                    name: 'auth-character-selector',
                    query: { redirect },
                });
            }
        }

        // No route permission check, so we can go ahead and return true
        if (!to.meta.permission) {
            return true;
        }

        // Route has permission attached to it, check if user "can" go there
        if (can(to.meta.permission)) {
            // User has permission
            return true;
        } else {
            useNotificatorStore().add({
                title: { key: 'notifications.auth.no_permission.title', parameters: {} },
                description: {
                    key: 'notifications.auth.no_permission.content',
                    parameters: { path: to.name ? toTitleCase(to.name?.toString()) : to.path },
                },
                type: 'warning',
            });

            if (accessToken.value !== null) {
                return navigateTo({
                    name: 'overview',
                });
            }
        }
    }

    // Only update the redirect query param if it isn't set already
    const redirect = to.query.redirect ?? to.fullPath;
    return navigateTo({
        name: 'auth-login',
        // save the location we were at to come back later
        query: { redirect },
    });
});

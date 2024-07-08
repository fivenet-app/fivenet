import { parseQuery, type RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    const authStore = useAuthStore();
    const { activeChar, username, lastCharID } = storeToRefs(authStore);

    // Default is that a page requires authentication, but if it doesn't exit quickly
    if (to.meta.requiresAuth === false) {
        if ((to.meta.redirectIfAuthed === undefined || to.meta.redirectIfAuthed) && username.value !== null) {
            const settingsStore = useSettingsStore();
            const { startpage } = storeToRefs(settingsStore);

            const redirect = from.query.redirect ?? startpage.value ?? '/overview';
            const path = redirect || '/overview';
            const url = new URL('https://example.com' + path);

            // @ts-ignore the route should be valid, as we test it against a valid URL list
            return navigateTo({
                path: url.pathname,
                query: parseQuery(url.search),
                hash: url.hash,
            });
        }
        return true;
    }

    // Auth token is not null and only needed by route
    if (to.meta.authTokenOnly && username.value !== null) {
        return true;
    }

    const redirect = to.query.redirect ?? to.fullPath;
    // Check if user has access token
    if (username.value !== null) {
        // If the user has an acitve char, check for perms otherwise, redirect to char selector
        if (activeChar.value === null) {
            // If we don't have an active char, but a last char ID set, try to choose it and immidiately continue
            if (lastCharID.value > 0) {
                const { setActiveChar, setPermissions, setJobProps } = authStore;
                try {
                    await authStore.chooseCharacter(authStore.lastCharID);

                    if (redirect !== undefined) {
                        const path = redirect || '/overview';
                        const url = new URL('https://example.com' + path);
                        // @ts-ignore the route should be valid, as we test it against a valid list of URLs
                        return await navigateTo({
                            path: url.pathname,
                            query: parseQuery(url.search),
                            hash: url.hash,
                        });
                    } else {
                        // @ts-ignore the route should be valid, as we test it against a valid list of URLs
                        const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
                        return await navigateTo(target);
                    }
                } catch (e) {
                    setActiveChar(null);
                    setPermissions([]);
                    setJobProps(undefined);

                    return navigateTo({
                        name: 'auth-character-selector',
                        query: { redirect: redirect },
                    });
                }
            }

            if (activeChar.value === null) {
                // Only update the redirect query param if it isn't set already
                return navigateTo({
                    name: 'auth-character-selector',
                    query: { redirect: redirect },
                });
            }
        }

        // No route permission check, so we can go ahead and return true
        if (!to.meta.permission) {
            return true;
        }

        // Route has permission attached to it, check if user "can" go there
        if (can(to.meta.permission).value) {
            // User has permission
            return true;
        } else {
            useNotificatorStore().add({
                title: { key: 'notifications.auth.no_permission.title', parameters: {} },
                description: {
                    key: 'notifications.auth.no_permission.content',
                    parameters: { path: to.name ? toTitleCase(to.name?.toString()) : to.path },
                },
                type: NotificationType.WARNING,
            });

            if (username.value !== null) {
                return navigateTo({
                    name: 'overview',
                });
            }
        }
    }

    // Only update the redirect query param if it isn't set already
    return navigateTo({
        name: 'auth-login',
        // save the location we were at to come back later
        query: { redirect: redirect },
    });
});

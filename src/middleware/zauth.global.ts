import { RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { toDate } from '~/utils/time';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    const authStore = useAuthStore();
    const { activeChar, accessToken, lastCharID } = storeToRefs(authStore);

    // Default is that a page requires authentication, but if it doesn't exit quickly
    if (to.meta.hasOwnProperty('requiresAuth') && !to.meta.requiresAuth) {
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
                await chooseCharacter();
            }

            if (activeChar.value === null) {
                // Only update the redirect query param if it isn't set already
                const redirect = to.query.redirect ?? to.fullPath;
                return navigateTo({
                    name: 'auth-character-selector',
                    query: { redirect: redirect },
                });
            }
        }

        // No route permission check, so we can go ahead and return
        if (!to.meta.permission) {
            return true;
        }

        // Route has permission attached to it, check if user "can" go there
        if (can(to.meta.permission)) {
            // User has permission
            return true;
        } else {
            useNotificatorStore().dispatchNotification({
                title: { key: 'notifications.auth.no_permission.title', parameters: {} },
                content: {
                    key: 'notifications.auth.no_permission.content',
                    parameters: [to.name ? toTitleCase(to.name?.toString()) : to.path],
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
        query: { redirect: redirect },
    });
});

async function chooseCharacter(): Promise<any> {
    const authStore = useAuthStore();
    const { setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;
    const { lastCharID } = storeToRefs(authStore);

    const { $grpc } = useNuxtApp();

    try {
        const call = $grpc.getAuthClient().chooseCharacter({
            charId: lastCharID.value,
        });
        const { response } = await call;

        setAccessToken(response.token, toDate(response.expires) as null | Date);
        setActiveChar(response.char!);
        setPermissions(response.permissions);
        if (response.jobProps) {
            setJobProps(response.jobProps!);
        } else {
            setJobProps(null);
        }
    } catch (e) {
        setActiveChar(null);
        setPermissions([]);
        setJobProps(null);
    }
}

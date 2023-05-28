import { RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import slug from '~/utils/slugify';
import { toDate } from '~/utils/time';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    // Default is that a page requires authentication
    if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
        const authStore = useAuthStore();
        const { activeChar, accessToken, permissions, lastCharID } = storeToRefs(authStore);
        const { setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;

        if (to.meta.authOnlyToken && accessToken.value !== null) {
            return true;
        }

        // Check if user has access token
        if (accessToken.value !== null) {
            // If the user has an acitve char, check for perms otherwise, redirect to char selector
            if (activeChar.value === null) {
                // If we don't have an active char, but a last char ID set, try to choose it and immidiately continue
                if (lastCharID.value > 0) {
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
                        // Only update the redirect query param if it isn't set already
                        const redirect = to.query.redirect ?? to.fullPath;
                        return navigateTo({
                            name: 'auth-character-selector',
                            query: { redirect: redirect },
                        });
                    }
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
            if (permissions.value.includes('superuser')) {
                return true;
            }

            // Route has permission attached to it, check if user has required permission
            if (to.meta.permission) {
                const perm = slug(to.meta.permission as string);
                if (permissions.value.includes(perm)) {
                    // User has permission
                    return true;
                } else {
                    useNotificationsStore().dispatchNotification({
                        title: "You don't have permission!",
                        content: 'No permission to go to ' + (to.name ? to.name?.toString() : to.path) + '.',
                        type: 'warning',
                    });

                    if (accessToken.value) {
                        return navigateTo({
                            name: 'overview',
                        });
                    }
                }
            } else {
                // No route permission check, so we can go ahead and return
                return true;
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
});

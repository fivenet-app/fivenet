import { ChooseCharacterRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import slug from '~/utils/slugify';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    // Default is that a page requires authentication
    if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
        const authStore = useAuthStore();
        const { activeChar, accessToken, permissions, lastCharID,
            setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;

        if (to.meta.authOnlyToken && accessToken !== null) {
            return true;
        }

        // Check if user has access token
        if (accessToken !== null) {
            // If the user has an acitve char, check for perms otherwise, redirect to char selector
            if (activeChar === null) {
                // If we don't have an active char, but a last char ID set, try to choose it and immidiately continue
                if (lastCharID > 0) {
                    const { $grpc } = useNuxtApp();

                    const req = new ChooseCharacterRequest();
                    req.setCharId(lastCharID);

                    try {
                        const resp = await $grpc.getAuthClient().chooseCharacter(req, null);

                        setAccessToken(resp.getToken(), toDate(resp.getExpires()) as null | Date);
                        setActiveChar(resp.getChar()!);
                        setPermissions(resp.getPermissionsList());
                        if (resp.hasJobProps()) {
                            setJobProps(resp.getJobProps()!);
                        } else {
                            setJobProps(null);
                        }
                    } catch (e) {
                        $grpc.handleRPCError(e as RpcError);

                        // Only update the redirect query param if it isn't set already
                        const redirect = to.query.redirect ?? to.fullPath;
                        return navigateTo({
                            name: 'auth-character-selector',
                            query: { redirect: redirect },
                        });
                    }
                }

                if (activeChar === null) {
                    // Only update the redirect query param if it isn't set already
                    const redirect = to.query.redirect ?? to.fullPath;
                    return navigateTo({
                        name: 'auth-character-selector',
                        query: { redirect: redirect },
                    });
                }
            }
            if (permissions.includes('superuser')) {
                return true;
            }

            // Route has permission attached to it, check if user has required permission
            if (to.meta.permission) {
                const perm = slug(to.meta.permission as string);
                if (permissions.includes(perm)) {
                    // User has permission
                    return true;
                } else {
                    useNotificationsStore().dispatchNotification({
                        title: "You don't have permission!",
                        content: 'No permission to go to ' + (to.name ? to.name?.toString() : to.path) + '.',
                        type: 'warning',
                    });

                    if (accessToken) {
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

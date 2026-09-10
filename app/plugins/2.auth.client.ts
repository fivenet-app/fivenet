import { parseQuery, type RouteLocationNormalized } from 'vue-router';
import { canAccessRoute, getRoutePermissionDeniedNotification } from '~/composables/auth/routePermission';
import { isSetupBypassRoute } from '~/composables/setup';

export function getLoginRedirect(to: RouteLocationNormalized) {
    return {
        name: 'auth-login',
        query: { redirect: getRedirectPath((to.query.redirect ?? to.fullPath) as string) },
        replace: true,
    } as const;
}

function loginRedirect(to: RouteLocationNormalized) {
    return navigateTo(getLoginRedirect(to));
}

export function getCharacterSelectorRedirect(to: RouteLocationNormalized) {
    return {
        name: 'auth-character-selector',
        query: { redirect: getRedirectPath((to.query.redirect ?? to.fullPath) as string) },
        replace: true,
    } as const;
}

function characterSelectorRedirect(to: RouteLocationNormalized) {
    return navigateTo(getCharacterSelectorRedirect(to));
}

export default defineNuxtPlugin({
    name: 'auth',

    async setup(nuxtApp) {
        addRouteMiddleware(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
            const appConfigResult = await nuxtApp.$appConfigPromise;
            if (appConfigResult === undefined) return;

            const authStore = useAuthStore();
            const { can, username } = useAuth();
            const appConfig = useAppConfig();
            const isSetupRoute = to.path.startsWith('/settings/setup');

            if (to.meta.requiresAuth === false) {
                // A public page may be reached after a reload. Bootstrap from
                // the account cookie before deciding whether it should redirect.
                const account = await authStore.ensureAccountSession();
                if (account.kind === 'ready' && to.meta.redirectIfAuthed !== false && username.value !== null) {
                    const url = getRedirect(from);
                    // @ts-expect-error route is validated by parseRedirectURL/getRedirect
                    return navigateTo({ path: url.pathname, query: parseQuery(url.search), hash: url.hash, replace: true });
                }

                return true;
            }

            const account = await authStore.ensureAccountSession();
            if (account.kind === 'needs-login' || account.kind === 'temporary-failure') return loginRedirect(to);

            const hasConfigAdminAccess = async (): Promise<boolean> => {
                if (can('internal.Superuser/ConfigAdmin').value) return true;

                try {
                    await authStore.refreshAccountSession();
                } catch (_) {
                    return false;
                }

                return can('internal.Superuser/ConfigAdmin').value;
            };

            if (isSetupRoute && appConfig.setupComplete !== false) {
                return (await hasConfigAdminAccess()) ? true : navigateTo({ name: 'overview', replace: true });
            }

            if (!isSetupRoute && !isSetupBypassRoute(to.path) && appConfig.setupComplete === false) {
                if (await hasConfigAdminAccess()) {
                    return navigateTo({
                        path: '/settings/setup',
                        query: { redirect: getRedirectPath((to.query.redirect ?? to.fullPath) as string) },
                        replace: true,
                    });
                }
            }

            // Existing authTokenOnly routes are account-only routes. Every
            // other protected route requires a server-confirmed character.
            if (!to.meta.authTokenOnly) {
                const character = await authStore.ensureCharacterSession();
                if (character.kind === 'needs-login') return loginRedirect(to);
                if (character.kind === 'needs-character' || character.kind === 'temporary-failure') {
                    return characterSelectorRedirect(to);
                }
            }

            if (!to.meta.permission || canAccessRoute(to)) return true;

            useNotificationsStore().add(getRoutePermissionDeniedNotification(to));
            return navigateTo({ name: 'overview', replace: true });
        });
    },
});

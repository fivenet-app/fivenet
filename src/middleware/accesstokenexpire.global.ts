import { NavigationGuard, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, _: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const authStore = useAuthStore();
        const { getAccessTokenExpiration, setAccessToken } = authStore;

        const expiration = getAccessTokenExpiration;

        // Check if we have an expiration time, make sure the token isn't expired (yet)
        if (expiration !== null) {
            const now = new Date();
            // Token expired, redirect to login
            if (expiration <= now) {
                console.info('Auth: Token expired, redirecting to login.');
                setAccessToken(null, null);

                // Only update the redirect query param if it isn't set already
                const redirect = to.query.redirect ?? to.fullPath;
                return navigateTo({
                    name: 'auth-login',
                    // save the location we were at to come back later
                    query: { redirect: redirect },
                });
            }
        }
    },
);

import { NavigationGuard, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';

export default defineNuxtRouteMiddleware(
    (to: RouteLocationNormalized, from: RouteLocationNormalized): ReturnType<NavigationGuard> => {
        const authStore = useAuthStore();
        const { getAccessTokenExpiration, setAccessToken } = authStore;

        const expiration = getAccessTokenExpiration;

        // Check if we have an expiration time, make sure the token isn't expired (yet)
        if (expiration !== null) {
            const now = new Date();
            // Token expired, redirect to login
            if (expiration <= now) {
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

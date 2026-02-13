import { defineStore } from 'pinia';

export const authUserTokenKey = 'fivenet:user_token_v1';

/**
 * Pinia store for managing user sessions.
 */
export const useAuthSessionStore = defineStore(
    'auth_session',
    () => {
        /**
         * User token of this session.
         */
        const userToken = ref<string | null>(null);

        /**
         * User token info based on the user token claims.
         */
        const userInfo = ref<JWTUserInfoClaims | null>(null);

        /**
         * Get the user token for this session. This will be used by gRPC interceptors to authenticate requests.
         */
        const getUserToken = (): string | null => userToken.value;

        /**
         * Set the user token for this session. This will be used by gRPC interceptors to authenticate requests.
         * @param token JWT user info token
         */
        const setUserToken = (token: string | null) => {
            userToken.value = token;
            if (token) {
                sessionStorage.setItem(authUserTokenKey, token);
                const uinfo = new JWTUserInfoClaims(token);
                userInfo.value = uinfo;
            } else {
                sessionStorage.removeItem(authUserTokenKey);
                userInfo.value = null;
            }
        };

        return {
            userToken,
            userInfo,

            getUserToken,
            setUserToken,
        };
    },
    {
        persist: {
            storage: piniaPluginPersistedstate.sessionStorage(),
            pick: ['userToken', 'userInfo'],
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            afterHydrate: (ctx: any) => {
                const store = ctx.store;
                if (store.userInfo && typeof store.userInfo.expiration === 'string') {
                    store.userInfo.expiration = new Date(store.expiration);
                }
            },
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthSessionStore, import.meta.hot));
}

import { defineStore } from 'pinia';

/**
 * Pinia store for managing user sessions.
 */
export const useAuthSessionStore = defineStore('auth_session', () => {
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
     * Sets the active character token for this tab. Account authentication
     * is cookie-backed; this token is intentionally not persisted.
     */
    const setUserToken = (token: string | null) => {
        userToken.value = token;
        if (token) {
            const uinfo = new JWTUserInfoClaims(token);
            userInfo.value = uinfo;
        } else {
            userInfo.value = null;
        }
    };

    return {
        userToken,
        userInfo,

        getUserToken,
        setUserToken,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthSessionStore, import.meta.hot));
}

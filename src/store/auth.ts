import { User } from '@fivenet/gen/resources/users/users_pb';
import { StoreDefinition } from 'pinia';
import { defineStore } from 'pinia';

export interface AuthState {
    accessToken: null | string;
    accessTokenExpiration: null | Date;
    lastCharID: number;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: Array<String>;
}

export const useAuthStore = defineStore('auth', {
    state: () => ({
        // Persisted to Local Storage
        accessToken: null as null | string,
        accessTokenExpiration: null as null | Date,
        lastCharID: 0 as number,
        // Temporary
        activeChar: null as null | User,
        loggingIn: false as boolean,
        loginError: null as null | string,
        permissions: [] as Array<String>,
    }),
    persist: {
        paths: ['accessToken', 'accessTokenExpiration', 'lastCharID'],
    },
    actions: {
        loginStart(): void {
            this.loggingIn = true;
        },
        loginStop(errorMessage: null | string): void {
            this.loggingIn = false;
            this.loginError = errorMessage;
        },
        setAccessToken(accessToken: null | string, expiration: null | number | Date): void {
            this.accessToken = accessToken;
            if (typeof expiration === 'number') expiration = new Date(expiration);
            this.accessTokenExpiration = expiration;
        },
        setActiveChar(char: null | User): void {
            this.activeChar = char;
            this.lastCharID = char ? char.getUserId() : this.lastCharID;
        },
        setPermissions(permissions: string[]): void {
            this.permissions = permissions;
        },
        async clear(): Promise<void> {
            this.setAccessToken(null, null);
            this.setActiveChar(null);
            this.setPermissions([]);
        },
    },
    getters: {
        getAccessToken: (state): null | string => state.accessToken,
        getAccessTokenExpiration(state): null | Date {
            if (typeof state.accessTokenExpiration === 'string')
                state.accessTokenExpiration = new Date(Date.parse(state.accessTokenExpiration));

            return state.accessTokenExpiration;
        },
        getActiveChar: (state): null | User => state.activeChar,
        getPermissions: (state): Array<String> => state.permissions,
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore as unknown as StoreDefinition, import.meta.hot));
}

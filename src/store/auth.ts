import { User } from '@fivenet/gen/resources/users/users_pb';
import { StoreDefinition } from 'pinia';
import { defineStore } from 'pinia';

export interface AuthState {
    accessToken: null | string;
    lastCharID: number;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: Array<String>;
}

export const useAuthStore = defineStore('authStore', {
    state: () => ({
        // Persisted to Local Storage
        accessToken: null as null | string,
        lastCharID: 0 as number,
        // Temporary
        activeChar: null as null | User,
        loggingIn: false as boolean,
        loginError: null as null | string,
        permissions: [] as Array<String>,
    }),
    persist: {
        paths: ['accessToken', 'lastCharID'],
    },
    actions: {
        loginStart(): void {
            this.loggingIn = true;
        },
        loginStop(errorMessage: null | string): void {
            this.loggingIn = false;
            this.loginError = errorMessage;
        },
        updateAccessToken(accessToken: null | string): void {
            this.accessToken = accessToken;
        },
        updateActiveChar(char: null | User): void {
            this.activeChar = char;
            this.lastCharID = char ? char.getUserId() : this.lastCharID;
        },
        updatePermissions(permissions: string[]): void {
            this.permissions = permissions;
        },
        async clear(): Promise<void> {
            this.updateAccessToken(null);
            this.updateActiveChar(null);
            this.updatePermissions([]);
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore as unknown as StoreDefinition, import.meta.hot));
}

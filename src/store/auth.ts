import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { JobProps } from '~~/gen/ts/resources/users/jobs';
import { User } from '~~/gen/ts/resources/users/users';
import { useNotificationsStore } from './notifications';

export type JobPropsState = {
    quickButtons: String[];
};

export interface AuthState {
    accessToken: null | string;
    accessTokenExpiration: null | Date;
    lastCharID: number;
    accountID: bigint;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: String[];
    jobProps: null | JobPropsState;
}

export const useAuthStore = defineStore('auth', {
    state: () => ({
        // Persisted to Local Storage
        accessToken: null as null | string,
        accessTokenExpiration: null as null | Date,
        lastCharID: 0 as number,
        // Temporary
        accountID: 0n as bigint,
        activeChar: null as null | User,
        loggingIn: false as boolean,
        loginError: null as null | string,
        permissions: [] as String[],
        jobProps: {
            quickButtons: [],
        } as null | JobPropsState,
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
        setAccessToken(accessToken: null | string, expiration: null | bigint | Date): void {
            this.accessToken = accessToken;
            if (typeof expiration === 'bigint') expiration = new Date(expiration.toString());
            this.accessTokenExpiration = expiration;
        },
        setActiveChar(char: null | User): void {
            this.activeChar = char;
            this.lastCharID = char ? char.userId : this.lastCharID;
        },
        setPermissions(permissions: string[]): void {
            this.permissions = permissions.sort();
        },
        setJobProps(jp: null | JobProps): void {
            if (jp === null) {
                this.jobProps = null;
            } else {
                this.jobProps = {
                    quickButtons: jp.quickButtons.split(';').filter((v) => v !== ''),
                };
            }
        },
        clearAuthInfo(): void {
            this.setAccessToken(null, null);
            this.setActiveChar(null);
            this.setPermissions([]);
            this.setJobProps(null);
        },

        // GRPC Calls
        async doLogin(username: string, password: string): Promise<void> {
            return new Promise(async (res, rej) => {
                // Start login
                this.loginStart();
                this.setActiveChar(null);
                this.setPermissions([]);

                const { $grpc } = useNuxtApp();
                try {
                    const call = $grpc.getUnAuthClient().login({
                        username: username,
                        password: password,
                    });
                    const { response } = await call;

                    this.loginStop(null);
                    this.setAccessToken(response.token, toDate(response.expires) as null | Date);

                    return res();
                } catch (e) {
                    this.loginStop((e as RpcError).message);
                    this.setAccessToken(null, null);
                    $grpc.handleError(e as RpcError);
                    return rej(e as RpcError);
                }
            });
        },
        async doLogout(): Promise<void> {
            return new Promise(async (res, rej) => {
                const { $grpc } = useNuxtApp();
                try {
                    await $grpc.getAuthClient().logout({});
                    this.clearAuthInfo();

                    return res();
                } catch (e) {
                    $grpc.handleError(e as RpcError);

                    useNotificationsStore().dispatchNotification({
                        title: { key: 'notifications.auth.error_logout.title', parameters: [] },
                        content: { key: 'notifications.auth.error_logout.content', parameters: [(e as RpcError).message] },
                        type: 'error',
                    });
                    this.clearAuthInfo();

                    return rej(e as RpcError);
                }
            });
        },
    },
    getters: {
        getAccessTokenExpiration(state): null | Date {
            if (typeof state.accessTokenExpiration === 'string')
                state.accessTokenExpiration = new Date(Date.parse(state.accessTokenExpiration));

            return state.accessTokenExpiration;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore as unknown as StoreDefinition, import.meta.hot));
}

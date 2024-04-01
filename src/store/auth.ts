import { defineStore, type StoreDefinition } from 'pinia';
import { parseQuery } from 'vue-router';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { type JobProps } from '~~/gen/ts/resources/users/jobs';
import { User } from '~~/gen/ts/resources/users/users';

export interface AuthState {
    accessToken: null | string;
    accessTokenExpiration: null | Date;
    username: null | string;
    lastCharID: number;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: string[];
    jobProps: null | JobProps;
}

export const useAuthStore = defineStore('auth', {
    state: () =>
        ({
            // Persisted to Local Storage
            accessToken: null,
            accessTokenExpiration: null,
            lastCharID: 0,
            username: null,
            // Temporary
            activeChar: null,
            loggingIn: false,
            loginError: null,
            permissions: [],
            jobProps: {
                job: '',
                theme: 'defaultTheme',
                radioFrequency: undefined,
                quickButtons: {},
                logoUrl: undefined,
            } as JobProps,
        }) as AuthState,
    persist: {
        paths: ['accessToken', 'accessTokenExpiration', 'lastCharID', 'username'],
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
            if (typeof expiration === 'bigint') {
                expiration = new Date(expiration.toString());
            }
            this.accessTokenExpiration = expiration;
        },
        setActiveChar(char: null | User): void {
            this.activeChar = char;
            this.lastCharID = char ? char.userId : this.lastCharID;
        },
        setPermissions(permissions: string[]): void {
            this.permissions.length = 0;
            this.permissions.push(...permissions.sort());
        },
        setJobProps(jp: undefined | JobProps): void {
            if (jp === undefined) {
                this.jobProps = null;
            } else {
                this.jobProps = jp;
            }
        },
        clearAuthInfo(): void {
            this.setAccessToken(null, null);
            this.setActiveChar(null);
            this.setPermissions([]);
            this.setJobProps(undefined);
        },

        // GRPC Calls
        async doLogin(username: string, password: string): Promise<void> {
            // Start login
            this.loginStart();
            this.setActiveChar(null);
            this.setPermissions([]);

            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getUnAuthClient().login({ username, password });
                const { response } = await call;

                this.loginStop(null);
                this.setAccessToken(response.token, toDate(response.expires));
                this.username = username;
            } catch (e) {
                this.loginStop((e as RpcError).message);
                this.setAccessToken(null, null);
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
        async doLogout(): Promise<void> {
            const { $grpc } = useNuxtApp();
            try {
                await $grpc.getAuthClient().logout({});
                this.clearAuthInfo();
            } catch (e) {
                $grpc.handleError(e as RpcError);

                useNotificatorStore().dispatchNotification({
                    title: { key: 'notifications.auth.error_logout.title', parameters: {} },
                    content: {
                        key: 'notifications.auth.error_logout.content',
                        parameters: { msg: (e as RpcError).message },
                    },
                    type: 'error',
                });
                this.clearAuthInfo();

                throw e;
            }
        },
        async chooseCharacter(charId?: number): Promise<void> {
            if (charId === undefined) {
                charId = this.lastCharID;
            }

            const { $grpc } = useNuxtApp();
            try {
                if (this.lastCharID !== charId) {
                    useClipboardStore().clear();
                }

                const call = $grpc.getAuthClient().chooseCharacter({ charId });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                this.setAccessToken(response.token, toDate(response.expires));
                this.setActiveChar(response.char);
                this.setPermissions(response.permissions);
                this.setJobProps(response.jobProps);

                if (useRoute().query.redirect !== undefined) {
                    const path = useRoute().query.redirect?.toString() || '/overview';
                    const url = new URL('https://example.com' + path);
                    // @ts-ignore the route should be valid, as we test it against a valid URL list
                    await navigateTo({
                        path: url.pathname,
                        query: parseQuery(url.search),
                        hash: url.hash,
                    });
                } else {
                    // @ts-ignore the route should be valid, as we test it against a valid URL list
                    const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
                    await navigateTo(target);
                }
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
    },
    getters: {
        isSuperuser(state): boolean {
            return state.permissions.includes('superuser');
        },
        getAccessTokenExpiration(state): null | Date {
            if (typeof state.accessTokenExpiration === 'string') {
                state.accessTokenExpiration = new Date(Date.parse(state.accessTokenExpiration));
            }

            return state.accessTokenExpiration;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore as unknown as StoreDefinition, import.meta.hot));
}

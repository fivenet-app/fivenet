import { defineStore, type StoreDefinition } from 'pinia';
import { parseQuery } from 'vue-router';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { Job, type JobProps } from '~~/gen/ts/resources/users/jobs';
import { User } from '~~/gen/ts/resources/users/users';
import type { SetSuperUserModeRequest } from '~~/gen/ts/services/auth/auth';

export interface AuthState {
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
        paths: ['accessTokenExpiration', 'lastCharID', 'username'],
    },
    actions: {
        loginStart(): void {
            this.loggingIn = true;
        },
        loginStop(errorMessage: null | string): void {
            this.loggingIn = false;
            this.loginError = errorMessage;
        },
        setAccessTokenExpiration(expiration: null | string | Date): void {
            if (typeof expiration === 'string') {
                expiration = new Date(expiration);
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
            this.setAccessTokenExpiration(null);
            this.setActiveChar(null);
            this.setPermissions([]);
            this.setJobProps(undefined);
            this.username = null;
        },

        // GRPC Calls
        async doLogin(username: string, password: string): Promise<void> {
            // Start login
            this.loginStart();
            this.setActiveChar(null);
            this.setPermissions([]);

            const { $grpc } = useNuxtApp();
            try {
                const call = $grpc.getAuthClient().login({ username, password });
                const { response } = await call;

                this.loginStop(null);

                this.username = username;
                if (!response.char) {
                    console.info('Simple login response received, redirecting to char selector');
                    this.setAccessTokenExpiration(toDate(response.expires));

                    const route = useRoute();

                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                } else {
                    console.info('Received fast-tracked char response for char id:', response.char.char?.userId);
                    this.setActiveChar(response.char.char ?? null);
                    this.setAccessTokenExpiration(toDate(response.char.expires));
                    this.setPermissions(response.char.permissions);
                    this.setJobProps(response.char.jobProps);

                    // @ts-ignore the route should be valid, as we test it against a valid URL list
                    const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
                    await navigateTo(target);
                }
            } catch (e) {
                this.loginStop((e as RpcError).message);
                this.setAccessTokenExpiration(null);
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

                useNotificatorStore().add({
                    title: { key: 'notifications.auth.error_logout.title', parameters: {} },
                    description: {
                        key: 'notifications.auth.error_logout.content',
                        parameters: { msg: (e as RpcError).message },
                    },
                    type: NotificationType.ERROR,
                });
                this.clearAuthInfo();

                throw e;
            }
        },
        async chooseCharacter(charId?: number, redirect?: boolean): Promise<void> {
            if (charId === undefined) {
                charId = this.lastCharID;
            }

            const { $grpc } = useNuxtApp();
            try {
                if (this.lastCharID !== charId) {
                    useClipboardStore().clear();
                }

                const call = $grpc.getAuthClient().chooseCharacter({ charId: charId });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                this.setActiveChar(response.char);
                this.setAccessTokenExpiration(toDate(response.expires));
                this.setPermissions(response.permissions);
                this.setJobProps(response.jobProps);

                if (redirect) {
                    const redirect = useRoute().query.redirect ?? useSettingsStore().startpage ?? '/overview';
                    const path = redirect || '/overview';
                    const url = new URL('https://example.com' + path);

                    // @ts-ignore the route should be valid, as we test it against a valid URL list
                    await navigateTo({
                        path: url.pathname,
                        query: parseQuery(url.search),
                        hash: url.hash,
                    });
                }
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
        async setSuperUserMode(superuser: boolean, job?: Job): Promise<void> {
            const { $grpc } = useNuxtApp();
            try {
                const req = {
                    superuser: superuser,
                } as SetSuperUserModeRequest;

                if (job) {
                    req.job = job!.name;
                }

                const call = $grpc.getAuthClient().setSuperUserMode(req);
                const { response } = await call;

                if (superuser) {
                    this.permissions.push('superuser');
                } else {
                    this.permissions = this.permissions.filter((p) => p !== 'superuser');
                }

                this.setAccessTokenExpiration(toDate(response.expires));
                this.setActiveChar(response.char!);
                this.setJobProps(response.jobProps);

                useNotificatorStore().add({
                    title: { key: 'notifications.superuser_menu.setsuperusermode.title', parameters: {} },
                    description: {
                        key: 'notifications.superuser_menu.setsuperusermode.content',
                        parameters: { job: job?.label ?? this.activeChar?.jobLabel ?? 'N/A' },
                    },
                    type: NotificationType.INFO,
                });

                await navigateTo({ name: 'overview' });
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

import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { JobProps } from '~~/gen/ts/resources/users/job_props';
import type { Job } from '~~/gen/ts/resources/users/jobs';
import type { User } from '~~/gen/ts/resources/users/users';
import type { SetSuperUserModeRequest } from '~~/gen/ts/services/auth/auth';

export const logger = useLogger('🔑 Auth');

export interface AuthState {
    accessTokenExpiration: null | Date;
    username: null | string;
    lastCharID: undefined | number;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: string[];
    jobProps: null | JobProps;
}

export const useAuthStore = defineStore('auth', {
    state: () =>
        ({
            // Persisted
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
        pick: ['accessTokenExpiration', 'lastCharID', 'username'],
    },
    actions: {
        loginStart(): void {
            this.loggingIn = true;
            this.loginError = null;
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
            if (char === null) {
                useGRPCWebsocketTransport().close();
            }
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
            useGRPCWebsocketTransport().close();
        },

        // GRPC Calls
        async doLogin(username: string, password: string): Promise<void> {
            // Start login
            this.loginStart();
            this.setActiveChar(null);
            this.setPermissions([]);

            try {
                const call = getGRPCAuthClient().login({ username: username, password: password });
                const { response } = await call;

                this.loginStop(null);

                this.username = username;

                if (response.char === undefined) {
                    logger.info('Login response without included char, redirecting to char selector');
                    this.setAccessTokenExpiration(toDate(response.expires));

                    const route = useRoute();

                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                } else {
                    logger.info('Received fast-tracked login response with char, id:', response.char.char?.userId);
                    this.setActiveChar(response.char.char ?? null);
                    this.setAccessTokenExpiration(toDate(response.char.expires));
                    this.setPermissions(response.char.permissions);
                    this.setJobProps(response.char.jobProps);

                    // @ts-expect-error route should be valid, as we test it against a valid URL list
                    const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
                    await navigateTo(target);
                }
            } catch (e) {
                this.loginStop((e as RpcError).message);
                this.setAccessTokenExpiration(null);
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async doLogout(): Promise<void> {
            this.loggingIn = true;

            try {
                await getGRPCAuthClient().logout({});
                this.clearAuthInfo();
            } catch (e) {
                handleGRPCError(e as RpcError);

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
            } finally {
                this.loginStop(null);
            }
        },
        async chooseCharacter(charId?: number, redirect?: boolean): Promise<void> {
            if (charId === undefined) {
                if (!this.lastCharID) {
                    const route = useRoute();

                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                    return;
                }

                charId = this.lastCharID;
            }

            try {
                const call = getGRPCAuthClient().chooseCharacter({ charId: charId });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                this.username = response.username;
                this.setActiveChar(response.char);
                this.setAccessTokenExpiration(toDate(response.expires));
                this.setPermissions(response.permissions);
                this.setJobProps(response.jobProps);

                if (redirect) {
                    const redirect = useRoute().query.redirect ?? useSettingsStore().startpage ?? '/overview';
                    const path = redirect || '/overview';
                    const url = new URL('https://example.com' + path);

                    // @ts-expect-error route should be valid, as we test it against a valid URL list
                    await navigateTo({
                        path: url.pathname,
                        query: parseQuery(url.search),
                        hash: url.hash,
                    });
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async setSuperUserMode(superuser: boolean, job?: Job): Promise<void> {
            try {
                const req = {
                    superuser: superuser,
                } as SetSuperUserModeRequest;

                if (job) {
                    req.job = job!.name;
                }

                const call = getGRPCAuthClient().setSuperUserMode(req);
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
                handleGRPCError(e as RpcError);
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
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

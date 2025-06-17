import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { webSocket } from '~/composables/grpc/grpcws/bridge';
import { useSettingsStore } from '~/stores/settings';
import type { JobProps } from '~~/gen/ts/resources/jobs/job_props';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';
import type { SetSuperuserModeRequest } from '~~/gen/ts/services/auth/auth';

const logger = useLogger('🔑 Auth');

export const useAuthStore = defineStore(
    'auth',
    () => {
        const { $grpc } = useNuxtApp();
        const notifications = useNotificationsStore();

        // State
        const accessTokenExpiration = ref<Date | null>(null);
        const username = ref<string | null>(null);
        const lastCharID = ref<number | undefined>(0);

        const activeChar = ref<User | null>(null);
        const loggingIn = ref<boolean>(false);
        const loginError = ref<RpcError | null>(null);
        const permissions = ref<string[]>([]);
        const jobProps = ref<JobProps | null>({
            job: '',
            livemapMarkerColor: '',
            radioFrequency: undefined,
            quickButtons: {
                penaltyCalculator: false,
                mathCalculator: false,
            },
            logoFileId: undefined,
            logoFile: undefined,
        });

        // Actions
        const loginStart = (): void => {
            loggingIn.value = true;
            loginError.value = null;
        };

        const loginStop = (error: RpcError | null): void => {
            loggingIn.value = false;
            loginError.value = error;
        };

        const setAccessTokenExpiration = (expiration: string | Date | null): void => {
            if (typeof expiration === 'string') {
                expiration = new Date(expiration);
            }
            accessTokenExpiration.value = expiration;
        };

        const setActiveChar = (char: User | null): void => {
            activeChar.value = char;
            lastCharID.value = char ? char.userId : lastCharID.value;
        };

        const setPermissions = (perms: string[]): void => {
            permissions.value.length = 0;
            permissions.value.push(...perms.sort());
        };

        const setJobProps = (jp: JobProps | undefined): void => {
            if (jp === undefined) {
                jobProps.value = null;
            } else {
                if (!jobProps.value) {
                    jobProps.value = jp;
                } else {
                    jobProps.value.job = jp.job;
                    jobProps.value.jobLabel = jp.jobLabel;
                    jobProps.value.livemapMarkerColor = jp.livemapMarkerColor;
                    jobProps.value.radioFrequency = jp.radioFrequency;
                    jobProps.value.quickButtons = jp.quickButtons;
                    jobProps.value.discordGuildId = jp.discordGuildId;
                    jobProps.value.logoFileId = jp.logoFileId;
                    jobProps.value.logoFile = jp.logoFile;
                }
            }
        };

        const clearAuthInfo = (): void => {
            setAccessTokenExpiration(null);
            setActiveChar(null);
            setPermissions([]);
            setJobProps(undefined);
            username.value = null;
            useGRPCWebsocketTransport().close();
        };

        // GRPC Calls
        const doLogin = async (user: string, pass: string): Promise<void> => {
            loginStart();
            setActiveChar(null);
            setPermissions([]);

            try {
                const call = $grpc.auth.auth.login({ username: user, password: pass });
                const { response } = await call;

                loginStop(null);

                username.value = user;

                if (response.char === undefined) {
                    logger.info('Login response without included char, redirecting to char selector');
                    setAccessTokenExpiration(toDate(response.expires));

                    const route = useRoute();
                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                } else {
                    logger.info('Received fast-tracked login response with char, id:', response.char.char?.userId);
                    setActiveChar(response.char.char ?? null);
                    setAccessTokenExpiration(toDate(response.char.expires));
                    setPermissions(response.char.permissions);
                    setJobProps(response.char.jobProps);

                    const startpage = useSettingsStore().startpage ?? '/overview';
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-ignore route should be valid, as we test it against a valid URL list
                    await navigateTo(startpage);
                }
            } catch (e) {
                const err = e as RpcError;
                loginStop(err);
                setAccessTokenExpiration(null);
                handleGRPCError(err);
            }
        };

        const doLogout = async (): Promise<void> => {
            loggingIn.value = true;

            try {
                await $grpc.auth.auth.logout({});
            } catch (e) {
                clearAuthInfo();
                handleGRPCError(e as RpcError);

                notifications.add({
                    title: { key: 'notifications.auth.error_logout.title', parameters: {} },
                    description: {
                        key: 'notifications.auth.error_logout.content',
                        parameters: { msg: (e as RpcError).message },
                    },
                    type: NotificationType.ERROR,
                });
            } finally {
                clearAuthInfo();
                loginStop(null);
            }
        };

        const chooseCharacter = async (charId?: number, redirect?: boolean): Promise<void> => {
            if (charId === undefined) {
                if (!lastCharID.value) {
                    const route = useRoute();

                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                    return;
                }

                charId = lastCharID.value;
            }

            try {
                const call = $grpc.auth.auth.chooseCharacter({
                    charId,
                });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                username.value = response.username;
                setActiveChar(response.char);
                setAccessTokenExpiration(toDate(response.expires));
                setPermissions(response.permissions);
                setJobProps(response.jobProps);

                if (redirect) {
                    const redirectQuery = useRoute().query.redirect;
                    const redirectPath =
                        (typeof redirectQuery === 'string' ? redirectQuery : redirectQuery?.join('/')) ??
                        useSettingsStore().startpage ??
                        '/overview';
                    const path = redirectPath || '/overview';
                    const url = new URL('https://example.com' + path);

                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-ignore route should be valid, as we test it against a valid URL list
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
        };

        const setSuperuserMode = async (superuser: boolean, job?: Job): Promise<void> => {
            try {
                const req: SetSuperuserModeRequest = {
                    superuser,
                };

                if (job) {
                    req.job = job.name;
                }

                const call = $grpc.auth.auth.setSuperuserMode(req);
                const { response } = await call;

                if (superuser) {
                    permissions.value.push('superuser');
                } else {
                    permissions.value = permissions.value.filter((p) => p !== 'superuser');
                }

                notifications.add({
                    title: { key: 'notifications.superuser_menu.setsuperusermode.title', parameters: {} },
                    description: {
                        key: 'notifications.superuser_menu.setsuperusermode.content',
                        parameters: {
                            job: job?.label ?? activeChar.value?.jobLabel ?? 'N/A',
                        },
                    },
                    type: NotificationType.INFO,
                });

                await navigateTo({ name: 'overview' });

                setAccessTokenExpiration(toDate(response.expires));
                setActiveChar(response.char!);
                setPermissions(response.permissions);
                setJobProps(response.jobProps);
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Getters
        const isSuperuser = computed<boolean>(() => {
            return permissions.value.includes('superuser');
        });

        // Watchers
        watch(username, (val) => {
            // Connect to the WebSocket if the user is logged in
            if (val !== null && val !== '') {
                if (webSocket.status.value !== 'OPEN' && webSocket.status.value !== 'CONNECTING') {
                    webSocket.open();
                }
            } else {
                webSocket.close();
            }
        });

        return {
            // State
            accessTokenExpiration,
            username,
            lastCharID,
            activeChar,
            loggingIn,
            loginError,
            permissions,
            jobProps,

            // Actions
            loginStart,
            loginStop,
            setAccessTokenExpiration,
            setActiveChar,
            setPermissions,
            setJobProps,
            clearAuthInfo,
            doLogin,
            doLogout,
            chooseCharacter,
            setSuperuserMode,

            // Getters
            isSuperuser,
        };
    },
    {
        persist: {
            pick: ['accessTokenExpiration', 'lastCharID', 'username'],
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            afterHydrate: (ctx: any) => {
                const store = ctx.store;
                if (typeof store.accessTokenExpiration === 'string') {
                    store.accessTokenExpiration = new Date(store.accessTokenExpiration);
                }
            },
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

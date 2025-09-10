import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { webSocket } from '~/composables/grpc/grpcws/bridge';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { JobProps } from '~~/gen/ts/resources/jobs/job_props';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';
import type { User } from '~~/gen/ts/resources/users/users';

const logger = useLogger('🔑 Auth');

export const useAuthStore = defineStore(
    'auth',
    () => {
        const settingsStore = useSettingsStore();
        const notifications = useNotificationsStore();

        // State
        const sessionExpiration = ref<Date | null>(null);
        const username = ref<string | null>(null);
        const lastCharID = ref<number | undefined>(0);

        const activeChar = ref<User | null>(null);
        const loggingIn = ref<boolean>(false);
        const loginError = ref<RpcError | null>(null);
        const permissions = ref<Permission[]>([]);
        const attributes = ref<RoleAttribute[]>([]);
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
        };

        const setActiveChar = (char: User | null): void => {
            activeChar.value = char;
            lastCharID.value = char ? char.userId : lastCharID.value;
        };

        const setPermissions = (perms: Permission[], attrs: RoleAttribute[]): void => {
            permissions.value.length = 0;
            permissions.value.push(...perms.sort());
            attributes.value.length = 0;
            attributes.value.push(...attrs.sort());
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
            logger.info('Clearing auth info');
            username.value = null;
            setActiveChar(null);
            setAccessTokenExpiration(null);
            setPermissions([], []);
            setJobProps(undefined);

            // Close the WebSocket connection when logging out
            useGRPCWebsocketTransport().close();
        };

        // GRPC Calls
        const doLogin = async (user: string, pass: string): Promise<void> => {
            loginStart();
            setActiveChar(null);
            setPermissions([], []);

            try {
                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.login({ username: user, password: pass });
                const { response } = await call;

                loginStop(null);

                username.value = user;

                if (response.char === undefined) {
                    logger.info('Login response (not fast-tracked), redirecting to char selector');
                    setAccessTokenExpiration(toDate(response.expires));

                    const route = useRoute();
                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                } else {
                    logger.info('Received fast-tracked login response with char, id:', response.char.char?.userId);

                    setAccessTokenExpiration(toDate(response.char.expires));
                    setActiveChar(response.char.char ?? null);
                    setPermissions(response.char.permissions, response.char.attributes);
                    setJobProps(response.char.jobProps);

                    const startpage = settingsStore.startpage ?? '/overview';
                    try {
                        await navigateTo(startpage);
                    } catch (_) {
                        logger.error('Failed to navigate to startpage, falling back to /overview');
                        await navigateTo('/overview');
                    }
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

            const authAuthClient = await getAuthAuthClient();

            try {
                await authAuthClient.logout({});
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
            if (charId === undefined || charId <= 0) {
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

            const authAuthClient = await getAuthAuthClient();

            try {
                const call = authAuthClient.chooseCharacter({
                    charId: charId,
                });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                username.value = response.username;
                setAccessTokenExpiration(toDate(response.expires));
                setActiveChar(response.char);
                setPermissions(response.permissions, response.attributes);
                setJobProps(response.jobProps);

                if (redirect) {
                    const redirectQuery = useRoute().query.redirect;
                    const redirectPath =
                        (typeof redirectQuery === 'string' ? redirectQuery : redirectQuery?.join('/')) ??
                        settingsStore.startpage ??
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
            const authAuthClient = await getAuthAuthClient();

            try {
                const call = authAuthClient.setSuperuserMode({
                    superuser: superuser,
                    job: job?.name,
                });
                const { response } = await call;

                await navigateTo({ name: 'overview' });

                if (superuser) {
                    permissions.value.push({
                        id: 0,
                        category: 'Superuser',
                        name: 'Superuser',
                        guardName: 'superuser-superuser',
                        val: true,
                    } as Permission);
                } else {
                    permissions.value = permissions.value.filter((p) => p.guardName === 'superuser-superuser');
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

                setActiveChar(response.char!);
                setPermissions(response.permissions, response.attributes);
                setJobProps(response.jobProps);
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Getters
        const isSuperuser = computed<boolean>(() => {
            return !!permissions.value.find((p) => p.guardName === 'superuser-superuser');
        });

        // Watchers
        watch(username, (val) => {
            // Connect to the WebSocket if the user is logged in
            if (val !== null && val !== '') {
                if (webSocket.status.value !== 'OPEN' && webSocket.status.value !== 'CONNECTING') {
                    logger.info('Username set, opening WebSocket connection, status:', webSocket.status.value);
                    webSocket.open();
                }
            } else {
                logger.info('Username cleared, closing WebSocket connection, status:', webSocket.status.value);
                webSocket.close();
            }
        });

        return {
            // State
            sessionExpiration,
            username,
            lastCharID,
            activeChar,
            loggingIn,
            loginError,
            permissions,
            attributes,
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
            pick: ['sessionExpiration', 'username', 'lastCharID'],
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            afterHydrate: (ctx: any) => {
                const store = ctx.store;
                if (typeof store.sessionExpiration === 'string') {
                    store.sessionExpiration = new Date(store.sessionExpiration);
                }
            },
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

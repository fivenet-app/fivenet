import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { JobProps } from '~~/gen/ts/resources/users/job_props';
import type { Job } from '~~/gen/ts/resources/users/jobs';
import type { User } from '~~/gen/ts/resources/users/users';
import type { SetSuperUserModeRequest } from '~~/gen/ts/services/auth/auth';

export const logger = useLogger('🔑 Auth');

export const useAuthStore = defineStore(
    'auth',
    () => {
        const { $grpc } = useNuxtApp();

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
            theme: 'defaultTheme',
            livemapMarkerColor: '',
            radioFrequency: undefined,
            quickButtons: {
                penaltyCalculator: false,
                mathCalculator: false,
                bodyCheckup: false,
            },
            logoUrl: undefined,
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
            if (char === null) {
                useGRPCWebsocketTransport().close();
            }
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
                    jobProps.value.theme = jp.theme;
                    jobProps.value.livemapMarkerColor = jp.livemapMarkerColor;
                    jobProps.value.radioFrequency = jp.radioFrequency;
                    jobProps.value.quickButtons = jp.quickButtons;
                    jobProps.value.discordGuildId = jp.discordGuildId;
                    jobProps.value.logoUrl = jp.logoUrl;
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

                    // @ts-expect-error route should be valid, as we test it against a valid URL list
                    const target = useRouter().resolve(useSettingsStore().startpage ?? '/overview');
                    await navigateTo(target);
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

                useNotificatorStore().add({
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
                    const redirectPath = useRoute().query.redirect ?? useSettingsStore().startpage ?? '/overview';
                    const path = redirectPath || '/overview';
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
        };

        const setSuperUserMode = async (superuser: boolean, job?: Job): Promise<void> => {
            try {
                const req: SetSuperUserModeRequest = {
                    superuser,
                };

                if (job) {
                    req.job = job.name;
                }

                const call = $grpc.auth.auth.setSuperUserMode(req);
                const { response } = await call;

                if (superuser) {
                    permissions.value.push('superuser');
                } else {
                    permissions.value = permissions.value.filter((p) => p !== 'superuser');
                }

                setAccessTokenExpiration(toDate(response.expires));
                setActiveChar(response.char!);
                setJobProps(response.jobProps);

                useNotificatorStore().add({
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
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Getters
        const isSuperuser = computed<boolean>(() => {
            return permissions.value.includes('superuser');
        });

        const getAccessTokenExpiration = computed<null | Date>(() => {
            if (typeof accessTokenExpiration.value === 'string') {
                accessTokenExpiration.value = new Date(Date.parse(accessTokenExpiration.value as unknown as string));
            }
            return accessTokenExpiration.value;
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
            setSuperUserMode,

            // Getters
            isSuperuser,
            getAccessTokenExpiration,
        };
    },
    {
        persist: {
            pick: ['accessTokenExpiration', 'lastCharID', 'username'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

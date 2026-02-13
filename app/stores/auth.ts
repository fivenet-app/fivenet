import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { webSocket } from '~/composables/grpcws/bridge';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { JobProps } from '~~/gen/ts/resources/jobs/props/props';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions/permissions';
import type { User } from '~~/gen/ts/resources/users/user';
import type { ImpersonateJobResponse } from '~~/gen/ts/services/auth/auth';

const logger = useLogger('🔑 Auth');

/**
 * Pinia store for managing user sessions, permissions, and state.
 */
export const useAuthStore = defineStore(
    'auth',
    () => {
        const settingsStore = useSettingsStore();
        const notifications = useNotificationsStore();
        const authSessionStore = useAuthSessionStore();

        // State
        /**
         * The username of the currently logged-in user.
         */
        const username = ref<string | null>(null);
        /**
         * Account ID.
         */
        const accountId = ref<number | null>(null);

        /**
         * The ID of the last selected character.
         */
        const lastCharID = ref<number | undefined>(0);
        /**
         * The currently active character.
         */
        const activeChar = ref<User | null>(null);
        /**
         * The job properties of the user.
         */
        const jobProps = ref<JobProps | null>({
            job: '',
            livemapMarkerColor: '',
            radioFrequency: undefined,
            quickButtons: {
                penaltyCalculator: false,
            },
            logoFileId: undefined,
            logoFile: undefined,
        });

        /**
         * Indicates whether a login operation is in progress.
         */
        const loggingIn = ref<boolean>(false);
        /**
         * Stores any error that occurred during login.
         */
        const loginError = ref<RpcError | null>(null);

        /**
         * The list of permissions assigned to the user.
         */
        const permissions = ref<Permission[]>([]);
        /**
         * The list of role attributes assigned to the user.
         */
        const attributes = ref<RoleAttribute[]>([]);

        /**
         * Set or unset the username.
         * @param val - The username of the user.
         */
        const setUsername = (val: string | null) => {
            // Connect to the WebSocket if the user is logged in
            if (val) {
                username.value = val;
                if (webSocket.status.value !== 'OPEN' && webSocket.status.value !== 'CONNECTING') {
                    logger.info('Username set, opening WebSocket connection, current status:', webSocket.status.value);
                    webSocket.open();
                }
            } else {
                username.value = null;
                logger.info('Username cleared, closing WebSocket connection, current status:', webSocket.status.value);
                webSocket.close();
            }
        };

        /**
         * Starts the login process by setting the loggingIn state to true and clearing any previous errors.
         */
        const loginStart = (): void => {
            loggingIn.value = true;
            loginError.value = null;
        };

        /**
         * Stops the login process and sets the provided error, if any.
         * @param error - The error that occurred during login, or null if successful.
         */
        const loginStop = (error: RpcError | null = null): void => {
            loggingIn.value = false;
            loginError.value = error;
        };

        /**
         * Updates the job properties of the user.
         * @param jp - The new job properties to set, or undefined to clear.
         */
        const setJobProps = (jp: JobProps | undefined): void => {
            if (!jp) {
                jobProps.value = null;
                return;
            }

            jobProps.value = {
                ...jobProps.value,
                ...jp,
            };
        };

        /**
         * Sets the currently active character and updates the last character ID.
         * @param char - The character to set as active, or null to clear.
         */
        const setActiveChar = (char: User | null = null): void => {
            lastCharID.value = char ? char.userId : lastCharID.value;
            activeChar.value = char;
        };

        /**
         * Updates the user's permissions and role attributes.
         * @param perms - The list of permissions to set.
         * @param attrs - The list of role attributes to set.
         */
        const setPermissions = (perms: Permission[], attrs: RoleAttribute[]): void => {
            permissions.value = [...perms.sort()];
            attributes.value = [...attrs.sort()];
        };

        /**
         * Logs in the user with the provided credentials.
         * @param user - The username of the user.
         * @param pass - The password of the user.
         */
        const doLogin = async (user: string, pass: string): Promise<void> => {
            // Prevent multiple simultaneous login attempts
            if (loggingIn.value) return;

            loginStart();
            setActiveChar();
            setPermissions([], []);

            try {
                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.login({ username: user, password: pass });
                const { response } = await call;
                refreshCookie('fivenet_authed');

                accountId.value = response.accountId;
                setUsername(user);
                loginStop();

                if (response.char === undefined) {
                    logger.info('Login response (not fast-tracked), redirecting to char selector');

                    const route = useRoute();
                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                } else {
                    logger.info('Received fast-tracked login response with char, id:', response.char.char?.userId);

                    setActiveChar(response.char.char ?? null);
                    setPermissions(response.char.permissions, response.char.attributes);
                    setJobProps(response.char.jobProps);
                    setUserToken(response.char.token);

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
                handleGRPCError(err);
            }
        };

        /**
         * Logout the currently logged-in user and clear authentication information.
         */
        const doLogout = async (): Promise<void> => {
            // User is about to logout, ignore ongoing logins/choose character actions
            loginStart();

            try {
                const authAuthClient = await getAuthAuthClient();

                await authAuthClient.logout({});

                refreshCookie('fivenet_authed');
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
                loginStop();
            }
        };

        /**
         * Selects a character for the user and optionally redirects to a specified page.
         * @param charId - The ID of the character to select. If undefined, the last character ID is used.
         * @param redirect - Whether to redirect the user after selecting the character.
         */
        const chooseCharacter = async (charId?: number, redirect?: boolean): Promise<void> => {
            // Prevent multiple simultaneous login attempts
            if (loggingIn.value) return;

            loginStart();

            if (charId === undefined || charId <= 0) {
                if (!lastCharID.value) {
                    if (!redirect) {
                        loginStop();
                        return;
                    }

                    const route = useRoute();
                    // Clear user token and unset login-block to avoid issues
                    setUserToken();
                    loginStop();
                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                    return;
                }

                charId = lastCharID.value;
            }

            try {
                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.chooseCharacter({
                    charId: charId,
                });
                const { response } = await call;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                setUsername(response.username);
                setUserToken(response.token);
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
            } finally {
                loginStop();
            }
        };

        /**
         * Impersonate job grade for the current user.
         * E.g., for testing permissions of a role.
         * @param grade - The job grade to impersonate.
         */
        const impersonateJob = async (grade: number): Promise<ImpersonateJobResponse> => {
            const authAuthClient = await getAuthAuthClient();

            const call = authAuthClient.impersonateJob({
                jobGrade: grade,
            });
            const { response } = await call;
            if (!response.char) {
                throw new Error('Server Error! No character in impersonate job response.');
            }

            setActiveChar(response.char);
            setPermissions(response.permissions, response.attributes);
            // Job props doesn't change on impersonation (user is part of the same job)
            setUserToken(response.token);

            return response;
        };

        /**
         * Sets the superuser mode for the user and updates permissions accordingly.
         * @param superuser - Whether to enable or disable superuser mode.
         * @param job - The job associated with the superuser mode, if any.
         */
        const setSuperuserMode = async (superuser: boolean, job?: Job): Promise<void> => {
            const authAuthClient = await getAuthAuthClient();

            try {
                const call = authAuthClient.setSuperuserMode({
                    superuser,
                    job: job?.name,
                });
                const { response } = await call;
                setUserToken(response.token, true);

                await navigateTo('/overview');

                // Update permissions based on superuser mode
                if (superuser) {
                    const superuserPermission: Permission = {
                        id: 0,
                        category: 'Superuser',
                        name: 'Superuser',
                        guardName: 'superuser-superuser',
                        val: true,
                    };
                    if (!permissions.value.some((p) => p.guardName === superuserPermission.guardName)) {
                        permissions.value.push(superuserPermission);
                    }
                } else {
                    permissions.value = permissions.value.filter((p) => p.guardName !== 'superuser-superuser');
                }

                // Notify user about the change
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

                // Update state with response data
                setActiveChar(response.char!);
                setPermissions(response.permissions, response.attributes);
                setJobProps(response.jobProps);
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Clears all authentication-related information from store and closes the WebSocket connection.
         */
        const clearAuthInfo = (): void => {
            logger.info('Clearing auth info');
            setUsername(null);
            setActiveChar(null);
            setPermissions([], []);
            setJobProps(undefined);
            setUserToken();

            // Close the WebSocket connection when logging out
            useGRPCWebsocketTransport().close();
        };

        /**
         * Set user token in session storage.
         * @param token - The user token to set. If undefined, the token is removed.
         */
        const setUserToken = async (token?: string, restartWS = false): Promise<void> => {
            if (!token) {
                authSessionStore.setUserToken(null);
                return;
            }

            const currentToken = authSessionStore.getUserToken();
            if (currentToken === token) {
                logger.debug('User token is the same as the current one, skipping update');
                return;
            }

            logger.debug('Setting user token in session storage');
            authSessionStore.setUserToken(token);

            if (activeChar.value !== null && restartWS) {
                logger.info('User token updated, send WebSocket re-auth message');
                useGRPCWebsocketTransport().updateUserToken(token);
            }
        };

        // Getters
        const isSuperuser = computed<boolean>(() => !!permissions.value.find((p) => p.guardName === 'superuser-superuser'));

        return {
            // State
            username,
            accountId,

            lastCharID,
            activeChar,
            jobProps,

            loggingIn,
            loginError,

            permissions,
            attributes,

            // Actions
            setUsername,
            loginStart,
            loginStop,
            doLogin,
            doLogout,

            clearAuthInfo,
            setUserToken,
            setActiveChar,
            setPermissions,
            setJobProps,

            chooseCharacter,
            impersonateJob,
            setSuperuserMode,

            // Getters
            isSuperuser,
        };
    },
    {
        persist: {
            pick: ['username', 'lastCharID'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

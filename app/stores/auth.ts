import type { RpcError, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { parseQuery } from 'vue-router';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { webSocket } from '~/composables/grpcws/bridge';
import { isSetupBypassRoute } from '~/composables/setup';
import { isUnauthenticatedError } from '~/utils/errors';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { JobProps } from '~~/gen/ts/resources/jobs/props/props';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions/permissions';
import type { User } from '~~/gen/ts/resources/users/user';
import type { ImpersonateJobResponse, RefreshAccountSessionResponse } from '~~/gen/ts/services/auth/auth';

const logger = useLogger('🔑 Auth');

const getAuthTabId = (): string => {
    if (typeof sessionStorage === 'undefined') return 'server';

    const key = 'fivenet-auth-tab-id';
    const existing = sessionStorage.getItem(key);
    if (existing) return existing;

    const id = typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : Math.random().toString(36);
    sessionStorage.setItem(key, id);
    return id;
};

export type AuthPhase =
    'anonymous' | 'bootstrapping' | 'account-ready' | 'selecting-character' | 'character-ready' | 'signing-out';

export type AuthEnsureResult =
    { kind: 'ready' } | { kind: 'needs-login' } | { kind: 'needs-character' } | { kind: 'temporary-failure'; error: unknown };

/**
 * Auth values that have been validated for use by data queries. These are kept
 * separate from the live session while an auth transition is in progress.
 */
export type AuthQueryContext = {
    accountId: number | null;
    characterId: number | undefined;
    job: string | undefined;
    grade: number | undefined;
    isSuperuser: boolean;
    canBeSuperuser: boolean;
    canBeConfigAdmin: boolean;
    revision: number;
};

/**
 * Pinia store for managing user sessions, permissions, and state.
 */
export const useAuthStore = defineStore(
    'auth',
    () => {
        const settingsStore = useSettingsStore();
        const notifications = useNotificationsStore();
        const authSessionStore = useAuthSessionStore();

        const characterAuthOptions = (charId?: number): RpcOptions | undefined => {
            const token = authSessionStore.getUserToken();
            const tokenAccountId = authSessionStore.userInfo?.accountId;
            const tokenUserId = authSessionStore.userInfo?.userId;
            if (!token || tokenAccountId === undefined || tokenUserId === undefined) return undefined;
            if (accountId.value === null || tokenAccountId !== accountId.value) return undefined;
            const expectedCharId = charId ?? activeChar.value?.userId;
            if (expectedCharId === undefined || tokenUserId !== expectedCharId) return undefined;

            return {
                meta: {
                    Authorization: `Bearer ${token}`,
                },
            };
        };

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
         * Whether the current account can enter superuser mode.
         */
        const canBeSuperuser = ref<boolean>(false);
        /**
         * Whether the current account can access config-admin setup without a character token.
         */
        const accountCanBeConfigAdmin = ref<boolean>(false);
        /** Whether the current account is currently in superuser (job admin mode) access. */
        const isSuperuser = computed<boolean>(() => !!permissions.value.find((p) => p.guardName === jobAdminPermGuard));
        /** Whether the current account can access config-admin gated screens and RPCs. */
        const canBeConfigAdmin = computed<boolean>(
            () => permissions.value.some((p) => p.guardName === configAdminPermGuard) || accountCanBeConfigAdmin.value,
        );
        /** The server-confirmed authentication state for this tab. */
        const phase = ref<AuthPhase>('anonymous');
        /** Non-HttpOnly marker maintained by the server for fast client-side session checks. */
        const authedCookie = useCookie<string | null>('fivenet_authed');
        /** Increments whenever credentials are cleared or replaced. */
        const sessionGeneration = ref<number>(0);
        const lastFailure = ref<{ kind: 'account-expired' | 'character-expired' | 'temporary'; at: number } | null>(null);
        /** Prevents query-backed views from fetching until route access was revalidated. */
        const isQueryTransitioning = ref(false);
        const queryContextRevision = ref(0);

        const buildQueryContext = (): AuthQueryContext => ({
            accountId: accountId.value,
            characterId: activeChar.value?.userId,
            job: activeChar.value?.job,
            grade: activeChar.value?.jobGrade,
            isSuperuser: isSuperuser.value,
            canBeSuperuser: canBeSuperuser.value,
            canBeConfigAdmin: canBeConfigAdmin.value,
            revision: queryContextRevision.value,
        });
        const queryContext = ref<AuthQueryContext>(buildQueryContext());

        const beginQueryTransition = (): void => {
            isQueryTransitioning.value = true;
        };

        /**
         * Publish live auth state to async-data keys. Call only after current
         * route access has been checked (and any redirect has completed).
         */
        const commitQueryContext = (): void => {
            queryContextRevision.value += 1;
            // Publish the complete new key while requests are still disabled.
            // Enabling first would start a request with the old key, which is
            // immediately cancelled when queryContext updates below.
            queryContext.value = buildQueryContext();
            isQueryTransitioning.value = false;
        };

        let chooseCharacterPromise: Promise<void> | undefined;
        let accountSessionPromise: Promise<RefreshAccountSessionResponse> | undefined;

        const broadcastSessionChange = (type: 'login' | 'logout' | 'changed'): void => {
            if (typeof BroadcastChannel === 'undefined') return;

            const channel = new BroadcastChannel('fivenet-auth');
            channel.postMessage({ type, source: getAuthTabId() });
            channel.close();
        };

        /**
         * Set or unset the username.
         * @param val - The username of the user.
         */
        const setUsername = (val: string | null, openSocket = true) => {
            // Connect to the WebSocket if the user is logged in
            if (val) {
                username.value = val;
                if (openSocket && webSocket.status.value !== 'OPEN' && webSocket.status.value !== 'CONNECTING') {
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
            permissions.value = [...perms].sort();
            attributes.value = [...attrs].sort();
            canBeSuperuser.value = perms.some((p) => p.guardName === superuserCanBePermGuard);
        };

        const setAccountCanBeConfigAdmin = (val: boolean): void => {
            accountCanBeConfigAdmin.value = val;
        };

        /**
         * Updates whether the current account can enter superuser mode.
         * @param val - Whether superuser mode is available.
         */
        const setCanBeSuperuser = (val: boolean): boolean => {
            const previousValue = canBeSuperuser.value;
            canBeSuperuser.value = val;
            if (!canBeSuperuser.value) {
                // Remove can be- and superuser permission if user can't be superuser anymore
                permissions.value = permissions.value.filter(
                    (p) => p.guardName !== jobAdminPermGuard && p.guardName !== superuserCanBePermGuard,
                );
            }

            return previousValue;
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
            clearAuthInfo();
            phase.value = 'bootstrapping';
            lastFailure.value = null;
            sessionGeneration.value += 1;
            const loginGeneration = sessionGeneration.value;

            try {
                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.login({ username: user, password: pass });
                const { response } = await call;
                if (sessionGeneration.value !== loginGeneration) return;

                beginQueryTransition();
                accountId.value = response.accountId;
                setAccountCanBeConfigAdmin(response.canBeConfigAdmin);
                authedCookie.value = 'true';
                setUsername(user);
                phase.value = 'account-ready';
                broadcastSessionChange('login');
                loginStop();

                if (response.char === undefined) {
                    logger.info('Login response (not fast-tracked), redirecting to char selector');
                    commitQueryContext();

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
                    phase.value = 'character-ready';
                    commitQueryContext();

                    const startpage = settingsStore.startpage ?? '/overview';
                    try {
                        await navigateTo(startpage);
                    } catch (_) {
                        logger.error('Failed to navigate to startpage, falling back to /overview');
                        await navigateTo('/overview');
                    }
                }
            } catch (e) {
                if (sessionGeneration.value !== loginGeneration) return;
                const err = e as RpcError;
                phase.value = 'anonymous';
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
            phase.value = 'signing-out';
            sessionGeneration.value += 1;

            try {
                // ChooseCharacter refreshes the account cookie. Let an already
                // running request finish before destroying that cookie, or its
                // response can recreate the session after logout succeeds.
                await chooseCharacterPromise?.catch(() => undefined);

                const authAuthClient = await getAuthAuthClient();

                await authAuthClient.logout({});
            } catch (e) {
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
                broadcastSessionChange('logout');
                loginStop();
            }
        };

        /**
         * Selects a character for the user and optionally redirects to a specified page.
         * @param charId - The ID of the character to select. If undefined, the last character ID is used.
         * @param redirect - Whether to redirect the user after selecting the character.
         */
        const chooseCharacterInternal = async (
            charId?: number,
            redirect?: boolean,
            hideError: boolean = false,
        ): Promise<void> => {
            loginStart();
            phase.value = 'selecting-character';
            const characterGeneration = sessionGeneration.value;

            if (charId === undefined || charId <= 0) {
                if (!lastCharID.value) {
                    if (!redirect) {
                        loginStop();
                        return;
                    }

                    const route = useRoute();
                    const appConfig = useAppConfig();

                    // Clear the user token because no character token can be restored.
                    setUserToken();

                    let refreshResp: RefreshAccountSessionResponse | null = null;
                    try {
                        refreshResp = await refreshAccountSession();
                        if (sessionGeneration.value !== characterGeneration) return;
                    } catch (e) {
                        if (isUnauthenticatedError(e)) {
                            invalidateSession();
                            loginStop();
                            return;
                        }

                        // Ignore refresh errors here; the selector fallback is still valid.
                    }

                    loginStop();

                    if (
                        appConfig.setupComplete === false &&
                        refreshResp?.canBeConfigAdmin &&
                        route.path !== '/settings/setup' &&
                        !isSetupBypassRoute(route.path)
                    ) {
                        await navigateTo({
                            path: '/settings/setup',
                            query: route.query,
                        });
                        return;
                    }

                    await navigateTo({
                        name: 'auth-character-selector',
                        query: route.query,
                    });
                    return;
                }

                charId = lastCharID.value;
            }

            try {
                if (accountId.value === null) {
                    try {
                        await refreshAccountSession();
                        if (sessionGeneration.value !== characterGeneration) return;
                    } catch (e) {
                        if (isUnauthenticatedError(e)) {
                            invalidateSession();
                            throw e;
                        }

                        // Ignore refresh errors here; chooseCharacter can still proceed with the current session state.
                    }
                }

                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.chooseCharacter(
                    {
                        charId: charId,
                    },
                    characterAuthOptions(charId),
                );
                const { response } = await call;
                if (sessionGeneration.value !== characterGeneration) return;
                if (!response.char) {
                    throw new Error('Server Error! No character in choose character response.');
                }

                beginQueryTransition();
                setUsername(response.username);
                setActiveChar(response.char);
                setUserToken(response.token);
                setPermissions(response.permissions, response.attributes);
                setAccountCanBeConfigAdmin(false);
                setJobProps(response.jobProps);
                phase.value = 'character-ready';
                commitQueryContext();

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
                if (!hideError) handleGRPCError(e as RpcError);
                if (phase.value === 'selecting-character')
                    phase.value = accountId.value === null ? 'anonymous' : 'account-ready';
                throw e;
            } finally {
                loginStop();
            }
        };

        const chooseCharacter = (charId?: number, redirect?: boolean, hideError: boolean = false): Promise<void> => {
            if (chooseCharacterPromise) return chooseCharacterPromise;

            const promise = chooseCharacterInternal(charId, redirect, hideError);
            chooseCharacterPromise = promise;
            void promise.then(
                () => {
                    if (chooseCharacterPromise === promise) chooseCharacterPromise = undefined;
                },
                () => {
                    if (chooseCharacterPromise === promise) chooseCharacterPromise = undefined;
                },
            );

            return promise;
        };

        /**
         * Impersonate job grade for the current user.
         * E.g., for testing permissions of a role.
         * @param grade - The job grade to impersonate.
         */
        const impersonateJob = async (grade: number): Promise<ImpersonateJobResponse> => {
            const authAuthClient = await getAuthAuthClient();

            const call = authAuthClient.impersonateJob(
                {
                    jobGrade: grade,
                },
                characterAuthOptions(),
            );
            const { response } = await call;
            if (!response.char) {
                throw new Error('Server Error! No character in impersonate job response.');
            }

            beginQueryTransition();
            setActiveChar(response.char);
            setPermissions(response.permissions, response.attributes);
            // Job props doesn't change on impersonation (user is part of the same job)
            setUserToken(response.token);
            commitQueryContext();

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
                const call = authAuthClient.setSuperuserMode(
                    {
                        superuser,
                        job: job?.name,
                    },
                    characterAuthOptions(),
                );
                const { response } = await call;
                // Update state with response data first so websocket reauth can pick up the active character.
                beginQueryTransition();
                setActiveChar(response.char!);
                setUserToken(response.token);
                setPermissions(response.permissions, response.attributes);
                setJobProps(response.jobProps);
                commitQueryContext();
                await navigateTo('/overview');

                // Notify user about the change
                if (superuser) {
                    notifications.add({
                        title: { key: 'notifications.superuser_menu.enabled_mode.title', parameters: {} },
                        description: {
                            key: 'notifications.superuser_menu.enabled_mode.content',
                            parameters: {
                                job: response.char?.jobLabel ?? job?.label ?? activeChar.value?.jobLabel ?? 'N/A',
                            },
                        },
                        type: NotificationType.INFO,
                    });
                } else {
                    notifications.add({
                        title: { key: 'notifications.superuser_menu.disabled_mode.title', parameters: {} },
                        description: {
                            key: 'notifications.superuser_menu.disabled_mode.content',
                            parameters: {},
                        },
                        type: NotificationType.INFO,
                    });
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /** Clears character-scoped state while retaining the account cookie session. */
        const clearCharacterSession = (kind?: 'character-expired' | 'temporary', batched = true): void => {
            if (batched) beginQueryTransition();
            setActiveChar(null);
            setPermissions([], []);
            setCanBeSuperuser(false);
            setAccountCanBeConfigAdmin(false);
            setJobProps(undefined);
            setUserToken();

            if (kind) lastFailure.value = { kind, at: Date.now() };
            phase.value = accountId.value === null ? 'anonymous' : 'account-ready';
            if (batched) commitQueryContext();
        };

        /**
         * Clears all authentication-related information from store and closes the WebSocket connection.
         */
        const clearAuthInfo = (): void => {
            if (
                phase.value === 'anonymous' &&
                username.value === null &&
                accountId.value === null &&
                activeChar.value === null &&
                authSessionStore.getUserToken() === null &&
                authedCookie.value !== 'true'
            ) {
                return;
            }

            logger.info('Clearing auth info');
            sessionGeneration.value += 1;
            beginQueryTransition();
            authedCookie.value = 'false';
            setUsername(null);
            accountId.value = null;
            clearCharacterSession(undefined, false);
            phase.value = 'anonymous';
            commitQueryContext();

            // Close the WebSocket connection when logging out
            useGRPCWebsocketTransport().close();
        };

        /** Clears server-rejected credentials while preserving the character preference. */
        const invalidateSession = (kind: 'account-expired' | 'character-expired' = 'account-expired'): void => {
            if (phase.value === 'anonymous' && accountId.value === null && username.value === null) return;

            sessionGeneration.value += 1;
            lastFailure.value = { kind, at: Date.now() };
            clearAuthInfo();
            broadcastSessionChange('changed');
        };

        /**
         * Set the in-memory character token for this tab.
         * @param token - The user token to set. If undefined, the token is removed.
         */
        const setUserToken = (token?: string): void => {
            if (!token) {
                authSessionStore.setUserToken(null);
                useGRPCWebsocketTransport().updateUserToken(null);
                return;
            }

            const currentToken = authSessionStore.getUserToken();
            if (currentToken === token) {
                logger.debug('User token is the same as the current one, skipping update');
                if (authSessionStore.userInfo?.accountId !== undefined) {
                    accountId.value = authSessionStore.userInfo.accountId;
                }
                return;
            }

            logger.debug('Setting in-memory user token');
            authSessionStore.setUserToken(token);
            if (authSessionStore.userInfo?.accountId !== undefined) {
                accountId.value = authSessionStore.userInfo.accountId;
            }

            logger.info('User token updated, send WebSocket re-auth message');
            useGRPCWebsocketTransport().updateUserToken(token);
        };

        const fetchAccountSession = async (): Promise<RefreshAccountSessionResponse> => {
            if (accountSessionPromise) return accountSessionPromise;

            const promise = (async (): Promise<RefreshAccountSessionResponse> => {
                const authAuthClient = await getAuthAuthClient();

                const call = authAuthClient.refreshAccountSession({});
                const { response } = await call;
                return response;
            })();

            accountSessionPromise = promise;
            void promise.then(
                () => {
                    if (accountSessionPromise === promise) accountSessionPromise = undefined;
                },
                () => {
                    if (accountSessionPromise === promise) accountSessionPromise = undefined;
                },
            );

            return promise;
        };

        const applyAccountSession = (response: RefreshAccountSessionResponse, openSocket: boolean): void => {
            beginQueryTransition();
            authedCookie.value = 'true';
            accountId.value = response.accountId;
            lastFailure.value = null;
            setAccountCanBeConfigAdmin(response.canBeConfigAdmin);
            if (response.username) {
                setUsername(response.username, openSocket);
            }
            if (activeChar.value === null) phase.value = 'account-ready';
            commitQueryContext();
        };

        const refreshAccountSession = async (): Promise<RefreshAccountSessionResponse | null> => {
            const generation = sessionGeneration.value;
            const response = await fetchAccountSession();
            if (sessionGeneration.value !== generation) return null;
            applyAccountSession(response, false);

            return response;
        };

        const restoreAccountSession = async (): Promise<RefreshAccountSessionResponse | null> => {
            const generation = sessionGeneration.value;
            const response = await fetchAccountSession();
            if (sessionGeneration.value !== generation) return null;
            applyAccountSession(response, true);

            return response;
        };

        const ensureAccountSession = async (force = false): Promise<AuthEnsureResult> => {
            if (!force && authedCookie.value === 'false') {
                if (phase.value !== 'anonymous' || accountId.value !== null || username.value !== null) {
                    clearAuthInfo();
                }
                return { kind: 'needs-login' };
            }

            // An unauthenticated refresh is a terminal result until the user
            // logs in. Without this guard, public-route middleware can retry
            // the same missing/expired cookie while redirecting to login.
            if (!force && phase.value === 'anonymous' && lastFailure.value?.kind === 'account-expired') {
                return { kind: 'needs-login' };
            }

            if (!force && (phase.value === 'account-ready' || phase.value === 'character-ready')) {
                return { kind: 'ready' };
            }

            phase.value = 'bootstrapping';
            try {
                await restoreAccountSession();
                return { kind: 'ready' };
            } catch (e) {
                if (isUnauthenticatedError(e)) {
                    invalidateSession('account-expired');
                    return { kind: 'needs-login' };
                }

                // The account cookie may still be valid; do not turn a
                // transport/server failure into a login redirect.
                phase.value = 'bootstrapping';
                lastFailure.value = { kind: 'temporary', at: Date.now() };
                return { kind: 'temporary-failure', error: e };
            }
        };

        const refreshCharacterSession = async (): Promise<AuthEnsureResult> => {
            const characterId = activeChar.value?.userId ?? lastCharID.value;
            if (!characterId) return { kind: 'needs-character' };

            try {
                await chooseCharacter(characterId, false, true);
                return activeChar.value !== null ? { kind: 'ready' } : { kind: 'needs-character' };
            } catch (e) {
                if (isUnauthenticatedError(e)) {
                    invalidateSession('account-expired');
                    return { kind: 'needs-login' };
                }

                clearCharacterSession('temporary');
                return { kind: 'temporary-failure', error: e };
            }
        };

        const ensureCharacterSession = async (): Promise<AuthEnsureResult> => {
            const account = await ensureAccountSession();
            if (account.kind !== 'ready') return account;
            if (phase.value === 'character-ready' && activeChar.value !== null) return { kind: 'ready' };
            if (!lastCharID.value) return { kind: 'needs-character' };

            try {
                await chooseCharacter(lastCharID.value, false, true);
                return activeChar.value !== null ? { kind: 'ready' } : { kind: 'needs-character' };
            } catch (e) {
                if (isUnauthenticatedError(e)) {
                    invalidateSession('account-expired');
                    return { kind: 'needs-login' };
                }

                clearCharacterSession('temporary');
                return { kind: 'temporary-failure', error: e };
            }
        };

        // Normal session changes may publish immediately. Auth event handlers
        // call beginQueryTransition first, which holds this snapshot stable
        // until they have revalidated the current route.
        watch(
            [
                accountId,
                activeChar,
                permissions,
                attributes,
                canBeSuperuser,
                accountCanBeConfigAdmin,
                isSuperuser,
                canBeConfigAdmin,
            ],
            () => {
                if (!isQueryTransitioning.value) commitQueryContext();
            },
            { flush: 'sync' },
        );

        return {
            username,
            accountId,

            lastCharID,
            activeChar,
            jobProps,

            loggingIn,
            loginError,

            permissions,
            attributes,
            canBeSuperuser,
            canBeConfigAdmin,
            phase,
            sessionGeneration,
            lastFailure,
            queryContext,
            isQueryTransitioning,

            loginStart,
            loginStop,
            doLogin,
            doLogout,

            clearAuthInfo,
            clearCharacterSession,
            invalidateSession,
            setUserToken,
            setCanBeSuperuser,
            setAccountCanBeConfigAdmin,
            setJobProps,
            beginQueryTransition,
            commitQueryContext,
            refreshAccountSession,
            restoreAccountSession,
            ensureAccountSession,
            ensureCharacterSession,
            refreshCharacterSession,

            chooseCharacter,
            impersonateJob,
            setSuperuserMode,

            isSuperuser,
        };
    },
    {
        persist: {
            pick: ['lastCharID'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}

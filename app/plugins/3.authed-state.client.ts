import { getAuthTabId, type AuthPhase } from '~/stores/auth';
import { useCentrumStore } from '~/stores/centrum';
import { useClipboardStore } from '~/stores/clipboard';
import { useHistoryStore } from '~/stores/history';
import { useLivemapStore } from '~/stores/livemap';
import { useMailerStore } from '~/stores/mailer';
import type { Notification } from '~/types/notifications';
import { createAuthContextRestartHandler } from '~/utils/authContextStream';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export type AuthBroadcastMessage = {
    type?: 'login' | 'logout' | 'changed';
    source?: string;
};

type AuthBroadcastStore = {
    ensureAccountSession: (force?: boolean) => Promise<unknown>;
    clearAuthInfo: () => void;
};

export async function handleAuthBroadcastMessage(
    message: AuthBroadcastMessage | undefined,
    tabId: string,
    authStore: AuthBroadcastStore,
    addNotification: (notification: Notification) => void,
): Promise<void> {
    if (message?.source && message.source === tabId) return;

    if (message?.type === 'login') {
        // A login can only be observed by checking the server-backed session;
        // never select a character from a broadcast hint.
        await authStore.ensureAccountSession(true);
        return;
    }

    // Logout/invalidated-session notifications already contain the
    // authoritative result. Clear local state without refreshing;
    // refreshing here would broadcast another `changed` event on failure and
    // create a request loop across tabs.
    authStore.clearAuthInfo();
    if (message?.type === 'logout') {
        addNotification({
            title: { key: 'notifications.auth.logged_out.title', parameters: {} },
            description: { key: 'notifications.auth.logged_out.content', parameters: {} },
            type: NotificationType.INFO,
        });
    }
}

export function getAuthStateRedirect(
    phase: AuthPhase,
    lastFailure: { kind: 'account-expired' | 'character-expired' | 'temporary'; at: number } | null,
    route: { path?: string; meta: { requiresAuth?: boolean; authTokenOnly?: boolean } },
): 'login' | 'character-selector' | undefined {
    if (route.path === '/auth/logout') return undefined;
    if (route.meta.requiresAuth === false) return undefined;
    if (phase === 'anonymous') return 'login';
    if (phase === 'account-ready' && lastFailure?.kind !== 'temporary' && !route.meta.authTokenOnly)
        return 'character-selector';

    return undefined;
}

export default defineNuxtPlugin({
    name: 'auth-state-coordinator',
    parallel: true,

    setup() {
        const authStore = useAuthStore();
        const clipboardStore = useClipboardStore();
        const historyStore = useHistoryStore();
        const centrumStore = useCentrumStore();
        const livemapStore = useLivemapStore();
        const mailerStore = useMailerStore();
        const notifications = useNotificationsStore();
        const auth = useAuth();
        const route = useRoute();
        let redirectPromise: Promise<unknown> | undefined;

        watch(
            () => authStore.accountId,
            (accountId) => {
                clipboardStore.setAccountScope(accountId);
                historyStore.setAccountScope(accountId);
                mailerStore.setAccountScope(accountId);
            },
            { immediate: true, flush: 'sync' },
        );

        // `keys.userState` is derived from the last route-validated auth
        // context. Restart shared streams only after that context is committed;
        // raw permission/character changes occur earlier during the transition.
        const restartForCommittedUserState = createAuthContextRestartHandler(auth.keys.userState.value, () => {
            void livemapStore.restartForAuthContext();
            void centrumStore.restartForAuthContext();
        });
        watch([() => auth.keys.userState.value, () => auth.isQueryTransitioning.value], restartForCommittedUserState, {
            flush: 'sync',
        });

        const redirectToLoginOnce = (): Promise<unknown> => {
            if (redirectPromise) return redirectPromise;

            redirectPromise = Promise.resolve(
                navigateTo({
                    name: 'auth-login',
                    query: { redirect: route.fullPath },
                    replace: true,
                }),
            ).finally(() => {
                redirectPromise = undefined;
            });

            return redirectPromise;
        };

        const redirectToCharacterSelectorOnce = (): Promise<unknown> => {
            if (redirectPromise) return redirectPromise;

            redirectPromise = Promise.resolve(
                navigateTo({
                    name: 'auth-character-selector',
                    query: { redirect: route.fullPath },
                    replace: true,
                }),
            ).finally(() => {
                redirectPromise = undefined;
            });

            return redirectPromise;
        };

        watch(
            () => authStore.phase as AuthPhase,
            async (phase) => {
                const redirect = getAuthStateRedirect(phase, authStore.lastFailure, route);
                if (redirect === 'login') {
                    await redirectToLoginOnce();
                    return;
                }

                if (redirect === 'character-selector') {
                    await redirectToCharacterSelectorOnce();
                }
            },
        );

        if (typeof BroadcastChannel !== 'undefined') {
            const channel = new BroadcastChannel('fivenet-auth');
            const tabId = getAuthTabId();
            channel.addEventListener('message', async (event: MessageEvent<AuthBroadcastMessage>) =>
                handleAuthBroadcastMessage(event.data, tabId, authStore, notifications.add),
            );

            onScopeDispose(() => channel.close());
        }
    },
});

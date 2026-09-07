import { describe, expect, it, vi } from 'vitest';
import { getAuthStateRedirect, handleAuthBroadcastMessage } from './3.authed-state.client';

const protectedCharacterRoute = {
    path: '/overview',
    meta: {
        requiresAuth: true,
        authTokenOnly: false,
    },
};

describe('auth state coordinator redirect decisions', () => {
    it('sends a confirmed anonymous state to login', () => {
        expect(getAuthStateRedirect('anonymous', null, protectedCharacterRoute)).toBe('login');
    });

    it('does not redirect a temporary restoration failure', () => {
        expect(
            getAuthStateRedirect('account-ready', { kind: 'temporary', at: Date.now() }, protectedCharacterRoute),
        ).toBeUndefined();
    });

    it('sends a character-required route to the selector when only the account is ready', () => {
        expect(getAuthStateRedirect('account-ready', null, protectedCharacterRoute)).toBe('character-selector');
    });

    it('does not redirect account-only routes to the character selector', () => {
        expect(
            getAuthStateRedirect('account-ready', null, {
                meta: {
                    requiresAuth: true,
                    authTokenOnly: true,
                },
            }),
        ).toBeUndefined();
    });

    it('does not redirect the logout page after auth is cleared', () => {
        expect(
            getAuthStateRedirect('anonymous', null, {
                path: '/auth/logout',
                meta: { requiresAuth: true, authTokenOnly: true },
            }),
        ).toBeUndefined();
    });
});

describe('auth state coordinator cross-tab messages', () => {
    it('ignores messages originating from this tab', async () => {
        const authStore = {
            ensureAccountSession: vi.fn(),
            clearAuthInfo: vi.fn(),
        };
        const addNotification = vi.fn();

        await handleAuthBroadcastMessage({ type: 'logout', source: 'tab-1' }, 'tab-1', authStore, addNotification);

        expect(authStore.clearAuthInfo).not.toHaveBeenCalled();
        expect(addNotification).not.toHaveBeenCalled();
    });

    it('refreshes once for a login from another tab', async () => {
        const authStore = {
            ensureAccountSession: vi.fn().mockResolvedValue({ kind: 'ready' }),
            clearAuthInfo: vi.fn(),
        };

        await handleAuthBroadcastMessage({ type: 'login', source: 'tab-2' }, 'tab-1', authStore, vi.fn());

        expect(authStore.ensureAccountSession).toHaveBeenCalledTimes(1);
        expect(authStore.ensureAccountSession).toHaveBeenCalledWith(true);
        expect(authStore.clearAuthInfo).not.toHaveBeenCalled();
    });

    it('clears on logout without refreshing and notifies the other tab', async () => {
        const authStore = {
            ensureAccountSession: vi.fn(),
            clearAuthInfo: vi.fn(),
        };
        const addNotification = vi.fn();

        await handleAuthBroadcastMessage({ type: 'logout', source: 'tab-2' }, 'tab-1', authStore, addNotification);

        expect(authStore.clearAuthInfo).toHaveBeenCalledTimes(1);
        expect(authStore.ensureAccountSession).not.toHaveBeenCalled();
        expect(addNotification).toHaveBeenCalledWith(
            expect.objectContaining({
                type: expect.anything(),
                title: { key: 'notifications.auth.logged_out.title', parameters: {} },
            }),
        );
    });

    it('clears an invalidated session without creating a refresh loop', async () => {
        const authStore = {
            ensureAccountSession: vi.fn(),
            clearAuthInfo: vi.fn(),
        };
        const addNotification = vi.fn();

        await handleAuthBroadcastMessage({ type: 'changed', source: 'tab-2' }, 'tab-1', authStore, addNotification);

        expect(authStore.clearAuthInfo).toHaveBeenCalledTimes(1);
        expect(authStore.ensureAccountSession).not.toHaveBeenCalled();
        expect(addNotification).not.toHaveBeenCalled();
    });
});

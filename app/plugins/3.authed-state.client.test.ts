import { describe, expect, it } from 'vitest';
import { getAuthStateRedirect } from './3.authed-state.client';

const protectedCharacterRoute = {
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
});

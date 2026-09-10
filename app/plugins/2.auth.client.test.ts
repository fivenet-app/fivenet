import { describe, expect, it } from 'vitest';
import { getCharacterSelectorRedirect, getLoginRedirect } from './2.auth.client';

const route = {
    fullPath: '/documents/?categories=[1]',
    query: { categories: '[1]' },
} as never;

describe('auth route redirects', () => {
    it('redirects account failures to login while preserving the target', () => {
        expect(getLoginRedirect(route)).toEqual({
            name: 'auth-login',
            query: { redirect: '/documents/?categories=[1]' },
            replace: true,
        });
    });

    it('redirects character failures to the selector while preserving the target', () => {
        expect(getCharacterSelectorRedirect(route)).toEqual({
            name: 'auth-character-selector',
            query: { redirect: '/documents/?categories=[1]' },
            replace: true,
        });
    });
});

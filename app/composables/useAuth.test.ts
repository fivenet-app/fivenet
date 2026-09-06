import { describe, expect, it } from 'vitest';
import { authKeys } from './useAuth';

describe('auth cache keys', () => {
    it('builds account and character keys', () => {
        expect(authKeys.account(123)).toBe('fivenet:account:123');
        expect(authKeys.character(123, 456)).toBe('fivenet:account:123:character:456');
    });

    it('represents missing identity values safely', () => {
        expect(authKeys.account(null)).toBe('fivenet:account:none');
        expect(authKeys.character(null, undefined)).toBe('fivenet:account:none:character:none');
    });

    it('changes capability keys when access state changes', () => {
        const regular = authKeys.capabilities(123, 456, false, true, false);
        const superuser = authKeys.capabilities(123, 456, true, true, false);
        const configAdmin = authKeys.capabilities(123, 456, false, true, true);

        expect(regular).not.toBe(superuser);
        expect(regular).not.toBe(configAdmin);
    });

    it('changes user-state keys when job or grade changes', () => {
        const police = authKeys.userState(123, 456, 'police', 1, false, true, false);
        const sheriff = authKeys.userState(123, 456, 'sheriff', 1, false, true, false);
        const seniorPolice = authKeys.userState(123, 456, 'police', 2, false, true, false);

        expect(police).not.toBe(sheriff);
        expect(police).not.toBe(seniorPolice);
    });
});

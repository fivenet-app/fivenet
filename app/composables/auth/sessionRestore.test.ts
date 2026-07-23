import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { restoreAuthTokenOnlySession } from './sessionRestore';

describe('restoreAuthTokenOnlySession', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('prefers character restoration when a usable last character exists', async () => {
        const chooseCharacter = vi.fn().mockResolvedValue(undefined);
        const restoreAccountSession = vi.fn().mockResolvedValue(undefined);

        await restoreAuthTokenOnlySession({
            activeChar: { value: null },
            lastCharID: 42,
            chooseCharacter,
            restoreAccountSession,
        });

        expect(chooseCharacter).toHaveBeenCalledWith(42, false);
        expect(restoreAccountSession).not.toHaveBeenCalled();
    });

    it('falls back to restoring the account session when lastCharID is empty', async () => {
        const chooseCharacter = vi.fn().mockResolvedValue(undefined);
        const restoreAccountSession = vi.fn().mockResolvedValue(undefined);

        await restoreAuthTokenOnlySession({
            activeChar: { value: null },
            lastCharID: 0,
            chooseCharacter,
            restoreAccountSession,
        });

        expect(chooseCharacter).not.toHaveBeenCalled();
        expect(restoreAccountSession).toHaveBeenCalledTimes(1);
    });

    it('falls back to restoring the account session when character restoration fails', async () => {
        const error = new Error('invalid character');
        const chooseCharacter = vi.fn().mockRejectedValue(error);
        const restoreAccountSession = vi.fn().mockResolvedValue(undefined);
        const warn = vi.spyOn(console, 'warn').mockImplementation(() => undefined);

        await restoreAuthTokenOnlySession({
            activeChar: { value: null },
            lastCharID: 42,
            chooseCharacter,
            restoreAccountSession,
        });

        expect(chooseCharacter).toHaveBeenCalledWith(42, false);
        expect(warn).toHaveBeenCalledWith(
            'Failed to restore last selected character, falling back to account session refresh.',
            error,
        );
        expect(restoreAccountSession).toHaveBeenCalledTimes(1);
    });
});

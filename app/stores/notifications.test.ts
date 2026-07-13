import { describe, expect, it, vi } from 'vitest';
import type { AccountGroupsChanged } from '~~/gen/ts/resources/userinfo/userinfo';
import { handleAccountGroupsChangedEvent, shouldRestartNotificationStream } from './notifications';

describe('handleAccountGroupsChangedEvent', () => {
    function createAuthStore(
        overrides: Partial<{
            isSuperuser: boolean;
            canBeConfigAdmin: boolean;
        }> = {},
    ) {
        return {
            isSuperuser: overrides.isSuperuser ?? false,
            canBeConfigAdmin: overrides.canBeConfigAdmin ?? false,
            setCanBeSuperuser: vi.fn().mockReturnValue(overrides.isSuperuser ?? false),
            setAccountCanBeConfigAdmin: vi.fn(),
            chooseCharacter: vi.fn().mockResolvedValue(undefined),
        };
    }

    function createEvent(overrides: Partial<AccountGroupsChanged>): AccountGroupsChanged {
        return {
            accountId: 1,
            changedAt: undefined as never,
            canBeSuperuser: false,
            canBeConfigAdmin: false,
            newGroups: undefined,
            ...overrides,
        };
    }

    it('refreshes when config-admin is granted without a superuser change', async () => {
        const authStore = createAuthStore({
            isSuperuser: false,
            canBeConfigAdmin: false,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: false,
                canBeConfigAdmin: true,
            }),
            authStore,
            { accountOnly: false },
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(false);
        expect(authStore.chooseCharacter).toHaveBeenCalledTimes(1);
        expect(authStore.chooseCharacter).toHaveBeenCalledWith(undefined, false);
    });

    it('refreshes when config-admin is revoked without a superuser change', async () => {
        const authStore = createAuthStore({
            isSuperuser: true,
            canBeConfigAdmin: true,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: true,
                canBeConfigAdmin: false,
            }),
            authStore,
            { accountOnly: false },
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(true);
        expect(authStore.chooseCharacter).toHaveBeenCalledTimes(1);
        expect(authStore.chooseCharacter).toHaveBeenCalledWith(undefined, false);
    });

    it('refreshes when superuser capability is revoked', async () => {
        const authStore = createAuthStore({
            isSuperuser: true,
            canBeConfigAdmin: true,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: false,
                canBeConfigAdmin: true,
            }),
            authStore,
            { accountOnly: false },
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(false);
        expect(authStore.chooseCharacter).toHaveBeenCalledTimes(1);
        expect(authStore.chooseCharacter).toHaveBeenCalledWith(undefined, false);
    });

    it('does not refresh when neither capability changes', async () => {
        const authStore = createAuthStore({
            isSuperuser: false,
            canBeConfigAdmin: false,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: false,
                canBeConfigAdmin: false,
            }),
            authStore,
            { accountOnly: false },
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(false);
        expect(authStore.chooseCharacter).not.toHaveBeenCalled();
    });

    it('updates account-level capabilities without reselecting a character in account-only scope', async () => {
        const authStore = createAuthStore({
            isSuperuser: false,
            canBeConfigAdmin: false,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: true,
                canBeConfigAdmin: true,
            }),
            authStore,
            { accountOnly: true },
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(true);
        expect(authStore.setAccountCanBeConfigAdmin).toHaveBeenCalledWith(true);
        expect(authStore.chooseCharacter).not.toHaveBeenCalled();
    });

    it('clears account-level config-admin state before refreshing character-scoped revocations', async () => {
        const authStore = createAuthStore({
            isSuperuser: false,
            canBeConfigAdmin: true,
        });

        await handleAccountGroupsChangedEvent(
            createEvent({
                canBeSuperuser: false,
                canBeConfigAdmin: false,
            }),
            authStore,
            { accountOnly: false },
        );

        expect(authStore.setAccountCanBeConfigAdmin).toHaveBeenCalledWith(false);
        expect(authStore.chooseCharacter).toHaveBeenCalledWith(undefined, false);
    });

    it('does not restart when the live stream controller has already changed', () => {
        const staleAbort = new AbortController();
        const liveAbort = new AbortController();

        staleAbort.abort();

        expect(shouldRestartNotificationStream(liveAbort, staleAbort)).toBe(false);
    });
});

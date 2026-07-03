import { describe, expect, it, vi } from 'vitest';
import type { AccountGroupsChanged } from '~~/gen/ts/resources/userinfo/userinfo';
import { handleAccountGroupsChangedEvent } from './notifications';

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
        );

        expect(authStore.setCanBeSuperuser).toHaveBeenCalledWith(false);
        expect(authStore.chooseCharacter).not.toHaveBeenCalled();
    });
});

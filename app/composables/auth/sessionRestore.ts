type AuthTokenOnlyRestoreStore = {
    activeChar: { value: unknown | null };
    lastCharID: number | undefined;
    chooseCharacter: (charId?: number, redirect?: boolean) => Promise<unknown>;
    restoreAccountSession: () => Promise<unknown>;
};

/**
 * Restores an auth-token-only session by preferring a character restore when possible.
 * Falls back to account-session refresh only when there is no usable last character.
 */
export async function restoreAuthTokenOnlySession(authStore: AuthTokenOnlyRestoreStore): Promise<void> {
    if (authStore.activeChar.value !== null) return;

    if (authStore.lastCharID !== undefined && authStore.lastCharID > 0) {
        try {
            await authStore.chooseCharacter(authStore.lastCharID, false);
            return;
        } catch (e) {
            console.warn('Failed to restore last selected character, falling back to account session refresh.', e);
        }
    }

    try {
        await authStore.restoreAccountSession();
    } catch (_) {
        // Keep the existing behavior: a failed account-session restore should not block route entry.
    }
}

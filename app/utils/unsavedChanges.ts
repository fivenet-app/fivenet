const confirmationPromises = new Map<symbol, Promise<boolean>>();

export const DEFAULT_UNSAVED_CHANGES_CONFIRMATION_GROUP = Symbol('utils-unsaved-changes');

export function runSharedUnsavedChangesConfirmation(
    confirmationGroup: symbol,
    createConfirmation: () => Promise<boolean>,
): Promise<boolean> {
    const existingConfirmation = confirmationPromises.get(confirmationGroup);
    if (existingConfirmation) return existingConfirmation;

    const confirmation = (async () => {
        try {
            return await createConfirmation();
        } finally {
            confirmationPromises.delete(confirmationGroup);
        }
    })();

    confirmationPromises.set(confirmationGroup, confirmation);

    return confirmation;
}

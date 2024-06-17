export function getErrorMessage(message: string): string {
    if (isTranslatedError(message)) {
        const errSplits = message.split(';');
        if (errSplits.length > 1) {
            return errSplits[1] ?? message;
        }
    }
    return message;
}

export function isTranslatedError(message: string): boolean {
    return message.startsWith('errors.');
}

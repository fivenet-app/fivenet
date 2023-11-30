export function getErrorMessage(message: string): string {
    if (message.startsWith('errors.')) {
        const errSplits = message.split(';');
        if (errSplits.length > 1) {
            return errSplits[1];
        }
    }
    return message;
}

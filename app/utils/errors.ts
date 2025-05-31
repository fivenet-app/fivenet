import type { TranslateItem } from '~/types/i18n';
import type { Error as CommonError } from '~~/gen/ts/resources/common/error';

export function getErrorMessage(err: RpcError): TranslateItem {
    if (isTranslatedError(err.message)) {
        const parsed = parseError(err.message);
        if (parsed?.content) {
            return parsed.content;
        }
    }

    return { key: err.message, parameters: {} };
}

export function parseError(err: RpcError): CommonError | undefined {
    try {
        return JSON.parse(err.message) as CommonError;
    } catch (_) {
        return undefined;
    }
}

export function isTranslatedError(message: string): boolean {
    return message.startsWith('{');
}

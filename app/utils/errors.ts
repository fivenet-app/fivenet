import type { Error as CommonError } from '~~/gen/ts/resources/common/error';

export function getErrorMessage(err: RpcError): TranslateItem {
    if (isTranslatedError(err.message)) {
        const parsed = JSON.parse(err.message) as CommonError;
        if (parsed.content) {
            return parsed.content;
        }
    }
    return { key: err.message, parameters: {} };
}

export function isTranslatedError(message: string): boolean {
    return message.startsWith('{');
}

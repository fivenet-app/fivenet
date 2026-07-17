import type { Error as CommonError } from '~~/gen/ts/resources/common/error';

export function getErrorMessage(err: RpcError): I18NItem {
    if (isTranslatedError(err.message)) {
        const parsed = parseErrorMessage(err.message);
        if (parsed?.content) {
            return parsed.content;
        }
    }

    return { key: err.message, parameters: {} };
}

export function parseError(err: RpcError): CommonError | undefined {
    return parseErrorMessage(err.message);
}

export function parseErrorMessage(message: string): CommonError | undefined {
    try {
        return JSON.parse(message) as CommonError;
    } catch (_) {
        return undefined;
    }
}

export function isTranslatedError(message: string): boolean {
    return message.trimStart().startsWith('{');
}

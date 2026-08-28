import type { Error as CommonError } from '~~/gen/ts/resources/common/error';
import type { I18NItem } from '~~/gen/ts/resources/common/i18n';

const templateKindTranslationKeys: Record<string, string> = {
    users: 'errors.documents.DocumentsService.TemplateKinds.users',
    documents: 'errors.documents.DocumentsService.TemplateKinds.documents',
    vehicles: 'errors.documents.DocumentsService.TemplateKinds.vehicles',
};

export function getErrorMessage(err: RpcError): I18NItem {
    if (isTranslatedError(err.message)) {
        const parsed = parseErrorMessage(err.message);
        if (parsed?.content) {
            return parsed.content;
        }
    }

    return { key: err.message, parameters: {} };
}

export function parseError(err: Error): CommonError | undefined {
    return parseErrorMessage(err.message);
}

export function parseErrorMessage(message: string): CommonError | undefined {
    try {
        return JSON.parse(message) as CommonError;
    } catch (_) {
        return undefined;
    }
}

export function localizeTemplateErrorParameters(error: CommonError, translate: (key: string) => string): CommonError {
    const localizeItem = (item: I18NItem | undefined): I18NItem | undefined => {
        if (!item) {
            return undefined;
        }

        const kind = item.parameters.kind;
        const translationKey = kind ? templateKindTranslationKeys[kind] : undefined;
        if (!translationKey) {
            return { ...item, parameters: { ...item.parameters } };
        }

        return {
            ...item,
            parameters: {
                ...item.parameters,
                kind: translate(translationKey),
            },
        };
    };

    return {
        ...error,
        title: localizeItem(error.title),
        content: localizeItem(error.content),
    };
}

export function isTranslatedError(message: string): boolean {
    return message.trimStart().startsWith('{');
}

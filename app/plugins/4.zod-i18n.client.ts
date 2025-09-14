// plugins/zod-i18n.ts
import { defineNuxtPlugin, useNuxtApp } from '#app';
import { z, ZodIssueCode, type ZodErrorMap } from 'zod';

export default defineNuxtPlugin(() => {
    const { $i18n } = useNuxtApp(); // global Composer from @nuxtjs/i18n
    // NOTE: don't destructure t; call $i18n.t(...) directly so it's always bound

    const errorMap: ZodErrorMap = (issue, ctx) => {
        const path = issue.path?.length ? issue.path.join('.') : undefined;

        // a tiny helper so we don’t repeat ourselves
        const t = (key: string, params?: Record<string, any>) =>
            // @nuxtjs/i18n uses Vue I18n under the hood; params become named values
            // Fallback to ctx.defaultError if key is missing
            ($i18n?.t(key, params) as string) || ctx.defaultError;

        switch (issue.code) {
            case ZodIssueCode.invalid_type:
                return {
                    message: t('zod.invalid_type', {
                        path,
                        expected: issue.expected,
                        received: issue.received,
                    }),
                };
            case ZodIssueCode.invalid_literal:
                return {
                    message: t('zod.invalid_literal', { path, expected: String(issue.expected) }),
                };
            case ZodIssueCode.unrecognized_keys:
                return {
                    message: t('zod.unrecognized_keys', { path, keys: issue.keys.join(', ') }),
                };
            case ZodIssueCode.invalid_union:
                return { message: t('zod.invalid_union', { path }) };
            case ZodIssueCode.invalid_union_discriminator:
                return {
                    message: t('zod.invalid_union_discriminator', {
                        path,
                        options: issue.options?.join(', '),
                    }),
                };
            case ZodIssueCode.invalid_enum_value:
                return {
                    message: t('zod.invalid_enum_value', {
                        path,
                        options: issue.options.join(', '),
                        received: String(issue.received),
                    }),
                };
            case ZodIssueCode.invalid_arguments:
                return { message: t('zod.invalid_arguments', { path }) };
            case ZodIssueCode.invalid_return_type:
                return { message: t('zod.invalid_return_type', { path }) };
            case ZodIssueCode.invalid_date:
                return { message: t('zod.invalid_date', { path }) };
            case ZodIssueCode.invalid_string:
                return {
                    message: t('zod.invalid_string', {
                        path,
                        validation: typeof issue.validation === 'string' ? issue.validation : 'invalid',
                    }),
                };
            case ZodIssueCode.too_small:
                return {
                    message: t('zod.too_small', {
                        path,
                        type: issue.type, // string | number | array | set | date
                        minimum: issue.minimum,
                        inclusive: issue.inclusive,
                        exact: issue.exact,
                    }),
                };
            case ZodIssueCode.too_big:
                return {
                    message: t('zod.too_big', {
                        path,
                        type: issue.type,
                        maximum: issue.maximum,
                        inclusive: issue.inclusive,
                        exact: issue.exact,
                    }),
                };
            case ZodIssueCode.custom:
                return { message: t('zod.custom', { path, message: issue.message ?? '' }) };
            case ZodIssueCode.invalid_intersection_types:
                return { message: t('zod.invalid_intersection_types', { path }) };
            case ZodIssueCode.not_multiple_of:
                return { message: t('zod.not_multiple_of', { path, multipleOf: issue.multipleOf }) };
            case ZodIssueCode.not_finite:
                return { message: t('zod.not_finite', { path }) };
            default:
                return { message: ctx.defaultError };
        }
    };

    // Set globally once per request/app instance
    z.setErrorMap(errorMap);
});

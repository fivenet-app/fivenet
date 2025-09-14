// plugins/zod-i18n.ts
import { defineNuxtPlugin, useNuxtApp } from '#app';
import { z } from 'zod/v4';

export default defineNuxtPlugin(() => {
    const nuxtApp = useNuxtApp();

    // --- utils ---
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const t = (key: string, params?: Record<string, any>) => (nuxtApp.$i18n?.t(key, params ?? {}) as unknown as string) ?? '';

    const hasKey = (key: string) => {
        const s = t(key);
        // vue-i18n returns the key itself when missing (unless you changed missing handler)
        return s && s !== key;
    };

    const joinValues = (values: unknown[], sep = ', ') => values.map((v) => (typeof v === 'string' ? v : String(v))).join(sep);

    const stringifyPrimitive = (v: unknown) => (typeof v === 'string' ? JSON.stringify(v) : String(v));

    const parsedType = (input: unknown): string => {
        if (input === null) return 'null';
        if (Array.isArray(input)) return 'array';
        return typeof input;
    };

    const getOriginLabel = (origin?: string | null) => {
        if (!origin) return t('zod.value') || 'value';
        const k = `zod.origins.${origin}`;
        return hasKey(k) ? t(k) : origin;
    };

    // i18n-backed Sizable: zod.sizable.<origin>.{unit,verb}
    const getSizing = (origin?: string | null): { unit: string; verb: string } | null => {
        if (!origin) return null;
        const unitKey = `zod.sizable.${origin}.unit`;
        const verbKey = `zod.sizable.${origin}.verb`;
        if (hasKey(unitKey) || hasKey(verbKey)) {
            return {
                unit: hasKey(unitKey) ? t(unitKey) : '',
                verb: hasKey(verbKey) ? t(verbKey) : '',
            };
        }
        return null;
    };

    // i18n-backed Nouns: zod.nouns.<format>
    const getNoun = (format: string) => {
        const k = `zod.nouns.${format}`;
        return hasKey(k) ? t(k) : format;
    };

    const getComparator = (dir: 'max' | 'min', inclusive: boolean) => {
        // keys you can translate to “≤ / ≥” OR “höchstens / mindestens”, etc.
        const key =
            dir === 'max'
                ? inclusive
                    ? 'zod.comparators.lte'
                    : 'zod.comparators.lt'
                : inclusive
                  ? 'zod.comparators.gte'
                  : 'zod.comparators.gt';
        const fallback = dir === 'max' ? (inclusive ? '<=' : '<') : inclusive ? '>=' : '>';
        return hasKey(key) ? t(key) : fallback;
    };

    // --- Zod v4 mapper that uses i18n nouns + sizable ---
    const customError: z.core.$ZodErrorMap = (issue): string => {
        switch (issue.code) {
            case 'invalid_type': {
                const received = parsedType(issue.input);
                if (received === 'undefined') return t('zod.required') || 'Required';
                return (
                    t('zod.invalid_type', {
                        expected: issue.expected,
                        received,
                    }) || 'Invalid input'
                );
            }

            case 'invalid_value': {
                const values = issue.values as unknown[] | undefined;
                if (values?.length === 1) {
                    return t('zod.invalid_value_single', { expected: stringifyPrimitive(values[0]) }) || 'Invalid input';
                }
                if (values && values.length > 1) {
                    return t('zod.invalid_value_one_of', { expected: joinValues(values, '|') }) || 'Invalid input';
                }
                return t('zod.invalid_value') || 'Invalid value';
            }

            case 'too_big': {
                const comp = getComparator('max', issue.inclusive ?? false);
                const originLabel = getOriginLabel(issue.origin);
                const sizing = getSizing(issue.origin);
                if (sizing) {
                    // e.g. DE: "{origin} darf {verb} {maximum} {unit} {comparator} enthalten"
                    return (
                        t('zod.too_big.have_full', {
                            origin: originLabel,
                            comparator: comp,
                            maximum: String(issue.maximum),
                            unit: sizing.unit,
                            verb: sizing.verb,
                        }) || 'Too big'
                    );
                }
                return (
                    t('zod.too_big.be_full', {
                        origin: originLabel,
                        comparator: comp,
                        maximum: String(issue.maximum),
                    }) || 'Too big'
                );
            }

            case 'too_small': {
                const comp = getComparator('min', issue.inclusive ?? false);
                const originLabel = getOriginLabel(issue.origin);
                const sizing = getSizing(issue.origin);
                if (sizing) {
                    // e.g. DE: "{origin} muss {verb} {comparator} {minimum} {unit} enthalten"
                    return (
                        t('zod.too_small.have_full', {
                            origin: originLabel,
                            comparator: comp,
                            minimum: String(issue.minimum),
                            unit: sizing.unit,
                            verb: sizing.verb,
                        }) || 'Too small'
                    );
                }
                // e.g. DE: "{origin} muss {comparator} {minimum} sein"
                return (
                    t('zod.too_small.be_full', {
                        origin: originLabel,
                        comparator: comp,
                        minimum: String(issue.minimum),
                    }) || 'Too small'
                );
            }

            case 'invalid_format': {
                const f = issue as z.core.$ZodStringFormatIssues;
                if (f.format === 'starts_with') return t('zod.starts_with', { prefix: f.prefix });
                if (f.format === 'ends_with') return t('zod.ends_with', { suffix: f.suffix });
                if (f.format === 'includes') return t('zod.includes', { includes: f.includes });
                if (f.format === 'regex') return t('zod.regex', { pattern: f.pattern });

                const noun = getNoun(f.format);
                return t('zod.invalid_format_noun', { noun }) || `Invalid ${noun}`;
            }

            case 'not_multiple_of':
                return t('zod.not_multiple_of', { divisor: issue.divisor }) || 'Invalid number';

            case 'unrecognized_keys': {
                const keys = (issue as z.core.$ZodIssueUnrecognizedKeys).keys;
                const count = keys.length;
                const msgKey = count === 1 ? 'zod.unrecognized_key' : 'zod.unrecognized_keys';
                return (
                    t(msgKey, { keys: joinValues(keys, ', '), count }) ||
                    (count === 1 ? 'Unrecognized key' : 'Unrecognized keys')
                );
            }

            case 'invalid_key':
                return t('zod.invalid_key', { origin: issue.origin }) || 'Invalid key';

            case 'invalid_union':
                return t('zod.invalid_union') || 'Invalid input';

            case 'invalid_element':
                return t('zod.invalid_element', { origin: issue.origin }) || 'Invalid element';

            case 'custom':
                return t('zod.custom', { message: (issue as z.core.$ZodIssueCustom).message ?? '' }) || 'Invalid input';

            default:
                return t('zod.invalid_input') || 'Invalid input';
        }
    };

    // init + keep in sync with language changes
    z.config({
        customError: customError,
    });
    nuxtApp.hook('i18n:localeSwitched', () => {
        z.config({
            customError: customError,
        });
    });
});

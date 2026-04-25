import { describe, expect, it } from 'vitest';
import { z } from 'zod';
import { secondsToDuration } from './duration';
import { zodDurationMinMaxPair, zodProtoDurationSchema } from './validation';

function getFirstIssueI18n(result: { error: { issues: unknown[] } }): string | undefined {
    const issue = result.error.issues[0] as { params?: { i18n?: string } } | undefined;
    return issue?.params?.i18n;
}

describe('zodProtoDurationSchema', () => {
    it('accepts a valid duration', () => {
        const schema = zodProtoDurationSchema({ required: true });
        const result = schema.safeParse(secondsToDuration(3600));

        expect(result.success).toBe(true);
    });

    it('rejects invalid duration shape', () => {
        const schema = zodProtoDurationSchema({ required: true });
        const result = schema.safeParse({ nope: true });

        expect(result.success).toBe(false);
        if (!result.success) {
            expect(getFirstIssueI18n(result)).toBe('zod.custom.duration.invalid');
        }
    });

    it('rejects missing required value', () => {
        const schema = zodProtoDurationSchema({ required: true });
        const result = schema.safeParse(undefined);

        expect(result.success).toBe(false);
        if (!result.success) {
            expect(getFirstIssueI18n(result)).toBe('zod.custom.duration.required');
        }
    });

    it('rejects value below minimum', () => {
        const schema = zodProtoDurationSchema({
            required: true,
            min: secondsToDuration(60),
        });
        const result = schema.safeParse(secondsToDuration(30));

        expect(result.success).toBe(false);
        if (!result.success) {
            expect(getFirstIssueI18n(result)).toBe('zod.custom.duration.min');
        }
    });

    it('rejects value above maximum', () => {
        const schema = zodProtoDurationSchema({
            required: true,
            max: secondsToDuration(120),
        });
        const result = schema.safeParse(secondsToDuration(180));

        expect(result.success).toBe(false);
        if (!result.success) {
            expect(getFirstIssueI18n(result)).toBe('zod.custom.duration.max');
        }
    });

    it('supports custom key overrides', () => {
        const schema = zodProtoDurationSchema({
            required: true,
            i18n: {
                required: 'zod.custom.duration.required_custom',
            },
        });
        const result = schema.safeParse(undefined);

        expect(result.success).toBe(false);
        if (!result.success) {
            expect(getFirstIssueI18n(result)).toBe('zod.custom.duration.required_custom');
        }
    });
});

describe('zodDurationMinMaxPair', () => {
    it('requires both fields', () => {
        const schema = zodDurationMinMaxPair();
        const result = schema.safeParse({});

        expect(result.success).toBe(false);
        if (!result.success) {
            const issueI18n = result.error.issues.map((i) => (i as { params?: { i18n?: string } }).params?.i18n);
            expect(issueI18n).toContain('zod.custom.duration.required');
        }
    });

    it('rejects min greater than max and attaches issue to maxDuration', () => {
        const schema = zodDurationMinMaxPair();
        const result = schema.safeParse({
            minDuration: secondsToDuration(3600),
            maxDuration: secondsToDuration(600),
        });

        expect(result.success).toBe(false);
        if (!result.success) {
            const issue = result.error.issues.find(
                (i) => (i as { params?: { i18n?: string } }).params?.i18n === 'zod.custom.duration.min_max_order',
            ) as { path?: (string | number)[] } | undefined;
            expect(issue?.path).toEqual(['maxDuration']);
        }
    });

    it('accepts equal or ascending ranges', () => {
        const schema = zodDurationMinMaxPair();

        const equal = schema.safeParse({
            minDuration: secondsToDuration(600),
            maxDuration: secondsToDuration(600),
        });
        const ascending = schema.safeParse({
            minDuration: secondsToDuration(600),
            maxDuration: secondsToDuration(1200),
        });

        expect(equal.success).toBe(true);
        expect(ascending.success).toBe(true);
    });

    it('allows a missing pair when requiredWhen returns false', () => {
        const schema = zodDurationMinMaxPair({
            requiredWhen: (value) => value.requiresExpiration === true,
        }).extend({
            requiresExpiration: z.boolean(),
        });

        const result = schema.safeParse({
            requiresExpiration: false,
        });

        expect(result.success).toBe(true);
    });

    it('requires both fields when requiredWhen returns true', () => {
        const schema = zodDurationMinMaxPair({
            requiredWhen: (value) => value.requiresExpiration === true,
        }).extend({
            requiresExpiration: z.boolean(),
        });

        const result = schema.safeParse({
            requiresExpiration: true,
        });

        expect(result.success).toBe(false);
        if (!result.success) {
            const issues = result.error.issues.map((issue) => ({
                path: issue.path,
                i18n: (issue as { params?: { i18n?: string } }).params?.i18n,
            }));
            expect(issues).toContainEqual({
                path: ['minDuration'],
                i18n: 'zod.custom.duration.required',
            });
            expect(issues).toContainEqual({
                path: ['maxDuration'],
                i18n: 'zod.custom.duration.required',
            });
        }
    });

    it('keeps the range order validation when extended with another field', () => {
        const schema = zodDurationMinMaxPair({
            requiredWhen: (value) => value.requiresExpiration === true,
        }).extend({
            requiresExpiration: z.boolean(),
        });

        const result = schema.safeParse({
            requiresExpiration: true,
            minDuration: secondsToDuration(3600),
            maxDuration: secondsToDuration(600),
        });

        expect(result.success).toBe(false);
        if (!result.success) {
            const issue = result.error.issues.find(
                (i) => (i as { params?: { i18n?: string } }).params?.i18n === 'zod.custom.duration.min_max_order',
            ) as { path?: (string | number)[] } | undefined;
            expect(issue?.path).toEqual(['maxDuration']);
        }
    });
});

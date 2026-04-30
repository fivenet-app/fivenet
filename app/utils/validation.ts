import { z } from 'zod';
import { durationToSeconds } from '~/utils/duration';
import type { Duration } from '~~/gen/ts/google/protobuf/duration';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';

export const pageNumberSchema = z.coerce.number().int().nonnegative().min(1).max(999_999_999).prefault(1);

type DurationI18nKeys = {
    invalid: string;
    required: string;
    min: string;
    max: string;
    minMaxOrder: string;
};

const defaultDurationI18nKeys: DurationI18nKeys = {
    invalid: 'zod.custom.duration.invalid',
    required: 'zod.custom.duration.required',
    min: 'zod.custom.duration.min',
    max: 'zod.custom.duration.max',
    minMaxOrder: 'zod.custom.duration.min_max_order',
};

function getDurationI18nKeys(overrides?: Partial<DurationI18nKeys>): DurationI18nKeys {
    return {
        ...defaultDurationI18nKeys,
        ...overrides,
    };
}

function isDuration(value: unknown): value is Duration {
    if (!value || typeof value !== 'object') {
        return false;
    }

    const duration = value as Partial<Duration>;
    return (
        typeof duration.seconds === 'number' &&
        Number.isFinite(duration.seconds) &&
        typeof duration.nanos === 'number' &&
        Number.isFinite(duration.nanos)
    );
}

function addI18nIssue(
    ctx: z.RefinementCtx,
    i18n: string,
    options?: {
        path?: (string | number)[];
        params?: Record<string, unknown>;
    },
): void {
    ctx.addIssue({
        code: 'custom',
        path: options?.path,
        params: {
            i18n,
            ...(options?.params ?? {}),
        },
    });
}

interface ProtoDurationSchemaOptions {
    required?: boolean;
    min?: Duration;
    max?: Duration;
    i18n?: Partial<DurationI18nKeys>;
}

export function zodProtoDurationSchema(options?: ProtoDurationSchemaOptions) {
    const required = options?.required ?? false;
    const keys = getDurationI18nKeys(options?.i18n);

    return z.custom<Duration | undefined>().superRefine((value, ctx) => {
        if (value === undefined || value === null) {
            if (required) {
                addI18nIssue(ctx, keys.required);
            }
            return;
        }

        if (!isDuration(value)) {
            addI18nIssue(ctx, keys.invalid);
            return;
        }

        if (value.seconds < 0 || value.nanos < 0 || value.nanos > 999_999_999 || (value.seconds === 0 && value.nanos < 0)) {
            addI18nIssue(ctx, keys.invalid);
            return;
        }

        const seconds = durationToSeconds(value);
        if (seconds < 0) {
            addI18nIssue(ctx, keys.invalid);
            return;
        }

        if (options?.min) {
            const minimum = durationToSeconds(options.min);
            if (seconds < minimum) {
                addI18nIssue(ctx, keys.min, { params: { minimum } });
                return;
            }
        }

        if (options?.max) {
            const maximum = durationToSeconds(options.max);
            if (seconds > maximum) {
                addI18nIssue(ctx, keys.max, { params: { maximum } });
            }
        }
    });
}

interface DurationMinMaxPairOptions {
    required?: boolean;
    requiredWhen?: (value: DurationMinMaxPairValue) => boolean;
    min?: Duration;
    max?: Duration;
    i18n?: Partial<DurationI18nKeys>;
}

type DurationMinMaxPairValue = {
    minDuration?: Duration;
    maxDuration?: Duration;
} & Record<string, unknown>;

export function zodDurationMinMaxPair(options?: DurationMinMaxPairOptions) {
    const keys = getDurationI18nKeys(options?.i18n);
    const durationSchema = zodProtoDurationSchema({
        required: false,
        min: options?.min,
        max: options?.max,
        i18n: options?.i18n,
    }).optional();

    return z
        .object({
            minDuration: durationSchema,
            maxDuration: durationSchema,
        })
        .superRefine((value, ctx) => {
            const pairRequired = options?.requiredWhen?.(value as DurationMinMaxPairValue) ?? options?.required ?? true;

            if (pairRequired) {
                let missingRequiredDuration = false;

                if (value.minDuration === undefined || value.minDuration === null) {
                    addI18nIssue(ctx, keys.required, { path: ['minDuration'] });
                    missingRequiredDuration = true;
                }

                if (value.maxDuration === undefined || value.maxDuration === null) {
                    addI18nIssue(ctx, keys.required, { path: ['maxDuration'] });
                    missingRequiredDuration = true;
                }

                if (missingRequiredDuration) {
                    return;
                }
            }

            if (!isDuration(value.minDuration) || !isDuration(value.maxDuration)) {
                return;
            }

            if (durationToSeconds(value.minDuration) > durationToSeconds(value.maxDuration)) {
                addI18nIssue(ctx, keys.minMaxOrder, { path: ['maxDuration'] });
            }
        });
}

export const userAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    userId: z.coerce.number(),
    user: z.custom<UserShort>().optional(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function userAccessEntries(t: (key: string, params?: any) => string) {
    return z.array(userAccessEntry).superRefine((entries, ctx) => {
        const seen = new Set();
        entries.forEach((entry, index) => {
            const key = `${entry.userId}`;

            if (seen.has(key)) {
                ctx.addIssue({
                    code: 'custom',
                    message: t('zod.custom.access_entry.duplicate_user'),
                    path: [index, 'userId'],
                });
            } else {
                seen.add(key);
            }
        });
    });
}

export const jobAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    job: z.coerce.string().nonempty(),
    minimumGrade: z.coerce.number().nonnegative(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

// Extend the jobsAccessEntries schema to validate duplicates
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function jobsAccessEntries(t: (key: string, params?: any) => string) {
    return z.array(jobAccessEntry).superRefine((entries, ctx) => {
        const seen = new Set();
        entries.forEach((entry, index) => {
            const key = `${entry.job}-${entry.minimumGrade}`;

            if (seen.has(key)) {
                ctx.addIssue({
                    code: 'custom',
                    message: t('zod.custom.access_entry.duplicate_job_grade'),
                    path: [index, 'minimumGrade'],
                });
            } else {
                seen.add(key);
            }
            return entries;
        });
    });
}

export const qualificationAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    qualificationId: z.coerce.number(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function qualificationAccessEntries(t: (key: string, params?: any) => string) {
    return z.array(qualificationAccessEntry).superRefine((entries, ctx) => {
        const seen = new Set();
        entries.forEach((entry, index) => {
            const key = `${entry.qualificationId}`;

            if (seen.has(key)) {
                ctx.addIssue({
                    code: 'custom',
                    message: t('zod.custom.access_entry.duplicate_qualification'),
                    path: [index, 'qualificationId'],
                });
            } else {
                seen.add(key);
            }
        });
    });
}

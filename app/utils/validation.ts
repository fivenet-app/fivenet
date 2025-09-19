import { z } from 'zod';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export const pageNumberSchema = z.coerce.number().int().nonnegative().min(1).max(999_999_999).prefault(1);

export const zodDurationSchema = z
    .number()
    .nonnegative()
    .max(999_999)
    .superRefine((duration, ctx) => {
        if (!duration || duration < 0) {
            ctx.addIssue({
                code: 'custom',
                params: {
                    i18n: 'zod.errors.custom_types.duration.invalid',
                },
            });
            return;
        }

        const d = duration.toString();
        if (!/^\d+(\.\d+)?$/.test(d)) {
            ctx.addIssue({
                code: 'custom',
                params: {
                    i18n: 'zod.errors.custom_types.duration.invalid',
                },
            });
            return;
        }

        const val = toDuration(d);
        if (val.seconds <= 0 || val.nanos < 0 || (val.seconds === 0 && val.nanos > 0)) {
            ctx.addIssue({
                code: 'custom',
                params: {
                    i18n: 'zod.errors.custom_types.duration.invalid',
                },
            });
            return;
        }
    });

export const userAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    userId: z.coerce.number(),
    user: z.custom<UserShort>().optional(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

export const jobAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    job: z.coerce.string().nonempty(),
    minimumGrade: z.coerce.number().nonnegative(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

export const qualificationAccessEntry = z.object({
    id: z.coerce.number(),
    targetId: z.coerce.number(),
    qualificationId: z.coerce.number(),
    access: z.coerce.number().nonnegative(),
    required: z.coerce.boolean().optional(),
});

import { z } from 'zod';

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

import { z } from 'zod';

export const zodDurationSchema = z
    .number()
    .nonnegative()
    .max(999_999)
    .superRefine((duration, ctx) => {
        if (!duration || duration < 0) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                params: {
                    i18n: 'zodI18n.errors.custom.duration.invalid',
                },
            });
            return false;
        }

        const d = duration.toString();
        if (!/^\d+(\.\d+)?$/.test(d)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                params: {
                    i18n: 'zodI18n.errors.custom.duration.invalid',
                },
            });
            return false;
        }

        const val = toDuration(d);
        if (val.seconds <= 0 || val.nanos < 0 || (val.seconds === 0 && val.nanos > 0)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                params: {
                    i18n: 'zodI18n.errors.custom.duration.invalid',
                },
            });
            return false;
        }

        return true;
    });

export function zodFileSingleSchema(
    fileSize: number,
    types: string[],
    optional?: boolean,
): z.ZodEffects<z.ZodType<FileList, z.ZodTypeDef, FileList>, FileList, FileList> {
    return z.custom<FileList>().superRefine((files, ctx) => {
        if (!files || files.length === 0 || !files[0]) {
            if (!optional) {
                ctx.addIssue({
                    code: z.ZodIssueCode.custom,
                    params: {
                        i18n: 'zodI18n.errors.custom.filelist.required',
                    },
                });
                return false;
            }

            return true;
        }

        if (!types.includes(files[0].type)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                params: {
                    i18n: 'filelist.errors.custom.wrong_file_type',
                    types: types,
                },
            });
            return false;
        }

        if (files[0].size > fileSize) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                params: {
                    i18n: 'zodI18n.errors.custom.filelist.wrong_file_type',
                    size: Math.ceil(fileSize / 10240),
                },
            });
            return false;
        }

        return true;
    });
}

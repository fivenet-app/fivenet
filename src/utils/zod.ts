import { z } from 'zod';

export const zodDurationSchema = z
    .string()
    .min(2)
    .max(10)
    .superRefine((duration, ctx) => {
        if (!duration) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'zodI18n.custom.duration.invalid',
            });
            return false;
        }

        if (!/^\d+(\.\d+)?s$/.test(duration)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'zodI18n.custom.duration.invalid',
            });
            return false;
        }

        const val = toDuration(duration);
        if (val.seconds < 0 || val.nanos < 0 || (val.seconds === 0 && val.nanos > 0)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'zodI18n.custom.duration.invalid',
            });
            return false;
        }

        return true;
    })
    .transform((val) => toDuration(val));

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
                    message: 'zodI18n.custom.filelist.required',
                });
                return false;
            }

            return true;
        }

        if (!types.includes(files[0].type)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'zodI18n.custom.filelist.wrong_file_type',
                params: {
                    types: types,
                },
            });
            return false;
        }

        if (files[0].size > fileSize) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'zodI18n.custom.filelist.wrong_file_type',
                params: {
                    size: Math.ceil(fileSize / 10240),
                },
            });
            return false;
        }

        return true;
    });
}

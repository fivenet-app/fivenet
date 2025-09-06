import { z } from 'zod';
import type { UserShort } from '~~/gen/ts/resources/users/users';

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
    job: z.string().nonempty(),
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

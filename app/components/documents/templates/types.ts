import { z } from 'zod';

export interface ObjectSpecsValue {
    req: boolean;
    min: number;
    max: number;
}

export const zWorkflow = z.object({
    reminders: z.union([
        z.object({
            reminder: z.literal(false),
            reminderSettings: z.object({
                reminders: z
                    .object({
                        duration: z.coerce.number().max(60).positive().optional(),
                        message: z.string().max(1024),
                    })
                    .array()
                    .max(3)
                    .default([]),
            }),
        }),

        z.object({
            reminder: z.literal(true),
            reminderSettings: z.object({
                reminders: z
                    .object({
                        duration: z.coerce.number().max(60).positive(),
                        message: z.string().min(3).max(1024),
                    })
                    .array()
                    .max(3)
                    .default([]),
            }),
        }),
    ]),

    autoClose: z.union([
        z.object({
            autoClose: z.literal(false),
            autoCloseSettings: z.object({
                duration: z.coerce.number().max(60).positive(),
                message: z.string().max(1024).default(''),
            }),
        }),

        z.object({
            autoClose: z.literal(true),
            autoCloseSettings: z.object({
                duration: z.coerce.number().max(60).positive(),
                message: z.string().min(3).max(1024).default(''),
            }),
        }),
    ]),
});

export type zWorkflowSchema = z.output<typeof zWorkflow>;

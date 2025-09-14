import { z } from 'zod';

export const zWorkflow = z.object({
    reminders: z.union([
        z.object({
            reminder: z.literal(false),
            reminderSettings: z.object({
                reminders: z
                    .object({
                        duration: z.coerce.number().max(60).positive().optional(),
                        message: z.coerce.string().max(1024),
                    })
                    .array()
                    .max(3)
                    .default([]),
                maxReminderCount: z.coerce.number().min(1).max(10).default(10),
            }),
        }),

        z.object({
            reminder: z.literal(true),
            reminderSettings: z.object({
                reminders: z
                    .object({
                        duration: z.coerce.number().max(60).positive(),
                        message: z.coerce.string().min(3).max(1024),
                    })
                    .array()
                    .max(3)
                    .default([]),
                maxReminderCount: z.coerce.number().min(1).max(10).default(10),
            }),
        }),
    ]),

    autoClose: z.union([
        z.object({
            autoClose: z.literal(false),
            autoCloseSettings: z.object({
                duration: z.coerce.number().max(60).positive(),
                message: z.coerce.string().max(1024).default(''),
            }),
        }),

        z.object({
            autoClose: z.literal(true),
            autoCloseSettings: z.object({
                duration: z.coerce.number().max(60).positive(),
                message: z.coerce.string().min(3).max(1024).default(''),
            }),
        }),
    ]),
});

export type zWorkflowSchema = z.output<typeof zWorkflow>;

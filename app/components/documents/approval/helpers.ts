import type { BadgeProps } from '@nuxt/ui';
import { z } from 'zod';
import { ApprovalAssigneeKind, ApprovalStatus, ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';

export function approvalTaskStatusToColor(status: ApprovalTaskStatus): BadgeProps['color'] {
    switch (status) {
        case ApprovalTaskStatus.APPROVED:
            return 'success';

        case ApprovalTaskStatus.CANCELLED:
            return 'warning';

        case ApprovalTaskStatus.DECLINED:
            return 'error';

        case ApprovalTaskStatus.EXPIRED:
            return 'gray';

        case ApprovalTaskStatus.UNSPECIFIED:
        case ApprovalTaskStatus.PENDING:
        default:
            return 'info';
    }
}

export function approvalStatusToColor(status: ApprovalStatus): BadgeProps['color'] {
    switch (status) {
        case ApprovalStatus.APPROVED:
            return 'success';

        case ApprovalStatus.REVOKED:
            return 'warning';

        case ApprovalStatus.DECLINED:
            return 'error';

        case ApprovalStatus.UNSPECIFIED:
        default:
            return 'info';
    }
}

export function getZApprovalTask(startJobGrade: number) {
    return z.union([
        z.object({
            ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.JOB_GRADE),
            userId: z.coerce.number(),
            job: z.coerce.string().optional(),
            minimumGrade: z.coerce.number().min(startJobGrade).optional(),
            label: z.string().max(120).default(''),
            signatureRequired: z.coerce.boolean().default(false),
            slots: z.coerce.number().min(1).max(10).optional().default(1),
            dueAt: z.date().optional(),
            comment: z.coerce.string().max(255).optional(),
        }),
        z.object({
            ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.JOB_GRADE),
            userId: z.coerce.number().optional().default(0),
            job: z.coerce.string(),
            minimumGrade: z.coerce.number().min(startJobGrade),
            label: z.string().max(120).default(''),
            signatureRequired: z.coerce.boolean().default(false),
            slots: z.coerce.number().min(1).max(10).optional().default(1),
            dueAt: z.date().optional(),
            comment: z.coerce.string().max(255).optional(),
        }),
    ]);
}

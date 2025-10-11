import type { BadgeProps } from '@nuxt/ui';
import { ApprovalStatus, ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';

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

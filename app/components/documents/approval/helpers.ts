import type { BadgeProps } from '@nuxt/ui';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';

export function approvalTaskStatusToColor(status: ApprovalTaskStatus): BadgeProps['color'] {
    switch (status) {
        case ApprovalTaskStatus.APPROVED:
            return 'green';

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

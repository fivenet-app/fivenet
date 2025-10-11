<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import { approvalTaskStatusToColor } from './helpers';

interface Props extends /* @vue-ignore */ BadgeProps {
    status?: ApprovalTaskStatus | undefined;
}

defineProps<Props>();
</script>

<template>
    <UBadge
        :label="$t(`enums.documents.ApprovalTaskStatus.${ApprovalTaskStatus[status ?? 0]}`)"
        :color="approvalTaskStatusToColor(status ?? 0)"
    />
</template>

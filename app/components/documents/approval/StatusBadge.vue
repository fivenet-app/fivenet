<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { ApprovalStatus } from '~~/gen/ts/resources/documents/approval';
import { approvalStatusToColor } from './helpers';

interface Props extends /* @vue-ignore */ BadgeProps {
    status?: ApprovalStatus | undefined;
}

defineProps<Props>();
</script>

<template>
    <UBadge
        :label="$t(`enums.documents.ApprovalStatus.${ApprovalStatus[status ?? 0]}`)"
        :color="approvalStatusToColor(status ?? 0)"
    />
</template>

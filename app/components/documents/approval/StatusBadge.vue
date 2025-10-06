<script lang="ts" setup>
import { ApprovalStatus } from '~~/gen/ts/resources/documents/approval';
import { approvalStatusToColor } from './helpers';

defineProps<{
    status: ApprovalStatus;
}>();
</script>

<template>
    <UBadge :label="$t(`enums.documents.ApprovalStatus.${ApprovalStatus[status]}`)" :color="approvalStatusToColor(status)" />
</template>

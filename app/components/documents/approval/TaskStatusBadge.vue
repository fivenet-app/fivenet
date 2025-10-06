<script lang="ts" setup>
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import { approvalTaskStatusToColor } from './helpers';

defineProps<{
    status: ApprovalTaskStatus;
}>();
</script>

<template>
    <UBadge
        :label="$t(`enums.documents.ApprovalTaskStatus.${ApprovalTaskStatus[status]}`)"
        :color="approvalTaskStatusToColor(status)"
    />
</template>

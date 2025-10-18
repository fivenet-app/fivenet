<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModalWithReason from '~/components/partials/ConfirmModalWithReason.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { type ApprovalPolicy, ApprovalStatus } from '~~/gen/ts/resources/documents/approval';
import type { ListApprovalsResponse } from '~~/gen/ts/services/documents/approval';
import StatusBadge from './StatusBadge.vue';

const props = defineProps<{
    documentId: number;
    policy?: ApprovalPolicy;
    taskId?: number;
}>();

const overlay = useOverlay();

const { can } = useAuth();

const approvalClient = await getDocumentsApprovalClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `approval-drawer-${props.documentId}-approvals`,
    () => listApprovals(),
);

async function listApprovals(): Promise<ListApprovalsResponse> {
    try {
        const call = approvalClient.listApprovals({
            documentId: props.documentId,
            taskId: props.taskId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function revokeApproval(approvalId: number, comment: string = '') {
    try {
        const call = approvalClient.revokeApproval({
            approvalId: approvalId,
            comment: comment,
        });
        await call;

        await refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const confirmModal = overlay.create(ConfirmModalWithReason);
</script>

<!-- eslint-disable vue/no-v-html -->
<template>
    <UCard :ui="{ body: 'p-0 sm:p-0', footer: 'p-0 sm:px-2' }">
        <template #header>
            <div class="flex items-center justify-between gap-1">
                <h3 class="flex-1 text-base leading-6 font-semibold">
                    {{ $t('common.approvals', 2) }}
                </h3>

                <slot name="header" />

                <UTooltip :text="$t('common.refresh')">
                    <UButton variant="link" icon="i-mdi-refresh" @click="() => refresh()" />
                </UTooltip>
            </div>
        </template>

        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.approvals', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.approvals', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="data?.approvals.length === 0" icon="i-mdi-approval" :type="$t('common.approvals', 2)" />

        <ul v-else class="divide-y divide-default" role="list">
            <li
                v-for="approval in data?.approvals"
                :key="approval.id"
                class="relative flex flex-col justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
            >
                <div class="flex min-w-0 flex-1 gap-x-2">
                    <div class="min-w-0 flex-1 flex-auto">
                        <p class="text-sm leading-6 font-semibold text-toned">
                            <span class="inline-flex items-center gap-2 text-highlighted">
                                <StatusBadge :status="approval.status" />

                                <span class="line-clamp-2 hover:line-clamp-5">
                                    {{ approval.comment || $t('common.no_comment') }}
                                </span>
                            </span>
                        </p>
                        <p class="inline-flex gap-1 text-xs leading-6 font-semibold text-toned">
                            {{ $t('common.created_by') }}

                            <CitizenInfoPopover
                                v-if="approval.userId"
                                class="ml-1"
                                :user="approval.user"
                                :user-id="approval.userId"
                                size="xs"
                            />
                        </p>
                    </div>

                    <div v-if="can('documents.ApprovalService/RevokeApproval').value" class="inline-flex items-center gap-2">
                        <UTooltip v-if="approval.status !== ApprovalStatus.REVOKED" :text="$t('common.revoke')">
                            <UButton
                                icon="i-mdi-cancel"
                                color="red"
                                variant="link"
                                @click="
                                    confirmModal.open({
                                        confirm: (reason) => revokeApproval(approval.id, reason),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </div>

                <div v-if="approval.payloadSvg || policy?.signatureRequired" class="flex items-center justify-center py-2">
                    <div
                        v-if="approval.payloadSvg"
                        class="mx-auto inline-block w-[clamp(220px,80vw,560px)] rounded-md bg-white p-2 shadow-sm sm:w-[clamp(280px,60vw,680px)]"
                    >
                        <div class="bg-svg" v-html="approval.payloadSvg" />
                    </div>
                    <UBadge v-else-if="policy?.signatureRequired" :label="$t('common.no_signature')" />
                </div>
            </li>
        </ul>
    </UCard>
</template>

<style scoped>
/* Make any injected <svg> responsive without editing its markup */
.bg-svg :is(svg) {
    display: block;
    width: 100%;
    height: auto;
    max-width: none; /* allow it to exceed intrinsic width */
}

/* If the incoming SVG has inline width/height attributes, this helps to override it */
.bg-svg :is(svg)[width] {
    width: 100% !important;
}
.bg-svg :is(svg)[height] {
    height: auto !important;
}
</style>

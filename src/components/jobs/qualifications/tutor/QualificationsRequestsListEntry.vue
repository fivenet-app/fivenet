<script lang="ts" setup>
import { useConfirmDialog } from '@vueuse/core';
import { CheckBoldIcon, CloseThickIcon, StarIcon, TrashCanIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationRequest, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { requestStatusToTextColor } from '~/components/jobs/qualifications/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import type { DeleteQualificationReqResponse } from '~~/gen/ts/services/qualifications/qualifications';

withDefaults(
    defineProps<{
        request: QualificationRequest;
        canSubmit?: boolean;
    }>(),
    {
        canSubmit: true,
    },
);

const emits = defineEmits<{
    (e: 'selectedRequestStatus', status?: RequestStatus): void;
    (e: 'gradeRequest'): void;
    (e: 'delete'): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteQualificationRequest(qualificationId: string, userId: number): Promise<DeleteQualificationReqResponse> {
    try {
        const call = $grpc.getQualificationsClient().deleteQualificationReq({
            qualificationId,
            userId,
        });
        const { response } = await call;

        emits('delete');

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();
onConfirm(async (request: QualificationRequest) => deleteQualificationRequest(request.qualificationId, request.userId));
</script>

<template>
    <tr>
        <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(request)" />

        <td>
            <CitizenInfoPopover :user="request.user" />
        </td>
        <td>
            <span v-if="request.userComment">{{ request.userComment }}</span>
        </td>
        <td>
            <div class="flex flex-initial flex-row gap-1 rounded-full">
                <template v-if="request.status !== undefined">
                    <span class="text-sm font-medium" :class="requestStatusToTextColor(request.status)">
                        <span class="font-semibold">{{
                            $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`)
                        }}</span>
                    </span>
                </template>
            </div>
        </td>
        <td>
            <p v-if="request.createdAt" class="mt-1 text-sm leading-5 text-gray-300">
                <GenericTime :value="request.createdAt" />
            </p>
        </td>
        <td>
            <p v-if="request.approvedAt" class="mt-1 text-sm leading-5 text-gray-300">
                <GenericTime :value="request.approvedAt" />
            </p>
        </td>
        <td>
            <CitizenInfoPopover v-if="request.approver" :user="request.approver" />
        </td>
        <td class="flex items-center gap-2">
            <UButton
                v-if="request.status !== RequestStatus.DENIED"
                :disabled="!canSubmit"
                class="flex-initial"
                :class="[!canSubmit ? 'disabled text-base-500 hover:text-base-400' : 'text-error-500 hover:text-error-400']"
                @click="$emit('selectedRequestStatus', RequestStatus.DENIED)"
            >
                <CloseThickIcon class="size-6" />
            </UButton>
            <UButton
                v-if="request.status !== RequestStatus.ACCEPTED"
                :disabled="!canSubmit"
                class="flex-initial"
                :class="[!canSubmit ? 'disabled text-base-500 hover:text-base-400' : 'text-success-500 hover:text-success-400']"
                @click="$emit('selectedRequestStatus', RequestStatus.ACCEPTED)"
            >
                <CheckBoldIcon class="size-6" />
            </UButton>
            <UButton
                v-if="request.status === RequestStatus.ACCEPTED"
                :disabled="!canSubmit"
                class="flex-initial"
                :class="[!canSubmit ? 'disabled text-base-500 hover:text-base-400' : 'text-yellow-500 hover:text-yellow-400']"
                @click="$emit('gradeRequest')"
            >
                <StarIcon class="size-6" />
            </UButton>
            <UButton
                v-if="can('QualificationsService.DeleteQualificationReq')"
                :disabled="!canSubmit"
                class="flex-initial text-primary-400 hover:text-primary-500"
                @click="reveal()"
            >
                <TrashCanIcon class="size-6" />
            </UButton>
        </td>
    </tr>
</template>

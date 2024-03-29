<script lang="ts" setup>
import { useConfirmDialog } from '@vueuse/core';
import { TrashCanIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationResult, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToTextColor } from '~/components/jobs/qualifications/helpers';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import type { DeleteQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';

defineProps<{
    result: QualificationResult;
}>();

const emits = defineEmits<{
    (e: 'delete'): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteQualificationResult(resultId: string): Promise<DeleteQualificationResultResponse> {
    try {
        const call = $grpc.getQualificationsClient().deleteQualificationResult({
            resultId,
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
onConfirm(async (resultId: string) => deleteQualificationResult(resultId));
</script>

<template>
    <tr>
        <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(result.id)" />

        <td>
            <CitizenInfoPopover :user="result.user" />
        </td>
        <td>
            <div class="flex flex-initial flex-row gap-1 rounded-full">
                <template v-if="result.status !== undefined">
                    <span class="text-sm font-medium" :class="resultStatusToTextColor(result.status)">
                        <span class="font-semibold">{{
                            $t(`enums.qualifications.ResultStatus.${ResultStatus[result.status]}`)
                        }}</span>
                    </span>
                </template>
            </div>
        </td>
        <td>
            <span v-if="result.score">{{ result.score }}</span>
        </td>
        <td>
            <p v-if="result.summary" class="mt-1 text-sm leading-5 text-gray-300">
                {{ result.summary }}
            </p>
        </td>
        <td>
            <p v-if="result.createdAt" class="mt-1 text-sm leading-5 text-gray-300">
                <GenericTime :value="result.createdAt" />
            </p>
        </td>
        <td>
            <CitizenInfoPopover :user="result.creator" />
        </td>
        <td class="flex flex-row justify-end">
            <button
                v-if="can('QualificationsService.DeleteQualificationResult')"
                type="button"
                class="flex-initial text-primary-500 hover:text-primary-400"
                @click="reveal()"
            >
                <TrashCanIcon class="size-5 text-primary-500" aria-hidden="true" />
            </button>
        </td>
    </tr>
</template>

<script lang="ts" setup>
import { TrashCanIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationResult, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToTextColor } from '~/components/jobs/qualifications/helpers';

defineProps<{
    result: QualificationResult;
}>();

defineEmits<{
    (e: 'delete'): void;
}>();
</script>

<template>
    <tr>
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
            <button type="button" class="flex-initial text-primary-500 hover:text-primary-400" @click="$emit('delete')">
                <TrashCanIcon class="h-5 w-5 text-primary-500" aria-hidden="true" />
            </button>
        </td>
    </tr>
</template>

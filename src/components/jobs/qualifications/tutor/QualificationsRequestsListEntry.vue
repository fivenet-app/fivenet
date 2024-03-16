<script lang="ts" setup>
import { CheckBoldIcon, CloseThickIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationRequest, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    request: QualificationRequest;
}>();

defineEmits<{
    (e: 'selected'): void;
}>();

const canSubmit = ref(false);
</script>

<template>
    <tr>
        <td>{{ request.qualification?.abbreviation }}: {{ request.qualification?.title }}</td>
        <td>
            <span v-if="request.userComment">{{ request.userComment }}</span>
        </td>
        <td>
            <div class="flex flex-initial flex-row gap-1 rounded-full">
                <template v-if="request.status !== undefined">
                    <span class="text-sm font-medium text-info-400">
                        <span class="font-semibold">{{
                            $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`)
                        }}</span>
                    </span>
                </template>
            </div>
        </td>
        <td>
            <p v-if="request.createdAt" class="mt-1 text-sm leading-5 text-gray-300">
                {{ $t('common.created_at') }} <GenericTime :value="request.createdAt" />
            </p>
        </td>
        <td class="flex items-center gap-2">
            <button
                type="button"
                :disabled="!canSubmit"
                class="rounded flex flex-1 justify-center px-2 py-2 text-sm font-semibold text-neutral"
                :class="[
                    !canSubmit
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                ]"
            >
                <CloseThickIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
            </button>
            <button
                type="button"
                :disabled="!canSubmit"
                class="rounded flex flex-1 justify-center px-2 py-2 text-sm font-semibold text-neutral"
                :class="[
                    !canSubmit
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                ]"
            >
                <CheckBoldIcon class="h-5 w-5 text-success-400" aria-hidden="true" />
            </button>
        </td>
    </tr>
</template>

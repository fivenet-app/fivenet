<script lang="ts" setup>
import { ChevronRightIcon, ListStatusIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationResult, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    result: QualificationResult;
}>();
</script>

<template>
    <li class="relative flex justify-between px-4 py-5">
        <div class="flex min-w-0 gap-x-4">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <NuxtLink :to="{ name: 'jobs-qualifications-id', params: { id: result.qualificationId } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ result.qualification?.abbreviation }}: {{ result.qualification?.title }}
                    </NuxtLink>
                </p>
                <p class="mt-1 flex text-xs leading-5 text-gray-300">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.score') }}: {{ result.score }}</span>
                        <span v-if="result.summary"> ({{ $t('common.summary') }}: {{ result.summary }})</span>
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div class="flex flex-row gap-1">
                    <div v-if="result.status" class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1">
                        <ListStatusIcon class="size-5 text-info-400" />
                        <span class="text-sm font-medium text-info-700">
                            <span class="font-semibold">{{ $t('common.result') }}:</span>
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[result.status ?? 0]}`) }}
                        </span>
                    </div>
                </div>
                <p v-if="result.createdAt" class="mt-1 text-xs leading-5 text-gray-300">
                    {{ $t('common.created_at') }} <GenericTime :value="result.createdAt" />
                </p>
            </div>
            <ChevronRightIcon class="size-5 flex-none text-gray-300" />
        </div>
    </li>
</template>

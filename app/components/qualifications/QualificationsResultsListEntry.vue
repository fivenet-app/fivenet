<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { QualificationResult } from '~~/gen/ts/resources/qualifications/qualifications';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToBadgeColor } from './helpers';

defineProps<{
    result: QualificationResult;
}>();
</script>

<template>
    <li
        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 relative flex justify-between border-white px-2 py-2 sm:px-4 dark:border-gray-900"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <ULink :to="{ name: 'qualifications-id', params: { id: result.qualificationId } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ result.qualification?.abbreviation }}: {{ result.qualification?.title }}
                    </ULink>
                </p>
                <p class="mt-1 flex text-xs leading-5">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.score') }}: {{ result.score }}</span>
                        <span v-if="result.summary"> ({{ $t('common.summary') }}: {{ result.summary }})</span>
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-2">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <UBadge v-if="result.status" class="inline-flex gap-1" :color="resultStatusToBadgeColor(result?.status ?? 0)">
                    <UIcon class="size-5" name="i-mdi-list-status" />
                    <span>
                        {{ $t('common.result') }}:
                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[result.status ?? 0]}`) }}
                    </span>
                </UBadge>

                <p v-if="result.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="result.createdAt" />
                </p>
            </div>

            <UIcon class="size-5 flex-none" name="i-mdi-chevron-right" />
        </div>
    </li>
</template>

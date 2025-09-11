<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { type QualificationResult, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToBadgeColor } from '../helpers';

defineProps<{
    result: QualificationResult;
}>();
</script>

<template>
    <li
        class="relative flex justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm leading-6 font-semibold text-toned">
                    <ULink class="text-highlighted" :to="{ name: 'qualifications-id', params: { id: result.qualificationId } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ result.qualification?.abbreviation }}:
                        {{ !result.qualification?.title ? $t('common.untitled') : result.qualification?.title }}
                    </ULink>
                </p>
                <p class="mt-1 flex text-xs leading-5">
                    <span class="inline-flex gap-1 truncate">
                        <span class="font-semibold">{{ $t('common.score') }}: {{ result.score }}</span>
                        <span v-if="result.summary"> ({{ $t('common.summary') }}: {{ result.summary }})</span>
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-2">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <UBadge
                    v-if="result.status"
                    class="inline-flex gap-1"
                    :color="resultStatusToBadgeColor(result?.status ?? 0)"
                    icon="i-mdi-list-status"
                >
                    {{ $t('common.result') }}:
                    {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[result.status ?? 0]}`) }}
                </UBadge>

                <p v-if="result.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="result.createdAt" />
                </p>
            </div>

            <UIcon class="size-5 flex-none" name="i-mdi-chevron-right" />
        </div>
    </li>
</template>

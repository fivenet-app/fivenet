<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { Qualification, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    qualification: Qualification;
}>();
</script>

<template>
    <li class="relative flex justify-between px-3 py-4">
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <ULink :to="{ name: 'jobs-qualifications-id', params: { id: qualification.id } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </ULink>
                </p>
                <p class="mt-1 flex gap-1 text-xs leading-5">
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ qualification.description ?? $t('common.na') }}
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-2">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div class="flex flex-row gap-1">
                    <UBadge v-if="qualification.result?.status" class="inline-flex gap-1">
                        <UIcon name="i-mdi-list-status" class="size-5" />
                        <span>
                            {{ $t('common.result') }}:
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </span>
                    </UBadge>

                    <UBadge v-if="qualification.closed" color="red" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-lock" class="h-auto w-5" />
                        <span>
                            {{ $t('common.close', 2) }}
                        </span>
                    </UBadge>
                    <UBadge v-else color="green" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-lock-open-variant" class="h-auto w-5" />
                        <span>
                            {{ $t('common.open', 2) }}
                        </span>
                    </UBadge>
                </div>

                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.createdAt" />
                </p>
            </div>

            <UIcon name="i-mdi-chevron-right" class="size-5 flex-none" />
        </div>
    </li>
</template>

<script lang="ts" setup>
import { ChevronRightIcon, ListStatusIcon, LockIcon, LockOpenVariantIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { Qualification, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    qualification: Qualification;
}>();
</script>

<template>
    <li class="relative flex justify-between px-4 py-5">
        <div class="flex min-w-0 gap-x-4">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <NuxtLink :to="{ name: 'jobs-qualifications-id', params: { id: qualification.id } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </NuxtLink>
                </p>
                <p class="mt-1 flex gap-1 text-xs leading-5">
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ qualification.description ?? $t('common.na') }}
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div class="flex flex-row gap-1">
                    <div
                        v-if="qualification.result?.status"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1"
                    >
                        <ListStatusIcon class="size-5 text-info-400" />
                        <span class="text-sm font-medium text-info-700">
                            <span class="font-semibold">{{ $t('common.result') }}:</span>
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </span>
                    </div>

                    <div
                        v-if="qualification.closed"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                    >
                        <LockIcon class="size-5 text-error-400" />
                        <span class="text-sm font-medium text-error-700">
                            {{ $t('common.close', 2) }}
                        </span>
                    </div>
                    <div v-else class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1">
                        <LockOpenVariantIcon class="size-5 text-success-500" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.open', 2) }}
                        </span>
                    </div>
                </div>
                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.createdAt" />
                </p>
            </div>
            <ChevronRightIcon class="size-5 flex-none" />
        </div>
    </li>
</template>

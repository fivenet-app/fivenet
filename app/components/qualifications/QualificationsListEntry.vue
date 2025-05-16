<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import type { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToBadgeColor } from './helpers';

defineProps<{
    qualification: Qualification;
}>();
</script>

<template>
    <li
        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 relative flex justify-between border-white px-2 py-2 sm:px-4 dark:border-gray-900"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <ULink
                        class="inline-flex items-center gap-2"
                        :to="{ name: 'qualifications-id', params: { id: qualification.id } }"
                    >
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        <span>{{ qualification.abbreviation }}: {{ qualification.title }}</span>

                        <UBadge v-if="qualification.public" class="inline-flex gap-1" color="black" size="xs">
                            <UIcon class="size-5" name="i-mdi-earth" />
                            <span>
                                {{ $t('common.public') }}
                            </span>
                        </UBadge>

                        <UBadge v-if="qualification?.deletedAt" class="inline-flex gap-1" color="amber" size="xs">
                            <UIcon class="size-5" name="i-mdi-calendar-remove" />
                            <span>
                                {{ $t('common.deleted') }}
                                <GenericTime :value="qualification?.deletedAt" type="long" />
                            </span>
                        </UBadge>
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
                    <UBadge
                        v-if="qualification.result?.status"
                        class="inline-flex gap-1"
                        :color="resultStatusToBadgeColor(qualification.result?.status ?? 0)"
                        size="sm"
                    >
                        <UIcon class="size-5" name="i-mdi-list-status" />
                        <span>
                            {{ $t('common.result') }}:
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </span>
                    </UBadge>

                    <OpenClosedBadge :closed="qualification.closed" />
                </div>

                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.createdAt" />
                </p>
            </div>

            <UIcon class="size-5 flex-none" name="i-mdi-chevron-right" />
        </div>
    </li>
</template>

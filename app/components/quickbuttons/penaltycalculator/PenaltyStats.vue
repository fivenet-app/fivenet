<script lang="ts" setup>
import type { PenaltiesSummary } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

const props = defineProps<{
    summary: PenaltiesSummary;
    reduction: number;
}>();

const leeway = computed(() => props.reduction / 100);
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <UPageGrid :ui="{ wrapper: 'grid-cols-1 sm:grid-cols-4' }">
            <UPageCard
                :ui="{
                    body: {
                        padding: 'px-4 py-4 sm:p-4',
                    },
                    icon: { wrapper: 'mb-1' },
                }"
                icon="i-mdi-attach-money"
            >
                <template #icon>
                    <div class="flex gap-1">
                        <UIcon name="i-mdi-attach-money" class="text-primary h-10 w-10" />
                        <div class="flex items-center gap-1.5 truncate text-base font-semibold text-gray-900 dark:text-white">
                            {{ $t('common.fine') }}
                        </div>
                    </div>
                </template>

                <template #description>
                    <div class="flex flex-col gap-1">
                        <div class="flex gap-1">
                            <span class="text-4xl font-semibold tracking-tight">
                                {{ summary.fine.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ') }}
                            </span>
                            <span class="text-sm text-gray-500 dark:text-gray-400">$</span>
                        </div>

                        <span v-if="leeway > 0 && summary.fine > 0"> ($-{{ (summary.fine * leeway).toFixed(0) }}) </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :ui="{
                    body: {
                        padding: 'flex-1 px-4 py-4 sm:p-4',
                    },
                    icon: { wrapper: 'mb-1' },
                }"
            >
                <template #icon>
                    <div class="flex gap-1">
                        <UIcon name="i-mdi-clock" class="text-primary h-10 w-10" />
                        <div class="flex items-center gap-1.5 truncate text-base font-semibold text-gray-900 dark:text-white">
                            {{ $t('common.detention_time') }}
                        </div>
                    </div>
                </template>

                <template #description>
                    <div class="flex flex-col gap-1">
                        <div class="inline-flex gap-1">
                            <span class="text-4xl font-semibold tracking-tight">
                                {{ summary.detentionTime }}
                            </span>
                            <span class="text-sm text-gray-500 dark:text-gray-400">
                                {{ $t('common.time_ago.month', summary.detentionTime) }}
                            </span>
                        </div>

                        <span v-if="leeway > 0 && summary.detentionTime > 0">
                            (-{{ (summary.detentionTime * leeway).toFixed(0) }}
                            {{ $t('common.time_ago.month', (summary.detentionTime * leeway).toFixed(0)) }})
                        </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :ui="{
                    body: {
                        padding: 'flex-1 px-4 py-4 sm:p-4',
                    },
                    icon: { wrapper: 'mb-1' },
                }"
            >
                <template #icon>
                    <div class="flex gap-1">
                        <UIcon name="i-mdi-car" class="text-primary h-10 w-10" />
                        <div class="flex items-center gap-1.5 truncate text-base font-semibold text-gray-900 dark:text-white">
                            {{ $t('common.traffic_infraction_points', 2) }}
                        </div>
                    </div>
                </template>

                <template #description>
                    <div class="flex flex-col gap-1">
                        <div class="inline-flex gap-1">
                            <span class="text-4xl font-semibold tracking-tight">
                                {{ summary.stvoPoints }}
                            </span>
                            <span class="text-sm text-gray-500 dark:text-gray-400">
                                {{ $t('common.points', summary.stvoPoints) }}
                            </span>
                        </div>

                        <span v-if="leeway > 0 && summary.stvoPoints > 0">
                            (-{{ (summary.stvoPoints * leeway).toFixed(0) }}
                            {{ $t('common.points', (summary.stvoPoints * leeway).toFixed(0)) }})
                        </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :ui="{
                    body: {
                        padding: 'flex-1 px-4 py-4 sm:p-4',
                    },
                    icon: { wrapper: 'mb-1' },
                }"
            >
                <template #icon>
                    <div class="flex gap-1">
                        <UIcon name="i-mdi-equal" class="text-primary h-10 w-10" />
                        <div class="flex items-center gap-1.5 truncate text-base font-semibold text-gray-900 dark:text-white">
                            {{ $t('common.total_count') }}
                        </div>
                    </div>
                </template>

                <template #description>
                    <div class="inline-flex gap-1">
                        <span class="text-4xl font-semibold tracking-tight">
                            {{ summary.count }}
                        </span>
                    </div>
                </template>
            </UPageCard>
        </UPageGrid>
    </div>
</template>

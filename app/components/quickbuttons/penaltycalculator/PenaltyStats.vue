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
        <UPageGrid class="grid-cols-1 md:grid-cols-4 lg:grid-cols-4 xl:grid-cols-4">
            <UPageCard
                :ui="{
                    leadingIcon: 'mb-1',
                }"
            >
                <template #leading>
                    <div class="flex gap-1 truncate">
                        <UIcon class="h-6 w-6 text-primary" name="i-mdi-attach-money" />

                        <div class="flex items-center gap-1.5 text-base font-semibold text-highlighted">
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
                            <span class="text-sm text-muted">$</span>
                        </div>

                        <span v-if="leeway > 0 && summary.fine > 0"> ($-{{ (summary.fine * leeway).toFixed(0) }}) </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :ui="{
                    body: 'flex-1',
                    leadingIcon: 'mb-1',
                }"
            >
                <template #leading>
                    <div class="flex gap-1 truncate">
                        <UIcon class="h-6 w-6 text-primary" name="i-mdi-clock" />

                        <div class="flex items-center gap-1.5 text-base font-semibold text-highlighted">
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
                            <span class="text-sm text-muted">
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
                    body: 'flex-1',
                    leadingIcon: 'mb-1',
                }"
            >
                <template #leading>
                    <div class="flex gap-1 truncate">
                        <UIcon class="h-6 w-6 text-primary" name="i-mdi-car" />

                        <div class="flex items-center gap-1.5 text-base font-semibold text-highlighted">
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
                            <span class="text-sm text-muted">
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
                    body: 'flex-1',
                    leadingIcon: 'mb-1',
                }"
            >
                <template #leading>
                    <div class="flex gap-1 truncate">
                        <UIcon class="h-6 w-6 text-primary" name="i-mdi-equal" />

                        <div class="flex items-center gap-1.5 text-base font-semibold text-highlighted">
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

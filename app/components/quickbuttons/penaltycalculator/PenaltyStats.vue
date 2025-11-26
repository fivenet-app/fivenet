<script lang="ts" setup>
import type { PenaltiesSummary } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

const props = defineProps<{
    summary: PenaltiesSummary;
    reduction: number;
}>();

const { display, quickButtons } = useAppConfig();

const leeway = computed(() => props.reduction / 100);

const highlight = computed(() => {
    if (!quickButtons.penaltyCalculator?.warnSettings?.enabled) return;

    const result = {
        fine:
            quickButtons.penaltyCalculator?.warnSettings.fine &&
            props.summary.fine >= quickButtons.penaltyCalculator?.warnSettings.fine,
        detentionTime:
            quickButtons.penaltyCalculator?.warnSettings.detentionTime &&
            props.summary.detentionTime >= quickButtons.penaltyCalculator?.warnSettings.detentionTime,
        stvoPoints:
            quickButtons.penaltyCalculator?.warnSettings.stvoPoints &&
            props.summary.stvoPoints >= quickButtons.penaltyCalculator?.warnSettings.stvoPoints,
    };

    if (!result.fine && !result.detentionTime && !result.stvoPoints) return;

    return result;
});

const formatter = new Intl.NumberFormat(display.intlLocale, {
    style: 'currency',
    currency: display.currencyName,
    trailingZeroDisplay: 'stripIfInteger',
});
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <UPageGrid class="grid-cols-1 md:grid-cols-4 lg:grid-cols-4 xl:grid-cols-4">
            <UPageCard
                :highlight="!!highlight?.fine"
                :highlight-color="highlight?.fine ? 'warning' : undefined"
                :ui="{
                    body: 'flex-1',
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
                                {{ formatter.format(summary.fine ?? 0) }}
                            </span>
                        </div>

                        <span v-if="leeway > 0 && summary.fine > 0"> ({{ formatter.format(-(summary.fine * leeway)) }}) </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :highlight="!!highlight?.detentionTime"
                :highlight-color="highlight?.detentionTime ? 'warning' : undefined"
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
                                {{ $t('common.month', summary.detentionTime) }}
                            </span>
                        </div>

                        <span v-if="leeway > 0 && summary.detentionTime > 0">
                            (-{{ (summary.detentionTime * leeway).toFixed(0) }}
                            {{ $t('common.month', summary.detentionTime * leeway) }})
                        </span>
                    </div>
                </template>
            </UPageCard>

            <UPageCard
                :highlight="!!highlight?.stvoPoints"
                :highlight-color="highlight?.stvoPoints ? 'warning' : undefined"
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

        <UAlert
            v-if="!!highlight && quickButtons.penaltyCalculator?.warnSettings?.warnMessage"
            class="mt-3"
            color="warning"
            icon="i-mdi-warning-circle"
            :description="quickButtons.penaltyCalculator?.warnSettings?.warnMessage"
        />
    </div>
</template>

<script setup lang="ts">
import { VisAxis, VisCrosshair, VisGroupedBar, VisTooltip, VisXYContainer } from '@unovis/vue';
import { format } from 'date-fns';

const props = defineProps<{
    title: string;
    total: number;
    color: string;
    data: SeriesPoint[];
    currency?: boolean;
    detention?: boolean;
}>();

type SeriesPoint = {
    date: Date;
    value: number;
};

const cardRef = useTemplateRef<HTMLElement | null>('cardRef');

const { width } = useElementSize(cardRef);

const { format: formatNumber } = useIntlNumberFormatWithOptions({
    style: 'decimal',
    currency: undefined,
});

const { format: formatCurrency } = useIntlNumberFormat();

const formatDetention = useDetentionTimeFormatter();

const x = (_: SeriesPoint, i: number) => i;
const y = (d: SeriesPoint) => d.value;

const formatDate = (date: Date): string => format(date, 'd MMM');

const xTicks = (i: number) => {
    if (!props.data?.[i]) {
        return '';
    }

    return formatDate(props.data[i].date);
};

const template = (d?: SeriesPoint) => {
    if (!d || !(d.date instanceof Date)) {
        return '';
    }

    return `${formatDate(d.date)}<br>${props.title}: ${formatNumber(d.value)}`;
};
</script>

<template>
    <UCard ref="cardRef" :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-3' }">
        <template #header>
            <div>
                <p class="mb-1.5 text-xs text-muted uppercase">{{ $t('common.total') }}: {{ title }}</p>
                <p class="text-3xl font-semibold text-highlighted">
                    <span v-if="!detention">
                        {{ props.currency ? formatCurrency(total) : formatNumber(total) }}
                    </span>
                    <span v-else>
                        {{ formatDetention(total) }}
                    </span>
                </p>
            </div>
        </template>

        <VisXYContainer class="h-72" :data="data ?? []" :padding="{ top: 10, right: 12, left: 8, bottom: 8 }" :width="width">
            <VisGroupedBar :x="x" :y="y" :color="color" />

            <VisAxis type="x" :x="x" :tick-format="xTicks" />

            <VisCrosshair :color="color" :template="template" />

            <VisTooltip />
        </VisXYContainer>
    </UCard>
</template>

<style scoped>
.unovis-xy-container {
    --vis-crosshair-line-stroke-color: var(--ui-primary);
    --vis-crosshair-circle-stroke-color: var(--ui-bg);

    --vis-axis-grid-color: var(--ui-border);
    --vis-axis-tick-color: var(--ui-border);
    --vis-axis-tick-label-color: var(--ui-text-dimmed);

    --vis-tooltip-background-color: var(--ui-bg);
    --vis-tooltip-border-color: var(--ui-border);
    --vis-tooltip-text-color: var(--ui-text-highlighted);
}
</style>

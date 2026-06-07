<script lang="ts" setup>
import { StackedBar } from '@unovis/ts';
import type { StackedBarDataRecord } from '@unovis/ts/components/stacked-bar/types';
import { VisAxis, VisCrosshair, VisStackedBar, VisTooltip, VisXYContainer } from '@unovis/vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import type { LabelCount } from '~~/gen/ts/resources/jobs/labels/labels';
import type { GetColleagueLabelsStatsResponse } from '~~/gen/ts/services/jobs/colleagues';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const jobsColleaguesClient = await getJobsColleaguesClient();

const bodyRef = useTemplateRef<HTMLElement | null>('bodyRef');
const { height: bodyHeight } = useElementSize(bodyRef);

const cardRef = useTemplateRef<HTMLElement | null>('cardRef');
const { width } = useElementSize(cardRef);

const canSubmit = ref<boolean>(true);

async function getColleagueLabelsStats(): Promise<GetColleagueLabelsStatsResponse> {
    canSubmit.value = false;
    try {
        const { response } = await jobsColleaguesClient.getColleagueLabelsStats({
            labelIds: [],
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const {
    data: stats,
    status,
    error,
    refresh,
} = useLazyAsyncData('jobs-colleagues-labels-stats', () =>
    getColleagueLabelsStats().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);

const totalCount = computed(() => stats.value?.count.reduce((stat, sum) => sum.count + stat, 0));

const x = (_: LabelCount, i: number) => i;
const y = [(d: LabelCount) => d.count];
const color = (d: LabelCount) => d.label?.color ?? 'neutral';

const tooltipTemplate = (d: StackedBarDataRecord<{ count: number; label?: { name?: string; color?: string } }>): string =>
    `${d.datum.label?.name ?? t('common.na')}: ${d.datum.count}`;
</script>

<template>
    <UModal :title="`${$t('common.label', 2)} - ${$t('common.total_count')}`" fullscreen :ui="{ body: 'flex flex-col' }">
        <template #body>
            <div ref="bodyRef" class="flex-1 overflow-hidden">
                <UCard ref="cardRef" :ui="{ root: 'overflow-visible mx-2 my-0.5', body: '!px-0 !pt-0 !pb-3' }">
                    <template #header>
                        <div>
                            <p class="mb-1.5 text-xs text-muted uppercase">
                                {{ $t('common.total_count') }}
                            </p>
                            <p class="text-3xl font-semibold text-highlighted">
                                {{
                                    $n(totalCount ?? 0, {
                                        maximumFractionDigits: 0,
                                    })
                                }}
                            </p>
                        </div>
                    </template>

                    <DataErrorBlock v-if="error" :error="error" :retry="refresh" />
                    <DataPendingBlock
                        v-else-if="isRequestPending(status)"
                        :message="$t('common.loading', [$t('common.stats', 2)])"
                    />
                    <VisXYContainer
                        v-else
                        :data="stats?.count ?? []"
                        :padding="{ top: 40 }"
                        :width="width"
                        :height="bodyHeight - 110"
                    >
                        <VisStackedBar orientation="horizontal" :x="x" :y="y" :color="color" />

                        <VisCrosshair color="var(--ui-primary)" />
                        <VisTooltip :triggers="{ [StackedBar.selectors.bar]: tooltipTemplate }" />

                        <VisAxis type="x" :grid-line="true" />
                        <VisAxis
                            type="y"
                            :grid-line="true"
                            tick-text-color="var(--ui-text-highlighted)"
                            :num-ticks="stats?.count.length ?? 0"
                            :tick-format="(i: number) => stats?.count[i]?.label?.name ?? i.toString()"
                        />
                    </VisXYContainer>
                </UCard>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <RefreshButton
                    variant="solid"
                    :loading="!canSubmit"
                    :disabled="!canSubmit"
                    icon-only
                    @click="() => refresh()"
                />
            </UFieldGroup>
        </template>
    </UModal>
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

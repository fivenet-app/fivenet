<script lang="ts" setup>
import { StackedBar } from '@unovis/ts';
import { VisAxis, VisStackedBar, VisTooltip, VisXYContainer } from '@unovis/vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { LabelCount } from '~~/gen/ts/resources/jobs/labels';
import type { GetColleagueLabelsStatsResponse } from '~~/gen/ts/services/jobs/jobs';

const { isOpen } = useOverlay();

const jobsJobsClient = await getJobsJobsClient();

const bodyRef = useTemplateRef('bodyRef');
const { height, width } = useElementSize(bodyRef);

const canSubmit = ref(true);

async function getColleagueLabelsStats(): Promise<GetColleagueLabelsStatsResponse> {
    canSubmit.value = false;
    try {
        const { response } = await jobsJobsClient.getColleagueLabelsStats({
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
    error,
    refresh,
} = useLazyAsyncData('jobs-colleagues-labels-stats', () =>
    getColleagueLabelsStats().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);

const totalCount = computed(() => stats.value?.count.reduce((stat, sum) => sum.count + stat, 0));

const x = (_: LabelCount, i: number) => i;
const y = [(d: LabelCount) => d.count];
const color = (d: LabelCount) => d.label?.color ?? 'neutral';

const tooltipTemplate = (d: LabelCount): string => (d.label?.name ? `${d.label?.name}: ${d.count}` : '');
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" fullscreen>
        <UCard
            :ui="{
                header: {
                    base: 'h-(--ui-header-height)',
                    padding: 'p-4',
                },
                body: {
                    base: 'flex-1 flex min-h-[calc(100dvh-(2*var(--ui-header-height)))] max-h-[calc(100dvh-(2*var(--ui-header-height)))] overflow-y-auto',
                    padding: '',
                },
                footer: {
                    base: 'min-h-(--ui-header-height) max-h-(--ui-header-height)',
                    padding: 'p-4',
                },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.label', 2) }} - {{ $t('common.total_count') }}: {{ totalCount }}
                    </h3>

                    <UButton class="-my-1" color="neutral" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div ref="bodyRef" class="flex-1">
                <DataErrorBlock v-if="error" :error="error" :retry="refresh" />

                <ClientOnly v-else>
                    <VisXYContainer
                        :data="stats?.count ?? []"
                        :margin="{ top: 16, left: 32, right: 32, bottom: 16 }"
                        :height="height"
                        :width="width"
                    >
                        <VisStackedBar orientation="horizontal" :x="x" :y="y" :color="color" />
                        <VisTooltip :triggers="{ [StackedBar.selectors.bar]: tooltipTemplate }" />

                        <VisAxis type="x" :grid-line="false" :label="$t('common.count')" />
                        <VisAxis
                            type="y"
                            :grid-line="false"
                            :num-ticks="stats?.count.length ?? 0"
                            :tick-format="(i: number) => stats?.count[i]?.label?.name ?? i.toString()"
                        />
                    </VisXYContainer>
                </ClientOnly>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="neutral" block @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UTooltip :text="$t('common.refresh')">
                        <UButton icon="i-mdi-refresh" :loading="!canSubmit" :disabled="!canSubmit" @click="refresh" />
                    </UTooltip>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

<style scoped>
.unovis-xy-container {
    --vis-crosshair-line-stroke-color: rgb(var(--color-primary-500));
    --vis-crosshair-circle-stroke-color: #fff;

    --vis-axis-grid-color: rgb(var(--color-gray-200));
    --vis-axis-tick-color: rgb(var(--color-gray-200));
    --vis-axis-tick-label-color: rgb(var(--color-gray-400));

    --vis-tooltip-background-color: #fff;
    --vis-tooltip-border-color: rgb(var(--color-gray-200));
    --vis-tooltip-text-color: rgb(var(--color-gray-900));
}

.dark {
    .unovis-xy-container {
        --vis-crosshair-line-stroke-color: rgb(var(--color-primary-400));
        --vis-crosshair-circle-stroke-color: rgb(var(--color-gray-900));

        --vis-axis-grid-color: rgb(var(--color-gray-800));
        --vis-axis-tick-color: rgb(var(--color-gray-800));
        --vis-axis-tick-label-color: rgb(var(--color-gray-400));

        --vis-tooltip-background-color: rgb(var(--color-gray-900));
        --vis-tooltip-border-color: rgb(var(--color-gray-800));
        --vis-tooltip-text-color: #fff;
    }
}
</style>

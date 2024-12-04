<script lang="ts" setup>
import { VisBulletLegend, VisDonut, VisSingleContainer } from '@unovis/vue';
import type { LabelCount } from '~~/gen/ts/resources/jobs/labels';
import type { GetColleagueLabelsStatsResponse } from '~~/gen/ts/services/jobs/jobs';

const { isOpen } = useModal();

async function getColleagueLabelsStats(): Promise<GetColleagueLabelsStatsResponse> {
    try {
        const { response } = await getGRPCJobsClient().getColleagueLabelsStats({
            labelIds: [],
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: stats } = useLazyAsyncData('jobs-colleagues-labels-stats', () => getColleagueLabelsStats());

const totalCount = computed(() => stats.value?.count.reduce((stat, sum) => sum.count + stat, 0));

const value = (d: LabelCount) => d.count;
const items = computed(
    () =>
        stats.value?.count?.map((d) => ({
            name: d.label?.name ? `${d.label?.name}: ${d.count}` : '',
            color: d.label?.color ?? 'gray',
        })) ?? [],
);
const color = (d: LabelCount) => d.label?.color ?? 'gray';

const bodyRef = useTemplateRef('bodyRef');
const { height, width } = useElementSize(bodyRef);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" fullscreen>
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.label', 2) }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div
                ref="bodyRef"
                class="max-h-[calc(98dvh-(2*var(--header-height)))] min-h-[calc(98dvh-(2*var(--header-height)))]"
            >
                <ClientOnly>
                    <VisSingleContainer
                        :data="stats?.count ?? []"
                        :padding="{ top: 2, left: 2, right: 2, bottom: 2 }"
                        :height="height - 64"
                        :width="width"
                    >
                        <VisDonut
                            :value="value"
                            :color="color"
                            :arc-width="30"
                            :central-label="`${$t('common.total_count')}: ${totalCount}`"
                        />
                        <VisBulletLegend :items="items" orientation="horizontal" />
                    </VisSingleContainer>
                </ClientOnly>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

<style scoped>
.unovis-single-container {
    --vis-nested-donut-central-label-text-color: rgb(var(--color-gray-400));
    --vis-legend-label-color: rgb(var(--color-gray-400));
    --vis-dark-legend-label-color: rgb(var(--color-gray-400));
    --vis-donut-central-label-text-color: rgb(var(--color-gray-400));
}

.dark {
    .unovis-single-container {
        --vis-nested-donut-central-label-text-color: rgb(var(--color-gray-200));
        --vis-legend-label-color: rgb(var(--color-gray-200));
        --vis-dark-legend-label-color: rgb(var(--color-gray-200));
        --vis-donut-central-label-text-color: rgb(var(--color-gray-200));
    }
}
</style>

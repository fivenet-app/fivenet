<script lang="ts" setup>
import UnitFeedItem from '~/components/centrum/units/UnitFeedItem.vue';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { ListUnitActivityResponse } from '~~/gen/ts/services/centrum/centrum';

const props = defineProps<{
    unitId: number;
}>();

const centrumCentrumClient = await getCentrumCentrumClient();

const offset = ref(0);

const { data, refresh } = useLazyAsyncData(`centrum-unit-${props.unitId}-activity-${offset.value}`, () => listUnitActivity());

async function listUnitActivity(): Promise<ListUnitActivityResponse> {
    try {
        const call = centrumCentrumClient.listUnitActivity({
            pagination: {
                offset: offset.value,
            },
            id: props.unitId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { pause, resume } = useIntervalFn(async () => {
    pause();
    await refresh();
    resume();
}, 3500);
</script>

<template>
    <div class="my-1 flex h-full flex-1 grow flex-col gap-2 px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold">{{ $t('common.feed') }}</h2>
        </div>

        <div class="flex flex-1 flex-col overflow-x-auto overflow-y-auto">
            <ul class="space-y-2" role="list">
                <UnitFeedItem
                    v-for="(activityItem, activityItemIdx) in data?.activity"
                    :key="activityItem.id"
                    :activity-length="data?.activity?.length ?? 0"
                    :item="activityItem"
                    :activity-item-idx="activityItemIdx"
                />
            </ul>
        </div>
    </div>
</template>

<script lang="ts" setup>
import List from '~/components/centrum/dispatches/List.vue';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';

useHead({
    title: 'common.dispatches',
});
definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
    showQuickButtons: false,
});

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`centrum-dispatches-${offset.value}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    return new Promise(async (res, rej) => {
        try {
            const req: ListDispatchesRequest = {
                pagination: {
                    offset: offset.value,
                },
                notStatus: [],
                status: [],
            };

            const call = $grpc.getCentrumClient().listDispatches(req);
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="w-full">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.dispatches')])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.dispatches')])" :retry="refresh" />
        <DataNoDataBlock v-else-if="data?.dispatches.length === 0" :type="$t('common.dispatches')" />

        <div v-else>
            <List :show-button="false" :hide-actions="true" :dispatches="data?.dispatches" />

            <TablePagination :pagination="data?.pagination" @offset-change="offset = $event" :refresh="refresh" />
        </div>
    </div>
</template>

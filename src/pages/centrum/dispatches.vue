<script lang="ts" setup>
import List from '~/components/centrum/dispatches/List.vue';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import TablePagination from '~/components/partials/elements/TablePagination.vue';

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
</script>

<template>
    <div class="w-full">
        <List :show-button="false" :hide-actions="true" :dispatches="data?.dispatches" />

        <TablePagination :pagination="data?.pagination" />
    </div>
</template>

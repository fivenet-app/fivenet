<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { TicketIcon } from 'mdi-vue3';
import type { ListDocumentActivityResponse } from '~~/gen/ts/services/docstore/docstore';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import ActivityListEntry from '~/components/documents/ActivityListEntry.vue';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`document-${props.documentId}-${offset.value}`, () =>
    listDocumentActivity(),
);

async function listDocumentActivity(): Promise<ListDocumentActivityResponse> {
    try {
        const call = $grpc.getDocStoreClient().listDocumentActivity({
            pagination: {
                offset: offset.value,
            },
            documentId: props.documentId,
            activityTypes: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="data === null || data.activity.length === 0"
            :icon="TicketIcon"
            :message="$t('common.not_found', [$t('common.activity')])"
        />
        <div v-else>
            <div class="sm:divide-y sm:divide-base-400 mb-1">
                <ActivityListEntry v-for="item in data.activity" :key="item.id" :entry="item" />
            </div>

            <TablePagination :pagination="data?.pagination" :refresh="refresh" @offset-change="offset = $event" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import DocumentActivityListEntry from '~/components/documents/activity/DocumentActivityListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { ListDocumentActivityResponse } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    documentId: number;
}>();

const { $grpc } = useNuxtApp();

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-${page.value}`, () => listDocumentActivity());

async function listDocumentActivity(): Promise<ListDocumentActivityResponse> {
    try {
        const call = $grpc.docstore.docStore.listDocumentActivity({
            pagination: {
                offset: offset.value,
            },
            documentId: props.documentId,
            activityTypes: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.document', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.document', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="!data || data.activity.length === 0"
            icon="i-mdi-ticket"
            :message="$t('common.not_found', [$t('common.activity')])"
        />

        <ul v-else class="mb-1 divide-y divide-gray-100 dark:divide-gray-800" role="list">
            <DocumentActivityListEntry v-for="item in data.activity" :key="item.id" :entry="item" />
        </ul>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>

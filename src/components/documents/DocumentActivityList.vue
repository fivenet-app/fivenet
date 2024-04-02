<script lang="ts" setup>
import { TicketIcon } from 'mdi-vue3';
import type { ListDocumentActivityResponse } from '~~/gen/ts/services/docstore/docstore';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DocumentActivityListEntry from '~/components/documents/DocumentActivityListEntry.vue';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`document-${props.documentId}-${page.value}`, () =>
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

        <template v-else>
            <div class="mb-1 sm:divide-y sm:divide-base-400">
                <DocumentActivityListEntry v-for="item in data.activity" :key="item.id" :entry="item" />
            </div>
        </template>

        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="data?.pagination?.pageSize ?? 0"
                :total="data?.pagination?.totalCount ?? 0"
            />
        </div>
    </div>
</template>

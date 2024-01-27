<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { FrequentlyAskedQuestionsIcon } from 'mdi-vue3';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/docstore/docstore';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DocumentRequestsListEntry from '~/components/documents/requests/DocumentRequestsListEntry.vue';
import { checkDocAccess } from '~/components/documents/helpers';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import type { Document } from '~~/gen/ts/resources/documents/documents';

const props = defineProps<{
    doc: Document;
    access: DocumentAccess;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const {
    data: requests,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.doc.id}-requests-${offset.value}`, () => listDocumnetReqs(props.doc.id));

async function listDocumnetReqs(documentId: string): Promise<ListDocumentReqsResponse> {
    try {
        const call = $grpc.getDocStoreClient().listDocumentReqs({
            pagination: {
                offset: offset.value,
            },
            documentId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canUpdate =
    can('DocStoreService.CreateDocumentReq') &&
    checkDocAccess(props.access, props.doc.creator, AccessLevel.EDIT, 'DocStoreService.CreateDocumentReq');
const canDelete =
    can('DocStoreService.DeleteDocumentReq') &&
    checkDocAccess(props.access, props.doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocumentReq');
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.request', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.request', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="requests === null || requests.requests.length === 0"
            :icon="FrequentlyAskedQuestionsIcon"
            :message="$t('common.not_found', [$t('common.request', 2)])"
        />

        <template v-else>
            <ul role="list" class="mb-6 divide-y divide-gray-100 rounded-md">
                <DocumentRequestsListEntry
                    v-for="request in requests.requests"
                    :key="request.id"
                    :request="request"
                    :can-update="canUpdate"
                    :can-delete="canDelete"
                    @refresh="refresh()"
                />
            </ul>
        </template>
    </div>
</template>

<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { TrashCanIcon } from 'mdi-vue3';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const {
    data: requests,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-requests-${offset.value}`, () => listDocumnetReqs(props.documentId));

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
</script>

<template>
    <ul v-if="requests !== null" class="list-disc">
        <li v-for="item in requests.requests" :key="item.id">
            {{ item.id }} - {{ DocActivityType[item.requestType] }} (Reason: {{ item.reason }})
            <button v-if="can('DocStoreService.DeleteDocumentReq')" type="button">
                <TrashCanIcon class="w-6 h-6" />
            </button>
        </li>
    </ul>
</template>

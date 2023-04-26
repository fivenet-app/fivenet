<script lang="ts" setup>
import { ref } from 'vue';
import { DocumentRelation } from '@fivenet/gen/resources/documents/documents_pb';
import { GetUserDocumentsRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { PaginationRequest } from '@fivenet/gen/resources/common/database/database_pb';
import DocumentRelations from '~/components/documents/DocumentRelations.vue';
import { DocumentTextIcon } from '@heroicons/vue/24/outline';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

const { data: relations, pending, refresh, error } = await useLazyAsyncData(`citizeninfo-documents-${props.userId}`, () => getUserDocuments());

const offset = ref(0);

async function getUserDocuments(): Promise<Array<DocumentRelation>> {
    return new Promise(async (res, rej) => {
        const req = new GetUserDocumentsRequest();
        req.setPagination((new PaginationRequest()).setOffset(offset.value))
        req.setUserId(props.userId);

        try {
            const resp = await $grpc.getDocStoreClient().
                getUserDocuments(req, null);

            return res(resp.getRelationsList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return;
        }
    });
}
</script>

<template>
    <div class="mt-2">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.user',1 )} ${$t('common.document', 2)}`])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [`${$t('common.user', 1)} ${$t('common.document', 2)}`])" :retry="refresh" />
        <button v-else-if="relations && relations.length == 0" type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled>
            <DocumentTextIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                {{ $t('common.not_found', [`${$t('common.user', 1)} ${$t('common.document', 2)}`]) }}
            </span>
        </button>
        <DocumentRelations v-else-if="relations" :relations="relations" />
    </div>
</template>

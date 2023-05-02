<script lang="ts" setup>
import { ref } from 'vue';
import { DocumentRelation } from '@fivenet/gen/resources/documents/documents_pb';
import { ListUserDocumentsRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { PaginationRequest } from '@fivenet/gen/resources/common/database/database_pb';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { DOC_RELATION_Util } from '@fivenet/gen/resources/documents/documents.pb_enums';
import { DocumentTextIcon, ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const { $grpc } = useNuxtApp();

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

const offset = ref(0);

const { data: relations, pending, refresh, error } = useLazyAsyncData(`user-${props.userId}-documents-${offset.value}`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<Array<DocumentRelation>> {
    return new Promise(async (res, rej) => {
        const req = new ListUserDocumentsRequest();
        req.setPagination((new PaginationRequest()).setOffset(offset.value))
        req.setUserId(props.userId);

        try {
            const resp = await $grpc.getDocStoreClient().
                listUserDocuments(req, null);

            return res(resp.getRelationsList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="mt-2">
        <DataPendingBlock v-if="pending"
            :message="$t('common.loading', [`${$t('common.user', 1)} ${$t('common.document', 2)}`])" />
        <DataErrorBlock v-else-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.user', 1)} ${$t('common.document', 2)}`])"
            :retry="refresh" />
        <button v-else-if="relations && relations.length == 0" type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled>
            <DocumentTextIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                {{ $t('common.not_found', [`${$t('common.document', 1)} ${$t('common.relation', 2)}`]) }}
            </span>
        </button>
        <div v-if="relations">
            <!-- Relations list (smallest breakpoint only) -->
            <div v-if="relations.length > 0" class="sm:hidden text-neutral">
                <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="relation in relations" :key="relation.getId()">
                        <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowsRightLeftIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                                    <span class="flex flex-col text-sm truncate">
                                        <span>
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: relation.getDocumentId() } }">
                                                {{ relation.getDocument()?.getTitle() }}<span
                                                    v-if="relation.getDocument()?.getCategory()"> (Category: {{
                                                        relation.getDocument()?.getCategory()?.getName() }})</span>
                                            </NuxtLink>
                                        </span>
                                        <span>
                                            <NuxtLink
                                                :to="{ name: 'citizens-id', params: { id: relation.getTargetUserId() } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ relation.getTargetUser()?.getFirstname() + ", " +
                                                    relation.getTargetUser()?.getLastname() }}
                                            </NuxtLink>
                                        </span>
                                        <span class="font-medium ">{{
                                            DOC_RELATION_Util.toEnumKey(relation.getRelation()) }}</span>
                                        <span class="truncate">{{ relation.getSourceUser()?.getFirstname()
                                            +
                                            ", " +
                                            relation.getSourceUser()?.getLastname() }}</span>
                                        <time datetime="">{{ toDateLocaleString(relation.getCreatedAt()) }}</time>
                                    </span>
                                </span>
                                <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                            </span>
                        </a>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div v-if="relations.length > 0" class="hidden sm:block">
                <div>
                    <div class="flex flex-col mt-2">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <table class="min-w-full bg-base-700 text-neutral">
                                <thead>
                                    <tr>
                                        <th class="px-6 py-3 text-sm font-semibold text-left " scope="col">
                                            {{ $t('common.document', 1) }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-left " scope="col">
                                            {{ $t('common.target') }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-right " scope="col">
                                            {{ $t('common.relation', 1) }}
                                        </th>
                                        <th class="hidden px-6 py-3 text-sm font-semibold text-left md:block" scope="col">
                                            {{ $t('common.creator') }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-right " scope="col">
                                            {{ $t('common.date') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-600 bg-base-800 text-neutral">
                                    <tr v-for="relation in relations" :key="relation.getId()">
                                        <td class="px-6 py-4 text-sm ">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: relation.getDocumentId() } }">
                                                {{ relation.getDocument()?.getTitle() }}<span
                                                    v-if="relation.getDocument()?.getCategory()"> ({{ $t('common.category',
                                                        1)
                                                    }}: {{
    relation.getDocument()?.getCategory()?.getName() }})</span>
                                            </NuxtLink>
                                        </td>
                                        <td class="px-6 py-4 text-sm ">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{ name: 'citizens-id', params: { id: relation.getTargetUserId() } }"
                                                    class="inline-flex space-x-2 text-sm truncate group">
                                                    {{ relation.getTargetUser()?.getFirstname() + ", " +
                                                        relation.getTargetUser()?.getLastname() }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                            <span class="font-medium ">{{
                                                DOC_RELATION_Util.toEnumKey(relation.getRelation()) }}</span>
                                        </td>
                                        <td class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{ name: 'citizens-id', params: { id: relation.getSourceUserId() } }"
                                                    class="inline-flex space-x-2 text-sm truncate group">
                                                    {{ relation.getSourceUser()?.getFirstname() + ", " +
                                                        relation.getSourceUser()?.getLastname() }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                            <time datetime="">{{ toDateLocaleString(relation.getCreatedAt()) }}</time>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

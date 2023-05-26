<script lang="ts" setup>
import { ref } from 'vue';
import { DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { DocumentTextIcon, ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';
import { LockClosedIcon, LockOpenIcon } from '@heroicons/vue/20/solid';
import { RpcError } from 'grpc-web';

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
        try {
            const call = $grpc.getDocStoreClient().
                listUserDocuments({
                    pagination: {
                        offset: offset.value,
                    },
                    userId: props.userId,
                    relations: [],
                });
            const { response } = await call;

            return res(response.relations);
        } catch (e) {
            $grpc.handleError(e as RpcError);
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
        <button v-else-if="relations && relations.length === 0" type="button"
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
                    <li v-for="relation in relations" :key="relation.id">
                        <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowsRightLeftIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                                    <span class="flex flex-col text-sm truncate">
                                        <span>
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: relation.documentId } }">
                                                {{ relation.document?.title }}<span
                                                    v-if="relation.document?.category"> (Category: {{
                                                        relation.document?.category?.name }})</span>
                                            </NuxtLink>
                                        </span>
                                        <span>
                                            <div v-if="relation.document?.closed"
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                                                <LockClosedIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                                <span class="text-sm font-medium text-error-700">
                                                    {{ $t('common.close', 2) }}
                                                </span>
                                            </div>
                                            <div v-else
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                                <LockOpenIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                                <span class="text-sm font-medium text-green-700">
                                                    {{ $t('common.open') }}
                                                </span>
                                            </div>
                                        </span>
                                        <span>
                                            <NuxtLink :to="{ name: 'citizens-id', params: { id: relation.targetUserId } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ relation.targetUser?.firstname + ", " +
                                                    relation.targetUser?.lastname }}
                                            </NuxtLink>
                                        </span>
                                        <span class="font-medium">
                                            {{
                                                $t(`enums.docstore.DOC_RELATION.${relation.relation}`)
                                            }}
                                        </span>
                                        <span class="truncate">
                                            {{ relation.sourceUser?.firstname }},
                                            {{ relation.sourceUser?.lastname }}
                                        </span>
                                        <time :datetime="toDateLocaleString(relation.createdAt, $d)">
                                            {{ useLocaleTimeAgo(toDate(relation.createdAt)!).value }}
                                        </time>
                                    </span>
                                </span>
                                <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                            </span>
                        </a>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div v-if="relations && relations.length > 0" class="hidden sm:block">
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
                                            {{ $t('common.close', 2) }}
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
                                    <tr v-for="relation in relations" :key="relation.id">
                                        <td class="px-6 py-4 text-sm">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: relation.documentId } }">
                                                {{ relation.document?.title }}<span
                                                    v-if="relation.document?.category"> ({{ $t('common.category',
                                                        1)
                                                    }}: {{ relation.document?.category?.name }})
                                                </span>
                                            </NuxtLink>
                                        </td>
                                        <td class="px-6 py-4 text-sm">
                                            <div v-if="relation.document?.closed"
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                                                <LockClosedIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                                <span class="text-sm font-medium text-error-700">
                                                    {{ $t('common.close', 2) }}
                                                </span>
                                            </div>
                                            <div v-else
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                                                <LockOpenIcon class="w-5 h-5 text-green-500" aria-hidden="true" />
                                                <span class="text-sm font-medium text-green-700">
                                                    {{ $t('common.open') }}
                                                </span>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{ name: 'citizens-id', params: { id: relation.targetUserId } }"
                                                    class="inline-flex space-x-2 text-sm truncate group">
                                                    {{ relation.targetUser?.firstname + ", " +
                                                        relation.targetUser?.lastname }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <span class="font-medium">
                                                {{
                                                    $t(`enums.docstore.DOC_RELATION.${relation.relation}`)
                                                }}
                                            </span>
                                        </td>
                                        <td class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{ name: 'citizens-id', params: { id: relation.sourceUserId } }"
                                                    class="inline-flex space-x-2 text-sm truncate group">
                                                    {{ relation.sourceUser?.firstname + ", " +
                                                        relation.sourceUser?.lastname }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <time :datetime="$d(relation.createdAt?.timestamp?.toDate()!, 'short')">
                                                {{ $d(toDate(relation.createdAt)!) }}
                                            </time>
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

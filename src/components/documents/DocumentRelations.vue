<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountMultiple, mdiArrowCollapse, mdiChevronRight } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { DOC_RELATION, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import DataNoDataBlock from '../partials/DataNoDataBlock.vue';

const { $grpc } = useNuxtApp();

const props = withDefaults(
    defineProps<{
        documentId: bigint;
        showDocument?: boolean;
        showSource?: boolean;
    }>(),
    {
        showDocument: true,
        showSource: true,
    }
);

const {
    data: relations,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-relations`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<Array<DocumentRelation>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getDocumentRelations({
                documentId: props.documentId,
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
    <div>
        <DataNoDataBlock
            v-if="relations && relations.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.relation', 2)}`"
            :icon="mdiAccountMultiple"
        />
        <!-- Relations list (smallest breakpoint only) -->
        <div v-if="relations && relations.length > 0" class="sm:hidden text-neutral">
            <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                <li v-for="relation in relations" :key="relation.id?.toString()">
                    <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <SvgIcon
                                    class="flex-shrink-0 w-5 h-5 text-gray-400"
                                    aria-hidden="true"
                                    type="mdi"
                                    :path="mdiArrowCollapse"
                                />
                                <span class="flex flex-col text-sm truncate">
                                    <span v-if="showDocument">
                                        <NuxtLink
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: relation.documentId.toString(),
                                                },
                                            }"
                                        >
                                            {{ relation.document?.title
                                            }}<span v-if="relation.document?.category">
                                                (Category:
                                                {{ relation.document?.category?.name }})
                                            </span>
                                        </NuxtLink>
                                    </span>
                                    <span>
                                        <NuxtLink
                                            :to="{
                                                name: 'citizens-id',
                                                params: {
                                                    id: relation.targetUserId,
                                                },
                                            }"
                                            class="inline-flex space-x-2 text-sm truncate group"
                                        >
                                            {{ relation.targetUser?.firstname + ', ' + relation.targetUser?.lastname }}
                                        </NuxtLink>
                                    </span>
                                    <span class="font-medium">
                                        {{ $t(`enums.docstore.DOC_RELATION.${DOC_RELATION[relation.relation]}`) }}
                                    </span>
                                    <span v-if="showSource" class="truncate"
                                        >{{ relation.sourceUser?.firstname + ', ' + relation.sourceUser?.lastname }}
                                    </span>
                                    <time :datetime="$d(toDate(relation.createdAt)!, 'short')">
                                        {{ $d(toDate(relation.createdAt)!, 'short') }}
                                    </time>
                                </span>
                            </span>
                            <SvgIcon
                                class="flex-shrink-0 w-5 h-5 text-base-200"
                                aria-hidden="true"
                                type="mdi"
                                :path="mdiChevronRight"
                            />
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
                                    <th v-if="showDocument" class="px-6 py-3 text-sm font-semibold text-left" scope="col">
                                        {{ $t('common.document', 1) }}
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-left" scope="col">
                                        {{ $t('common.target') }}
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                        {{ $t('common.relation', 1) }}
                                    </th>
                                    <th
                                        v-if="showSource"
                                        class="hidden px-6 py-3 text-sm font-semibold text-left md:block"
                                        scope="col"
                                    >
                                        {{ $t('common.creator') }}
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                        {{ $t('common.date') }}
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-600 bg-base-800 text-neutral">
                                <tr v-for="relation in relations" :key="relation.id?.toString()">
                                    <td v-if="showDocument" class="px-6 py-4 text-sm">
                                        <NuxtLink
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: relation.documentId.toString(),
                                                },
                                            }"
                                        >
                                            {{ relation.document?.title
                                            }}<span v-if="relation.document?.category">
                                                ({{ $t('common.category', 1) }}: {{ relation.document?.category?.name }})</span
                                            >
                                        </NuxtLink>
                                    </td>
                                    <td class="px-6 py-4 text-sm">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{
                                                    name: 'citizens-id',
                                                    params: {
                                                        id: relation.targetUserId,
                                                    },
                                                }"
                                                class="inline-flex space-x-2 text-sm truncate group"
                                            >
                                                {{ relation.targetUser?.firstname + ', ' + relation.targetUser?.lastname }}
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                        <span class="font-medium">
                                            {{ $t(`enums.docstore.DOC_RELATION.${DOC_RELATION[relation.relation]}`) }}
                                        </span>
                                    </td>
                                    <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{
                                                    name: 'citizens-id',
                                                    params: {
                                                        id: relation.sourceUserId,
                                                    },
                                                }"
                                                class="inline-flex space-x-2 text-sm truncate group"
                                            >
                                                {{ relation.sourceUser?.firstname + ', ' + relation.sourceUser?.lastname }}
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                        <time :datetime="$d(toDate(relation.createdAt)!, 'short')">
                                            {{ $d(toDate(relation.createdAt)!, 'short') }}
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
</template>

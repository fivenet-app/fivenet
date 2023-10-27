<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { AccountMultipleIcon, ArrowCollapseIcon, ChevronRightIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import { DocRelation, DocumentRelation } from '~~/gen/ts/resources/documents/documents';

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
    },
);

const {
    data: relations,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-relations`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<DocumentRelation[]> {
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
            :icon="AccountMultipleIcon"
        />

        <template v-if="relations && relations.length > 0">
            <!-- Relations list (smallest breakpoint only) -->
            <div class="sm:hidden text-neutral">
                <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="relation in relations" :key="relation.id?.toString()">
                        <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowCollapseIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
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
                                                <span
                                                    v-if="relation.document?.category"
                                                    class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30 mr-1"
                                                >
                                                    {{ relation.document?.category?.name }}
                                                </span>
                                                {{ relation.document?.title }}
                                            </NuxtLink>
                                        </span>
                                        <span>
                                            <CitizenInfoPopover :user="relation.targetUser" />
                                        </span>
                                        <span class="font-medium">
                                            {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                        </span>
                                        <span v-if="showSource" class="truncate">
                                            <CitizenInfoPopover :user="relation.sourceUser" />
                                        </span>
                                        <Time :value="relation.createdAt" />
                                    </span>
                                </span>
                                <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                            </span>
                        </a>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div class="hidden sm:block">
                <div>
                    <div class="flex flex-col mt-2">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <table class="min-w-full bg-base-600 text-neutral divide-y divide-base-600">
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
                                <tbody class="divide-y divide-y divide-gray-600 bg-base-700 text-neutral">
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
                                                <span
                                                    v-if="relation.document?.category"
                                                    class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30 mr-1"
                                                >
                                                    {{ relation.document?.category?.name }}
                                                </span>
                                                {{ relation.document?.title }}
                                            </NuxtLink>
                                        </td>
                                        <td class="px-6 py-4 text-sm">
                                            <div class="flex">
                                                <CitizenInfoPopover :user="relation.targetUser" />
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <span class="font-medium">
                                                {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                            </span>
                                        </td>
                                        <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                            <div class="flex">
                                                <CitizenInfoPopover :user="relation.sourceUser" />
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <Time :value="relation.createdAt" />
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>

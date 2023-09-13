<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { ArrowCollapseIcon, ChevronRightIcon, FileDocumentMultipleIcon } from 'mdi-vue3';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import { DOC_REFERENCE, DocumentReference } from '~~/gen/ts/resources/documents/documents';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';

const { $grpc } = useNuxtApp();

const props = withDefaults(
    defineProps<{
        documentId: bigint;
        showSource?: boolean;
    }>(),
    {
        showSource: true,
    },
);

const {
    data: references,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references`, () => getDocumentReferences());

async function getDocumentReferences(): Promise<DocumentReference[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getDocumentReferences({
                documentId: props.documentId,
            });
            const { response } = await call;

            return res(response.references);
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
            v-if="references && references.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.reference', 2)}`"
            :mdi="FileDocumentMultipleIcon"
        />

        <template v-if="references && references.length > 0">
            <!-- Relations list (smallest breakpoint only) -->
            <div class="sm:hidden text-neutral">
                <ul role="list" class="mt-2 divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="reference in references" :key="reference.id?.toString()">
                        <NuxtLink
                            :to="{
                                name: 'documents-id',
                                params: { id: reference.targetDocumentId.toString() },
                            }"
                            class="block px-4 py-4 bg-base-800 hover:bg-base-700"
                        >
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowCollapseIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                                    <span class="flex flex-col text-sm truncate">
                                        <span>
                                            {{ reference.targetDocument?.title
                                            }}<span v-if="reference.targetDocument?.category"
                                                >&nbsp;({{ $t('common.category', 1) }}:
                                                {{ reference.targetDocument?.category?.name }})</span
                                            >
                                        </span>
                                        <span class="font-medium">
                                            {{ DOC_REFERENCE[reference.reference] }}
                                        </span>
                                        <span v-if="showSource" class="truncate">
                                            {{ reference.sourceDocument?.title
                                            }}<span v-if="reference.sourceDocument?.category">
                                                ({{ $t('common.category', 1) }}:
                                                {{ reference.sourceDocument?.category?.name }})</span
                                            >
                                        </span>
                                        <span>
                                            <CitizenInfoPopover :user="reference.creator" />
                                        </span>
                                        <Time :value="reference.createdAt" />
                                    </span>
                                </span>
                                <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                            </span>
                        </NuxtLink>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div class="hidden sm:block">
                <div>
                    <div class="flex flex-col mt-2">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <table class="min-w-full bg-base-700 text-neutral">
                                <thead>
                                    <tr>
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
                                            {{ $t('common.source') }}
                                        </th>
                                        <th class="hidden px-6 py-3 text-sm font-semibold text-left md:block" scope="col">
                                            {{ $t('common.creator') }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                            {{ $t('common.date') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-600 bg-base-800 text-neutral">
                                    <tr v-for="reference in references" :key="reference.id?.toString()">
                                        <td class="px-6 py-4 text-sm">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'documents-id',
                                                        params: {
                                                            id: reference.targetDocumentId.toString(),
                                                        },
                                                    }"
                                                    class="inline-flex space-x-2 text-sm truncate group"
                                                >
                                                    {{ reference.targetDocument?.title
                                                    }}<span v-if="reference.targetDocument?.category"
                                                        >&nbsp;({{ $t('common.category', 1) }}:
                                                        {{ reference.targetDocument?.category?.name }})</span
                                                    >
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <span class="font-medium">
                                                {{ $t(`enums.docstore.DOC_REFERENCE.${DOC_REFERENCE[reference.reference]}`) }}
                                            </span>
                                        </td>
                                        <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'documents-id',
                                                        params: {
                                                            id: reference.sourceDocumentId.toString(),
                                                        },
                                                    }"
                                                    class="inline-flex space-x-1 text-sm truncate group"
                                                >
                                                    {{ reference.sourceDocument?.title
                                                    }}<span v-if="reference.sourceDocument?.category"
                                                        >&nbsp;({{ $t('common.category', 1) }}:
                                                        {{ reference.sourceDocument?.category?.name }})</span
                                                    >
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <div class="flex">
                                                <NuxtLink :to="{ name: 'citizens-id', params: { id: reference.creatorId! } }">
                                                    {{ reference.creator?.firstname }},
                                                    {{ reference.creator?.lastname }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <Time :value="reference.createdAt" />
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

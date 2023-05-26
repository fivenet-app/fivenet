<script lang="ts" setup>
import { DOC_REFERENCE_Util } from '~~/gen/ts/resources/documents/documents.pb_enums';
import { DocumentReference } from '~~/gen/ts/resources/documents/documents';
import { GetDocumentReferencesRequest } from '~~/gen/ts/services/docstore/docstore';
import { ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const { $grpc } = useNuxtApp();

const props = defineProps({
    documentId: {
        required: true,
        type: Number,
    },
    showSource: {
        required: false,
        type: Boolean,
        default: true,
    },
});

const { data: references, pending, refresh, error } = useLazyAsyncData(`document-${props.documentId}-references`, () => getDocumentReferences());

async function getDocumentReferences(): Promise<Array<DocumentReference>> {
    return new Promise(async (res, rej) => {
        const req = new GetDocumentReferencesRequest();
        req.setDocumentId(props.documentId);

        try {
            const resp = await $grpc.getDocStoreClient().
                getDocumentReferences(req);

            return res(resp.getReferencesList());
        } catch (e) {
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div>
        <span v-if="references && references.length === 0" class="text-neutral">{{ $t('common.not_found',
            [`${$t('common.document', 1)}
                    ${$t('common.reference', 2)}`]) }}</span>
        <!-- Relations list (smallest breakpoint only) -->
        <div v-if="references && references.length > 0" class="sm:hidden text-neutral">
            <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                <li v-for="reference in references" :key="reference.id">
                    <NuxtLink :to="{ name: 'documents-id', params: { id: reference.getTardocumentId } }"
                        class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <ArrowsRightLeftIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                                <span class="flex flex-col text-sm truncate">
                                    <span>
                                        {{ reference.getTardocument?.title }}<span
                                            v-if="reference.getTardocument?.category">&nbsp;({{
                                                $t('common.category', 1) }}: {{
        reference.getTardocument?.category?.name }})</span>
                                    </span>
                                    <span class="font-medium">
                                        {{ DOC_REFERENCE_Util.toEnumKey(reference.getReference()) }}
                                    </span>
                                    <span v-if="showSource" class="truncate">
                                        {{ reference.getSourceDocument()?.title }}<span
                                            v-if="reference.getSourceDocument()?.category"> ({{ $t('common.category',
                                                1) }}: {{
        reference.getSourceDocument()?.category?.name }})</span>
                                    </span>
                                    <span>
                                        <NuxtLink :to="{ name: 'citizens-id', params: { id: reference.getCreatorId() } }">
                                            {{ reference.getCreator()?.firstname }}, {{
                                                reference.getCreator()?.lastname }}
                                        </NuxtLink>
                                    </span>
                                    <time datetime="">{{ $d(reference.createdAt?.timestamp?.toDate()!, 'short')
                                    }}</time>
                                </span>
                            </span>
                            <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </NuxtLink>
                </li>
            </ul>
        </div>

        <!-- Relations table (small breakpoint and up) -->
        <div v-if="references && references.length > 0" class="hidden sm:block">
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
                                    <th v-if="showSource" class="hidden px-6 py-3 text-sm font-semibold text-left md:block"
                                        scope="col">
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
                                <tr v-for="reference in references" :key="reference.id">
                                    <td class="px-6 py-4 text-sm">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: reference.getTardocumentId } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ reference.getTardocument?.title }}<span
                                                    v-if="reference.getTardocument?.category">&nbsp;({{
                                                        $t('common.category', 1) }}: {{
        reference.getTardocument?.category?.name }})</span>
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                        <span class="font-medium">
                                            {{
                                                $t(`enums.docstore.DOC_REFERENCE.${DOC_REFERENCE_Util.toEnumKey(reference.getReference())!}`)
                                            }}
                                        </span>
                                    </td>
                                    <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: reference.getSourceDocumentId() } }"
                                                class="inline-flex space-x-1 text-sm truncate group">
                                                {{ reference.getSourceDocument()?.title }}<span
                                                    v-if="reference.getSourceDocument()?.category">&nbsp;({{
                                                        $t('common.category', 1) }}: {{
        reference.getSourceDocument()?.category?.name }})</span>
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{ name: 'citizens-id', params: { id: reference.getCreatorId() } }">
                                                {{ reference.getCreator()?.firstname }}, {{
                                                    reference.getCreator()?.lastname }}
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                        <time datetime="">{{ $d(reference.createdAt?.timestamp?.toDate()!,
                                            'short') }}</time>
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

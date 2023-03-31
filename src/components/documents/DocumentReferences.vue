<script lang="ts" setup>
import { DOC_REFERENCE_Util } from '@fivenet/gen/resources/documents/documents.pb_enums';
import { DocumentReference } from '@fivenet/gen/resources/documents/documents_pb';
import { ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';
import { toDateLocaleString } from '../../utils/time';

defineProps({
    references: {
        required: true,
        type: Array<DocumentReference>,
    },
    showSource: {
        required: false,
        type: Boolean,
        default: true,
    },
});
</script>

<template>
    <div>
        <span v-if="references.length == 0" class="text-neutral">No Document References found.</span>
        <!-- Relations list (smallest breakpoint only) -->
        <div v-if="references.length > 0" class="sm:hidden text-neutral">
            <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                <li v-for="reference in references" :key="reference.getId()">
                    <NuxtLink :to="{ name: 'documents-id', params: { id: reference.getTargetDocumentId() } }"
                        class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <ArrowsRightLeftIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                                <span class="flex flex-col text-sm truncate">
                                    <span>
                                        {{ reference.getTargetDocument()?.getTitle() }}<span
                                                    v-if="reference.getTargetDocument()?.getCategory()">&nbsp;(Category: {{
                                                        reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                    </span>
                                    <span class="font-medium ">{{
                                        DOC_REFERENCE_Util.toEnumKey(reference.getReference()) }}</span>
                                    <span v-if="showSource" class="truncate">
                                        {{ reference.getSourceDocument()?.getTitle() }}<span
                                            v-if="reference.getSourceDocument()?.getCategory()"> (Category: {{
                                                reference.getSourceDocument()?.getCategory()?.getName() }})</span>
                                    </span>
                                    <span>
                                        <NuxtLink :to="{ name: 'citizens-id', params: { id: reference.getCreatorId() }} ">
                                            {{ reference.getCreator()?.getFirstname() }}, {{ reference.getCreator()?.getLastname() }}
                                        </NuxtLink>
                                    </span>
                                    <time datetime="">{{ toDateLocaleString(reference.getCreatedAt()) }}</time>
                                </span>
                            </span>
                            <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                        </span>
                    </NuxtLink>
                </li>
            </ul>
        </div>

        <!-- Relations table (small breakpoint and up) -->
        <div v-if="references.length > 0" class="hidden sm:block">
            <div>
                <div class="flex flex-col mt-2">
                    <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                        <table class="min-w-full bg-base-700 text-neutral">
                            <thead>
                                <tr>
                                    <th class="px-6 py-3 text-sm font-semibold text-left" scope="col">
                                        Target
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                        Relation
                                    </th>
                                    <th v-if="showSource" class="hidden px-6 py-3 text-sm font-semibold text-left md:block" scope="col">
                                        Source
                                    </th>
                                    <th class="hidden px-6 py-3 text-sm font-semibold text-left md:block" scope="col">
                                        Creator
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                        Date
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-base-600 bg-base-800 text-neutral">
                                <tr v-for="reference in references" :key="reference.getId()">
                                    <td class="px-6 py-4 text-sm ">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: reference.getTargetDocumentId() } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ reference.getTargetDocument()?.getTitle() }}<span
                                                    v-if="reference.getTargetDocument()?.getCategory()">&nbsp;(Category: {{
                                                        reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                        <span class="font-medium ">{{
                                            DOC_REFERENCE_Util.toEnumKey(reference.getReference()) }}</span>
                                    </td>
                                    <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                        <div class="flex">
                                            <NuxtLink
                                                :to="{ name: 'documents-id', params: { id: reference.getSourceDocumentId() } }"
                                                class="inline-flex space-x-1 text-sm truncate group">
                                                {{ reference.getSourceDocument()?.getTitle() }}<span
                                                    v-if="reference.getSourceDocument()?.getCategory()">&nbsp;(Category: {{
                                                        reference.getSourceDocument()?.getCategory()?.getName() }})</span>
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                        <div class="flex">
                                            <NuxtLink :to="{ name: 'citizens-id', params: { id: reference.getCreatorId() }} ">
                                                {{ reference.getCreator()?.getFirstname() }}, {{ reference.getCreator()?.getLastname() }}
                                            </NuxtLink>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                        <time datetime="">{{ toDateLocaleString(reference.getCreatedAt()) }}</time>
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

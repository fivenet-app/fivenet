<script lang="ts" setup>
import { DOC_REFERENCE_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';
import { DocumentReference } from '@arpanet/gen/resources/documents/documents_pb';
import { ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';
import { getDateLocaleString } from '../../utils/time';

defineProps({
    references: {
        required: true,
        type: Array<DocumentReference>,
    }
});
</script>

<template>
    <div>
        <span v-if="references.length == 0" class="text-neutral">No Document References found.</span>
        <!-- Relations list (smallest breakpoint only) -->
        <div v-if="references.length > 0" class="shadow sm:hidden">
            <ul role="list" class="mt-2 divide-y divide-gray-200 overflow-hidden shadow sm:hidden">
                <li v-for="reference in references" :key="reference.getId()">
                    <a href="#" class="block bg-white px-4 py-4 hover:bg-gray-50">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <ArrowsRightLeftIcon class="h-5 w-5 flex-shrink-0 text-gray-400" aria-hidden="true" />
                                <span class="flex flex-col truncate text-sm ">
                                    <span>
                                        {{ reference.getTargetDocument()?.getTitle() }}<span
                                            v-if="reference.getTargetDocument()?.getCategory()"> (Category: {{
                                                reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                    </span>
                                    <span class="font-medium ">{{
                                        DOC_REFERENCE_Util.toEnumKey(reference.getReference()) }}</span>
                                    <span class="truncate">
                                        {{ reference.getSourceDocument()?.getTitle() }}<span
                                            v-if="reference.getTargetDocument()?.getCategory()"> (Category: {{
                                                reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                    </span>
                                    <time datetime="">{{ getDateLocaleString(reference.getCreatedAt()) }}</time>
                                </span>
                            </span>
                            <ChevronRightIcon class="h-5 w-5 flex-shrink-0 text-gray-400" aria-hidden="true" />
                        </span>
                    </a>
                </li>
            </ul>
        </div>

        <!-- Relations table (small breakpoint and up) -->
        <div v-if="references.length > 0" class="hidden sm:block">
            <div>
                <div class="mt-2 flex flex-col">
                    <div class="min-w-full overflow-hidden overflow-x-auto align-middle shadow sm:rounded-lg">
                        <table class="min-w-full bg-base-700 text-neutral">
                            <thead>
                                <tr>
                                    <th class="px-6 py-3 text-left text-sm font-semibold"
                                        scope="col">
                                        Target</th>
                                    <th class="px-6 py-3 text-right text-sm font-semibold"
                                        scope="col">
                                        Relation</th>
                                    <th class="hidden px-6 py-3 text-left text-sm font-semibold md:block"
                                        scope="col">Source</th>
                                    <th class="px-6 py-3 text-right text-sm font-semibold"
                                        scope="col">
                                        Date</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-base-600 bg-base-800 text-neutral">
                                <tr v-for="reference in references" :key="reference.getId()">
                                    <td class="px-6 py-4 text-sm ">
                                        <div class="flex">
                                            <router-link
                                                :to="{ name: 'Documents: Info', params: { id: reference.getSourceDocumentId() } }"
                                                class="group inline-flex space-x-2 truncate text-sm">
                                                {{ reference.getSourceDocument()?.getTitle() }}<span
                                                    v-if="reference.getTargetDocument()?.getCategory()"> (Category: {{
                                                        reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                            </router-link>
                                        </div>
                                    </td>
                                    <td class="whitespace-nowrap px-6 py-4 text-right text-sm ">
                                        <span class="font-medium ">{{
                                            DOC_REFERENCE_Util.toEnumKey(reference.getReference()) }}</span>
                                    </td>
                                    <td class="hidden whitespace-nowrap px-6 py-4 text-sm  md:block">
                                        <div class="flex">
                                            <router-link
                                                :to="{ name: 'Documents: Info', params: { id: reference.getTargetDocumentId() } }"
                                                class="group inline-flex space-x-1 truncate text-sm">
                                                {{ reference.getTargetDocument()?.getTitle() }}<span
                                                    v-if="reference.getTargetDocument()?.getCategory()"> (Category: {{
                                                        reference.getTargetDocument()?.getCategory()?.getName() }})</span>
                                            </router-link>
                                        </div>
                                    </td>
                                    <td class="whitespace-nowrap px-6 py-4 text-right text-sm ">
                                        <time datetime="">{{ getDateLocaleString(reference.getCreatedAt()) }}</time>
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

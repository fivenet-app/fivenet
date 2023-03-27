<script lang="ts" setup>
import { DOC_RELATION_Util } from '@arpanet/gen/resources/documents/documents.pb_enums';
import { DocumentRelation } from '@arpanet/gen/resources/documents/documents_pb';
import { ArrowsRightLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';
import { toDateLocaleString } from '../../utils/time';

defineProps({
    relations: {
        required: true,
        type: Array<DocumentRelation>,
    },
    showDocument: {
        required: false,
        type: Boolean,
        default: true,
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
        <span v-if="relations.length == 0" class="text-neutral">No Document Relations found.</span>
        <!-- Relations list (smallest breakpoint only) -->
        <div v-if="relations.length > 0" class="sm:hidden text-neutral">
            <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                <li v-for="relation in relations" :key="relation.getId()">
                    <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <ArrowsRightLeftIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                                <span class="flex flex-col text-sm truncate">
                                    <span v-if="showDocument">
                                        <router-link
                                            :to="{ name: 'Documents: Info', params: { id: relation.getDocumentId() } }">
                                            {{ relation.getDocument()?.getTitle() }}<span
                                                v-if="relation.getDocument()?.getCategory()"> (Category: {{
                                                    relation.getDocument()?.getCategory()?.getName() }})</span>
                                        </router-link>
                                    </span>
                                    <span>
                                        <router-link
                                            :to="{ name: 'Citizens: Info', params: { id: relation.getTargetUserId() } }"
                                            class="inline-flex space-x-2 text-sm truncate group">
                                            {{ relation.getTargetUser()?.getFirstname() + ", " +
                                                relation.getTargetUser()?.getLastname() }}
                                        </router-link>
                                    </span>
                                    <span class="font-medium ">{{
                                        DOC_RELATION_Util.toEnumKey(relation.getRelation()) }}</span>
                                    <span v-if="showSource" class="truncate">{{ relation.getSourceUser()?.getFirstname() +
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
                                    <th v-if="showDocument" class="px-6 py-3 text-sm font-semibold text-left " scope="col">
                                        Document
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-left " scope="col">
                                        Target
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right " scope="col">
                                        Relation
                                    </th>
                                    <th v-if="showSource" class="hidden px-6 py-3 text-sm font-semibold text-left md:block"
                                        scope="col">
                                        Creator
                                    </th>
                                    <th class="px-6 py-3 text-sm font-semibold text-right " scope="col">
                                        Date
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-600 bg-base-800 text-neutral">
                                <tr v-for="relation in relations" :key="relation.getId()">
                                    <td v-if="showDocument" class="px-6 py-4 text-sm ">
                                        <router-link
                                            :to="{ name: 'Documents: Info', params: { id: relation.getDocumentId() } }">
                                            {{ relation.getDocument()?.getTitle() }}<span
                                                v-if="relation.getDocument()?.getCategory()"> (Category: {{
                                                    relation.getDocument()?.getCategory()?.getName() }})</span>
                                        </router-link>
                                    </td>
                                    <td class="px-6 py-4 text-sm ">
                                        <div class="flex">
                                            <router-link
                                                :to="{ name: 'Citizens: Info', params: { id: relation.getTargetUserId() } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ relation.getTargetUser()?.getFirstname() + ", " +
                                                    relation.getTargetUser()?.getLastname() }}
                                            </router-link>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-right whitespace-nowrap ">
                                        <span class="font-medium ">{{
                                            DOC_RELATION_Util.toEnumKey(relation.getRelation()) }}</span>
                                    </td>
                                    <td v-if="showSource" class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                        <div class="flex">
                                            <router-link
                                                :to="{ name: 'Citizens: Info', params: { id: relation.getSourceUserId() } }"
                                                class="inline-flex space-x-2 text-sm truncate group">
                                                {{ relation.getSourceUser()?.getFirstname() + ", " +
                                                    relation.getSourceUser()?.getLastname() }}
                                            </router-link>
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
</template>

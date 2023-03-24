<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue';
import { watchDebounced } from '@vueuse/shared';
import { getDocStoreClient } from '../../grpc/grpc';
import { FindDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { Document } from '@arpanet/gen/resources/documents/documents_pb';
import { OrderBy, PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import TablePagination from '../partials/TablePagination.vue';
import { CalendarIcon, BriefcaseIcon, UserIcon, DocumentMagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { getDateLocaleString, getDateRelativeString } from '../../utils/time';

const search = ref({ title: '', });
// TODO Implement order by for documents
const orderBys = ref<Array<OrderBy>>([]);
const offset = ref(0);
const totalCount = ref(0);
const listEnd = ref(0);
const documents = ref<Array<Document>>([]);

function findDocuments(pos: number) {
    if (pos < 0) pos = 0;

    const req = new FindDocumentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos));
    req.setSearch(search.value.title);
    req.setOrderbyList([]);

    getDocStoreClient().
        findDocuments(req, null).
        then((resp) => {
            const pag = resp.getPagination();
            if (pag !== undefined) {
                totalCount.value = pag.getTotalCount();
                offset.value = pag.getOffset();
                listEnd.value = pag.getEnd();
            }
            documents.value = resp.getDocumentsList();
        });
}

const searchInput = ref<HTMLInputElement | null>(null);

function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watchDebounced(search.value, () => findDocuments(0), { debounce: 650, maxWait: 1500 });

onBeforeMount(() => {
    findDocuments(0);
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findDocuments(0)">
                        <div class="grid grid-cols-5 gap-4">
                            <div class="col-span-4 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">Title Search</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="search.title" ref="searchInput"
                                        type="text" name="search" id="search"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="form-control">
                                <div class="relative mt-2 flex items-center">
                                    <router-link :to="{ name: 'Documents: Create' }"
                                        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Create</router-link>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <button v-if="documents.length == 0" type="button" @click="focusSearch()"
                            class="relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                            <DocumentMagnifyingGlassIcon class="text-white mx-auto h-12 w-12" />
                            <span class="mt-2 block text-sm font-semibold text-gray-300">No Documents found! Either update
                                your search
                                query or create the first document using the above "Create"-button.</span>
                        </button>
                        <div v-else>
                            <ul role="list" class="divide-y divide-gray-200">
                                <li v-for="doc in documents" :key="doc.getId()">
                                    <router-link :to="{ name: 'Documents: Info', params: { id: doc.getId() } }"
                                        class="block hover:bg-gray-50">
                                        <div class="px-4 py-4 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <p class="truncate text-sm font-medium text-indigo-600">{{ doc.getTitle() }}
                                                </p>
                                                <div class="ml-2 flex flex-shrink-0">
                                                    <p
                                                        class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
                                                        {{ doc.getState() }}</p>
                                                </div>
                                            </div>
                                            <div class="mt-2 sm:flex sm:justify-between">
                                                <div class="sm:flex">
                                                    <p class="max-w-2xl truncate text-gray-300">{{ doc.getContent() }}</p>
                                                </div>
                                                <div class="sm:flex">
                                                    <p class="flex items-center text-sm text-gray-300">
                                                        <UserIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                            aria-hidden="true" />
                                                        {{ doc.getCreator()?.getFirstname() }}, {{
                                                            doc.getCreator()?.getLastname()
                                                        }}
                                                    </p>
                                                    <p class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6">
                                                        <BriefcaseIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                            aria-hidden="true" />
                                                        {{ doc.getCreator()?.getJobLabel() }}
                                                    </p>
                                                </div>
                                                <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                                    <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                        aria-hidden="true" />
                                                    <p>
                                                        Created <time :datetime="getDateLocaleString(doc.getCreatedAt())">{{
                                                            getDateRelativeString(doc.getCreatedAt()) }}</time>
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </router-link>
                                </li>
                            </ul>

                            <TablePagination :offset="offset" :entries="documents.length" :end="listEnd" :total="totalCount"
                                :callback="findDocuments" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

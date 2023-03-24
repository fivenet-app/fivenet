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
    <div class="py-2 sm:min-h-[66rem]">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findDocuments(0)">
                        <label for="search" class="block text-sm font-medium leading-6 text-neutral mb-2">Search</label>
                        <div class="flex flex-row gap-4 mx-auto items-center">
                            <div class="flex-1 form-control">
                                <div class="relative">
                                    <input v-model="search.title" ref="searchInput" type="text" name="search" id="search"
                                        placeholder="Title"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral shadow-sm placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-initial form-control">
                                <div class="relative">
                                    <router-link :to="{ name: 'Documents: Create' }"
                                        class="inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">Create</router-link>
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
                            <DocumentMagnifyingGlassIcon class="text-neutral mx-auto h-12 w-12" />
                            <span class="mt-2 block text-sm font-semibold text-gray-300">No Documents found! Either update
                                your search
                                query or create the first document using the above "Create"-button.</span>
                        </button>
                        <div v-else>
                            <ul class="flex flex-col">
                                <li v-for="doc in documents" :key="doc.getId()" class="flex-initial hover:bg-base-800 bg-base-850 rounded-lg my-1">
                                    <router-link :to="{ name: 'Documents: Info', params: { id: doc.getId() } }">
                                        <div class="mx-2 mt-1 mb-4">
                                            <div class="flex flex-row">
                                            <p class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                                                {{ doc.getTitle() }}
                                            </p>
                                            <p class="px-2 py-2 text-sm text-neutral ml-auto">
                                            <p
                                                class="inline-flex rounded-full bg-primary-100 px-2 text-xs font-semibold leading-5 text-primary-700">
                                                {{ doc.getState() }}</p>
                                            </p>
                                        </div>
                                        <div class="flex flex-row gap-2 text-base-200">
                                            <div class="flex flex-1 justify-start flex-row items-center">
                                                <UserIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                    aria-hidden="true" />
                                                {{ doc.getCreator()?.getFirstname() }}, {{
                                                    doc.getCreator()?.getLastname()
                                                }}
                                            </div>
                                            <div class="flex flex-1 justify-center flex-row items-center">
                                                <BriefcaseIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                    aria-hidden="true" />
                                                {{ doc.getCreator()?.getJobLabel() }}
                                            </div>
                                            <div class="flex flex-1 justify-end flex-row items-center">
                                                <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
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
                                :callback="findDocuments" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

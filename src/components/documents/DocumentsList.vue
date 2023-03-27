<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue';
import { watchDebounced } from '@vueuse/shared';
import { getDocStoreClient } from '../../grpc/grpc';
import { FindDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { Document } from '@arpanet/gen/resources/documents/documents_pb';
import { OrderBy, PaginationRequest, PaginationResponse } from '@arpanet/gen/resources/common/database/database_pb';
import TablePagination from '../partials/TablePagination.vue';
import { CalendarIcon, BriefcaseIcon, UserIcon, DocumentMagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { toDateLocaleString, toDateRelativeString } from '../../utils/time';
import TemplatesModal from './TemplatesModal.vue';

const search = ref({ title: '', });
// TODO Implement order by for documents
const orderBys = ref<Array<OrderBy>>([]);
const pagination = ref<PaginationResponse>();
const documents = ref<Array<Document>>([]);

function findDocuments(pos: number) {
    if (pos < 0) pos = 0;

    const req = new FindDocumentsRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos));
    req.setOrderbyList([]);
    req.setSearch(search.value.title);

    getDocStoreClient().
        findDocuments(req, null).
        then((resp) => {
            pagination.value = resp.getPagination();
            documents.value = resp.getDocumentsList();
        });
}

const searchInput = ref<HTMLInputElement | null>(null);

function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

const templatesOpen = ref(false);

watchDebounced(search.value, () => findDocuments(0), { debounce: 650, maxWait: 1500 });

onBeforeMount(() => {
    findDocuments(0);
});
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" />
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findDocuments(0)">
                        <label for="search" class="block mb-2 text-sm font-medium leading-6 text-neutral">Search</label>
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <input v-model="search.title" ref="searchInput" type="text" name="search" id="search"
                                    placeholder="Title"
                                    class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                            </div>
                            <div class="flex-initial form-control">
                                <button @click="templatesOpen = true"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                    Create
                                </button>
                            </div>
                            <div class="flex-initial">
                                <router-link :to="{ name: 'Documents: Templates' }"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                    Templates
                                </router-link>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <button v-if="documents.length == 0" type="button" @click="focusSearch()"
                            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                            <DocumentMagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">No Documents found! Either update
                                your search
                                query or create the first document using the above "Create"-button.</span>
                        </button>
                        <div v-else>
                            <ul class="flex flex-col">
                                <li v-for="doc in documents" :key="doc.getId()"
                                    class="flex-initial my-1 rounded-lg hover:bg-base-800 bg-base-850">
                                    <router-link :to="{ name: 'Documents: Info', params: { id: doc.getId() } }">
                                        <div class="mx-2 mt-1 mb-4">
                                            <div class="flex flex-row">
                                                <p class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                                                    {{ doc.getTitle() }}
                                                </p>
                                                <p class="px-2 py-2 ml-auto text-sm text-neutral">
                                                <p
                                                    class="inline-flex px-2 text-xs font-semibold leading-5 rounded-full bg-primary-100 text-primary-700">
                                                    {{ doc.getState() }}</p>
                                                </p>
                                            </div>
                                            <div class="flex flex-row gap-2 text-base-200">
                                                <div class="flex flex-row items-center justify-start flex-1">
                                                    <UserIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true" />
                                                    {{ doc.getCreator()?.getFirstname() }}, {{
                                                        doc.getCreator()?.getLastname()
                                                    }}
                                                </div>
                                                <div class="flex flex-row items-center justify-center flex-1">
                                                    <BriefcaseIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true" />
                                                    {{ doc.getCreator()?.getJobLabel() }}
                                                </div>
                                                <div class="flex flex-row items-center justify-end flex-1">
                                                    <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400"
                                                        aria-hidden="true" />
                                                    <p>
                                                        Created <time :datetime="toDateLocaleString(doc.getCreatedAt())">{{
                                                            toDateRelativeString(doc.getCreatedAt()) }}</time>
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </router-link>
                                </li>
                            </ul>

                            <TablePagination :pagination="pagination!" :callback="findDocuments" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

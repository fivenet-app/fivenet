<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import { watchDebounced } from '@vueuse/shared';
import { getDocStoreClient } from '../../grpc/grpc';
import { FindDocumentsRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { Document } from '@arpanet/gen/resources/documents/documents_pb';
import { OrderBy } from '@arpanet/gen/resources/common/database/database_pb';
import TablePagination from '../partials/TablePagination.vue';
import { CalendarIcon, MapPinIcon, UsersIcon } from '@heroicons/vue/20/solid';

const search = ref({ title: '', });
// TODO Implement order by for documents
const orderBys = ref<Array<OrderBy>>([]);
const offset = ref(0);
const totalCount = ref(0);
const listEnd = ref(0);
const documents = ref<Array<Document>>([]);

function findDocuments(pos: number) {
    if (pos < 0) return;

    const req = new FindDocumentsRequest();
    req.setOffset(pos);
    req.setSearch(search.value.title);
    req.setOrderbyList([]);

    getDocStoreClient().
        findDocuments(req, null).
        then((resp) => {
            totalCount.value = resp.getTotalCount();
            offset.value = resp.getOffset();
            listEnd.value = resp.getEnd();
            documents.value = resp.getDocumentsList();
        });
}

watchDebounced(search.value, () => findDocuments(0), { debounce: 750, maxWait: 1500 });

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
                            <div class="col-span-3 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">Search</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="search.title" v-on:keyup.enter="findDocuments(0)" type="text"
                                        name="search" id="search"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="form-control">
                                <div class="relative mt-2 flex items-center">
                                    <router-link to="/documents/create"
                                        class="rounded-md bg-white/10 py-2.5 px-3.5 text-sm font-semibold text-white shadow-sm hover:bg-white/20">Create</router-link>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="overflow-hidden bg-white shadow sm:rounded-md mt-3">
                <ul role="list" class="divide-y divide-gray-200">
                    <li v-for="doc in documents" :key="doc.getId()">
                        <router-link :to="{ name: 'Documents: Info', params: { id: doc.getId() } }"
                            class="block hover:bg-gray-50">
                            <div class="px-4 py-4 sm:px-6">
                                <div class="flex items-center justify-between">
                                    <p class="truncate text-sm font-medium text-indigo-600">{{ doc.getTitle() }}</p>
                                    <div class="ml-2 flex flex-shrink-0">
                                        <p
                                            class="inline-flex rounded-full bg-green-100 px-2 text-xs font-semibold leading-5 text-green-800">
                                            {{ doc.getContentType() }}</p>
                                    </div>
                                </div>
                                <div class="mt-2 sm:flex sm:justify-between">
                                    <div class="sm:flex">
                                        <p class="flex items-center text-sm text-gray-500">
                                            <UsersIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                aria-hidden="true" />
                                            {{ doc.getContent() }}
                                        </p>
                                        <p class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6">
                                            <MapPinIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                                aria-hidden="true" />
                                            {{ doc.getCreator()?.getJob() }}
                                        </p>
                                    </div>
                                    <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                        <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400"
                                            aria-hidden="true" />
                                        <p>
                                            Created at
                                            {{ ' ' }}
                                            <time :datetime="doc.getCreatedAt()?.getTimestamp()?.toDate().toDateString()">{{
                                                doc.getCreatedAt()?.getTimestamp()?.toDate() }}</time>
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </router-link>
                    </li>
                </ul>
            </div>

            <TablePagination :offset="offset" :entries="documents.length" :end="listEnd" :total="totalCount"
                :callback="findDocuments" />
        </div>
    </div>
</template>

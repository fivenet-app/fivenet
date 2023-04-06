<script lang="ts" setup>
import { DocumentCategory } from '@fivenet/gen/resources/documents/category_pb';
import { CompleteDocumentCategoryRequest } from '@fivenet/gen/services/completor/completor_pb';
import { RpcError } from 'grpc-web';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import Cards from '~/components/partials/Cards.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { CardElements } from '~~/src/utils/types';
import CategoryModal from './CategoryModal.vue';

const { $grpc } = useNuxtApp();

const { data: categories, pending, refresh, error } = await useLazyAsyncData(`documents-categories`, () => getCategories());
const items = ref<CardElements>([]);

async function getCategories(): Promise<Array<DocumentCategory>> {
    return new Promise(async (res, rej) => {
        const req = new CompleteDocumentCategoryRequest();

        try {
            const resp = await $grpc.getCompletorClient().
                completeDocumentCategory(req, null);

            return res(resp.getCategoriesList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(categories, () => categories.value?.forEach((v) => {
    items.value.push({ title: v?.getName(), description: v?.getDescription() });
}));

const category = ref<{ name: string, description: string }>({ name: '', description: '' });

const chosenCategory = ref<DocumentCategory>();
const open = ref(false);

function openCategory(idx: number): void {
    chosenCategory.value = categories.value![idx];
    open.value = true;
}
</script>

<template>
    <div>
        <CategoryModal :category="chosenCategory" :open="open" @close="open = false" />
        <div class="py-2">
            <div class="px-2 sm:px-6 lg:px-8">
                <div v-can="'DocStoreService.CreateDocumentCategory'" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <form @submit.prevent="refresh()">
                            <div class="flex flex-row gap-4 mx-auto">
                                <div class="flex-1 form-control">
                                    <label for="search"
                                        class="block text-sm font-medium leading-6 text-neutral">Category</label>
                                    <div class="relative flex items-center mt-2">
                                        <input v-model="category.name" type="text" name="search" id="search"
                                            placeholder="Category"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="search"
                                        class="block text-sm font-medium leading-6 text-neutral">Description</label>
                                    <div class="relative flex items-center mt-2">
                                        <input v-model="category.description" type="text" name="description"
                                            id="description" placeholder="Description"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <div class="relative flex items-center mt-2">
                                        <button type="button"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6">
                                            Create
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
                <div class="flow-root mt-2">
                    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                        <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                            <DataPendingBlock v-if="pending" message="Loading categories..." />
                            <DataErrorBlock v-else-if="error" title="Unable to load categories!" :retry="refresh" />
                            <button v-else-if="categories && categories.length == 0" type="button"
                                class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
                                <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                                <span class="block mt-2 text-sm font-semibold text-base-200">
                                    No categories for your job and rank found.
                                </span>
                            </button>
                            <div v-else>
                                <Cards :items="items" :show-icon="true" @selected="openCategory($event)" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

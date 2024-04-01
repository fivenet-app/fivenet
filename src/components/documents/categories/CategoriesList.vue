<script lang="ts" setup>
import { TagIcon } from 'mdi-vue3';
import CardsList from '~/components/partials/CardsList.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type CardElements } from '~/utils/types';
import { Category } from '~~/gen/ts/resources/documents/category';
import CategoriesModal from '~/components/documents/categories/CategoriesModal.vue';

const { $grpc } = useNuxtApp();

const { data: categories, pending, refresh, error } = useLazyAsyncData(`documents-categories`, () => listCategories());
const items = ref<CardElements>([]);

async function listCategories(): Promise<Category[]> {
    try {
        const call = $grpc.getDocStoreClient().listCategories({});
        const { response } = await call;

        return response.category;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(categories, () => {
    if (items.value) {
        items.value.length = 0;
    }
    categories.value?.forEach((v) => {
        items.value.push({ title: v?.name, description: v?.description });
    });
});

const selectedCategory = ref<Category>();
const open = ref(false);

async function openCategory(idx: number): Promise<void> {
    selectedCategory.value = categories.value![idx];
    open.value = true;
}
</script>

<template>
    <div class="py-2 pb-14">
        <CategoriesModal :category="selectedCategory" :open="open" @close="open = false" @updated="refresh()" />

        <div class="px-1 sm:px-2 lg:px-4">
            <div v-if="can('DocStoreService.CreateCategory')" class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <div class="mx-auto flex flex-row gap-4">
                        <div class="flex-1">
                            <div class="relative mt-2 flex items-center">
                                <button
                                    type="button"
                                    class="inline-flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    @click="
                                        selectedCategory = undefined;
                                        open = true;
                                    "
                                >
                                    {{ $t('components.documents.categories.modal.create_category') }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.category', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.category', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="categories && categories.length === 0"
                            :icon="TagIcon"
                            :type="$t('common.category', 2)"
                        />
                        <div v-else class="flex justify-center">
                            <CardsList :items="items" :show-icon="true" @selected="openCategory($event)" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

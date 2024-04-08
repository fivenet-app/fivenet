<script lang="ts" setup>
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

const modal = useModal();
</script>

<template>
    <div>
        <UDashboardNavbar :title="$t('pages.documents.categories.title')">
            <template #right>
                <UButton
                    v-if="can('DocStoreService.CreateCategory')"
                    color="gray"
                    trailing-icon="i-mdi-plus"
                    @click="
                        modal.open(CategoriesModal, {
                            onUpdate: refresh,
                        })
                    "
                >
                    {{ $t('components.documents.categories.modal.create_category') }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <div class="px-1 sm:px-2">
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
                            icon="i-mdi-tag"
                            :type="$t('common.category', 2)"
                        />
                        <div v-else class="flex justify-center">
                            <CardsList
                                :items="items"
                                :show-icon="true"
                                @selected="
                                    categories &&
                                        modal.open(CategoriesModal, {
                                            category: categories[$event],
                                            onUpdate: refresh,
                                        })
                                "
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

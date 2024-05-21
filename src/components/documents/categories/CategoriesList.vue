<script lang="ts" setup>
import CardsList from '~/components/partials/CardsList.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type CardElements } from '~/utils/types';
import { Category } from '~~/gen/ts/resources/documents/category';
import CategoriesModal from '~/components/documents/categories/CategoriesModal.vue';

const { data: categories, pending: loading, refresh, error } = useLazyAsyncData(`documents-categories`, () => listCategories());

const items = ref<CardElements>([]);

async function listCategories(): Promise<Category[]> {
    try {
        const call = getGRPCDocStoreClient().listCategories({});
        const { response } = await call;

        return response.category;
    } catch (e) {
        handleGRPCError(e as RpcError);
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
    <UDashboardNavbar :title="$t('pages.documents.categories.title')">
        <template #right>
            <UButtonGroup class="inline-flex">
                <UButton color="black" icon="i-mdi-arrow-back" to="/documents">
                    {{ $t('common.back') }}
                </UButton>

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
            </UButtonGroup>
        </template>
    </UDashboardNavbar>

    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.category', 2)])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.category', 2)])" :retry="refresh" />
    <DataNoDataBlock v-else-if="categories && categories.length === 0" icon="i-mdi-tag" :type="$t('common.category', 2)" />

    <div v-else class="flex justify-center">
        <CardsList
            :items="items"
            :show-icon="true"
            class="m-2"
            @selected="
                categories &&
                    modal.open(CategoriesModal, {
                        category: categories[$event],
                        onUpdate: refresh,
                    })
            "
        />
    </div>
</template>

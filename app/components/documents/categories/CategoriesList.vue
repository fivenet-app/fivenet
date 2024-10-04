<script lang="ts" setup>
import CategoriesModal from '~/components/documents/categories/CategoriesModal.vue';
import CardsList from '~/components/partials/CardsList.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { CardElements } from '~/utils/types';
import type { Category } from '~~/gen/ts/resources/documents/category';

const { data: categories, pending: loading, refresh, error } = useLazyAsyncData(`documents-categories`, () => listCategories());

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

const items = ref<CardElements>([]);

watch(categories, () => {
    if (items.value) {
        items.value.length = 0;
    }

    items.value =
        categories.value?.map((v) => ({ title: v?.name, description: v?.description, icon: v.icon, color: v.color })) ?? [];
});

const modal = useModal();
</script>

<template>
    <UDashboardNavbar :title="$t('pages.documents.categories.title')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" to="/documents">
                {{ $t('common.back') }}
            </UButton>

            <UButtonGroup class="inline-flex">
                <UButton
                    v-if="can('DocStoreService.CreateCategory').value"
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

    <UDashboardPanelContent>
        <div v-if="loading" class="flex justify-center">
            <UPageGrid>
                <UPageCard v-for="idx in 6" :key="idx">
                    <template #title>
                        <USkeleton class="h-6 w-[275px]" />
                    </template>
                    <template #description>
                        <USkeleton class="h-11 w-[350px]" />
                    </template>
                </UPageCard>
            </UPageGrid>
        </div>
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.category', 2)])" :retry="refresh" />
        <DataNoDataBlock v-else-if="categories && categories.length === 0" icon="i-mdi-tag" :type="$t('common.category', 2)" />

        <div v-else class="flex justify-center">
            <CardsList
                :items="items"
                @selected="
                    categories &&
                        modal.open(CategoriesModal, {
                            category: categories[$event],
                            onUpdate: refresh,
                        })
                "
            />
        </div>
    </UDashboardPanelContent>
</template>

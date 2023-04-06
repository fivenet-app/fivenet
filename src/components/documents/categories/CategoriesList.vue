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

const { $grpc } = useNuxtApp();

const router = useRouter();

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
    items.value.push({ title: v?.getName(), description: v?.getDescription(), href: { name: 'documents-categories-id', params: { id: v.getId() } } });
}));

async function goToCategory(id: number): Promise<void> {
    await router.push({ name: 'documents-categories-id', params: { id: id, } });
}
</script>

<template>
    <div class="py-2">
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
            <Cards :items="items" :show-icon="true" @selected="goToCategory($event)" />
        </div>
    </div>
</template>

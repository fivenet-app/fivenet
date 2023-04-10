<script lang="ts" setup>
import { DocumentCategory } from '@fivenet/gen/resources/documents/category_pb';
import { RpcError } from 'grpc-web';
import Cards from '~/components/partials/Cards.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { CardElements } from '~/utils/types';
import CategoryModal from './CategoryModal.vue';
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { CreateDocumentCategoryRequest, ListDocumentCategoriesRequest } from '@fivenet/gen/services/docstore/docstore_pb';

const { $grpc } = useNuxtApp();

const { data: categories, pending, refresh, error } = await useLazyAsyncData(`documents-categories`, () => getCategories());
const items = ref<CardElements>([]);

async function getCategories(): Promise<Array<DocumentCategory>> {
    return new Promise(async (res, rej) => {
        const req = new ListDocumentCategoriesRequest();

        try {
            const resp = await $grpc.getDocStoreClient().
                listDocumentCategories(req, null);

            return res(resp.getCategoryList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(categories, () => {
    if (items.value) {
        items.value.length = 0;
    }
    categories.value?.forEach((v) => {
        items.value.push({ title: v?.getName(), description: v?.getDescription() });
    });
});

const chosenCategory = ref<DocumentCategory>();
const open = ref(false);

async function openCategory(idx: number): Promise<void> {
    chosenCategory.value = categories.value![idx];
    open.value = true;
}

async function createDocumentCategory(name: string, description: string): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CreateDocumentCategoryRequest();
        const cat = new DocumentCategory();
        cat.setName(name);
        cat.setDescription(description);
        req.setCategory(cat);

        try {
            const resp = await $grpc.getDocStoreClient().
                createDocumentCategory(req, null);

            refresh();

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { handleSubmit } = useForm({
    validationSchema: toTypedSchema(
        object({
            name: string().required().min(3).max(255),
            description: string().required().max(255),
        }),
    ),
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createDocumentCategory(values.name, values.description));
</script>

<template>
    <div>
        <CategoryModal :category="chosenCategory" :open="open" @close="open = false" @deleted="refresh()" />
        <div class="py-2">
            <div class="px-2 sm:px-6 lg:px-8">
                <div v-can="'DocStoreService.CreateDocumentCategory'" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <form @submit="onSubmit">
                            <div class="flex flex-row gap-4 mx-auto">
                                <div class="flex-1 form-control">
                                    <label for="name"
                                        class="block text-sm font-medium leading-6 text-neutral">Category</label>
                                    <div class="relative flex items-center mt-2">
                                        <Field type="text" name="name" id="name" placeholder="Category"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                        <ErrorMessage name="description" as="p" class="mt-2 text-sm text-red-500" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="description"
                                        class="block text-sm font-medium leading-6 text-neutral">Description</label>
                                    <div class="relative flex items-center mt-2">
                                        <Field type="text" name="description" id="description" placeholder="Description"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                        <ErrorMessage name="description" as="p" class="mt-2 text-sm text-red-500" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <div class="relative flex items-center mt-2">
                                        <button type="submit"
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

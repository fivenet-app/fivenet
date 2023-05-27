<script lang="ts" setup>
import { DocumentCategory } from '~~/gen/ts/resources/documents/category';
import Cards from '~/components/partials/Cards.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { CardElements } from '~/utils/types';
import CategoryModal from './CategoryModal.vue';
import { defineRule } from 'vee-validate';
import { max, min, required } from '@vee-validate/rules';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();

const { data: categories, pending, refresh, error } = useLazyAsyncData(`documents-categories`, () => getCategories());
const items = ref<CardElements>([]);

async function getCategories(): Promise<Array<DocumentCategory>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().
                listDocumentCategories({});
            const { response } = await call;

            return res(response.category);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(categories, () => {
    if (items.value) {
        items.value.length = 0;
    }
    categories.value?.forEach((v) => {
        items.value.push({ title: v?.name, description: v?.description });
    });
});

const chosenCategory = ref<DocumentCategory>();
const open = ref(false);

async function openCategory(idx: number): Promise<void> {
    chosenCategory.value = categories.value![idx];
    open.value = true;
}

async function createDocumentCategory(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().
                createDocumentCategory({
                    category: {
                        id: 0,
                        name: values.name,
                        description: values.description,
                    },
                });

            refresh();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    name: string;
    description: string;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 0, max: 255 },
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createDocumentCategory(values));
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
                                    <label for="name" class="block text-sm font-medium leading-6 text-neutral">
                                        {{ $t('common.category', 1) }}
                                    </label>
                                    <div class="relative flex items-center mt-2">
                                        <VeeField type="text" name="name" :placeholder="$t('common.category', 1)"
                                            :label="$t('common.category', 1)"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                        <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="description" class="block text-sm font-medium leading-6 text-neutral">
                                        {{ $t('common.description') }}
                                    </label>
                                    <div class="relative flex items-center mt-2">
                                        <VeeField type="text" name="description" :placeholder="$t('common.description')"
                                            :label="$t('common.description')"
                                            class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="submit" class="block text-sm font-medium leading-6 text-neutral">
                                        &nbsp;
                                    </label>
                                    <div class="relative flex items-center mt-2">
                                        <button type="submit"
                                            class="block w-full px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                            {{ $t('common.create') }}
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
                            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.category', 2)])" />
                            <DataErrorBlock v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.category', 2)])" :retry="refresh" />
                            <button v-else-if="categories && categories.length === 0" type="button"
                                class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
                                <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                                <span class="block mt-2 text-sm font-semibold text-base-200">
                                    {{ $t('common.not_found',
                                        [$t('components.documents.categories.categories_list.categories_for_your_job',
                                            [$t('common.category', 2), $t('common.job', 1), $t('common.rank')])]) }}
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

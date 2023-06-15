<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiTag } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useNotificationsStore } from '~/store/notifications';
import { DocumentCategory } from '~~/gen/ts/resources/documents/category';

const { $grpc } = useNuxtApp();
const notifications = useNotificationsStore();

const emit = defineEmits<{
    (e: 'deleted'): void;
    (e: 'close'): void;
}>();

const props = defineProps<{
    open: boolean;
    category?: DocumentCategory;
}>();

async function deleteCategory(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().deleteDocumentCategory({
                ids: [props.category?.id!],
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.category_deleted.title', parameters: [] },
                content: { key: 'notifications.category_deleted.content', parameters: [] },
                type: 'success',
            });
            emit('close');
            emit('deleted');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateCategory(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        props.category!.name = values.name;
        props.category!.description = values.description;

        try {
            await $grpc.getDocStoreClient().updateDocumentCategory({
                category: props.category!,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.category_updated.title', parameters: [] },
                content: { key: 'notifications.category_updated.content', parameters: [] },
                type: 'success',
            });
            emit('close');

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

const onSubmit = handleSubmit(async (values): Promise<void> => await updateCategory(values));
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <div>
                            <DialogPanel
                                class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-6xl sm:p-6"
                            >
                                <div>
                                    <div>
                                        <div
                                            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-800"
                                        >
                                            <SvgIcon
                                                class="h-6 w-6 text-primary-500"
                                                aria-hidden="true"
                                                type="mdi"
                                                :path="mdiTag"
                                            />
                                        </div>
                                        <div class="mt-3 text-center sm:mt-5">
                                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                                {{ $t('common.category', 1) }}:
                                                {{ category?.name }}
                                            </DialogTitle>
                                            <div class="mt-2">
                                                <div class="sm:flex-auto">
                                                    <form @submit="onSubmit">
                                                        <div class="flex flex-row gap-4 mx-auto">
                                                            <div class="flex-1 form-control">
                                                                <label
                                                                    for="name"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                    >{{ $t('common.category', 1) }}
                                                                </label>
                                                                <div class="relative flex items-center mt-2">
                                                                    <VeeField
                                                                        type="text"
                                                                        name="name"
                                                                        :placeholder="$t('common.category', 1)"
                                                                        :value="category?.name"
                                                                        :label="$t('common.category', 1)"
                                                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    />
                                                                    <VeeErrorMessage
                                                                        name="category"
                                                                        as="p"
                                                                        class="mt-2 text-sm text-error-400"
                                                                    />
                                                                </div>
                                                            </div>
                                                            <div class="flex-1 form-control">
                                                                <label
                                                                    for="description"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                    >{{ $t('common.description') }}
                                                                </label>
                                                                <div class="relative flex items-center mt-2">
                                                                    <VeeField
                                                                        type="text"
                                                                        name="description"
                                                                        :placeholder="$t('common.description')"
                                                                        :value="category?.description"
                                                                        :label="$t('common.description')"
                                                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    />
                                                                    <VeeErrorMessage
                                                                        name="description"
                                                                        as="p"
                                                                        class="mt-2 text-sm text-error-400"
                                                                    />
                                                                </div>
                                                            </div>
                                                            <div class="flex-1 form-control">
                                                                <div class="relative flex items-center mt-2">
                                                                    <button
                                                                        type="submit"
                                                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    >
                                                                        {{ $t('common.create') }}
                                                                    </button>
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </form>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                        <button
                                            type="button"
                                            v-can="'DocStoreService.DeleteDocumentCategory'"
                                            class="flex-1 rounded-md bg-red-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-red-400"
                                            @click="deleteCategory()"
                                            ref="cancelButtonRef"
                                        >
                                            {{ $t('common.delete') }}
                                        </button>
                                        <button
                                            type="button"
                                            class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                            @click="$emit('close')"
                                            ref="cancelButtonRef"
                                        >
                                            {{ $t('common.close', 1) }}
                                        </button>
                                    </div>
                                </div>
                            </DialogPanel>
                        </div>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

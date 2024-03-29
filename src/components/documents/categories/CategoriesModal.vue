<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon, TagIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { Category } from '~~/gen/ts/resources/documents/category';

const { $grpc } = useNuxtApp();
const notifications = useNotificatorStore();

const emit = defineEmits<{
    (e: 'update'): void;
    (e: 'close'): void;
}>();

const props = defineProps<{
    open: boolean;
    category?: Category;
}>();

interface FormData {
    name: string;
    description: string;
}

async function createCategory(values: FormData): Promise<void> {
    try {
        await $grpc.getDocStoreClient().createCategory({
            category: {
                id: '0',
                name: values.name,
                description: values.description,
            },
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.category_created.title', parameters: {} },
            content: { key: 'notifications.category_created.content', parameters: {} },
            type: 'success',
        });

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateCategory(values: FormData): Promise<void> {
    props.category!.name = values.name;
    props.category!.description = values.description;

    try {
        await $grpc.getDocStoreClient().updateCategory({
            category: props.category!,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.category_updated.title', parameters: {} },
            content: { key: 'notifications.category_updated.content', parameters: {} },
            type: 'success',
        });

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteCategory(): Promise<void> {
    if (props.category === undefined) {
        return;
    }

    try {
        await $grpc.getDocStoreClient().deleteCategory({
            ids: [props.category.id!],
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.category_deleted.title', parameters: {} },
            content: { key: 'notifications.category_deleted.content', parameters: {} },
            type: 'success',
        });
        emit('close');
        emit('update');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 0, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await (props.category === undefined ? createCategory(values) : updateCategory(values)).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
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
                        <DialogPanel
                            class="relative w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:max-w-lg sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <form @submit.prevent="onSubmitThrottle">
                                <div>
                                    <div class="mx-auto flex size-12 items-center justify-center rounded-full bg-primary-100">
                                        <TagIcon class="size-5 text-primary-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-neutral">
                                            <template v-if="category">
                                                {{ $t('components.documents.categories.modal.update_category') }}:
                                                {{ category?.name }}
                                            </template>
                                            <template v-else>
                                                {{ $t('components.documents.categories.modal.create_category') }}
                                            </template>
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="text-sm text-gray-100">
                                                <div class="flex-1">
                                                    <label for="name" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.category', 1) }}
                                                    </label>
                                                    <VeeField
                                                        type="text"
                                                        name="name"
                                                        :placeholder="$t('common.category', 1)"
                                                        :label="$t('common.category', 1)"
                                                        :value="category?.name"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                                <div class="flex-1">
                                                    <label
                                                        for="description"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.description') }}
                                                    </label>
                                                    <VeeField
                                                        as="textarea"
                                                        name="description"
                                                        :placeholder="$t('common.description')"
                                                        :label="$t('common.description')"
                                                        :value="category?.description"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage
                                                        name="description"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="mt-5 sm:mt-6 sm:flex sm:gap-3">
                                    <button
                                        v-if="category !== undefined && can('DocStoreService.DeleteCategory')"
                                        type="button"
                                        class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-error-500 hover:bg-error-400 focus-visible:outline-error-500',
                                        ]"
                                        @click="deleteCategory()"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.delete') }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ category === undefined ? $t('common.create') : $t('common.update') }}
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full justify-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-200 sm:mt-0"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close') }}
                                    </button>
                                </div>
                            </form>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

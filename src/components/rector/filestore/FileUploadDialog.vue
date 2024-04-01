<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import type { File, FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { UploadFileResponse } from '~~/gen/ts/services/rector/filestore';

defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'uploadedFile', file: FileInfo): void;
}>();

const { $grpc } = useNuxtApp();

interface FormData {
    prefix: string;
    name: string;
    file: Blob;
}

async function uploadFile(values: FormData): Promise<UploadFileResponse | undefined> {
    const file = {} as File;

    file.data = new Uint8Array(await values.file.arrayBuffer());

    try {
        const { response } = await $grpc.getRectorFilestoreClient().uploadFile({
            prefix: values.prefix,
            name: values.name,
            file,
        });

        if (response.file) {
            emit('uploadedFile', response.file);
        }

        emit('close');

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const prefixes = ['jobassets'];

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        prefix: { required: true },
        name: { required: true, min: 3, max: 128 },
        file: { required: true, mimes: ['image/jpeg', 'image/jpg', 'image/png'], size: 2000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> => await uploadFile(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
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
                            class="relative h-112 w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-2xl sm:p-6"
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
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('common.upload') }}
                            </DialogTitle>
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="prefix" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.category') }}
                                        </label>
                                        <VeeField
                                            v-slot="{ field }"
                                            name="prefix"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.category')"
                                            :label="$t('common.category')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        >
                                            <select
                                                v-bind="field"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            >
                                                <option
                                                    v-for="prefix in prefixes"
                                                    :key="prefix"
                                                    :selected="field.value === prefix"
                                                    :value="prefix"
                                                >
                                                    {{ prefix }}
                                                </option>
                                            </select>
                                        </VeeField>
                                        <VeeErrorMessage name="prefix" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="name" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.name') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="name"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.name')"
                                            :label="$t('common.name')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="file" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.image') }}
                                        </label>
                                        <template v-if="isNUIAvailable()">
                                            <p class="text-sm text-neutral">
                                                {{ $t('system.not_supported_on_tablet.title') }}
                                            </p>
                                        </template>
                                        <template v-else>
                                            <VeeField
                                                v-slot="{ handleChange, handleBlur }"
                                                name="file"
                                                :placeholder="$t('common.image')"
                                                :label="$t('common.image')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            >
                                                <input
                                                    type="file"
                                                    accept="image/jpeg,image/jpg,image/png"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="handleChange"
                                                    @blur="handleBlur"
                                                />
                                            </VeeField>
                                            <VeeErrorMessage name="file" as="p" class="mt-2 text-sm text-error-400" />
                                        </template>
                                    </div>
                                </div>
                                <div class="absolute bottom-0 left-0 flex w-full">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-neutral px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.save') }}
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

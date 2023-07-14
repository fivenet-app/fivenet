<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCarEmergency } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { digits, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatch';

defineProps<{
    open: boolean;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'created', dsp: Dispatch): void;
}>();

const location = ref<{ x: number; y: number }>({ x: 0, y: 0 });
defineExpose({ location });

const { $grpc } = useNuxtApp();

async function createDispatch(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().createDispatch({
                dispatch: {
                    id: 0n,
                    job: '',
                    message: values.message,
                    description: values.description,
                    anon: values.anon as boolean,
                    attributes: {
                        list: [],
                    },
                    x: location.value.x,
                    y: location.value.y,
                    units: [],
                },
            });
            const { response } = await call;

            emits('created', response.dispatch!);
            emits('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    message: string;
    description: string;
    anon: boolean;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        message: { required: true, min: 3, max: 255 },
        description: { required: false, min: 6, max: 512 },
    },
    initialValues: {
        anon: false,
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createDispatch(values));
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
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-6xl sm:p-6"
                        >
                            <form @submit="onSubmit">
                                <div>
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-800">
                                        <SvgIcon
                                            class="h-6 w-6 text-primary-500"
                                            aria-hidden="true"
                                            type="mdi"
                                            :path="mdiCarEmergency"
                                        />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                            Create Dispatch
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="my-2 space-y-24">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="message"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.message') }}
                                                    </label>
                                                    <VeeField
                                                        type="text"
                                                        name="message"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.message')"
                                                        :label="$t('common.message')"
                                                    />
                                                    <VeeErrorMessage
                                                        name="message"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>
                                            <div class="my-2 space-y-20">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="description"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.description') }}
                                                    </label>
                                                    <VeeField
                                                        type="text"
                                                        name="description"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.description')"
                                                        :label="$t('common.description')"
                                                    />
                                                    <VeeErrorMessage
                                                        name="description"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>
                                            <div class="my-2 space-y-20">
                                                <div class="flex-1 form-control">
                                                    <label for="anon" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.anon') }}
                                                    </label>
                                                    <div class="flex h-6 items-center">
                                                        <VeeField
                                                            type="checkbox"
                                                            name="anon"
                                                            class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600 h-6 w-6"
                                                            :placeholder="$t('common.anon')"
                                                            :label="$t('common.anon')"
                                                        />
                                                    </div>
                                                    <VeeErrorMessage name="anon" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                        ref="cancelButtonRef"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex-1 rounded-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        {{ $t('common.create') }}
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

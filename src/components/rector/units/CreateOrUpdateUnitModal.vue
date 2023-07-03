<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiGroup } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { Unit } from '~~/gen/ts/resources/dispatch/units';

const props = defineProps<{
    open: boolean;
    unit?: Unit;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

async function createOrUpdateUnit(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().createOrUpdateUnit({
                unit: {
                    id: props.unit?.id ?? 0n,
                    name: values.name,
                    initials: values.initials,
                    color: values.color.replaceAll('#', ''),
                    description: values.description,
                    users: [],
                },
            });
            await call;

            emits('refresh');

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
    initials: string;
    description: string;
    color: string;
}

const { handleSubmit, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 24 },
        initials: { required: true, min: 2, max: 4 },
        description: { required: false, max: 255 },
        color: { required: true, max: 7 },
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => createOrUpdateUnit(values));

onMounted(() => {
    if (props.unit) {
        setValues({
            name: props.unit.name,
            initials: props.unit.initials,
            description: props.unit.description,
            color: `#${props.unit.color}`,
        });
    }
});
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
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
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
                            class="relative transform overflow-hidden rounded-lg bg-base-850 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6"
                        >
                            <form @submit="onSubmit">
                                <div>
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
                                        <SvgIcon
                                            type="mdi"
                                            :path="mdiGroup"
                                            class="h-6 w-6 text-green-600"
                                            aria-hidden="true"
                                        />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-white">
                                            {{ $t('pages.rector.units.create_unit') }}
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="text-sm text-gray-100">
                                                <div class="flex-1 form-control">
                                                    <label for="name" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.name') }}
                                                    </label>
                                                    <VeeField
                                                        name="name"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.name')"
                                                        :label="$t('common.name')"
                                                    />
                                                    <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="initials"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.initials') }}
                                                    </label>
                                                    <VeeField
                                                        name="initials"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.initials')"
                                                        :label="$t('common.initials')"
                                                    />
                                                    <VeeErrorMessage
                                                        name="initials"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="description"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.description') }}
                                                    </label>
                                                    <VeeField
                                                        name="description"
                                                        type="text"
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
                                                <div class="flex-1 form-control">
                                                    <label for="color" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.color') }}
                                                    </label>
                                                    <VeeField
                                                        name="color"
                                                        type="color"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :label="$t('common.color')"
                                                    />
                                                    <VeeErrorMessage name="color" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                                    <button
                                        type="submit"
                                        class="inline-flex w-full justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:col-start-2"
                                    >
                                        {{ $t('pages.rector.units.create_unit') }}
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                        @click="$emit('close')"
                                        ref="cancelButtonRef"
                                    >
                                        Cancel
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

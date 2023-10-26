<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, GroupIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ColorInput from 'vue-color-input';
import { Unit } from '~~/gen/ts/resources/dispatch/units';

const props = defineProps<{
    open: boolean;
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update'): void;
}>();

const { $grpc } = useNuxtApp();

async function createOrUpdateUnit(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().createOrUpdateUnit({
                unit: {
                    id: props.unit?.id ?? 0n,
                    job: '',
                    name: values.name,
                    initials: values.initials,
                    color: values.color.replaceAll('#', ''),
                    description: values.description,
                    users: [],
                },
            });
            await call;

            if (props.unit) {
                props.unit.name = values.name;
                props.unit.initials = values.initials;
                props.unit.color = values.color.replaceAll('#', '');
                props.unit.description = values.description;
            }

            emit('update');
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
    initials: string;
    description: string;
    color: string;
}

const { handleSubmit, meta, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 24 },
        initials: { required: true, min: 2, max: 4 },
        description: { required: false, max: 255 },
        color: { required: true, max: 7 },
    },
    initialValues: {
        color: '#000000',
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createOrUpdateUnit(values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

function updateUnitInForm(): void {
    if (props.unit !== undefined) {
        setValues({
            name: props.unit.name,
            initials: props.unit.initials,
            description: props.unit.description,
            color: `#${props.unit.color}`,
        });
    }
}

onBeforeMount(async () => updateUnitInForm());
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
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
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
                            class="relative transform overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 w-full sm:max-w-lg sm:p-6"
                        >
                            <div class="absolute right-0 top-0 pr-4 pt-4 block">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                </button>
                            </div>
                            <form @submit.prevent="onSubmitThrottle">
                                <div>
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-success-100">
                                        <GroupIcon class="h-6 w-6 text-success-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-neutral">
                                            <span v-if="unit && unit?.id">
                                                {{ $t('components.centrum.units.update_unit') }}
                                            </span>
                                            <span v-else>
                                                {{ $t('components.centrum.units.create_unit') }}
                                            </span>
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
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :label="$t('common.color')"
                                                        v-slot="{ field }"
                                                    >
                                                        <ColorInput
                                                            v-model="field.value"
                                                            disable-alpha
                                                            format="hex"
                                                            position="top"
                                                        />
                                                    </VeeField>
                                                    <VeeErrorMessage name="color" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                                    <button
                                        type="submit"
                                        class="inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                        </template>
                                        <span v-if="unit && unit?.id">
                                            {{ $t('components.centrum.units.update_unit') }}
                                        </span>
                                        <span v-else>
                                            {{ $t('components.centrum.units.create_unit') }}
                                        </span>
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full justify-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.cancel') }}
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

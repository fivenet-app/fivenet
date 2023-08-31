<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { CloseIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';

const props = defineProps<{
    open: boolean;
    unit: Unit;
    status?: UNIT_STATUS;
    location?: Coordinate;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const status: number = props.status ?? props.unit?.status?.status ?? UNIT_STATUS.UNKNOWN;

const statuses = ref<{ status: UNIT_STATUS; selected?: boolean }[]>([
    { status: UNIT_STATUS.AVAILABLE },
    { status: UNIT_STATUS.BUSY },
    { status: UNIT_STATUS.ON_BREAK },
    { status: UNIT_STATUS.UNAVAILABLE },
]);
statuses.value.forEach((s) => {
    if (s.status.valueOf() === status) {
        s.selected = true;
    }
});

async function updateUnitStatus(id: bigint, values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateUnitStatus({
                unitId: id,
                status: values.status,
                code: values.code,
                reason: values.reason,
            });
            await call;

            emits('close');

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
    status: number;
    code?: string;
    reason?: string;
}

const { handleSubmit, setFieldValue } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        code: { required: false },
        reason: { required: false, min: 3, max: 255 },
    },
    initialValues: {
        status: status,
    },
    validateOnMount: true,
});

const onSubmit = handleSubmit(async (values): Promise<void> => await updateUnitStatus(props.unit.id, values));

watch(props, () => {
    if (props.status) {
        setFieldValue('status', props.status.valueOf());
    }
});
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-2xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                <form
                                    @submit="onSubmit"
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-white">
                                                    {{ $t('components.centrum.update_unit_status.title') }}: {{ unit.name }} ({{
                                                        unit.initials
                                                    }})
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-white"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-white/10 divide-y divide-white/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="status"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.status') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    name="status"
                                                                    as="div"
                                                                    :placeholder="$t('common.status')"
                                                                    :label="$t('common.status')"
                                                                    v-slot="{ field }"
                                                                >
                                                                    <select
                                                                        v-bind="field"
                                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    >
                                                                        <option
                                                                            v-for="status in statuses"
                                                                            :selected="status.selected"
                                                                            :value="status.status"
                                                                        >
                                                                            {{
                                                                                $t(
                                                                                    `enums.centrum.UNIT_STATUS.${
                                                                                        UNIT_STATUS[
                                                                                            status.status ?? (0 as number)
                                                                                        ]
                                                                                    }`,
                                                                                )
                                                                            }}
                                                                        </option>
                                                                    </select>
                                                                </VeeField>
                                                                <VeeErrorMessage
                                                                    name="status"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="code"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.code') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="code"
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.code')"
                                                                    :label="$t('common.code')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="code"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="reason"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.reason') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="reason"
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.reason')"
                                                                    :label="$t('common.reason')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="reason"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                type="submit"
                                                class="w-full relative inline-flex items-center rounded-l-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                            >
                                                {{ $t('common.update') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10"
                                                @click="$emit('close')"
                                            >
                                                {{ $t('common.close', 1) }}
                                            </button>
                                        </span>
                                    </div>
                                </form>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

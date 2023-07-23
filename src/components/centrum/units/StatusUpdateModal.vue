<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCarEmergency } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';

const props = defineProps<{
    open: boolean;
    unit: Unit;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const location = ref<{ x: number; y: number }>({ x: 0, y: 0 });
defineExpose({ location });

const { $grpc } = useNuxtApp();

const statuses = ref<{ status: UNIT_STATUS; selected?: boolean }[]>([
    { status: UNIT_STATUS.AVAILABLE },
    { status: UNIT_STATUS.BUSY },
    { status: UNIT_STATUS.ON_BREAK },
    { status: UNIT_STATUS.UNAVAILABLE },
]);
statuses.value.forEach((s) => {
    if (s.status === props.unit.status?.status) {
        s.selected = true;
    }
});

async function updateUnitStatus(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateUnitStatus({
                unitId: props.unit.id,
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
    status: UNIT_STATUS;
    code?: string;
    reason: string;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        code: { required: false },
        reason: { required: true, min: 3, max: 255 },
    },
    initialValues: {
        status: props.unit.status?.status,
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await updateUnitStatus(values));
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
                                            Update Unit Status
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="my-2 space-y-24">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="status"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.status') }}
                                                    </label>
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
                                                                {{ UNIT_STATUS[status.status] }}
                                                            </option>
                                                        </select>
                                                    </VeeField>
                                                    <VeeErrorMessage name="status" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                            <div class="my-2 space-y-20">
                                                <div class="flex-1 form-control">
                                                    <label for="code" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.code') }}
                                                    </label>
                                                    <VeeField
                                                        type="text"
                                                        name="code"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.code')"
                                                        :label="$t('common.code')"
                                                    />
                                                    <VeeErrorMessage name="code" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                            <div class="my-2 space-y-20">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="reason"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.reason') }}
                                                    </label>
                                                    <VeeField
                                                        type="text"
                                                        name="reason"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.reason')"
                                                        :label="$t('common.reason')"
                                                    />
                                                    <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
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
                                        {{ $t('common.update') }}
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

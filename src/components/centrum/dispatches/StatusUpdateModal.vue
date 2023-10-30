<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, HoopHouseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { dispatchStatusToBGColor, dispatchStatuses } from '../helpers';

const props = defineProps<{
    open: boolean;
    dispatchId: bigint;
    status?: StatusDispatch;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const status: number = props.status ?? StatusDispatch.NEW;

async function updateDispatchStatus(dispatchId: bigint, values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateDispatchStatus({
                dispatchId: dispatchId,
                status: values.status,
                code: values.code,
                reason: values.reason,
            });
            await call;

            emit('close');

            setFieldValue('status', values.status.valueOf());

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

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
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

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateDispatchStatus(props.dispatchId, values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

watch(props, () => {
    if (props.status) {
        setFieldValue('status', props.status.valueOf());
    }
});
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
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
                                    @submit.prevent="onSubmitThrottle"
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="inline-flex text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('components.centrum.update_dispatch_status.title') }}:
                                                    <IDCopyBadge class="ml-2" :id="dispatchId" prefix="DSP" />
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
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
                                                    <dl class="border-b border-neutral/10 divide-y divide-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
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
                                                                    class="w-full grid grid-cols-2 gap-0.5"
                                                                    :placeholder="$t('common.status')"
                                                                    :label="$t('common.status')"
                                                                    v-slot="{ field }"
                                                                >
                                                                    <button
                                                                        v-for="(item, idx) in dispatchStatuses"
                                                                        :key="item.name"
                                                                        type="button"
                                                                        class="text-neutral bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                                                        :class="[
                                                                            idx >= dispatchStatuses.length - 1
                                                                                ? 'col-span-2'
                                                                                : '',
                                                                            item.class,
                                                                            field.value == item.status
                                                                                ? 'disabled bg-base-500 hover:bg-base-400'
                                                                                : item.status
                                                                                ? dispatchStatusToBGColor(item.status)
                                                                                : item.class,
                                                                            ,
                                                                        ]"
                                                                        :disabled="field.value == item.status"
                                                                        @click="
                                                                            setFieldValue('status', item.status?.valueOf() ?? 0)
                                                                        "
                                                                    >
                                                                        <component
                                                                            :is="item.icon ?? HoopHouseIcon"
                                                                            class="text-base-100 group-hover:text-neutral h-5 w-5 shrink-0"
                                                                            aria-hidden="true"
                                                                        />
                                                                        <span class="mt-1">
                                                                            {{
                                                                                item.status
                                                                                    ? $t(
                                                                                          `enums.centrum.StatusDispatch.${
                                                                                              StatusDispatch[item.status ?? 0]
                                                                                          }`,
                                                                                      )
                                                                                    : $t(item.name)
                                                                            }}
                                                                        </span>
                                                                    </button>
                                                                </VeeField>
                                                                <VeeErrorMessage
                                                                    name="status"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
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
                                                                    @focusin="focusTablet(true)"
                                                                    @focusout="focusTablet(false)"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="code"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
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
                                                                    @focusin="focusTablet(true)"
                                                                    @focusout="focusTablet(false)"
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
                                                class="flex justify-center w-full relative rounded-l-md py-2.5 px-3.5 text-sm font-semibold text-neutral"
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
                                                {{ $t('common.update') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
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

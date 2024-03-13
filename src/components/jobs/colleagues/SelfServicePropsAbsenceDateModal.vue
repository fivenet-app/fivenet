<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import type { JobsUserProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    open: boolean;
    userId: number;
    userProps?: JobsUserProps;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
    absenceBegin?: string;
    absenceEnd?: string;
}

async function setAbsenceDate(values: FormData): Promise<void> {
    const userProps: JobsUserProps = {
        userId: props.userId,
        absenceBegin: values.absenceBegin ? toTimestamp(fromString(values.absenceBegin)) : {},
        absenceEnd: values.absenceEnd ? toTimestamp(fromString(values.absenceEnd)) : {},
    };

    try {
        const call = $grpc.getJobsClient().setJobsUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:absenceDates', {
            userId: props.userId,
            absenceBegin: response.props?.absenceBegin,
            absenceEnd: response.props?.absenceEnd,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            content: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setFieldValue, resetForm, controlledValues } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
        absenceBegin: { required: true },
        absenceEnd: { required: true },
    },
    validateOnMount: true,
    initialValues: {
        absenceBegin: toDatetimeLocal(
            props.userProps?.absenceBegin && toDate(props.userProps.absenceBegin).getTime() > new Date().getTime()
                ? toDate(props.userProps.absenceBegin)
                : new Date(),
        ).split('T')[0],
        absenceEnd:
            props.userProps?.absenceEnd && toDate(props.userProps.absenceEnd).getTime() > new Date().getTime()
                ? toDatetimeLocal(toDate(props.userProps.absenceEnd)).split('T')[0]
                : undefined,
    },
});

function updateAbsenceDateField(): void {
    resetForm();

    if (props.userProps?.absenceBegin && toDate(props.userProps.absenceBegin).getTime() > new Date().getTime()) {
        setFieldValue('absenceBegin', toDatetimeLocal(toDate(props.userProps.absenceBegin)).split('T')[0]);
    } else {
        setFieldValue('absenceBegin', toDatetimeLocal(new Date()).split('T')[0]);
    }

    if (props.userProps?.absenceEnd && toDate(props.userProps.absenceEnd).getTime() > new Date().getTime()) {
        setFieldValue('absenceEnd', toDatetimeLocal(toDate(props.userProps.absenceEnd)).split('T')[0]);
    } else {
        setFieldValue('absenceEnd', undefined);
    }
}

watch(props, () => updateAbsenceDateField());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setAbsenceDate(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
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
                <div class="fixed inset-0 bg-base-900 bg-opacity-75 transition-opacity" />
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
                            class="relative h-96 w-full transform overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-2xl sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('components.jobs.self_service.set_absence_date') }}
                            </DialogTitle>
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2 space-y-24">
                                    <div class="form-control flex-1">
                                        <label for="job" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.reason') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="reason"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.reason')"
                                            :label="$t('common.reason')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-24">
                                    <div class="form-control flex-1">
                                        <label for="absenceBegin" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.from') }}
                                        </label>
                                        <VeeField
                                            type="date"
                                            name="absenceBegin"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.from')"
                                            :label="$t('common.from')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="absenceBegin" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-24">
                                    <div class="form-control flex-1">
                                        <label for="absenceEnd" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.to') }}
                                        </label>
                                        <VeeField
                                            type="date"
                                            name="absenceEnd"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.to')"
                                            :label="$t('common.to')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="absenceEnd" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="absolute bottom-0 left-0 flex w-full">
                                    <button
                                        type="button"
                                        class="rounded-bd flex-1 bg-neutral text-gray-900 hover:bg-gray-200 px-3.5 py-2.5 text-sm font-semibold"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="rounded-bd flex flex-1 justify-center px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
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

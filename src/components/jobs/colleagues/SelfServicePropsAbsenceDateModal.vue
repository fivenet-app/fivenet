<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import type { JobsUserProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    userId: number;
    userProps?: JobsUserProps;
}>();

const emit = defineEmits<{
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

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

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setFieldValue, resetForm } = useForm<FormData>({
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
        await setAbsenceDate(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.jobs.self_service.set_absence_date') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <form @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="reason" class="block text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </label>
                            <VeeField
                                type="text"
                                name="reason"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.reason')"
                                :label="$t('common.reason')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="absenceBegin" class="block text-sm font-medium leading-6">
                                {{ $t('common.from') }}
                            </label>
                            <VeeField
                                type="date"
                                name="absenceBegin"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.from')"
                                :label="$t('common.from')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="absenceBegin" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="absenceEnd" class="block text-sm font-medium leading-6">
                                {{ $t('common.to') }}
                            </label>
                            <VeeField
                                type="date"
                                name="absenceEnd"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.to')"
                                :label="$t('common.to')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="absenceEnd" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                </form>
            </div>

            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                <UButton
                    class="flex-1 rounded-md bg-neutral-50 px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                    @click="isOpen = false"
                >
                    {{ $t('common.close', 1) }}
                </UButton>
                <UButton
                    type="submit"
                    class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold"
                    :disabled="!meta.valid || !canSubmit"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 size-5 animate-spin" />
                    </template>
                    {{ $t('common.save') }}
                </UButton>
            </div>
        </UCard>
    </UModal>
</template>

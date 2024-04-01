<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { max, min, numeric, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
    trafficPoints: number;
    reset?: boolean;
}

async function setTrafficPoints(values: FormData): Promise<void> {
    if (!values.reset && values.trafficPoints === 0) {
        return;
    }

    const points = values.reset ? 0 : (props.user.props?.trafficInfractionPoints ?? 0) + values.trafficPoints;

    const userProps: UserProps = {
        userId: props.user.userId,
        trafficInfractionPoints: points,
    };

    try {
        const call = $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:trafficInfractionPoints', response.props?.trafficInfractionPoints ?? 0);

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
defineRule('numeric', numeric);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
        trafficPoints: { required: true, numeric: true, min: 0, max: 5 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setTrafficPoints(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
    setFieldValue('reset', false);
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
                                    class="rounded-md bg-neutral-50 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('components.citizens.citizen_info_profile.set_traffic_points') }}
                            </DialogTitle>
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="reason" class="block text-sm font-medium leading-6 text-neutral">
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
                                <div class="my-2 space-y-20">
                                    <div class="flex-1">
                                        <label for="trafficPoints" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.traffic_infraction_points') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="trafficPoints"
                                            min="0"
                                            max="9999999"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.traffic_infraction_points')"
                                            :label="$t('common.traffic_infraction_points')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="trafficPoints" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="absolute bottom-0 left-0 flex w-full">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-neutral-50 px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="button"
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-error-500 hover:bg-error-400 focus-visible:outline-error-500',
                                        ]"
                                        @click="
                                            setFieldValue('reset', true);
                                            onSubmitThrottle($event);
                                        "
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.reset') }}
                                    </button>
                                    <button
                                        type="button"
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400',
                                        ]"
                                        @click="
                                            setFieldValue('reset', false);
                                            onSubmitThrottle($event);
                                        "
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.add') }}
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

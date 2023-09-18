<script lang="ts" setup>
import { Dialog, DialogPanel, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { defineRule } from 'vee-validate';
import { useNotificationsStore } from '~/store/notifications';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();
const notifications = useNotificationsStore();

const props = defineProps<{
    open: boolean;
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

async function setWantedState(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        const userProps: UserProps = {
            userId: props.user.userId,
            wanted: props.user.props ? !props.user.props.wanted : true,
        };

        try {
            await $grpc.getCitizenStoreClient().setUserProps({
                props: userProps,
                reason: values.reason,
            });

            if (!props.user.props) {
                props.user.props = userProps;
            } else {
                props.user.props!.wanted = userProps.wanted;
            }

            notifications.dispatchNotification({
                title: { key: 'notifications.action_successfull.title', parameters: [] },
                content: { key: 'notifications.action_successfull.content', parameters: [] },
                type: 'success',
            });

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
    reason: string;
}

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setWantedState(values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn((e) => {
    canSubmit.value = false;
    onSubmit(e);
}, 1000);
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
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 sm:w-full sm:max-w-2xl sm:p-6 h-96"
                        >
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1 form-control">
                                        <label for="job" class="block text-sm font-medium leading-6 text-neutral">
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
                                <div class="absolute bottom-0 w-full left-0 sm:flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex-1 rounded-bd py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
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

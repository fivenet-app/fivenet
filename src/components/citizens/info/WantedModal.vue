<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();
const notifications = useNotificatorStore();

const props = defineProps<{
    open: boolean;
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:wantedStatus', value: boolean): void;
}>();

interface FormData {
    reason: string;
}

async function setWantedState(values: FormData): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        wanted: props.user.props ? !props.user.props.wanted : true,
    };

    try {
        const call = $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:wantedStatus', response.props?.wanted ?? false);

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

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setWantedState(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
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
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
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
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:max-w-2xl sm:p-6 h-96"
                        >
                            <div class="absolute right-0 top-0 pr-4 pt-4 block">
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
                                {{
                                    user.props?.wanted
                                        ? $t('components.citizens.citizen_info_profile.revoke_wanted')
                                        : $t('components.citizens.citizen_info_profile.set_wanted')
                                }}
                            </DialogTitle>
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
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="absolute bottom-0 w-full left-0 flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex-1 flex justify-center rounded-bd py-2.5 px-3.5 text-sm font-semibold text-neutral"
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

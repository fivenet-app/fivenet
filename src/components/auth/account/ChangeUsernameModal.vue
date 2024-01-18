<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
// eslint-disable-next-line camelcase
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { AccountEditIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { clearAuthInfo } = authStore;

const notifications = useNotificatorStore();

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

interface FormData {
    currentUsername: string;
    newUsername: string;
}

async function changeUsername(values: FormData): Promise<void> {
    try {
        const call = $grpc.getAuthClient().changeUsername({
            current: values.currentUsername,
            new: values.newUsername,
        });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.auth.change_username.title', parameters: {} },
            content: { key: 'notifications.auth.change_username.content', parameters: {} },
            type: 'success',
        });

        setTimeout(async () => {
            await navigateTo({ name: 'auth-logout' });
            clearAuthInfo();
        }, 1500);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        currentUsername: { required: true, min: 3, max: 24, alpha_dash: true },
        newUsername: { required: true, min: 3, max: 24, alpha_dash: true },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await changeUsername(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="canSubmit && $emit('close')">
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
                            class="relative w-full transform overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-lg sm:p-6"
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
                            <div>
                                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-700">
                                    <AccountEditIcon class="h-5 w-5 text-primary-500" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        {{ $t('components.auth.change_username_modal.change_username') }}
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <form class="my-2 space-y-6" @submit.prevent="onSubmitThrottle">
                                            <div>
                                                <label for="currentUsername" class="sr-only">{{
                                                    $t('components.auth.change_username_modal.current_username')
                                                }}</label>
                                                <div>
                                                    <VeeField
                                                        name="currentUsername"
                                                        type="text"
                                                        autocomplete="current-username"
                                                        :placeholder="
                                                            $t('components.auth.change_username_modal.current_username')
                                                        "
                                                        :label="$t('components.auth.change_username_modal.current_username')"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    />
                                                    <VeeErrorMessage
                                                        name="currentUsername"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>
                                            <div>
                                                <label for="newUsername" class="sr-only">{{
                                                    $t('components.auth.change_username_modal.new_username')
                                                }}</label>
                                                <div>
                                                    <VeeField
                                                        name="newUsername"
                                                        type="text"
                                                        autocomplete="new-username"
                                                        :placeholder="$t('components.auth.change_username_modal.new_username')"
                                                        :label="$t('components.auth.change_username_modal.new_username')"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    />
                                                    <VeeErrorMessage
                                                        name="newUsername"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>

                                            <div>
                                                <button
                                                    type="submit"
                                                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                                    :disabled="!meta.valid || !canSubmit"
                                                    :class="[
                                                        !meta.valid || !canSubmit
                                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                    ]"
                                                >
                                                    <template v-if="!canSubmit">
                                                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
                                                    </template>
                                                    {{ $t('components.auth.change_username_modal.change_username') }}
                                                </button>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="$emit('close')"
                                >
                                    {{ $t('common.close', 1) }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { AccountKeyIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { setAccessToken } = authStore;

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const newPassword = ref('');

interface FormData {
    currentPassword: string;
    newPassword: string;
}

async function changePassword(values: FormData): Promise<void> {
    try {
        const call = $grpc.getAuthClient().changePassword({
            current: values.currentPassword,
            new: values.newPassword,
        });
        const { response } = await call;

        setAccessToken(response.token, toDate(response.expires) as null | Date);

        notifications.dispatchNotification({
            title: { key: 'notifications.auth.changed_password.title', parameters: {} },
            content: { key: 'notifications.auth.changed_password.content', parameters: {} },
            type: 'success',
        });

        await navigateTo({ name: 'overview' });
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
        currentPassword: { required: true, min: 6, max: 70 },
        newPassword: { required: true, min: 6, max: 70 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await changePassword(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
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
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:max-w-lg sm:p-6"
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
                            <div>
                                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-700">
                                    <AccountKeyIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        {{ $t('components.auth.change_password_modal.change_password') }}
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <form class="my-2 space-y-6" @submit.prevent="onSubmitThrottle">
                                            <div>
                                                <label for="currentPassword" class="sr-only">{{
                                                    $t('components.auth.change_password_modal.current_password')
                                                }}</label>
                                                <div>
                                                    <VeeField
                                                        name="currentPassword"
                                                        type="password"
                                                        autocomplete="current-password"
                                                        :placeholder="
                                                            $t('components.auth.change_password_modal.current_password')
                                                        "
                                                        :label="$t('components.auth.change_password_modal.current_password')"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    />
                                                    <VeeErrorMessage
                                                        name="currentPassword"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>
                                            <div>
                                                <label for="newPassword" class="sr-only">{{
                                                    $t('components.auth.change_password_modal.new_password')
                                                }}</label>
                                                <div>
                                                    <VeeField
                                                        v-model:model-value="newPassword"
                                                        name="newPassword"
                                                        type="password"
                                                        autocomplete="new-password"
                                                        :placeholder="$t('components.auth.change_password_modal.new_password')"
                                                        :label="$t('components.auth.change_password_modal.new_password')"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    />
                                                    <PasswordStrengthMeter :input="newPassword" class="mt-2" />
                                                    <VeeErrorMessage
                                                        name="newPassword"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                            </div>

                                            <div>
                                                <button
                                                    type="submit"
                                                    class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
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
                                                    {{ $t('components.auth.change_password_modal.change_password') }}
                                                </button>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                <button
                                    type="button"
                                    class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
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

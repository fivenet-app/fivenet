<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { KeyIcon } from '@heroicons/vue/24/solid';
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { ChangePasswordRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const notifications = useNotificationsStore();

defineProps({
    open: {
        required: true,
        type: Boolean,
    },
});

defineEmits<{
    (e: 'close'): void,
}>();

const { t } = useI18n();

const newPassword = ref<string>('');

async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new ChangePasswordRequest();
        req.setCurrent(currentPassword);
        req.setNew(newPassword);

        try {
            const resp = await $grpc.getAuthClient()
                .changePassword(req, null);

            store.updateAccessToken(resp.getToken());

            notifications.dispatchNotification({ title: t('notifications.changed_password.title'), content: t('notifications.changed_password.content'), type: 'success' });
            await navigateTo({ name: 'overview' });
            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { handleSubmit } = useForm({
    validationSchema: toTypedSchema(
        object({
            currentPassword: string().required().min(6).max(70),
            newPassword: string().required().min(6).max(70),
        }),
    ),
});

const onSubmit = handleSubmit(async (values): Promise<void> => await changePassword(values.currentPassword, values.newPassword));
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
                            <div>
                                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-800">
                                    <KeyIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                                </div>
                                <div class="mt-3 text-center sm:mt-5">
                                    <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                        Change Password
                                    </DialogTitle>
                                    <div class="mt-2">
                                        <form @submit="onSubmit" class="my-2 space-y-6">
                                            <div>
                                                <label for="currentPassword" class="sr-only">Password</label>
                                                <div>
                                                    <Field id="currentPassword" name="currentPassword" type="password"
                                                        autocomplete="current-password" :placeholder="$t('components.auth.change_password_modal.current_password')"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                                    <ErrorMessage name="currentPassword" as="p"
                                                        class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>
                                            <div>
                                                <label for="newPassword" class="sr-only">Password</label>
                                                <div>
                                                    <Field id="newPassword" name="newPassword" type="password"
                                                        autocomplete="new-password" :placeholder="$t('components.auth.change_password_modal.new_password')"
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        v-model:model-value="newPassword" />
                                                    <PartialsPasswordStrengthMeter :input="newPassword" class="mt-2" />
                                                    <ErrorMessage name="newPassword" as="p"
                                                        class="mt-2 text-sm text-error-400" />
                                                </div>
                                            </div>

                                            <div>
                                                <button type="submit"
                                                    class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                                                    {{ $t('components.auth.change_password_modal.change_password') }}
                                                </button>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                            <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                <button type="button"
                                    class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="$emit('close')" ref="cancelButtonRef">
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

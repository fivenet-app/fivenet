<script lang="ts" setup>
// eslint-disable-next-line camelcase
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { clearAuthInfo } = authStore;

const notifications = useNotificatorStore();

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

        notifications.add({
            title: { key: 'notifications.auth.change_username.title', parameters: {} },
            description: { key: 'notifications.auth.change_username.content', parameters: {} },
            type: 'success',
        });

        useTimeoutFn(async () => {
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
        await changeUsername(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
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
                        {{ $t('components.auth.change_username_modal.change_username') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <UForm :state="{}">
                <UFormGroup name="currentUsername" :label="$t('components.auth.change_username_modal.current_username')">
                    <VeeField
                        name="currentUsername"
                        type="text"
                        autocomplete="current-username"
                        :placeholder="$t('components.auth.change_username_modal.current_username')"
                        :label="$t('components.auth.change_username_modal.current_username')"
                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="currentUsername" as="p" class="mt-2 text-sm text-error-400" />
                </UFormGroup>

                <UFormGroup name="currentUsername" :label="$t('components.auth.change_username_modal.new_username')">
                    <VeeField
                        name="newUsername"
                        type="text"
                        autocomplete="new-username"
                        :placeholder="$t('components.auth.change_username_modal.new_username')"
                        :label="$t('components.auth.change_username_modal.new_username')"
                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="newUsername" as="p" class="mt-2 text-sm text-error-400" />
                </UFormGroup>
            </UForm>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('components.auth.change_username_modal.change_username') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

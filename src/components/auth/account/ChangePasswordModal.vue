<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { setAccessToken } = authStore;

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

        notifications.add({
            title: { key: 'notifications.auth.changed_password.title', parameters: {} },
            description: { key: 'notifications.auth.changed_password.content', parameters: {} },
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
        await changePassword(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
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
                        {{ $t('components.auth.change_password_modal.change_password') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <UForm :state="{}">
                <UFormGroup name="currentPassword" :label="$t('components.auth.change_password_modal.current_password')">
                    <VeeField
                        name="currentPassword"
                        type="password"
                        autocomplete="current-password"
                        :placeholder="$t('components.auth.change_password_modal.current_password')"
                        :label="$t('components.auth.change_password_modal.current_password')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <VeeErrorMessage name="currentPassword" as="p" class="mt-2 text-sm text-error-400" />
                </UFormGroup>

                <UFormGroup name="newPassword" :label="$t('components.auth.change_password_modal.new_password')">
                    <VeeField
                        v-model:model-value="newPassword"
                        name="newPassword"
                        type="password"
                        autocomplete="new-password"
                        :placeholder="$t('components.auth.change_password_modal.new_password')"
                        :label="$t('components.auth.change_password_modal.new_password')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    />
                    <PasswordStrengthMeter :input="newPassword" class="mt-2" />
                    <VeeErrorMessage name="newPassword" as="p" class="mt-2 text-sm text-error-400" />
                </UFormGroup>
            </UForm>

            <template #footer>
                <div class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                    <UButton @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton :disabled="!meta.valid || !canSubmit" :loading="!canSubmit" @click="onSubmitThrottle">
                        {{ $t('components.auth.change_password_modal.change_password') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </UModal>
</template>

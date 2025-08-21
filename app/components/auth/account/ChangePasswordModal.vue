<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useAuthStore } from '~/stores/auth';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { setAccessTokenExpiration } = authStore;

const authAuthClient = await getAuthAuthClient();

const schema = z.object({
    currentPassword: z.string().min(6).max(70),
    newPassword: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    currentPassword: '',
    newPassword: '',
});

async function changePassword(values: Schema): Promise<void> {
    try {
        const call = authAuthClient.changePassword({
            current: values.currentPassword,
            new: values.newPassword,
        });
        const { response } = await call;
        isOpen.value = false;

        setAccessTokenExpiration(toDate(response.expires));

        notifications.add({
            title: { key: 'notifications.auth.changed_password.title', parameters: {} },
            description: { key: 'notifications.auth.changed_password.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'overview' });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const currentPasswordVisibility = ref(false);
function toggleCurrentPasswordVisibility() {
    currentPasswordVisibility.value = !currentPasswordVisibility.value;
}

const newPasswordVisibility = ref(false);
function toggleNewPasswordVisibility() {
    newPasswordVisibility.value = !newPasswordVisibility.value;
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await changePassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" :prevent-close="!canSubmit">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.auth.ChangePasswordModal.change_password') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup name="currentPassword" :label="$t('components.auth.ChangePasswordModal.current_password')">
                    <UInput
                        v-model="state.currentPassword"
                        name="currentPassword"
                        :type="currentPasswordVisibility ? 'text' : 'password'"
                        autocomplete="current-password"
                        :placeholder="$t('components.auth.ChangePasswordModal.current_password')"
                        :ui="{ icon: { trailing: { pointer: '' } } }"
                    >
                        <template #trailing>
                            <UButton
                                color="gray"
                                variant="link"
                                :icon="currentPasswordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                                :padded="false"
                                @click="toggleCurrentPasswordVisibility"
                            />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup name="newPassword" :label="$t('components.auth.ChangePasswordModal.new_password')">
                    <UInput
                        v-model="state.newPassword"
                        name="newPassword"
                        :type="newPasswordVisibility ? 'text' : 'password'"
                        autocomplete="new-password"
                        :placeholder="$t('components.auth.ChangePasswordModal.new_password')"
                        :ui="{ icon: { trailing: { pointer: '' } } }"
                    >
                        <template #trailing>
                            <UButton
                                color="gray"
                                variant="link"
                                :icon="newPasswordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                                :padded="false"
                                @click="toggleNewPasswordVisibility"
                            />
                        </template>
                    </UInput>
                    <PasswordStrengthMeter class="mt-2" :input="state.newPassword" />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('components.auth.ChangePasswordModal.change_password') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

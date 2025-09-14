<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useAuthStore } from '~/stores/auth';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { setAccessTokenExpiration } = authStore;

const authAuthClient = await getAuthAuthClient();

const schema = z.object({
    currentPassword: z.coerce.string().min(6).max(70),
    newPassword: z.coerce.string().min(6).max(70),
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
        emit('close', false);

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
const newPasswordVisibility = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await changePassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.auth.ChangePasswordModal.change_password')" :prevent-close="!canSubmit">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="currentPassword" :label="$t('components.auth.ChangePasswordModal.current_password')">
                    <UInput
                        v-model="state.currentPassword"
                        name="currentPassword"
                        :type="currentPasswordVisibility ? 'text' : 'password'"
                        autocomplete="current-password"
                        :placeholder="$t('components.auth.ChangePasswordModal.current_password')"
                        :ui="{ trailing: 'pe-1' }"
                    >
                        <template #trailing>
                            <UButton
                                color="neutral"
                                variant="link"
                                :icon="currentPasswordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                                :aria-label="currentPasswordVisibility ? 'Hide password' : 'Show password'"
                                :aria-pressed="currentPasswordVisibility"
                                aria-controls="password"
                                @click="currentPasswordVisibility = !currentPasswordVisibility"
                            />
                        </template>
                    </UInput>
                </UFormField>

                <UFormField name="newPassword" :label="$t('components.auth.ChangePasswordModal.new_password')">
                    <UInput
                        v-model="state.newPassword"
                        name="newPassword"
                        :type="newPasswordVisibility ? 'text' : 'password'"
                        autocomplete="new-password"
                        :placeholder="$t('components.auth.ChangePasswordModal.new_password')"
                        :ui="{ trailing: 'pe-1' }"
                    >
                        <template #trailing>
                            <UButton
                                color="neutral"
                                variant="link"
                                :icon="newPasswordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                                :aria-label="newPasswordVisibility ? 'Hide password' : 'Show password'"
                                :aria-pressed="newPasswordVisibility"
                                aria-controls="password"
                                @click="newPasswordVisibility = !newPasswordVisibility"
                            />
                        </template>
                    </UInput>
                    <PasswordStrengthMeter class="mt-1" :input="state.newPassword" />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('components.auth.ChangePasswordModal.change_password')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

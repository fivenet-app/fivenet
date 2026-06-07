<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

const authAuthClient = await getAuthAuthClient();

const schema = z.object({
    currentPassword: passwordSchema,
    newPassword: passwordSchema,
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    currentPassword: '',
    newPassword: '',
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function changePassword(values: Schema): Promise<void> {
    try {
        const call = authAuthClient.changePassword({
            currentPassword: values.currentPassword,
            newPassword: values.newPassword,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.auth.changed_password.title', parameters: {} },
            description: { key: 'notifications.auth.changed_password.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'auth-logout', query: { redirect: '/auth/login' } });
        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const currentPasswordVisibility = ref<boolean>(false);
const newPasswordVisibility = ref<boolean>(false);

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await changePassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="$t('components.auth.ChangePasswordModal.change_password')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.auth.ChangePasswordModal.change_password') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="currentPassword" :label="$t('components.auth.ChangePasswordModal.current_password')">
                    <UInput
                        v-model="state.currentPassword"
                        class="w-full"
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
                        class="w-full"
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
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('components.auth.ChangePasswordModal.change_password')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

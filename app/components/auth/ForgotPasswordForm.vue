<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useNotificatorStore } from '~/store/notificator';
import { getErrorMessage } from '~/utils/errors';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emits = defineEmits<{
    (e: 'toggle'): void;
}>();

const notifications = useNotificatorStore();

async function forgotPassword(values: Schema): Promise<void> {
    try {
        await getGRPCAuthClient().forgotPassword({
            regToken: values.registrationToken.toString(),
            new: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.forgot_password.title', parameters: {} },
            description: { key: 'notifications.auth.forgot_password.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('toggle');
    } catch (e) {
        accountError.value = getErrorMessage((e as RpcError).message);
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const accountError = ref('');

const schema = z.object({
    registrationToken: z.string().length(6).trim(),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const registrationToken = useRouteQuery('registrationToken', '');

const state = reactive<Schema>({
    registrationToken: registrationToken.value,
    password: '',
});

const passwordVisibility = ref(false);

function togglePasswordVisibility() {
    passwordVisibility.value = !passwordVisibility.value;
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await forgotPassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
        <UAlert icon="i-mdi-info-circle">
            <template #description>
                <I18nT keypath="components.auth.ForgotPassword.subtitle">
                    <template #command>
                        <UKbd size="md" :ui="{ size: { md: '' } }" class="h-7 min-w-[24px] text-[13px]">/fivenet</UKbd>
                    </template>
                </I18nT>
            </template>
        </UAlert>

        <UFormGroup name="registrationToken" :label="$t('components.auth.ForgotPassword.registration_token')">
            <UInput
                v-model="state.registrationToken"
                type="text"
                inputmode="numeric"
                aria-describedby="hint"
                pattern="[0-9]*"
                autocomplete="registrationToken"
                :placeholder="$t('components.auth.ForgotPassword.registration_token')"
            />
        </UFormGroup>

        <UFormGroup name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                :type="passwordVisibility ? 'text' : 'password'"
                autocomplete="new-password"
                :placeholder="$t('common.password')"
                :ui="{ icon: { trailing: { pointer: '' } } }"
            >
                <template #trailing>
                    <UButton
                        color="gray"
                        variant="link"
                        :icon="passwordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                        :padded="false"
                        @click="togglePasswordVisibility"
                    />
                </template>
            </UInput>
            <PasswordStrengthMeter :input="state.password" class="mt-2" />
        </UFormGroup>

        <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
            {{ $t('components.auth.ForgotPassword.submit_button') }}
        </UButton>

        <UAlert
            v-if="accountError"
            class="mt-2"
            :title="$t('components.auth.ForgotPassword.create_error')"
            :description="isTranslatedError(accountError) ? $t(accountError) : accountError"
            color="red"
            :close-button="{
                icon: 'i-mdi-window-close',
                color: 'gray',
                variant: 'link',
                padded: false,
            }"
            @close="accountError = ''"
        />
    </UForm>
</template>

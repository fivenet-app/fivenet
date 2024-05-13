<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const accountError = ref('');

const schema = z.object({
    registrationToken: z.string().length(6),
    username: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    registrationToken: '',
    username: '',
    password: '',
});

async function createAccount(values: Schema): Promise<void> {
    try {
        await $grpc.getAuthClient().createAccount({
            regToken: values.registrationToken.toString(),
            username: values.username,
            password: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.account_created.title', parameters: {} },
            description: { key: 'notifications.auth.account_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'auth-login' });
    } catch (e) {
        accountError.value = (e as RpcError).message;
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createAccount(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <h2 class="pb-4 text-center text-3xl">
            {{ $t('components.auth.RegistrationForm.title') }}
        </h2>

        <p class="pb-4 text-sm">
            {{ $t('components.auth.RegistrationForm.subtitle') }}
        </p>

        <UForm :schema="schema" :state="state" class="space-y-2" @submit="onSubmitThrottle">
            <UFormGroup name="registrationToken" :label="$t('components.auth.ForgotPassword.registration_token')">
                <UInput
                    v-model="state.registrationToken"
                    type="text"
                    inputmode="numeric"
                    aria-describedby="hint"
                    pattern="[0-9]*"
                    autocomplete="registrationToken"
                    :placeholder="$t('components.auth.ForgotPassword.registration_token')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </UFormGroup>

            <UFormGroup name="username" :label="$t('common.username')">
                <UInput
                    v-model="state.username"
                    type="text"
                    autocomplete="username"
                    :placeholder="$t('common.username')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </UFormGroup>

            <UFormGroup name="password" :label="$t('common.password')">
                <UInput
                    v-model="state.password"
                    type="password"
                    autocomplete="current-password"
                    :placeholder="$t('common.password')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
                <PasswordStrengthMeter :input="state.password" class="mt-2" />
            </UFormGroup>

            <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                {{ $t('components.auth.RegistrationForm.submit_button') }}
            </UButton>
        </UForm>

        <div class="mt-6">
            <UButton block :to="{ name: 'auth-login' }">
                {{ $t('components.auth.RegistrationForm.back_to_login_button') }}
            </UButton>
        </div>

        <UAlert
            v-if="accountError"
            class="mt-2"
            :title="$t('components.auth.RegistrationForm.create_error')"
            :message="accountError.startsWith('errors.') ? $t(accountError) : accountError"
            color="red"
            :close-button="{
                icon: 'i-mdi-window-close',
                color: 'gray',
                variant: 'link',
                padded: false,
            }"
            @close="accountError = ''"
        />
    </div>
</template>

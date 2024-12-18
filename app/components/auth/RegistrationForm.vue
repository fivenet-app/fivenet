<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';

const notifications = useNotificatorStore();

const accountError = ref<RpcError | undefined>();

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

const registrationToken = useRouteQuery('registrationToken', '');

const state = reactive<Schema>({
    registrationToken: registrationToken.value,
    username: '',
    password: '',
});

async function createAccount(values: Schema): Promise<void> {
    try {
        await getGRPCAuthClient().createAccount({
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
        const err = e as RpcError;
        accountError.value = err;
        handleGRPCError(err);
        throw e;
    }
}

const passwordVisibility = ref(false);

function togglePasswordVisibility() {
    passwordVisibility.value = !passwordVisibility.value;
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

        <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
            <UAlert icon="i-mdi-info-circle">
                <template #description>
                    <I18nT keypath="components.auth.RegistrationForm.subtitle">
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

            <UFormGroup name="username" :label="$t('common.username')">
                <UInput v-model="state.username" type="text" autocomplete="username" :placeholder="$t('common.username')" />
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
                {{ $t('components.auth.RegistrationForm.submit_button') }}
            </UButton>
        </UForm>

        <div class="mt-6">
            <UButton block color="gray" :to="{ name: 'auth-login' }" :disabled="!canSubmit">
                {{ $t('components.auth.RegistrationForm.back_to_login_button') }}
            </UButton>
        </div>

        <DataErrorBlock
            v-if="accountError"
            class="mt-2"
            :title="$t('components.auth.RegistrationForm.create_error')"
            :error="accountError"
            :close="() => (accountError = undefined)"
        />
    </div>
</template>

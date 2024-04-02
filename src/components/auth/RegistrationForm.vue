<script lang="ts" setup>
// eslint-disable-next-line camelcase
import { alpha_dash, digits, max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { defineRule } from 'vee-validate';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const accountError = ref('');
const curPassword = ref('');

interface FormData {
    registrationToken: number;
    username: string;
    password: string;
}

async function createAccount(values: FormData): Promise<void> {
    try {
        await $grpc.getUnAuthClient().createAccount({
            regToken: values.registrationToken.toString(),
            username: values.username,
            password: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.account_created.title', parameters: {} },
            description: { key: 'notifications.auth.account_created.content', parameters: {} },
            type: 'success',
        });

        await navigateTo({ name: 'auth-login' });
    } catch (e) {
        accountError.value = (e as RpcError).message;
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        registrationToken: { required: true, digits: 6 },
        username: { required: true, min: 3, max: 24, alpha_dash: true },
        password: { required: true, min: 6, max: 70 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createAccount(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <div>
        <h2 class="pb-4 text-center text-3xl">
            {{ $t('components.auth.registration_form.title') }}
        </h2>

        <p class="pb-4 text-sm">
            {{ $t('components.auth.registration_form.subtitle') }}
        </p>

        <form class="my-2 space-y-6" @submit.prevent="onSubmitThrottle">
            <div>
                <label for="registrationToken" class="sr-only">
                    {{ $t('components.auth.registration_form.registration_token') }}
                </label>
                <div>
                    <VeeField
                        name="registrationToken"
                        type="text"
                        inputmode="numeric"
                        aria-describedby="hint"
                        pattern="[0-9]*"
                        autocomplete="registrationToken"
                        :placeholder="$t('components.auth.registration_form.registration_token')"
                        :label="$t('components.auth.registration_form.registration_token')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>
            <div>
                <label for="username" class="sr-only">
                    {{ $t('common.username') }}
                </label>
                <div>
                    <VeeField
                        name="username"
                        type="text"
                        autocomplete="username"
                        :placeholder="$t('common.username')"
                        :label="$t('common.username')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="username" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>
            <div>
                <label for="password" class="sr-only">
                    {{ $t('common.password') }}
                </label>
                <div>
                    <VeeField
                        v-model:model-value="curPassword"
                        name="password"
                        type="password"
                        autocomplete="current-password"
                        :placeholder="$t('common.password')"
                        :label="$t('common.password')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <PasswordStrengthMeter :input="curPassword" class="mt-2" />
                    <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>

            <div>
                <UButton type="submit" block :disabled="!meta.valid || !canSubmit" :loading="!canSubmit">
                    {{ $t('components.auth.registration_form.submit_button') }}
                </UButton>
            </div>
        </form>

        <div class="mt-6">
            <UButton :to="{ name: 'auth-login' }" block>
                {{ $t('components.auth.registration_form.back_to_login_button') }}
            </UButton>
        </div>

        <UAlert
            v-if="accountError"
            :title="$t('components.auth.registration_form.create_error')"
            :message="accountError.startsWith('errors.') ? $t(accountError) : accountError"
            color="red"
        />
    </div>
</template>

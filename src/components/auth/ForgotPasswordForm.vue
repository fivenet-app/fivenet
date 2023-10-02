<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { digits, max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import Alert from '~/components/partials/elements/Alert.vue';
import { useNotificatorStore } from '~/store/notificator';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

defineEmits<{
    (e: 'back'): void;
}>();

const newPassword = ref('');

async function forgotPassword(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getUnAuthClient().forgotPassword({
                regToken: values.registrationToken.toString(),
                new: values.password,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.auth.account_created.title', parameters: [] },
                content: { key: 'notifications.auth.account_created.content', parameters: [] },
                type: 'success',
            });

            return res();
        } catch (e) {
            accountError.value = (e as RpcError).message;
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const accountError = ref('');

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    registrationToken: number;
    password: string;
}

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        registrationToken: { required: true, digits: 6 },
        password: { required: true, min: 6, max: 70 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await forgotPassword(values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.forgot_password.title') }}
    </h2>

    <p class="pb-4 text-sm text-white">
        {{ $t('components.auth.forgot_password.subtitle') }}
    </p>

    <form @submit.prevent="onSubmitThrottle" class="my-2 space-y-6">
        <div>
            <label for="registrationToken" class="sr-only">
                {{ $t('components.auth.forgot_password.registration_token') }}
            </label>
            <div>
                <VeeField
                    name="registrationToken"
                    type="text"
                    inputmode="numeric"
                    aria-describedby="hint"
                    pattern="[0-9]*"
                    autocomplete="registrationToken"
                    :placeholder="$t('components.auth.forgot_password.registration_token')"
                    :label="$t('components.auth.forgot_password.registration_token')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6"
                />
                <VeeErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <VeeField
                    name="password"
                    type="password"
                    autocomplete="current-password"
                    :placeholder="$t('common.password')"
                    :label="$t('common.password')"
                    v-model:model-value="newPassword"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                />
                <PasswordStrengthMeter :input="newPassword" class="mt-2" />
                <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button
                type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                :disabled="!meta.valid || !canSubmit"
                :class="[
                    !meta.valid || !canSubmit
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                ]"
            >
                <template v-if="!canSubmit">
                    <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                </template>
                {{ $t('components.auth.forgot_password.submit_button') }}
            </button>
        </div>
    </form>

    <div class="mt-6">
        <button
            type="button"
            @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
        >
            {{ $t('components.auth.forgot_password.back_to_login_button') }}
        </button>
    </div>

    <Alert v-if="accountError" :title="$t('components.auth.forgot_password.create_error')" :message="accountError" />
</template>

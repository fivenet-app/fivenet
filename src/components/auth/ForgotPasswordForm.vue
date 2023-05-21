<script lang="ts" setup>
import { ForgotPasswordRequest } from '@fivenet/gen/services/auth/auth_pb';
import { digits, max, min, required } from '@vee-validate/rules';
import { RpcError } from 'grpc-web';
import { ErrorMessage, Field, Form, defineRule } from 'vee-validate';
import Alert from '~/components/partials/Alert.vue';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

defineEmits<{
    (e: 'back'): void,
}>();

const { t } = useI18n();

const form = ref<{ currPassword: string; registrationToken: number; }>({
    registrationToken: 0,
    currPassword: '',
});

async function forgotPassword(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new ForgotPasswordRequest();
        req.setRegToken(form.value.registrationToken.toString());
        req.setNew(form.value.currPassword);

        try {
            await $grpc.getUnAuthClient().
                forgotPassword(req, null);

            notifications.dispatchNotification({
                title: t('notifications.auth.account_created.title'),
                content: t('notifications.auth.account_created.content'),
                type: 'success'
            });

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            accountError.value = (e as RpcError).message;
            return rej(e as RpcError);
        }
    });
}

const accountError = ref('');

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.forgot_password.title') }}
    </h2>

    <p class="pb-4 text-sm text-white">
        {{ $t('components.auth.forgot_password.subtitle') }}
    </p>

    <Form @submit.prevent="forgotPassword" class="my-2 space-y-6">
        <div>
            <label for="registrationToken" class="sr-only">
                {{ $t('components.auth.forgot_password.registration_token') }}
            </label>
            <div>
                <Field name="registrationToken" type="text" inputmode="numeric" v-model="form.registrationToken"
                    aria-describedby="hint" pattern="[0-9]*" autocomplete="registrationToken"
                    :placeholder="$t('components.auth.forgot_password.registration_token')"
                    :label="$t('components.auth.forgot_password.registration_token')" :rules="{ required: true, digits: 6 }"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                <ErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <Field name="password" type="password" autocomplete="current-password" :placeholder="$t('common.password')"
                    :label="$t('common.password')" :rules="{ required: true, min: 6, max: 70 }"
                    v-model:model-value="form.currPassword"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <PartialsPasswordStrengthMeter :input="form.currPassword" class="mt-2" />
                <ErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('components.auth.forgot_password.submit_button') }}
            </button>
        </div>
    </Form>

    <div class="mt-6">
        <button type="button" @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            {{ $t('components.auth.forgot_password.back_to_login_button') }}
        </button>
    </div>

    <Alert v-if="accountError" :title="$t('components.auth.forgot_password.create_error')" :message="accountError" />
</template>

<script lang="ts" setup>
import { ForgotPasswordRequest } from '@fivenet/gen/services/auth/auth_pb';
import { digits, max, min, required } from '@vee-validate/rules';
import { RpcError } from 'grpc-web';
import { defineRule } from 'vee-validate';
import Alert from '~/components/partials/Alert.vue';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

defineEmits<{
    (e: 'back'): void,
}>();

const { t } = useI18n();

const newPassword = ref('');

async function forgotPassword(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new ForgotPasswordRequest();
        req.setRegToken(values.registrationToken.toString());
        req.setNew(values.password);

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

interface FormData {
    registrationToken: number;
    password: string;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        registrationToken: { required: true, digits: 6 },
        password: { required: true, min: 6, max: 70 },
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await forgotPassword(values));
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.forgot_password.title') }}
    </h2>

    <p class="pb-4 text-sm text-white">
        {{ $t('components.auth.forgot_password.subtitle') }}
    </p>

    <form @submit="onSubmit" class="my-2 space-y-6">
        <div>
            <label for="registrationToken" class="sr-only">
                {{ $t('components.auth.forgot_password.registration_token') }}
            </label>
            <div>
                <VeeField name="registrationToken" type="text" inputmode="numeric"
                    aria-describedby="hint" pattern="[0-9]*" autocomplete="registrationToken"
                    :placeholder="$t('components.auth.forgot_password.registration_token')"
                    :label="$t('components.auth.forgot_password.registration_token')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                <VeeErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <VeeField name="password" type="password" autocomplete="current-password"
                    :placeholder="$t('common.password')" :label="$t('common.password')"
                    v-model:model-value="newPassword"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <PartialsPasswordStrengthMeter :input="newPassword" class="mt-2" />
                <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('components.auth.forgot_password.submit_button') }}
            </button>
        </div>
    </form>

    <div class="mt-6">
        <button type="button" @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            {{ $t('components.auth.forgot_password.back_to_login_button') }}
        </button>
    </div>

    <Alert v-if="accountError" :title="$t('components.auth.forgot_password.create_error')" :message="accountError" />
</template>

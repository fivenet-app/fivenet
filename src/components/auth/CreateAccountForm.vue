<script lang="ts" setup>
import { CreateAccountRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { ErrorMessage, Field, Form, defineRule } from 'vee-validate';
import Alert from '~/components/partials/Alert.vue';
import { useNotificationsStore } from '~/store/notifications';
import { alpha_dash, digits, max, min, required } from '@vee-validate/rules';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

defineEmits<{
    (e: 'back'): void,
}>();

const { t } = useI18n();

const form = ref<{ username: string; password: string; registrationToken: number; }>({
    username: '',
    password: '',
    registrationToken: 0,
});

async function createAccount(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CreateAccountRequest();
        req.setRegToken(form.value.registrationToken.toString());
        req.setUsername(form.value.username);
        req.setPassword(form.value.password);

        try {
            await $grpc.getUnAuthClient().
                createAccount(req, null);

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
defineRule('alpha_dash', alpha_dash);
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.create_account.title') }}
    </h2>

    <p class="pb-4 text-sm text-white">
        {{ $t('components.auth.create_account.subtitle') }}
    </p>

    <Form @submit.prevent="createAccount" class="my-2 space-y-6">
        <div>
            <label for="registrationToken" class="sr-only">
                {{ $t('components.auth.create_account.registration_token') }}
            </label>
            <div>
                <Field name="registrationToken" type="text" inputmode="numeric" :rules="{ required: true, digits: 6 }"
                    aria-describedby="hint" pattern="[0-9]*" autocomplete="registrationToken"
                    :placeholder="$t('components.auth.create_account.registration_token')"
                    :label="$t('components.auth.create_account.registration_token')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                <ErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="username" class="sr-only">
                {{ $t('common.username') }}
            </label>
            <div>
                <Field name="username" type="text" autocomplete="username" :placeholder="$t('common.username')"
                    :label="$t('common.username')"
                    :rules="{ required: true, min: 3, max: 24, alpha_dash: true }"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <ErrorMessage name="Username" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <Field name="password" type="password" autocomplete="current-password" :placeholder="$t('common.password')"
                    :label="$t('common.password')" v-model:model-value="form.password"
                    :rules="{ required: true, min: 6, max: 70 }"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <PartialsPasswordStrengthMeter :input="form.password" class="mt-2" />
                <ErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('components.auth.create_account.submit_button') }}
            </button>
        </div>
    </Form>

    <div class="mt-6">
        <button type="button" @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            {{ $t('components.auth.create_account.back_to_login_button') }}
        </button>
    </div>

    <Alert v-if="accountError" :title="$t('components.auth.create_account.create_error')" :message="accountError" />
</template>

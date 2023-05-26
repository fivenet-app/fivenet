<script lang="ts" setup>
import Alert from '~/components/partials/Alert.vue';
import { useNotificationsStore } from '~/store/notifications';
import { alpha_dash, digits, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

defineEmits<{
    (e: 'back'): void,
}>();

const { t } = useI18n();

const accountError = ref('');
const curPassword = ref('');

async function createAccount(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getUnAuthClient().
                createAccount({
                    regToken: values.registrationToken.toString(),
                    username: values.username,
                    password: values.password,
                });

            notifications.dispatchNotification({
                title: t('notifications.auth.account_created.title'),
                content: t('notifications.auth.account_created.content'),
                type: 'success'
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            accountError.value = (e as RpcError).message;
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);

interface FormData {
    registrationToken: number;
    username: string;
    password: string;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        registrationToken: { required: true, digits: 6 },
        username: { required: true, min: 3, max: 24, alpha_dash: true },
        password: { required: true, min: 6, max: 70 },
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createAccount(values));
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.create_account.title') }}
    </h2>

    <p class="pb-4 text-sm text-white">
        {{ $t('components.auth.create_account.subtitle') }}
    </p>

    <form @submit="onSubmit" class="my-2 space-y-6">
        <div>
            <label for="registrationToken" class="sr-only">
                {{ $t('components.auth.create_account.registration_token') }}
            </label>
            <div>
                <VeeField name="registrationToken" type="text" inputmode="numeric" aria-describedby="hint" pattern="[0-9]*"
                    autocomplete="registrationToken" :placeholder="$t('components.auth.create_account.registration_token')"
                    :label="$t('components.auth.create_account.registration_token')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                <VeeErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="username" class="sr-only">
                {{ $t('common.username') }}
            </label>
            <div>
                <VeeField name="username" type="text" autocomplete="username" :placeholder="$t('common.username')"
                    :label="$t('common.username')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <VeeErrorMessage name="Username" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <VeeField name="password" type="password" autocomplete="current-password"
                    :placeholder="$t('common.password')" :label="$t('common.password')" v-model:model-value="curPassword"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <PartialsPasswordStrengthMeter :input="curPassword" class="mt-2" />
                <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('components.auth.create_account.submit_button') }}
            </button>
        </div>
    </form>

    <div class="mt-6">
        <button type="button" @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            {{ $t('components.auth.create_account.back_to_login_button') }}
        </button>
    </div>

    <Alert v-if="accountError" :title="$t('components.auth.create_account.create_error')" :message="accountError" />
</template>

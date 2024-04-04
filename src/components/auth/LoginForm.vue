<script lang="ts" setup>
// eslint-disable-next-line camelcase
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';

defineEmits<{
    (e: 'toggle'): void;
}>();

const appConfig = useAppConfig();

const authStore = useAuthStore();
const { loginError } = storeToRefs(authStore);
const { doLogin } = authStore;

const configStore = useSettingsStore();
const { isNUIAvailable } = storeToRefs(configStore);

const { cookiesEnabledIds } = useCookieControl();

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);

interface FormData {
    username: string;
    password: string;
}

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        username: { required: true, min: 3, max: 24, alpha_dash: true },
        password: { required: true, min: 6, max: 70 },
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await doLogin(values.username, values.password).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const socialLoginEnabled = ref(false);
if (isNUIAvailable.value) {
    socialLoginEnabled.value = true;
} else if (cookiesEnabledIds.value?.includes('social_login')) {
    socialLoginEnabled.value = true;
}

watch(
    () => cookiesEnabledIds.value,
    (current, previous) => {
        if (!previous?.includes('social_login') && current?.includes('social_login')) {
            socialLoginEnabled.value = true;
        } else {
            socialLoginEnabled.value = false;
        }
    },
    { deep: true },
);
</script>

<template>
    <div>
        <h2 class="pb-4 text-center text-3xl">
            {{ $t('components.auth.login.title') }}
        </h2>

        <form class="my-2 space-y-6" @submit.prevent="onSubmitThrottle">
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
                        name="password"
                        type="password"
                        autocomplete="current-password"
                        :placeholder="$t('common.password')"
                        :label="$t('common.password')"
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>

            <div>
                <UButton type="submit" block :disabled="!meta.valid || !canSubmit" :loading="!canSubmit">
                    {{ $t('common.login') }}
                </UButton>
            </div>
        </form>

        <div class="my-4 space-y-2">
            <template v-if="!isNUIAvailable">
                <p v-if="!socialLoginEnabled" class="mt-2 text-sm text-error-400">
                    {{ $t('components.auth.login.social_login_disabled') }}
                </p>
                <div v-for="provider in appConfig.login.providers" :key="provider.name">
                    <UButton v-if="!socialLoginEnabled" block :disabled="!socialLoginEnabled">
                        {{ provider.label }} {{ $t('common.login') }}
                    </UButton>
                    <UButton
                        v-else
                        block
                        :external="true"
                        :to="`/api/oauth2/login/${provider.name}`"
                        :disabled="!socialLoginEnabled"
                    >
                        {{ provider.label }} {{ $t('common.login') }}
                    </UButton>
                </div>
            </template>
        </div>

        <UAlert
            v-if="loginError"
            :title="$t('components.auth.login.login_error')"
            :message="loginError.startsWith('errors.') ? $t(loginError) : loginError"
            color="red"
        />

        <div class="mt-6">
            <UButton block @click="$emit('toggle')">
                {{ $t('components.auth.login.forgot_password') }}
            </UButton>
        </div>
        <div v-if="appConfig.login.signupEnabled" class="mt-6">
            <UButton block :to="{ name: 'auth-registration' }">
                {{ $t('components.auth.login.register_account') }}
            </UButton>
        </div>
    </div>
</template>

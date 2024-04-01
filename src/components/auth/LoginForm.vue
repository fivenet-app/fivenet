<script lang="ts" setup>
// eslint-disable-next-line camelcase
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
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
        <h2 class="pb-4 text-center text-3xl text-neutral">
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </div>

            <div>
                <button
                    type="submit"
                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    :disabled="!meta.valid || !canSubmit"
                    :class="[
                        !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                    </template>
                    {{ $t('common.login') }}
                </button>
            </div>
        </form>

        <div class="my-4 space-y-2">
            <template v-if="!isNUIAvailable">
                <p v-if="!socialLoginEnabled" class="mt-2 text-sm text-error-400">
                    {{ $t('components.auth.login.social_login_disabled') }}
                </p>
                <div v-for="provider in appConfig.login.providers" :key="provider.name">
                    <button
                        v-if="!socialLoginEnabled"
                        type="button"
                        class="flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                        :class="!socialLoginEnabled ? 'disabled' : ''"
                    >
                        {{ provider.label }} {{ $t('common.login') }}
                    </button>
                    <NuxtLink
                        v-else
                        :external="true"
                        :to="`/api/oauth2/login/${provider.name}`"
                        class="flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                        :class="!socialLoginEnabled ? 'disabled' : ''"
                    >
                        {{ provider.label }} {{ $t('common.login') }}
                    </NuxtLink>
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
            <button
                type="button"
                class="flex w-full justify-center rounded-md bg-secondary-600 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                @click="$emit('toggle')"
            >
                {{ $t('components.auth.login.forgot_password') }}
            </button>
        </div>
        <div v-if="appConfig.login.signupEnabled" class="mt-6">
            <NuxtLink
                :to="{ name: 'auth-registration' }"
                type="button"
                class="flex w-full justify-center rounded-md bg-secondary-600 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
            >
                {{ $t('components.auth.login.register_account') }}
            </NuxtLink>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import Alert from '~/components/partials/elements/Alert.vue';
import { useAuthStore } from '~/store/auth';
import { useConfigStore } from '~/store/config';

const authStore = useAuthStore();
const { loginError } = storeToRefs(authStore);
const { doLogin } = authStore;

const configStore = useConfigStore();
const { appConfig, clientConfig } = storeToRefs(configStore);

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
        await doLogin(values.username, values.password).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const socialLoginEnabled = ref(false);
if (clientConfig.value.NUIEnabled) {
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
    <h2 class="pb-4 text-3xl text-center text-neutral">
        {{ $t('components.auth.login.title') }}
    </h2>

    <form @submit.prevent="onSubmitThrottle" class="my-2 space-y-6">
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
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
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
                {{ $t('common.login') }}
            </button>
        </div>
    </form>

    <div class="my-4 space-y-2">
        <template v-if="!clientConfig.NUIEnabled">
            <p v-if="!socialLoginEnabled" class="mt-2 text-sm text-error-400">
                {{ $t('pages.auth.login.social_login_disabled') }}
            </p>
            <div v-for="prov in appConfig.login.providers" class="">
                <button
                    v-if="!socialLoginEnabled"
                    type="button"
                    class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    :class="!socialLoginEnabled ? 'disabled' : ''"
                >
                    {{ prov.label }} {{ $t('common.login') }}
                </button>
                <NuxtLink
                    v-else
                    :external="true"
                    :to="`/api/oauth2/login/${prov.name}`"
                    class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    :class="!socialLoginEnabled ? 'disabled' : ''"
                >
                    {{ prov.label }} {{ $t('common.login') }}
                </NuxtLink>
            </div>
        </template>
    </div>

    <Alert
        v-if="loginError"
        :title="$t('components.auth.login.login_error')"
        :message="loginError.startsWith('errors.') ? $t(loginError) : loginError"
    />
</template>

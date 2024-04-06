<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
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

const schema = z.object({
    username: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    username: '',
    password: '',
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await doLogin(event.data.username, event.data.password).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
            {{ $t('components.auth.LoginForm.title') }}
        </h2>

        <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
            <UFormGroup name="username" :label="$t('common.username')">
                <UInput
                    v-model="state.username"
                    type="text"
                    autocomplete="username"
                    :placeholder="$t('common.username')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </UFormGroup>

            <UFormGroup name="password" :label="$t('common.password')">
                <UInput
                    v-model="state.password"
                    type="password"
                    autocomplete="current-password"
                    :placeholder="$t('common.password')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </UFormGroup>

            <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                {{ $t('common.login') }}
            </UButton>
        </UForm>

        <div class="my-4 space-y-2">
            <template v-if="!isNUIAvailable">
                <p v-if="!socialLoginEnabled" class="mt-2 text-sm text-error-400">
                    {{ $t('components.auth.LoginForm.social_login_disabled') }}
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
            class="mt-2"
            :title="$t('components.auth.LoginForm.login_error')"
            :message="loginError.startsWith('errors.') ? $t(loginError) : loginError"
            color="red"
        />

        <div class="mt-6">
            <UButton block @click="$emit('toggle')">
                {{ $t('components.auth.LoginForm.forgot_password') }}
            </UButton>
        </div>
        <div v-if="appConfig.login.signupEnabled" class="mt-6">
            <UButton block :to="{ name: 'auth-registration' }">
                {{ $t('components.auth.LoginForm.register_account') }}
            </UButton>
        </div>
    </div>
</template>

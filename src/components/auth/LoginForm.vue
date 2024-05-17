<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';
import { useCookiesStore } from '~/store/cookies';

defineEmits<{
    (e: 'toggle'): void;
}>();

const appConfig = useAppConfig();

const authStore = useAuthStore();
const { loginError } = storeToRefs(authStore);
const { doLogin } = authStore;

const settingsStore = useSettingsStore();
const { isNUIAvailable } = storeToRefs(settingsStore);

const cookiesStore = useCookiesStore();
const { hasCookiesAccepted } = storeToRefs(cookiesStore);

const schema = z.object({
    username: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    username: '',
    password: '',
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await doLogin(event.data.username, event.data.password).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const socialLoginEnabled = ref(hasCookiesAccepted.value && !isNUIAvailable.value);

watch(hasCookiesAccepted, () => (socialLoginEnabled.value = hasCookiesAccepted.value && !isNUIAvailable.value));
</script>

<template>
    <div>
        <h2 class="pb-4 text-center text-3xl">
            {{ $t('components.auth.LoginForm.title') }}
        </h2>

        <UForm :schema="schema" :state="state" class="space-y-2" @submit="onSubmitThrottle">
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

        <div class="space-y-2">
            <template v-if="!isNUIAvailable">
                <p v-if="!socialLoginEnabled" class="text-sm text-error-400">
                    {{ $t('components.auth.LoginForm.social_login_disabled') }}
                </p>

                <template v-else>
                    <UDivider label="OR" orientation="horizontal" class="mt-2" />

                    <div v-for="provider in appConfig.login.providers" :key="provider.name">
                        <UButton
                            block
                            color="white"
                            :external="true"
                            :to="`/api/oauth2/login/${provider.name}`"
                            :disabled="!canSubmit"
                            :icon="provider.icon?.startsWith('i-') ? provider.icon : undefined"
                        >
                            <img
                                v-if="!provider.icon?.startsWith('i-')"
                                :src="provider.icon"
                                :alt="provider.name"
                                class="size-5"
                            />
                            {{ $t('components.auth.LoginForm.login_with', [provider.label]) }}
                        </UButton>
                    </div>
                </template>
            </template>
        </div>

        <UAlert
            v-if="loginError"
            class="mt-2"
            icon="i-mdi-alert"
            :title="$t('components.auth.LoginForm.login_error')"
            :message="loginError.startsWith('errors.') ? $t(loginError) : loginError"
            color="red"
            :close-button="{
                icon: 'i-mdi-window-close',
                color: 'gray',
                variant: 'link',
                padded: false,
            }"
            @close="loginError = ''"
        />

        <div class="space-y-2">
            <UDivider orientation="horizontal" class="mb-4 mt-4" />

            <UButton block color="gray" :disabled="!canSubmit" @click="$emit('toggle')">
                {{ $t('components.auth.LoginForm.forgot_password') }}
            </UButton>

            <UButton
                v-if="appConfig.login.signupEnabled"
                block
                color="gray"
                :disabled="!canSubmit"
                :to="{ name: 'auth-registration' }"
            >
                {{ $t('components.auth.LoginForm.register_account') }}
            </UButton>
        </div>
    </div>
</template>

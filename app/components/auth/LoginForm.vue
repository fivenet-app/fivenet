<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useAuthStore } from '~/store/auth';
import { useCookiesStore } from '~/store/cookies';
import { useSettingsStore } from '~/store/settings';

const { login } = useAppConfig();

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
    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
        <UFormGroup name="username" :label="$t('common.username')">
            <UInput v-model="state.username" type="text" autocomplete="username" :placeholder="$t('common.username')" />
        </UFormGroup>

        <UFormGroup name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                type="password"
                autocomplete="current-password"
                :placeholder="$t('common.password')"
            />
        </UFormGroup>

        <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
            {{ $t('common.login') }}
        </UButton>

        <div v-if="!isNUIAvailable || login.providers.length > 0" class="space-y-2">
            <p v-if="!socialLoginEnabled" class="text-sm text-error-400">
                {{ $t('components.auth.LoginForm.social_login_disabled') }}
            </p>

            <template v-else>
                <UDivider :label="$t('common.or')" orientation="horizontal" class="mt-2" />

                <div v-for="provider in login.providers" :key="provider.name">
                    <UButton
                        block
                        color="white"
                        :external="true"
                        :to="`/api/oauth2/login/${provider.name}`"
                        :disabled="!canSubmit"
                        :icon="provider.icon?.startsWith('i-') ? provider.icon : undefined"
                    >
                        <img v-if="!provider.icon?.startsWith('i-')" :src="provider.icon" :alt="provider.name" class="size-5" />
                        {{ $t('components.auth.LoginForm.login_with', [provider.label]) }}
                    </UButton>
                </div>
            </template>
        </div>

        <UAlert
            v-if="loginError"
            class="mt-2"
            icon="i-mdi-alert"
            :title="$t('components.auth.LoginForm.login_error')"
            :description="isTranslatedError(loginError) ? $t(loginError) : loginError"
            color="red"
            :close-button="{
                icon: 'i-mdi-window-close',
                color: 'gray',
                variant: 'link',
                padded: false,
            }"
            @close="loginError = ''"
        />
    </UForm>
</template>

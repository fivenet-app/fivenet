<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { useCookiesStore } from '~/stores/cookies';
import { useSettingsStore } from '~/stores/settings';

const props = defineProps<{
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', val: boolean): void;
}>();

const canSubmit = useVModel(props, 'modelValue', emit);

const { login } = useAppConfig();

const authStore = useAuthStore();
const { loginError } = storeToRefs(authStore);
const { doLogin } = authStore;

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const cookiesStore = useCookiesStore();
const { hasCookiesAccepted, isConsentModalOpen } = storeToRefs(cookiesStore);

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

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await doLogin(event.data.username, event.data.password).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const socialLoginEnabled = ref(hasCookiesAccepted.value && !nuiEnabled.value);

watch(hasCookiesAccepted, () => (socialLoginEnabled.value = hasCookiesAccepted.value && !nuiEnabled.value));

const passwordVisibility = ref(false);

function togglePasswordVisibility() {
    passwordVisibility.value = !passwordVisibility.value;
}
</script>

<template>
    <UForm class="space-y-4" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UFormGroup name="username" :label="$t('common.username')">
            <UInput v-model="state.username" type="text" autocomplete="username" :placeholder="$t('common.username')" />
        </UFormGroup>

        <UFormGroup name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                :type="passwordVisibility ? 'text' : 'password'"
                autocomplete="current-password"
                :placeholder="$t('common.password')"
                :ui="{ icon: { trailing: { pointer: '' } } }"
            >
                <template #trailing>
                    <UButton
                        color="gray"
                        variant="link"
                        :icon="passwordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                        :padded="false"
                        @click="togglePasswordVisibility"
                    />
                </template>
            </UInput>
        </UFormGroup>

        <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
            {{ $t('common.login') }}
        </UButton>

        <div v-if="!nuiEnabled && login.providers.length > 0" class="space-y-2">
            <UAlert
                v-if="!socialLoginEnabled"
                :description="$t('components.auth.LoginForm.social_login_disabled')"
                color="info"
                variant="subtle"
                :actions="[
                    {
                        label: $t('components.CookieControl.name'),
                        icon: 'i-mdi-cookie',
                        color: 'info',
                        variant: 'outline',
                        click: () => (isConsentModalOpen = true),
                    },
                ]"
            />

            <template v-else>
                <UDivider class="mt-2" :label="$t('common.or')" orientation="horizontal" />

                <div v-for="provider in login.providers" :key="provider.name">
                    <UButton
                        block
                        color="white"
                        :external="true"
                        :to="`/api/oauth2/login/${provider.name}`"
                        :disabled="!canSubmit"
                        :icon="provider.icon?.startsWith('i-') ? provider.icon : undefined"
                    >
                        <NuxtImg
                            v-if="!provider.icon?.startsWith('i-')"
                            class="size-5"
                            :src="provider.icon"
                            :alt="provider.name"
                            placeholder-class="size-5"
                            loading="lazy"
                        />
                        {{ $t('components.auth.LoginForm.login_with', [provider.label]) }}
                    </UButton>
                </div>
            </template>
        </div>

        <DataErrorBlock
            v-if="loginError"
            class="mt-2"
            :title="$t('components.auth.LoginForm.login_error')"
            :error="loginError"
            :close="() => (loginError = null)"
        />
    </UForm>
</template>

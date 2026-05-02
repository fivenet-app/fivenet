<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { openTokenMgmt } from '~/composables/nui';
import { useSettingsStore } from '~/stores/settings';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const authAuthClient = await getAuthAuthClient();

const accountError = ref<RpcError | undefined>();

const schema = z.object({
    registrationToken: z.coerce.string().length(6).trim(),
    username: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
    password: z.coerce.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const registrationToken = useRouteQuery('registrationToken', '');

const state = reactive<Schema>({
    registrationToken: registrationToken.value.trim(),
    username: '',
    password: '',
});

async function createAccount(values: Schema): Promise<void> {
    try {
        await authAuthClient.createAccount({
            regToken: values.registrationToken,
            username: values.username,
            password: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.account_created.title', parameters: {} },
            description: { key: 'notifications.auth.account_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo('/auth/login');
    } catch (e) {
        const err = e as RpcError;
        accountError.value = err;
        handleGRPCError(err);
        throw e;
    }
}

const passwordVisibility = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createAccount(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="space-y-2">
        <h2 class="text-center text-3xl">
            {{ $t('components.auth.registration_form.title') }}
        </h2>

        <UForm class="space-y-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UAlert
                icon="i-mdi-info-circle"
                variant="subtle"
                :actions="
                    !nuiEnabled
                        ? []
                        : [
                              {
                                  label: $t('components.auth.open_token_mgmt'),
                                  onClick: () => openTokenMgmt(),
                              },
                          ]
                "
            >
                <template #description>
                    <I18nT
                        :keypath="
                            !nuiEnabled
                                ? 'components.auth.registration_form.subtitle'
                                : 'components.auth.registration_form.subtitle_nui'
                        "
                    >
                        <template #command>
                            <UKbd class="h-7 min-w-[24px] text-[13px] normal-case" size="md" value="/fivenet" />
                        </template>
                    </I18nT>
                </template>
            </UAlert>

            <UFormField name="registrationToken" :label="$t('components.auth.forgot_password.registration_token')">
                <UInput
                    v-model="state.registrationToken"
                    type="text"
                    inputmode="numeric"
                    aria-describedby="hint"
                    pattern="[0-9]*"
                    autocomplete="registrationToken"
                    :placeholder="$t('components.auth.forgot_password.registration_token')"
                    :ui="{ root: 'w-full' }"
                />
            </UFormField>

            <UFormField name="username" :label="$t('common.username')">
                <UInput
                    v-model="state.username"
                    type="text"
                    autocomplete="username"
                    :placeholder="$t('common.username')"
                    :ui="{ root: 'w-full' }"
                />
            </UFormField>

            <UFormField name="password" :label="$t('common.password')">
                <UInput
                    v-model="state.password"
                    :type="passwordVisibility ? 'text' : 'password'"
                    autocomplete="new-password"
                    :placeholder="$t('common.password')"
                    :ui="{ trailing: 'pe-1', root: 'w-full' }"
                >
                    <template #trailing>
                        <UButton
                            color="neutral"
                            variant="link"
                            size="sm"
                            :icon="passwordVisibility ? 'i-mdi-eye' : 'i-mdi-eye-closed'"
                            :aria-label="passwordVisibility ? 'Hide password' : 'Show password'"
                            :aria-pressed="passwordVisibility"
                            aria-controls="password"
                            @click="passwordVisibility = !passwordVisibility"
                        />
                    </template>
                </UInput>

                <PasswordStrengthMeter class="mt-1" :input="state.password" />
            </UFormField>

            <UButton
                type="submit"
                block
                :disabled="!canSubmit"
                :loading="!canSubmit"
                :label="$t('components.auth.registration_form.submit_button')"
            />

            <USeparator orientation="horizontal" color="primary" />

            <UButton
                block
                color="neutral"
                trailing-icon="i-mdi-login"
                :to="'/auth/login'"
                :disabled="!canSubmit"
                :label="$t('components.auth.registration_form.back_to_login_button')"
            />

            <DataErrorBlock
                v-if="accountError"
                :title="$t('components.auth.registration_form.create_error')"
                :error="accountError"
                :close="() => (accountError = undefined)"
            />
        </UForm>
    </div>
</template>

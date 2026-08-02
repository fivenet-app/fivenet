<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'toggleTab'): void;
}>();

const canSubmit = defineModel<boolean>({ required: true });

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const authAuthClient = await getAuthAuthClient();

const accountError = ref<RpcError | undefined>();

const schema = z.object({
    registrationToken: z.coerce
        .string()
        .length(6)
        .regex(/^[0-9]+$/),
    password: passwordSchema,
});

type Schema = z.output<typeof schema>;

const registrationToken = useRouteQuery('registrationToken', '');

const state = reactive<Schema>({
    registrationToken: registrationToken.value,
    password: '',
});

async function forgotPassword(values: Schema): Promise<void> {
    try {
        await authAuthClient.forgotPassword({
            regToken: values.registrationToken.toString(),
            new: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.forgot_password.title', parameters: {} },
            description: { key: 'notifications.auth.forgot_password.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        canSubmit.value = true;
        emit('toggleTab');
    } catch (e) {
        const err = e as RpcError;
        accountError.value = err;
        handleGRPCError(err);
        throw e;
    }
}

const passwordVisibility = ref<boolean>(false);

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await forgotPassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
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
                            ? 'components.auth.forgot_password.subtitle'
                            : 'components.auth.forgot_password.subtitle_nui'
                    "
                >
                    <template #command>
                        <UKbd class="h-7 min-w-[24px] text-[13px] normal-case" size="md" value="/fivenet" />
                    </template>
                </I18nT>
            </template>
        </UAlert>

        <UFormField
            name="registrationToken"
            :label="$t('components.auth.forgot_password.registration_token')"
            :ui="{
                root: 'w-full',
            }"
        >
            <div class="inline-flex w-full items-center justify-center">
                <UPinInput
                    :model-value="state.registrationToken.split('').map(Number)"
                    type="number"
                    :length="6"
                    :separator="3"
                    autocomplete="registrationToken"
                    size="xl"
                    @update:model-value="($event) => (state.registrationToken = $event.map(String).join(''))"
                />
            </div>
        </UFormField>

        <UFormField name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                :type="passwordVisibility ? 'text' : 'password'"
                autocomplete="new-password"
                :placeholder="$t('common.password')"
                aria-describedby="password-strength"
                :ui="{ trailing: 'pe-1', root: 'w-full' }"
            >
                <template #trailing>
                    <UButton
                        color="neutral"
                        variant="link"
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
            :label="$t('components.auth.forgot_password.submit_button')"
        />

        <DataErrorBlock
            v-if="accountError"
            :title="$t('components.auth.forgot_password.create_error')"
            :error="accountError"
            :close="() => (accountError = undefined)"
        />
    </UForm>
</template>

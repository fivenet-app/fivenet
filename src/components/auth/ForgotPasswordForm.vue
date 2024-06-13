<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import { useNotificatorStore } from '~/store/notificator';
import { getErrorMessage } from '~/utils/errors';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emits = defineEmits<{
    (e: 'toggle'): void;
}>();

const notifications = useNotificatorStore();

async function forgotPassword(values: Schema): Promise<void> {
    try {
        await getGRPCAuthClient().forgotPassword({
            regToken: values.registrationToken.toString(),
            new: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.ForgotPassword.title', parameters: {} },
            description: { key: 'notifications.auth.ForgotPassword.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('toggle');
    } catch (e) {
        accountError.value = (e as RpcError).message;
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const accountError = ref('');

const schema = z.object({
    registrationToken: z.string().length(6),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    registrationToken: '',
    password: '',
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await forgotPassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
        <UAlert icon="i-mdi-info-circle">
            <template #description>
                <I18nT keypath="components.auth.ForgotPassword.subtitle">
                    <template #command>
                        <UKbd size="md" :ui="{ size: { md: '' } }" class="h-7 min-w-[24px] text-[13px]">/fivenet</UKbd>
                    </template>
                </I18nT>
            </template>
        </UAlert>

        <UFormGroup name="registrationToken" :label="$t('components.auth.ForgotPassword.registration_token')">
            <UInput
                v-model="state.registrationToken"
                type="text"
                inputmode="numeric"
                aria-describedby="hint"
                pattern="[0-9]*"
                autocomplete="registrationToken"
                :placeholder="$t('components.auth.ForgotPassword.registration_token')"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </UFormGroup>

        <UFormGroup name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                type="password"
                autocomplete="new-password"
                :placeholder="$t('common.password')"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <PasswordStrengthMeter :input="state.password" class="mt-2" />
        </UFormGroup>

        <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
            {{ $t('components.auth.ForgotPassword.submit_button') }}
        </UButton>

        <UAlert
            v-if="accountError"
            class="mt-2"
            :title="$t('components.auth.ForgotPassword.create_error')"
            :message="getErrorMessage(accountError)"
            color="red"
            :close-button="{
                icon: 'i-mdi-window-close',
                color: 'gray',
                variant: 'link',
                padded: false,
            }"
            @close="accountError = ''"
        />
    </UForm>
</template>

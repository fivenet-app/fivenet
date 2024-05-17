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

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

async function forgotPassword(values: Schema): Promise<void> {
    try {
        await $grpc.getAuthClient().forgotPassword({
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
        $grpc.handleError(e as RpcError);
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
    <div>
        <h2 class="pb-4 text-center text-3xl">
            {{ $t('components.auth.ForgotPassword.title') }}
        </h2>

        <p class="pb-4 text-sm">
            {{ $t('components.auth.ForgotPassword.subtitle') }}
        </p>

        <UForm :schema="schema" :state="state" class="space-y-2" @submit="onSubmitThrottle">
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
                    autocomplete="current-password"
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

        <div class="space-y-2">
            <UDivider orientation="horizontal" class="mb-4 mt-4" />

            <UButton block @click="$emit('toggle')">
                {{ $t('components.auth.ForgotPassword.back_to_login_button') }}
            </UButton>
        </div>
    </div>
</template>

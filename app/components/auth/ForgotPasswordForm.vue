<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import PasswordStrengthMeter from '~/components/auth/PasswordStrengthMeter.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', val: boolean): void;
    (e: 'toggle'): void;
}>();

const { $grpc } = useNuxtApp();

const canSubmit = useVModel(props, 'modelValue', emit);

const notifications = useNotificationsStore();

const accountError = ref<RpcError | undefined>();

const schema = z.object({
    registrationToken: z.string().length(6).trim(),
    password: z.string().min(6).max(70),
});

type Schema = z.output<typeof schema>;

const registrationToken = useRouteQuery('registrationToken', '');

const state = reactive<Schema>({
    registrationToken: registrationToken.value,
    password: '',
});

async function forgotPassword(values: Schema): Promise<void> {
    try {
        await $grpc.auth.auth.forgotPassword({
            regToken: values.registrationToken.toString(),
            new: values.password,
        });

        notifications.add({
            title: { key: 'notifications.auth.forgot_password.title', parameters: {} },
            description: { key: 'notifications.auth.forgot_password.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('toggle');
    } catch (e) {
        const err = e as RpcError;
        accountError.value = err;
        handleGRPCError(err);
        throw e;
    }
}

const passwordVisibility = ref(false);

function togglePasswordVisibility() {
    passwordVisibility.value = !passwordVisibility.value;
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await forgotPassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm class="space-y-4" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UAlert icon="i-mdi-info-circle">
            <template #description>
                <I18nT keypath="components.auth.ForgotPassword.subtitle">
                    <template #command>
                        <UKbd class="h-7 min-w-[24px] text-[13px]" size="md" :ui="{ size: { md: '' } }">/fivenet</UKbd>
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
            />
        </UFormGroup>

        <UFormGroup name="password" :label="$t('common.password')">
            <UInput
                v-model="state.password"
                :type="passwordVisibility ? 'text' : 'password'"
                autocomplete="new-password"
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
            <PasswordStrengthMeter class="mt-2" :input="state.password" />
        </UFormGroup>

        <UButton type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
            {{ $t('components.auth.ForgotPassword.submit_button') }}
        </UButton>

        <DataErrorBlock
            v-if="accountError"
            class="mt-2"
            :title="$t('components.auth.ForgotPassword.create_error')"
            :error="accountError"
            :close="() => (accountError = undefined)"
        />
    </UForm>
</template>

<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

const authAuthClient = await getAuthAuthClient();

const schema = z.object({
    currentUsername: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
    newUsername: z
        .string()
        .min(3)
        .max(24)
        .regex(/^[0-9A-Za-zÄÖÜß_-]{3,24}$/),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    currentUsername: '',
    newUsername: '',
});

async function changeUsername(values: Schema): Promise<void> {
    try {
        const call = authAuthClient.changeUsername({
            current: values.currentUsername,
            new: values.newUsername,
        });
        await call;
        emit('close', false);

        notifications.add({
            title: { key: 'notifications.auth.change_username.title', parameters: {} },
            description: { key: 'notifications.auth.change_username.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'auth-logout', query: { redirect: '/auth/login' } });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await changeUsername(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.auth.ChangeUsernameModal.change_username')" :prevent-close="!canSubmit">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="currentUsername" :label="$t('components.auth.ChangeUsernameModal.current_username')">
                    <UInput
                        v-model="state.currentUsername"
                        type="text"
                        autocomplete="current-username"
                        :placeholder="$t('components.auth.ChangeUsernameModal.current_username')"
                    />
                </UFormField>

                <UFormField name="newUsername" :label="$t('components.auth.ChangeUsernameModal.new_username')">
                    <UInput
                        v-model="state.newUsername"
                        type="text"
                        autocomplete="new-username"
                        :placeholder="$t('components.auth.ChangeUsernameModal.new_username')"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('components.auth.ChangeUsernameModal.change_username')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

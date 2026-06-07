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
    currentUsername: usernameSchema,
    newUsername: usernameSchema,
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    currentUsername: '',
    newUsername: '',
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function changeUsername(values: Schema): Promise<void> {
    try {
        const call = authAuthClient.changeUsername({
            currentUsername: values.currentUsername,
            newUsername: values.newUsername,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.auth.change_username.title', parameters: {} },
            description: { key: 'notifications.auth.change_username.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'auth-logout', query: { redirect: '/auth/login' } });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await changeUsername(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="$t('components.auth.change_username_modal.change_username')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.auth.change_username_modal.change_username') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="currentUsername" :label="$t('components.auth.change_username_modal.current_username')">
                    <UInput
                        v-model="state.currentUsername"
                        class="w-full"
                        type="text"
                        autocomplete="current-username"
                        :placeholder="$t('components.auth.change_username_modal.current_username')"
                    />
                </UFormField>

                <UFormField name="newUsername" :label="$t('components.auth.change_username_modal.new_username')">
                    <UInput
                        v-model="state.newUsername"
                        class="w-full"
                        type="text"
                        autocomplete="new-username"
                        :placeholder="$t('components.auth.change_username_modal.new_username')"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('components.auth.change_username_modal.change_username')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

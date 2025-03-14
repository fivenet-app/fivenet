<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

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
        const call = $grpc.auth.auth.changeUsername({
            current: values.currentUsername,
            new: values.newUsername,
        });
        await call;
        isOpen.value = false;

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
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" :prevent-close="!canSubmit">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.auth.ChangeUsernameModal.change_username') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup name="currentUsername" :label="$t('components.auth.ChangeUsernameModal.current_username')">
                    <UInput
                        v-model="state.currentUsername"
                        type="text"
                        autocomplete="current-username"
                        :placeholder="$t('components.auth.ChangeUsernameModal.current_username')"
                    />
                </UFormGroup>

                <UFormGroup name="newUsername" :label="$t('components.auth.ChangeUsernameModal.new_username')">
                    <UInput
                        v-model="state.newUsername"
                        type="text"
                        autocomplete="new-username"
                        :placeholder="$t('components.auth.ChangeUsernameModal.new_username')"
                    />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('components.auth.ChangeUsernameModal.change_username') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const schema = z.object({
    reason: z.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

async function setWantedState(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        wanted: props.user.props ? !props.user.props.wanted : true,
    };

    try {
        const call = $grpc.citizenstore.citizenStore.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:wantedStatus', response.props?.wanted ?? false);

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setWantedState(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{
                                user.props?.wanted
                                    ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                                    : $t('components.citizens.CitizenInfoProfile.set_wanted')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')" class="flex-1" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

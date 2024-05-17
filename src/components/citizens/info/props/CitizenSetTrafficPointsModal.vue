<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { useNotificatorStore } from '~/store/notificator';
import { User, UserProps } from '~~/gen/ts/resources/users/users';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const schema = z.object({
    reason: z.string().min(3).max(255),
    trafficInfractionPoints: z.coerce.number().int().nonnegative().lt(99999),
    reset: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    trafficInfractionPoints: 0,
    reset: false,
});

async function setTrafficPoints(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        trafficInfractionPoints: values.trafficInfractionPoints,
    };

    if (values.reset) {
        userProps.trafficInfractionPoints = 0;
    }

    try {
        const call = $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:trafficInfractionPoints', response.props?.trafficInfractionPoints ?? 0);

        state.reset = false;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setTrafficPoints(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm ref="form" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.citizens.CitizenInfoProfile.set_traffic_points') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')">
                        <UInput
                            v-model="state.reason"
                            type="text"
                            :placeholder="$t('common.reason')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="trafficInfractionPoints" :label="$t('common.traffic_infraction_points')">
                        <UInput
                            v-model="state.trafficInfractionPoints"
                            type="number"
                            min="0"
                            max="9999999"
                            :placeholder="$t('common.traffic_infraction_points')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.add') }}
                        </UButton>

                        <UButton
                            type="submit"
                            block
                            class="flex-1"
                            color="red"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

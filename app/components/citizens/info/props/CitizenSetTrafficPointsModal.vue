<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const { isOpen } = useOverlay();

const notifications = useNotificationsStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    reason: z.string().min(3).max(255),
    trafficInfractionPoints: z.coerce.number().int().nonnegative().lt(99999),
    reset: z.coerce.boolean(),
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
        const call = citizensCitizensClient.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:trafficInfractionPoints', response.props?.trafficInfractionPoints ?? 0);

        state.reset = false;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(isOpen, () => {
    if (!isOpen.value) {
        return;
    }

    state.trafficInfractionPoints = props.user.props?.trafficInfractionPoints ?? 0;
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setTrafficPoints(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal>
        <UForm ref="form" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ $t('components.citizens.CitizenInfoProfile.set_traffic_points') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="isOpen = false"
                        />
                    </div>
                </template>

                <div>
                    <UFormField name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormField>

                    <UFormField name="trafficInfractionPoints" :label="$t('common.traffic_infraction_points')">
                        <UInput
                            v-model="state.trafficInfractionPoints"
                            type="number"
                            :min="0"
                            :max="9999999"
                            :placeholder="$t('common.traffic_infraction_points')"
                        />
                    </UFormField>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.add') }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            block
                            color="error"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton class="flex-1" color="neutral" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

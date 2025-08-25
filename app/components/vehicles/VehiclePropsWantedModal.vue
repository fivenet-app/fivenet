<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getVehiclesVehiclesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { VehicleProps } from '~~/gen/ts/resources/vehicles/props';

const props = defineProps<{
    plate: string;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const vehicleProps = defineModel<VehicleProps>('vehicleProps');

const notifications = useNotificationsStore();

const vehiclesVehiclesClient = await getVehiclesVehiclesClient();

const schema = z.object({
    reason: z.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

async function setWantedState(values: Schema): Promise<void> {
    const vProps: VehicleProps = {
        plate: props.plate,
        wanted: vehicleProps.value ? !vehicleProps.value.wanted : true,
        wantedReason: values.reason,
    };

    try {
        const call = vehiclesVehiclesClient.setVehicleProps({
            props: vProps,
        });
        const { response } = await call;

        vehicleProps.value = response.props;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('close');
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
    <UModal>
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ vehicleProps?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="$emit('close')"
                        />
                    </div>
                </template>

                <template #body>
                    <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormField>
                </template>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="neutral" block @click="$emit('close')">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

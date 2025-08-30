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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="vehicleProps?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close')" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>

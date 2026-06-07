<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getVehiclesVehiclesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { VehicleProps } from '~~/gen/ts/resources/vehicles/props/props';

const props = defineProps<{
    plate: string;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const vehicleProps = defineModel<VehicleProps>('vehicleProps');

const notifications = useNotificationsStore();

const vehiclesVehiclesClient = await getVehiclesVehiclesClient();

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

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

        emits('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setWantedState(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emits('close', false);
}
</script>

<template>
    <UModal
        :title="vehicleProps?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ vehicleProps?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted') }}
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
                <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" type="text" :placeholder="$t('common.reason')" />
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
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

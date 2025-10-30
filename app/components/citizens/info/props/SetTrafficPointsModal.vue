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
    (e: 'close', v: boolean): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const notifications = useNotificationsStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
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

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(
    () => props.user,
    () => {
        state.trafficInfractionPoints = props.user.props?.trafficInfractionPoints ?? 0;
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setTrafficPoints(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.citizens.CitizenInfoProfile.set_traffic_points')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                </UFormField>

                <UFormField name="trafficInfractionPoints" :label="$t('common.traffic_infraction_points')">
                    <UInputNumber
                        v-model="state.trafficInfractionPoints"
                        :min="0"
                        :max="9999999"
                        :step="1"
                        :placeholder="$t('common.traffic_infraction_points')"
                        class="w-full"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.add')"
                    @click="formRef?.submit()"
                />

                <UButton
                    color="error"
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.reset')"
                    @click="
                        state.reset = true;
                        formRef?.submit();
                    "
                />

                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </UModal>
</template>

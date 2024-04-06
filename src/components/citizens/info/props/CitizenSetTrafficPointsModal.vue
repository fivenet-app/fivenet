<script lang="ts" setup>
import { max, min, numeric, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
    trafficPoints: number;
    reset?: boolean;
}

async function setTrafficPoints(values: FormData): Promise<void> {
    if (!values.reset && values.trafficPoints === 0) {
        return;
    }

    const points = values.reset ? 0 : (props.user.props?.trafficInfractionPoints ?? 0) + values.trafficPoints;

    const userProps: UserProps = {
        userId: props.user.userId,
        trafficInfractionPoints: points,
    };

    try {
        const call = $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:trafficInfractionPoints', response.props?.trafficInfractionPoints ?? 0);

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('numeric', numeric);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
        trafficPoints: { required: true, numeric: true, min: 0, max: 5 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setTrafficPoints(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
    setFieldValue('reset', false);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
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
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="reason" class="block text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </label>
                            <VeeField
                                type="text"
                                name="reason"
                                class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.reason')"
                                :label="$t('common.reason')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2 space-y-20">
                        <div class="flex-1">
                            <label for="trafficPoints" class="block text-sm font-medium leading-6">
                                {{ $t('common.traffic_infraction_points') }}
                            </label>
                            <VeeField
                                type="text"
                                name="trafficPoints"
                                min="0"
                                max="9999999"
                                class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.traffic_infraction_points')"
                                :label="$t('common.traffic_infraction_points')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="trafficPoints" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="
                            setFieldValue('reset', true);
                            onSubmitThrottle($event);
                        "
                    >
                        {{ $t('common.reset') }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="
                            setFieldValue('reset', false);
                            onSubmitThrottle($event);
                        "
                    >
                        {{ $t('common.add') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { RequestStatus, type QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        request: QualificationRequest;
        status?: RequestStatus;
    }>(),
    {
        status: RequestStatus.PENDING,
    },
);

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

interface FormData {
    status: RequestStatus;
    approverComment: string;
}

async function createOrUpdateQualificationRequest(
    qualificationId: string,
    userId: number,
    values: FormData,
): Promise<CreateOrUpdateQualificationRequestResponse> {
    try {
        const call = $grpc.getQualificationsClient().createOrUpdateQualificationRequest({
            request: {
                qualificationId,
                userId,
                status: values.status,
                approverComment: values.approverComment,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emits('refresh');
        emits('close');

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        approverComment: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
    initialValues: {
        status: props.status,
    },
});

watch(props, () => setFieldValue('status', props.status ?? RequestStatus.PENDING));

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> =>
        await createOrUpdateQualificationRequest(props.request.qualificationId, props.request.userId, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const availableStatus = [RequestStatus.ACCEPTED, RequestStatus.DENIED, RequestStatus.PENDING];
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.qualifications.request_modal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit="onSubmitThrottle">
                    <div class="flex-1">
                        <label for="status" class="block text-sm font-medium leading-6">
                            {{ $t('common.status') }}
                        </label>
                        <VeeField
                            v-slot="{ handleChange, field, value }"
                            as="div"
                            name="status"
                            :placeholder="$t('common.status')"
                            :label="$t('common.status')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        >
                            <USelectMenu
                                :model-value="value"
                                :options="availableStatus"
                                :placeholder="
                                    field.value
                                        ? $t(
                                              `enums.qualifications.RequestStatus.${RequestStatus[availableStatus.find((t) => t === field.value) ?? 0]}`,
                                          )
                                        : $t('common.na')
                                "
                                @change="handleChange"
                            >
                                <template #label>
                                    <span v-if="field.value" class="truncate">{{
                                        $t(`enums.qualifications.RequestStatus.${RequestStatus[field.value]}`)
                                    }}</span>
                                </template>
                                <template #option="{ option }">
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.RequestStatus.${RequestStatus[option]}`)
                                    }}</span>
                                </template>
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.attributes', 1)]) }}
                                </template>
                            </USelectMenu>
                        </VeeField>
                        <VeeErrorMessage name="status" as="p" class="mt-2 text-sm text-error-400" />
                    </div>

                    <div class="flex-1">
                        <label for="approverComment" class="block text-sm font-medium leading-6">
                            {{ $t('common.message') }}
                        </label>
                        <VeeField
                            as="textarea"
                            name="approverComment"
                            class="block h-36 w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="$t('common.message')"
                            :label="$t('common.message')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="approverComment" as="p" class="mt-2 text-sm text-error-400" />
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
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.submit') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

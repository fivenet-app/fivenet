<script lang="ts" setup>
import { max, max_value, min, min_value, numeric, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = defineProps<{
    qualificationId: string;
    userId: number;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    status: ResultStatus;
    score: number;
    summary: string;
}

async function createOrUpdateQualificationRequest(
    qualificationId: string,
    userId: number,
    values: FormData,
): Promise<CreateOrUpdateQualificationResultResponse> {
    try {
        const call = $grpc.getQualificationsClient().createOrUpdateQualificationResult({
            result: {
                id: '0',
                qualificationId,
                userId,
                status: values.status,
                score: values.score,
                summary: values.summary,
                creatorId: 0,
                creatorJob: '',
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emits('refresh');
        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('min_value', min_value);
defineRule('max_value', max_value);
defineRule('numeric', numeric);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        score: { required: true, min_value: 0, max_value: 100, numeric: true },
        summary: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
    initialValues: {
        status: ResultStatus.PENDING,
        score: 0,
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> =>
        await createOrUpdateQualificationRequest(props.qualificationId, props.userId, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const availableStatus = [ResultStatus.SUCCESSFUL, ResultStatus.FAILED, ResultStatus.PENDING];
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.qualifications.request_modal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
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
                                              `enums.qualifications.ResultStatus.${ResultStatus[availableStatus.find((t) => t === field.value) ?? 0]}`,
                                          )
                                        : $t('common.na')
                                "
                                @change="handleChange"
                            >
                                <template #label>
                                    <span v-if="field.value" class="truncate">{{
                                        $t(`enums.qualifications.ResultStatus.${ResultStatus[field.value]}`)
                                    }}</span>
                                </template>
                                <template #option="{ option }">
                                    {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[option]}`) }}
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
                        <label for="score" class="block text-sm font-medium leading-6">
                            {{ $t('common.score') }}
                        </label>
                        <VeeField
                            type="number"
                            name="score"
                            min="0"
                            max="100"
                            :placeholder="$t('common.score')"
                            :label="$t('common.score')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="score" as="p" class="mt-2 text-sm text-error-400" />
                    </div>

                    <div class="flex-1">
                        <label for="summary" class="block text-sm font-medium leading-6">
                            {{ $t('common.summary') }}
                        </label>
                        <VeeField
                            as="textarea"
                            name="summary"
                            class="block h-24 w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="$t('common.summary')"
                            :label="$t('common.summary')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="summary" as="p" class="mt-2 text-sm text-error-400" />
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

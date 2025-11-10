<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { RequestStatus, type QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { requestStatusToBadgeColor } from '../helpers';

const props = withDefaults(
    defineProps<{
        request: QualificationRequest;
        status?: RequestStatus;
    }>(),
    {
        status: RequestStatus.PENDING,
    },
);

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const notifications = useNotificationsStore();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const availableStatus = [
    { status: RequestStatus.ACCEPTED },
    { status: RequestStatus.DENIED },
    { status: RequestStatus.PENDING },
];

const schema = z.object({
    status: z.enum(RequestStatus),
    approverComment: z.coerce.string().max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: props.status ?? RequestStatus.PENDING,
    approverComment: '',
});

async function createOrUpdateQualificationRequest(
    qualificationId: number,
    userId: number,
    values: Schema,
): Promise<CreateOrUpdateQualificationRequestResponse> {
    try {
        const call = qualificationsQualificationsClient.createOrUpdateQualificationRequest({
            request: {
                qualificationId: qualificationId,
                userId: userId,
                status: values.status,
                approverComment: values.approverComment,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('refresh');
        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(props, () => (state.status = props.status ?? RequestStatus.PENDING));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateQualificationRequest(props.request.qualificationId, props.request.userId, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.qualifications.request_modal.title')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="status" :label="$t('common.status')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.status"
                            :items="availableStatus"
                            value-key="status"
                            :placeholder="$t('common.status')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            class="w-full"
                        >
                            <template #default>
                                <UBadge class="truncate" :color="requestStatusToBadgeColor(state.status)">
                                    {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[state.status]}`) }}
                                </UBadge>
                            </template>

                            <template #item-label="{ item }">
                                <UBadge class="truncate" :color="requestStatusToBadgeColor(item.status)">
                                    {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[item.status]}`) }}
                                </UBadge>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.status')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="approverComment" :label="$t('common.message')">
                    <UTextarea
                        v-model="state.approverComment"
                        name="approverComment"
                        :rows="3"
                        :placeholder="$t('common.message')"
                        class="w-full"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.submit')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

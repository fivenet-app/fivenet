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
const { t } = useI18n();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const availableStatus = computed(() => {
    switch (props.status) {
        case RequestStatus.ACCEPTED:
            return [{ status: RequestStatus.ACCEPTED }];
        case RequestStatus.DENIED:
            return [{ status: RequestStatus.DENIED }];
        default:
            return [{ status: RequestStatus.ACCEPTED }, { status: RequestStatus.DENIED }, { status: RequestStatus.PENDING }];
    }
});

const schema = z
    .object({
        status: z.enum(RequestStatus),
        approverComment: z.coerce.string().max(255),
        notifyUser: z.coerce.boolean().default(true),
    })
    .superRefine((values, ctx) => {
        if (values.status === RequestStatus.DENIED && values.approverComment.trim().length === 0) {
            ctx.addIssue({
                code: 'custom',
                path: ['approverComment'],
                message: t('components.qualifications.request_tutor_modal.denial_reason_required'),
            });
        }
    });

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: props.status ?? RequestStatus.PENDING,
    approverComment: '',
    notifyUser: true,
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);

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
            skipNotification: !values.notifyUser,
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

watch(props, () => {
    state.status = props.status ?? RequestStatus.PENDING;
    syncSnapshot();
});

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await createOrUpdateQualificationRequest(props.request.qualificationId, props.request.userId, event.data);
});

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="$t('components.qualifications.request_tutor_modal.title')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.qualifications.request_tutor_modal.title') }}
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
            <UForm ref="formRef" :schema="schema" :state="state" class="space-y-2" @submit="submit">
                <UFormField class="flex-1" name="status" :label="$t('common.status')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.status"
                            class="w-full"
                            :items="availableStatus"
                            value-key="status"
                            :placeholder="$t('common.status')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template #default>
                                <UBadge
                                    class="truncate"
                                    :color="requestStatusToBadgeColor(state.status)"
                                    :label="$t(`enums.qualifications.RequestStatus.${RequestStatus[state.status]}`)"
                                />
                            </template>

                            <template #item-label="{ item }">
                                <UBadge
                                    class="truncate"
                                    :color="requestStatusToBadgeColor(item.status)"
                                    :label="$t(`enums.qualifications.RequestStatus.${RequestStatus[item.status]}`)"
                                />
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
                        class="w-full"
                        name="approverComment"
                        :rows="3"
                        :placeholder="$t('common.message')"
                    />
                </UFormField>

                <USwitch v-model="state.notifyUser" :label="$t('components.jobs.groups.details.notify_user')" />
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
                    :loading="isSubmitting"
                    :label="$t('common.submit')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

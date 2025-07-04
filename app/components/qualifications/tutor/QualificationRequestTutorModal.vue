<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { RequestStatus, type QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { requestStatusToBgColor } from '../helpers';

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
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const availableStatus = [
    { status: RequestStatus.ACCEPTED },
    { status: RequestStatus.DENIED },
    { status: RequestStatus.PENDING },
];

const schema = z.object({
    status: z.nativeEnum(RequestStatus),
    approverComment: z.string().max(255),
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
        const call = $grpc.qualifications.qualifications.createOrUpdateQualificationRequest({
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
        isOpen.value = false;

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
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.qualifications.request_modal.title') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup class="flex-1" name="status" :label="$t('common.status')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.status"
                                :options="availableStatus"
                                value-attribute="status"
                                :placeholder="$t('common.status')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
                                    <span class="size-2 rounded-full" :class="requestStatusToBgColor(state.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.RequestStatus.${RequestStatus[state.status]}`)
                                    }}</span>
                                </template>

                                <template #option="{ option }">
                                    <span class="size-2 rounded-full" :class="requestStatusToBgColor(option.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.RequestStatus.${RequestStatus[option.status]}`)
                                    }}</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.status')]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="approverComment" :label="$t('common.message')">
                        <UTextarea
                            v-model="state.approverComment"
                            name="approverComment"
                            :rows="3"
                            :placeholder="$t('common.message')"
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.submit') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

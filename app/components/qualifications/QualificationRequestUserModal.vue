<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = defineProps<{
    qualificationId: number;
}>();

const emit = defineEmits<{
    (e: 'updatedRequest', value?: QualificationRequest): void;
}>();

const { isOpen } = useOverlay();

const notifications = useNotificationsStore();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const schema = z.object({
    userComment: z.string().min(0).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    userComment: '',
});

async function createOrUpdateQualificationRequest(
    qualificationId: number,
    values: Schema,
): Promise<CreateOrUpdateQualificationRequestResponse> {
    try {
        const call = qualificationsQualificationsClient.createOrUpdateQualificationRequest({
            request: {
                qualificationId: qualificationId,
                userId: 0,
                userComment: values.userComment,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('updatedRequest', response.request);
        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateQualificationRequest(props.qualificationId, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <UModal>
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ $t('components.qualifications.request_modal.title') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="isOpen = false"
                        />
                    </div>
                </template>

                <div>
                    <UFormField class="flex-1" name="userComment" :label="$t('common.message')">
                        <UTextarea v-model="state.userComment" name="userComment" :placeholder="$t('common.message')" />
                    </UFormField>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="neutral" block @click="isOpen = false">
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

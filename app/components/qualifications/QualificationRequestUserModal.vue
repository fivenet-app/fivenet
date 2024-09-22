<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = defineProps<{
    qualificationId: string;
}>();

const emits = defineEmits<{
    (e: 'updatedRequest', value?: QualificationRequest): void;
}>();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const schema = z.object({
    userComment: z.string().min(0).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    userComment: '',
});

async function createOrUpdateQualificationRequest(
    qualificationId: string,
    values: Schema,
): Promise<CreateOrUpdateQualificationRequestResponse> {
    try {
        const call = getGRPCQualificationsClient().createOrUpdateQualificationRequest({
            request: {
                qualificationId: qualificationId,
                userId: 0,
                userComment: values.userComment,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('updatedRequest', response.request);
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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
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
                    <UFormGroup name="userComment" :label="$t('common.message')" class="flex-1">
                        <UTextarea v-model="state.userComment" name="userComment" :placeholder="$t('common.message')" />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.submit') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>

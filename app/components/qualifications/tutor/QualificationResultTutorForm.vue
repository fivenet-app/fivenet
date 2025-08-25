<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ExamGrading } from '~~/gen/ts/resources/qualifications/exam';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { resultStatusToBgColor } from '../helpers';

const props = withDefaults(
    defineProps<{
        qualificationId: number;
        userId?: number;
        resultId?: number;
        score?: number;
        viewOnly?: boolean;
        grading?: ExamGrading | undefined;
    }>(),
    {
        userId: undefined,
        resultId: undefined,
        score: undefined,
        viewOnly: false,
        grading: undefined,
    },
);

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { activeChar } = useAuth();

const completorStore = useCompletorStore();

const notifications = useNotificationsStore();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const availableStatus = [
    { status: ResultStatus.SUCCESSFUL },
    { status: ResultStatus.FAILED },
    { status: ResultStatus.PENDING },
];

const usersLoading = ref(false);
const selectedUser = ref<undefined | UserShort>(undefined);

const schema = z.object({
    status: z.nativeEnum(ResultStatus),
    score: z.coerce.number().min(0).max(1000),
    summary: z.string().max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: ResultStatus.SUCCESSFUL,
    score: props.score ?? 0,
    summary: '',
});

async function createOrUpdateQualificationResult(
    qualificationId: number,
    values: Schema,
): Promise<CreateOrUpdateQualificationResultResponse> {
    try {
        const call = qualificationsQualificationsClient.createOrUpdateQualificationResult({
            result: {
                id: props.resultId ?? 0,
                qualificationId: qualificationId,
                userId: props.userId ?? selectedUser.value?.userId ?? 0,
                status: values.status,
                score: values.score,
                summary: values.summary,
                creatorId: activeChar.value!.userId,
                creatorJob: activeChar.value!.job,
            },
            grading: props.grading,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('refresh');
        emit('close');

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(
    () => props.score,
    () => (state.score = props.score ?? 0),
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateQualificationResult(props.qualificationId, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl leading-6 font-semibold">
                        {{ $t('components.qualifications.result_modal.title') }}
                    </h3>

                    <UButton class="-my-1" color="neutral" variant="ghost" icon="i-mdi-window-close" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <slot />

                <template v-if="!viewOnly">
                    <UFormField v-if="userId === undefined" class="flex-1" name="selectedUser" :label="$t('common.citizen')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="selectedUser"
                                class="flex-1"
                                :searchable="
                                    async (q: string) => {
                                        usersLoading = true;
                                        const users = await completorStore.completeCitizens({
                                            search: q,
                                            userIds: selectedUser ? [selectedUser.userId] : [],
                                        });
                                        usersLoading = false;
                                        return users;
                                    }
                                "
                                searchable-lazy
                                :searchable-placeholder="$t('common.search_field')"
                                :search-attributes="['firstname', 'lastname']"
                                :placeholder="$t('common.citizen', 1)"
                                trailing
                                by="userId"
                                leading-icon="i-mdi-user"
                            >
                                <template #item-label>
                                    <template v-if="selectedUser">
                                        {{ usersToLabel([selectedUser]) }}
                                    </template>
                                </template>

                                <template #item="{ option: user }">
                                    {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField class="flex-1" name="status" :label="$t('common.status')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.status"
                                :items="availableStatus"
                                value-key="status"
                                :placeholder="$t('common.status')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #item-label>
                                    <span class="size-2 rounded-full" :class="resultStatusToBgColor(state.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.ResultStatus.${ResultStatus[state.status]}`)
                                    }}</span>
                                </template>

                                <template #item="{ option }">
                                    <span class="size-2 rounded-full" :class="resultStatusToBgColor(option.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.ResultStatus.${ResultStatus[option.status]}`)
                                    }}</span>
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.status')]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField class="flex-1" name="score" :label="$t('common.score')">
                        <UInput
                            v-model="state.score"
                            name="score"
                            type="number"
                            :min="0"
                            :max="100"
                            :step="0.1"
                            :placeholder="$t('common.score')"
                            :label="$t('common.score')"
                            trailing-icon="i-mdi-star-check"
                        />
                    </UFormField>

                    <UFormField class="flex-1" name="summary" :label="$t('common.summary')">
                        <UTextarea v-model="state.summary" name="summary" :rows="3" :placeholder="$t('common.summary')" />
                    </UFormField>
                </template>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="neutral" block @click="$emit('close')">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UButton v-if="!viewOnly" class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                        {{ $t('common.submit') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UForm>
</template>

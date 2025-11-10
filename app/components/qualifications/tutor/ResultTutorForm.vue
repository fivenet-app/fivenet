<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ExamGrading } from '~~/gen/ts/resources/qualifications/exam';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { resultStatusToBadgeColor } from '../helpers';

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
    (e: 'close', v: boolean): void;
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

const selectedUser = ref<undefined | UserShort>(undefined);

const schema = z.object({
    status: z.enum(ResultStatus),
    score: z.coerce.number().min(0).max(1000),
    summary: z.coerce.string().max(255),
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
        emit('close', false);

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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal>
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-2xl leading-6 font-semibold">
                    {{ $t('components.qualifications.result_modal.title') }}
                </h3>

                <UButton color="neutral" variant="ghost" icon="i-mdi-window-close" @click="$emit('close', false)" />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <slot />

                <template v-if="!viewOnly">
                    <UFormField v-if="userId === undefined" class="flex-1" name="selectedUser" :label="$t('common.citizen')">
                        <SelectMenu
                            v-model="selectedUser"
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeCitizens({
                                        search: q,
                                        userIds: selectedUser ? [selectedUser.userId] : [],
                                    })
                            "
                            searchable-key="completor-citizens"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['firstname', 'lastname']"
                            :placeholder="$t('common.citizen', 1)"
                            trailing
                            leading-icon="i-mdi-user"
                        >
                            <template v-if="selectedUser" #default>
                                {{ userToLabel(selectedUser) }}
                            </template>

                            <template #item-label="{ item }">
                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField class="flex-1" name="status" :label="$t('common.status')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.status"
                                :items="availableStatus"
                                value-key="status"
                                class="w-full"
                                :placeholder="$t('common.status')"
                                :search-input="{ placeholder: $t('common.search_field') }"
                            >
                                <template #default>
                                    <UBadge class="truncate" :color="resultStatusToBadgeColor(state.status)">
                                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[state.status]}`) }}
                                    </UBadge>
                                </template>

                                <template #item-label="{ item }">
                                    <UBadge class="truncate" :color="resultStatusToBadgeColor(item.status)">
                                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[item.status]}`) }}
                                    </UBadge>
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.status')]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField class="flex-1" name="score" :label="$t('common.score')">
                        <UInputNumber
                            v-model="state.score"
                            class="w-full"
                            name="score"
                            :min="0"
                            :max="100"
                            :step="0.1"
                            :placeholder="$t('common.score')"
                            :label="$t('common.score')"
                            trailing-icon="i-mdi-star-check"
                        />
                    </UFormField>

                    <UFormField class="flex-1" name="summary" :label="$t('common.summary')">
                        <UTextarea
                            v-model="state.summary"
                            name="summary"
                            :rows="3"
                            :placeholder="$t('common.summary')"
                            class="w-full"
                        />
                    </UFormField>
                </template>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    v-if="!viewOnly"
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.submit')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { resultStatusToBgColor } from '../helpers';

const props = withDefaults(
    defineProps<{
        qualificationId: string;
        userId?: number;
        resultId?: string;
        score?: number;
        viewOnly?: boolean;
    }>(),
    {
        userId: undefined,
        resultId: undefined,
        score: undefined,
        viewOnly: false,
    },
);

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { activeChar } = useAuth();

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

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
    qualificationId: string,
    values: Schema,
): Promise<CreateOrUpdateQualificationResultResponse> {
    try {
        const call = getGRPCQualificationsClient().createOrUpdateQualificationResult({
            result: {
                id: props.resultId ?? '0',
                qualificationId: qualificationId,
                userId: props.userId ?? selectedUser.value?.userId ?? 0,
                status: values.status,
                score: values.score,
                summary: values.summary,
                creatorId: activeChar.value!.userId,
                creatorJob: activeChar.value!.job,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('refresh');
        emits('close');

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
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.qualifications.result_modal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <slot />

                <template v-if="!viewOnly">
                    <UFormGroup v-if="userId === undefined" name="selectedUser" :label="$t('common.citizen')" class="flex-1">
                        <ClientOnly>
                            <USelectMenu
                                v-model="selectedUser"
                                :searchable="
                                    async (query: string) => {
                                        usersLoading = true;
                                        const users = await completorStore.completeCitizens({
                                            search: query,
                                        });
                                        usersLoading = false;
                                        return users;
                                    }
                                "
                                searchable-lazy
                                :searchable-placeholder="$t('common.search_field')"
                                :search-attributes="['firstname', 'lastname']"
                                class="flex-1"
                                :placeholder="$t('common.citizen', 1)"
                                trailing
                                by="userId"
                                leading-icon="i-mdi-user"
                            >
                                <template #label>
                                    <template v-if="selectedUser">
                                        {{ usersToLabel([selectedUser]) }}
                                    </template>
                                </template>
                                <template #option="{ option: user }">
                                    {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                </template>
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup name="status" :label="$t('common.status')" class="flex-1">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.status"
                                :options="availableStatus"
                                value-attribute="status"
                                :placeholder="$t('common.status')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
                                    <span class="size-2 rounded-full" :class="resultStatusToBgColor(state.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.ResultStatus.${ResultStatus[state.status]}`)
                                    }}</span>
                                </template>

                                <template #option="{ option }">
                                    <span class="size-2 rounded-full" :class="resultStatusToBgColor(option.status)" />
                                    <span class="truncate">{{
                                        $t(`enums.qualifications.ResultStatus.${ResultStatus[option.status]}`)
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

                    <UFormGroup name="score" :label="$t('common.score')" class="flex-1">
                        <UInput
                            v-model="state.score"
                            name="score"
                            type="number"
                            :min="0"
                            :max="100"
                            :placeholder="$t('common.score')"
                            :label="$t('common.score')"
                            trailing-icon="i-mdi-star-check"
                        />
                    </UFormGroup>

                    <UFormGroup name="summary" :label="$t('common.summary')" class="flex-1">
                        <UTextarea v-model="state.summary" name="summary" :rows="3" :placeholder="$t('common.summary')" />
                    </UFormGroup>
                </template>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="$emit('close')">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UButton v-if="!viewOnly" type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                        {{ $t('common.submit') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UForm>
</template>

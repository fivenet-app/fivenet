<script lang="ts" setup>
import { isPast } from 'date-fns';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { ExamQuestions, ExamUser } from '~~/gen/ts/resources/qualifications/exam';
import type { GetExamInfoResponse, TakeExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import ExamViewQuestions from './ExamViewQuestions.vue';

const props = defineProps<{
    qualificationId: number;
}>();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const { data, status, refresh, error } = useLazyAsyncData(`qualification-${props.qualificationId}-examinfo`, () =>
    getExamInfo(props.qualificationId),
);

async function getExamInfo(qualificationId: number): Promise<GetExamInfoResponse> {
    try {
        const call = qualificationsQualificationsClient.getExamInfo({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        examUser.value = response.examUser;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function takeExam(cancel = false): Promise<TakeExamResponse> {
    try {
        const call = qualificationsQualificationsClient.takeExam({
            qualificationId: props.qualificationId,
            cancel: cancel,
        });
        const { response } = await call;

        exam.value = response.exam;
        examUser.value = response.examUser;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const exam = ref<ExamQuestions | undefined>();
const examUser = ref<ExamUser | undefined>();

watch(data, async () => {
    if (data.value?.examUser?.endsAt !== undefined && data.value?.examUser?.endedAt === undefined) {
        await takeExam(false);
    }
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.exam.title')">
        <template #right>
            <PartialsBackButton :fallback-to="`/qualifications/${qualificationId}`" />
        </template>
    </UDashboardNavbar>

    <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.exam', 1)])" />
    <DataErrorBlock
        v-else-if="error"
        :title="$t('common.unable_to_load', [$t('common.exam', 1)])"
        :error="error"
        :retry="refresh"
    />
    <DataNoDataBlock
        v-else-if="!data"
        icon="i-mdi-account-school"
        :message="$t('common.not_found', [$t('common.qualification', 1)])"
    />

    <ExamViewQuestions
        v-else-if="exam && examUser && examUser?.endsAt"
        :qualification-id="qualificationId"
        :exam="exam"
        :exam-user="examUser"
        :qualification="data.qualification"
    />

    <template v-else>
        <UDashboardToolbar>
            <template #default>
                <div class="flex justify-between gap-2">
                    <div class="flex gap-2">
                        <UBadge v-if="data?.qualification?.examSettings?.time" class="inline-flex gap-1">
                            <UIcon class="size-4" name="i-mdi-clock" />
                            {{ $t('common.duration') }}: {{ fromDuration(data.qualification.examSettings.time) }}s
                        </UBadge>
                        <UBadge class="inline-flex gap-1">
                            <UIcon class="size-4" name="i-mdi-question-mark" />
                            {{ $t('common.count') }}: {{ data?.questionCount }}
                            {{ $t('common.question', data?.questionCount ?? 1) }}
                        </UBadge>
                    </div>
                    <div class="flex gap-2">
                        <UBadge v-if="data.examUser?.startedAt">
                            {{ $t('common.begins_at') }}
                            {{ $d(toDate(data.examUser?.startedAt), 'long') }}
                        </UBadge>
                        <UBadge v-if="data?.examUser?.endsAt">
                            {{ $t('common.ends_at') }}
                            {{ $d(toDate(data?.examUser?.endsAt), 'long') }}
                        </UBadge>
                    </div>
                </div>
            </template>
        </UDashboardToolbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <UCard :ui="{ rounded: '' }">
                <UAlert v-if="data?.examUser?.endedAt || isPast(toDate(data?.examUser?.endsAt))">
                    <h3 class="text-lg">
                        {{ $t('components.qualifications.exam_view.times_up') }}
                    </h3>
                </UAlert>
                <UButton
                    v-else-if="!data?.examUser?.endedAt"
                    class="w-full"
                    size="xl"
                    color="gray"
                    icon="i-mdi-play"
                    block
                    @click="takeExam(false)"
                >
                    {{ $t('components.qualifications.take_test') }}
                </UButton>
            </UCard>
        </UDashboardPanelContent>
    </template>
</template>

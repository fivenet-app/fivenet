<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { GetExamInfoResponse, TakeExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import ExamViewQuestions from './ExamViewQuestions.vue';
import type { ExamQuestions, ExamUser } from '~~/gen/ts/resources/qualifications/exam';

const props = defineProps<{
    qualificationId: string;
}>();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: string): Promise<GetExamInfoResponse> {
    try {
        const call = getGRPCQualificationsClient().getExamInfo({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function takeExam(cancel = false): Promise<TakeExamResponse> {
    try {
        const call = getGRPCQualificationsClient().takeExam({
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
    if (exam.value === undefined) {
        await takeExam(false);
    }
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.exam.title')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" :to="`/qualifications/${qualificationId}`">
                {{ $t('common.back') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.exam', 1)])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.exam', 1)])" :retry="refresh" />
    <DataNoDataBlock
        v-else-if="!data"
        icon="i-mdi-account-school"
        :message="$t('common.not_found', [$t('common.qualification', 1)])"
    />

    <ExamViewQuestions
        v-else-if="exam && data?.qualification && data?.user && examUser?.endsAt"
        :qualification-id="qualificationId"
        :exam="exam"
        :exam-user="data.user"
        :qualification="data.qualification"
    />

    <UCard v-else>
        <template #header>
            <div class="flex gap-2">
                <UBadge v-if="data?.qualification?.examSettings?.time" class="inline-flex gap-1">
                    <UIcon name="i-mdi-clock" class="size-4" />
                    {{ $t('common.duration') }}: {{ fromDuration(data.qualification.examSettings.time) }}
                </UBadge>
                <UBadge class="inline-flex gap-1">
                    <UIcon name="i-mdi-question-mark" class="size-4" />
                    {{ $t('common.count') }}: {{ data?.questionCount }} {{ $t('common.question', data?.questionCount ?? 1) }}
                </UBadge>
            </div>
        </template>

        <UButton
            v-if="!data?.user || !data.user.endedAt"
            color="gray"
            icon="i-mdi-play"
            block
            class="w-full"
            @click="takeExam(false)"
        >
            {{ $t('components.qualifications.take_test') }}
        </UButton>
    </UCard>
</template>

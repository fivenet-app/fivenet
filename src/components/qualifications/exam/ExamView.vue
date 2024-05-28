<script lang="ts" setup>
import type { GetExamResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = defineProps<{
    qualificationId: string;
}>();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: string): Promise<GetExamResponse> {
    try {
        const call = getGRPCQualificationsClient().getExam({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

// TODO
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.exam.title')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" :to="`/qualifications/${qualificationId}`">
                {{ $t('common.back') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <UCard>
        {{ data?.exam }}

        <h3>Settings</h3>
        {{ data?.exam?.settings }}

        <h3>Questions</h3>

        {{ data?.exam?.questions }}

        <!-- TODO -->
    </UCard>
</template>

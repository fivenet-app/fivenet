<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { GetQualificationResponse } from '~~/gen/ts/services/jobs/qualifications';

const props = defineProps<{
    id: string;
}>();

const { $grpc } = useNuxtApp();

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-qualification-${props.id}`, () => getQualification(props.id));

async function getQualification(qualificationId: string): Promise<GetQualificationResponse> {
    try {
        const call = $grpc.getJobsQualificationsClient().getQualification({
            qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 1)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!data?.qualification" />

            <template v-else>
                <!-- TODO -->
                <NuxtLink :to="{ name: 'jobs-qualifications-id-edit', params: { id: data.qualification.id } }">
                    {{ data.qualification.title }}
                </NuxtLink>
            </template>
        </div>
    </div>
</template>

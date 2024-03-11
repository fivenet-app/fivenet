<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import QualificationsListEntry from '~/components/jobs/qualifications/QualificationsListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationsResponse } from '~~/gen/ts/services/jobs/qualifications';

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-qualifications-${offset.value}`, () => listQualifications());

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = $grpc.getJobsQualificationsClient().listQualifications({
            pagination: {
                offset: offset.value,
            },
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
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data?.qualifications.length === 0" />

            <template v-else>
                <ul role="list" class="divide-y divide-gray-100">
                    <QualificationsListEntry
                        v-for="training in data?.qualifications"
                        :key="training.id"
                        :qualification="training"
                    />
                </ul>
            </template>
        </div>
    </div>
</template>

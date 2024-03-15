<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import QualificationsListEntry from '~/components/jobs/qualifications/QualificationsListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`qualifications-${offset.value}`, () => listQualifications());

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualifications({
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
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="border-b border-gray-200 bg-base-600 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                <div class="ml-4 mt-4">
                    <h3 class="text-base font-semibold leading-6 text-gray-200">
                        {{ $t('components.jobs.qualifications.all_qualifications') }}
                    </h3>
                    <p class="mt-1 text-sm text-gray-500"></p>
                </div>
                <div class="ml-4 mt-4 flex-shrink-0">
                    <NuxtLink
                        v-if="can('QualificationsService.CreateQualification')"
                        :to="{ name: 'jobs-qualifications-create' }"
                        class="relative inline-flex items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
                    >
                        {{ $t('components.jobs.qualifications.create_new_qualification') }}
                    </NuxtLink>
                </div>
            </div>
        </div>
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
                        v-for="qualification in data?.qualifications"
                        :key="qualification.id"
                        :qualification="qualification"
                    />
                </ul>
            </template>
        </div>
    </div>
</template>

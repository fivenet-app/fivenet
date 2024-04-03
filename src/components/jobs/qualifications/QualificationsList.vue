<script lang="ts" setup>
import { SchoolIcon } from 'mdi-vue3';
import QualificationsListEntry from '~/components/jobs/qualifications/QualificationsListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const { $grpc } = useNuxtApp();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`qualifications-${page.value}`, () => listQualifications());

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

watch(offset, async () => refresh());
</script>

<template>
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="border-b border-gray-200 bg-background px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                <div class="ml-4 mt-4">
                    <h3 class="text-base font-semibold leading-6 text-gray-200">
                        {{ $t('components.qualifications.all_qualifications') }}
                    </h3>
                </div>
                <div class="ml-4 mt-4 shrink-0">
                    <NuxtLink
                        v-if="can('QualificationsService.CreateQualification')"
                        :to="{ name: 'jobs-qualifications-create' }"
                        class="relative inline-flex items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                    >
                        {{ $t('components.qualifications.create_new_qualification') }}
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
            <DataNoDataBlock
                v-else-if="data?.qualifications.length === 0"
                :message="$t('common.not_found', [$t('common.qualifications', 2)])"
                icon="i-mdi-school"
            />

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
        <div class="border-t border-gray-200 bg-background px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex items-center">
                <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
                    <UPagination
                        v-model="page"
                        :page-count="data?.pagination?.pageSize ?? 0"
                        :total="data?.pagination?.totalCount ?? 0"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

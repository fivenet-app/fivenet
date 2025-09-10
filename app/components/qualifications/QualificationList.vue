<script lang="ts" setup>
import { z } from 'zod';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import QualificationListEntry from '~/components/qualifications/QualificationListEntry.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const schema = z.object({
    search: z.string().max(64).optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'abbreviation',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'abbreviation', desc: false }] }),
    page: pageNumberSchema,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const query = useSearchForm('qualifications_list', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    () => `qualifications-${JSON.stringify(query.sorting)}-${query.page}`,
    () => listQualifications(),
);

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualifications({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
            search: query.search,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });
</script>

<template>
    <UCard>
        <template #header>
            <div class="flex items-center justify-between gap-1">
                <h3 class="flex-1 text-2xl leading-6 font-semibold">
                    {{ $t('components.qualifications.all_qualifications') }}
                </h3>

                <UForm :schema="schema" :state="query" @submit="refresh">
                    <UFormField name="search">
                        <UInput
                            v-model="query.search"
                            type="text"
                            name="search"
                            :placeholder="$t('common.search')"
                            leading-icon="i-mdi-search"
                        />
                    </UFormField>
                </UForm>

                <SortButton v-model="query.sorting" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <div>
            <DataPendingBlock
                v-if="isRequestPending(status)"
                :message="$t('common.loading', [$t('common.qualifications', 2)])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.qualifications.length === 0"
                :message="$t('common.not_found', [$t('common.qualifications', 2)])"
                icon="i-mdi-school"
            />

            <ul v-else class="divide-y divide-default" role="list">
                <QualificationListEntry
                    v-for="qualification in data?.qualifications"
                    :key="qualification.id"
                    :qualification="qualification"
                />
            </ul>
        </div>

        <template #footer>
            <Pagination
                v-model="query.page"
                :pagination="data?.pagination"
                :status="status"
                :refresh="refresh"
                disable-border
            />
        </template>
    </UCard>
</template>

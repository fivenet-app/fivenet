<script lang="ts" setup>
import { z } from 'zod';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import QualificationsListEntry from '~/components/qualifications/QualificationsListEntry.vue';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const { $grpc } = useNuxtApp();

const schema = z.object({
    search: z.string().max(64).optional(),
    sort: z.custom<TableSortable>().default({
        column: 'abbreviation',
        direction: 'asc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('qualifications_list', schema);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualifications-${query.sort.column}:${query.sort.direction}-${query.page}`, () => listQualifications());

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = $grpc.qualifications.qualifications.listQualifications({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sort,
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
    <UCard
        :ui="{
            body: { padding: '' },
        }"
    >
        <template #header>
            <div class="flex items-center justify-between gap-1">
                <h3 class="flex-1 text-2xl font-semibold leading-6">
                    {{ $t('components.qualifications.all_qualifications') }}
                </h3>

                <UForm :schema="schema" :state="query" @submit="refresh">
                    <UFormGroup name="search">
                        <UInput
                            v-model="query.search"
                            type="text"
                            name="search"
                            :placeholder="$t('common.search')"
                            leading-icon="i-mdi-search"
                        />
                    </UFormGroup>
                </UForm>

                <SortButton v-model="query.sort" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <div>
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.qualifications', 2)])" />
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

            <ul v-else class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                <QualificationsListEntry
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
                :loading="loading"
                :refresh="refresh"
                disable-border
            />
        </template>
    </UCard>
</template>

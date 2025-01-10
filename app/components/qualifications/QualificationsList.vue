<script lang="ts" setup>
import { z } from 'zod';
import Pagination from '~/components/partials/Pagination.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import QualificationsListEntry from '~/components/qualifications/QualificationsListEntry.vue';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import SortButton from '../partials/SortButton.vue';

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'abbreviation',
    direction: 'asc',
});

const schema = z.object({
    search: z.string().max(64).optional(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualifications-${sort.value.column}:${sort.value.direction}-${page.value}`, () => listQualifications(), {
    watch: [sort],
});

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = getGRPCQualificationsClient().listQualifications({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            search: query.search,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
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

                <SortButton v-model="sort" :fields="[{ label: $t('common.id'), value: 'id' }]" />
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

            <ul v-else role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
                <QualificationsListEntry
                    v-for="qualification in data?.qualifications"
                    :key="qualification.id"
                    :qualification="qualification"
                />
            </ul>
        </div>

        <template #footer>
            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" disable-border />
        </template>
    </UCard>
</template>

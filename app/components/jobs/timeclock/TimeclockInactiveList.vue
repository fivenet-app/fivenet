<script lang="ts" setup>
import type { Form } from '#ui/types';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { Perms } from '~~/gen/ts/perms';
import type { ListInactiveEmployeesResponse } from '~~/gen/ts/services/jobs/timeclock';

const { t } = useI18n();

const { can } = useAuth();

const form = ref<null | Form<Schema>>();

const schema = z.object({
    days: z.coerce.number().min(1).max(31),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    days: 14,
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'rank',
    direction: 'asc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-timeclock-inactive-${sort.value.column}:${sort.value.direction}-${page.value}-${state.days}`,
    () => listInactiveEmployees(state),
    {
        watch: [sort],
    },
);

async function listInactiveEmployees(values: Schema): Promise<ListInactiveEmployeesResponse> {
    try {
        const call = getGRPCJobsTimeclockClient().listInactiveEmployees({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            days: values.days,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(
    state,
    async () => {
        const valid = await form.value?.validate();
        if (valid) {
            refresh();
        }
    },
    { debounce: 200, maxWait: 1250 },
);
watch(offset, async () => refresh());

const columns = [
    {
        key: 'name',
        label: t('common.name'),
        sortable: true,
    },
    {
        key: 'rank',
        label: t('common.rank', 1),
        sortable: true,
    },
    {
        key: 'phoneNumber',
        label: t('common.phone_number'),
    },
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
        permission: 'JobsService.GetColleague' as Perms,
    },
].filter((c) => c.permission === undefined || can(c.permission).value);
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <div class="flex w-full flex-col">
                <UButton
                    v-if="can('JobsTimeclockService.ListTimeclock').value"
                    :to="{ name: 'jobs-timeclock' }"
                    icon="i-mdi-arrow-left"
                    class="mb-2 place-self-end"
                >
                    {{ $t('common.timeclock') }}
                </UButton>

                <UForm ref="form" :schema="schema" :state="state" class="flex w-full flex-row gap-2" @submit="refresh()">
                    <UFormGroup name="days" :label="$t('common.time_ago.day', 2)" class="flex-1">
                        <UInput
                            v-model="state.days"
                            name="days"
                            type="number"
                            :min="1"
                            :max="31"
                            :placeholder="$t('common.time_ago.day', 2)"
                        />
                    </UFormGroup>
                </UForm>
            </div>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.colleague', 2)])" :retry="refresh" />
    <UTable
        v-else
        v-model:sort="sort"
        :loading="loading"
        :columns="columns"
        :rows="data?.colleagues"
        :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
        sort-mode="manual"
        class="flex-1"
    >
        <template #name-data="{ row: colleague }">
            <div class="inline-flex items-center gap-1 text-gray-900 dark:text-white">
                <ProfilePictureImg
                    :src="colleague.avatar?.url"
                    :name="`${colleague.firstname} ${colleague.lastname}`"
                    size="sm"
                    :enable-popup="true"
                />

                <span> {{ colleague.firstname }} {{ colleague.lastname }} </span>
            </div>
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.job_grade') }}</dt>
                <dd class="mt-1 truncate">
                    {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                </dd>
            </dl>
        </template>
        <template #rank-data="{ row: colleague }">
            {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
        </template>
        <template #phoneNumber-data="{ row: colleague }">
            <PhoneNumberBlock :number="colleague.phoneNumber" />
        </template>
        <template #actions-data="{ row: colleague }">
            <div :key="colleague.id">
                <UTooltip v-if="checkIfCanAccessColleague(colleague, 'JobsService.GetColleague')" :text="$t('common.show')">
                    <ULink
                        icon="i-mdi-eye"
                        :to="{
                            name: 'jobs-colleagues-id-info',
                            params: { id: colleague.userId ?? 0 },
                        }"
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>

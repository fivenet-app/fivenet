<script lang="ts" setup>
import type { Form } from '#ui/types';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { getJobsTimeclockClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import type { ListInactiveEmployeesResponse } from '~~/gen/ts/services/jobs/timeclock';
import ColleagueName from '../colleagues/ColleagueName.vue';

const { t } = useI18n();

const { can } = useAuth();

const jobsTimeclockClient = await getJobsTimeclockClient();

const schema = z.object({
    days: z.coerce.number().min(1).max(31),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    days: 14,
});

const page = useRouteQuery('page', '1', { transform: Number });

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'rank',
    direction: 'asc',
});

const { data, status, refresh, error } = useLazyAsyncData(
    `jobs-timeclock-inactive-${sort.value.column}:${sort.value.direction}-${page.value}-${state.days}`,
    () => listInactiveEmployees(state),
    {
        watch: [sort],
    },
);

async function listInactiveEmployees(values: Schema): Promise<ListInactiveEmployeesResponse> {
    try {
        const call = jobsTimeclockClient.listInactiveEmployees({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
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

const form = ref<null | Form<Schema>>();

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
        key: 'absence',
        label: t('common.absence_date'),
    },
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
        permission: 'jobs.JobsService/GetColleague' as Perms,
    },
].filter((c) => c.permission === undefined || can(c.permission).value);

const { game } = useAppConfig();
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <div class="flex w-full flex-col">
                <UButton
                    v-if="can('jobs.TimeclockService/ListTimeclock').value"
                    class="mb-2 place-self-end"
                    :to="{ name: 'jobs-timeclock' }"
                    icon="i-mdi-arrow-left"
                >
                    {{ $t('common.timeclock') }}
                </UButton>

                <UForm ref="form" class="flex w-full flex-row gap-2" :schema="schema" :state="state" @submit="refresh()">
                    <UFormGroup class="flex-1" name="days" :label="$t('common.time_ago.day', 2)">
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

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
        :error="error"
        :retry="refresh"
    />
    <UTable
        v-else
        v-model:sort="sort"
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :rows="data?.colleagues"
        :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
        sort-mode="manual"
    >
        <template #name-data="{ row: colleague }">
            <div class="inline-flex items-center gap-1 text-gray-900 dark:text-white">
                <ProfilePictureImg
                    :src="colleague.avatar"
                    :name="`${colleague.firstname} ${colleague.lastname}`"
                    size="sm"
                    :enable-popup="true"
                />

                <ColleagueName :colleague="colleague" />
            </div>
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.job_grade') }}</dt>
                <dd class="mt-1 truncate">
                    {{ colleague.jobGradeLabel }}
                    <template v-if="colleague.job !== game.unemployedJobName"> ({{ colleague.jobGrade }}) </template>
                </dd>
            </dl>
        </template>

        <template #rank-data="{ row: colleague }">
            {{ colleague.jobGradeLabel }}
            <template v-if="colleague.job !== game.unemployedJobName"> ({{ colleague.jobGrade }}) </template>
        </template>

        <template #phoneNumber-data="{ row: colleague }">
            <PhoneNumberBlock :number="colleague.phoneNumber" />
        </template>

        <template #absence-data="{ row: colleague }">
            <dl v-if="colleague.props?.absenceEnd" class="font-normal">
                <dd class="truncate">
                    {{ $t('common.from') }}:
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                </dd>
                <dd class="truncate">
                    {{ $t('common.to') }}: <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </dd>
            </dl>
        </template>

        <template v-if="can('jobs.JobsService/GetColleague').value" #actions-data="{ row: colleague }">
            <div :key="colleague.id">
                <UTooltip
                    v-if="checkIfCanAccessColleague(colleague, 'jobs.JobsService/GetColleague')"
                    :text="$t('common.show')"
                >
                    <UButton
                        variant="link"
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

    <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>

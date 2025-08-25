<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { Form, TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { getJobsTimeclockClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { ListInactiveEmployeesResponse } from '~~/gen/ts/services/jobs/timeclock';
import ColleagueName from '../colleagues/ColleagueName.vue';

const appConfig = useAppConfig();

const { t } = useI18n();

const { can } = useAuth();

const jobsTimeclockClient = await getJobsTimeclockClient();

const schema = z.object({
    days: z.coerce.number().min(1).max(31),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'rank',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'rank', desc: false }] }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    days: 14,
    sorting: {
        columns: [
            {
                id: 'rank',
                desc: false,
            },
        ],
    },
});

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, refresh, error } = useLazyAsyncData(
    () => `jobs-timeclock-inactive-${JSON.stringify(state.sorting)}-${page.value}-${state.days}`,
    () => listInactiveEmployees(state),
);

async function listInactiveEmployees(values: Schema): Promise<ListInactiveEmployeesResponse> {
    try {
        const call = jobsTimeclockClient.listInactiveEmployees({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            sort: values.sorting,
            days: values.days,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const formRef = ref<null | Form<Schema>>();

watchDebounced(
    state,
    async () => {
        const valid = await formRef.value?.validate();
        if (valid) {
            refresh();
        }
    },
    { debounce: 200, maxWait: 1250 },
);

const columns = computed(() =>
    (
        [
            {
                accessorKey: 'name',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.name'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) =>
                    h('div', { class: 'inline-flex items-center gap-1 text-highlighted' }, [
                        h(ProfilePictureImg, {
                            src: row.original.avatar,
                            name: `${row.original.firstname} ${row.original.lastname}`,
                            size: 'sm',
                            enablePopup: true,
                        }),
                        h(ColleagueName, { colleague: row.original }),
                    ]),
            },
            {
                accessorKey: 'jobGrade',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.rank', 1),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) =>
                    h('span', null, [
                        row.original.jobGradeLabel,
                        row.original.job !== game.unemployedJobName ? h('span', null, ` (${row.original.jobGrade})`) : null,
                    ]),
            },
            {
                accessorKey: 'phoneNumber',
                header: t('common.phone_number'),
                cell: ({ row }) => h(PhoneNumberBlock, { number: row.original.phoneNumber }),
            },
            {
                accessorKey: 'absence',
                header: t('common.absence_date'),
                cell: ({ row }) =>
                    row.original.props?.absenceEnd
                        ? h('dl', { class: 'font-normal' }, [
                              h('dd', { class: 'truncate' }, [
                                  `${t('common.from')}:`,
                                  h(GenericTime, { value: row.original.props?.absenceBegin, type: 'date' }),
                              ]),
                              h('dd', { class: 'truncate' }, [
                                  `${t('common.to')}:`,
                                  h(GenericTime, { value: row.original.props?.absenceEnd, type: 'date' }),
                              ]),
                          ])
                        : null,
            },
            {
                accessorKey: 'dateofbirth',
                header: t('common.date_of_birth'),
                cell: ({ row }) => row.original.dateofbirth,
            },
            can('jobs.JobsService/GetColleague').value
                ? {
                      accessorKey: 'actions',
                      header: t('common.action', 2),
                      cell: ({ row }) =>
                          checkIfCanAccessColleague(row.original, 'jobs.JobsService/GetColleague')
                              ? h(UTooltip, { text: t('common.show') }, () =>
                                    h(UButton, {
                                        variant: 'link',
                                        icon: 'i-mdi-eye',
                                        to: {
                                            name: 'jobs-colleagues-id-info',
                                            params: { id: row.original.userId ?? 0 },
                                        },
                                    }),
                                )
                              : null,
                  }
                : undefined,
        ] as TableColumn<Colleague>[]
    ).filter((c) => c !== undefined),
);

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

                <UForm ref="formRef" class="flex w-full flex-row gap-2" :schema="schema" :state="state" @submit="refresh()">
                    <UFormField class="flex-1" name="days" :label="$t('common.time_ago.day', 2)">
                        <UInput
                            v-model="state.days"
                            name="days"
                            type="number"
                            :min="1"
                            :max="31"
                            :placeholder="$t('common.time_ago.day', 2)"
                        />
                    </UFormField>
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
        v-model:sorting="state.sorting.columns"
        :columns="columns"
        :data="data?.colleagues"
        :loading="isRequestPending(status)"
        :empty="$t('common.not_found', [$t('common.colleague', 2)])"
        :pagination-options="{ manualPagination: true }"
        :sorting-options="{ manualSorting: true }"
    />

    <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>

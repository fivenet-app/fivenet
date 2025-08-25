<script lang="ts" setup>
import { useAppConfig } from '#app';
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { isFuture } from 'date-fns';
import { h } from 'vue';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useSettingsStore } from '~/stores/settings';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { Label } from '~~/gen/ts/resources/jobs/labels';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { GetColleagueLabelsResponse, ListColleaguesResponse } from '~~/gen/ts/services/jobs/jobs';
import ColleagueLabelStatsModal from './ColleagueLabelStatsModal.vue';
import ColleagueName from './ColleagueName.vue';
import JobLabelsModal from './JobLabelsModal.vue';
import SelfServicePropsAbsenceDateModal from './SelfServicePropsAbsenceDateModal.vue';

const { t } = useI18n();

const overlay = useOverlay();

const { attr, can, activeChar } = useAuth();

const jobsJobsClient = await getJobsJobsClient();

const schema = z.object({
    name: z.string().max(50).default(''),
    absent: z.coerce.boolean().default(false),
    labels: z.coerce.number().array().max(3).default([]),
    namePrefix: z.string().max(12).optional(),
    nameSuffix: z.string().max(12).optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_colleagues', schema);

const settingsStore = useSettingsStore();
const { jobsService } = storeToRefs(settingsStore);

const { data, status, refresh, error } = useLazyAsyncData(
    () =>
        `jobs-colleagues-${JSON.stringify(query.sorting)}-${query.page}-${query.name}-${query.absent}-${query.labels.join(',')}-${query.namePrefix}-${query.nameSuffix}`,
    () => listColleagues(),
);

async function listColleagues(): Promise<ListColleaguesResponse> {
    try {
        const call = jobsJobsClient.listColleagues({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
            search: query.name,
            userIds: [],
            absent: query.absent,
            labelIds: query.labels,
            namePrefix: query.namePrefix,
            nameSuffix: query.nameSuffix,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function getColleagueLabels(search?: string): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await jobsJobsClient.getColleagueLabels({
            search: search,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

function updateAbsenceDates(value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void {
    const colleague = data.value?.colleagues.find((c) => c.userId === value.userId);
    if (colleague === undefined) {
        return;
    }

    if (colleague.props === undefined) {
        colleague.props = {
            userId: colleague.userId,
            job: colleague.job,
            absenceBegin: value.absenceBegin,
            absenceEnd: value.absenceEnd,
        };
    } else {
        colleague.props.absenceBegin = value.absenceBegin;
        colleague.props.absenceEnd = value.absenceEnd;
    }
}

function toggleLabelInSearch(label: Label): void {
    const idx = query.labels.findIndex((l) => l === label.id);
    if (idx > -1) {
        query.labels.splice(idx, 1);
    } else {
        query.labels.push(label.id);
    }
}

const appConfig = useAppConfig();

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
            },
            {
                accessorKey: 'jobGrade',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.rank'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) => h('span', {}, row.original.jobGradeLabel),
            },
            {
                accessorKey: 'absence',
                header: t('common.absence_date'),
            },
            {
                accessorKey: 'phoneNumber',
                header: t('common.phone_number'),
                cell: ({ row }) => h(PhoneNumberBlock, { number: row.original.phoneNumber }),
            },
            {
                accessorKey: 'email',
                header: t('common.mail'),
                cell: ({ row }) =>
                    h(EmailInfoPopover, {
                        email: row.original.email,
                        variant: 'link',
                        color: 'primary',
                        truncate: true,
                        trailing: false,
                        hideNaText: true,
                    }),
            },
            {
                accessorKey: 'dateofbirth',
                header: t('common.date_of_birth'),
                cell: ({ row }) => h('span', {}, row.original.dateofbirth),
            },
            can(['jobs.JobsService/GetColleague', 'jobs.JobsService/SetColleagueProps']).value
                ? {
                      accessorKey: 'actions',
                      header: t('common.action', 2),
                      cell: ({ row }) =>
                          h('div', { class: 'flex flex-col justify-end md:flex-row' }, [
                              h(
                                  UTooltip,
                                  {
                                      text: t('components.jobs.self_service.set_absence_date'),
                                      vIf:
                                          canDo.value.setJobsUserProps &&
                                          (row.original.userId === activeChar.value!.userId ||
                                              attr('jobs.JobsService/SetColleagueProps', 'Types', 'AbsenceDate').value) &&
                                          checkIfCanAccessColleague(row.original, 'jobs.JobsService/SetColleagueProps'),
                                  },
                                  [
                                      h(UButton, {
                                          variant: 'link',
                                          icon: 'i-mdi-island',
                                          onClick: () =>
                                              selfServicePropsAbsenceDateModal.open({
                                                  userId: row.original.userId,
                                                  'onUpdate:absenceDates': ($event) => updateAbsenceDates($event),
                                              }),
                                      }),
                                  ],
                              ),
                              h(
                                  UTooltip,
                                  {
                                      text: t('common.show'),
                                      vIf:
                                          canDo.value.getColleague &&
                                          checkIfCanAccessColleague(row.original, 'jobs.JobsService/GetColleague'),
                                  },
                                  [
                                      h(UButton, {
                                          variant: 'link',
                                          icon: 'i-mdi-eye',
                                          to: {
                                              name: 'jobs-colleagues-id-info',
                                              params: { id: row.original.userId ?? 0 },
                                          },
                                      }),
                                  ],
                              ),
                          ]),
                  }
                : undefined,
        ] as TableColumn<Colleague>[]
    ).filter((c) => c !== undefined),
);

const canDo = computed(() => ({
    getColleague: can('jobs.JobsService/GetColleague').value,
    setJobsUserProps: can('jobs.JobsService/SetColleagueProps').value,
}));

const { game } = useAppConfig();

const selfServicePropsAbsenceDateModal = overlay.create(SelfServicePropsAbsenceDateModal);
const jobLabelsModal = overlay.create(JobLabelsModal);
const colleagueLabelStatsModal = overlay.create(ColleagueLabelStatsModal);

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
            <div class="flex gap-2">
                <UFormField class="flex-1" name="name" :label="$t('common.search')">
                    <UInput
                        ref="input"
                        v-model="query.name"
                        type="text"
                        name="name"
                        :placeholder="$t('common.name')"
                        block
                        leading-icon="i-mdi-search"
                        @keydown.esc="$event.target.blur()"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormField>

                <UFormField
                    class="flex flex-initial flex-col"
                    name="absent"
                    :label="$t('common.absent')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <USwitch v-model="query.absent" />
                    </div>
                </UFormField>

                <UFormField v-if="jobsService.cardView" :label="$t('common.sort_by')">
                    <SortButton
                        v-model="query.sorting"
                        :fields="[
                            { label: $t('common.rank'), value: 'rank' },
                            { label: $t('common.name'), value: 'name' },
                        ]"
                    />
                </UFormField>

                <UFormField
                    v-if="
                        can('jobs.JobsService/ManageLabels').value ||
                        attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value
                    "
                    label="&nbsp"
                    :ui="{ container: 'inline-flex gap-1' }"
                >
                    <UButton
                        v-if="can('jobs.JobsService/ManageLabels').value"
                        :label="$t('common.label', 2)"
                        icon="i-mdi-tag"
                        @click="jobLabelsModal.open({})"
                    />

                    <UTooltip v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value" :text="$t('common.stats')">
                        <UButton icon="i-mdi-chart-donut" color="neutral" @click="colleagueLabelStatsModal.open({})" />
                    </UTooltip>
                </UFormField>
            </div>

            <UAccordion
                class="mt-2"
                color="neutral"
                variant="soft"
                size="sm"
                :items="[{ label: $t('common.advanced_search'), slot: 'search' as const }]"
            >
                <template #search>
                    <div class="flex flex-row flex-wrap gap-2">
                        <UFormField
                            v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value"
                            class="flex flex-1 flex-col"
                            name="labels"
                            :label="$t('common.label', 2)"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.labels"
                                    class="flex-1"
                                    multiple
                                    :searchable="async (q: string) => (await getColleagueLabels(q))?.labels ?? []"
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['name']"
                                    option-attribute="name"
                                    clear-search-on-close
                                    value-key="id"
                                >
                                    <template #item-label="{ item }">
                                        <span v-if="item.length" class="inline-flex flex-wrap gap-1 truncate">
                                            <UBadge
                                                v-for="label in item"
                                                :key="label.id"
                                                class="truncate"
                                                :class="isColorBright(label.color) ? 'text-black!' : 'text-white!'"
                                                :style="{ backgroundColor: label.color }"
                                                :label="label.name"
                                            />
                                        </span>
                                        <span v-else>&nbsp;</span>
                                    </template>

                                    <template #item="{ item }">
                                        <UBadge
                                            class="truncate"
                                            :class="isColorBright(item.color) ? 'text-black!' : 'text-white!'"
                                            :style="{ backgroundColor: item.color }"
                                        >
                                            {{ item.name }}
                                        </UBadge>
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.label', 2)]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField
                            class="flex flex-col"
                            name="namePrefix"
                            :label="$t('common.prefix')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.namePrefix" type="text" />
                        </UFormField>

                        <UFormField
                            class="flex flex-col"
                            name="nameSuffix"
                            :label="$t('common.suffix')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.nameSuffix" type="text" />
                        </UFormField>

                        <UFormField
                            class="flex flex-initial flex-col"
                            name="cards"
                            :label="$t('common.card_view')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <div class="flex flex-1 items-center">
                                <USwitch v-model="jobsService.cardView" />
                            </div>
                        </UFormField>
                    </div>
                </template>
            </UAccordion>
        </UForm>
    </UDashboardToolbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
        :error="error"
        :retry="refresh"
    />
    <template v-else>
        <UTable
            v-if="!jobsService.cardView"
            v-model:sorting="query.sorting.columns"
            :loading="isRequestPending(status)"
            :columns="columns"
            :data="data?.colleagues"
            :empty="$t('common.not_found', [$t('common.colleague', 2)])"
            :sorting-options="{ manualSorting: true }"
            :pagination-options="{ manualPagination: true }"
        >
            <template #name-cell="{ row }">
                <div class="inline-flex items-center text-highlighted">
                    <ProfilePictureImg
                        class="mr-2"
                        :src="row.original?.avatar"
                        :name="`${row.original.firstname} ${row.original.lastname}`"
                        size="sm"
                        :enable-popup="true"
                        :alt="$t('common.avatar')"
                    />

                    <ColleagueName :colleague="row.original" />
                </div>
            </template>

            <template #absence-cell="{ row }">
                <dl
                    v-if="row.original.props?.absenceEnd && isFuture(toDate(row.original.props?.absenceEnd))"
                    class="font-normal"
                >
                    <dd class="truncate">
                        {{ $t('common.from') }}:
                        <GenericTime :value="row.original.props?.absenceBegin" type="date" />
                    </dd>
                    <dd class="truncate">
                        {{ $t('common.to') }}: <GenericTime :value="row.original.props?.absenceEnd" type="date" />
                    </dd>
                </dl>
            </template>
        </UTable>

        <div v-else class="relative flex-1 overflow-x-auto">
            <UPageGrid
                :ui="{
                    wrapper: 'grid grid-cols-1 p-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6',
                }"
            >
                <UPageCard v-for="colleague in data?.colleagues" :key="colleague.userId">
                    <template #title>
                        <ColleagueName :colleague="colleague" />
                    </template>

                    <template #header>
                        <div class="flex items-center justify-center overflow-hidden">
                            <ProfilePictureImg
                                :src="colleague?.avatar"
                                :name="`${colleague.firstname} ${colleague.lastname}`"
                                size="3xl"
                                :enable-popup="true"
                                :alt="$t('common.avatar')"
                                :rounded="false"
                                img-class="h-40 w-40 md:h-56 md:w-56 xl:h-64 xl:w-64 max-w-full"
                            />
                        </div>
                    </template>

                    <template #description>
                        <div class="flex flex-col gap-1 truncate">
                            <span>
                                {{ colleague.jobGradeLabel }}
                                <template v-if="colleague.job !== game.unemployedJobName">
                                    ({{ colleague.jobGrade }})
                                </template>
                            </span>

                            <div>
                                <PhoneNumberBlock :number="colleague.phoneNumber" />
                            </div>

                            <span class="inline-flex items-center gap-1">
                                <UIcon class="h-5 w-5 shrink-0" name="i-mdi-birthday-cake" />

                                <span>{{ colleague.dateofbirth }}</span>
                            </span>

                            <span class="flex items-center gap-1">
                                <UIcon class="h-5 w-5 shrink-0" name="i-mdi-email" />

                                <EmailInfoPopover :email="colleague.email" variant="link" :trailing="false" />
                            </span>

                            <div
                                v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value"
                                class="flex flex-row gap-1"
                            >
                                <UIcon class="h-5 w-5 shrink-0" name="i-mdi-tag" />

                                <span v-if="!colleague.props?.labels?.list.length">
                                    {{ $t('common.none', [$t('common.label', 2)]) }}
                                </span>
                                <div v-else class="flex max-w-full flex-row flex-wrap gap-1">
                                    <UButton
                                        v-for="label in colleague.props?.labels?.list"
                                        :key="label.name"
                                        class="justify-between gap-2"
                                        :class="isColorBright(hexToRgb(label.color, RGBBlack)!) ? 'text-black!' : 'text-white!'"
                                        :style="{ backgroundColor: label.color }"
                                        size="xs"
                                        @click="toggleLabelInSearch(label)"
                                    >
                                        <span class="truncate">
                                            {{ label.name }}
                                        </span>
                                    </UButton>
                                </div>
                            </div>

                            <span
                                v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))"
                                class="inline-flex items-center gap-1"
                            >
                                <UIcon class="size-5" name="i-mdi-island" />
                                <GenericTime :value="colleague.props?.absenceBegin" type="shortDate" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                            </span>
                        </div>
                    </template>

                    <template
                        v-if="
                            (canDo.setJobsUserProps &&
                                (colleague.userId === activeChar!.userId ||
                                    attr('jobs.JobsService/SetColleagueProps', 'Types', 'AbsenceDate').value) &&
                                checkIfCanAccessColleague(colleague, 'jobs.JobsService/SetColleagueProps')) ||
                            (canDo.getColleague && checkIfCanAccessColleague(colleague, 'jobs.JobsService/GetColleague'))
                        "
                        #footer
                    >
                        <UButtonGroup class="inline-flex w-full">
                            <UTooltip
                                v-if="
                                    canDo.setJobsUserProps &&
                                    (colleague.userId === activeChar!.userId ||
                                        attr('jobs.JobsService/SetColleagueProps', 'Types', 'AbsenceDate').value) &&
                                    checkIfCanAccessColleague(colleague, 'jobs.JobsService/SetColleagueProps')
                                "
                                class="flex-1"
                                :text="$t('components.jobs.self_service.set_absence_date')"
                            >
                                <UButton
                                    :label="$t('components.jobs.self_service.set_absence_date')"
                                    icon="i-mdi-island"
                                    block
                                    truncate
                                    @click="
                                        selfServicePropsAbsenceDateModal.open({
                                            userId: colleague.userId,
                                            'onUpdate:absenceDates': ($event) => updateAbsenceDates($event),
                                        })
                                    "
                                />
                            </UTooltip>

                            <UTooltip
                                v-if="
                                    canDo.getColleague && checkIfCanAccessColleague(colleague, 'jobs.JobsService/GetColleague')
                                "
                                class="flex-1"
                                :text="$t('common.show')"
                            >
                                <UButton
                                    :label="$t('common.show')"
                                    icon="i-mdi-eye"
                                    block
                                    truncate
                                    :to="{
                                        name: 'jobs-colleagues-id-info',
                                        params: { id: colleague.userId ?? 0 },
                                    }"
                                />
                            </UTooltip>
                        </UButtonGroup>
                    </template>
                </UPageCard>
            </UPageGrid>
        </div>
    </template>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>

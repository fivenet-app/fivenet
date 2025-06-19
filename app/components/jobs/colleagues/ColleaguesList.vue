<script lang="ts" setup>
import type { TableColumn } from '#ui/types';
import { isFuture } from 'date-fns';
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
import type { Label } from '~~/gen/ts/resources/jobs/labels';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { GetColleagueLabelsResponse, ListColleaguesResponse } from '~~/gen/ts/services/jobs/jobs';
import ColleagueName from './ColleagueName.vue';
import ColleaguesLabelStatsModal from './ColleaguesLabelStatsModal.vue';
import JobsLabelsModal from './JobsLabelsModal.vue';
import SelfServicePropsAbsenceDateModal from './SelfServicePropsAbsenceDateModal.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const { attr, can, activeChar } = useAuth();

const schema = z.object({
    name: z.string().max(50).default(''),
    absent: z.coerce.boolean().default(false),
    labels: z.coerce.number().array().max(3).default([]),
    namePrefix: z.string().max(12).optional(),
    nameSuffix: z.string().max(12).optional(),
    sort: z.custom<TableSortable>().default({
        column: 'rank',
        direction: 'asc',
    }),
    page: z.coerce.number().min(1).default(1),
});

const query = useSearchForm('jobs_colleagues', schema);

const settingsStore = useSettingsStore();
const { jobsService } = storeToRefs(settingsStore);

const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (query.page - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-colleagues-${query.sort.column}:${query.sort.direction}-${query.page}-${query.name}-${query.absent}-${query.labels.join(',')}-${query.namePrefix}-${query.nameSuffix}`,
    () => listColleagues(),
    {
        transform: (input) => ({ pagination: input.pagination, colleagues: wrapRows(input?.colleagues, columns) }),
    },
);

async function listColleagues(): Promise<ListColleaguesResponse> {
    try {
        const call = $grpc.jobs.jobs.listColleagues({
            pagination: {
                offset: offset.value,
            },
            sort: query.sort,
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
        const { response } = await $grpc.jobs.jobs.getColleagueLabels({
            search: search,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
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

const columns = [
    {
        key: 'name',
        label: t('common.name'),
        sortable: true,
    },
    {
        key: 'rank',
        label: t('common.rank'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
        sortable: true,
    },
    {
        key: 'absence',
        label: t('common.absence_date'),
    },
    {
        key: 'phoneNumber',
        label: t('common.phone_number'),
    },
    {
        key: 'email',
        label: t('common.mail'),
    },
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    can(['jobs.JobsService/GetColleague', 'jobs.JobsService/SetColleagueProps']).value
        ? {
              key: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined) as TableColumn[];

const canDo = computed(() => ({
    getColleague: can('jobs.JobsService/GetColleague').value,
    setJobsUerProps: can('jobs.JobsService/SetColleagueProps').value,
}));

const { game } = useAppConfig();

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
            <div class="flex gap-2">
                <UFormGroup class="flex-1" name="name" :label="$t('common.search')">
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
                </UFormGroup>

                <UFormGroup
                    class="flex flex-initial flex-col"
                    name="absent"
                    :label="$t('common.absent')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <UToggle v-model="query.absent">
                            <span class="sr-only">
                                {{ $t('common.absent') }}
                            </span>
                        </UToggle>
                    </div>
                </UFormGroup>

                <UFormGroup v-if="jobsService.cardView" :label="$t('common.sort_by')">
                    <SortButton
                        v-model="query.sort"
                        :fields="[
                            { label: $t('common.rank'), value: 'rank' },
                            { label: $t('common.name'), value: 'name' },
                        ]"
                    />
                </UFormGroup>

                <UFormGroup
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
                        @click="modal.open(JobsLabelsModal, {})"
                    />

                    <UTooltip v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value" :text="$t('common.stats')">
                        <UButton icon="i-mdi-chart-donut" color="white" @click="modal.open(ColleaguesLabelStatsModal, {})" />
                    </UTooltip>
                </UFormGroup>
            </div>

            <UAccordion
                class="mt-2"
                color="white"
                variant="soft"
                size="sm"
                :items="[{ label: $t('common.advanced_search'), slot: 'search' }]"
            >
                <template #search>
                    <div class="flex flex-row flex-wrap gap-2">
                        <UFormGroup
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
                                    value-attribute="id"
                                >
                                    <template #label="{ selected }">
                                        <span v-if="selected.length" class="inline-flex flex-wrap gap-1 truncate">
                                            <UBadge
                                                v-for="label in selected"
                                                :key="label.id"
                                                class="truncate"
                                                :class="isColorBright(label.color) ? '!text-black' : '!text-white'"
                                                :style="{ backgroundColor: label.color }"
                                                :label="label.name"
                                            />
                                        </span>
                                        <span v-else>&nbsp;</span>
                                    </template>

                                    <template #option="{ option }">
                                        <UBadge
                                            class="truncate"
                                            :class="isColorBright(option.color) ? '!text-black' : '!text-white'"
                                            :style="{ backgroundColor: option.color }"
                                        >
                                            {{ option.name }}
                                        </UBadge>
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.label', 2)]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup
                            class="flex flex-col"
                            name="namePrefix"
                            :label="$t('common.prefix')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.namePrefix" type="text" />
                        </UFormGroup>

                        <UFormGroup
                            class="flex flex-col"
                            name="nameSuffix"
                            :label="$t('common.suffix')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.nameSuffix" type="text" />
                        </UFormGroup>

                        <UFormGroup
                            class="flex flex-initial flex-col"
                            name="cards"
                            :label="$t('common.card_view')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <div class="flex flex-1 items-center">
                                <UToggle v-model="jobsService.cardView">
                                    <span class="sr-only">
                                        {{ $t('common.card_view') }}
                                    </span>
                                </UToggle>
                            </div>
                        </UFormGroup>
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
            v-model:sort="query.sort"
            class="flex-1"
            :loading="loading"
            :columns="columns"
            :rows="data?.colleagues"
            :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
            sort-mode="manual"
        >
            <template #name-data="{ row: colleague }">
                <div class="inline-flex items-center text-gray-900 dark:text-white">
                    <ProfilePictureImg
                        class="mr-2"
                        :src="colleague?.avatar"
                        :name="`${colleague.firstname} ${colleague.lastname}`"
                        size="sm"
                        :enable-popup="true"
                        :alt="$t('common.avatar')"
                    />

                    <ColleagueName :colleague="colleague" />
                </div>

                <dl class="font-normal lg:hidden">
                    <dt class="sr-only">{{ $t('common.job_grade') }}</dt>
                    <dd class="mt-1 truncate">
                        {{ colleague.jobGradeLabel }}
                        <template v-if="colleague.job !== game.unemployedJobName"> ({{ colleague.jobGrade }})</template>
                    </dd>
                </dl>
            </template>

            <template #rank-data="{ row: colleague }">
                {{ colleague.jobGradeLabel }}
                <template v-if="colleague.job !== game.unemployedJobName"> ({{ colleague.jobGrade }})</template>
            </template>

            <template #absence-data="{ row: colleague }">
                <dl v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))" class="font-normal">
                    <dd class="truncate">
                        {{ $t('common.from') }}:
                        <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                    </dd>
                    <dd class="truncate">
                        {{ $t('common.to') }}: <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                    </dd>
                </dl>
            </template>

            <template #phoneNumber-data="{ row: colleague }">
                <div>
                    <PhoneNumberBlock :number="colleague.phoneNumber" />
                </div>

                <dl class="font-normal lg:hidden">
                    <dt class="sr-only">{{ $t('common.date_of_birth') }}</dt>
                    <dd class="mt-1 truncate">
                        {{ colleague.dateofbirth.value }}
                    </dd>
                </dl>
            </template>

            <template #dateofbirth-data="{ row: colleague }">
                {{ colleague.dateofbirth.value }}
            </template>

            <template #email-data="{ row: colleague }">
                <EmailInfoPopover
                    :email="colleague.email"
                    variant="link"
                    color="primary"
                    truncate
                    :trailing="false"
                    :padded="false"
                    hide-na-text
                />
            </template>

            <template #actions-data="{ row: colleague }">
                <div :key="colleague.id" class="flex flex-col justify-end md:flex-row">
                    <UTooltip
                        v-if="
                            canDo.setJobsUerProps &&
                            (colleague.userId === activeChar!.userId ||
                                attr('jobs.JobsService/SetColleagueProps', 'Types', 'AbsenceDate').value) &&
                            checkIfCanAccessColleague(colleague, 'jobs.JobsService/SetColleagueProps')
                        "
                        :text="$t('components.jobs.self_service.set_absence_date')"
                    >
                        <UButton
                            variant="link"
                            icon="i-mdi-island"
                            @click="
                                modal.open(SelfServicePropsAbsenceDateModal, {
                                    userId: colleague.userId,
                                    'onUpdate:absenceDates': ($event) => updateAbsenceDates($event),
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip
                        v-if="canDo.getColleague && checkIfCanAccessColleague(colleague, 'jobs.JobsService/GetColleague')"
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

                                <span>{{ colleague.dateofbirth.value }}</span>
                            </span>

                            <span class="flex items-center gap-1">
                                <UIcon class="h-5 w-5 shrink-0" name="i-mdi-email" />

                                <EmailInfoPopover :email="colleague.email" variant="link" :trailing="false" :padded="false" />
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
                                        :class="isColorBright(hexToRgb(label.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                        :style="{ backgroundColor: label.color }"
                                        size="xs"
                                        :ui="{ padding: { xs: 'px-2 py-1' } }"
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
                            (canDo.setJobsUerProps &&
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
                                    canDo.setJobsUerProps &&
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
                                        modal.open(SelfServicePropsAbsenceDateModal, {
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

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>

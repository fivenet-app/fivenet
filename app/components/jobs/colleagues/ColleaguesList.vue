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
import { useSettingsStore } from '~/store/settings';
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
    name: z.string().max(50),
    absent: z.boolean(),
    labels: z.custom<Label>().array().max(3),
    namePrefix: z.string().max(12).optional(),
    nameSuffix: z.string().max(12).optional(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    name: '',
    absent: false,
    labels: [],
    namePrefix: undefined,
    nameSuffix: undefined,
});

const settingsStore = useSettingsStore();
const { jobsService } = storeToRefs(settingsStore);

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
    `jobs-colleagues-${sort.value.column}:${sort.value.direction}-${page.value}-${query.name}-${query.absent}-${query.labels.join(',')}-${query.namePrefix}-${query.nameSuffix}`,
    () => listColleagues(),
    {
        transform: (input) => ({ pagination: input.pagination, colleagues: wrapRows(input?.colleagues, columns) }),
        watch: [sort],
    },
);

async function listColleagues(): Promise<ListColleaguesResponse> {
    try {
        const call = $grpc.jobs.jobs.listColleagues({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            search: query.name,
            absent: query.absent,
            labelIds: query.labels.map((l) => l.id),
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
    const idx = query.labels.findIndex((l) => l.id === label.id);
    if (idx > -1) {
        query.labels.splice(idx, 1);
    } else {
        query.labels.push(label);
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
    can(['JobsService.GetColleague', 'JobsService.SetJobsUserProps']).value
        ? {
              key: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined) as TableColumn[];

const { game } = useAppConfig();

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
            <div class="flex gap-2">
                <UFormGroup name="name" :label="$t('common.search')" class="flex-1">
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
                    name="absent"
                    :label="$t('common.absent')"
                    class="flex flex-initial flex-col"
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
                        v-model="sort"
                        :fields="[
                            { label: $t('common.rank'), value: 'rank' },
                            { label: $t('common.name'), value: 'name' },
                        ]"
                    />
                </UFormGroup>

                <UFormGroup
                    v-if="
                        can('JobsService.ManageColleagueLabels').value ||
                        attr('JobsService.GetColleague', 'Types', 'Labels').value
                    "
                    label="&nbsp"
                    :ui="{ container: 'inline-flex gap-1' }"
                >
                    <UButton
                        v-if="can('JobsService.ManageColleagueLabels').value"
                        :label="$t('common.label', 2)"
                        icon="i-mdi-tag"
                        @click="modal.open(JobsLabelsModal, {})"
                    />

                    <UTooltip v-if="attr('JobsService.GetColleague', 'Types', 'Labels').value" :text="$t('common.stats')">
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
                            v-if="attr('JobsService.GetColleague', 'Types', 'Labels').value"
                            name="labels"
                            :label="$t('common.label', 2)"
                            class="flex flex-1 flex-col"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.labels"
                                    class="flex-1"
                                    multiple
                                    :searchable="
                                        async (q: string) => {
                                            return (await getColleagueLabels(q)).labels;
                                        }
                                    "
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['name']"
                                    option-attribute="name"
                                    by="name"
                                    clear-search-on-close
                                >
                                    <template #label>
                                        <span v-if="query.labels.length" class="truncate">
                                            <span v-for="(label, idx) in query.labels" :key="label.id">
                                                <span class="truncate" :style="{ backgroundColor: label.color }">{{
                                                    label.name
                                                }}</span>
                                                <span v-if="idx < query.labels.length - 1">, </span>
                                            </span>
                                        </span>
                                        <span v-else>&nbsp;</span>
                                    </template>

                                    <template #option="{ option }">
                                        <span class="truncate" :style="{ backgroundColor: option.color }">{{
                                            option.name
                                        }}</span>
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
                            name="namePrefix"
                            :label="$t('common.prefix')"
                            class="flex flex-col"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.namePrefix" type="text" />
                        </UFormGroup>

                        <UFormGroup
                            name="nameSuffix"
                            :label="$t('common.suffix')"
                            class="flex flex-col"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <UInput v-model="query.nameSuffix" type="text" />
                        </UFormGroup>

                        <UFormGroup
                            name="cards"
                            :label="$t('common.card_view')"
                            class="flex flex-initial flex-col"
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
            v-model:sort="sort"
            :loading="loading"
            :columns="columns"
            :rows="data?.colleagues"
            :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
            sort-mode="manual"
            class="flex-1"
        >
            <template #name-data="{ row: colleague }">
                <div class="inline-flex items-center text-gray-900 dark:text-white">
                    <ProfilePictureImg
                        :src="colleague?.avatar?.url"
                        :name="`${colleague.firstname} ${colleague.lastname}`"
                        size="sm"
                        :enable-popup="true"
                        :alt="$t('common.avatar')"
                        class="mr-2"
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
                <PhoneNumberBlock :number="colleague.phoneNumber" />

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
                            can('JobsService.SetJobsUserProps').value &&
                            (colleague.userId === activeChar!.userId ||
                                attr('JobsService.SetJobsUserProps', 'Types', 'AbsenceDate').value) &&
                            checkIfCanAccessColleague(colleague, 'JobsService.SetJobsUserProps')
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
                        v-if="
                            can('JobsService.GetColleague').value &&
                            checkIfCanAccessColleague(colleague, 'JobsService.GetColleague')
                        "
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
                                :src="colleague?.avatar?.url"
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
                        <div class="flex flex-col gap-1">
                            <span>
                                {{ colleague.jobGradeLabel }}
                                <template v-if="colleague.job !== game.unemployedJobName">
                                    ({{ colleague.jobGrade }})
                                </template>
                            </span>

                            <PhoneNumberBlock :number="colleague.phoneNumber" />

                            <span class="inline-flex items-center gap-1">
                                <UIcon name="i-mdi-birthday-cake" class="h-5 w-5" />

                                <span>{{ colleague.dateofbirth.value }}</span>
                            </span>

                            <span class="inline-flex items-center gap-1">
                                <UIcon name="i-mdi-email" class="h-5 w-5" />

                                <EmailInfoPopover
                                    :email="colleague.email"
                                    variant="link"
                                    truncate
                                    :trailing="false"
                                    :padded="false"
                                />
                            </span>

                            <div v-if="attr('JobsService.GetColleague', 'Types', 'Labels').value" class="flex flex-row gap-1">
                                <UIcon name="i-mdi-tag" class="h-5 w-5 shrink-0" />

                                <span v-if="!colleague.props?.labels?.list.length">
                                    {{ $t('common.none', [$t('common.label', 2)]) }}
                                </span>
                                <div v-else class="flex max-w-full flex-row flex-wrap gap-1">
                                    <UButton
                                        v-for="label in colleague.props?.labels?.list"
                                        :key="label.name"
                                        :style="{ backgroundColor: label.color }"
                                        class="justify-between gap-2"
                                        :class="
                                            isColourBright(hexToRgb(label.color, RGBBlack)!) ? '!text-black' : '!text-white'
                                        "
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
                                <UIcon name="i-mdi-island" class="size-5" />
                                <GenericTime :value="colleague.props?.absenceBegin" type="shortDate" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                            </span>
                        </div>
                    </template>

                    <template #footer>
                        <UButtonGroup class="inline-flex w-full">
                            <UTooltip
                                v-if="
                                    can('JobsService.SetJobsUserProps').value &&
                                    (colleague.userId === activeChar!.userId ||
                                        attr('JobsService.SetJobsUserProps', 'Types', 'AbsenceDate').value) &&
                                    checkIfCanAccessColleague(colleague, 'JobsService.SetJobsUserProps')
                                "
                                :text="$t('components.jobs.self_service.set_absence_date')"
                                class="flex-1"
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
                                    can('JobsService.GetColleague').value &&
                                    checkIfCanAccessColleague(colleague, 'JobsService.GetColleague')
                                "
                                :text="$t('common.show')"
                                class="flex-1"
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

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>

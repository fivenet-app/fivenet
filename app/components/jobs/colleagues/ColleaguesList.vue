<script lang="ts" setup>
import { isFuture } from 'date-fns';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import EmailBlock from '~/components/partials/citizens/EmailBlock.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { ListColleaguesResponse } from '~~/gen/ts/services/jobs/jobs';
import SelfServicePropsAbsenceDateModal from './SelfServicePropsAbsenceDateModal.vue';

const { t } = useI18n();

const modal = useModal();

const { attr, can, activeChar } = useAuth();

const schema = z.object({
    name: z.string().max(50),
    absent: z.boolean(),
    cards: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    name: '',
    absent: false,
    cards: false,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = ref<TableSortable>({
    column: 'rank',
    direction: 'asc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-colleagues-${sort.value.column}:${sort.value.direction}-${page.value}-${query.name}-${query.absent}`,
    () => listColleagues(),
    {
        transform: (input) => ({ ...input, entries: wrapRows(input?.colleagues, columns) }),
        watch: [sort],
    },
);

async function listColleagues(): Promise<ListColleaguesResponse> {
    try {
        const call = getGRPCJobsClient().listColleagues({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            search: query.name,
            absent: query.absent,
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
].filter((c) => c !== undefined) as { key: string; label: string; sortable?: boolean }[];

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm :schema="schema" :state="query" class="flex w-full gap-2" @submit="refresh()">
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

            <UFormGroup
                name="cards"
                :label="$t('common.card_view')"
                class="flex flex-initial flex-col"
                :ui="{ container: 'flex-1 flex' }"
            >
                <div class="flex flex-1 items-center">
                    <UToggle v-model="query.cards">
                        <span class="sr-only">
                            {{ $t('common.card_view') }}
                        </span>
                    </UToggle>
                </div>
            </UFormGroup>
        </UForm>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.colleague', 2)])" :retry="refresh" />
    <template v-else>
        <UTable
            v-if="!query.cards"
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
                    <span>{{ colleague.firstname }} {{ colleague.lastname }}</span>
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
                <EmailBlock :email="colleague.email" hide-icon hide-na-text />
            </template>

            <template #actions-data="{ row: colleague }">
                <div :key="colleague.id" class="flex flex-col justify-end md:flex-row">
                    <UButton
                        v-if="
                            can('JobsService.SetJobsUserProps').value &&
                            (colleague.userId === activeChar!.userId ||
                                attr('JobsService.SetJobsUserProps', 'Types', 'AbsenceDate').value) &&
                            checkIfCanAccessColleague(colleague, 'JobsService.SetJobsUserProps')
                        "
                        variant="link"
                        icon="i-mdi-island"
                        @click="
                            modal.open(SelfServicePropsAbsenceDateModal, {
                                userId: colleague.userId,
                                'onUpdate:absenceDates': ($event) => updateAbsenceDates($event),
                            })
                        "
                    />

                    <UButton
                        v-if="
                            can('JobsService.GetColleague').value &&
                            checkIfCanAccessColleague(colleague, 'JobsService.GetColleague')
                        "
                        variant="link"
                        icon="i-mdi-eye"
                        :to="{
                            name: 'jobs-colleagues-id-info',
                            params: { id: colleague.userId ?? 0 },
                        }"
                    />
                </div>
            </template>
        </UTable>

        <div v-else class="relative flex-1 overflow-x-auto">
            <UPageGrid
                :ui="{
                    wrapper: 'grid grid-cols-1 p-4 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4',
                }"
            >
                <UPageCard
                    v-for="colleague in data?.colleagues"
                    :key="colleague.userId"
                    :title="`${colleague.firstname} ${colleague.lastname}`"
                >
                    <template #header>
                        <div class="flex items-center justify-center">
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
                                {{ colleague.jobGradeLabel
                                }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                            </span>

                            <PhoneNumberBlock :number="colleague.phoneNumber" />

                            <span class="inline-flex items-center gap-1">
                                <UIcon name="i-mdi-birthday-cake" class="h-5 w-5" />

                                <span>{{ colleague.dateofbirth.value }}</span>
                            </span>

                            <EmailBlock :email="colleague.email" />

                            <span
                                v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))"
                                class="inline-flex items-center gap-1"
                            >
                                <UIcon name="i-mdi-island" class="size-5" />
                                <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                            </span>
                            <span v-else class="h-7"></span>
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

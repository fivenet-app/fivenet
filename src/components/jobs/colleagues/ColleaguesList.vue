<script lang="ts" setup>
import { watchDebounced } from '@vueuse/core';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import type { Perms } from '~~/gen/ts/perms';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useAuthStore } from '~/store/auth';
import SelfServicePropsAbsenceDateModal from './SelfServicePropsAbsenceDateModal.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const query = ref<{
    name: string;
    absent?: boolean;
}>({
    name: '',
    absent: false,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`jobs-colleagues-${page.value}-${query.value.name}`, async () => {
    try {
        const call = $grpc.getJobsClient().listColleagues({
            pagination: {
                offset: offset.value,
            },
            searchName: query.value.name,
            absent: query.value.absent,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
});

watch(offset, async () => refresh());
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });

function updateAbsenceDates(value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void {
    const colleague = data.value?.colleagues.find((c) => c.userId === value.userId);
    if (colleague === undefined) {
        return;
    }

    if (colleague.props === undefined) {
        colleague.props = {
            userId: colleague.userId,
            absenceBegin: value.absenceBegin,
            absenceEnd: value.absenceEnd,
        };
    } else {
        colleague.props.absenceBegin = value.absenceBegin;
        colleague.props.absenceEnd = value.absenceEnd;
    }
}

const today = new Date();
today.setHours(0);
today.setMinutes(0);
today.setSeconds(0);
today.setMilliseconds(0);

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'rank',
        label: t('common.rank'),
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
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
    },
    can(['JobsService.GetColleague', 'JobsService.SetJobsUserProps'] as Perms[])
        ? {
              key: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined);

const modal = useModal();
</script>

<template>
    <div class="py-2 pb-4">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="flex-1">
                                <label for="searchName" class="block text-sm font-medium leading-6">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative mt-2">
                                    <UInput
                                        ref="searchInput"
                                        v-model="query.name"
                                        type="text"
                                        name="searchName"
                                        :placeholder="$t('common.name')"
                                        block
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="flex-initial">
                                <label for="absent" class="block text-sm font-medium leading-6">
                                    {{ $t('common.absent') }}
                                </label>
                                <div class="relative mt-3 flex items-center">
                                    <UToggle v-model="query.absent">
                                        <span class="sr-only">
                                            {{ $t('common.absent') }}
                                        </span>
                                    </UToggle>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataErrorBlock
                            v-if="error"
                            :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
                            :retry="refresh"
                        />
                        <UTable
                            v-else
                            :loading="loading"
                            :columns="columns"
                            :rows="data?.colleagues"
                            :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.vehicle', 2)]) }"
                            :page-count="(data?.pagination?.totalCount ?? 0) / (data?.pagination?.pageSize ?? 1)"
                            :total="data?.pagination?.totalCount"
                        >
                            <template #name-data="{ row: colleague }">
                                {{ colleague.firstname }} {{ colleague.lastname }}
                                <dl class="font-normal lg:hidden">
                                    <dt class="sr-only">{{ $t('common.job_grade') }}</dt>
                                    <dd class="mt-1 truncate">
                                        {{ colleague.jobGradeLabel
                                        }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                                    </dd>
                                </dl>
                            </template>
                            <template #rank-data="{ row: colleague }">
                                {{ colleague.jobGradeLabel
                                }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                            </template>
                            <template #absence-data="{ row: colleague }">
                                <dl
                                    v-if="
                                        colleague.props?.absenceEnd &&
                                        toDate(colleague.props?.absenceEnd).getTime() >= today.getTime()
                                    "
                                    class="font-normal"
                                >
                                    <dd class="truncate text-accent-200">
                                        {{ $t('common.from') }}:
                                        <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                                    </dd>
                                    <dd class="truncate text-accent-200">
                                        {{ $t('common.to') }}: <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                                    </dd>
                                </dl>
                            </template>
                            <template #phone-data="{ row: colleague }">
                                <PhoneNumberBlock :number="colleague.phoneNumber" />
                            </template>
                            <template #actions-data="{ row: colleague }">
                                <UButton
                                    v-if="
                                        can('JobsService.SetJobsUserProps') &&
                                        checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.SetJobsUserProps')
                                    "
                                    variant="link"
                                    icon="i-mdi-island"
                                    @click="
                                        modal.open(SelfServicePropsAbsenceDateModal, {
                                            userId: colleague.userId,
                                        })
                                    "
                                />

                                <UButton
                                    v-if="
                                        can('JobsService.GetColleague') &&
                                        checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.GetColleague')
                                    "
                                    variant="link"
                                    icon="i-mdi-eye"
                                    :to="{
                                        name: 'jobs-colleagues-id',
                                        params: { id: colleague.userId ?? 0 },
                                    }"
                                />
                            </template>
                        </UTable>

                        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
                            <UPagination
                                v-model="page"
                                :page-count="data?.pagination?.pageSize ?? 0"
                                :total="data?.pagination?.totalCount ?? 0"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

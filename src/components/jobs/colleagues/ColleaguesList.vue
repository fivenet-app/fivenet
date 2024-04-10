<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import type { Perms } from '~~/gen/ts/perms';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useAuthStore } from '~/store/auth';
import SelfServicePropsAbsenceDateModal from './SelfServicePropsAbsenceDateModal.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import Pagination from '~/components/partials/Pagination.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const modal = useModal();

const schema = z.object({
    name: z.string().max(50),
    absent: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    name: '',
    absent: false,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`jobs-colleagues-${page.value}-${query.name}`, async () => {
    try {
        const call = $grpc.getJobsClient().listColleagues({
            pagination: {
                offset: offset.value,
            },
            search: query.name,
            absent: query.absent,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
});

watch(offset, async () => refresh());
watchDebounced(query, () => refresh(), { debounce: 200, maxWait: 1250 });

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
].filter((c) => c !== undefined) as { key: string; label: string; sortable?: boolean }[];

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <div class="py-2 pb-4">
        <div class="px-1 sm:px-2">
            <div class="sm:flex sm:items-center">
                <UForm :schema="schema" :state="query" class="flex w-full gap-2" @submit="refresh()">
                    <UFormGroup name="name" :label="$t('common.search')" class="flex-1">
                        <UInput
                            ref="input"
                            v-model="query.name"
                            type="text"
                            name="name"
                            :placeholder="$t('common.name')"
                            block
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        >
                            <template #trailing>
                                <UKbd value="/" />
                            </template>
                        </UInput>
                    </UFormGroup>

                    <UFormGroup name="absent" :label="$t('common.absent')">
                        <UToggle v-model="query.absent">
                            <span class="sr-only">
                                {{ $t('common.absent') }}
                            </span>
                        </UToggle>
                    </UFormGroup>
                </UForm>
            </div>
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
                    :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
                >
                    <template #name-data="{ row: colleague }">
                        <div class="inline-flex items-center text-gray-900 dark:text-white">
                            <ProfilePictureImg
                                :url="colleague.props?.mugShot?.url"
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
                                {{ colleague.jobGradeLabel
                                }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                            </dd>
                        </dl>
                    </template>
                    <template #rank-data="{ row: colleague }">
                        {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                    </template>
                    <template #absence-data="{ row: colleague }">
                        <dl
                            v-if="
                                colleague.props?.absenceEnd && toDate(colleague.props?.absenceEnd).getTime() >= today.getTime()
                            "
                            class="font-normal"
                        >
                            <dd class="truncate">
                                {{ $t('common.from') }}:
                                <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                            </dd>
                            <dd class="truncate">
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
                                    'onUpdate:absenceDates': ($event) => updateAbsenceDates($event),
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
                                name: 'jobs-colleagues-id-actvitiy',
                                params: { id: colleague.userId ?? 0 },
                            }"
                        />
                    </template>
                </UTable>

                <Pagination v-model="page" :pagination="data?.pagination" />
            </div>
        </div>
    </div>
</template>

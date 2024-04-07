<script lang="ts" setup>
import { min, numeric, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type { ListInactiveEmployeesResponse } from '~~/gen/ts/services/jobs/timeclock';
import { checkIfCanAccessColleague } from '../colleagues/helpers';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const query = ref<{
    days: number;
}>({
    days: 14,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`jobs-timeclock-inactive-${page.value}-${query.value.days}`, () => listInactiveEmployees());

async function listInactiveEmployees(): Promise<ListInactiveEmployeesResponse> {
    try {
        const call = $grpc.getJobsTimeclockClient().listInactiveEmployees({
            pagination: {
                offset: offset.value,
            },
            days: query.value.days,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    days: number;
}

defineRule('required', required);
defineRule('min', min);
defineRule('numeric', numeric);

const { meta } = useForm<FormData>({
    validationSchema: {
        days: { required: true, min: 1, numeric: true },
    },
});

watchDebounced(
    query.value,
    async () => {
        if (meta.value.valid) {
            refresh();
        }
    },
    { debounce: 600, maxWait: 1400 },
);
watch(offset, async () => refresh());

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'rank',
        label: t('common.rank', 1),
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
].filter((c) => c.permission === undefined || can(c.permission));
</script>

<template>
    <div>
        <UDashboardToolbar>
            <template #default>
                <div class="flex w-full flex-col">
                    <UButton
                        v-if="can('JobsTimeclockService.ListTimeclock')"
                        :to="{ name: 'jobs-timeclock' }"
                        icon="i-mdi-arrow-left"
                        class="place-self-start"
                    >
                        {{ $t('common.timeclock') }}
                    </UButton>

                    <UForm :state="{}" class="flex w-full flex-row gap-2" @submit="refresh()">
                        <UFormGroup name="days" :label="$t('common.time_ago.day', 2)" class="flex-1">
                            <UInput
                                v-model="query.days"
                                name="days"
                                type="number"
                                min="3"
                                max="31"
                                :placeholder="$t('common.time_ago.day', 2)"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>
                    </UForm>
                </div>
            </template>
        </UDashboardToolbar>

        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.colleague', 2)])" :retry="refresh" />
        <UTable
            v-else
            :loading="loading"
            :columns="columns"
            :rows="data?.colleagues"
            :empty-state="{ icon: 'i-mdi-account', label: $t('common.not_found', [$t('common.colleague', 2)]) }"
        >
            <template #name-data="{ row: colleague }">
                <div class="inline-flex items-center">
                    <ProfilePictureImg
                        :url="colleague.avatar?.url"
                        :name="`${colleague.firstname} ${colleague.lastname}`"
                        size="sm"
                        :enable-popup="true"
                        class="mr-2"
                    />

                    {{ colleague.firstname }} {{ colleague.lastname }}
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
                <NuxtLink
                    v-if="checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.GetColleague')"
                    icon="i-mdi-eye"
                    :to="{
                        name: 'jobs-colleagues-id-actvitiy',
                        params: { id: colleague.userId ?? 0 },
                    }"
                />
            </template>
        </UTable>

        <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="data?.pagination?.pageSize ?? 0"
                :total="data?.pagination?.totalCount ?? 0"
            />
        </div>
    </div>
</template>

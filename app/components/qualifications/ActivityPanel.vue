<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { type QualificationActivity, QualificationActivityType } from '~~/gen/ts/resources/qualifications/activity/activity';
import { RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListQualificationActivityResponse } from '~~/gen/ts/services/qualifications/qualifications';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import {
    qualificationActivityTypeColor,
    qualificationActivityTypeIcon,
    requestStatusToBadgeColor,
    resultStatusToBadgeColor,
} from './helpers';

const props = defineProps<{ qualificationId: number }>();
const { t } = useI18n();
const completorStore = useCompletorStore();
const qualificationsClient = await getQualificationsQualificationsClient();
const page = ref(1);
const types = ref<QualificationActivityType[]>([]);
const user = ref<UserShort>();
const dateRange = ref<DateRange>();
const sorting = ref<{ columns: SortByColumn[] }>({ columns: [{ id: 'createdAt', desc: true }] });

const activityKey = computed(
    () =>
        `qualification-activity-${props.qualificationId}-${page.value}-${types.value.join(',')}-${user.value?.userId ?? 0}-${dateRange.value?.start.toISOString() ?? ''}-${dateRange.value?.end.toISOString() ?? ''}-${JSON.stringify(sorting.value)}`,
);
const { data, status, error, refresh } = useAuthedLazyAsyncData(
    'userState',
    activityKey,
    ({ signal }) => listActivity(signal),
    {
        watch: [() => props.qualificationId, page, types, () => user.value?.userId, dateRange, sorting],
    },
);
const activityItems = computed<QualificationActivity[]>(() => data.value?.activity ?? []);

function activityTypeLabel(activityType: QualificationActivityType): string {
    return t(`enums.qualifications.QualificationActivityType.${QualificationActivityType[activityType]}`);
}

const activityTypeItems = computed(() =>
    Object.values(QualificationActivityType)
        .filter(
            (value): value is QualificationActivityType =>
                typeof value === 'number' && value !== QualificationActivityType.UNSPECIFIED,
        )
        .map((value) => ({
            label: activityTypeLabel(value),
            value,
            icon: qualificationActivityTypeIcon(value),
            ui: {
                itemLeadingIcon: qualificationActivityTypeColor(value),
            },
        })),
);

async function listActivity(signal: AbortSignal): Promise<ListQualificationActivityResponse> {
    const { response } = await qualificationsClient.listQualificationActivity(
        {
            qualificationId: props.qualificationId,
            pagination: { offset: calculateOffset(page.value, data.value?.pagination) },
            types: types.value,
            userId: user.value?.userId,
            sort: sorting.value,
            from: toTimestamp(dateRange.value?.start),
            to: toTimestamp(dateRange.value?.end),
        },
        { abort: signal },
    );
    return response;
}

watch(
    [types, () => user.value?.userId, dateRange, sorting],
    async () => {
        if (page.value === 1) await refresh();
        else page.value = 1;
    },
    { deep: true },
);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <div class="my-2 grid w-full gap-3 md:grid-cols-[repeat(3,minmax(0,1fr))_auto]">
                    <UFormField class="w-full" :label="$t('common.citizen', 1)">
                        <SelectMenu
                            v-model="user"
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeCitizens({
                                        search: q,
                                        userIds: user?.userId ? [user.userId] : [],
                                    })
                            "
                            searchable-key="qualification-activity-user"
                            :filter-fields="['firstname', 'lastname']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            clear
                        >
                            <template v-if="user" #default>
                                {{ userToLabel(user) }}
                            </template>

                            <template #item-label="{ item }">
                                {{ userToLabel(item) }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField class="w-full" :label="$t('common.type')">
                        <USelectMenu
                            v-model="types"
                            class="w-full"
                            multiple
                            :items="activityTypeItems"
                            value-key="value"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                        </USelectMenu>
                    </UFormField>

                    <UFormField class="w-full" :label="$t('common.date')">
                        <InputDateRangePopover v-model="dateRange" class="w-full" />
                    </UFormField>

                    <UFormField label="&nbsp;">
                        <SortButton v-model="sorting" :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]" />
                    </UFormField>
                </div>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="relative flex-1 overflow-x-auto">
                <DataErrorBlock
                    v-if="error"
                    class="w-full"
                    :title="$t('common.not_found', [`${$t('common.qualification', 1)} ${$t('common.activity')}`])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!isRequestPending(status) && activityItems.length === 0"
                    class="w-full"
                    icon="i-mdi-pulse"
                    :type="`${$t('common.qualification', 1)} ${$t('common.activity')}`"
                />

                <ul v-else class="divide-y divide-default">
                    <li v-for="entry in activityItems" :key="entry.id" class="flex gap-3 p-3">
                        <UIcon
                            :name="qualificationActivityTypeIcon(entry.type)"
                            :class="[qualificationActivityTypeColor(entry.type), 'mt-1 size-7 shrink-0']"
                        />

                        <div class="min-w-0 flex-1 space-y-1">
                            <div class="flex items-center justify-between gap-2">
                                <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                                    <span class="text-sm font-medium">{{ activityTypeLabel(entry.type) }}</span>
                                    <UBadge
                                        v-if="entry.data?.requestStatus !== undefined"
                                        :color="requestStatusToBadgeColor(entry.data.requestStatus)"
                                        :label="
                                            $t(`enums.qualifications.RequestStatus.${RequestStatus[entry.data.requestStatus]}`)
                                        "
                                    />
                                    <UBadge
                                        v-if="entry.data?.resultStatus !== undefined"
                                        :color="resultStatusToBadgeColor(entry.data.resultStatus)"
                                        :label="
                                            $t(`enums.qualifications.ResultStatus.${ResultStatus[entry.data.resultStatus]}`)
                                        "
                                    />
                                    <UBadge
                                        v-if="entry.data?.score !== undefined"
                                        color="neutral"
                                        :label="`${entry.data.score}`"
                                    />
                                </div>
                                <GenericTime class="shrink-0 text-sm text-dimmed" :value="entry.createdAt" type="long" />
                            </div>

                            <div class="flex items-center justify-between gap-2 text-sm">
                                <div v-if="entry.targetUser" class="inline-flex items-center gap-1">
                                    <span class="font-semibold">{{ $t('common.citizen', 1) }}:</span>
                                    <CitizenInfoPopover :user="entry.targetUser" show-avatar />
                                </div>
                                <div v-else />
                                <div class="inline-flex shrink-0 items-center gap-1">
                                    <span>{{ $t('common.created_by') }}</span>
                                    <CitizenInfoPopover v-if="entry.actorUser" :user="entry.actorUser" show-avatar />
                                    <span v-else>{{ $t('common.system') }}</span>
                                </div>
                            </div>
                        </div>
                    </li>
                </ul>
            </div>
        </template>

        <template #footer>
            <Pagination
                v-if="data?.pagination"
                v-model="page"
                :pagination="data.pagination"
                :status="status"
                :refresh="refresh"
            />
        </template>
    </UDashboardPanel>
</template>

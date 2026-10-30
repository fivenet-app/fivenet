<script lang="ts" setup>
import { z } from 'zod';
import type { Form } from '@nuxt/ui';
import ActivityFeedFilterLayout from '~/components/partials/data/ActivityFeedFilterLayout.vue';
import ActivityFeedSkeleton from '~/components/partials/data/ActivityFeedSkeleton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { UserActivityType } from '~~/gen/ts/resources/users/activity/activity';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizens/citizens';
import ActivityFeedEntry from './ActivityFeedEntry.vue';
import { citizenUserActivityIndicatorColor, citizenUserActivityTypeBGColor, citizenUserActivityTypeIcon } from './helpers';

const props = defineProps<{ userId: number }>();
const { t } = useI18n();
const { attr, activeChar } = useAuth();
const citizensCitizensClient = await getCitizensCitizensClient();
const activityTypes = Object.keys(UserActivityType)
    .map((t) => UserActivityType[t as keyof typeof UserActivityType])
    .filter((at) => typeof at === 'number' && at >= 3);
const propertyActivityKeys: Partial<Record<UserActivityType, string>> = {
    [UserActivityType.LICENSES]: 'user.licenses',
    [UserActivityType.WANTED]: 'user.wanted',
    [UserActivityType.TRAFFIC_INFRACTION_POINTS]: 'user.traffic_infraction_points',
    [UserActivityType.MUGSHOT]: 'user.mugshot',
    [UserActivityType.FILE]: 'user.files',
    [UserActivityType.LABELS]: 'user.labels',
    [UserActivityType.JOB]: 'user.job',
    [UserActivityType.FINE]: 'user.open_fines',
};
const visibleActivityTypes = activityTypes.filter((activityType) => {
    const propertyKey = propertyActivityKeys[activityType];
    return propertyKey === undefined || attr('citizens.CitizensService/ListCitizens', 'Properties', propertyKey).value;
});
const options = visibleActivityTypes.map((at) => ({
    label: t(`enums.users.UserActivityType.${UserActivityType[at]}`),
    icon: citizenUserActivityTypeIcon(at),
    value: at,
    ui: { itemLeadingIcon: citizenUserActivityTypeBGColor(at) },
}));
const schema = z.object({
    types: z.enum(UserActivityType).array().max(activityTypes.length).default([]),
    dateRange: z.custom<DateRange>().optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([{ id: 'createdAt', desc: true }]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});
type Schema = z.output<typeof schema>;
const query = useSearchForm('citizen_activity', schema);
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);
const activityKey = computed(() => `citizeninfo-activity-${props.userId}-${JSON.stringify(validatedQuery.value)}`);
const { data, status, refresh, error } = useAuthedLazyAsyncData('userState', activityKey, ({ signal }) =>
    listUserActivity(validatedQuery.value, signal),
);
const timelineItems = computed(() =>
    (data.value?.activity ?? []).map((activity) => ({
        ...activity,
        icon: citizenUserActivityTypeIcon(activity.type),
        ui: { indicator: citizenUserActivityIndicatorColor(activity) },
    })),
);
async function listUserActivity(values: Schema, signal: AbortSignal): Promise<ListUserActivityResponse> {
    try {
        const { response } = await citizensCitizensClient.listUserActivity(
            {
                pagination: { offset: calculateOffset(values.page, data.value?.pagination) },
                sort: values.sorting,
                userId: props.userId,
                types: values.types,
                from: toTimestamp(values.dateRange?.start),
                to: toTimestamp(values.dateRange?.end),
                propertyKeys: [],
            },
            { abort: signal },
        );
        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
const denyView = computed(
    () =>
        props.userId === activeChar.value?.userId && !attr('citizens.CitizensService/ListUserActivity', 'Fields', 'Own').value,
);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0' }">
        <template v-if="!denyView" #header>
            <UDashboardToolbar>
                <template #default>
                    <UForm
                        ref="formRef"
                        class="my-2 flex w-full flex-col gap-2 sm:flex-row sm:flex-wrap"
                        :schema="schema"
                        :state="query"
                        @submit="commitValidatedQuery"
                    >
                        <ActivityFeedFilterLayout :field-count="3">
                            <UFormField class="w-full min-w-0" name="types" :label="$t('common.type', 2)"
                                ><ClientOnly
                                    ><USelectMenu
                                        v-model="query.types"
                                        class="w-full min-w-40 flex-1"
                                        multiple
                                        nullable
                                        :items="options"
                                        value-key="value"
                                        :search-input="{ placeholder: $t('common.type', 2) }"
                                        ><template #empty>{{
                                            $t('common.not_found', [$t('common.type', 2)])
                                        }}</template></USelectMenu
                                    ></ClientOnly
                                ></UFormField
                            >
                            <UFormField class="w-full min-w-0" name="dateRange" :label="$t('common.date')"
                                ><InputDateRangePopover v-model="query.dateRange" class="w-full" clearable time
                            /></UFormField>
                            <UFormField :label="$t('common.sort')"
                                ><SortButton
                                    v-model="query.sorting"
                                    :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]"
                            /></UFormField>
                        </ActivityFeedFilterLayout>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>
        <template #body>
            <UContainer v-if="denyView" class="my-2"
                ><UAlert
                    variant="subtle"
                    color="error"
                    icon="i-mdi-denied"
                    :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
                    :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
            /></UContainer>
            <ActivityFeedSkeleton v-else-if="isRequestPending(status)" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data || data.activity.length === 0"
                :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                icon="i-mdi-pulse"
            />
            <UTimeline v-else :items="timelineItems" size="md" :ui="{ wrapper: '!mt-0 !pb-2' }"
                ><template #indicator="{ item }"><UIcon :name="item.icon" class="size-6" /></template
                ><template #wrapper="{ item: activity }"><ActivityFeedEntry :activity="activity" /></template
            ></UTimeline>
        </template>
        <template #footer
            ><Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh"
        /></template>
    </UDashboardPanel>
</template>

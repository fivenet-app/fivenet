<script lang="ts" setup>
import type { Form } from '@nuxt/ui';
import { z } from 'zod';
import ActivityPanelEntry from '~/components/qualifications/ActivityPanelEntry.vue';
import ActivityFeedFilterLayout from '~/components/partials/data/ActivityFeedFilterLayout.vue';
import ActivityFeedSkeleton from '~/components/partials/data/ActivityFeedSkeleton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { type QualificationActivity, QualificationActivityType } from '~~/gen/ts/resources/qualifications/activity/activity';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListQualificationActivityResponse } from '~~/gen/ts/services/qualifications/qualifications';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import {
    qualificationActivityTypeColor,
    qualificationActivityTypeIcon,
    qualificationActivityTypeIndicatorColor,
} from './helpers';

const props = defineProps<{ qualificationId: number }>();
const { t } = useI18n();
const completorStore = useCompletorStore();
const qualificationsClient = await getQualificationsQualificationsClient();

const schema = z.object({
    types: z.enum(QualificationActivityType).array().default([]),
    user: z.custom<UserShort>().optional(),
    dateRange: z.custom<DateRange>().optional(),
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

type Schema = z.output<typeof schema>;

const query = useSearchForm(`qualification_activity_${props.qualificationId}`, schema);
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const activityKey = computed(() => `qualification-activity-${props.qualificationId}-${JSON.stringify(validatedQuery.value)}`);
const { data, status, error, refresh } = useAuthedLazyAsyncData('userState', activityKey, ({ signal }) =>
    listActivity(validatedQuery.value, signal),
);
const activityItems = computed<QualificationActivity[]>(() => data.value?.activity ?? []);
const timelineItems = computed(() =>
    activityItems.value.map((activity) => ({
        ...activity,
        icon: qualificationActivityTypeIcon(activity.type),
        ui: { indicator: qualificationActivityTypeIndicatorColor(activity.type) },
    })),
);

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

async function listActivity(values: Schema, signal: AbortSignal): Promise<ListQualificationActivityResponse> {
    const { response } = await qualificationsClient.listQualificationActivity(
        {
            qualificationId: props.qualificationId,
            pagination: { offset: calculateOffset(values.page, data.value?.pagination) },
            types: values.types,
            userId: values.user?.userId,
            sort: values.sorting,
            from: toTimestamp(values.dateRange?.start),
            to: toTimestamp(values.dateRange?.end),
        },
        { abort: signal },
    );
    return response;
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0' }">
        <template #header>
            <UDashboardToolbar>
                <UForm ref="formRef" class="my-2 w-full" :schema="schema" :state="query" @submit="commitValidatedQuery">
                    <ActivityFeedFilterLayout :field-count="4">
                        <UFormField class="w-full" name="user" :label="$t('common.citizen', 1)">
                            <SelectMenu
                                v-model="query.user"
                                class="w-full"
                                :searchable="
                                    async (q: string) =>
                                        await completorStore.completeCitizens({
                                            search: q,
                                            userIds: query.user?.userId ? [query.user.userId] : [],
                                        })
                                "
                                searchable-key="qualification-activity-user"
                                :filter-fields="['firstname', 'lastname']"
                                :search-input="{ placeholder: $t('common.search_field') }"
                                clear
                            >
                                <template v-if="query.user" #default>
                                    {{ userToLabel(query.user) }}
                                </template>

                                <template #item-label="{ item }">
                                    {{ userToLabel(item) }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                            </SelectMenu>
                        </UFormField>

                        <UFormField class="w-full" name="types" :label="$t('common.type')">
                            <USelectMenu
                                v-model="query.types"
                                class="w-full"
                                multiple
                                nullable
                                :items="activityTypeItems"
                                value-key="value"
                                :search-input="{ placeholder: $t('common.search_field') }"
                            >
                                <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                            </USelectMenu>
                        </UFormField>

                        <UFormField class="w-full" name="dateRange" :label="$t('common.date')">
                            <InputDateRangePopover v-model="query.dateRange" class="w-full" clearable time />
                        </UFormField>

                        <UFormField :label="$t('common.sort')">
                            <SortButton
                                v-model="query.sorting"
                                :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]"
                            />
                        </UFormField>
                    </ActivityFeedFilterLayout>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="relative flex-1">
                <ActivityFeedSkeleton v-if="isRequestPending(status)" />

                <DataErrorBlock
                    v-else-if="error"
                    class="w-full"
                    :title="$t('common.not_found', [`${$t('common.qualification', 1)} ${$t('common.activity')}`])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="activityItems.length === 0"
                    class="w-full"
                    icon="i-mdi-pulse"
                    :type="`${$t('common.qualification', 1)} ${$t('common.activity')}`"
                />

                <UTimeline v-else :items="timelineItems" size="md" :ui="{ wrapper: '!mt-0 !pb-2' }">
                    <template #indicator="{ item }">
                        <UIcon :name="item.icon" class="size-6" />
                    </template>

                    <template #wrapper="{ item }">
                        <ActivityPanelEntry :activity="item" />
                    </template>
                </UTimeline>
            </div>
        </template>

        <template #footer>
            <Pagination
                v-if="data?.pagination"
                v-model="query.page"
                :pagination="data.pagination"
                :status="status"
                :refresh="refresh"
            />
        </template>
    </UDashboardPanel>
</template>

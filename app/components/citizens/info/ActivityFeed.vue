<script lang="ts" setup>
import { z } from 'zod';
import ActivityFeedFilterLayout from '~/components/partials/data/ActivityFeedFilterLayout.vue';
import ActivityFeedItem from '~/components/partials/data/ActivityFeedItem.vue';
import ActivityFeedSkeleton from '~/components/partials/data/ActivityFeedSkeleton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import type { Form } from '@nuxt/ui';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { UserActivityType } from '~~/gen/ts/resources/users/activity/activity';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizens/citizens';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import LabelBadge from '~/components/citizens/labels/LabelBadge.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/relations/relations';
import { citizenUserActivityIconColor, citizenUserActivityTypeBGColor, citizenUserActivityTypeIcon } from './helpers';

const props = defineProps<{
    userId: number;
}>();

const { t } = useI18n();
const numberFormatter = useDisplayNumberFormat();

const { attr, activeChar } = useAuth();

const citizensCitizensClient = await getCitizensCitizensClient();

const activityTypes = Object.keys(UserActivityType)
    .map((t) => UserActivityType[t as keyof typeof UserActivityType])
    .filter((at) => {
        if (typeof at === 'string') {
            return false;
        } else if (typeof at === 'number' && at < 3) {
            return false;
        }
        return true;
    });

const options = activityTypes.map((at) => ({
    label: t(`enums.users.UserActivityType.${UserActivityType[at]}`),
    icon: citizenUserActivityTypeIcon(at),
    value: at,
    ui: {
        itemLeadingIcon: citizenUserActivityTypeBGColor(at),
    },
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

const query = useSearchForm('citizen_activity', schema);

type Schema = z.output<typeof schema>;

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const activityKey = computed(() => `citizeninfo-activity-${props.userId}-${JSON.stringify(validatedQuery.value)}`);

const { data, status, refresh, error } = useAuthedLazyAsyncData('userState', activityKey, ({ signal }) =>
    listUserActivity(validatedQuery.value, signal),
);

async function listUserActivity(values: Schema, signal: AbortSignal): Promise<ListUserActivityResponse> {
    try {
        const call = citizensCitizensClient.listUserActivity(
            {
                pagination: {
                    offset: calculateOffset(values.page, data.value?.pagination),
                },
                sort: values.sorting,
                userId: props.userId,
                types: values.types,
                from: toTimestamp(values.dateRange?.start),
                to: toTimestamp(values.dateRange?.end),
            },
            { abort: signal },
        );
        const { response } = await call;

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
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
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
                            <UFormField class="w-full min-w-0" name="types" :label="$t('common.type', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.types"
                                        class="w-full min-w-40 flex-1"
                                        multiple
                                        nullable
                                        :items="options"
                                        value-key="value"
                                        :search-input="{ placeholder: $t('common.type', 2) }"
                                    >
                                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full min-w-0" name="dateRange" :label="$t('common.date')">
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
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer v-if="denyView" class="my-2">
                <UAlert
                    variant="subtle"
                    color="error"
                    icon="i-mdi-denied"
                    :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
                    :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
                />
            </UContainer>

            <ActivityFeedSkeleton v-else-if="isRequestPending(status)" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data || data?.activity.length === 0"
                :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                icon="i-mdi-pulse"
            />

            <div v-else class="relative flex-1">
                <ul class="min-w-full divide-y divide-default" role="list">
                    <ActivityFeedItem
                        v-for="activity in data?.activity"
                        :key="activity.id"
                        :icon="citizenUserActivityTypeIcon(activity.type)"
                        :icon-class="citizenUserActivityIconColor(activity)"
                    >
                        <template #title>
                            <template
                                v-if="activity.type === UserActivityType.NAME && activity.data?.data.oneofKind === 'nameChange'"
                            >
                                <I18nT keypath="components.citizens.CitizenInfoActivityFeedEntry.name_change">
                                    <template #old>
                                        <span class="font-semibold">{{ activity.data.data.nameChange.old }}</span>
                                    </template>
                                    <template #new>
                                        <span class="font-semibold">{{ activity.data.data.nameChange.new }}</span>
                                    </template>
                                </I18nT>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.DOCUMENT &&
                                    activity.data?.data.oneofKind === 'documentRelation'
                                "
                            >
                                {{
                                    $t(
                                        `components.citizens.CitizenInfoActivityFeedEntry.document_relation.${activity.data.data.documentRelation.added ? 'added' : 'removed'}`,
                                    )
                                }}
                                <DocumentInfoPopover :document-id="activity.data.data.documentRelation.documentId" load-on-open>
                                    <template #title>
                                        <IDCopyBadge
                                            :id="activity.data.data.documentRelation.documentId"
                                            prefix="DOC"
                                            size="xs"
                                            disable-tooltip
                                            variant="link"
                                            hide-icon
                                        />
                                    </template>
                                </DocumentInfoPopover>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.WANTED &&
                                    activity.data?.data.oneofKind === 'wantedChange'
                                "
                            >
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.bool_set_citizen') }}
                                <span class="font-semibold">
                                    {{
                                        activity.data.data.wantedChange.wanted
                                            ? $t('common.wanted')
                                            : `${$t('common.not')} ${$t('common.wanted')}`
                                    }}
                                </span>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.JOB && activity.data?.data.oneofKind === 'jobChange'
                                "
                            >
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_job_set') }}
                                <span class="font-semibold">
                                    {{ activity.data.data.jobChange.jobLabel }}
                                    <span v-if="activity.data.data.jobChange.grade">
                                        ({{ $t('common.rank') }}: {{ activity.data.data.jobChange.gradeLabel }})
                                    </span>
                                </span>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.TRAFFIC_INFRACTION_POINTS &&
                                    activity.data?.data.oneofKind === 'trafficInfractionPointsChange'
                                "
                            >
                                <I18nT
                                    keypath="components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text"
                                >
                                    <template #old>
                                        <span class="font-semibold">{{
                                            activity.data.data.trafficInfractionPointsChange.old
                                        }}</span>
                                    </template>
                                    <template #new>
                                        <span class="font-semibold">{{
                                            activity.data.data.trafficInfractionPointsChange.new
                                        }}</span>
                                    </template>
                                </I18nT>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.MUGSHOT &&
                                    activity.data?.data.oneofKind === 'mugshotChange'
                                "
                            >
                                {{
                                    $t(
                                        `components.citizens.CitizenInfoActivityFeedEntry.user_props_mugshot_${activity.data.data.mugshotChange.new ? 'set' : 'removed'}`,
                                    )
                                }}
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.LABELS &&
                                    activity.data?.data.oneofKind === 'labelsChange'
                                "
                            >
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_labels_updated') }}
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.LICENSES &&
                                    activity.data?.data.oneofKind === 'licensesChange'
                                "
                            >
                                {{
                                    activity.data.data.licensesChange.added
                                        ? $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.added')
                                        : $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.removed')
                                }}:
                                {{ activity.data.data.licensesChange.licenses.map((license) => license.label).join(', ') }}
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.JAIL && activity.data?.data.oneofKind === 'jailChange'
                                "
                            >
                                <template v-if="activity.data.data.jailChange.seconds > 0">
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.jailed') }}
                                    {{ fromSecondsToFormattedDuration(activity.data.data.jailChange.seconds) }}
                                </template>
                                <template v-else-if="activity.data.data.jailChange.seconds === 0">
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.unjailed') }}
                                </template>
                                <template v-else>
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.escaped') }}
                                </template>
                            </template>

                            <template
                                v-else-if="
                                    activity.type === UserActivityType.FINE && activity.data?.data.oneofKind === 'fineChange'
                                "
                            >
                                <template v-if="activity.data.data.fineChange.removed">
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.removed') }}
                                </template>
                                <template v-else-if="activity.data.data.fineChange.amount < 0">
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.paid') }}
                                </template>
                                <template v-else>
                                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.created') }}
                                </template>
                                {{ numberFormatter.format(Math.abs(activity.data.data.fineChange.amount)) }}
                            </template>

                            <template v-else>{{ UserActivityType[activity.type] }}</template>
                        </template>

                        <template #timestamp>
                            <GenericTime :value="activity.createdAt" type="long" />
                        </template>

                        <div
                            v-if="activity.type === UserActivityType.LABELS && activity.data?.data.oneofKind === 'labelsChange'"
                            class="flex flex-wrap gap-1"
                        >
                            <LabelBadge
                                v-for="label in activity.data.data.labelsChange.removed"
                                :key="label.id"
                                class="line-through"
                                :label="label"
                            />
                            <LabelBadge v-for="label in activity.data.data.labelsChange.added" :key="label.id" :label="label" />
                            <LabelBadge
                                v-for="label in activity.data.data.labelsChange.removedIds"
                                :id="label"
                                :key="label"
                                :expired="activity.data.data.labelsChange.expired"
                            />
                            <LabelBadge
                                v-for="label in activity.data.data.labelsChange.addedIds"
                                :id="label.id"
                                :key="label.id"
                                :expires-at="label.expiresAt"
                            />
                        </div>

                        <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                            <p class="inline-flex min-w-0 gap-1 text-sm break-words">
                                <template
                                    v-if="
                                        activity.type === UserActivityType.DOCUMENT &&
                                        activity.data?.data.oneofKind === 'documentRelation'
                                    "
                                >
                                    <span class="font-semibold">{{ $t('common.type') }}:</span>
                                    <span>{{
                                        $t(
                                            `enums.documents.DocRelation.${DocRelation[activity.data.data.documentRelation.relation]}`,
                                        )
                                    }}</span>
                                </template>
                                <template
                                    v-else-if="
                                        activity.type !== UserActivityType.JAIL ||
                                        (activity.data?.data.oneofKind === 'jailChange' &&
                                            activity.data.data.jailChange.seconds >= 0)
                                    "
                                >
                                    <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                                    <span
                                        v-if="
                                            activity.type === UserActivityType.LABELS &&
                                            activity.data?.data.oneofKind === 'labelsChange' &&
                                            activity.data.data.labelsChange.expired
                                        "
                                    >
                                        {{ $t('common.expired') }}
                                    </span>
                                    <!-- eslint-disable-next-line vue/no-v-html -->
                                    <span v-else v-html="activity.reason" />
                                </template>
                            </p>

                            <p v-if="activity.sourceUser" class="inline-flex shrink-0 text-sm">
                                {{ $t('common.created_by') }}
                                <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                            </p>
                        </div>
                    </ActivityFeedItem>
                </ul>
            </div>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>

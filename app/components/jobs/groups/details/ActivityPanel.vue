<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import QualificationBadge from '~/components/partials/qualifications/QualificationBadge.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import { type GroupActivity, GroupActivityType } from '~~/gen/ts/resources/jobs/groups/activity';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListGroupActivityResponse } from '~~/gen/ts/services/jobs/groups';
import { groupActivityTypeColor, groupActivityTypeIcon, groupRuleLabel } from '../helpers';

const props = defineProps<{
    groupId: number;
    canView: boolean;
}>();

const { t } = useI18n();
const completorStore = useCompletorStore();
const jobsGroupsClient = await getJobsGroupsClient();

const page = ref(1);
const selectedTypes = ref<GroupActivityType[]>([]);
const selectedUser = ref<UserShort>();
const dateRange = ref<DateRange>();

const activityKey = computed(
    () =>
        `jobs-group-activity-${props.groupId}-${page.value}-${selectedTypes.value.join(',')}-${selectedUser.value?.userId ?? 0}-${
            dateRange.value?.start.toISOString() ?? ''
        }-${dateRange.value?.end.toISOString() ?? ''}`,
);

const {
    data: activity,
    status: activityStatus,
    error: activityError,
    refresh: refreshActivity,
} = useLazyAsyncData(activityKey, () => listGroupActivity(), {
    watch: [() => props.groupId, page, selectedTypes, () => selectedUser.value?.userId, dateRange],
});

const activityItems = computed<GroupActivity[]>(() => activity.value?.activity ?? []);

const activityTypeItems = computed(() =>
    Object.values(GroupActivityType)
        .filter((value): value is GroupActivityType => typeof value === 'number' && value !== GroupActivityType.UNSPECIFIED)
        .map((value) => ({
            label: t(`enums.jobs.groups.GroupActivityType.${GroupActivityType[value] ?? 'UNSPECIFIED'}`),
            value,
            icon: groupActivityTypeIcon(value),
            ui: {
                itemLeadingIcon: groupActivityTypeColor(value),
            },
        })),
);

async function listGroupActivity(): Promise<ListGroupActivityResponse> {
    const { response } = await jobsGroupsClient.listGroupActivity({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, activity.value?.pagination),
        },
        sort: { columns: [{ id: 'created_at', desc: true }] },
        types: selectedTypes.value,
        userId: selectedUser.value?.userId,
        from: toTimestamp(dateRange.value?.start),
        to: toTimestamp(dateRange.value?.end),
    });

    return response;
}

async function applyFilters(): Promise<void> {
    if (page.value === 1) {
        await refreshActivity();
        return;
    }

    page.value = 1;
}

async function clearFilters(): Promise<void> {
    selectedTypes.value = [];
    selectedUser.value = undefined;
    dateRange.value = undefined;
    if (page.value === 1) {
        await refreshActivity();
        return;
    }

    page.value = 1;
}

function activityLabel(activity: GroupActivity): string {
    return t(`enums.jobs.groups.GroupActivityType.${GroupActivityType[activity.type] ?? 'UNSPECIFIED'}`);
}

function activityRuleLabel(activity: GroupActivity): string | undefined {
    if (activity.data?.data.oneofKind === 'rule') {
        const rule = activity.data.data.rule;
        const ruleId = rule.id || activity.ruleId;
        return ruleId ? `#${ruleId} - ${groupRuleLabel(rule, t)}` : groupRuleLabel(rule, t);
    }

    if (activity.ruleId) {
        return `#${activity.ruleId}`;
    }

    return undefined;
}

function activityRuleQualification(activity: GroupActivity, qualificationId: number): QualificationShort | undefined {
    if (activity.data?.data.oneofKind !== 'rule' || activity.data.data.rule.rule.oneofKind !== 'qualification') {
        return undefined;
    }

    return activity.data.data.rule.rule.qualification.qualifications.find(
        (qualification) => qualification.id === qualificationId,
    );
}

watch(
    () => props.groupId,
    () => {
        page.value = 1;
    },
);
</script>

<template>
    <div v-if="canView" class="grid gap-4">
        <UCard variant="subtle">
            <div class="grid gap-3 lg:grid-cols-[minmax(180px,260px)_minmax(0,1fr)_minmax(260px,1fr)_auto] lg:items-end">
                <UFormField :label="$t('common.type')">
                    <USelectMenu
                        v-model="selectedTypes"
                        class="w-full"
                        multiple
                        :items="activityTypeItems"
                        value-key="value"
                        :search-input="{ placeholder: $t('common.search_field') }"
                    />
                </UFormField>

                <UFormField :label="$t('common.colleague', 1)">
                    <SelectMenu
                        v-model="selectedUser"
                        class="w-full"
                        :searchable="
                            async (q: string) =>
                                await completorStore.completeColleagues(q, selectedUser?.userId ? [selectedUser.userId] : [])
                        "
                        searchable-key="jobs-group-activity-user"
                        :filter-fields="['firstname', 'lastname']"
                        :search-input="{ placeholder: $t('common.search_field') }"
                        :placeholder="$t('common.colleague', 1)"
                        clear
                    >
                        <template v-if="selectedUser" #default>
                            {{ userToLabel(selectedUser) }}
                        </template>
                        <template #item-label="{ item }">
                            {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                        </template>
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                        </template>
                    </SelectMenu>
                </UFormField>

                <UFormField :label="$t('common.date')">
                    <InputDateRangePopover v-model="dateRange" class="w-full" clearable time />
                </UFormField>

                <UFieldGroup class="inline-flex w-full sm:w-auto">
                    <UButton
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-filter-remove"
                        :label="$t('common.clear')"
                        @click="clearFilters"
                    />
                    <UButton icon="i-mdi-filter" :label="$t('common.apply')" @click="applyFilters" />
                </UFieldGroup>
            </div>
        </UCard>

        <DataErrorBlock
            v-if="activityError"
            :title="$t('common.unable_to_load', [$t('common.activity')])"
            :error="activityError"
            :retry="refreshActivity"
        />
        <div v-else-if="isRequestPending(activityStatus) || activityItems.length > 0">
            <ul class="divide-y divide-default" role="list">
                <template v-if="isRequestPending(activityStatus)">
                    <li v-for="idx in 5" :key="idx" class="px-2 py-4">
                        <div class="flex space-x-3">
                            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                                <USkeleton class="size-full" />
                            </div>

                            <div class="flex-1 space-y-1">
                                <div class="flex items-center justify-between gap-3">
                                    <h3 class="text-sm font-medium">
                                        <USkeleton class="h-5 w-[250px]" />
                                    </h3>

                                    <p>
                                        <USkeleton class="h-5 w-[150px]" />
                                    </p>
                                </div>

                                <div class="flex items-center justify-between gap-3">
                                    <p class="flex flex-col gap-1 text-sm">
                                        <USkeleton class="h-8 w-[200px]" />
                                    </p>
                                    <p class="inline-flex items-center gap-1 text-sm">
                                        <USkeleton class="h-5 w-[150px]" />
                                    </p>
                                </div>
                            </div>
                        </div>
                    </li>
                </template>

                <template v-else>
                    <li
                        v-for="entry in activityItems"
                        :key="entry.id"
                        class="border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                    >
                        <div class="flex space-x-3">
                            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                                <UIcon
                                    :class="[groupActivityTypeColor(entry.type), 'size-full']"
                                    :name="groupActivityTypeIcon(entry.type)"
                                    inline
                                />
                            </div>

                            <div class="flex-1 space-y-1">
                                <div class="flex items-center justify-between gap-3">
                                    <h3 class="text-sm font-medium">
                                        {{ activityLabel(entry) }}
                                    </h3>

                                    <p v-if="entry.createdAt" class="text-sm text-muted">
                                        <GenericTime :value="entry.createdAt" type="long" />
                                    </p>
                                </div>

                                <div class="flex items-center justify-between gap-3">
                                    <div class="flex flex-col gap-1 text-sm">
                                        <template v-if="entry.reason">
                                            <span class="inline-flex gap-1">
                                                <span class="font-semibold">{{ $t('common.reason') }}:</span>
                                                <span>{{ entry.reason }}</span>
                                            </span>
                                        </template>

                                        <template v-if="entry.targetUserId">
                                            <span class="inline-flex items-center gap-1">
                                                <span class="font-semibold">{{ $t('common.target') }}:</span>
                                                <ColleagueInfoPopover :user="entry.targetUser" :user-id="entry.targetUserId" />
                                            </span>
                                        </template>

                                        <template v-if="activityRuleLabel(entry)">
                                            <span class="inline-flex gap-1">
                                                <span class="font-semibold">
                                                    {{ $t('components.jobs.groups.details.rule') }}:
                                                </span>
                                                <span>{{ activityRuleLabel(entry) }}</span>
                                            </span>
                                        </template>

                                        <div
                                            v-if="
                                                entry.data?.data.oneofKind === 'rule' &&
                                                entry.data.data.rule.rule.oneofKind === 'qualification' &&
                                                entry.data.data.rule.rule.qualification.qualificationIds.length > 0
                                            "
                                            class="flex flex-wrap gap-1"
                                        >
                                            <QualificationBadge
                                                v-for="qualificationId in entry.data.data.rule.rule.qualification
                                                    .qualificationIds"
                                                :key="qualificationId"
                                                :qualification-id="qualificationId"
                                                :qualification="activityRuleQualification(entry, qualificationId)"
                                            />
                                        </div>
                                    </div>

                                    <p class="inline-flex items-center gap-1 text-sm">
                                        <span>{{ $t('common.created_by') }}</span>
                                        <ColleagueInfoPopover :user="entry.actorUser" :user-id="entry.actorUserId" />
                                    </p>
                                </div>
                            </div>
                        </div>
                    </li>
                </template>
            </ul>
        </div>
        <DataNoDataBlock v-else :type="$t('common.activity')" icon="i-mdi-history" :padded="false" />

        <Pagination v-model="page" :pagination="activity?.pagination" :status="activityStatus" :refresh="refreshActivity" />
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>

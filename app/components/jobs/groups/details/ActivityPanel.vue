<script lang="ts" setup>
import { z } from 'zod';
import ActivityPanelEntry from '~/components/jobs/groups/details/ActivityPanelEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDateRangePopover, { type DateRange } from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import { type GroupActivity, GroupActivityType } from '~~/gen/ts/resources/jobs/groups/activity';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListGroupActivityResponse } from '~~/gen/ts/services/jobs/groups';
import { groupActivityTypeColor, groupActivityTypeIcon, groupActivityTypeIndicatorColor } from '../helpers';

const props = defineProps<{
    groupId: number;
    canView: boolean;
}>();

const { t } = useI18n();
const completorStore = useCompletorStore();
const jobsGroupsClient = await getJobsGroupsClient();

const schema = z.object({
    types: z.enum(GroupActivityType).array().default([]),
    user: z.custom<UserShort>().optional(),
    dateRange: z.custom<DateRange>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({ types: [], user: undefined, dateRange: undefined });

const page = ref(1);

const activityKey = computed(
    () =>
        `jobs-group-activity-${props.groupId}-${page.value}-${state.types.join(',')}-${state.user?.userId ?? 0}-${
            state.dateRange?.start.toISOString() ?? ''
        }-${state.dateRange?.end.toISOString() ?? ''}`,
);

const {
    data: activity,
    status: activityStatus,
    error: activityError,
    refresh: refreshActivity,
} = useAuthedLazyAsyncData('userState', activityKey, ({ signal }) => listGroupActivity(signal), {
    watch: [() => props.groupId, page, () => state.types, () => state.user?.userId, () => state.dateRange],
});

const activityItems = computed<GroupActivity[]>(() => activity.value?.activity ?? []);
const timelineItems = computed(() =>
    activityItems.value.map((activity) => ({
        ...activity,
        icon: groupActivityTypeIcon(activity.type),
        ui: { indicator: groupActivityTypeIndicatorColor(activity.type) },
    })),
);

function activityTypeLabel(activityType: GroupActivityType): string {
    return t(`enums.jobs.groups.GroupActivityType.${GroupActivityType[activityType] ?? 'UNSPECIFIED'}`);
}

const activityTypeItems = computed(() =>
    Object.values(GroupActivityType)
        .filter((value): value is GroupActivityType => typeof value === 'number' && value !== GroupActivityType.UNSPECIFIED)
        .map((value) => ({
            label: activityTypeLabel(value),
            value,
            icon: groupActivityTypeIcon(value),
            ui: {
                itemLeadingIcon: groupActivityTypeColor(value),
            },
        })),
);

async function listGroupActivity(signal: AbortSignal): Promise<ListGroupActivityResponse> {
    const { response } = await jobsGroupsClient.listGroupActivity(
        {
            groupId: props.groupId,
            pagination: {
                offset: calculateOffset(page.value, activity.value?.pagination),
            },
            sort: { columns: [{ id: 'createdAt', desc: true }] },
            types: state.types,
            userId: state.user?.userId,
            from: toTimestamp(state.dateRange?.start),
            to: toTimestamp(state.dateRange?.end),
        },
        { abort: signal },
    );

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
    state.types = [];
    state.user = undefined;
    state.dateRange = undefined;
    if (page.value === 1) {
        await refreshActivity();
        return;
    }

    page.value = 1;
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
            <UForm :schema="schema" :state="state" @submit="applyFilters">
                <div class="grid gap-3 lg:grid-cols-[minmax(180px,260px)_minmax(0,1fr)_minmax(260px,1fr)_auto] lg:items-end">
                    <UFormField :label="$t('common.type')">
                        <USelectMenu
                            v-model="state.types"
                            class="w-full"
                            multiple
                            :items="activityTypeItems"
                            value-key="value"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        />
                    </UFormField>

                    <UFormField :label="$t('common.colleague', 1)">
                        <SelectMenu
                            v-model="state.user"
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeColleagues(q, state.user?.userId ? [state.user.userId] : [])
                            "
                            searchable-key="jobs-group-activity-user"
                            :filter-fields="['firstname', 'lastname']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :placeholder="$t('common.colleague', 1)"
                            clear
                        >
                            <template v-if="state.user" #default>
                                {{ userToLabel(state.user) }}
                            </template>
                            <template #item-label="{ item }">
                                {{ userToLabel(item) }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField :label="$t('common.date')">
                        <InputDateRangePopover v-model="state.dateRange" class="w-full" clearable time />
                    </UFormField>

                    <UFieldGroup class="inline-flex w-full sm:w-auto">
                        <UButton
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-filter-remove"
                            :label="$t('common.clear')"
                            @click="clearFilters"
                        />
                        <UButton type="submit" icon="i-mdi-filter" :label="$t('common.apply')" />
                    </UFieldGroup>
                </div>
            </UForm>
        </UCard>

        <DataErrorBlock
            v-if="activityError"
            :title="$t('common.unable_to_load', [$t('common.activity')])"
            :error="activityError"
            :retry="refreshActivity"
        />
        <div v-else-if="isRequestPending(activityStatus) || activityItems.length > 0">
            <template v-if="isRequestPending(activityStatus)">
                <ul class="divide-y divide-default" role="list">
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
                </ul>
            </template>

            <UTimeline v-else :items="timelineItems" size="md" :ui="{ wrapper: '!mt-0 !pb-2' }">
                <template #indicator="{ item }">
                    <UIcon :name="item.icon" class="size-6" />
                </template>

                <template #wrapper="{ item }">
                    <ActivityPanelEntry :activity="item" />
                </template>
            </UTimeline>
        </div>
        <DataNoDataBlock v-else :type="$t('common.activity')" icon="i-mdi-history" :padded="false" />

        <Pagination v-model="page" :pagination="activity?.pagination" :status="activityStatus" :refresh="refreshActivity" />
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>

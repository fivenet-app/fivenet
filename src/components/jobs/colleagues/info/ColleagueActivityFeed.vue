<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import { BulletinBoardIcon, CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import ColleagueActivityFeedEntry from '~/components/jobs/colleagues/info/ColleagueActivityFeedEntry.vue';
import type { ListColleagueActivityResponse } from '~~/gen/ts/services/jobs/jobs';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';

const props = withDefaults(
    defineProps<{
        userId?: number;
        showTargetUser?: boolean;
    }>(),
    {
        userId: undefined,
        showTargetUser: false,
    },
);

const { $grpc } = useNuxtApp();

const selectedUsers = ref<Colleague[]>([]);
const selectedUsersIds = computed(() =>
    props.userId !== undefined ? [props.userId] : selectedUsers.value.map((u) => u.userId),
);

const offset = ref(0);
const { data, pending, refresh, error } = useLazyAsyncData(
    `jobs-colleague-${selectedUsersIds.value.join(',')}-${offset.value}`,
    () => listColleagueActivity(selectedUsersIds.value),
);

async function listColleagueActivity(userIds: number[]): Promise<ListColleagueActivityResponse> {
    try {
        const call = $grpc.getJobsClient().listColleagueActivity({
            userIds,
            pagination: { offset: offset.value },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, () => refresh());

const queryColleagueNameRaw = ref<string>('');
const queryColleagueName = computed(() => queryColleagueNameRaw.value.trim());

const { data: colleagues, refresh: refreshColleagues } = useLazyAsyncData(
    `jobs-colleagues-${offset.value}-${queryColleagueName.value}`,
    async () => {
        try {
            const call = $grpc.getJobsClient().listColleagues({
                pagination: {
                    offset: offset.value,
                },
                searchName: queryColleagueName.value,
            });
            const { response } = await call;

            return response;
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
        }
    },
);

watchDebounced(
    queryColleagueName,
    async () => {
        await refreshColleagues();
        if (props.userId === undefined && selectedUsers.value) {
            colleagues.value?.colleagues.unshift(...selectedUsers.value);
        }
    },
    {
        debounce: 500,
        maxWait: 1250,
    },
);

const accessAttrs = attrList('JobsService.GetColleague', 'Access');
const colleagueSearchAttrs = ['own', 'lower_rank', 'same_rank', 'any'];

watch(props, async () => refresh());
watchDebounced(selectedUsers, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});

function charsGetDisplayValue(chars: Colleague[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}
</script>

<template>
    <div class="py-2 pb-4">
        <div class="px-1 sm:px-2 lg:px-4">
            <div
                v-if="userId === undefined && accessAttrs.some((a) => colleagueSearchAttrs.includes(a))"
                class="mb-4 sm:flex sm:items-center"
            >
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="flex-1">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <Combobox v-model="selectedUsers" as="div" class="mt-2 w-full" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (chars: any) => (chars ? charsGetDisplayValue(chars) : $t('common.na'))
                                                    "
                                                    :placeholder="$t('common.target')"
                                                    @change="queryColleagueNameRaw = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="colleagues?.colleagues !== undefined && colleagues.colleagues.length > 0"
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="colleague in colleagues?.colleagues"
                                                    :key="colleague.identifier"
                                                    v-slot="{ active, selected }"
                                                    :value="colleague"
                                                    as="char"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ colleague.firstname }} {{ colleague.lastname }} ({{
                                                                colleague?.dateofbirth
                                                            }})
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="size-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 align-middle">
                        <DataPendingBlock
                            v-if="pending"
                            :message="$t('common.loading', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
                        />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.not_found', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.activity.length === 0"
                            :icon="markRaw(BulletinBoardIcon)"
                            :type="`${$t('common.colleague', 1)} ${$t('common.activity')}`"
                        />
                        <div v-else>
                            <ul role="list" class="divide-y divide-gray-200">
                                <li v-for="activity in data?.activity" :key="activity.id" class="py-4">
                                    <ColleagueActivityFeedEntry :activity="activity" :show-target-user="showTargetUser" />
                                </li>
                            </ul>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

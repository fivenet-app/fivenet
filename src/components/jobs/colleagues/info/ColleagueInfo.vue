<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { BulletinBoardIcon, CloseIcon, IslandIcon, ListStatusIcon, MenuIcon, SchoolIcon, TimelineClockIcon } from 'mdi-vue3';
import type { DefineComponent } from 'vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { GetColleagueResponse } from '~~/gen/ts/services/jobs/jobs';
import ColleagueActivityFeed from '~/components/jobs/colleagues/info/ColleagueActivityFeed.vue';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';
import ConductList from '~/components/jobs/conduct/ConductList.vue';
import type { Perms } from '~~/gen/ts/perms';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';
import GenericTime from '~/components/partials/elements/GenericTime.vue';

const props = defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const {
    data: colleague,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`jobs-colleague-${props.userId}`, () => getColleague(props.userId));

async function getColleague(userId: number): Promise<GetColleagueResponse> {
    try {
        const call = $grpc.getJobsClient().getColleague({
            userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const tabs: { id: string; name: string; icon: DefineComponent; permission: Perms }[] = [
    {
        id: 'activity',
        name: 'common.activity',
        icon: markRaw(BulletinBoardIcon),
        permission: 'JobsService.ListColleagueActivity' as Perms,
    },
    {
        id: 'timeclock',
        name: 'common.timeclock',
        icon: markRaw(TimelineClockIcon),
        permission: 'JobsTimeclockService.ListTimeclock' as Perms,
    },
    {
        id: 'qualifications',
        name: 'pages.qualifications.title',
        to: { name: 'jobs-qualifications' },
        icon: markRaw(SchoolIcon),
        permission: 'QualificationsService.ListQualifications' as Perms,
    },
    {
        id: 'conduct',
        name: 'pages.jobs.conduct.title',
        icon: markRaw(ListStatusIcon),
        permission: 'JobsConductService.ListConductEntries' as Perms,
    },
].filter((tab) => can(tab.permission));

const selectedTab = ref(0);

function changeTab(index: number) {
    selectedTab.value = index;
}

const open = ref(false);

const absenceDateModal = ref(false);

const today = new Date();
today.setHours(0);
today.setMinutes(0);
today.setSeconds(0);
today.setMilliseconds(0);
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.colleague', 1)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.colleague', 1)])"
            :message="$t(error.message)"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="colleague === null || !colleague.colleague" />

        <template v-else>
            <SelfServicePropsAbsenceDateModal
                :open="absenceDateModal"
                :user-id="colleague.colleague.userId"
                :user-props="colleague.colleague.props"
                @close="absenceDateModal = false"
            />

            <div class="mb-6">
                <div class="flex gap-4 my-4 px-4">
                    <ProfilePictureImg
                        :url="colleague.colleague.avatar?.url"
                        :name="`${colleague.colleague.firstname} ${colleague.colleague.lastname}`"
                        size="xl"
                        :rounded="false"
                        :enable-popup="true"
                    />
                    <div class="w-full">
                        <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                            <h1 class="flex-1 break-words py-1 pl-0.5 pr-0.5 text-4xl font-bold text-neutral sm:pl-1">
                                {{ colleague.colleague.firstname }} {{ colleague.colleague.lastname }}
                            </h1>

                            <button
                                v-if="
                                    can('JobsService.SetJobsUserProps') &&
                                    checkIfCanAccessColleague(activeChar!, colleague.colleague, 'JobsService.SetJobsUserProps')
                                "
                                type="button"
                                class="place-self-end inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                @click="absenceDateModal = true"
                            >
                                <IslandIcon class="w-5 h-auto" aria-hidden="true" />
                                {{ $t('components.jobs.self_service.set_absence_date') }}
                            </button>
                        </div>
                        <div class="my-2 flex flex-row items-center gap-2">
                            <span class="rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800">
                                {{ colleague.colleague.jobLabel }}
                                <span v-if="colleague.colleague.jobGrade > 0">
                                    ({{ $t('common.rank') }}: {{ colleague.colleague.jobGradeLabel }})</span
                                >
                            </span>

                            <span
                                v-if="
                                    colleague.colleague.props?.absenceEnd &&
                                    toDate(colleague.colleague.props?.absenceEnd).getTime() >= today.getTime()
                                "
                                class="inline-flex gap-1 items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800"
                            >
                                <IslandIcon class="h-5 w-5" aria-hidden="true" />
                                <GenericTime :value="colleague.colleague.props?.absenceBegin" type="date" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="colleague.colleague.props?.absenceEnd" type="date" />
                            </span>
                        </div>
                    </div>
                </div>

                <nav class="bg-base-700 lg:rounded-lg">
                    <div class="mx-auto ml-2 max-w-7xl px-4 sm:px-6 lg:px-8">
                        <div class="flex h-16 items-center justify-between">
                            <div class="flex items-center md:overflow-x-auto">
                                <div class="-ml-2 flex md:hidden">
                                    <!-- Mobile menu button -->
                                    <button
                                        type="button"
                                        class="relative inline-flex items-center justify-center rounded-md bg-base-500 p-2 text-accent-200 hover:bg-base-400 hover:bg-opacity-75 hover:text-neutral focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2 focus:ring-offset-base-600"
                                        @click="open = !open"
                                    >
                                        <span class="absolute -inset-0.5" />
                                        <span class="sr-only">{{ $t('components.partials.sidebar.open_navigation') }}</span>
                                        <MenuIcon v-if="!open" class="block h-5 w-5" aria-hidden="true" />
                                        <CloseIcon v-else class="block h-5 w-5" aria-hidden="true" />
                                    </button>
                                </div>
                                <div class="md:block hidden">
                                    <div class="flex items-baseline space-x-2">
                                        <template v-for="(tab, index) in tabs" :key="tab.id">
                                            <span class="flex-1">
                                                <button
                                                    type="button"
                                                    class="group flex shrink-0 items-center gap-2 rounded-md p-3 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                                    :class="
                                                        selectedTab === index
                                                            ? 'bg-accent-100/20 font-bold text-primary-300'
                                                            : ''
                                                    "
                                                    @click="selectedTab = index"
                                                >
                                                    <component
                                                        :is="tab.icon"
                                                        :class="[
                                                            selectedTab === index ? '' : 'group-hover:text-base-300',
                                                            'h-5 w-5',
                                                        ]"
                                                        aria-hidden="true"
                                                    />
                                                    <span>
                                                        {{ $t(tab.name) }}
                                                    </span>
                                                </button>
                                            </span>
                                        </template>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="-ml-3 -mr-3 md:hidden" :class="open ? 'block' : 'hidden'">
                            <div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
                                <template v-for="(tab, index) in tabs" :key="tab.id">
                                    <button
                                        type="button"
                                        class="group flex w-full shrink-0 items-center items-center gap-2 rounded-md p-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                        :class="selectedTab === index ? 'bg-accent-100/20 font-bold text-primary-300' : ''"
                                        @click="selectedTab = index"
                                    >
                                        <component
                                            :is="tab.icon"
                                            :class="[selectedTab === index ? '' : 'group-hover:text-base-300', 'h-5 w-5']"
                                            aria-hidden="true"
                                        />
                                        <span>
                                            {{ $t(tab.name) }}
                                        </span>
                                    </button>
                                </template>
                            </div>
                        </div>
                    </div>
                </nav>

                <TabGroup :selected-index="selectedTab" @change="changeTab">
                    <TabList class="hidden">
                        <Tab v-for="tab in tabs" :key="tab.id"></Tab>
                    </TabList>
                    <TabPanels class="mt-2 bg-transparent">
                        <TabPanel v-if="can('JobsService.GetColleague')">
                            <ColleagueActivityFeed :user-id="userId" />
                        </TabPanel>
                        <TabPanel v-if="can('JobsTimeclockService.ListTimeclock')">
                            <TimeclockOverviewBlock :user-id="userId" />
                        </TabPanel>
                        <TabPanel v-if="can('QualificationsService.ListQualifications')">
                            <JobsQualificationsResultsList class="mt-4" :user-id="userId" />
                        </TabPanel>
                        <TabPanel v-if="can('JobsConductService.ListConductEntries')">
                            <ConductList :user-id="userId" :hide-user-search="true" />
                        </TabPanel>
                    </TabPanels>
                </TabGroup>
            </div>
        </template>
    </div>
</template>

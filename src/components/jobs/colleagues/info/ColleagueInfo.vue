<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { BulletinBoardIcon, CloseIcon, IslandIcon, ListStatusIcon, MenuIcon, SchoolIcon, TimelineClockIcon } from 'mdi-vue3';
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

const { t } = useI18n();

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

const tabs: { key: string; label: string; icon: string; permission: Perms }[] = [
    {
        key: 'activity',
        label: t('common.activity'),
        icon: 'i-mdi-bulletin-board',
        permission: 'JobsService.ListColleagueActivity' as Perms,
    },
    {
        key: 'timeclock',
        label: t('common.timeclock'),
        icon: 'i-mdi-timeline-clock',
        permission: 'JobsTimeclockService.ListTimeclock' as Perms,
    },
    {
        key: 'qualifications',
        label: t('pages.qualifications.title'),
        icon: 'i-mdi-school',
        permission: 'QualificationsService.ListQualifications' as Perms,
    },
    {
        key: 'conduct',
        label: t('pages.jobs.conduct.title'),
        icon: 'i-mdi-list-status',
        permission: 'JobsConductService.ListConductEntries' as Perms,
    },
].filter((tab) => can(tab.permission));

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
                <div class="my-4 flex gap-4 px-4">
                    <ProfilePictureImg
                        :url="colleague.colleague.avatar?.url"
                        :name="`${colleague.colleague.firstname} ${colleague.colleague.lastname}`"
                        size="xl"
                        :rounded="false"
                        :enable-popup="true"
                    />
                    <div class="w-full">
                        <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                            <h1 class="flex-1 break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                                {{ colleague.colleague.firstname }} {{ colleague.colleague.lastname }}
                            </h1>

                            <UButton
                                v-if="
                                    can('JobsService.SetJobsUserProps') &&
                                    checkIfCanAccessColleague(activeChar!, colleague.colleague, 'JobsService.SetJobsUserProps')
                                "
                                class="inline-flex items-center gap-x-1.5 place-self-end rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400"
                                @click="absenceDateModal = true"
                            >
                                <IslandIcon class="h-auto w-5" />
                                {{ $t('components.jobs.self_service.set_absence_date') }}
                            </UButton>
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
                                class="inline-flex items-center gap-1 rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800"
                            >
                                <IslandIcon class="size-5" />
                                <GenericTime :value="colleague.colleague.props?.absenceBegin" type="date" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="colleague.colleague.props?.absenceEnd" type="date" />
                            </span>
                        </div>
                    </div>
                </div>

                <UTabs :items="tabs" class="w-full">
                    <template #default="{ item, selected }">
                        <div class="flex items-center gap-2 relative truncate">
                            <UIcon :name="item.icon" class="w-4 h-4 flex-shrink-0" />

                            <span class="truncate">{{ item.label }}</span>

                            <span
                                v-if="selected"
                                class="absolute -right-4 w-2 h-2 rounded-full bg-primary-500 dark:bg-primary-400"
                            />
                        </div>
                    </template>

                    <template #item="{ item }">
                        <div v-if="item.key === 'activity'" class="space-y-3">
                            <ColleagueActivityFeed :user-id="userId" />
                        </div>
                        <div v-else-if="item.key === 'timeclock'" class="space-y-3">
                            <TimeclockOverviewBlock :user-id="userId" />
                        </div>
                        <div v-else-if="item.key === 'qualifications'" class="space-y-3">
                            <JobsQualificationsResultsList class="mt-4" :user-id="userId" />
                        </div>
                        <div v-else-if="item.key === 'conduct'" class="space-y-3">
                            <ConductList :user-id="userId" :hide-user-search="true" />
                        </div>
                    </template>
                </UTabs>
            </div>
        </template>
    </div>
</template>

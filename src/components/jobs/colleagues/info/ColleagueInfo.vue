<script lang="ts" setup>
import { IslandIcon } from 'mdi-vue3';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { GetColleagueResponse } from '~~/gen/ts/services/jobs/jobs';
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

const modal = useModal();

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
            <div class="mb-6">
                <div class="my-4 flex gap-4 px-4">
                    <ProfilePictureImg
                        :url="colleague.colleague.avatar?.url"
                        :name="`${colleague.colleague.firstname} ${colleague.colleague.lastname}`"
                        size="xl"
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
                                icon="i-mdi-island"
                                size="md"
                                @click="
                                    modal.open(SelfServicePropsAbsenceDateModal, {
                                        userId: colleague.colleague.userId,
                                        userProps: colleague.colleague.props,
                                    })
                                "
                            >
                                {{ $t('components.jobs.self_service.set_absence_date') }}
                            </UButton>
                        </div>
                        <div class="my-2 flex flex-row items-center gap-2">
                            <UBadge>
                                {{ colleague.colleague.jobLabel }}
                                <span v-if="colleague.colleague.jobGrade > 0" class="ml-1">
                                    ({{ $t('common.rank') }}: {{ colleague.colleague.jobGradeLabel }})</span
                                >
                            </UBadge>

                            <UBadge
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
                            </UBadge>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>

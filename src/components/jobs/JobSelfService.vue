<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CameraIcon, IslandIcon } from 'mdi-vue3';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import SelfServicePropsProfilePictureModal from '~/components/jobs/colleagues/SelfServicePropsProfilePictureModal.vue';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';

defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { data: colleagueSelf } = useLazyAsyncData('jobs-selfcolleague', async () => {
    try {
        const call = $grpc.getJobsClient().getSelf({});
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
});

function updateAbsenceDate(value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void {
    if (colleagueSelf.value === null) {
        return;
    }

    if (colleagueSelf.value.colleague!.props === undefined) {
        colleagueSelf.value.colleague!.props = {
            userId: colleagueSelf.value!.colleague!.userId,
            absenceBegin: value.absenceBegin,
            absenceEnd: value.absenceEnd,
        };
    } else {
        colleagueSelf.value.colleague!.props.absenceBegin = value.absenceBegin;
        colleagueSelf.value.colleague!.props.absenceEnd = value.absenceEnd;
    }
}

const absenceDateModal = ref(false);
const profilePictureModal = ref(false);
</script>

<template>
    <GenericContainer class="flex-1 text-neutral">
        <SelfServicePropsAbsenceDateModal
            :open="absenceDateModal"
            :user-id="userId"
            :user-props="colleagueSelf?.colleague?.props"
            @close="absenceDateModal = false"
            @update:absence-dates="updateAbsenceDate($event)"
        />
        <SelfServicePropsProfilePictureModal :open="profilePictureModal" @close="profilePictureModal = false" />

        <h3 class="text-lg font-semibold">{{ $t('components.jobs.self_service.title') }}</h3>

        <div class="flex flex-col md:flex-row items-center flex-initial gap-1">
            <button
                v-if="
                    colleagueSelf?.colleague &&
                    can('JobsService.SetJobsUserProps') &&
                    checkIfCanAccessColleague(activeChar!, colleagueSelf.colleague, 'JobsService.SetJobsUserProps')
                "
                type="button"
                class="w-full inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="absenceDateModal = true"
            >
                <IslandIcon class="mr-2 h-5 w-auto" aria-hidden="true" />
                <span>{{ $t('components.jobs.self_service.set_absence_date') }}</span>
            </button>
            <button
                type="button"
                class="w-full inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="profilePictureModal = true"
            >
                <CameraIcon class="mr-2 h-5 w-auto" aria-hidden="true" />
                <span>{{ $t('components.jobs.self_service.set_profile_picture') }}</span>
            </button>
        </div>
    </GenericContainer>
</template>

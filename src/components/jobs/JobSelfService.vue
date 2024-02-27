<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CameraIcon, IslandIcon } from 'mdi-vue3';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import SelfServicePropsProfilePictureModal from '~/components/jobs/colleagues/SelfServicePropsProfilePictureModal.vue';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

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

function updateAbsenceDate(value?: Timestamp): void {
    if (colleagueSelf.value === null) {
        return;
    }

    if (colleagueSelf.value.colleague!.props === undefined) {
        colleagueSelf.value.colleague!.props = {
            userId: colleagueSelf.value!.colleague!.userId,
            absenceDate: value,
        };
    } else {
        colleagueSelf.value.colleague!.props.absenceDate = value;
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
            @update:absence-date="updateAbsenceDate($event)"
        />
        <SelfServicePropsProfilePictureModal :open="profilePictureModal" @close="profilePictureModal = false" />

        <h3 class="text-lg font-semibold">{{ $t('components.jobs.self_service.title') }}</h3>

        <div class="flex items-center flex-initial gap-1">
            <button
                type="button"
                class="w-full inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="absenceDateModal = true"
            >
                <IslandIcon class="mr-2 h-5 w-auto" />
                <span>{{ $t('components.jobs.self_service.set_absence_date') }}</span>
            </button>
            <button
                type="button"
                class="w-full inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="profilePictureModal = true"
            >
                <CameraIcon class="mr-2 h-5 w-auto" />
                <span>{{ $t('components.jobs.self_service.set_profile_picture') }}</span>
            </button>
        </div>
    </GenericContainer>
</template>

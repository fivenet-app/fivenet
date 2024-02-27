<script lang="ts" setup>
import { CameraIcon, IslandIcon } from 'mdi-vue3';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import JobsUserPropsAbsenceDateModal from '~/components/jobs/colleagues/JobsUserPropsAbsenceDateModal.vue';

defineProps<{
    userId: number;
}>();

const absenceDateModal = ref(false);
const profilePictureModal = ref(false);
</script>

<template>
    <GenericContainer class="flex-1 text-neutral">
        <JobsUserPropsAbsenceDateModal
            :open="absenceDateModal"
            :user-props="{
                userId: 0,
            }"
            @close="absenceDateModal = false"
        />

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

<script lang="ts" setup>
import { RadioHandheldIcon } from 'mdi-vue3';
import JobMotd from '~/components/jobs/JobMotd.vue';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';
import SquareImg from '~/components/partials/SquareImg.vue';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.jobs.overview.title',
});
definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
});

const authStore = useAuthStore();

const { activeChar, jobProps } = storeToRefs(authStore);

const showRadioFrequency = ref(false);
</script>

<template>
    <div>
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="grid-col-2 grid gap-2">
                    <div class="sm:flex-auto">
                        <div class="flex flex-row gap-2">
                            <GenericContainer class="flex-1">
                                <div class="flex flex-row gap-4">
                                    <SquareImg
                                        v-if="jobProps && jobProps.logoUrl"
                                        :url="jobProps?.logoUrl.url"
                                        :text="$t('common.logo')"
                                        size="xl"
                                    />

                                    <div>
                                        <h1 class="text-3xl font-semibold leading-6 text-neutral">
                                            {{ activeChar?.jobLabel }}
                                        </h1>
                                        <h2 class="mt-2 text-xl font-semibold leading-6 text-neutral">
                                            {{ $t('common.rank') }}: {{ activeChar?.jobGradeLabel }}
                                        </h2>
                                    </div>
                                </div>

                                <JobMotd />
                            </GenericContainer>

                            <GenericContainer v-if="jobProps?.radioFrequency" class="text-neutral">
                                <h3 class="text-lg font-semibold">
                                    {{ $t('common.radio_frequency') }}
                                </h3>
                                <p class="flex items-center text-center text-lg font-bold">
                                    <RadioHandheldIcon class="h-auto w-6" />
                                    <span
                                        :class="showRadioFrequency ? '' : 'blur-md'"
                                        @click="showRadioFrequency = !showRadioFrequency"
                                        >{{ jobProps?.radioFrequency }}.00</span
                                    >
                                </p>
                                <button
                                    v-if="isNUIAvailable()"
                                    type="button"
                                    class="mt-1 w-full rounded-md bg-primary-500 px-2 py-1 text-xs font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                    @click="setRadioFrequency(jobProps.radioFrequency)"
                                >
                                    {{ $t('common.connect') }}
                                </button>
                            </GenericContainer>
                        </div>

                        <TimeclockOverviewBlock v-if="can('JobsTimeclockService.ListTimeclock')" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

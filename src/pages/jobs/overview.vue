<script lang="ts" setup>
import JobMotd from '~/components/jobs/JobMotd.vue';
import JobSelfService from '~/components/jobs/JobSelfService.vue';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';
import SquareImg from '~/components/partials/elements/SquareImg.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.jobs.overview.title',
});
definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
});

const authStore = useAuthStore();

const { activeChar, jobProps } = storeToRefs(authStore);

const { game } = useAppConfig();

const showRadioFrequency = ref(false);
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <div class="px-2 py-2">
                <div class="grid gap-2">
                    <div class="flex flex-row gap-2">
                        <UCard class="flex-1">
                            <template #header>
                                <div class="flex flex-row gap-2">
                                    <SquareImg
                                        v-if="jobProps && jobProps.logoUrl"
                                        :src="jobProps?.logoUrl.url"
                                        :alt="`${jobProps.jobLabel} ${$t('common.logo')}`"
                                        size="xl"
                                        :no-blur="true"
                                    />

                                    <div>
                                        <h1 class="text-3xl font-semibold leading-6">
                                            {{ activeChar?.jobLabel }}
                                        </h1>
                                        <h2
                                            v-if="activeChar?.job !== game.unemployedJobName"
                                            class="mt-2 text-xl font-semibold leading-6"
                                        >
                                            {{ $t('common.rank') }}: {{ activeChar?.jobGradeLabel }}
                                        </h2>
                                    </div>
                                </div>
                            </template>

                            <JobMotd />
                        </UCard>

                        <UCard v-if="jobProps?.radioFrequency">
                            <template #header>
                                <h3 class="inline-flex items-center gap-1 text-lg font-semibold">
                                    <UIcon name="i-mdi-radio-handheld" class="size-7" />

                                    <span>{{ $t('common.radio_frequency') }}</span>
                                </h3>
                            </template>

                            <div class="flex flex-col gap-2">
                                <div class="flex items-center justify-center text-lg font-bold">
                                    <span
                                        :class="showRadioFrequency ? '' : 'blur'"
                                        @click="showRadioFrequency = !showRadioFrequency"
                                        >{{ jobProps?.radioFrequency }}.00</span
                                    >
                                </div>

                                <UButton
                                    v-if="isNUIAvailable()"
                                    block
                                    variant="soft"
                                    @click="setRadioFrequency(jobProps.radioFrequency)"
                                >
                                    {{ $t('common.connect') }}
                                </UButton>
                            </div>
                        </UCard>
                    </div>

                    <JobSelfService v-if="can('JobsService.ListColleagues') && activeChar" :user-id="activeChar.userId" />

                    <TimeclockOverviewBlock v-if="can('JobsTimeclockService.ListTimeclock')" />
                </div>
            </div>
        </template>
    </PagesJobsLayout>
</template>

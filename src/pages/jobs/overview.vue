<script lang="ts" setup>
import JobMotd from '~/components/jobs/JobMotd.vue';
import JobSelfService from '~/components/jobs/JobSelfService.vue';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';
import SquareImg from '~/components/partials/elements/SquareImg.vue';
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
        <div class="px-1 py-2 sm:px-2 lg:px-4">
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
                                    <h2 class="mt-2 text-xl font-semibold leading-6">
                                        {{ $t('common.rank') }}: {{ activeChar?.jobGradeLabel }}
                                    </h2>
                                </div>
                            </div>
                        </template>

                        <JobMotd />
                    </UCard>

                    <UCard v-if="jobProps?.radioFrequency">
                        <template #header>
                            <h3 class="text-lg font-semibold">
                                {{ $t('common.radio_frequency') }}
                            </h3>
                        </template>

                        <div class="flex flex-col gap-2">
                            <div class="flex items-center text-center text-lg font-bold">
                                <UIcon name="i-mdi-radio-handheld" class="h-auto w-7" />
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

                <div v-if="activeChar" class="flex flex-row gap-2">
                    <JobSelfService :user-id="activeChar.userId" />
                </div>

                <TimeclockOverviewBlock v-if="can('JobsTimeclockService.ListTimeclock')" />
            </div>
        </div>
    </div>
</template>

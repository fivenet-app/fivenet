<script lang="ts" setup>
import JobMotd from '~/components/jobs/JobMotd.vue';
import JobSelfService from '~/components/jobs/JobSelfService.vue';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import { useSettingsStore } from '~/stores/settings';

useHead({
    title: 'pages.jobs.overview.title',
});

definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
});

const { can, activeChar, jobProps } = useAuth();

const { game } = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const showRadioFrequency = ref(false);
</script>

<template>
    <div class="overflow-y-auto px-2 py-2">
        <div class="grid gap-2">
            <div class="flex flex-row gap-2">
                <UCard class="flex-1">
                    <template #header>
                        <div class="flex flex-row gap-2">
                            <GenericImg
                                v-if="jobProps && jobProps.logoFile"
                                :src="jobProps.logoFile.filePath"
                                :alt="`${jobProps.jobLabel} ${$t('common.logo')}`"
                                size="xl"
                                :no-blur="true"
                            />

                            <div>
                                <h1 class="text-3xl leading-6 font-semibold">
                                    {{ activeChar?.jobLabel }}
                                </h1>
                                <h2
                                    v-if="activeChar?.job !== game.unemployedJobName"
                                    class="mt-2 text-xl leading-6 font-semibold"
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
                            <UIcon class="size-7" name="i-mdi-radio-handheld" />

                            <span>{{ $t('common.radio_frequency') }}</span>
                        </h3>
                    </template>

                    <div class="flex flex-col gap-2">
                        <div class="flex items-center justify-center text-lg font-bold">
                            <span :class="showRadioFrequency ? '' : 'blur'" @click="showRadioFrequency = !showRadioFrequency"
                                >{{ jobProps?.radioFrequency }}.00</span
                            >
                        </div>

                        <UButton v-if="nuiEnabled" block variant="soft" @click="setRadioFrequency(jobProps.radioFrequency)">
                            {{ $t('common.connect') }}
                        </UButton>
                    </div>
                </UCard>
            </div>

            <JobSelfService
                v-if="can('jobs.JobsService/ListColleagues').value && activeChar !== null"
                :user-id="activeChar.userId"
            />

            <TimeclockOverviewBlock v-if="can('jobs.TimeclockService/ListTimeclock').value" />
        </div>
    </div>
</template>

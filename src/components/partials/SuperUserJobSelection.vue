<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { Job } from '~~/gen/ts/resources/users/jobs';

const authStore = useAuthStore();
const { setSuperUserMode } = authStore;
const { activeChar, isSuperuser } = storeToRefs(authStore);

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const selectedJob = ref<undefined | Job>(jobs.value.find((j) => j.name === activeChar.value?.job));

watchOnce(jobs, () => (selectedJob.value = jobs.value.find((j) => j.name === activeChar.value?.job)));

watch(selectedJob, () => {
    if (activeChar.value?.job === selectedJob.value?.name) {
        return;
    }

    setSuperUserMode(isSuperuser.value, selectedJob.value);
});
</script>

<template>
    <UInputMenu
        v-model="selectedJob"
        class="relative"
        option-attribute="label"
        :search-attributes="['name', 'label']"
        :options="jobs"
        :popper="{ placement: 'top' }"
        :placeholder="selectedJob ? selectedJob.label : $t('common.na')"
        :search="
            async (q?: string) => (await listJobs()).filter((j) => q === undefined || j.name.includes(q) || j.label.includes(q))
        "
    />
</template>

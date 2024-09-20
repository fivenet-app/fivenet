<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import type { Job } from '~~/gen/ts/resources/users/jobs';

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
        :placeholder="$t('common.job')"
        :search="
            async (q?: string) => (await listJobs()).filter((j) => q === undefined || j.name.includes(q) || j.label.includes(q))
        "
        search-lazy
        :search-placeholder="$t('common.search_field')"
    >
        <template #option-empty="{ query: search }">
            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
        </template>
        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
    </UInputMenu>
</template>

<script lang="ts" setup>
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';

defineOptions({
    inheritAttrs: false,
});

const authStore = useAuthStore();
const { setSuperuserMode } = authStore;
const { activeChar, isSuperuser } = storeToRefs(authStore);

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const selectedJob = ref<undefined | Job>(
    jobs.value.find((j) => j.name === activeChar.value?.job) ?? {
        name: activeChar.value?.job ?? 'na',
        label: activeChar.value?.jobLabel ?? 'N/A',
        grades: [],
    },
);

watchOnce(jobs, () => (selectedJob.value = jobs.value.find((j) => j.name === activeChar.value?.job)));

watch(selectedJob, async () => {
    if (activeChar.value?.job === selectedJob.value?.name) {
        return;
    }

    await setSuperuserMode(isSuperuser.value, selectedJob.value);
});
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model="selectedJob"
            class="relative"
            :filter-fields="['name', 'label']"
            :items="jobs"
            :placeholder="$t('common.job', 1)"
            :search="
                async (q?: string) => {
                    q = q?.toLowerCase()?.trim();
                    return (await listJobs()).filter(
                        (j) => q === undefined || j.name.toLowerCase().includes(q) || j.label.toLowerCase().includes(q),
                    );
                }
            "
            :ui-menu="{ height: 'max-h-40' }"
            v-bind="$attrs"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
        </UInputMenu>
    </ClientOnly>
</template>

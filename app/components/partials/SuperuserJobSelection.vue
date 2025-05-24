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

watch(selectedJob, () => {
    if (activeChar.value?.job === selectedJob.value?.name) {
        return;
    }

    setSuperuserMode(isSuperuser.value, selectedJob.value);
});
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model="selectedJob"
            class="relative"
            option-attribute="label"
            :search-attributes="['name', 'label']"
            :options="jobs"
            :popper="{ placement: 'top' }"
            :placeholder="$t('common.job', 1)"
            :search="
                async (q?: string) => {
                    q = q?.toLowerCase()?.trim();
                    return (await listJobs()).filter(
                        (j) => q === undefined || j.name.toLowerCase().includes(q) || j.label.toLowerCase().includes(q),
                    );
                }
            "
            search-lazy
            :search-placeholder="$t('common.search_field')"
            :ui-menu="{ height: 'max-h-40' }"
            v-bind="$attrs"
        >
            <template #option-empty="{ query: search }">
                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
            </template>

            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
        </UInputMenu>
    </ClientOnly>
</template>

<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';

defineProps<{
    collapsed?: boolean | undefined;
}>();

defineOptions({
    inheritAttrs: false,
});

const { t } = useI18n();

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

watch(activeChar, () => {
    if (!activeChar.value) {
        selectedJob.value = undefined;
        return;
    }

    selectedJob.value = jobs.value.find((j) => j.name === activeChar.value?.job) ?? {
        name: activeChar.value?.job ?? 'na',
        label: activeChar.value?.jobLabel ?? 'N/A',
        grades: [],
    };
});

watchOnce(jobs, () => (selectedJob.value = jobs.value.find((j) => j.name === activeChar.value?.job)));

watch(selectedJob, async () => {
    if (activeChar.value?.job === selectedJob.value?.name) return;

    await setSuperuserMode(isSuperuser.value, selectedJob.value);
});

const items = computed(
    () =>
        [
            {
                label: t('common.superuser') + ' (' + (isSuperuser.value ? t('common.active') : t('common.inactive')) + ')',
                icon: 'i-mdi-square-root',
                type: 'link' as const,
                active: isSuperuser.value,
                onSelect: () => {
                    authStore.setSuperuserMode(!isSuperuser.value);
                },
            },
        ] satisfies NavigationMenuItem[],
);

onBeforeMount(async () => listJobs());
</script>

<template>
    <UInputMenu
        v-if="isSuperuser"
        v-model="selectedJob"
        class="relative -mb-3.5"
        variant="soft"
        :filter-fields="['name', 'label']"
        :placeholder="$t('common.job', 1)"
        :items="jobs"
        searchable-key="superuser-job-selection"
        autocomplete="off"
        name="job"
        v-bind="$attrs"
    >
        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
    </UInputMenu>

    <UNavigationMenu orientation="vertical" tooltip popover :collapsed="collapsed" :items="items" />
</template>

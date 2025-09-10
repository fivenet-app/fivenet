<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import InputMenu from './InputMenu.vue';

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

watchOnce(jobs, () => (selectedJob.value = jobs.value.find((j) => j.name === activeChar.value?.job)));

watch(selectedJob, async () => {
    if (activeChar.value?.job === selectedJob.value?.name) return;

    await setSuperuserMode(isSuperuser.value, selectedJob.value);
});

const items = computed(
    () =>
        [
            {
                label: t('common.superuser'),
                icon: 'i-mdi-square-root',
                type: 'link' as const,
                active: isSuperuser.value,
                onSelect: () => {
                    authStore.setSuperuserMode(!isSuperuser.value);
                },
            },
        ] satisfies NavigationMenuItem[],
);
</script>

<template>
    <InputMenu
        v-if="isSuperuser"
        v-model="selectedJob"
        class="relative -mb-3.5"
        variant="soft"
        :filter-fields="['name', 'label']"
        :placeholder="$t('common.job', 1)"
        :items="jobs"
        :searchable="
            async (q?: string) => {
                q = q?.toLowerCase()?.trim();
                return (await listJobs()).filter(
                    (j) => q === undefined || j.name.toLowerCase().includes(q) || j.label.toLowerCase().includes(q),
                );
            }
        "
        searchable-key="superuser-job-selection"
        v-bind="$attrs"
    >
        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
    </InputMenu>

    <UNavigationMenu orientation="vertical" tooltip popover :collapsed="collapsed" :items="items" />
</template>

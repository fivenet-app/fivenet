<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import JobAssetsFileList from '~/components/jobs/JobAssetsFileList.vue';
import type { JobAsset } from '~~/gen/ts/resources/jobs/job_asset';
import type { ListJobAssetsResponse } from '~~/gen/ts/services/settings/job_assets';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';

const { can } = useAuth();
const canUpdate = can('settings.JobsAssetsService/UpdateJobAsset');
const canDelete = can('settings.JobsAssetsService/DeleteJobAsset');

async function refreshAssets(): Promise<void> {
    await refresh();
}

defineExpose({ refreshAssets });

const jobsClient = await getSettingsJobsassetsClient();

const emit = defineEmits<{
    countChanged: [count: number];
}>();

const editingAsset = ref<JobAsset>();

const { data, status, error, refresh } = useAuthedLazyAsyncData('userState', 'jobs-job-assets', ({ signal }) =>
    listAssets(signal),
);

async function listAssets(signal: AbortSignal): Promise<ListJobAssetsResponse> {
    const { response } = await jobsClient.listJobAssets({ pagination: { offset: 0 } }, { abort: signal });
    return response;
}

const assets = computed<JobAsset[]>(() => data.value?.assets ?? []);

watch(
    () => assets.value.length,
    (count) => emit('countChanged', count),
    { immediate: true },
);

function openEditor(asset: JobAsset): void {
    editingAsset.value = asset;
}

function closeEditor(): void {
    editingAsset.value = undefined;
}
</script>

<template>
    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.file', 2)])" :error="error" :retry="refresh" />
    <DataPendingBlock v-else-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.file', 2)])" />
    <JobAssetsFileList
        v-else
        :assets="assets"
        :can-update="canUpdate"
        :can-delete="canDelete"
        @edit="openEditor"
        @refresh="refresh"
    />

    <JobAssetsUpdate :asset="editingAsset" @close="closeEditor" @updated="refresh" />
</template>

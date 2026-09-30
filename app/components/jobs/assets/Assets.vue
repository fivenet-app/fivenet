<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import AssetsFileList from '~/components/jobs/assets/AssetsFileList.vue';
import AssetsUpdate from '~/components/jobs/assets/AssetsUpdate.vue';
import type { JobAsset } from '~~/gen/ts/resources/jobs/job_asset';
import type { ListJobAssetsResponse } from '~~/gen/ts/services/settings/job_assets';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';

const emit = defineEmits<{
    (e: 'countChanged', count: number): void;
}>();

const { can } = useAuth();
const overlay = useOverlay();

const canUpdate = can('settings.JobsAssetsService/UpdateJobAsset');
const canDelete = can('settings.JobsAssetsService/DeleteJobAsset');

const imageCacheBuster = ref(new Date().getTime());

const { data, status, error, refresh } = useAuthedLazyAsyncData('userState', 'jobs-job-assets', ({ signal }) =>
    listAssets(signal),
);

async function refreshAssets(): Promise<void> {
    await refresh();
    imageCacheBuster.value = new Date().getTime();
}

defineExpose({ refresh: refreshAssets, status });

const assetsUpdateModal = overlay.create(AssetsUpdate);

async function listAssets(signal: AbortSignal): Promise<ListJobAssetsResponse> {
    const jobsClient = await getSettingsJobsassetsClient();

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
    assetsUpdateModal.open({
        asset: asset,
        onUpdated: refreshAssets,
    });
}
</script>

<template>
    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.file', 2)])"
        :error="error"
        :retry="refreshAssets"
    />
    <DataPendingBlock v-else-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.file', 2)])" />
    <AssetsFileList
        v-else
        :assets="assets"
        :cache-buster="imageCacheBuster"
        :can-update="canUpdate"
        :can-delete="canDelete"
        @edit="openEditor"
        @refresh="refreshAssets"
    />
</template>

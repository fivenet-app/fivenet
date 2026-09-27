<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import type { JobAsset } from '~~/gen/ts/resources/jobs/job_asset';
import type { ListJobAssetsResponse } from '~~/gen/ts/services/settings/job_assets';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';
import { getJobAssetDisplayName } from '~/utils/files';

const props = defineProps<{
    fileLimitReached: boolean;
    onSelect: (url: string) => void;
}>();

const emit = defineEmits<{
    close: [open: boolean];
}>();

const jobsClient = await getSettingsJobsassetsClient();

const {
    data: jobAssetsData,
    status: jobAssetsStatus,
    error: jobAssetsError,
    refresh: refreshJobAssets,
} = useAuthedLazyAsyncData('userState', 'editor-job-assets', ({ signal }) => listJobAssets(signal));
const jobAssets = computed<JobAsset[]>(() => jobAssetsData.value?.assets ?? []);

async function listJobAssets(signal: AbortSignal): Promise<ListJobAssetsResponse> {
    const { response } = await jobsClient.listJobAssets({ pagination: { offset: 0 } }, { abort: signal });
    return response;
}

function selectJobAsset(asset: JobAsset): void {
    if (props.fileLimitReached) return;
    const path = asset.file?.filePath;
    if (!path) return;
    props.onSelect(`/api/filestore/${path}`);
    emit('close', false);
}
</script>

<template>
    <UModal :title="$t('common.job') + ' ' + $t('common.file', 2)" :ui="{ content: 'max-w-4xl' }">
        <template #body>
            <DataErrorBlock
                v-if="jobAssetsError"
                :title="$t('common.unable_to_load', [$t('common.file', 2)])"
                :error="jobAssetsError"
                :retry="refreshJobAssets"
            />
            <DataPendingBlock
                v-else-if="isRequestPending(jobAssetsStatus)"
                :message="$t('common.loading', [$t('common.file', 2)])"
            />
            <DataNoDataBlock
                v-else-if="!jobAssets.length"
                :message="$t('common.not_found', [$t('common.file', 2)])"
                :padded="false"
            />
            <div v-else class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
                <UButton
                    v-for="asset in jobAssets"
                    :key="asset.id"
                    class="group h-auto flex-col items-stretch overflow-hidden p-0 text-left"
                    variant="soft"
                    :disabled="props.fileLimitReached"
                    @click="selectJobAsset(asset)"
                >
                    <GenericImg
                        :src="`/api/filestore/${asset.file?.filePath}`"
                        :alt="getJobAssetDisplayName(asset)"
                        class="aspect-square w-full"
                        img-class="h-full w-full object-cover"
                        :rounded="false"
                    />
                    <span class="w-full truncate p-2 text-xs">{{ getJobAssetDisplayName(asset) }}</span>
                </UButton>
            </div>
        </template>

        <template #footer>
            <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
        </template>
    </UModal>
</template>

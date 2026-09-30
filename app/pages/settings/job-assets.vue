<script lang="ts" setup>
import Assets from '~/components/jobs/assets/Assets.vue';
import AssetsUpload from '~/components/jobs/assets/AssetsUpload.vue';
import Pagination from '~/components/partials/Pagination.vue';

const { jobAssets } = useAppConfig();
const { can } = useAuth();
const canCreate = can('settings.JobsAssetsService/CreateJobAsset');
const assetCount = ref(0);
const jobAssetsList = ref<InstanceType<typeof Assets>>();

async function refreshAssets(): Promise<void> {
    await jobAssetsList.value?.refresh();
}

const overlay = useOverlay();
const assetsUploadModal = overlay.create(AssetsUpload, {});

function openUpload(): void {
    assetsUploadModal.open({
        assetCount: assetCount.value,
        onUploaded: refreshAssets,
    });
}

useHead({
    title: 'pages.settings.job_assets.title',
});

definePageMeta({
    title: 'pages.settings.job_assets.title',
    requiresAuth: true,
    permission: 'settings.JobsAssetsService/ListJobAssets',
});
</script>

<template>
    <UDashboardPanel
        id="settings-job-assets"
        :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }"
    >
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.job_assets.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UDashboardToolbar>
                <template #left>
                    <UBadge :label="`${assetCount} / ${jobAssets.maxFiles}`" color="neutral" variant="subtle" />
                </template>

                <template #right>
                    <UButton
                        v-if="canCreate && assetCount < jobAssets.maxFiles"
                        icon="i-mdi-upload"
                        :label="$t('common.upload')"
                        @click="openUpload"
                    />
                </template>
            </UDashboardToolbar>

            <Assets ref="jobAssetsList" @count-changed="assetCount = $event" />
        </template>

        <template #footer>
            <Pagination :status="jobAssetsList?.status" :refresh="refreshAssets" hide-buttons hide-text />
        </template>
    </UDashboardPanel>
</template>

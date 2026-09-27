<script lang="ts" setup>
import JobAssets from '~/components/jobs/assets/Assets.vue';
import JobAssetsUpload from '~/components/jobs/assets/AssetsUpload.vue';

const { jobAssets } = useAppConfig();
const { can } = useAuth();
const canCreate = can('settings.JobsAssetsService/CreateJobAsset');
const assetCount = ref(0);
const uploadOpen = ref(false);
const jobAssetsList = ref<InstanceType<typeof JobAssets>>();

function handleUploaded(): void {
    jobAssetsList.value?.refreshAssets();
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
                        @click="uploadOpen = true"
                    />
                </template>
            </UDashboardToolbar>

            <JobAssets ref="jobAssetsList" @count-changed="assetCount = $event" />

            <JobAssetsUpload v-model:open="uploadOpen" :asset-count="assetCount" @uploaded="handleUploaded" />
        </template>
    </UDashboardPanel>
</template>

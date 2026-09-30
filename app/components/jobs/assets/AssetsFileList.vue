<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import type { JobAsset } from '~~/gen/ts/resources/jobs/job_asset';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';
import { copyToClipboardWrapper } from '~/utils/clipboard';

const props = defineProps<{
    assets: JobAsset[];
    cacheBuster: number;
    canUpdate: boolean;
    canDelete: boolean;
}>();

const emit = defineEmits<{
    (e: 'edit', asset: JobAsset): void;
    (e: 'refresh'): void;
}>();

const jobsClient = await getSettingsJobsassetsClient();
const { t } = useI18n();
const overlay = useOverlay();
const confirmModal = overlay.create(ConfirmModal);
const notifications = useNotificationsStore();

function getAssetURL(asset: JobAsset): string | undefined {
    const filePath = asset.file?.filePath;
    if (!filePath || typeof window === 'undefined') return undefined;

    return new URL(`/api/filestore/${filePath}`, window.location.origin).href;
}

async function copyImageLink(asset: JobAsset): Promise<void> {
    const url = getAssetURL(asset);
    if (!url) return;

    await copyToClipboardWrapper(url);
    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

async function deleteAsset(asset: JobAsset): Promise<void> {
    if (!props.canDelete) return;
    try {
        await jobsClient.deleteJobAsset({ id: asset.id });
        emit('refresh');
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

function confirmDelete(asset: JobAsset): void {
    confirmModal.open({
        title: t('common.delete'),
        description: t('pages.settings.job_assets.delete_description'),
        confirm: () => deleteAsset(asset),
    });
}
</script>

<template>
    <div class="grid gap-4 p-4 sm:p-6">
        <div v-if="assets.length" class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
            <div v-for="asset in assets" :key="asset.id" class="group relative overflow-hidden rounded-lg bg-elevated">
                <div class="relative">
                    <GenericImg
                        :src="`/api/filestore/${asset.file?.filePath}?date=${cacheBuster}`"
                        :alt="getJobAssetDisplayName(asset)"
                        class="aspect-square w-full"
                        img-class="h-full w-full object-cover"
                        :rounded="false"
                        enable-popup
                    />
                    <UTooltip :text="$t('common.copy_image_link')">
                        <UButton
                            class="absolute top-2 right-2 opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100 sm:focus:opacity-100"
                            size="xs"
                            color="neutral"
                            variant="soft"
                            icon="i-mdi-link-variant"
                            :aria-label="$t('common.copy_image_link')"
                            @click.stop="copyImageLink(asset)"
                        />
                    </UTooltip>
                </div>

                <div class="flex items-center gap-1 p-2 text-xs">
                    <span class="min-w-0 flex-1 truncate">{{ getJobAssetDisplayName(asset) }}</span>
                    <UButton v-if="canUpdate" size="xs" variant="ghost" icon="i-mdi-pencil" @click="emit('edit', asset)" />
                    <UButton
                        v-if="canDelete"
                        size="xs"
                        color="error"
                        variant="ghost"
                        icon="i-mdi-delete"
                        @click="confirmDelete(asset)"
                    />
                </div>
            </div>
        </div>
        <DataNoDataBlock v-else :message="$t('common.not_found', [$t('common.file', 2)])" :padded="false" />
    </div>
</template>

<script lang="ts" setup>
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';

const props = defineProps<{
    assetCount: number;
}>();

const emit = defineEmits<{
    uploaded: [];
}>();

const open = defineModel<boolean>('open', { default: false });
const { fileUpload, jobAssets } = useAppConfig();
const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);
const jobsClient = await getSettingsJobsassetsClient();
const { uploadImages } = useImageUpload();
const selectedFiles = ref<File[]>([]);

const { resizeAndUpload } = useFileUploader(() => jobsClient.uploadJobAsset({}), 'jobassets', 0);

async function doUploadSelected(): Promise<void> {
    if (selectedFiles.value.length === 0) return;

    try {
        const result = await uploadImages({
            files: selectedFiles.value,
            currentFileCount: props.assetCount,
            fileLimit: jobAssets.maxFiles,
            uploadOne: (file) => resizeAndUpload(file),
            invalidTypeNotification: {
                title: { key: 'common.error', parameters: {} },
                description: { key: 'common.allowed_file_types', parameters: {} },
            },
            fileLimitNotification: {
                title: { key: 'common.error', parameters: {} },
                description: { key: 'common.file', parameters: {} },
            },
            successNotification: {
                title: { key: 'notifications.action_successful.title', parameters: {} },
                description: { key: 'notifications.action_successful.content', parameters: {} },
            },
        });
        if (!result.ok) return;

        selectedFiles.value = [];
        open.value = false;
        emit('uploaded');
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

const { submit: uploadSelected, isSubmitting } = useSubmitGuard(doUploadSelected);

watch(open, (isOpen) => {
    if (!isOpen) selectedFiles.value = [];
});
</script>

<template>
    <UModal v-model:open="open" :title="$t('common.upload')" :close="!isSubmitting" :dismissible="!isSubmitting">
        <template #body>
            <NotSupportedTabletBlock v-if="nuiEnabled" />
            <div v-else class="grid gap-4">
                <UFileUpload
                    v-model="selectedFiles"
                    multiple
                    class="w-full"
                    :accept="fileUpload.types.images.join(',')"
                    :max-files="jobAssets.maxFiles - assetCount"
                    :placeholder="$t('common.image')"
                    :label="$t('common.file_upload_label')"
                    :description="$t('common.allowed_file_types')"
                />
                <PartialsContentGuidelinesAlert />
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    :disabled="isSubmitting"
                    :label="$t('common.cancel')"
                    @click="open = false"
                />

                <UButton
                    class="flex-1"
                    icon="i-mdi-upload"
                    :disabled="selectedFiles.length === 0 || isSubmitting"
                    :loading="isSubmitting"
                    :label="$t('common.upload')"
                    @click="uploadSelected"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

<script lang="ts" setup>
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    assetCount: number;
}>();

const emit = defineEmits<{
    uploaded: [];
}>();

const open = defineModel<boolean>('open', { default: false });
const { fileUpload, jobAssets } = useAppConfig();
const jobsClient = await getSettingsJobsassetsClient();
const notifications = useNotificationsStore();
const { uploadImages } = useImageUpload();
const selectedFiles = ref<File[]>([]);

const { resizeAndUpload } = useFileUploader(() => jobsClient.uploadJobAsset({}), 'jobassets', 0);

async function uploadSelected(): Promise<void> {
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
        });
        if (!result.ok) return;

        selectedFiles.value = [];
        open.value = false;
        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        emit('uploaded');
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}
</script>

<template>
    <UModal v-model:open="open" :title="$t('common.upload')">
        <template #body>
            <UFileUpload
                v-model="selectedFiles"
                multiple
                class="w-full"
                :accept="fileUpload.types.images.join(',')"
                :max-files="jobAssets.maxFiles - assetCount"
                :label="$t('common.file_upload_label')"
            />
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" :label="$t('common.cancel')" @click="open = false" />

                <UButton
                    class="flex-1"
                    icon="i-mdi-upload"
                    :disabled="selectedFiles.length === 0"
                    :label="$t('common.upload')"
                    @click="uploadSelected"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>

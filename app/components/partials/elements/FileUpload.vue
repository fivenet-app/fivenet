<script lang="ts" setup>
import type { ClientStreamingCall, RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { z } from 'zod';
import type { File as FileGRPC } from '~~/gen/ts/resources/file/file';
import type { UploadPacket, UploadResponse } from '~~/gen/ts/resources/file/filestore';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ConfirmModal from '../ConfirmModal.vue';
import NotSupportedTabletBlock from '../NotSupportedTabletBlock.vue';
import GenericImg from './GenericImg.vue';

const props = defineProps<{
    modelValue: FileGRPC | undefined;
    disabled?: boolean;

    uploadFn: (opts?: RpcOptions) => ClientStreamingCall<UploadPacket, UploadResponse>;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    deleteFn: (fileId: number, reason?: string) => UnaryCall<any, any>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: FileGRPC | undefined): void;
}>();

const modelValue = useVModel(props, 'modelValue', emit);

const modal = useModal();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

const appConfig = useAppConfig();

const _schema = z.object({
    fileUrl: z.custom<File>().array().min(1).max(1),
});

type Schema = z.output<typeof _schema>;

const state = reactive<Schema>({
    fileUrl: [] as File[],
});

const { resizeAndUpload } = useFileUploader(props.uploadFn, 'documents', 0);

async function uploadFile(files: File[]): Promise<void> {
    for (const f of files) {
        if (!f.type.startsWith('image/')) continue;

        const resp = await resizeAndUpload(f);

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        modelValue.value = resp.file;
        return;
    }
}

async function deleteFile(fileId: number | undefined): Promise<void> {
    if (fileId === undefined) return;

    try {
        await props.deleteFn(fileId);

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

function handleFileChanges(event: File[]) {
    state.fileUrl = event;
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <NotSupportedTabletBlock v-if="nuiEnabled" />
        <div v-else class="flex flex-col gap-1">
            <div class="flex flex-1 flex-row gap-1">
                <UInput
                    class="flex-1"
                    name="jobLogo"
                    type="file"
                    :accept="appConfig.fileUpload.types.images.join(',')"
                    block
                    :placeholder="$t('common.image')"
                    :disabled="disabled"
                    @change="handleFileChanges"
                />

                <UTooltip
                    v-if="!modelValue?.id || state.fileUrl.length > 0"
                    :ui="{ placement: 'top' }"
                    :text="$t('common.upload')"
                >
                    <UButton
                        icon="i-mdi-upload"
                        color="primary"
                        :disabled="state.fileUrl.length === 0 || disabled"
                        @click="uploadFile(state.fileUrl)"
                    />
                </UTooltip>
                <UTooltip v-else :ui="{ placement: 'top' }" :text="$t('common.delete')">
                    <UButton
                        icon="i-mdi-delete"
                        color="error"
                        :disabled="disabled"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteFile(modelValue?.id),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </div>

        <div class="flex w-full flex-col items-center justify-center gap-2">
            <GenericImg
                v-if="modelValue?.filePath"
                size="3xl"
                :src="`${modelValue.filePath}?date=${new Date().getTime()}`"
                :no-blur="true"
            />

            <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
        </div>
    </div>
</template>

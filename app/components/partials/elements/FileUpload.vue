<script lang="ts" setup>
import type { ClientStreamingCall, RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { z } from 'zod';
import type { File as FileGRPC } from '~~/gen/ts/resources/file/file';
import type { UploadFileRequest, UploadFileResponse } from '~~/gen/ts/resources/file/filestore';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ConfirmModal from '../ConfirmModal.vue';
import NotSupportedTabletBlock from '../NotSupportedTabletBlock.vue';
import GenericImg from './GenericImg.vue';

type FileLike = string | FileGRPC;
type MaybeFileLike = FileLike | undefined;

const props = defineProps<{
    modelValue: MaybeFileLike;
    disabled?: boolean;

    uploadFn: (opts?: RpcOptions) => ClientStreamingCall<UploadFileRequest, UploadFileResponse>;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    deleteFn: () => UnaryCall<any, any>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MaybeFileLike): void;
}>();

const modelValue = useVModel(props, 'modelValue', emit);

const overlay = useOverlay();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const _schema = z.object({
    fileUrl: z.custom<File>().array().min(1).max(1).default([]),
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
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (modelValue.value === undefined || typeof modelValue.value === 'string') {
            modelValue.value = resp.url;
        } else {
            modelValue.value = resp.file;
        }
        return;
    }
}

async function deleteFile(): Promise<void> {
    try {
        await props.deleteFn();

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
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

const filePath = computed(() => {
    if (typeof modelValue.value === 'string') {
        return modelValue.value;
    } else if (modelValue.value) {
        return modelValue.value.filePath;
    }
    return '';
});

const confirmModal = overlay.create(ConfirmModal);
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
                    @change="($event) => handleFileChanges($event)"
                />

                <UTooltip v-if="!filePath || state.fileUrl.length > 0" :text="$t('common.upload')">
                    <UButton
                        icon="i-mdi-upload"
                        color="primary"
                        :disabled="state.fileUrl.length === 0 || disabled"
                        @click="uploadFile(state.fileUrl)"
                    />
                </UTooltip>
                <UTooltip v-else :text="$t('common.delete')">
                    <UButton
                        icon="i-mdi-delete"
                        color="error"
                        :disabled="disabled"
                        @click="
                            confirmModal.open({
                                confirm: async () => deleteFile(),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </div>

        <div class="flex w-full flex-col items-center justify-center gap-2">
            <GenericImg v-if="filePath" size="3xl" :src="`${filePath}?date=${new Date().getTime()}`" :no-blur="true" />

            <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
        </div>
    </div>
</template>

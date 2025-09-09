<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
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

const schema = z.object({
    file: z.instanceof(File).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    file: undefined,
});

const { resizeAndUpload } = useFileUploader(props.uploadFn, 'documents', 0);

async function uploadFile(f: File): Promise<void> {
    if (!f.type.startsWith('image/')) return;

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

const filePath = computed(() => {
    if (typeof modelValue.value === 'string') {
        return modelValue.value;
    } else if (modelValue.value) {
        return modelValue.value.filePath;
    }
    return '';
});

async function onSubmit(event: FormSubmitEvent<Schema>) {
    await uploadFile(event.data.file!);
    state.file = undefined;
}

const confirmModal = overlay.create(ConfirmModal);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <div class="flex flex-col gap-2">
        <NotSupportedTabletBlock v-if="nuiEnabled" />
        <UForm v-else ref="formRef" :schema="schema" :state="state" class="flex flex-col gap-1" @submit="onSubmit">
            <div class="flex flex-1 flex-row gap-1">
                <UFileUpload
                    v-model="state.file"
                    class="flex-1"
                    name="jobLogo"
                    :accept="appConfig.fileUpload.types.images.join(',')"
                    block
                    :placeholder="$t('common.image')"
                    :disabled="disabled"
                />

                <UTooltip v-if="!filePath || !state.file" :text="$t('common.upload')">
                    <UButton
                        type="submit"
                        icon="i-mdi-upload"
                        color="primary"
                        :disabled="!state.file || disabled"
                        :loading="formRef?.loading"
                    />
                </UTooltip>

                <UTooltip v-else :text="$t('common.delete')">
                    <UButton
                        icon="i-mdi-delete"
                        color="error"
                        :disabled="disabled"
                        :loading="formRef?.loading"
                        @click="
                            confirmModal.open({
                                confirm: async () => deleteFile(),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </UForm>

        <div class="flex w-full flex-col items-center justify-center gap-2">
            <GenericImg
                v-if="filePath"
                size="3xl"
                :src="`${filePath}?date=${new Date().getTime()}`"
                :no-blur="true"
                img-class="h-30 w-auto"
            />

            <UAlert variant="subtle" icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
        </div>
    </div>
</template>

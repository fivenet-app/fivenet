<script lang="ts" setup>
import type { ClientStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import type { UploadNamespaces } from '~/composables/useFileUploader';
import { useSettingsStore } from '~/stores/settings';
import type { UploadFileResponse, UploadFileRequest } from '~~/gen/ts/resources/file/filestore';

type UploadFunction = (options?: RpcOptions) => ClientStreamingCall<UploadFileRequest, UploadFileResponse>;

const props = defineProps<{
    uploadFunc: UploadFunction;
    namespace: UploadNamespaces;
    parentId: number | string | Ref<number> | Ref<string>;
    onUploaded?: (response: UploadFileResponse) => void | Promise<void>;
}>();

const { fileUpload } = useAppConfig();
const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);
const { resizeAndUpload } = useFileUploader(props.uploadFunc, props.namespace, props.parentId);

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    file: z.file().mime(fileUpload.types.images).max(fileUpload.fileSizes.fileStore).min(1).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    file: undefined,
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            reason: value.reason,
            file: value.file
                ? {
                      name: value.file.name,
                      size: value.file.size,
                      type: value.file.type,
                  }
                : null,
        }),
});

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    if (!event.data.file) return;

    try {
        const response = await resizeAndUpload(event.data.file, event.data.reason);
        await props.onUploaded?.(response);
        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
});

const emit = defineEmits<{ (e: 'close', value: boolean): void }>();
const formRef = useTemplateRef('formRef');

async function closeDrawer(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;
    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal
        :title="$t('common.upload')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        :ui="{ body: 'mx-auto w-full max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">{{ $t('common.upload') }}</h3>
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeDrawer"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="grid gap-2" :schema="schema" :state="state" @submit="submit">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" :disabled="isSubmitting" :placeholder="$t('common.reason')" />
                </UFormField>

                <UFormField class="flex-1" name="file" :label="$t('common.file')">
                    <NotSupportedTabletBlock v-if="nuiEnabled" />
                    <div v-else class="grid gap-2">
                        <UFileUpload
                            v-model="state.file"
                            class="mx-auto max-w-md flex-1"
                            position="inside"
                            name="file"
                            :disabled="isSubmitting"
                            :accept="fileUpload.types.images.join(',')"
                            :placeholder="$t('common.image')"
                            :label="$t('common.file_upload_label')"
                            :description="$t('common.allowed_file_types')"
                            :ui="{
                                base: 'min-h-50',
                                files: 'min-h-50 w-full',
                                file: 'flex min-h-50 items-center justify-center',
                                fileLeadingAvatar: 'size-50 max-h-50 rounded-lg',
                            }"
                        />

                        <PartialsContentGuidelinesAlert />
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeDrawer"
                />
                <UButton
                    class="flex-1"
                    :disabled="!state.file || !canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.upload')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </ResponsiveDrawerModal>
</template>

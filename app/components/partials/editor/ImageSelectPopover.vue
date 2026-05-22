<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { Editor } from '@tiptap/core';
import { z } from 'zod';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        editor: Editor;
        files?: FileGrpc[];
        fileLimit?: number;
        disabled?: boolean;
        uploadHandler?: (file: File[]) => Promise<boolean>;
        openFileList?: () => Promise<void>;
    }>(),
    {
        files: () => [],
        fileLimit: 5,
        disabled: false,
        uploadHandler: undefined,
        openFileList: undefined,
    },
);

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'openFileList'): void;
}>();

const { fileUpload } = useAppConfig();

const notifications = useNotificationsStore();

const schema = z.object({
    url: z.url().max(512),
});

type Schema = z.output<typeof schema>;

const imageState = reactive<Schema>({
    url: '',
});

// Try to download image from remote url
async function setViaURL(urlOrBlob: string | File): Promise<void> {
    canSubmit.value = false;

    // Check if file limit is reached
    if (props.files.length >= props.fileLimit) {
        logger.warn('File limit reached, cannot upload more files');
        notifications.add({
            title: { key: 'components.partials.tiptap_editor.notifications.file_limit_reached.title', parameters: {} },
            description: {
                key: 'components.partials.tiptap_editor.notifications.file_limit_reached.content',
                parameters: {},
            },
            type: NotificationType.ERROR,
        });
        return;
    }

    // Use image proxy for external URLs
    if (typeof urlOrBlob === 'string') {
        let dataUrl: string | undefined = undefined;
        // If Image Proxy is enabled use it to load the image
        if (urlOrBlob.startsWith('http')) {
            const url = new URL(urlOrBlob);
            // Check if image is already served by our host and one of the paths
            const isSameHost = url.host === window.location.host;
            const isServedPath = safeImagePaths.some((path) => url.pathname.startsWith(path + '/'));
            if (isSameHost && isServedPath) {
                url.pathname = url.pathname.replace(/(?<!:)\/\//, '/');
                dataUrl = urlOrBlob;
            } else {
                dataUrl = `/api/image_proxy/${urlOrBlob}`;
            }
        } else {
            dataUrl = urlOrBlob;
        }

        return setImage(dataUrl);
    } else if (props.uploadHandler) {
        try {
            const response = await props.uploadHandler([urlOrBlob]);
            if (!response) {
                console.warn('Editor - Image upload response: Upload handler returned false');
                return;
            }

            notifications.add({
                title: { key: 'notifications.editor.file_upload.success.title', parameters: {} },
                description: { key: 'notifications.editor.file_upload.success.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } catch (e) {
            console.error('Editor - Image upload failed', e);

            notifications.add({
                title: { key: 'notifications.editor.file_upload.failed.title', parameters: {} },
                description: {
                    key: 'notifications.editor.file_upload.failed.content',
                    parameters: { error: (e as Error)?.message?.toString() ?? 'N/A' },
                },
                type: NotificationType.ERROR,
            });
        }
    } else {
        console.warn('Editor - No upload handler provided for image upload');
    }
}

function setImage(url: string | undefined): void {
    if (!url) return;

    props.editor
        .chain()
        .setEnhancedImage({
            src: url,
            alt: url,
        })
        .run();

    emit('close', false);
}

async function onFileHandler(file: File | null | undefined): Promise<void> {
    if (!file) {
        canSubmit.value = true;
        return;
    }

    await setViaURL(file);

    canSubmit.value = true;
    emit('close', false);
}

const fileLimitReached = computed(() => props.files.length >= props.fileLimit);

const open = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await setViaURL(event.data.url).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));

    open.value = false;
    imageState.url = '';
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UPopover v-model:open="open">
        <UTooltip :text="$t('components.partials.tiptap_editor.image')">
            <UButton icon="i-mdi-images" :color="fileLimitReached ? 'error' : 'neutral'" variant="ghost" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <UAlert
                    v-if="fileLimitReached"
                    class="w-96"
                    color="error"
                    icon="i-mdi-close-circle"
                    :title="$t('components.partials.tiptap_editor.file_limit_reached.title')"
                    :description="$t('components.partials.tiptap_editor.file_limit_reached.description')"
                />

                <UFieldGroup class="w-full">
                    <UButton
                        class="flex-1"
                        color="neutral"
                        block
                        icon="i-mdi-images"
                        :label="$t('components.partials.tiptap_editor.file_list')"
                        :disabled="disabled || !canSubmit"
                        :loading="disabled || !canSubmit"
                        @click="$emit('openFileList')"
                    />
                </UFieldGroup>

                <USeparator class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <UForm ref="formRef" :schema="schema" :state="imageState" @submit="onSubmitThrottle">
                    <UFormField name="url" :label="$t('components.partials.tiptap_editor.url')">
                        <UInput
                            v-model="imageState.url"
                            class="w-full"
                            type="text"
                            name="url"
                            :disabled="disabled || !canSubmit || fileLimitReached"
                        />
                    </UFormField>

                    <UFormField class="mt-2 w-full">
                        <UButton
                            class="flex-1"
                            type="submit"
                            icon="i-mdi-image"
                            :disabled="disabled || !canSubmit || !imageState.url || fileLimitReached"
                            :loading="disabled || !canSubmit"
                            :label="$t('components.partials.tiptap_editor.insert')"
                            block
                            @click="formRef?.submit()"
                        />
                    </UFormField>
                </UForm>

                <USeparator class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <UFileUpload
                    name="file"
                    block
                    :disabled="disabled || !canSubmit || fileLimitReached"
                    :accept="fileUpload.types.images.join(',')"
                    :placeholder="$t('common.image')"
                    :label="$t('common.file_upload_label')"
                    :description="$t('common.allowed_file_types')"
                    @update:model-value="($event) => onFileHandler($event)"
                />
            </div>
        </template>
    </UPopover>
</template>

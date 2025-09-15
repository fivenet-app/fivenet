<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { Editor } from '@tiptap/vue-3';
import { z } from 'zod';
import { safeImagePaths } from '~/types/editor';
import { remoteImageURLToBase64Data } from './helpers';

const props = withDefaults(
    defineProps<{
        editor: Editor;
        fileLimit?: number;
        disabled?: boolean;
        uploadHandler?: (file: File[]) => Promise<void>;
        openFileList?: () => Promise<void>;
    }>(),
    {
        fileLimit: 10,
        disabled: false,
        uploadHandler: undefined,
        openFileList: undefined,
    },
);

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'openFileList'): void;
}>();

const { featureGates, fileUpload } = useAppConfig();

const schema = z.object({
    url: z.url().max(512),
});

type Schema = z.output<typeof schema>;

const imageState = reactive({
    url: '',
});

// Try to download image from remote url
async function setViaURL(urlOrBlob: string | File): Promise<void> {
    canSubmit.value = false;

    if (typeof urlOrBlob === 'string') {
        let dataUrl: string | undefined = undefined;
        // If Image Proxy is enabled use it to load the image
        if (featureGates.imageProxy && urlOrBlob.startsWith('http')) {
            const url = new URL(urlOrBlob);
            // Check if image is already served by our host and one of the paths
            const isSameHost = url.host === window.location.host;
            const isServedPath = safeImagePaths.some((path) => url.pathname.startsWith(path));
            if (isSameHost && isServedPath) {
                url.pathname = url.pathname.replace(/(?<!:)\/\//, '/');
                dataUrl = urlOrBlob;
            } else if (props.uploadHandler) {
                dataUrl = `/api/image_proxy/${urlOrBlob}`;
            } else {
                dataUrl = await remoteImageURLToBase64Data(`/api/image_proxy/${url.toString()}`);
            }
        } else {
            dataUrl = urlOrBlob;
        }

        return setImage(dataUrl);
    } else if (props.uploadHandler) {
        await props.uploadHandler([urlOrBlob]);
    } else {
        setImage(await blobToBase64(urlOrBlob));
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
        <UTooltip :text="$t('components.partials.TiptapEditor.image')">
            <UButton icon="i-mdi-image" color="neutral" variant="ghost" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="p-4">
                <UButtonGroup class="w-full">
                    <UButton
                        class="flex-1"
                        color="neutral"
                        block
                        icon="i-mdi-images"
                        :label="$t('components.partials.TiptapEditor.file_list')"
                        :disabled="disabled || !canSubmit"
                        :loading="disabled || !canSubmit"
                        @click="$emit('openFileList')"
                    />
                </UButtonGroup>

                <USeparator class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <UForm ref="formRef" :schema="schema" :state="imageState" @submit="onSubmitThrottle">
                    <UFormField name="url" :label="$t('common.url')">
                        <UInput v-model="imageState.url" type="text" name="url" class="w-full" :disabled="disabled" />
                    </UFormField>

                    <UFormField class="mt-2 w-full">
                        <UButton
                            class="flex-1"
                            type="submit"
                            icon="i-mdi-image"
                            :disabled="disabled || !canSubmit || !imageState.url"
                            :loading="disabled || !canSubmit"
                            :label="$t('common.insert')"
                            block
                            @click="formRef?.submit()"
                        />
                    </UFormField>
                </UForm>

                <USeparator class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <UFileUpload
                    name="file"
                    block
                    :disabled="disabled || !canSubmit"
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

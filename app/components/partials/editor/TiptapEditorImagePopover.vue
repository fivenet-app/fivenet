<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import type { Editor } from '@tiptap/core';
import { z } from 'zod';
import { remoteImageURLToBase64Data } from './helpers';

const { isOpen } = useModal();

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

defineEmits<{
    (e: 'openFileList'): void;
}>();

const { featureGates, fileUpload } = useAppConfig();

const schema = z.object({
    url: z.string().url(),
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
            if (props.uploadHandler) {
                dataUrl = `/api/image_proxy/${urlOrBlob}`;
            } else {
                const url = new URL(urlOrBlob);
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
    if (!url) {
        return;
    }

    props.editor
        .chain()
        .setEnhancedImage({
            src: url,
            alt: url,
        })
        .run();

    isOpen.value = false;
}

async function onFilesHandler(files: FileList | File[] | null): Promise<void> {
    if (!files || !files[0]) {
        canSubmit.value = true;
        return;
    }

    await setViaURL(files[0]);

    canSubmit.value = true;
    isOpen.value = false;
}

const dropZoneRef = useTemplateRef('dropZoneRef');

const { chooseFiles } = useFileSelection({
    dropzone: dropZoneRef,
    onFiles: onFilesHandler,
    // Specify the types of data to be received.
    allowedDataTypes: fileUpload.types.images,
    multiple: false,
});

const open = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await setViaURL(event.data.url).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UPopover v-model:open="open">
        <UTooltip :text="$t('components.partials.TiptapEditor.image')" :popper="{ placement: 'top' }">
            <UButton icon="i-mdi-image" color="white" variant="ghost" />
        </UTooltip>

        <template #panel>
            <div class="p-4">
                <UButtonGroup class="w-full">
                    <UButton
                        class="flex-1"
                        color="black"
                        block
                        icon="i-mdi-images"
                        :label="$t('components.partials.TiptapEditor.file_list')"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="$emit('openFileList')"
                    />
                </UButtonGroup>

                <UDivider class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <UForm :schema="schema" :state="imageState" @submit="onSubmitThrottle">
                    <UFormGroup :label="$t('common.url')">
                        <UInput v-model="imageState.url" type="text" />
                    </UFormGroup>

                    <UButtonGroup class="mt-2 w-full">
                        <UButton
                            class="flex-1"
                            type="submit"
                            icon="i-mdi-image"
                            :label="$t('common.insert')"
                            :disabled="!canSubmit || !imageState.url"
                            :loading="!canSubmit"
                        />
                    </UButtonGroup>
                </UForm>

                <UDivider class="my-2" :label="$t('common.or')" orientation="horizontal" />

                <ULink class="w-full" @click="chooseFiles">
                    <div ref="dropZoneRef" class="flex w-full items-center justify-center">
                        <label
                            class="flex h-48 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-100 hover:bg-gray-200 dark:border-gray-600 dark:bg-gray-800 dark:hover:border-gray-600 dark:hover:bg-gray-700"
                            for="dropzone-file"
                        >
                            <div class="flex flex-col items-center justify-center pb-4 pt-3">
                                <UIcon
                                    class="size-14"
                                    :class="!canSubmit && 'animate-spin'"
                                    :name="canSubmit ? 'i-mdi-file-upload-outline' : 'i-mdi-loading'"
                                />

                                <p class="mb-2 px-2 text-base text-gray-500 dark:text-gray-400">
                                    <span class="font-semibold">{{ $t('common.file_click_to_upload') }}</span>
                                    {{ $t('common.file_drag_n_drop') }}
                                </p>
                                <p class="text-sm text-gray-500 dark:text-gray-400">{{ $t('common.allowed_file_types') }}</p>
                            </div>
                        </label>
                    </div>
                </ULink>
            </div>
        </template>
    </UPopover>
</template>

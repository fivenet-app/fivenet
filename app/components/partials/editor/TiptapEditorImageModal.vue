<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import type { Editor } from '@tiptap/core';
import { z } from 'zod';
import { remoteImageURLToBase64Data } from './helpers';

const { isOpen } = useModal();

const props = defineProps<{
    editor: Editor;
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
async function setViaURL(urlOrBlob: string | Blob): Promise<void> {
    canSubmit.value = false;
    let dataUrl: string | undefined = undefined;
    if (typeof urlOrBlob === 'string') {
        // If Image Proxy is enabled use it to load the image
        if (featureGates.imageProxy) {
            const url = new URL(urlOrBlob);
            dataUrl = await remoteImageURLToBase64Data('/api/image_proxy/' + url.toString());
        } else {
            dataUrl = urlOrBlob;
        }
    } else {
        dataUrl = await blobToBase64(urlOrBlob);
    }

    setImage(dataUrl);
}

function setImage(url: string | undefined): void {
    if (!url) {
        return;
    }

    props.editor
        .chain()
        .setImage({
            src: url,
        })
        .run();

    isOpen.value = false;
}

async function onFilesHandler(files: FileList | File[] | null): Promise<void> {
    if (!files || !files[0]) {
        return;
    }

    return setViaURL(files[0]);
}

const dropZoneRef = useTemplateRef('dropZoneRef');

const { chooseFiles } = useFileSelection({
    dropzone: dropZoneRef,
    onFiles: onFilesHandler,
    // Specify the types of data to be received.
    allowedDataTypes: fileUpload.types.images,
    multiple: false,
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await setViaURL(event.data.url).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :prevent-close="!canSubmit">
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.image') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :schema="schema" :state="imageState" @submit="onSubmitThrottle">
                    <UFormGroup :label="$t('common.url')">
                        <UInput v-model="imageState.url" type="text" />
                    </UFormGroup>

                    <UButtonGroup class="mt-2 w-full">
                        <UButton
                            type="submit"
                            icon="i-mdi-image"
                            class="flex-1"
                            :label="$t('common.insert')"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                        />
                    </UButtonGroup>
                </UForm>

                <UDivider :label="$t('common.or')" orientation="horizontal" class="mb-2 mt-2" />

                <ULink class="w-full" @click="chooseFiles">
                    <div ref="dropZoneRef" class="flex w-full items-center justify-center">
                        <label
                            for="dropzone-file"
                            class="flex h-64 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-100 hover:bg-gray-200 dark:border-gray-600 dark:bg-gray-800 dark:hover:border-gray-600 dark:hover:bg-gray-700"
                        >
                            <div class="flex flex-col items-center justify-center pb-6 pt-5">
                                <UIcon :name="canSubmit ? 'i-mdi-file-upload-outline' : 'i-mdi-loading'" class="size-14" />

                                <p class="mb-2 text-base text-gray-500 dark:text-gray-400">
                                    <span class="font-semibold">{{ $t('common.file_click_to_upload') }}</span>
                                    {{ $t('common.file_drag_n_drop') }}
                                </p>
                                <p class="text-sm text-gray-500 dark:text-gray-400">PNG, JPG (MAX. 2000x1500px)</p>
                            </div>
                        </label>
                    </div>
                </ULink>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton block class="flex-1" color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

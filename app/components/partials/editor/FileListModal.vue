<script lang="ts" setup>
import type { Editor } from '@tiptap/vue-3';
import type { File } from '~~/gen/ts/resources/file/file';
import DataNoDataBlock from '../data/DataNoDataBlock.vue';
import GenericImg from '../elements/GenericImg.vue';

defineProps<{
    editor: Editor;
    files: File[];
}>();

const { isOpen } = useModal();
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
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
                        {{ $t('components.partials.TiptapEditor.file_list') }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div class="mx-auto flex w-full max-w-screen-xl flex-1 flex-col">
                <DataNoDataBlock v-if="files.length === 0" :message="$t('components.partials.TiptapEditor.file_list_empty')" />

                <UPageGrid v-else class="flex-1">
                    <UPageCard
                        v-for="file in files"
                        :key="file.id"
                        :title="file.filePath"
                        icon="i-mdi-file-document"
                        :ui="{ title: '!line-clamp-3 !whitespace-normal' }"
                    >
                        <template #icon>
                            <div class="flex-1">
                                <GenericImg
                                    v-if="file.contentType.startsWith('image/')"
                                    :src="file.filePath"
                                    :alt="file.filePath"
                                    size="3xl"
                                />
                                <UIcon
                                    v-else
                                    class="h-20 w-20 text-3xl"
                                    :name="file.contentType.startsWith('video/') ? 'i-mdi-video' : 'i-mdi-file-document'"
                                />
                            </div>

                            <div>
                                <UTooltip :text="$t('common.insert_file')">
                                    <UButton
                                        :icon="
                                            file.contentType.startsWith('image/') ? 'i-mdi-file-image-plus' : 'i-mdi-file-plus'
                                        "
                                        variant="link"
                                        size="xl"
                                        @click="
                                            editor
                                                .chain()
                                                .focus()
                                                .setEnhancedImage({
                                                    src: `/api/filestore/${file.filePath}`,
                                                    alt: file.filePath,
                                                    fileId: file.id,
                                                })
                                                .run();
                                            isOpen = false;
                                        "
                                    />
                                </UTooltip>
                            </div>
                        </template>

                        <template #description>
                            <ul>
                                <li>{{ file.contentType }}</li>
                                <li>{{ formatBytes(file.byteSize) }}</li>
                            </ul>
                        </template>
                    </UPageCard>
                </UPageGrid>

                <UAlert
                    class="mt-4"
                    icon="i-mdi-information-outline"
                    :description="$t('components.partials.TiptapEditor.file_list_hint')"
                />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" block color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>

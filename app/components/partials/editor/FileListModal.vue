<script lang="ts" setup>
import type { Editor } from '@tiptap/core';
import type { File } from '~~/gen/ts/resources/file/file';
import DataNoDataBlock from '../data/DataNoDataBlock.vue';
import GenericImg from '../elements/GenericImg.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    editor: Editor;
    files: File[];
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

function deleteFile(id: number): void {
    props.editor.chain().focus().removeEnhancedImageByFileId(id).run();

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}
</script>

<template>
    <UModal :title="$t('components.partials.tiptap_editor.file_list')" fullscreen>
        <template #body>
            <div class="mx-auto flex w-full flex-1 flex-col">
                <DataNoDataBlock v-if="files.length === 0" :message="$t('components.partials.tiptap_editor.file_list_empty')" />

                <UPageGrid v-else class="flex-1 lg:grid-cols-2">
                    <UPageCard
                        v-for="file in files"
                        :key="file.id"
                        :title="file.filePath"
                        icon="i-mdi-file-document"
                        :ui="{ title: 'line-clamp-3! whitespace-normal!', leading: 'size-full' }"
                    >
                        <template #leading>
                            <div class="flex-1">
                                <GenericImg
                                    v-if="file.contentType.startsWith('image/')"
                                    :src="file.filePath"
                                    :alt="file.filePath"
                                    size="3xl"
                                    enable-popup
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
                                            $emit('close', false);
                                        "
                                    />
                                </UTooltip>

                                <UTooltip :text="$t('common.remove_file')">
                                    <UButton
                                        color="error"
                                        icon="i-mdi-file-image-remove"
                                        size="xl"
                                        variant="link"
                                        @click="deleteFile(file.id)"
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
                    :description="$t('components.partials.tiptap_editor.file_list_hint')"
                />
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>

<script lang="ts" setup>
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import type { File as FilestoreFile } from '~~/gen/ts/resources/file/file';
import { copyToClipboardWrapper } from '~/utils/clipboard';

const props = defineProps<{ files: { file?: FilestoreFile; kind: string }[]; canDelete?: boolean }>();
const emit = defineEmits<{ (e: 'delete', file: FilestoreFile): void }>();
const { t } = useI18n();
const { notifyLinkCopied } = useClipboardNotifications();

function isImage(file: FilestoreFile | undefined): boolean {
    return !!file?.contentType?.startsWith('image/');
}

function kindLabel(kind: string): string {
    switch (kind) {
        case 'mugshot':
            return t('common.mugshot');
        case 'attachment':
            return t('common.attachment');
        default:
            return t('common.file');
    }
}

async function copyLink(file: FilestoreFile): Promise<void> {
    if (typeof window === 'undefined') return;
    await copyToClipboardWrapper(new URL(`/api/filestore/${file.filePath}`, window.location.origin).href);
    notifyLinkCopied();
}
</script>

<template>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        <div v-for="entry in props.files" :key="entry.file?.id" class="group relative overflow-hidden rounded-lg bg-elevated">
            <GenericImg
                v-if="isImage(entry.file)"
                :src="`/api/filestore/${entry.file?.filePath}`"
                :alt="kindLabel(entry.kind)"
                class="aspect-square w-full"
                img-class="h-full w-full object-cover"
                :rounded="false"
                enable-popup
            />
            <div v-else class="flex aspect-square flex-col items-center justify-center gap-2 p-3 text-center">
                <UIcon name="i-mdi-file-outline" class="size-10" />
                <span class="line-clamp-2 text-xs">{{ kindLabel(entry.kind) }}</span>
            </div>
            <div class="flex items-center gap-1 p-2 text-xs">
                <span class="min-w-0 flex-1 truncate">{{ kindLabel(entry.kind) }}</span>
                <UTooltip :text="$t('common.copy_file_link')">
                    <UButton
                        size="xs"
                        variant="ghost"
                        icon="i-mdi-link-variant"
                        :aria-label="$t('common.copy_file_link')"
                        @click.stop="entry.file && copyLink(entry.file)"
                    />
                </UTooltip>
                <UButton
                    v-if="props.canDelete && entry.kind !== 'mugshot'"
                    size="xs"
                    color="error"
                    variant="ghost"
                    icon="i-mdi-delete"
                    :aria-label="$t('common.delete')"
                    @click.stop="entry.file && emit('delete', entry.file)"
                />
            </div>
        </div>
    </div>
</template>

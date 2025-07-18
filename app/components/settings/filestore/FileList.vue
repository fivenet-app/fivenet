<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import FileUploadModal from '~/components/settings/filestore/FileUploadModal.vue';
import { useSettingsStore } from '~/stores/settings';
import type { File } from '~~/gen/ts/resources/file/file';
import type { DeleteFileResponse } from '~~/gen/ts/resources/file/filestore';
import type { ListFilesResponse } from '~~/gen/ts/services/filestore/filestore';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const prefix = ref('');

const page = useRouteQuery('page', '1', { transform: Number });

const {
    data: files,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`files-${page.value}-${prefix.value}`, () => listFiles(prefix.value));

async function listFiles(prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = $grpc.filestore.filestore.listFiles({
            pagination: {
                offset: calculateOffset(page.value, files.value?.pagination),
            },
            path: prefix,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteFile(path: string): Promise<DeleteFileResponse> {
    try {
        const { response } = $grpc.filestore.filestore.deleteFileByPath({
            path: path,
        });

        const idx = files.value?.files.findIndex((f) => f.filePath === path);
        if (idx !== undefined && idx > -1 && files.value !== null) {
            files.value?.files.splice(idx, 1);
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function addUploadedFile(file: File): void {
    const idx = files.value?.files.findIndex((f) => f.filePath === file.filePath);
    if (idx === undefined) {
        return;
    }

    if (idx > -1) {
        files.value?.files.unshift(file);
    } else {
        files.value?.files.splice(idx, 1, file);
    }
}

const modal = useModal();

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'preview',
        label: t('common.preview'),
    },
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'fileSize',
        label: t('common.file_size'),
    },
    {
        key: 'updatedAt',
        label: t('common.updated_at'),
    },
    {
        key: 'contentType',
        label: t('common.type'),
    },
];

const previewTypes = ['jpg', 'jpeg', 'png', 'webp'];
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardNavbar :title="$t('pages.settings.settings.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent>
            <StreamerModeAlert />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <UDashboardNavbar :title="$t('pages.settings.filestore.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />

                <UButton
                    trailing-icon="i-mdi-upload"
                    @click="
                        modal.open(FileUploadModal, {
                            onUploaded: addUploadedFile,
                        })
                    "
                >
                    {{ $t('common.upload') }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.file', 2)])"
            :error="error"
            :retry="refresh"
        />

        <UTable
            v-else
            class="flex-1"
            :loading="loading"
            :columns="columns"
            :rows="files?.files"
            :empty-state="{ icon: 'i-mdi-file-multiple', label: $t('common.not_found', [$t('common.file', 2)]) }"
        >
            <template #actions-data="{ row: file }">
                <UTooltip :text="$t('common.show')">
                    <UButton
                        variant="link"
                        icon="i-mdi-link-variant"
                        external
                        target="_blank"
                        :to="`/api/filestore/${file.filePath}`"
                    />
                </UTooltip>

                <UTooltip :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteFile(file.filePath),
                            })
                        "
                    />
                </UTooltip>
            </template>

            <template #name-data="{ row: file }">
                <span class="text-gray-900 dark:text-white">
                    {{ file.filePath }}
                </span>
            </template>

            <template #preview-data="{ row: file }">
                <UIcon
                    v-if="!previewTypes.some((ext) => file.filePath.endsWith(ext))"
                    class="size-8"
                    name="i-mdi-file-outline"
                />
                <NuxtImg v-else class="max-h-24 max-w-32" :src="`/api/filestore/${file.filePath}`" loading="lazy" />
            </template>

            <template #fileSize-data="{ row: file }">
                {{ formatBytes(file.size) }}
            </template>

            <template #updatedAt-data="{ row: file }">
                <GenericTime :value="toDate(file.lastModified)" />
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="files?.pagination" :loading="loading" :refresh="refresh" />
    </template>
</template>

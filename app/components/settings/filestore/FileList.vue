<script lang="ts" setup>
import { NuxtImg, UIcon } from '#components';
import type { TableColumn } from '@nuxt/ui';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import FileUploadModal from '~/components/settings/filestore/FileUploadModal.vue';
import { useSettingsStore } from '~/stores/settings';
import { getFilestoreFilestoreClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import type { DeleteFileResponse } from '~~/gen/ts/resources/file/filestore';
import type { ListFilesResponse } from '~~/gen/ts/services/filestore/filestore';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const filestoreFilestoreClient = await getFilestoreFilestoreClient();

const prefix = ref('');

const page = useRouteQuery('page', '1', { transform: Number });

const {
    data: files,
    status,
    refresh,
    error,
} = useLazyAsyncData(`files-${page.value}-${prefix.value}`, () => listFiles(prefix.value));

async function listFiles(prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = filestoreFilestoreClient.listFiles({
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
        const { response } = filestoreFilestoreClient.deleteFileByPath({
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

const overlay = useOverlay();

const fileUploadModal = overlay.create(FileUploadModal, {
    props: {
        onUploaded: addUploadedFile,
    },
});
const confirmModal = overlay.create(ConfirmModal);

const previewTypes = ['jpg', 'jpeg', 'png', 'webp'];

const columns = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.filePath),
            },
            {
                accessorKey: 'preview',
                header: t('common.preview'),
                cell: ({ row }) =>
                    !previewTypes.some((ext) => row.original.filePath.endsWith(ext))
                        ? h(UIcon, { class: 'size-8', name: 'i-mdi-file-outline' })
                        : h(NuxtImg, {
                              class: 'max-h-24 max-w-32',
                              src: `/api/filestore/${row.original.filePath}`,
                              loading: 'lazy',
                          }),
            },
            {
                accessorKey: 'fileSize',
                header: t('common.file_size'),
                cell: ({ row }) => formatBytes(row.original.byteSize),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created_at'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
        ] as TableColumn<File>[],
);
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

                <UButton trailing-icon="i-mdi-upload" @click="fileUploadModal.open()">
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
            :loading="isRequestPending(status)"
            :columns="columns"
            :data="files?.files"
            :pagination-options="{ manualPagination: true }"
            :sorting-options="{ manualSorting: true }"
            :empty="$t('common.not_found', [$t('common.file', 2)])"
            sticky
        >
            <template #actions-cell="{ row: file }">
                <UTooltip :text="$t('common.show')">
                    <UButton
                        variant="link"
                        icon="i-mdi-link-variant"
                        external
                        target="_blank"
                        :to="`/api/filestore/${file.original.filePath}`"
                    />
                </UTooltip>

                <UTooltip :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            confirmModal.open({
                                confirm: async () => deleteFile(file.original.filePath),
                            })
                        "
                    />
                </UTooltip>
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="files?.pagination" :status="status" :refresh="refresh" />
    </template>
</template>

<script lang="ts" setup>
import { UButton, UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
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

const schema = z.object({
    prefix: z.coerce
        .string()
        .max(64)
        .optional()
        .default('')
        .transform((val) => val.slice(0, 255)),
    page: pageNumberSchema,
});

const query = useSearchForm('settings-filelist', schema);

const {
    data: files,
    status,
    refresh,
    error,
} = useLazyAsyncData(`files-${query.page}-${query.prefix}`, () => listFiles(query.page, query.prefix));

async function listFiles(page: number, prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = filestoreFilestoreClient.listFiles({
            pagination: {
                offset: calculateOffset(page, files.value?.pagination),
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
    if (idx === undefined) return;

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
                id: 'actions',
                cell: ({ row }) =>
                    h('div', { class: 'flex items-center gap-1' }, [
                        row.original.isDir
                            ? h(
                                  UTooltip,
                                  { text: t('common.directory') },
                                  {
                                      default: () =>
                                          h(UButton, {
                                              variant: 'link',
                                              icon: 'i-mdi-subdirectory-arrow-right',
                                              onClick: () => {
                                                  query.prefix += row.original.filePath.replace(query.prefix, '') + '/';
                                              },
                                          }),
                                  },
                              )
                            : null,
                        h(
                            UTooltip,
                            { text: t('common.show') },
                            {
                                default: () =>
                                    h(
                                        UButton,
                                        {
                                            variant: 'link',
                                            icon: 'i-mdi-link-variant',
                                            external: true,
                                            target: '_blank',
                                            to: `/api/filestore/${row.original.filePath}`,
                                        },
                                        {},
                                    ),
                            },
                        ),
                        h(
                            UTooltip,
                            { text: t('common.delete') },
                            {
                                default: () =>
                                    h(
                                        UButton,
                                        {
                                            variant: 'link',
                                            icon: 'i-mdi-delete',
                                            color: 'error',
                                            onClick: () => {
                                                confirmModal.open({
                                                    confirm: async () => deleteFile(row.original.filePath),
                                                });
                                            },
                                        },
                                        {},
                                    ),
                            },
                        ),
                    ]),
            },
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.filePath),
            },
            {
                accessorKey: 'preview',
                header: t('common.preview'),
                cell: ({ row }) =>
                    row.original.isDir
                        ? h(UIcon, { class: 'size-8', name: 'i-mdi-folder' })
                        : !previewTypes.some((ext) => row.original.filePath.endsWith(ext))
                          ? h(UIcon, { class: 'size-8', name: 'i-mdi-file-outline' })
                          : h(
                                'div',
                                { class: 'flex justify-center items-center' },
                                h(GenericImg, {
                                    src: `/api/filestore/${row.original.filePath.replace(/^\//, '')}`,
                                    alt: row.original.filePath,
                                    loading: 'lazy',
                                    rounded: false,
                                    enablePopup: true,
                                    imgClass: 'max-h-24 max-w-32',
                                }),
                            ),
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

function goBackDirectory(): void {
    if (!query.prefix) return;

    const parts = query.prefix.split('/').filter((p) => p.length > 0);
    parts.pop();
    query.prefix = parts.length > 0 ? parts.join('/') + '/' : '';
}

watchDebounced(query, async () => (await formRef.value?.validate()) && refresh(), { debounce: 200, maxWait: 1250 });

const formRef = useTemplateRef('formRef');

const inputRef = useTemplateRef('inputRef');

defineShortcuts({
    '/': () => inputRef.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.settings.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />

                    <UButton
                        v-if="!streamerMode"
                        :label="$t('common.upload')"
                        trailing-icon="i-mdi-upload"
                        @click="fileUploadModal.open()"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm
                    ref="formRef"
                    class="my-2 flex w-full flex-1 flex-col gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="refresh()"
                >
                    <div class="flex flex-1 flex-row gap-2">
                        <UFormField class="flex-1" name="prefix" :label="$t('common.search')">
                            <UButtonGroup class="w-full">
                                <UInput
                                    ref="inputRef"
                                    v-model="query.prefix"
                                    type="text"
                                    name="prefix"
                                    :placeholder="$t('common.path')"
                                    class="w-full"
                                    leading-icon="i-mdi-search"
                                >
                                    <template #trailing>
                                        <UKbd value="/" />
                                    </template>
                                </UInput>
                                <UButton icon="i-mdi-subdirectory-arrow-left" variant="subtle" @click="goBackDirectory" />
                            </UButtonGroup>
                        </UFormField>
                    </div>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />
            <DataErrorBlock
                v-else-if="error"
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
                :empty="$t('common.not_found', [$t('common.file', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
                sticky
            />
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="files?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>

<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import FileUploadModal from '~/components/rector/filestore/FileUploadModal.vue';
import { useSettingsStore } from '~/store/settings';
import type { FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { DeleteFileResponse, ListFilesResponse } from '~~/gen/ts/services/rector/filestore';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const prefix = ref('');

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData('chars', () => listFiles(prefix.value));

async function listFiles(prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = getGRPCRectorFilestoreClient().listFiles({
            pagination: { offset: offset.value },
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
        const { response } = getGRPCRectorFilestoreClient().deleteFile({
            path,
        });

        const idx = data.value?.files.findIndex((f) => f.name === path);
        if (idx !== undefined && idx > -1 && data.value !== null) {
            data.value!.files = data.value?.files.splice(idx, 1) ?? [];
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

function addUploadedFile(file: FileInfo): void {
    const idx = data.value?.files.findIndex((f) => f.name === file.name);
    if (idx === undefined) {
        return;
    }

    if (idx > -1) {
        data.value?.files.unshift(file);
    } else {
        data.value?.files.splice(idx, 1, file);
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
        <UDashboardNavbar :title="$t('pages.rector.settings.title')">
            <template #right>
                <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                    {{ $t('common.back') }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="pb-24">
            <UDashboardSection
                :title="$t('system.streamer_mode.title')"
                :description="$t('system.streamer_mode.description')"
            />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <UDashboardNavbar :title="$t('pages.rector.filestore.title')">
            <template #right>
                <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                    {{ $t('common.back') }}
                </UButton>

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
            :title="$t('common.unable_to_load', [`${$t('common.data', 1)} ${$t('common.prop')}`])"
            :retry="refresh"
        />

        <UTable
            v-else
            :loading="loading"
            :columns="columns"
            :rows="data?.files"
            :empty-state="{ icon: 'i-mdi-file-multiple', label: $t('common.not_found', [$t('common.file', 2)]) }"
            class="flex-1"
        >
            <template #actions-data="{ row: file }">
                <UButton
                    variant="link"
                    icon="i-mdi-link-variant"
                    :external="true"
                    target="_blank"
                    :to="`/api/filestore/${file.name}`"
                />
                <UButton
                    variant="link"
                    icon="i-mdi-trash-can"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteFile(file.name),
                        })
                    "
                />
            </template>
            <template #name-data="{ row: file }">
                <span class="text-gray-900 dark:text-white">
                    {{ file.name }}
                </span>
            </template>
            <template #preview-data="{ row: file }">
                <span v-if="!previewTypes.some((ext) => file.name.endsWith(ext))"> </span>
                <img v-else :src="`/api/filestore/${file.name}`" class="max-h-24 max-w-32" />
            </template>
            <template #fileSize-data="{ row: file }">
                {{ formatBytes(file.size) }}
            </template>
            <template #updatedAt-data="{ row: file }">
                <GenericTime :value="toDate(file.lastModified)" />
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" infinite />
    </template>
</template>

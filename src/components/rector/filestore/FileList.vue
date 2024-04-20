<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useSettingsStore } from '~/store/settings';
import type { DeleteFileResponse, ListFilesResponse } from '~~/gen/ts/services/rector/filestore';
import FileUploadModal from '~/components/rector/filestore/FileUploadModal.vue';
import type { FileInfo } from '~~/gen/ts/resources/filestore/file';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const prefix = ref('');

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData('chars', () => listFiles(prefix.value));

async function listFiles(prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = $grpc.getRectorFilestoreClient().listFiles({
            pagination: { offset: offset.value },
            path: prefix,
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteFile(path: string): Promise<DeleteFileResponse> {
    try {
        const { response } = $grpc.getRectorFilestoreClient().deleteFile({
            path,
        });

        const idx = data.value?.files.findIndex((f) => f.name === path);
        if (idx !== undefined && idx > -1 && data.value !== null) {
            data.value.files = data.value.files.splice(idx, 1);
        }

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
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
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardPanelContent class="pb-24">
            <UDashboardSection
                :title="$t('system.streamer_mode.title')"
                :description="$t('system.streamer_mode.description')"
            />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <div class="sm:flex sm:items-center">
            <div class="w-full sm:flex-auto">
                <UButton
                    block
                    @click="
                        modal.open(FileUploadModal, {
                            onUploaded: addUploadedFile,
                        })
                    "
                >
                    {{ $t('common.upload') }}
                </UButton>
            </div>
        </div>

        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.data', 1)} ${$t('common.prop')}`])"
            :retry="refresh"
        />

        <template v-else>
            <UTable
                :loading="loading"
                :columns="columns"
                :rows="data?.files"
                :empty-state="{ icon: 'i-mdi-file-multiple', label: $t('common.not_found', [$t('common.file', 2)]) }"
            >
                <template #actions-data="{ row: file }">
                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
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
                <template #fileSize-data="{ row: file }">
                    {{ formatBytes(file.size) }}
                </template>
                <template #updatedAt-data="{ row: file }">
                    <GenericTime :value="toDate(file.lastModified)" />
                </template>
            </UTable>

            <Pagination v-model="page" :pagination="data?.pagination" infinite />
        </template>
    </template>
</template>

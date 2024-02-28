<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { FileMultipleIcon } from 'mdi-vue3';
import GenericContainerPanel from '~/components/partials/GenericContainerPanel.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useSettingsStore } from '~/store/settings';
import type { ListFilesResponse } from '~~/gen/ts/services/rector/rector';
import FileListEntry from '~/components/rector/filestore/FileListEntry.vue';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const offset = ref(0n);
const prefix = ref('');

const { data, pending, refresh, error } = useLazyAsyncData('chars', () => listFiles(prefix.value));

async function listFiles(prefix: string): Promise<ListFilesResponse> {
    try {
        const { response } = $grpc.getRectorClient().listFiles({
            pagination: { offset: offset.value },
            path: prefix,
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, () => refresh());
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <template v-if="streamerMode">
                <GenericContainerPanel>
                    <template #title>
                        {{ $t('system.streamer_mode.title') }}
                    </template>
                    <template #description>
                        {{ $t('system.streamer_mode.description') }}
                    </template>
                </GenericContainerPanel>
            </template>
            <template v-else>
                <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.data', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [`${$t('common.data', 1)} ${$t('common.prop')}`])"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="data === null" :icon="FileMultipleIcon" :type="$t('common.data', 1)" />

                <template v-else>
                    <table class="min-w-full divide-y divide-base-600">
                        <thead>
                            <tr>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.action', 2) }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.name') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.file_size') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.updated_at') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.type') }}
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-base-800">
                            <FileListEntry
                                v-for="(file, idx) in data.files"
                                :key="file.name"
                                :file="file"
                                @deleted="data.files.splice(idx, 1)"
                            />
                        </tbody>
                    </table>

                    <TablePagination :pagination="data.pagination" :refresh="refresh" @offset-change="offset = $event" />
                </template>
            </template>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog } from '@vueuse/core';
import { EyeIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { DeleteFileResponse } from '~~/gen/ts/services/rector/rector';

defineProps<{
    file: FileInfo;
}>();

const emits = defineEmits<{
    (e: 'deleted', path: string): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteFile(path: string): Promise<DeleteFileResponse> {
    try {
        const { response } = $grpc.getRectorClient().deleteFile({
            path,
        });

        emits('deleted', path);

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (path) => deleteFile(path));
</script>

<template>
    <tr :key="file.name">
        <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(file.name)" />

        <td class="whitespace-nowrap px-1 py-1 text-left text-sm text-accent-200 inline-flex items-center">
            <NuxtLink
                :external="true"
                target="_blank"
                :to="`/api/filestore/${file.name}`"
                class="text-primary-400 hover:text-primary-600"
            >
                <EyeIcon class="ml-auto mr-1.5 w-5 h-auto" aria-hidden="true" />
            </NuxtLink>
            <button type="button" class="text-primary-400 hover:text-primary-600" @click="reveal(file.name)">
                <TrashCanIcon class="w-5 h-5" aria-hidden="true" />
            </button>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            {{ file.name }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-sm text-accent-200">
            {{ formatBytesBigInt(file.size) }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-sm text-accent-200">
            <GenericTime :value="toDate(file.lastModified)" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-sm text-accent-200">
            {{ file.contentType }}
        </td>
    </tr>
</template>

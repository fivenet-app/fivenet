<script lang="ts" setup>
import { useConfirmDialog } from '@vueuse/core';
import { EyeIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { DeleteFileResponse } from '~~/gen/ts/services/rector/filestore';

defineProps<{
    file: FileInfo;
}>();

const emits = defineEmits<{
    (e: 'deleted', path: string): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteFile(path: string): Promise<DeleteFileResponse> {
    try {
        const { response } = $grpc.getRectorFilestoreClient().deleteFile({
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

        <td class="inline-flex items-center whitespace-nowrap p-1 text-left text-sm text-accent-200">
            <NuxtLink
                :external="true"
                target="_blank"
                :to="`/api/filestore/${file.name}`"
                class="text-primary-400 hover:text-primary-600"
            >
                <EyeIcon class="ml-auto mr-1.5 h-auto w-5" aria-hidden="true" />
            </NuxtLink>
            <button type="button" class="text-primary-400 hover:text-primary-600" @click="reveal(file.name)">
                <TrashCanIcon class="size-5" aria-hidden="true" />
            </button>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            {{ file.name }}
        </td>
        <td class="whitespace-nowrap p-1 text-left text-sm text-accent-200">
            {{ formatBytesBigInt(file.size) }}
        </td>
        <td class="whitespace-nowrap p-1 text-left text-sm text-accent-200">
            <GenericTime :value="toDate(file.lastModified)" />
        </td>
        <td class="whitespace-nowrap p-1 text-left text-sm text-accent-200">
            {{ file.contentType }}
        </td>
    </tr>
</template>

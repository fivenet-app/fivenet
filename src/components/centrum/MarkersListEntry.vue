<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog } from '@vueuse/core';
import { MapMarkerIcon, TrashCanIcon } from 'mdi-vue3';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/livemap';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useLivemapStore } from '~/store/livemap';

defineProps<{
    marker: MarkerMarker;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker } = livemapStore;

async function deleteMarker(id: string): Promise<void> {
    try {
        const call = $grpc.getLivemapperClient().deleteMarker({
            id,
        });
        await call;

        deleteMarkerMarker(id);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteMarker(id));
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(marker.info!.id)" />

    <tr :key="marker.info!.id" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td
            class="relative items-center justify-start whitespace-nowrap py-1 pl-0 pr-0 text-left text-sm font-medium sm:pr-0.5"
        >
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.go_to_location')"
                @click="$emit('goto', { x: marker.info!.x, y: marker.info!.y })"
            >
                <MapMarkerIcon class="ml-auto mr-1.5 w-5 h-auto" aria-hidden="true" />
            </button>
            <button
                v-if="can('LivemapperService.DeleteMarker')"
                type="button"
                :title="$t('common.delete')"
                class="inline-flex flex-row items-center text-primary-400 hover:text-primary-600"
                @click="reveal(marker.info!.id)"
            >
                <TrashCanIcon class="h-5 w-5" />
                <span class="sr-only">{{ $t('common.delete') }}</span>
            </button>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <GenericTime :value="marker.info!.createdAt" type="short" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <GenericTime v-if="marker.expiresAt" :value="marker.expiresAt" type="short" />
            <span v-else>
                {{ $t('common.na') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ marker.info!.name }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-100">
            {{ $t(`enums.livemap.MarkerType.${MarkerType[marker.type]}`) }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <p class="max-h-14 overflow-y-scroll break-words">
                {{ marker.info?.description ?? $t('common.na') }}
            </p>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="marker.creator">
                <CitizenInfoPopover :user="marker.creator" />
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ marker.creator?.jobLabel ?? $t('common.na') }}
        </td>
    </tr>
</template>

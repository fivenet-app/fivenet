<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LPopup } from '@vue-leaflet/vue-leaflet';
import { TrashCanIcon } from 'mdi-vue3';
import { useConfirmDialog } from '@vueuse/core';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

defineProps<{
    marker: Marker;
}>();

const { $grpc } = useNuxtApp();

async function deleteMarker(id: string): Promise<void> {
    try {
        const call = $grpc.getLivemapperClient().deleteMarker({
            id,
        });
        await call;
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

    <LPopup :options="{ closeButton: true }">
        <div v-if="can('LivemapperService.DeleteMarker')" class="mb-1 flex items-center gap-2">
            <button
                type="button"
                :title="$t('common.delete')"
                class="inline-flex items-center text-primary-500 hover:text-primary-400"
                @click="reveal(marker.info!.id)"
            >
                <TrashCanIcon class="h-5 w-5" />
                <span class="ml-1">{{ $t('common.delete') }}</span>
            </button>
        </div>
        <span class="font-semibold"> {{ $t('common.marker') }}: {{ marker.info?.name }} </span>
        <ul role="list" class="flex flex-col">
            <li>
                <span class="font-semibold">{{ $t('common.job') }}:</span>
                {{ marker.info?.jobLabel ?? $t('common.na') }}
            </li>
            <li>
                <span class="font-semibold">{{ $t('common.description') }}:</span>
                {{ marker.info?.description ?? $t('common.na') }}
            </li>
            <li class="inline-flex gap-1">
                <span class="flex-initial">
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                </span>
                <span class="flex-1">
                    <CitizenInfoPopover v-if="marker.creator" :user="marker.creator" />
                    <template v-else>
                        {{ $t('common.unknown') }}
                    </template>
                </span>
            </li>
        </ul>
    </LPopup>
</template>

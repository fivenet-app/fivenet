<script lang="ts" setup>
import { LPopup } from '@vue-leaflet/vue-leaflet';
import { MapMarkerIcon, TrashCanIcon } from 'mdi-vue3';
import { useConfirmDialog } from '@vueuse/core';
import { type MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
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

    <LPopup :options="{ closeButton: true }">
        <div class="mb-1 flex items-center gap-2">
            <button
                v-if="marker.info?.x && marker.info?.y"
                type="button"
                class="inline-flex items-center text-primary-500 hover:text-primary-400"
                @click="$emit('goto', { x: marker.info?.x, y: marker.info?.y })"
            >
                <MapMarkerIcon class="size-5" aria-hidden="true" />
                <span class="ml-1">{{ $t('common.mark') }}</span>
            </button>
            <button
                v-if="can('LivemapperService.DeleteMarker')"
                type="button"
                :title="$t('common.delete')"
                class="inline-flex items-center text-primary-500 hover:text-primary-400"
                @click="reveal(marker.info!.id)"
            >
                <TrashCanIcon class="size-5" aria-hidden="true" />
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
                <span class="font-semibold">{{ $t('common.expires_at') }}:</span>
                <GenericTime v-if="marker.expiresAt" :value="marker.expiresAt" />
                <span v-else>{{ $t('common.na') }}</span>
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

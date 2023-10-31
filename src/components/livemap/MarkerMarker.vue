<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LCircleMarker, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { useConfirmDialog } from '@vueuse/core';
import L from 'leaflet';
import { TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';

const props = withDefaults(
    defineProps<{
        marker: Marker;
        size?: number;
    }>(),
    {
        size: 20,
    },
);

defineEmits<{
    (e: 'selected'): void;
}>();

const iconAnchor: L.PointExpression = [props.size / 2, props.size];
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];
const icon = new L.DivIcon({
    html: `<div>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.8 16 17.6" fill="${
            props.marker.info?.color ? '#' + props.marker.info?.color : 'currentColor'
        }" class="w-full h-full">
                <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10zm0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6z"/>
            </svg>
        </div>`,
    iconSize: [props.size, props.size],
    iconAnchor,
    popupAnchor,
}) as L.Icon;

const { $grpc } = useNuxtApp();

async function deleteMarker(id: bigint): Promise<void> {
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

    <LCircleMarker
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id?.toString()"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius"
        :color="marker.info?.color ? '#' + marker.info?.color : '#fff'"
        :fill-opacity="(marker.data.data.circle.oapcity ?? 5) / 100"
    >
        <LPopup :options="{ closeButton: true }">
            <ul>
                <li>{{ marker.info?.name }}</li>
                <li>{{ $t('common.description') }}: {{ marker.info?.description }}</li>
            </ul>
            <template v-if="can('LivemapperService.DeleteMarker')">
                <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                    <TrashCanIcon class="w-6 h-6" />
                    <span>{{ $t('common.delete') }}</span>
                </button>
            </template>
        </LPopup>
    </LCircleMarker>

    <LMarker
        v-else
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        :icon="icon"
        @click="$emit('selected')"
    >
        <LPopup :options="{ closeButton: true }">
            <ul>
                <li>{{ marker.info?.name }}</li>
                <li v-if="marker.info?.description">{{ $t('common.description') }}: {{ marker.info?.description }}</li>
                <li class="italic">
                    <span class="font-semibold">{{ $t('common.sent_by') }}</span
                    >:
                    <span v-if="marker.creator">
                        {{ marker.creator?.firstname }}, {{ marker.creator?.lastname }} (<PhoneNumber
                            :number="marker.creator.phoneNumber"
                        />)
                    </span>
                    <span v-else>
                        {{ $t('common.unknown') }}
                    </span>
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

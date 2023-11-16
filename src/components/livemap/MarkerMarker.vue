<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LCircleMarker, LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { useConfirmDialog } from '@vueuse/core';
import { type PointExpression } from 'leaflet';
import { HelpIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';

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

const iconAnchor: PointExpression = [props.size / 2, props.size];
const popupAnchor: PointExpression = [0, (props.size / 2) * -1];

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

    <LCircleMarker
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius"
        :color="marker.info?.color ? '#' + marker.info?.color : '#fff'"
        :fill-opacity="(marker.data.data.circle.oapcity ?? 5) / 100"
    >
        <LPopup :options="{ closeButton: true }">
            <ul>
                <li class="inline-flex items-center">
                    {{ marker.info?.name }}
                    <template v-if="can('LivemapperService.DeleteMarker')">
                        <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                            <TrashCanIcon class="w-6 h-6" />
                            <span class="sr-only">{{ $t('common.delete') }}</span>
                        </button>
                    </template>
                </li>
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
    </LCircleMarker>

    <LMarker
        v-else-if="marker.data?.data.oneofKind === 'icon'"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :name="marker.info!.name"
        @click="$emit('selected')"
    >
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <component
                :is="
                    markerIcons.find((i) => marker.data?.data.oneofKind === 'icon' && i.name === marker.data?.data.icon.icon) ??
                    HelpIcon
                "
                class="w-6 h-6"
                :style="{ color: marker.info?.color ? '#' + marker.info?.color : 'currentColor' }"
            />
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <ul>
                <li class="inline-flex items-center">
                    {{ marker.info?.name }}
                    <template v-if="can('LivemapperService.DeleteMarker')">
                        <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                            <TrashCanIcon class="w-6 h-6" />
                            <span class="sr-only">{{ $t('common.delete') }}</span>
                        </button>
                    </template>
                </li>
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

    <LMarker v-else :lat-lng="[marker.info!.y, marker.info!.x]" :name="marker.info!.name" @click="$emit('selected')">
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <div>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 -0.8 16 17.6"
                    :fill="marker.info?.color ? '#' + marker.info?.color : 'currentColor'"
                    class="w-full h-full"
                >
                    <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10zm0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6z" />
                </svg>
            </div>
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <ul>
                <li class="inline-flex items-center">
                    {{ marker.info?.name }}
                    <template v-if="can('LivemapperService.DeleteMarker')">
                        <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                            <TrashCanIcon class="w-6 h-6" />
                            <span class="sr-only">{{ $t('common.delete') }}</span>
                        </button>
                    </template>
                </li>
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

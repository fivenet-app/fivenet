<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LCircle, LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { useConfirmDialog } from '@vueuse/core';
import { type PointExpression } from 'leaflet';
import { HelpIcon, MapMarkerQuestionIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Marker } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

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

    <LCircle
        v-if="marker.data?.data.oneofKind === 'circle'"
        :key="marker.info!.id"
        :lat-lng="[marker.info!.y, marker.info!.x]"
        :radius="marker.data?.data.circle.radius / 0.6931471805599453"
        :color="marker.info?.color ? '#' + marker.info?.color : '#ffffff'"
        :fill-opacity="(marker.data.data.circle.oapcity ?? 3) / 100"
    >
        <LPopup :options="{ closeButton: true }">
            <div v-if="can('LivemapperService.DeleteMarker')" class="mb-1 flex items-center gap-2">
                <button
                    type="button"
                    :title="$t('common.delete')"
                    class="inline-flex items-center text-primary-500 hover:text-primary-400"
                    @click="reveal()"
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
                <li>
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    <CitizenInfoPopover v-if="marker.creator" :user="marker.creator" />
                    <span v-else>
                        {{ $t('common.unknown') }}
                    </span>
                </li>
            </ul>
        </LPopup>
    </LCircle>

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
                class="h-6 w-6"
                :style="{ color: marker.info?.color ? '#' + marker.info?.color : 'currentColor' }"
            />
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <ul>
                <li class="inline-flex items-center">
                    <span class="font-semibold">
                        {{ marker.info?.name }}
                    </span>
                    <template v-if="can('LivemapperService.DeleteMarker')">
                        <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                            <TrashCanIcon class="h-5 w-5" />
                            <span class="sr-only">{{ $t('common.delete') }}</span>
                        </button>
                    </template>
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ marker.info?.description ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    <CitizenInfoPopover v-if="marker.creator" :user="marker.creator" />
                    <span v-else>
                        {{ $t('common.unknown') }}
                    </span>
                </li>
            </ul>
        </LPopup>
    </LMarker>

    <LMarker v-else :lat-lng="[marker.info!.y, marker.info!.x]" :name="marker.info!.name" @click="$emit('selected')">
        <LIcon :icon-size="[size, size]" :icon-anchor="iconAnchor" :popup-anchor="popupAnchor">
            <MapMarkerQuestionIcon :fill="marker.info?.color ? '#' + marker.info?.color : 'currentColor'" class="h-5 w-5" />
        </LIcon>

        <LPopup :options="{ closeButton: true }">
            <ul>
                <li class="inline-flex items-center">
                    <span class="font-semibold">
                        {{ marker.info?.name }}
                    </span>
                    <template v-if="can('LivemapperService.DeleteMarker')">
                        <button type="button" :title="$t('common.delete')" class="flex flex-row items-center" @click="reveal()">
                            <TrashCanIcon class="h-5 w-5" />
                            <span class="sr-only">{{ $t('common.delete') }}</span>
                        </button>
                    </template>
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ marker.info?.description ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    <CitizenInfoPopover v-if="marker.creator" :user="marker.creator" />
                    <span v-else>
                        {{ $t('common.unknown') }}
                    </span>
                </li>
            </ul>
        </LPopup>
    </LMarker>
</template>

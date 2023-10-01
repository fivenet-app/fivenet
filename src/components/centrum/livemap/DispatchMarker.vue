<script lang="ts" setup>
import { LIcon, LMarker, LPopup } from '@vue-leaflet/vue-leaflet';
import { BellIcon } from 'mdi-vue3';
import { dispatchStatusAnimate, dispatchStatusToBGColor, dispatchStatusToFillColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        size?: number;
    }>(),
    {
        size: 22,
    },
);

const emit = defineEmits<{
    (e: 'selected', dsp: Dispatch): void;
}>();

const iconAnchor: L.PointExpression | undefined = undefined;
const popupAnchor: L.PointExpression = [0, (props.size / 2) * -1];

function selected(_: bigint | string) {
    emit('selected', props.dispatch);
}
</script>

<template>
    <LMarker :key="dispatch.id?.toString()" :latLng="[dispatch.y, dispatch.x]" :name="dispatch.message" :z-index-offset="15">
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="uppercase flex flex-col items-center dsp-status-error">
                <span
                    class="rounded-md bg-white text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0"
                >
                    DSP-{{ props.dispatch.id }}
                </span>
                <BellIcon
                    class="w-full h-full"
                    :class="[
                        dispatchStatusAnimate(props.dispatch.status?.status ?? 0) ? 'animate-wiggle' : '',
                        dispatchStatusToFillColor(dispatch.status?.status ?? 0),
                    ]"
                />
            </div>
        </LIcon>
        <LPopup :options="{ closeButton: true }">
            <IDCopyBadge class="mb-1" prefix="DSP" :id="dispatch.id" :action="selected" />
            <ul>
                <li>
                    <span class="font-semibold">{{ $t('common.status') }}</span
                    >:
                    <span :class="dispatchStatusToBGColor(dispatch.status?.status ?? 0)">
                        {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                    </span>
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.message') }}</span
                    >: {{ dispatch.message }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.description') }}</span
                    >: {{ dispatch.description ?? 'N/A' }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.sent_at') }}</span
                    >: {{ $d(toDate(dispatch.createdAt), 'short') }}
                </li>
                <li class="italic">
                    <span class="font-semibold">{{ $t('common.sent_by') }}</span
                    >:
                    <span v-if="dispatch.anon">
                        {{ $t('common.anon') }}
                    </span>
                    <span v-else-if="dispatch.user">
                        {{ dispatch.user?.firstname }}, {{ dispatch.user?.lastname }} (<PhoneNumber
                            :number="dispatch.user.phoneNumber"
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

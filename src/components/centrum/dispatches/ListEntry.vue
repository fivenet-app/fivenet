<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
import { AccountMultiplePlusIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { dispatchStatusAnimate, dispatchStatusToBGColor } from '../helpers';
import AssignDispatchModal from './AssignDispatchModal.vue';
import Details from './Details.vue';

const props = defineProps<{
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status ?? 0));
const dispatchAnimated = computed(() => (dispatchStatusAnimate(props.dispatch.status?.status ?? 0) ? 'animate-pulse' : ''));

const dispatchAssistanceSound = useSound('/sounds/centrum/morse-sos.mp3', {
    volume: 0.15,
});
const debouncedPlay = useDebounceFn(() => dispatchAssistanceSound.play(), 950);

watch(props, () => {
    if (props.dispatch.status?.status === StatusDispatch.NEED_ASSISTANCE) {
        debouncedPlay();
    }
});

const openDetails = ref(false);
const openAssign = ref(false);

const openMessage = ref(false);
</script>

<template>
    <tr>
        <Details :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />
        <AssignDispatchModal v-if="openAssign" :open="openAssign" :dispatch="dispatch" @close="openAssign = false" />

        <td
            class="relative items-center whitespace-nowrap pl-0 py-1 pr-0 text-left text-sm font-medium sm:pr-0.5 justify-start"
        >
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="openAssign = true"
            >
                <AccountMultiplePlusIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                <MapMarkerIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.detail', 2)"
                @click="openDetails = true"
            >
                <DotsVerticalIcon class="w-6 h-auto ml-auto mr-1.5" />
            </button>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.id }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <Time :value="dispatch.createdAt" type="compact" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-100" :class="dispatchBackground">
            <span :class="dispatchAnimated">
                {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[props.dispatch.status?.status ?? 0]}`) }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.postal ?? $t('common.na') }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.units.length === 0" class="italic">{{ $t('enums.centrum.StatusDispatch.UNASSIGNED') }}</span>
            <span v-else class="mr-1">
                {{ dispatch.units.map((unit) => unit.unit?.initials ?? $t('common.na')).join(', ') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.anon">
                {{ $t('common.anon') }}
            </span>
            <span v-else-if="dispatch.creator">
                <CitizenInfoPopover :user="dispatch.creator" />
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="px-1 py-1 text-sm text-gray-300">
            <p class="break-all" :class="openMessage ? '' : 'max-h-24 max-w-sm'">
                {{
                    openMessage
                        ? dispatch.message
                        : dispatch.message.substring(0, 40) + (dispatch.message.length > 40 ? '...' : '')
                }}
            </p>
            <button
                v-if="dispatch.message.length > 40"
                type="button"
                @click="openMessage = !openMessage"
                class="flex justify-center px-1 py-1 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 bg-accent-500 hover:bg-accent-400 focus-visible:outline-accent-500"
            >
                {{ openMessage ? $t('common.read_less') : $t('common.read_more') }}
            </button>
        </td>
    </tr>
</template>

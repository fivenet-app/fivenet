<script lang="ts" setup>
import { unitStatusToBGColor } from '~/components/centrum/helpers';
import { RGB } from '~/utils/colour';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';
import Details from './Details.vue';

const props = defineProps<{
    unit: Unit;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const unitColorHex = hexToRgb('#' + props.unit.color ?? '000000') ?? ({ r: 0, g: 0, b: 0 } as RGB);

watch(
    () => props.unit,
    () => {
        console.log('LIST ENTRY UNIT:', UNIT_STATUS[props.unit.status?.status ?? 0], props.unit.users);
    },
);

const open = ref(false);
</script>

<template>
    <li class="col-span-1 flex rounded-md shadow-sm" @click="open = true">
        <Details :open="open" @close="open = false" :unit="unit" @goto="$emit('goto', $event)" />

        <div
            class="flex flex-shrink-0 items-center justify-center rounded-l-md text-sm font-medium border-l border-t border-b w-12"
            :class="isColourBright(unitColorHex) ? 'text-black' : 'text-white'"
            :style="'background-color: #' + unit.color ?? '000000'"
        >
            {{ unit.initials }}
        </div>
        <div class="flex flex-1 items-center justify-between truncate border border-gray-200 bg-gray">
            <div class="flex-1 px-1 py-2 text-sm">
                <span class="font-medium text-gray-100">{{ unit.name }}</span>
                <p class="text-gray-400">{{ $t('common.member', unit.users.length) }}</p>
            </div>
        </div>
        <div
            class="flex w-[5rem] flex-shrink-0 items-center justify-center rounded-r-md text-sm font-medium text-white border-r border-t border-b text-center"
            :class="unitStatusToBGColor(unit.status?.status ?? 0)"
        >
            {{ $t(`enums.centrum.UNIT_STATUS.${UNIT_STATUS[unit.status?.status ?? (0 as number)]}`) }}
        </div>
    </li>
</template>

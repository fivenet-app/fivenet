<script lang="ts" setup>
import { hexToRgb, isColorBright, rgbBlack, stringToColor } from '~/utils/color';

const props = withDefaults(
    defineProps<{
        state?: string;
        stateShort?: string;
        plate?: string;
        year?: string;
    }>(),
    {
        state: 'Los Santos',
        stateShort: 'LS',
        plate: 'UNKNOWN',
        year: '2018',
    },
);

const backgroundColor = stringToColor(props.plate);

const inverseColor = hexToRgb(backgroundColor, rgbBlack)!;
const year =
    props.year ??
    '201' +
        props.plate
            .charCodeAt(props.plate.length - 1)
            .toString()
            .charAt(0);
</script>

<template>
    <div class="flex max-w-48 min-w-20 flex-col items-center justify-center rounded-lg border" :style="{ backgroundColor }">
        <div class="grid w-full grid-cols-2 justify-items-center rounded-t-lg bg-error-600">
            <div class="hidden text-xs text-yellow-300 select-none md:block">{{ state }}</div>
            <div class="text-xs text-yellow-300 select-none md:hidden">{{ stateShort }}</div>
            <div class="text-xs text-white select-none">{{ year }}</div>
        </div>
        <div class="text-base" :class="isColorBright(inverseColor) ? 'text-black' : 'text-white'">
            {{ plate }}
        </div>
    </div>
</template>

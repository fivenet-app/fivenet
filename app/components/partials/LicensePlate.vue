<script lang="ts" setup>
import { RGBBlack, hexToRgb, isColourBright, stringToColour } from '~/utils/colour';

const props = withDefaults(
    defineProps<{
        state?: string;
        stateShort?: string;
        plate: string;
        year?: string;
    }>(),
    {
        state: 'Los Santos',
        stateShort: 'LS',
        plate: 'UNKNOWN',
        year: '2018',
    },
);

const backgroundColor = stringToColour(props.plate);

const inverseColor = hexToRgb(backgroundColor, RGBBlack)!;
const year =
    props.year ??
    '201' +
        props.plate
            .charCodeAt(props.plate.length - 1)
            .toString()
            .charAt(0);
</script>

<template>
    <div class="flex min-w-20 max-w-48 flex-col items-center justify-center rounded-lg border" :style="{ backgroundColor }">
        <div class="grid w-full grid-cols-2 justify-items-center rounded-t-lg bg-error-600">
            <div class="hidden select-none text-xs text-yellow-300 md:block">{{ state }}</div>
            <div class="select-none text-xs text-yellow-300 md:hidden">{{ stateShort }}</div>
            <div class="select-none text-xs text-white">{{ year }}</div>
        </div>
        <div class="text-base" :class="isColourBright(inverseColor) ? 'text-black' : 'text-white'">
            {{ plate }}
        </div>
    </div>
</template>

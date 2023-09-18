<script lang="ts" setup>
import { RGB, hexToRgb, isColourBright, stringToColour } from '~/utils/colour';

const props = withDefaults(
    defineProps<{
        state?: string;
        plate: string;
        year?: string;
    }>(),
    {
        state: 'Los Santos',
        plate: 'UNKNOWN',
    },
);

const backgroundColor = stringToColour(props.plate);

const inverseColor = hexToRgb(backgroundColor) ?? ({ r: 0, g: 0, b: 0 } as RGB);
const year = props.year ?? '201' + props.plate.charAt(props.plate.length - 1);
</script>

<template>
    <div
        class="flex flex-col items-center justify-center border-2 rounded-lg bg-[blue] max-w-[12rem]"
        :style="{ backgroundColor }"
    >
        <div class="w-full grid grid-cols-2 bg-error-600 justify-items-center rounded-t-lg">
            <div class="select-none text-xs text-warn-400">{{ state }}</div>
            <div class="select-none text-xs">{{ year }}</div>
        </div>
        <div class="text-xl" :class="isColourBright(inverseColor) ? 'text-black' : 'text-white'">
            {{ plate }}
        </div>
    </div>
</template>

<script lang="ts" setup>
import { type RGB, hexToRgb, isColourBright, stringToColour } from '~/utils/colour';

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
    },
);

const backgroundColor = stringToColour(props.plate);

const inverseColor = hexToRgb(backgroundColor) ?? ({ r: 0, g: 0, b: 0 } as RGB);
const year =
    props.year ??
    '201' +
        props.plate
            .charCodeAt(props.plate.length - 1)
            .toString()
            .charAt(0);
</script>

<template>
    <div
        class="flex min-w-[6.85rem] min-w-fit max-w-[12rem] flex-col items-center justify-center rounded-lg border-2 bg-[blue]"
        :style="{ backgroundColor }"
    >
        <div class="grid w-full grid-cols-2 justify-items-center rounded-t-lg bg-error-600">
            <div class="hidden select-none text-xs text-warn-400 md:block">{{ state }}</div>
            <div class="select-none text-xs text-warn-400 md:hidden">{{ stateShort }}</div>
            <div class="select-none text-xs">{{ year }}</div>
        </div>
        <div class="text-base" :class="isColourBright(inverseColor) ? 'text-black' : 'text-neutral'">
            {{ plate }}
        </div>
    </div>
</template>

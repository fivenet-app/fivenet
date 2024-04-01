<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { MinusThickIcon, PlusThickIcon } from 'mdi-vue3';
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <LControl position="topleft">
        <div class="inline-flex w-full flex-col">
            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                @click="if (zoom + 1 <= 7) zoom++;"
            >
                <PlusThickIcon class="size-5" aria-hidden="true" />
            </UButton>

            <UButton
                class="inset-0 inline-flex select-none items-center justify-center rounded-md border-2 border-black/20 bg-clip-padding"
            >
                {{ zoom }}
            </UButton>

            <UButton
                class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                @click="if (zoom - 1 >= 1) zoom--;"
            >
                <MinusThickIcon class="size-5" aria-hidden="true" />
            </UButton>
        </div>
    </LControl>
</template>

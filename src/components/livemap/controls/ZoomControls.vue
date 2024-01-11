<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { MinusThickIcon, PlusThickIcon } from 'mdi-vue3';
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <LControl position="topleft">
        <div class="w-full inline-flex flex-col">
            <button
                type="button"
                class="text-primary-500 p-1.5 rounded-md bg-neutral text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0 inline-flex items-center justify-center"
                @click="zoom++"
            >
                <PlusThickIcon class="h-5 w-5" aria-hidden="true" />
            </button>

            <button
                type="button"
                class="text-primary-500 p-1.5 rounded-md bg-neutral text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0 inline-flex items-center justify-center"
                @click="zoom--"
            >
                <MinusThickIcon class="h-5 w-5" aria-hidden="true" />
            </button>
        </div>
    </LControl>
</template>

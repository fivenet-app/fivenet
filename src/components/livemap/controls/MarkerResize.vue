<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';

const settingsStore = useSettingsStore();
const { livemapMarkerSize } = storeToRefs(settingsStore);
</script>

<template>
    <div class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center">
        <span class="text-lg mr-2 text-[#6f7683]">{{ livemapMarkerSize }}</span>
        <input
            name="livemapMarkerSize"
            type="range"
            class="h-1.5 w-full cursor-grab rounded-full my-auto"
            min="14"
            max="34"
            step="2"
            :value="livemapMarkerSize"
            @change="livemapMarkerSize = ($event.target as any).value"
        />
    </div>
</template>

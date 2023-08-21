<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <div class="p-2 bg-neutral border border-[#6b7280] flex flex-row justify-center">
        <span class="text-lg mr-2 text-[#6f7683]">{{ $t('components.livemap.center_selected_marker') }}</span>
        <input v-model="livemap.centerSelectedMarker" class="my-auto" name="centerSelectedMarker" type="checkbox" />
    </div>
</template>

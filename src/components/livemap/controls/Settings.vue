<script lang="ts" setup>
import { CogIcon } from 'mdi-vue3';
import SettingsModal from './SettingsModal.vue';

const open = ref(false);
</script>

<template>
    <SettingsModal :open="open" @close="open = false" />

    <button
        type="button"
        class="w-30 h-30 rounded-md bg-neutral text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0"
        @click="open = true"
        :title="$t('common.setting', 2)"
    >
        <span class="sr-only">{{ $t('common.setting', 2) }}</span>
        <CogIcon class="h-6 w-6" aria-hidden="true" />
    </button>
</template>

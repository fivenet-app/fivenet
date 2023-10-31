<script lang="ts" setup>
import { CogIcon } from 'mdi-vue3';
import SettingsModal from '~/components/livemap/controls/SettingsModal.vue';

const open = ref(false);
</script>

<template>
    <SettingsModal :open="open" @close="open = false" />

    <button
        type="button"
        class="rounded-md bg-neutral text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0 inline-flex items-center justify-center"
        :title="$t('common.setting', 2)"
        @click="open = true"
    >
        <CogIcon class="h-6 w-6" aria-hidden="true" />
        <span class="pr-0.5">{{ $t('common.setting', 2) }}</span>
    </button>
</template>

<script lang="ts" setup>
import { CogIcon } from 'mdi-vue3';
import SettingsModal from '~/components/livemap/controls/SettingsModal.vue';

const open = ref(false);
</script>

<template>
    <SettingsModal :open="open" @close="open = false" />

    <button
        type="button"
        class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-neutral bg-clip-padding text-black hover:bg-[#f4f4f4] focus:outline-none"
        :title="$t('common.setting', 2)"
        @click="open = true"
    >
        <CogIcon class="size-5" aria-hidden="true" />
        <span class="pr-0.5">{{ $t('common.setting', 2) }}</span>
    </button>
</template>

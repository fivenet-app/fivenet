<script lang="ts" setup>
import { ref } from 'vue';
import ClipboardModal from './ClipboardModal.vue';
import { ClipboardDocumentListIcon } from '@heroicons/vue/24/solid';

const open = ref(false);
</script>

<template>
    <ClipboardModal :open="open" @close="open = false" />
    <button type="button" :title="$t('common.clipboard')"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-10 right-8 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
        @click="open = true">
        <ClipboardDocumentListIcon class="w-10 h-auto" />
    </button>
</template>

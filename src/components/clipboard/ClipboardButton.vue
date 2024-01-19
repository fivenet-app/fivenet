<script lang="ts" setup>
import { ClipboardListIcon } from 'mdi-vue3';
import { ref } from 'vue';
import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';

const open = ref(false);
</script>

<template>
    <ClipboardModal :open="open" @close="open = false" />

    <button
        type="button"
        :title="$t('common.clipboard')"
        class="fixed bottom-2 right-6 z-10 flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-neutral shadow-float hover:bg-primary-400"
        @click="open = true"
    >
        <ClipboardListIcon class="h-auto w-10" />
    </button>
</template>

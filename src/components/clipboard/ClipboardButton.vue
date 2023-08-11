<script lang="ts" setup>
import { ClipboardListIcon } from 'mdi-vue3';
import { ref } from 'vue';
import Modal from '~/components/clipboard/modal/Modal.vue';

const open = ref(false);
</script>

<template>
    <Modal :open="open" @close="open = false" />
    <button
        type="button"
        :title="$t('common.clipboard')"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-2 right-6 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
        @click="open = true"
    >
        <ClipboardListIcon class="w-10 h-auto" />
    </button>
</template>

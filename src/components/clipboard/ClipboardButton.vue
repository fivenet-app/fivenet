<script lang="ts" setup>
import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';

const modal = useModal();
</script>

<template>
    <UButton
        :title="$t('common.clipboard')"
        icon="i-mdi-clipboard-list"
        :ui="{ rounded: 'rounded-full' }"
        class="fixed bottom-2 right-6 z-10 inline-flex size-12 items-center justify-center shadow-float"
        @click="modal.open(ClipboardModal, {})"
    >
    </UButton>
</template>

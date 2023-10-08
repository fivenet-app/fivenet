<script lang="ts" setup>
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import List from './List.vue';

defineProps<{
    open: boolean;
}>();
defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <ComponentModal :open="open" @close="$emit('close')" dialog-class="!max-w-2xl">
        <List />
    </ComponentModal>
</template>

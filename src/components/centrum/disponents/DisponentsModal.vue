<script lang="ts" setup>
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import DisponentsList from '~/components/centrum/disponents/DisponentsList.vue';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <ComponentModal :open="open" dialog-class="!max-w-2xl" @close="$emit('close')">
        <DisponentsList />
    </ComponentModal>
</template>

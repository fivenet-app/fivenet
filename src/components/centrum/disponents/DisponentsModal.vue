<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import DisponentsList from '~/components/centrum/disponents/DisponentsList.vue';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <GenericModal :open="open" dialog-class="!max-w-2xl" @close="$emit('close')">
        <DisponentsList />
    </GenericModal>
</template>

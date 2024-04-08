<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/quickbuttons/bodycheckup/BodyCheckup.vue';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <GenericModal :open="open" :title="$t('components.bodycheckup.title')" @close="$emit('close')">
        <BodyCheckup />
    </GenericModal>
</template>

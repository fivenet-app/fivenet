<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/quickbuttons/bodycheckup/BodyCheckup.vue';

const { isOpen } = useModal();
</script>

<template>
    <GenericModal :open="isOpen" :title="$t('components.bodycheckup.title')" @close="isOpen = false">
        <BodyCheckup />
    </GenericModal>
</template>

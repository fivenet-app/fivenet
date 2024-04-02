<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/quickbuttons/bodycheckup/BodyCheckup.vue';

const open = ref(true);
</script>

<template>
    <GenericModal :open="open" :unmount="false" :title="$t('components.bodycheckup.title')" @close="open = false">
        <BodyCheckup />
    </GenericModal>
</template>

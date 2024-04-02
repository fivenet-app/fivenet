<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/quickbuttons/bodycheckup/BodyCheckup.vue';

const open = ref(false);
</script>

<template>
    <UButton
        class="flex size-12 items-center justify-center rounded-full shadow-float"
        icon="i-mdi-human"
        @click="open = true"
    />

    <GenericModal :open="open" :unmount="false" :title="$t('components.bodycheckup.title')" @close="open = false">
        <BodyCheckup />
    </GenericModal>
</template>

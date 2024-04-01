<script lang="ts" setup>
import { HumanIcon } from 'mdi-vue3';
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/quickbuttons/bodycheckup/BodyCheckup.vue';

const open = ref(false);
</script>

<template>
    <UButton
        class="flex size-12 items-center justify-center rounded-full bg-accent-500 text-neutral shadow-float hover:bg-accent-400"
        @click="open = true"
    >
        <HumanIcon class="h-auto w-10" aria-hidden="true" />
    </UButton>

    <GenericModal :open="open" :unmount="false" :title="$t('components.bodycheckup.title')" @close="open = false">
        <BodyCheckup />
    </GenericModal>
</template>

<script lang="ts" setup>
import { HumanIcon } from 'mdi-vue3';
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/bodycheckup/BodyCheckup.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-neutral shadow-float hover:bg-primary-400"
        @click="open = true"
    >
        <HumanIcon class="h-auto w-10" />
    </button>

    <GenericModal :open="open" @close="open = false">
        <BodyCheckup />
    </GenericModal>
</template>

<script lang="ts" setup>
import { HumanIcon } from 'mdi-vue3';
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import BodyCheckup from '~/components/bodycheckup/BodyCheckup.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex h-12 w-12 items-center justify-center rounded-full bg-accent-500 text-neutral shadow-float hover:bg-accent-400"
        @click="open = true"
    >
        <HumanIcon class="w-10 h-auto" />
    </button>

    <GenericModal :open="open" :unmount="false" @close="open = false">
        <BodyCheckup />
    </GenericModal>
</template>

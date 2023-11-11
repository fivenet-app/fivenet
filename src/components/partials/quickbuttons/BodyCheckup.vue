<script lang="ts" setup>
import { HumanIcon } from 'mdi-vue3';
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import BodyCheckup from '~/components/bodycheckup/BodyCheckup.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex items-center justify-center w-12 h-12 rounded-full bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
        @click="open = true"
    >
        <HumanIcon class="w-10 h-auto" />
    </button>

    <ComponentModal :open="open" @close="open = false">
        <BodyCheckup />
    </ComponentModal>
</template>

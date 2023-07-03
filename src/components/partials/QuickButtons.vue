<script lang="ts" setup>
import QuickButton, { ButtonComponents as QuickButtonTypes } from '~/components/partials/QuickButton.vue';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <QuickButton v-for="button in jobProps?.quickButtons" :comp="(button as QuickButtonTypes)" />
</template>

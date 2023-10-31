<script lang="ts" setup>
import QuickButton, { type ButtonComponents as QuickButtonTypes } from '~/components/partials/QuickButton.vue';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <QuickButton v-for="button in jobProps?.quickButtons" :key="button" :comp="button as QuickButtonTypes" />
</template>

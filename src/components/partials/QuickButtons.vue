<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import QuickButton, { ButtonComponents as QuickButtonTypes } from './QuickButton.vue';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <QuickButton v-for="prop in jobProps?.quickButtons" :comp="(prop as QuickButtonTypes)" />
</template>

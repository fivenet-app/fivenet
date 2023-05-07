<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import ComponentButton, { ButtonComponents } from './ComponentButton.vue';

const authStore = useAuthStore();

const jobProps = computed(() => authStore.getJobProps);
</script>

<template>
    <ComponentButton v-for="prop in jobProps?.componentButtons" :comp="(prop as ButtonComponents)" />
</template>

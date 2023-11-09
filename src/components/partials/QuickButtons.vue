<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import PenaltyCalculator from '~/components/partials/quickbuttons/PenaltyCalculator.vue';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <PenaltyCalculator v-if="jobProps?.quickButtons?.penaltyCalculator" />
</template>

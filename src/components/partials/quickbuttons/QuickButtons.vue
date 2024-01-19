<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import PenaltyCalculatorButton from '~/components/partials/quickbuttons/PenaltyCalculatorButton.vue';
import BodyCheckupButton from '~/components/partials/quickbuttons/BodyCheckupButton.vue';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <div class="fixed bottom-2 right-20 z-10 flex flex-row-reverse items-center justify-center gap-2">
        <PenaltyCalculatorButton v-if="jobProps?.quickButtons?.penaltyCalculator" />
        <BodyCheckupButton v-if="jobProps?.quickButtons?.bodyCheckup" />
    </div>
</template>

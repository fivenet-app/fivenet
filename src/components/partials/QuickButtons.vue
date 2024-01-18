<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import PenaltyCalculator from '~/components/partials/quickbuttons/PenaltyCalculator.vue';
import BodyCheckup from '~/components/partials/quickbuttons/BodyCheckup.vue';

const authStore = useAuthStore();

const { jobProps } = storeToRefs(authStore);
</script>

<template>
    <div class="fixed bottom-2 right-20 z-10 flex flex-row-reverse items-center justify-center gap-2">
        <PenaltyCalculator v-if="jobProps?.quickButtons?.penaltyCalculator" />
        <BodyCheckup v-if="jobProps?.quickButtons?.bodyCheckup" />
    </div>
</template>

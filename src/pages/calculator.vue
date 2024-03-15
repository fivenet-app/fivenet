<script lang="ts" setup>
import Calculator from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

useHead({
    title: 'components.penaltycalculator.title',
});
definePageMeta({
    title: 'components.penaltycalculator.title',
    requiresAuth: true,
});
</script>

<template>
    <div class="w-full">
        <Calculator />
    </div>
</template>

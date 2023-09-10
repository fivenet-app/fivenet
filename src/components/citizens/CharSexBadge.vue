<script lang="ts" setup>
import { GenderFemaleIcon, GenderMaleIcon } from 'mdi-vue3';

defineProps<{
    sex: string;
}>();
</script>

<template>
    <span
        v-if="sex === 'f'"
        class="bg-pink-200 text-pink-600 inline-flex items-center rounded-full px-2.5 py-0.5 text-sm font-medium"
    >
        <GenderFemaleIcon class="h-5 w-5" />
    </span>

    <span
        v-else-if="sex === 'm'"
        class="bg-blue-200 text-blue-600 inline-flex items-center rounded-full px-2 py-0.25 text-sm font-medium"
    >
        <GenderMaleIcon class="h-5 w-5" />
    </span>
</template>
